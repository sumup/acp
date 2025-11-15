package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sumup/acp"
)

func main() {
	service := newMemoryService("USD", defaultCatalog())
	addr := ":8080"

	opts := []acp.Option{acp.WithMiddleware(logging, cors)}
	webhookOpts, err := webhookOptionsFromEnv()
	if err != nil {
		log.Fatalf("invalid webhook configuration: %v", err)
	}
	if webhookOpts != nil {
		opts = append(opts, acp.WithWebhookOptions(*webhookOpts))
	}

	log.Printf("ACP sample server listening on %s", addr)
	log.Printf("Try: curl -XPOST %s/checkout_sessions -d @- <<'JSON' ...", "http://localhost:8080")

	handler := acp.NewCheckoutHandler(service, opts...)
	if webhookOpts != nil {
		service.enableWebhooks(handler.SendWebhook)
		log.Printf("Webhooks enabled; delivering to %s", webhookOpts.Endpoint)
	} else {
		log.Printf("Webhooks disabled; set ACP_WEBHOOK_ENDPOINT, ACP_WEBHOOK_HEADER, and ACP_WEBHOOK_SECRET to enable delivery")
	}
	log.Fatal(http.ListenAndServe(addr, handler))
}

func webhookOptionsFromEnv() (*acp.WebhookOptions, error) {
	endpoint := strings.TrimSpace(os.Getenv("ACP_WEBHOOK_ENDPOINT"))
	header := strings.TrimSpace(os.Getenv("ACP_WEBHOOK_HEADER"))
	secret := os.Getenv("ACP_WEBHOOK_SECRET")

	if endpoint == "" && header == "" && secret == "" {
		return nil, nil
	}
	if endpoint == "" || header == "" || secret == "" {
		return nil, fmt.Errorf("ACP_WEBHOOK_* variables must all be set to enable webhook delivery")
	}
	return &acp.WebhookOptions{
		Endpoint:   endpoint,
		HeaderName: header,
		SecretKey:  []byte(secret),
	}, nil
}

// logging adds basic request logs without external dependencies.
func logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start).Truncate(time.Millisecond))
	}
}

// cors allows the browser-based testbed to call the sample server directly.
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept,API-Version")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code before forwarding to the real writer.
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

type product struct {
	SKU     string
	Title   string
	Price   int
	TaxRate float64
}

func defaultCatalog() []product {
	return []product{
		{SKU: "latte", Title: "Oat Milk Latte", Price: 650, TaxRate: 0.07},
		{SKU: "beans", Title: "Espresso Beans (1kg)", Price: 2400, TaxRate: 0.00},
		{SKU: "mug", Title: "Stoneware Mug", Price: 1500, TaxRate: 0.07},
	}
}

type sessionState struct {
	session *acp.CheckoutSession
	order   *acp.Order
}

type webhookSender func(context.Context, acp.EventData) error

type memoryService struct {
	mu        sync.RWMutex
	currency  string
	catalog   map[string]product
	sessions  map[string]*sessionState
	sessionID uint64
	orderID   uint64
	webhook   webhookSender
}

func newMemoryService(currency string, catalog []product) *memoryService {
	index := make(map[string]product, len(catalog))
	for _, p := range catalog {
		index[p.SKU] = p
	}
	return &memoryService{
		currency: strings.ToUpper(currency),
		catalog:  index,
		sessions: make(map[string]*sessionState),
	}
}

// CreateSession builds a new checkout session with default data.
func (s *memoryService) CreateSession(ctx context.Context, req acp.CheckoutSessionCreateRequest) (*acp.CheckoutSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &acp.CheckoutSession{
		ID:                 s.nextSessionID(),
		Currency:           s.currency,
		Buyer:              cloneBuyer(req.Buyer),
		FulfillmentAddress: cloneAddress(req.FulfillmentAddress),
		FulfillmentOptions: defaultFulfillmentOptions(),
		Messages:           defaultMessages(),
		Links: []acp.Link{
			{Type: acp.PrivacyPolicy, Url: "https://merchant.example/privacy"},
			{Type: acp.TermsOfUse, Url: "https://merchant.example/terms"},
		},
		PaymentProvider: &acp.PaymentProvider{
			Provider:                "sumup",
			SupportedPaymentMethods: []acp.SupportedPaymentMethods{acp.Card},
		},
	}

	if err := s.rebuildFinancials(session, req.Items); err != nil {
		return nil, err
	}
	session.Status = deriveStatus(session)

	state := &sessionState{session: session}
	s.sessions[session.ID] = state
	return cloneSession(session), nil
}

