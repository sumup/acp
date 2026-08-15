package acpcheckout

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sumup/acp"
)

// WebhookEventType enumerates the supported checkout webhook events.
type WebhookEventType string

const (
	// WebhookEventTypeOrderCreated is emitted when an order is first created.
	WebhookEventTypeOrderCreated WebhookEventType = "order_create"
	// WebhookEventTypeOrderUpdated is emitted when an existing order changes.
	WebhookEventTypeOrderUpdated WebhookEventType = "order_update"
)

// EventDataType labels the payload for a webhook event.
type EventDataType string

const (
	// EventDataTypeOrder indicates the webhook data payload is an order object.
	EventDataTypeOrder EventDataType = "order"
)

// OrderStatus is an extensible order lifecycle state.
//
// ACP defines common values such as created, confirmed, manual_review,
// processing, shipped, completed, and canceled. Receivers must accept values
// added by later protocol versions.
type OrderStatus string

// EventData is implemented by webhook payloads.
type EventData interface {
	eventType() WebhookEventType
}

// OrderCreate emits order data after the order is created.
type OrderCreate struct {
	// Type is the webhook payload discriminator and is always "order".
	Type EventDataType `json:"type"`
	// ID is the order identifier.
	ID string `json:"id"`
	// CheckoutSessionID identifies the checkout session that produced this order.
	CheckoutSessionID string `json:"checkout_session_id"`
	// OrderNumber is a human-readable order reference.
	OrderNumber *string `json:"order_number,omitempty"`
	// PermalinkURL is the buyer-facing order URL.
	PermalinkURL string `json:"permalink_url"`
	// Status is the order lifecycle state at event emission time.
	Status OrderStatus `json:"status"`
	// EstimatedDelivery is the expected delivery window.
	EstimatedDelivery *EstimatedDelivery `json:"estimated_delivery,omitempty"`
	// Confirmation contains receipt and confirmation metadata.
	Confirmation *OrderConfirmation `json:"confirmation,omitempty"`
	// Support contains merchant support contact details.
	Support *SupportInfo `json:"support,omitempty"`
	// LineItems describes ordered items and item-level fulfillment progress.
	LineItems []OrderLineItem `json:"line_items,omitempty"`
	// Fulfillments describes shipping/pickup/digital delivery state.
	Fulfillments []Fulfillment `json:"fulfillments,omitempty"`
	// Adjustments describes refunds, credits, returns, and disputes.
	Adjustments []Adjustment `json:"adjustments,omitempty"`
	// Totals is an order-level totals breakdown.
	Totals []Total `json:"totals,omitempty"`
}

func (OrderCreate) eventType() WebhookEventType { return WebhookEventTypeOrderCreated }

// OrderUpdated emits order data whenever the order status changes.
type OrderUpdated struct {
	// Type is the webhook payload discriminator and is always "order".
	Type EventDataType `json:"type"`
	// ID is the order identifier.
	ID string `json:"id"`
	// CheckoutSessionID identifies the checkout session that produced this order.
	CheckoutSessionID string `json:"checkout_session_id"`
	// OrderNumber is a human-readable order reference.
	OrderNumber *string `json:"order_number,omitempty"`
	// PermalinkURL is the buyer-facing order URL.
	PermalinkURL string `json:"permalink_url"`
	// Status is the order lifecycle state at event emission time.
	Status OrderStatus `json:"status"`
	// EstimatedDelivery is the expected delivery window.
	EstimatedDelivery *EstimatedDelivery `json:"estimated_delivery,omitempty"`
	// Confirmation contains receipt and confirmation metadata.
	Confirmation *OrderConfirmation `json:"confirmation,omitempty"`
	// Support contains merchant support contact details.
	Support *SupportInfo `json:"support,omitempty"`
	// LineItems describes ordered items and item-level fulfillment progress.
	LineItems []OrderLineItem `json:"line_items,omitempty"`
	// Fulfillments describes shipping/pickup/digital delivery state.
	Fulfillments []Fulfillment `json:"fulfillments,omitempty"`
	// Adjustments describes refunds, credits, returns, and disputes.
	Adjustments []Adjustment `json:"adjustments,omitempty"`
	// Totals is an order-level totals breakdown.
	Totals []Total `json:"totals,omitempty"`
}

