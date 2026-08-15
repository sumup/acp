package acpcart

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
)

type stubProvider struct {
	create func(context.Context, CartCreateRequest) (*Cart, error)
}

func (s stubProvider) CreateCart(ctx context.Context, req CartCreateRequest) (*Cart, error) {
	return s.create(ctx, req)
}

func (stubProvider) GetCart(context.Context, string) (*Cart, error) {
	panic("unexpected GetCart call")
}

func (stubProvider) UpdateCart(context.Context, string, CartUpdateRequest) (*Cart, error) {
	panic("unexpected UpdateCart call")
}

func (stubProvider) CancelCart(context.Context, string) (*Cart, error) {
	panic("unexpected CancelCart call")
}

func TestHandlerCreateCart(t *testing.T) {
	t.Parallel()

	provider := stubProvider{create: func(ctx context.Context, req CartCreateRequest) (*Cart, error) {
		requestCtx := acp.RequestContextFromContext(ctx)
		if requestCtx == nil || requestCtx.IdempotencyKey != "idem-1" {
			t.Fatalf("unexpected request context: %#v", requestCtx)
		}
		if len(req.LineItems) != 1 || req.LineItems[0].Id != "item_123" {
			t.Fatalf("unexpected request: %#v", req)
		}
		return &Cart{Id: "cart_1", Currency: "usd", LineItems: []LineItem{}, Totals: []Total{}}, nil
	}}
	handler := NewHandler(provider, acpauth.StaticTokenAuthorizer("test-key"))

	req := httptest.NewRequest(http.MethodPost, "/carts", bytes.NewBufferString(`{"line_items":[{"id":"item_123"}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("API-Version", acp.APIVersion)
	req.Header.Set("Idempotency-Key", "idem-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Idempotency-Key"); got != "idem-1" {
		t.Fatalf("Idempotency-Key = %q", got)
	}
}

func TestHandlerRequiresIdempotencyKey(t *testing.T) {
	t.Parallel()

	handler := NewHandler(stubProvider{}, acpauth.StaticTokenAuthorizer("test-key"))
	req := httptest.NewRequest(http.MethodPost, "/carts", bytes.NewBufferString(`{"line_items":[{"id":"item_123"}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("API-Version", acp.APIVersion)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var got acp.Error
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Code != acp.IdempotencyKeyRequired {
		t.Fatalf("code = %q", got.Code)
	}
}
