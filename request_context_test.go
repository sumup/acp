package acp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestContextFromRequest(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", nil)
	req.Header.Set("Authorization", " Bearer api_key ")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "acp-test/1.0")
	req.Header.Set("Idempotency-Key", "idem-123")
	req.Header.Set("Request-Id", "req-123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature", "sig-123")
	req.Header.Set("Timestamp", "2025-01-02T03:04:05Z")
	req.Header.Set("API-Version", "2025-01-01")

	got := requestContextFromRequest(req)
	if got == nil {
		t.Fatalf("expected request context")
	}
	if got.Authorization != "Bearer api_key" {
		t.Fatalf("unexpected authorization %q", got.Authorization)
	}
	if got.AcceptLanguage != "en-US" {
		t.Fatalf("unexpected accept-language %q", got.AcceptLanguage)
	}
	if got.UserAgent != "acp-test/1.0" {
		t.Fatalf("unexpected user-agent %q", got.UserAgent)
	}
	if got.IdempotencyKey != "idem-123" {
		t.Fatalf("unexpected idempotency key %q", got.IdempotencyKey)
	}
	if got.RequestID != "req-123" {
		t.Fatalf("unexpected request id %q", got.RequestID)
	}
	if got.Signature != "sig-123" {
		t.Fatalf("unexpected signature %q", got.Signature)
	}
	if got.Timestamp != "2025-01-02T03:04:05Z" {
		t.Fatalf("unexpected timestamp %q", got.Timestamp)
	}
	if got.APIVersion != "2025-01-01" {
		t.Fatalf("unexpected api version %q", got.APIVersion)
	}
}

func TestRequestContextRoundTrip(t *testing.T) {
	t.Parallel()

	requestCtx := &RequestContext{Authorization: "Bearer 123"}
	ctx := contextWithRequestContext(context.Background(), requestCtx)
	got := RequestContextFromContext(ctx)
	if got == nil {
		t.Fatalf("expected request context on context")
	}
	if got.Authorization != "Bearer 123" {
		t.Fatalf("unexpected authorization %q", got.Authorization)
	}
	if RequestContextFromContext(context.Background()) != nil {
		t.Fatalf("expected nil when request context not set")
	}
}
