package acp

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// SendWebhook posts webhook events to the OpenAI endpoint configured via [WithWebhookOptions].
func (h *CheckoutHandler) SendWebhook(ctx context.Context, data EventData) error {
	if h.cfg.webhook == nil {
		return errors.New("checkout: webhook options must be configured")
	}
	body, err := json.Marshal(webhookEvent{
		Type: data.eventType(),
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("checkout: marshal webhook payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.cfg.webhook.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("checkout: build webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Version", APIVersion)
	req.Header.Set(h.cfg.webhook.header, signWebhookPayload(h.cfg.webhook.secret, body))

	resp, err := h.cfg.webhook.client.Do(req)
	if err != nil {
		return fmt.Errorf("checkout: send webhook: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("checkout: webhook endpoint %s returned %s: %s", h.cfg.webhook.endpoint, resp.Status, strings.TrimSpace(string(snippet)))
	}
	return nil
}

func signWebhookPayload(secret, payload []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
