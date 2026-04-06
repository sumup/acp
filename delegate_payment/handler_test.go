package delegate_payment_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/delegate_payment"
	"go.uber.org/mock/gomock"
)

func TestDelegatedPaymentHandler_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockProvider := NewMockDelegatedPaymentProvider(ctrl)
	mockProvider.EXPECT().
		DelegatePayment(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req delegate_payment.DelegatePaymentRequest) (*delegate_payment.DelegatePaymentResponse, error) {
			requestCtx := acp.RequestContextFromContext(ctx)
			if requestCtx == nil {
				t.Fatal("expected request context on handler context")
			}
			if requestCtx.Authorization != "Bearer test-key" {
				t.Fatalf("unexpected authorization %q", requestCtx.Authorization)
			}
			if requestCtx.APIVersion != acp.APIVersion {
				t.Fatalf("unexpected api version %q", requestCtx.APIVersion)
			}
			return &delegate_payment.DelegatePaymentResponse{
				Id:       "vt_1",
				Created:  time.Now().UTC(),
				Metadata: map[string]string{"source": "test"},
			}, nil
		})

	h := delegate_payment.NewDelegatedPaymentHandler(mockProvider, acpauth.StaticTokenAuthorizer("test-key"))

	body := []byte(`{"payment_method":{"type":"card","card_number_type":"fpan","number":"4242424242424242","display_card_funding_type":"credit","metadata":{}},"allowance":{"reason":"one_time","max_amount":1000,"currency":"usd","checkout_session_id":"cs_1","merchant_id":"m_1","expires_at":"2026-12-31T00:00:00Z"},"risk_signals":[{"type":"card_testing","action":"authorized","score":0}],"metadata":{}}`)
	req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("API-Version", "2026-01-30")
	req.Header.Set("Idempotency-Key", "idem-1")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestDelegatedPaymentHandler(t *testing.T) {
	t.Parallel()

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)

		h := delegate_payment.NewDelegatedPaymentHandler(NewMockDelegatedPaymentProvider(ctrl), acpauth.StaticTokenAuthorizer("test-key"))

		body := []byte(`{"payment_method":{"type":"card","card_number_type":"fpan","number":"4242424242424242","display_card_funding_type":"credit","metadata":{}},"allowance":{"reason":"one_time","max_amount":1000,"currency":"usd","checkout_session_id":"cs_1","merchant_id":"m_1","expires_at":"2026-12-31T00:00:00Z"},"risk_signals":[{"type":"card_testing","action":"authorized","score":0}],"metadata":{}}`)
		req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("API-Version", "2026-01-30")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 got %d body=%s", rec.Code, rec.Body.String())
		}

		var got delegate_payment.Error
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if got.Code != delegate_payment.ErrorCode("missing_authorization") {
			t.Fatalf("expected missing_authorization got %q", got.Code)
		}
		if got.Message != "Authorization header is required" {
			t.Fatalf("unexpected message: %q", got.Message)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)

		h := delegate_payment.NewDelegatedPaymentHandler(NewMockDelegatedPaymentProvider(ctrl), acpauth.StaticTokenAuthorizer("test-key"))

		body := []byte(`{"payment_method":{"type":"card","card_number_type":"fpan","number":"4242424242424242","display_card_funding_type":"credit","display_last4":"42","metadata":{}},"allowance":{"reason":"one_time","max_amount":1000,"currency":"usd","checkout_session_id":"cs_1","merchant_id":"m_1","expires_at":"2026-12-31T00:00:00Z"},"risk_signals":[{"type":"card_testing","action":"authorized","score":0}],"metadata":{}}`)
		req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-key")
		req.Header.Set("API-Version", "2026-01-30")
		req.Header.Set("Idempotency-Key", "idem-1")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422 got %d body=%s", rec.Code, rec.Body.String())
		}

		var got delegate_payment.Error
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if got.Code != delegate_payment.InvalidCard {
			t.Fatalf("expected invalid_card got %q", got.Code)
		}
		if got.Param == nil || *got.Param != "$.payment_method.display_last4" {
			t.Fatalf("unexpected param: %#v", got.Param)
		}
	})
}
