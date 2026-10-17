package acpcart_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpcart"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	t.Run("cart lifecycle", func(t *testing.T) {
		t.Parallel()

		provider := NewMockProvider(gomock.NewController(t))
		cart := &acpcart.Cart{Id: "cart_1", Currency: "usd", LineItems: []acpcart.LineItem{}, Totals: []acpcart.Total{}}
		gomock.InOrder(
			provider.EXPECT().CreateCart(gomock.Any(), acpcart.CartCreateRequest{
				LineItems: []acpcart.Item{{Id: "item_123"}},
			}).DoAndReturn(func(ctx context.Context, _ acpcart.CartCreateRequest) (*acpcart.Cart, error) {
				metadata := acp.RequestContextFromContext(ctx)
				if metadata == nil || metadata.IdempotencyKey != "idem-1" || metadata.RequestID != "req-1" {
					t.Fatalf("unexpected request context: %#v", metadata)
				}
				if ctx.Value(contextKey{}) != "test" {
					t.Fatal("request context was not propagated")
				}
				return cart, nil
			}),
			provider.EXPECT().UpdateCart(gomock.Any(), cart.Id, acpcart.CartUpdateRequest{
				LineItems: []acpcart.Item{{Id: "item_456"}},
			}).Return(cart, nil),
			provider.EXPECT().GetCart(gomock.Any(), cart.Id).Return(cart, nil),
			provider.EXPECT().CancelCart(gomock.Any(), cart.Id).Return(cart, nil),
		)
		handler := acpcart.NewHandler(provider, acpauth.StaticTokenAuthorizer("test-key"))
		for _, step := range []struct {
			method string
			path   string
			body   string
			status int
		}{
			{http.MethodPost, "/carts", `{"line_items":[{"id":"item_123"}]}`, http.StatusCreated},
			{http.MethodPut, "/carts/" + cart.Id, `{"line_items":[{"id":"item_456"}]}`, http.StatusOK},
			{http.MethodGet, "/carts/" + cart.Id, "", http.StatusOK},
			{http.MethodPost, "/carts/" + cart.Id + "/cancel", "", http.StatusOK},
		} {
			req := newRequest(t, step.method, step.path, step.body)
			if step.method == http.MethodGet {
				req.Header.Del("Idempotency-Key")
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != step.status {
				t.Fatalf("%s %s: status = %d, body = %s", step.method, step.path, rec.Code, rec.Body.String())
			}
			var got acpcart.Cart
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got.Id != cart.Id || got.Currency != cart.Currency {
				t.Fatalf("unexpected cart: %#v", got)
			}
			if rec.Header().Get("API-Version") != acp.APIVersion || rec.Header().Get("Request-Id") != "req-1" {
				t.Fatalf("response headers = %v", rec.Header())
			}
			if got := rec.Header().Get("Idempotency-Key"); got != req.Header.Get("Idempotency-Key") {
				t.Fatalf("Idempotency-Key = %q", got)
			}
		}
	})

	t.Run("invalid requests do not reach the provider", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			name   string
			body   string
			header string
			status int
			code   acp.ErrorCode
		}{
			{"missing authorization", `{"line_items":[{"id":"item_123"}]}`, "Authorization", http.StatusBadRequest, acp.MissingAuthorization},
			{"missing API version", `{"line_items":[{"id":"item_123"}]}`, "API-Version", http.StatusBadRequest, acp.MissingAPIVersion},
			{"missing idempotency key", `{"line_items":[{"id":"item_123"}]}`, "Idempotency-Key", http.StatusBadRequest, acp.IdempotencyKeyRequired},
			{"missing items", `{}`, "", http.StatusBadRequest, acp.ErrorCode(acp.InvalidRequest)},
			{"empty items", `{"line_items":[]}`, "", http.StatusUnprocessableEntity, acp.ErrorCode(acp.InvalidRequest)},
			{"invalid JSON", `{`, "", http.StatusBadRequest, acp.ErrorCode(acp.InvalidRequest)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				handler := acpcart.NewHandler(NewMockProvider(gomock.NewController(t)), acpauth.StaticTokenAuthorizer("test-key"))
				req := newRequest(t, http.MethodPost, "/carts", tc.body)
				if tc.header != "" {
					req.Header.Del(tc.header)
				}
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)
				if rec.Code != tc.status {
					t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
				}
				var got acp.Error
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.Code != tc.code {
					t.Fatalf("code = %q, want %q", got.Code, tc.code)
				}
			})
		}
	})

	t.Run("wrapped provider errors preserve protocol status", func(t *testing.T) {
		t.Parallel()

		provider := NewMockProvider(gomock.NewController(t))
		provider.EXPECT().GetCart(gomock.Any(), "missing").Return(nil, fmt.Errorf("get cart: %w", acp.NewHTTPError(http.StatusNotFound, acp.InvalidRequest, acp.ErrorCode("cart_not_found"), "cart not found")))
		handler := acpcart.NewHandler(provider, acpauth.StaticTokenAuthorizer("test-key"))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, newRequest(t, http.MethodGet, "/carts/missing", ""))
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
	})
}

type contextKey struct{}

func newRequest(t *testing.T, method, path, body string) *http.Request {
	t.Helper()

	ctx := context.WithValue(t.Context(), contextKey{}, "test")
	req := httptest.NewRequestWithContext(ctx, method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("API-Version", acp.APIVersion)
	req.Header.Set("Idempotency-Key", "idem-1")
	req.Header.Set("Request-Id", "req-1")
	return req
}
