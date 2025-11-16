package acp

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
	WebhookEventTypeOrderCreated WebhookEventType = "order_created"
	WebhookEventTypeOrderUpdated WebhookEventType = "order_updated"
)

// EventDataType labels the payload for a webhook event.
type EventDataType string

const (
	EventDataTypeOrder EventDataType = "order"
)

// OrderStatus defines model for webhook data status.
type OrderStatus string

const (
	OrderStatusCreated      OrderStatus = "created"
	OrderStatusManualReview OrderStatus = "manual_review"
	OrderStatusConfirmed    OrderStatus = "confirmed"
	OrderStatusCanceled     OrderStatus = "canceled"
	OrderStatusShipped      OrderStatus = "shipped"
	OrderStatusFulfilled    OrderStatus = "fulfilled"
)

// RefundType captures the source of refunded funds.
type RefundType string

const (
	RefundTypeStoreCredit     RefundType = "store_credit"
	RefundTypeOriginalPayment RefundType = "original_payment"
)

// Refund describes a refund emitted in webhook events.
type Refund struct {
	Type   RefundType `json:"type"`
	Amount int        `json:"amount"`
}

// EventData is implemented by webhook payloads.
type EventData interface {
	eventType() WebhookEventType
}

// OrderCreate emits order data after the order is created.
type OrderCreate struct {
	Type              EventDataType `json:"type"`
	CheckoutSessionID string        `json:"checkout_session_id"`
	PermalinkURL      string        `json:"permalink_url"`
	Status            OrderStatus   `json:"status"`
	Refunds           []Refund      `json:"refunds"`
}

func (OrderCreate) eventType() WebhookEventType { return WebhookEventTypeOrderCreated }

// OrderUpdated emits order data whenever the order status changes.
type OrderUpdated struct {
	Type              EventDataType `json:"type"`
	CheckoutSessionID string        `json:"checkout_session_id"`
	PermalinkURL      string        `json:"permalink_url"`
	Status            OrderStatus   `json:"status"`
	Refunds           []Refund      `json:"refunds"`
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
	req.Header.Set("API-Version", APIVersion)
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
