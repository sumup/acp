package acp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthenticationMiddlewareRequiresAuthorizationHeader(t *testing.T) {
	t.Parallel()

	handler := NewDelegatedPaymentHandler(successService(), WithAuthenticator(AuthenticatorFunc(func(ctx context.Context, key string) error {
		return nil
	})))

	req := newDelegatePaymentHTTPRequest(t)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rec.Code, rec.Body.String())
	}
	var payload Error
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Code != MissingAuthorization {
		t.Fatalf("expected error code %s got %s", MissingAuthorization, payload.Code)
	}
}

func TestAuthenticationMiddlewareValidatesBearerFormat(t *testing.T) {
	t.Parallel()

	handler := NewDelegatedPaymentHandler(successService(), WithAuthenticator(AuthenticatorFunc(func(ctx context.Context, key string) error {
		return nil
	})))

	req := newDelegatePaymentHTTPRequest(t)
	req.Header.Set("Authorization", "Token abc")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
	var payload Error
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Code != InvalidAuthorization {
		t.Fatalf("expected error code %s got %s", InvalidAuthorization, payload.Code)
	}
}

func TestAuthenticationMiddlewareRejectsInvalidAPIKey(t *testing.T) {
	t.Parallel()

	handler := NewDelegatedPaymentHandler(successService(), WithAuthenticator(AuthenticatorFunc(func(ctx context.Context, key string) error {
		return errors.New("invalid api key")
	})))

	req := newDelegatePaymentHTTPRequest(t)
	req.Header.Set("Authorization", "Bearer bad-key")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
	var payload Error
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Code != InvalidAuthorization {
		t.Fatalf("expected error code %s got %s", InvalidAuthorization, payload.Code)
	}
}

func TestAuthenticationMiddlewareSurfacesAuthenticatorErrors(t *testing.T) {
	t.Parallel()

	authErr := NewHTTPError(http.StatusServiceUnavailable, ServiceUnavailable, ErrorCode(ServiceUnavailable), "auth service unavailable")
	handler := NewDelegatedPaymentHandler(successService(), WithAuthenticator(AuthenticatorFunc(func(ctx context.Context, key string) error {
		return authErr
	})))

	req := newDelegatePaymentHTTPRequest(t)
	req.Header.Set("Authorization", "Bearer auth-down")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 got %d", rec.Code)
	}
	var payload Error
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Type != ServiceUnavailable {
		t.Fatalf("expected error type %s got %s", ServiceUnavailable, payload.Type)
	}
}

func TestAuthenticationMiddlewareAllowsValidRequests(t *testing.T) {
	t.Parallel()

	handler := NewDelegatedPaymentHandler(successService(), WithAuthenticator(AuthenticatorFunc(func(ctx context.Context, key string) error {
		if key != "valid-key" {
			return errors.New("invalid")
		}
		return nil
	})))

	req := newDelegatePaymentHTTPRequest(t)
	req.Header.Set("Authorization", "Bearer valid-key")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rec.Code, rec.Body.String())
	}
}

func newDelegatePaymentHTTPRequest(t *testing.T) *http.Request {
	t.Helper()

	body, err := json.Marshal(sampleDelegatePaymentRequest())
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func successService() *delegatedStubService {
	return &delegatedStubService{
		delegate: func(ctx context.Context, req PaymentRequest) (*VaultToken, error) {
			return &VaultToken{
				ID:       "vt_success",
				Created:  time.Now().UTC(),
				Metadata: map[string]string{"source": "test"},
			}, nil
		},
	}
}
