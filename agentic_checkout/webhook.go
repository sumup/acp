package agentic_checkout

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

type OrderLineItem map[string]any
type Fulfillment map[string]any
type Adjustment map[string]any

// EventData is implemented by webhook payloads.
type EventData interface {
	eventType() WebhookEventType
}

// OrderCreate emits order data after the order is created.
type OrderCreate struct {
	// Type is the webhook payload discriminator and is always "order".
	Type EventDataType `json:"type"`
	// ID is the order identifier.
	ID *string `json:"id,omitempty"`
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
	ID *string `json:"id,omitempty"`
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

type WebhookOption func(*webhookSender)

// WebhookWithClient allows overriding the HTTP client used for delivering webhook events.
func WebhookWithClient(client *http.Client) WebhookOption {
	return func(ws *webhookSender) {
		ws.client = client
	}
}

// WebhookSender is an interface of a webhook delivery implementation.
type WebhookSender interface {
	// Send sends webhook to the configured webhook endpoint.
	Send(context.Context, EventData) error
}

// GetWebhookSender returns a configued [WebhookSender] that allows your implementation to deliver webhooks
// back to the agent.
//
// Parameters:
//
//   - endpoint is the absolute URL provided by OpenAI for receiving webhook events.
//   - merchantName controls the signature header name (the header name is Merchant_Name-Signature).
//   - secret is the HMAC secret provided by OpenAI for signing webhook payloads.
func (h *CheckoutHandler) GetWebhookSender(endpoint, merchantName string, secret []byte, opts ...WebhookOption) (WebhookSender, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint url: %w", err)
	}

	merchantName = strings.TrimSpace(merchantName)
	if merchantName == "" {
		return nil, fmt.Errorf("merchant name is required")
	}
	header := fmt.Sprintf("%s-Signature", strings.ReplaceAll(merchantName, " ", "_"))

	if len(secret) == 0 {
		return nil, fmt.Errorf("secret is required")
	}
	secretKey := append([]byte(nil), secret...)

	sender := &webhookSender{
		endpoint:        endpointURL,
		signatureHeader: header,
		secret:          secretKey,
		client:          http.DefaultClient,
	}

	for _, opt := range opts {
		opt(sender)
	}

	return sender, nil
}

type webhookSender struct {
	endpoint        *url.URL
	signatureHeader string
	secret          []byte
	client          *http.Client
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
	req.Header.Set("API-Version", apiVersionHeaderValue)
	req.Header.Set(s.signatureHeader, signWebhookPayload(s.secret, body))

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

func signWebhookPayload(secret, payload []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
