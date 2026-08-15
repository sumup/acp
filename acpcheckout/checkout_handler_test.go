package acpcheckout_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpcheckout"
	"go.uber.org/mock/gomock"
)

func TestCheckoutHandler_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockProvider := NewMockProvider(ctrl)
	mockProvider.EXPECT().
		CreateSession(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req acpcheckout.CheckoutSessionCreateRequest) (*acpcheckout.CheckoutSessionBase, error) {
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
			return &acpcheckout.CheckoutSessionBase{
				Id:                 "cs_1",
				Currency:           "USD",
				Status:             acpcheckout.CheckoutSessionBaseStatusInProgress,
				FulfillmentOptions: []acpcheckout.CheckoutSessionBase_FulfillmentOptions_Item{},
				LineItems:          []acpcheckout.LineItem{},
				Links:              []acpcheckout.Link{},
				Messages:           []acpcheckout.CheckoutSessionBase_Messages_Item{},
			}, nil
		})

	h := acpcheckout.NewHandler(mockProvider, acpauth.StaticTokenAuthorizer("test-key"))

	body := []byte(`{"capabilities":{},"currency":"USD","line_items":[{"id":"sku_1"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("API-Version", acp.APIVersion)
	req.Header.Set("Idempotency-Key", "idem-1")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestCheckoutHandler_WithServeMux(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockProvider := NewMockProvider(ctrl)
	mockProvider.EXPECT().
		CreateSession(gomock.Any(), gomock.Any()).
		Return(&acpcheckout.CheckoutSessionBase{
			Id:                 "cs_1",
			Currency:           "USD",
			Status:             acpcheckout.CheckoutSessionBaseStatusInProgress,
			FulfillmentOptions: []acpcheckout.CheckoutSessionBase_FulfillmentOptions_Item{},
			LineItems:          []acpcheckout.LineItem{},
			Links:              []acpcheckout.Link{},
			Messages:           []acpcheckout.CheckoutSessionBase_Messages_Item{},
		}, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	h := acpcheckout.NewHandler(mockProvider, acpauth.StaticTokenAuthorizer("test-key"), acpcheckout.WithServeMux(mux))

	healthReq := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthRec := httptest.NewRecorder()
	h.ServeHTTP(healthRec, healthReq)

	if healthRec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d body=%s", healthRec.Code, healthRec.Body.String())
	}

	checkoutBody := []byte(`{"capabilities":{},"currency":"USD","line_items":[{"id":"sku_1"}]}`)
	checkoutReq := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(checkoutBody))
	checkoutReq.Header.Set("Content-Type", "application/json")
	checkoutReq.Header.Set("Authorization", "Bearer test-key")
	checkoutReq.Header.Set("API-Version", acp.APIVersion)
	checkoutReq.Header.Set("Idempotency-Key", "idem-1")
	checkoutRec := httptest.NewRecorder()

	h.ServeHTTP(checkoutRec, checkoutReq)

	if checkoutRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", checkoutRec.Code, checkoutRec.Body.String())
	}
}

func TestCheckoutHandler(t *testing.T) {
	t.Parallel()

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)

		h := acpcheckout.NewHandler(NewMockProvider(ctrl), acpauth.StaticTokenAuthorizer("test-key"))

		body := []byte(`{"capabilities":{},"currency":"USD","line_items":[{"id":"sku_1"}]}`)
		req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("API-Version", acp.APIVersion)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 got %d body=%s", rec.Code, rec.Body.String())
		}

		var got acp.Error
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if got.Code != "missing_authorization" {
			t.Fatalf("expected missing_authorization got %q", got.Code)
		}
	})

	t.Run("semantic validation error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)

		h := acpcheckout.NewHandler(NewMockProvider(ctrl), acpauth.StaticTokenAuthorizer("test-key"))

		body := []byte(`{"capabilities":{},"currency":"USD","line_items":[]}`)
		req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-key")
		req.Header.Set("API-Version", acp.APIVersion)
		req.Header.Set("Idempotency-Key", "idem-1")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422 got %d body=%s", rec.Code, rec.Body.String())
		}

		var got acp.Error
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if got.Param == nil || *got.Param != "$.line_items" {
			t.Fatalf("unexpected param: %#v", got.Param)
		}
	})
}
