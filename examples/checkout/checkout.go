package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpcheckout"
)

type memoryService struct {
	mu       sync.Mutex
	sessions map[string]*acpcheckout.CheckoutSessionBase
}

func newMemoryService() *memoryService {
	return &memoryService{sessions: make(map[string]*acpcheckout.CheckoutSessionBase)}
}

func (s *memoryService) CreateSession(_ context.Context, req acpcheckout.CheckoutSessionCreateRequest) (*acpcheckout.CheckoutSessionBase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lineItems, totals := priceItems(req.LineItems)
	session := &acpcheckout.CheckoutSessionBase{
		Id:                 "cs_demo_1",
		Buyer:              req.Buyer,
		Capabilities:       req.Capabilities,
		Currency:           req.Currency,
		Status:             acpcheckout.CheckoutSessionBaseStatusReadyForPayment,
		FulfillmentOptions: []acpcheckout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          lineItems,
		Links:              []acpcheckout.Link{},
		Messages:           []acpcheckout.CheckoutSessionBase_Messages_Item{},
		Totals:             totals,
	}
	s.sessions[session.Id] = session
	return cloneSession(session), nil
}

func (s *memoryService) UpdateSession(_ context.Context, id string, req acpcheckout.CheckoutSessionUpdateRequest) (*acpcheckout.CheckoutSessionBase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.sessions[id]
	if !ok {
		return nil, sessionNotFound(id)
	}
	updated := *current
	if req.Buyer != nil {
		updated.Buyer = req.Buyer
	}
	if req.LineItems != nil {
		updated.LineItems, updated.Totals = priceItems(*req.LineItems)
	}
	s.sessions[id] = &updated
	return cloneSession(&updated), nil
}

func (s *memoryService) GetSession(_ context.Context, id string) (*acpcheckout.CheckoutSessionBase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, sessionNotFound(id)
	}
	return cloneSession(session), nil
}

func (s *memoryService) CompleteSession(_ context.Context, id string, req acpcheckout.CheckoutSessionCompleteRequest) (acpcheckout.CheckoutSessionWithOrder, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return acpcheckout.CheckoutSessionWithOrder{}, sessionNotFound(id)
	}
	buyer := session.Buyer
	if req.Buyer != nil {
		buyer = req.Buyer
	}
	return acpcheckout.CheckoutSessionWithOrder{
		Id:                 id,
		Buyer:              buyer,
		Currency:           session.Currency,
		Capabilities:       session.Capabilities,
		FulfillmentOptions: []acpcheckout.CheckoutSessionWithOrder_FulfillmentOptions_Item{},
		LineItems:          append([]acpcheckout.LineItem{}, session.LineItems...),
		Links:              append([]acpcheckout.Link{}, session.Links...),
		Messages:           []acpcheckout.CheckoutSessionWithOrder_Messages_Item{},
		Status:             acpcheckout.CheckoutSessionWithOrderStatusCompleted,
		Totals:             append([]acpcheckout.Total{}, session.Totals...),
		Order: acpcheckout.Order{
			Id:                "ord_demo_1",
			CheckoutSessionId: id,
			PermalinkUrl:      "https://merchant.example/orders/ord_demo_1",
		},
	}, nil
}

func (s *memoryService) CancelSession(_ context.Context, id string, _ *acpcheckout.CancelSessionRequest) (*acpcheckout.CheckoutSessionBase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.sessions[id]
	if !ok {
		return nil, sessionNotFound(id)
	}
	canceled := *current
	canceled.Status = acpcheckout.CheckoutSessionBaseStatusCanceled
	s.sessions[id] = &canceled
	return cloneSession(&canceled), nil
}

func priceItems(items []acpcheckout.Item) ([]acpcheckout.LineItem, []acpcheckout.Total) {
	lineItems := make([]acpcheckout.LineItem, 0, len(items))
	amount := 0
	for _, item := range items {
		name, unitAmount := catalogItem(item.Id)
		item.Name = &name
		item.UnitAmount = &unitAmount
		amount += unitAmount
		lineItems = append(lineItems, acpcheckout.LineItem{
			Id:         "line_" + item.Id,
			Item:       item,
			Name:       &name,
			Quantity:   1,
			UnitAmount: &unitAmount,
			Totals:     []acpcheckout.Total{total(unitAmount)},
		})
	}
	return lineItems, []acpcheckout.Total{total(amount)}
}

func catalogItem(id string) (string, int) {
	switch id {
	case "latte":
		return "Oat Milk Latte", 650
	case "mug":
		return "Stoneware Mug", 1500
	default:
		return id, 0
	}
}

func total(amount int) acpcheckout.Total {
	return acpcheckout.Total{
		Type:        acpcheckout.TotalTypeTotal,
		DisplayText: "Total",
		Amount:      amount,
	}
}

func cloneSession(session *acpcheckout.CheckoutSessionBase) *acpcheckout.CheckoutSessionBase {
	cloned := *session
	cloned.LineItems = append([]acpcheckout.LineItem{}, session.LineItems...)
	cloned.Links = append([]acpcheckout.Link{}, session.Links...)
	cloned.Totals = append([]acpcheckout.Total{}, session.Totals...)
	return &cloned
}

func sessionNotFound(id string) error {
	return acp.NewHTTPError(
		http.StatusNotFound,
		acp.InvalidRequest,
		acp.ErrorCode("checkout_session_not_found"),
		"checkout session not found: "+id,
	)
}

func main() {
	handler := acpcheckout.NewHandler(newMemoryService(), acpauth.StaticTokenAuthorizer("demo-key"))
	log.Println("checkout example listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