func (OrderUpdated) eventType() WebhookEventType { return WebhookEventTypeOrderUpdated }

type webhookEvent struct {
	Type WebhookEventType `json:"type"`
	Data any              `json:"data"`
}

// WebhookOption configures a [WebhookSender] returned by [Handler.GetWebhookSender].
type WebhookOption func(*webhookSender)

// WebhookWithClient allows overriding the HTTP client used for delivering webhook events.
func WebhookWithClient(client *http.Client) WebhookOption {
	return func(ws *webhookSender) {
		ws.client = client
	}
}

// WebhookSender delivers checkout webhook events to an ACP agent callback endpoint.
type WebhookSender interface {
	// Send sends webhook to the configured webhook endpoint.
	Send(context.Context, EventData) error
}

// NewWebhookSender returns a configured [WebhookSender] that delivers webhooks back to the agent.
//
// Parameters:
//
//   - endpoint is the absolute URL provided by the agent for receiving webhook events.
//   - secret is the HMAC secret provided by the agent for signing webhook payloads.
//
// Deprecated: use package acpwebhook's NewSender with its WebhookEvent for the
// complete 2026-04-17 order model.
func NewWebhookSender(endpoint string, secret []byte, opts ...WebhookOption) (WebhookSender, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint url: %w", err)
	}
	if !endpointURL.IsAbs() || endpointURL.Host == "" || (endpointURL.Scheme != "http" && endpointURL.Scheme != "https") {
		return nil, fmt.Errorf("endpoint must be an absolute HTTP(S) URL")
	}

	if len(secret) == 0 {
		return nil, fmt.Errorf("secret is required")
	}
	secretKey := append([]byte(nil), secret...)

	sender := &webhookSender{
		endpoint: endpointURL,
		secret:   secretKey,
		client:   http.DefaultClient,
		now:      time.Now,
	}

	for _, opt := range opts {
		opt(sender)
	}
	if sender.client == nil {
		return nil, fmt.Errorf("HTTP client is required")
	}
	if sender.now == nil {
		return nil, fmt.Errorf("clock is required")
	}

	return sender, nil
}

// GetWebhookSender returns a configured [WebhookSender] that delivers webhooks back to the agent.
//
// The merchantName argument is retained for source compatibility but is no
// longer used: ACP 2026-04-17 fixes the header name to Merchant-Signature.
//
// Deprecated: use [NewWebhookSender] or package acpwebhook's NewSender.
func (*Handler) GetWebhookSender(endpoint, _ string, secret []byte, opts ...WebhookOption) (WebhookSender, error) {
	return NewWebhookSender(endpoint, secret, opts...)
}

type webhookSender struct {
	endpoint *url.URL
	secret   []byte
	client   *http.Client
	now      func() time.Time
}

// Send sends webhook to the configured webhook endpoint.
func (s *webhookSender) Send(ctx context.Context, data EventData) error {
	body, err := json.Marshal(webhookEvent{
		Type: data.eventType(),
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("checkout: marshal webhook payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("checkout: build webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Version", acp.APIVersion)
	timestamp := s.now().UTC()
	req.Header.Set("Timestamp", timestamp.Format(time.RFC3339))
	req.Header.Set("Merchant-Signature", signWebhookPayload(s.secret, timestamp.Unix(), body))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("checkout: send webhook: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("checkout: webhook endpoint %s returned %s: %s", s.endpoint, resp.Status, strings.TrimSpace(string(snippet)))
	}
	return nil
}

func signWebhookPayload(secret []byte, timestamp int64, payload []byte) string {
	seconds := strconv.FormatInt(timestamp, 10)
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(seconds))
	_, _ = mac.Write([]byte("."))
	_, _ = mac.Write(payload)
	return fmt.Sprintf("t=%s,v1=%x", seconds, mac.Sum(nil))
}
