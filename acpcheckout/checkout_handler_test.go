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
	mockProvider := NewMockCheckoutProvider(ctrl)
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

	h := acpcheckout.NewCheckoutHandler(mockProvider, acpauth.StaticTokenAuthorizer("test-key"))

	body := []byte(`{"capabilities":{},"currency":"USD","line_items":[{"id":"sku_1"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
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

func TestCheckoutHandler(t *testing.T) {
	t.Parallel()

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)

		h := acpcheckout.NewCheckoutHandler(NewMockCheckoutProvider(ctrl), acpauth.StaticTokenAuthorizer("test-key"))

		body := []byte(`{"capabilities":{},"currency":"USD","line_items":[{"id":"sku_1"}]}`)
		req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("API-Version", "2026-01-30")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 got %d body=%s", rec.Code, rec.Body.String())
		}

		var got acpcheckout.Error
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

		h := acpcheckout.NewCheckoutHandler(NewMockCheckoutProvider(ctrl), acpauth.StaticTokenAuthorizer("test-key"))

		body := []byte(`{"capabilities":{},"currency":"USD","line_items":[]}`)
		req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-key")
		req.Header.Set("API-Version", "2026-01-30")
		req.Header.Set("Idempotency-Key", "idem-1")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422 got %d body=%s", rec.Code, rec.Body.String())
		}

		var got acpcheckout.Error
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if got.Param == nil || *got.Param != "$.line_items" {
			t.Fatalf("unexpected param: %#v", got.Param)
		}
	})
}
