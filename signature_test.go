package acp

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sumup/acp/signature"
)

func TestSignatureMiddlewareAllowsValidRequest(t *testing.T) {
	t.Parallel()

	key := []byte("secret")
	ts := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	handler := NewCheckoutHandler(&stubService{
		create: func(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSession, error) {
			return &CheckoutSession{
				ID:                 "cs_123",
				Status:             CheckoutSessionStatusInProgress,
				Currency:           "usd",
				LineItems:          []LineItem{},
				FulfillmentOptions: make([]FulfillmentOption, 0),
				Totals:             []Total{},
				Messages:           make([]Message, 0),
				Links:              []Link{},
			}, nil
		},
	}, WithSignatureVerifier(signature.HMACVerifier{Key: key}), checkoutWithClock(func() time.Time {
		return ts.Add(30 * time.Second)
	}))

	body := []byte(`{"items":[{"id":"sku_1","quantity":1}]}`)
	canonical, err := signature.CanonicalizeJSONBody(body)
	if err != nil {
		t.Fatalf("canonicalize: %v", err)
	}
	signature := signFixture(key, ts, canonical)

	req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature", signature)
	req.Header.Set("Timestamp", ts.Format(time.RFC3339Nano))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestSignatureMiddlewareRejectsInvalidSignature(t *testing.T) {
	t.Parallel()

	key := []byte("secret")
	ts := time.Now().UTC()
	handler := NewCheckoutHandler(&stubService{
		create: func(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSession, error) {
			return &CheckoutSession{}, nil
		},
	}, WithSignatureVerifier(signature.HMACVerifier{Key: key}), checkoutWithClock(func() time.Time {
		return ts
	}))

	body := []byte(`{"items":[{"id":"sku_1","quantity":1}]}`)
	req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature", "bogus")
	req.Header.Set("Timestamp", ts.Format(time.RFC3339Nano))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
}

func TestSignatureMiddlewareRejectsSkew(t *testing.T) {
	t.Parallel()

	key := []byte("secret")
	ts := time.Now().UTC()
	handler := NewCheckoutHandler(&stubService{
		create: func(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSession, error) {
			return &CheckoutSession{}, nil
		},
	}, WithSignatureVerifier(signature.HMACVerifier{Key: key}), WithMaxClockSkew(time.Minute), checkoutWithClock(func() time.Time {
		return ts.Add(2 * time.Minute)
	}))

	body := []byte(`{"items":[{"id":"sku_1","quantity":1}]}`)
	canonical, err := signature.CanonicalizeJSONBody(body)
	if err != nil {
		t.Fatalf("canonicalize: %v", err)
	}
	signature := signFixture(key, ts, canonical)

	req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature", signature)
	req.Header.Set("Timestamp", ts.Format(time.RFC3339Nano))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
	if want, got := "stale_timestamp", getErrorCode(rec.Body.Bytes()); want != got {
		t.Fatalf("expected code %s got %s", want, got)
	}
}

func TestSignatureMiddlewareRequiresHeadersWhenEnforced(t *testing.T) {
	t.Parallel()

	handler := NewCheckoutHandler(&stubService{
		get: func(ctx context.Context, id string) (*CheckoutSession, error) {
			return &CheckoutSession{}, nil
		},
	}, WithSignatureVerifier(signature.HMACVerifier{Key: []byte("secret")}), WithRequireSignedRequests(), checkoutWithClock(time.Now))

	req := httptest.NewRequest(http.MethodGet, "/checkout_sessions/cs_123", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
	if want, got := "signature_required", getErrorCode(rec.Body.Bytes()); want != got {
		t.Fatalf("expected code %s got %s", want, got)
	}
}

func signFixture(key []byte, ts time.Time, canonical []byte) string {
	payload := signature.BuildSigningPayload(ts, canonical)
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func getErrorCode(body []byte) string {
	var resp Error
	if err := json.Unmarshal(body, &resp); err != nil {
		return ""
	}
	return string(resp.Code)
}
