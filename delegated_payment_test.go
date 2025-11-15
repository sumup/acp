package acp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sumup/acp/secret"
)

func TestDelegatedPaymentHandler(t *testing.T) {
	t.Parallel()

	reqPayload := sampleDelegatePaymentRequest()
	service := &delegatedStubService{
		delegate: func(ctx context.Context, req PaymentRequest) (*VaultToken, error) {
			if req.Allowance.MerchantID != "acme" {
				t.Fatalf("unexpected merchant id %s", req.Allowance.MerchantID)
			}
			return &VaultToken{
				ID:       "vt_token",
				Created:  time.Now().UTC(),
				Metadata: map[string]string{"source": "test"},
			}, nil
		},
	}
	handler := NewDelegatedPaymentHandler(service)

	body, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 got %d body=%s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("API-Version"); got != APIVersion {
		t.Fatalf("expected API-Version header %s got %s", APIVersion, got)
	}
	var resp VaultToken
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.ID != "vt_token" {
		t.Fatalf("unexpected response id %s", resp.ID)
	}
}

func TestDelegatedPaymentHandlerErrors(t *testing.T) {
	t.Parallel()

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()

		handler := NewDelegatedPaymentHandler(&delegatedStubService{})
		req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", strings.NewReader("{"))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 got %d", rec.Code)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		payload := sampleDelegatePaymentRequest()
		payload.Metadata = nil
		body, _ := json.Marshal(payload)
		handler := NewDelegatedPaymentHandler(&delegatedStubService{})
		req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 got %d", rec.Code)
		}
	})

	t.Run("service error surfaces", func(t *testing.T) {
		t.Parallel()

		handler := NewDelegatedPaymentHandler(&delegatedStubService{
			delegate: func(ctx context.Context, req PaymentRequest) (*VaultToken, error) {
				return nil, NewHTTPError(http.StatusConflict, InvalidRequest, IdempotencyConflict, "idempotency conflict")
			},
		})
		body, _ := json.Marshal(sampleDelegatePaymentRequest())
		req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("expected 409 got %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "idempotency conflict") {
			t.Fatalf("unexpected body %s", rec.Body.String())
		}
	})

	t.Run("rate limited response sets retry-after header", func(t *testing.T) {
		t.Parallel()

		handler := NewDelegatedPaymentHandler(&delegatedStubService{
			delegate: func(ctx context.Context, req PaymentRequest) (*VaultToken, error) {
				return nil, NewRateLimitExceededError("too many attempts", WithRetryAfter(3*time.Second))
			},
		})
		body, _ := json.Marshal(sampleDelegatePaymentRequest())
		req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusTooManyRequests {
			t.Fatalf("expected 429 got %d", rec.Code)
		}
		if got := rec.Header().Get("Retry-After"); got != "3" {
			t.Fatalf("expected Retry-After header 3 got %s", got)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		t.Parallel()

		handler := NewDelegatedPaymentHandler(&delegatedStubService{})
		req := httptest.NewRequest(http.MethodGet, "/agentic_commerce/delegate_payment", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405 got %d", rec.Code)
		}
	})
}

type delegatedStubService struct {
	delegate func(context.Context, PaymentRequest) (*VaultToken, error)
}

func (s *delegatedStubService) DelegatePayment(ctx context.Context, req PaymentRequest) (*VaultToken, error) {
	if s.delegate != nil {
		return s.delegate(ctx, req)
	}
	return nil, NewHTTPError(http.StatusNotImplemented, InvalidRequest, ErrorCode("not_implemented"), "delegate payment not implemented")
}

func sampleDelegatePaymentRequest() PaymentRequest {
	expMonth := "11"
	expYear := "2026"
	displayLast4 := "4242"
	checks := []CardChecksPerformed{CardChecksPerformedAVS}

	return PaymentRequest{
		PaymentMethod: PaymentMethodCard{
			Type:                   PaymentMethodCardTypeCard,
			CardNumberType:         CardCardNumberTypeFPAN,
			Number:                 secret.New("4242424242424242"),
			ExpMonth:               &expMonth,
			ExpYear:                &expYear,
			DisplayLast4:           &displayLast4,
			DisplayCardFundingType: CardFundingTypeCredit,
			Metadata:               map[string]string{"issuer": "acme"},
			ChecksPerformed:        checks,
		},
		Allowance: Allowance{
			Reason:            AllowanceReasonOneTime,
			MaxAmount:         2000,
			Currency:          "usd",
			CheckoutSessionID: "csn_123",
			MerchantID:        "acme",
			ExpiresAt:         time.Now().Add(time.Hour).UTC(),
		},
		RiskSignals: []RiskSignal{
			{
				Type:   RiskSignalTypeCardTesting,
				Score:  10,
				Action: RiskSignalActionManualReview,
			},
		},
		Metadata: map[string]string{
			"campaign": "q4",
		},
	}
}
