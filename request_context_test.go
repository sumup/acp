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
	req.Header.Set("API-Version", APIVersion)

	got, err := RequestContextFromRequest(req)
	if err != nil {
		t.Fatalf("RequestContextFromRequest() error = %v", err)
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
	if got.APIVersion != APIVersion {
		t.Fatalf("unexpected api version %q", got.APIVersion)
	}
}

func TestRequestContextRoundTrip(t *testing.T) {
	t.Parallel()

	requestCtx := &RequestContext{Authorization: "Bearer 123"}
	ctx := ContextWithRequestContext(context.Background(), requestCtx)
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

func TestRequestContext_validate(t *testing.T) {
	t.Parallel()

	t.Run("valid request", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/checkout_sessions/cs_123", nil)
		req.Header.Set("Authorization", "Bearer test-key")
		req.Header.Set("API-Version", APIVersion)

		ctx, err := RequestContextFromRequest(req)
		if err != nil {
			t.Fatalf("RequestContextFromRequest() error = %v", err)
		}
		if ctx == nil {
			t.Fatal("expected request context")
		}
	})

	t.Run("missing authorization", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/checkout_sessions/cs_123", nil)
		req.Header.Set("API-Version", APIVersion)

		ctx, err := RequestContextFromRequest(req)
		if err != nil {
			t.Fatalf("RequestContextFromRequest() error = %v", err)
		}
		if ctx.Authorization != "" {
			t.Fatalf("RequestContextFromRequest() authorization = %q, want empty", ctx.Authorization)
		}
	})

	t.Run("invalid api version format", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/checkout_sessions/cs_123", nil)
		req.Header.Set("Authorization", "Bearer test-key")
		req.Header.Set("API-Version", "20260130")

		_, err := RequestContextFromRequest(req)
		if err == nil {
			t.Fatal("RequestContextFromRequest() error = nil, want error")
		}
		acpErr := err.(*Error)
		if acpErr.Message != "API-Version header must be in YYYY-MM-DD format" {
			t.Fatalf("RequestContextFromRequest() message = %q", acpErr.Message)
		}
	})

	t.Run("missing api version", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/checkout_sessions/cs_123", nil)
		req.Header.Set("Authorization", "Bearer test-key")

		_, err := RequestContextFromRequest(req)
		if err == nil {
			t.Fatal("RequestContextFromRequest() error = nil, want error")
		}
		acpErr := err.(*Error)
		if acpErr.Message != "API-Version header is required" {
			t.Fatalf("RequestContextFromRequest() message = %q", acpErr.Message)
		}
		if acpErr.Code != MissingAPIVersion {
			t.Fatalf("RequestContextFromRequest() code = %q", acpErr.Code)
		}
		if len(acpErr.SupportedVersions) != 1 || acpErr.SupportedVersions[0] != APIVersion {
			t.Fatalf("RequestContextFromRequest() supported versions = %#v", acpErr.SupportedVersions)
		}
	})

	t.Run("unsupported api version", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/checkout_sessions/cs_123", nil)
		req.Header.Set("Authorization", "Bearer test-key")
		req.Header.Set("API-Version", "2026-01-30")

		_, err := RequestContextFromRequest(req)
		if err == nil {
			t.Fatal("RequestContextFromRequest() error = nil, want error")
		}
		acpErr := err.(*Error)
		if acpErr.Code != UnsupportedAPIVersion {
			t.Fatalf("RequestContextFromRequest() code = %q", acpErr.Code)
		}
		if len(acpErr.SupportedVersions) != 1 || acpErr.SupportedVersions[0] != APIVersion {
			t.Fatalf("RequestContextFromRequest() supported versions = %#v", acpErr.SupportedVersions)
		}
	})
}