// UpdateSession refreshes an existing session with the provided fields.
func (s *memoryService) UpdateSession(ctx context.Context, id string, req acp.CheckoutSessionUpdateRequest) (*acp.CheckoutSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.sessions[id]
	if !ok {
		return nil, acp.NewHTTPError(http.StatusNotFound, acp.InvalidRequest, acp.ErrorCode("not_found"), "checkout session not found")
	}

	session := state.session
	if req.Buyer != nil {
		session.Buyer = cloneBuyer(req.Buyer)
	}
	if req.FulfillmentAddress != nil {
		session.FulfillmentAddress = cloneAddress(req.FulfillmentAddress)
	}
	if req.FulfillmentOptionId != nil {
		session.FulfillmentOptionId = req.FulfillmentOptionId
	}
	if req.Items != nil {
		if err := s.rebuildFinancials(session, *req.Items); err != nil {
			return nil, err
		}
	}
	session.Status = deriveStatus(session)
	return cloneSession(session), nil
}

// GetSession returns the current copy of a stored session.
func (s *memoryService) GetSession(ctx context.Context, id string) (*acp.CheckoutSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.sessions[id]
	if !ok {
		return nil, acp.NewHTTPError(http.StatusNotFound, acp.InvalidRequest, acp.ErrorCode("not_found"), "checkout session not found")
	}
	return cloneSession(state.session), nil
}

// CompleteSession marks a session as completed and emits a mock order.
func (s *memoryService) CompleteSession(ctx context.Context, id string, req acp.CheckoutSessionCompleteRequest) (*acp.SessionWithOrder, error) {
	s.mu.Lock()
	var (
		event   acp.EventData
		webhook = s.webhook
	)
	defer func() {
		s.mu.Unlock()
		if event == nil || webhook == nil {
			return
		}
		if err := webhook(ctx, event); err != nil {
			log.Printf("webhook delivery failed: %v", err)
		}
	}()

	state, ok := s.sessions[id]
	if !ok {
		return nil, acp.NewHTTPError(http.StatusNotFound, acp.InvalidRequest, acp.ErrorCode("not_found"), "checkout session not found")
	}
	session := state.session
	if session.Status == acp.CheckoutSessionStatusCanceled {
		return nil, acp.NewHTTPError(http.StatusConflict, acp.InvalidRequest, acp.ErrorCode("canceled"), "cannot complete a canceled session")
	}
	if len(session.LineItems) == 0 {
		return nil, acp.NewHTTPError(http.StatusBadRequest, acp.InvalidRequest, acp.ErrorCode("empty_cart"), "add items before completing the session")
	}
	if state.order != nil {
		return state.toOrderSession(), nil
	}

	session.Status = acp.CheckoutSessionStatusCompleted
	order := &acp.Order{
		ID:                s.nextOrderID(),
		CheckoutSessionId: session.ID,
		PermalinkUrl:      fmt.Sprintf("https://merchant.example/orders/%s", session.ID),
	}
	state.order = order

	event = acp.OrderCreate{
		Type:              acp.EventDataTypeOrder,
		CheckoutSessionID: session.ID,
		PermalinkURL:      order.PermalinkUrl,
		Status:            acp.OrderStatusCreated,
	}

	return state.toOrderSession(), nil
}

// CancelSession marks a session as canceled unless it already has an order.
func (s *memoryService) CancelSession(ctx context.Context, id string) (*acp.CheckoutSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.sessions[id]
	if !ok {
		return nil, acp.NewHTTPError(http.StatusNotFound, acp.InvalidRequest, acp.ErrorCode("not_found"), "checkout session not found")
	}
	if state.order != nil {
		return nil, acp.NewHTTPError(http.StatusConflict, acp.InvalidRequest, acp.ErrorCode("completed"), "completed sessions cannot be canceled")
	}

	state.session.Status = acp.CheckoutSessionStatusCanceled
	return cloneSession(state.session), nil
}

func (s *memoryService) rebuildFinancials(session *acp.CheckoutSession, items []acp.Item) error {
	lines, err := s.buildLineItems(items)
	if err != nil {
		return err
	}
	session.LineItems = lines
	session.Totals = buildTotals(lines, session.Currency)
	session.Messages = defaultMessages()
	return nil
}

