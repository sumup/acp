package acp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckoutHandlerSendWebhook(t *testing.T) {
	t.Parallel()

	var received struct {
		body   []byte
		header http.Header
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, _ := io.ReadAll(r.Body)
		received.body = payload
		received.header = r.Header.Clone()
		w.WriteHeader(http.StatusAccepted)
	}))
	t.Cleanup(srv.Close)

	handler := NewCheckoutHandler(&stubService{}, WithWebhookOptions(WebhookOptions{
		Endpoint:   srv.URL,
		HeaderName: "Merchant_Name-Signature",
		SecretKey:  []byte("super-secret"),
		Client:     srv.Client(),
	}))

	event := OrderCreate{
		Type:              "order",
		CheckoutSessionID: "cs_123",
		PermalinkURL:      "https://merchant.example/orders/cs_123",
		Status:            OrderStatusCreated,
	}
	if err := handler.SendWebhook(context.Background(), event); err != nil {
		t.Fatalf("SendWebhook() error = %v", err)
	}

	if got := received.header.Get("API-Version"); got != APIVersion {
		t.Fatalf("missing API-Version header, got %q", got)
	}
	sig := received.header.Get("Merchant_Name-Signature")
	expectedSig := signWebhookPayload([]byte("super-secret"), received.body)
	if sig != expectedSig {
		t.Fatalf("unexpected signature header %q", sig)
	}

	var decoded struct {
		Type WebhookEventType `json:"type"`
		Data OrderCreate      `json:"data"`
	}
	if err := json.Unmarshal(received.body, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if decoded.Type != WebhookEventTypeOrderCreated {
		t.Fatalf("unexpected webhook type %s", decoded.Type)
	}
	if decoded.Data.Type != EventDataTypeOrder {
		t.Fatalf("expected data.type order got %s", decoded.Data.Type)
	}
	if decoded.Data.CheckoutSessionID != event.CheckoutSessionID {
		t.Fatalf("unexpected checkout_session_id %s", decoded.Data.CheckoutSessionID)
	}
}
