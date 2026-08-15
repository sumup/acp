package acpwebhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSenderSignsEvent(t *testing.T) {
	t.Parallel()

	secret := []byte("test-secret")
	now := time.Unix(1_709_123_456, 0)
	var received bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
			return
		}
		mac := hmac.New(sha256.New, secret)
		_, _ = fmt.Fprintf(mac, "%d.%s", now.Unix(), body)
		want := fmt.Sprintf("t=%d,v1=%x", now.Unix(), mac.Sum(nil))
		if got := r.Header.Get("Merchant-Signature"); got != want {
			t.Errorf("Merchant-Signature = %q, want %q", got, want)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q", got)
		}
		received = true
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	sender, err := NewSender(server.URL, secret, withClock(func() time.Time { return now }))
	if err != nil {
		t.Fatalf("NewSender() error = %v", err)
	}
	event := WebhookEvent{
		Type: "order_create",
		Data: EventDataOrder{
			Type:              EventDataOrderTypeOrder,
			Id:                "order_1",
			CheckoutSessionId: "checkout_1",
			PermalinkUrl:      "https://merchant.example/orders/order_1",
			Status:            "created",
		},
	}
	if err := sender.Send(context.Background(), event); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if !received {
		t.Fatal("webhook was not received")
	}
}

func TestNewSenderRejectsRelativeEndpoint(t *testing.T) {
	t.Parallel()

	if _, err := NewSender("/webhooks", []byte("secret")); err == nil {
		t.Fatal("NewSender() error = nil")
	}
}