func (s *memoryService) buildLineItems(items []acp.Item) ([]acp.LineItem, error) {
	if len(items) == 0 {
		return nil, acp.NewHTTPError(http.StatusBadRequest, acp.InvalidRequest, acp.ErrorCode(acp.InvalidRequest), "items cannot be empty")
	}

	lines := make([]acp.LineItem, 0, len(items))
	for idx, item := range items {
		product, ok := s.catalog[item.ID]
		if !ok {
			return nil, acp.NewHTTPError(http.StatusBadRequest, acp.InvalidRequest, acp.ErrorCode("unknown_item"), fmt.Sprintf("items[%d]: %q is not sold by this merchant", idx, item.ID))
		}
		base := product.Price * item.Quantity
		discount := 0
		tax := int(math.Round(product.TaxRate * float64(base)))
		subtotal := base - discount
		total := subtotal + tax
		lines = append(lines, acp.LineItem{
			ID:         fmt.Sprintf("li_%s_%d", item.ID, idx),
			Item:       item,
			BaseAmount: base,
			Discount:   discount,
			Subtotal:   subtotal,
			Tax:        tax,
			Total:      total,
		})
	}
	return lines, nil
}

func (s *memoryService) nextSessionID() string {
	id := atomic.AddUint64(&s.sessionID, 1)
	return fmt.Sprintf("cs_%06d", id)
}

func (s *memoryService) nextOrderID() string {
	id := atomic.AddUint64(&s.orderID, 1)
	return fmt.Sprintf("ord_%06d", id)
}

func deriveStatus(session *acp.CheckoutSession) acp.CheckoutSessionStatus {
	switch {
	case session.Status == acp.CheckoutSessionStatusCanceled:
		return acp.CheckoutSessionStatusCanceled
	case session.Status == acp.CheckoutSessionStatusCompleted:
		return acp.CheckoutSessionStatusCompleted
	case len(session.LineItems) == 0:
		return acp.CheckoutSessionStatusInProgress
	case session.PaymentProvider != nil:
		return acp.CheckoutSessionStatusReadyForPayment
	default:
		return acp.CheckoutSessionStatusNotReadyForPayment
	}
}

func buildTotals(lines []acp.LineItem, currency string) []acp.Total {
	var (
		itemsBase int
		tax       int
		total     int
	)
	for _, line := range lines {
		itemsBase += line.BaseAmount
		tax += line.Tax
		total += line.Total
	}

	totals := []acp.Total{
		{Type: acp.TotalTypeItemsBaseAmount, Amount: itemsBase, DisplayText: formatMoney(currency, itemsBase)},
	}
	if tax > 0 {
		totals = append(totals, acp.Total{
			Type:        acp.TotalTypeTax,
			Amount:      tax,
			DisplayText: formatMoney(currency, tax),
		})
	}
	totals = append(totals, acp.Total{
		Type:        acp.TotalTypeTotal,
		Amount:      total,
		DisplayText: formatMoney(currency, total),
	})
	return totals
}

func formatMoney(currency string, cents int) string {
	value := float64(cents) / 100
	return fmt.Sprintf("%s %.2f", currency, value)
}

func defaultMessages() []acp.Message {
	info := acp.MessageInfo{
		Type:        "info",
		Content:     "This sample server keeps sessions in memory. Restarting the process wipes the cart.",
		ContentType: acp.MessageInfoContentTypePlain,
	}
	var msg acp.Message
	_ = msg.FromMessageInfo(info)
	return []acp.Message{msg}
}

func defaultFulfillmentOptions() []acp.FulfillmentOption {
	soon := time.Now().Add(48 * time.Hour)
	later := soon.Add(24 * time.Hour)
	shipping := acp.FulfillmentOptionShipping{
		ID:                   "ship_standard",
		Title:                "Standard Shipping",
		Subtitle:             strPtr("2-4 business days"),
		Subtotal:             formatMoney("USD", 500),
		Tax:                  formatMoney("USD", 0),
		Total:                formatMoney("USD", 500),
		Type:                 "shipping",
		EarliestDeliveryTime: &soon,
		LatestDeliveryTime:   &later,
	}
	digital := acp.FulfillmentOptionDigital{
		ID:       "pickup",
		Title:    "In-store pickup",
		Subtitle: strPtr("Collect in person"),
		Subtotal: formatMoney("USD", 0),
		Tax:      formatMoney("USD", 0),
		Total:    formatMoney("USD", 0),
		Type:     "digital",
	}

	opts := make([]acp.FulfillmentOption, 0, 2)
	var shippingUnion acp.FulfillmentOption
	_ = shippingUnion.FromFulfillmentOptionShipping(shipping)
	opts = append(opts, shippingUnion)

	var digitalUnion acp.FulfillmentOption
	_ = digitalUnion.FromFulfillmentOptionDigital(digital)
	opts = append(opts, digitalUnion)

	return opts
}

func strPtr(v string) *string {
	return &v
}

func cloneBuyer(b *acp.Buyer) *acp.Buyer {
	if b == nil {
		return nil
	}
	copy := *b
	return &copy
}

func cloneAddress(a *acp.Address) *acp.Address {
	if a == nil {
		return nil
	}
	copy := *a
	return &copy
}

func clonePaymentProvider(p *acp.PaymentProvider) *acp.PaymentProvider {
	if p == nil {
		return nil
	}
	copy := *p
	if p.SupportedPaymentMethods != nil {
		copy.SupportedPaymentMethods = append([]acp.SupportedPaymentMethods(nil), p.SupportedPaymentMethods...)
	}
	return &copy
}

func cloneLineItems(src []acp.LineItem) []acp.LineItem {
	if len(src) == 0 {
		return nil
	}
	dst := make([]acp.LineItem, len(src))
	copy(dst, src)
	return dst
}

func cloneTotals(src []acp.Total) []acp.Total {
	if len(src) == 0 {
		return nil
	}
	dst := make([]acp.Total, len(src))
	copy(dst, src)
	return dst
}

func cloneLinks(src []acp.Link) []acp.Link {
	if len(src) == 0 {
		return nil
	}
	dst := make([]acp.Link, len(src))
	copy(dst, src)
	return dst
}

func cloneSession(src *acp.CheckoutSession) *acp.CheckoutSession {
	if src == nil {
		return nil
	}
	dst := *src
	dst.Buyer = cloneBuyer(src.Buyer)
	dst.FulfillmentAddress = cloneAddress(src.FulfillmentAddress)
	dst.PaymentProvider = clonePaymentProvider(src.PaymentProvider)
	dst.LineItems = cloneLineItems(src.LineItems)
	dst.Totals = cloneTotals(src.Totals)
	dst.Links = cloneLinks(src.Links)
	dst.FulfillmentOptions = append([]acp.FulfillmentOption(nil), src.FulfillmentOptions...)
	dst.Messages = append([]acp.Message(nil), src.Messages...)
	return &dst
}

func convertFulfillmentOptions(src []acp.FulfillmentOption) []acp.FulfillmentOption {
	if len(src) == 0 {
		return nil
	}
	dst := make([]acp.FulfillmentOption, len(src))
	for i := range src {
		data, err := src[i].MarshalJSON()
		if err != nil {
			continue
		}
		_ = dst[i].UnmarshalJSON(data)
	}
	return dst
}

func convertMessages(src []acp.Message) []acp.Message {
	if len(src) == 0 {
		return nil
	}
	dst := make([]acp.Message, len(src))
	for i := range src {
		data, err := src[i].MarshalJSON()
		if err != nil {
			continue
		}
		_ = dst[i].UnmarshalJSON(data)
	}
	return dst
}

func (s *sessionState) toOrderSession() *acp.SessionWithOrder {
	order := &acp.SessionWithOrder{
		CheckoutSession: acp.CheckoutSession{
			ID:                  s.session.ID,
			Buyer:               cloneBuyer(s.session.Buyer),
			Currency:            s.session.Currency,
			FulfillmentAddress:  cloneAddress(s.session.FulfillmentAddress),
			FulfillmentOptionId: s.session.FulfillmentOptionId,
			FulfillmentOptions:  convertFulfillmentOptions(s.session.FulfillmentOptions),
			LineItems:           cloneLineItems(s.session.LineItems),
			Links:               cloneLinks(s.session.Links),
			Messages:            convertMessages(s.session.Messages),
			PaymentProvider:     clonePaymentProvider(s.session.PaymentProvider),
			Status:              acp.CheckoutSessionStatusCompleted,
			Totals:              cloneTotals(s.session.Totals),
		},
		Order: *s.order,
	}
	return order
}

func (s *memoryService) enableWebhooks(sender webhookSender) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhook = sender
}
