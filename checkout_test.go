package acp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckoutHandlerRoutes(t *testing.T) {
	t.Parallel()

	session := &CheckoutSession{
		Id:                 "cs_123",
		Status:             CheckoutSessionBaseStatusInProgress,
		Currency:           "USD",
		LineItems:          []LineItem{},
		FulfillmentOptions: make([]FulfillmentOption, 0),
		Totals:             []Total{},
		Messages:           make([]Message, 0),
		Links:              []Link{},
	}

	orderSession := &CheckoutSessionWithOrder{
		Id:                 session.Id,
		Status:             CheckoutSessionWithOrderStatusInProgress,
		Currency:           session.Currency,
		LineItems:          session.LineItems,
		Totals:             session.Totals,
		Links:              session.Links,
		FulfillmentOptions: make([]FulfillmentOption, 0),
		Messages:           make([]Message, 0),
		Order: Order{
			Id:                "ord_123",
			CheckoutSessionId: "cs_123",
			PermalinkUrl:      "https://example.com/orders/123",
		},
	}

	tests := map[string]struct {
		method     string
		path       string
		body       any
		setupStub  func(*stubService)
		wantStatus int
	}{
		"create session": {
			method: http.MethodPost,
			path:   "/checkout_sessions",
			body: CheckoutSessionCreateRequest{
				Items: []Item{{Id: "sku_1", Quantity: 1}},
			},
			setupStub: func(s *stubService) {
				s.create = func(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSession, error) {
					if len(req.Items) != 1 {
						t.Fatalf("expected 1 item")
					}
					return session, nil
				}
			},
			wantStatus: http.StatusCreated,
		},
		"get session": {
			method: http.MethodGet,
			path:   "/checkout_sessions/cs_123",
			setupStub: func(s *stubService) {
				s.get = func(ctx context.Context, id string) (*CheckoutSession, error) {
					if id != "cs_123" {
						t.Fatalf("unexpected id %s", id)
					}
					return session, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		"update session": {
			method: http.MethodPost,
			path:   "/checkout_sessions/cs_123",
			body: CheckoutSessionUpdateRequest{
				Items: &[]Item{{Id: "sku_1", Quantity: 2}},
			},
			setupStub: func(s *stubService) {
				s.update = func(ctx context.Context, id string, req CheckoutSessionUpdateRequest) (*CheckoutSession, error) {
					if id != "cs_123" {
						t.Fatalf("unexpected id %s", id)
					}
					return session, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		"complete session": {
			method: http.MethodPost,
			path:   "/checkout_sessions/cs_123/complete",
			body: CheckoutSessionCompleteRequest{
				PaymentData: PaymentData{Token: "tok", Provider: "sumup"},
			},
			setupStub: func(s *stubService) {
				s.complete = func(ctx context.Context, id string, req CheckoutSessionCompleteRequest) (*CheckoutSessionWithOrder, error) {
					return orderSession, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		"cancel session": {
			method: http.MethodPost,
			path:   "/checkout_sessions/cs_123/cancel",
			setupStub: func(s *stubService) {
				s.cancel = func(ctx context.Context, id string) (*CheckoutSession, error) {
					return session, nil
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stub := &stubService{}
			if tt.setupStub != nil {
				tt.setupStub(stub)
			}
			handler := NewCheckoutHandler(stub)
			var bodyReader *bytes.Reader
			if tt.body != nil {
				payload, err := json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("marshal body: %v", err)
				}
				bodyReader = bytes.NewReader(payload)
			} else {
				bodyReader = bytes.NewReader(nil)
			}
			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d got %d, body=%s", tt.wantStatus, rec.Code, rec.Body.String())
			}
			if got := rec.Header().Get("API-Version"); got != APIVersion {
				t.Fatalf("missing API-Version header")
			}
		})
	}
}

func TestCheckoutHandlerErrors(t *testing.T) {
	t.Parallel()

	t.Run("invalid JSON", func(t *testing.T) {
		handler := NewCheckoutHandler(&stubService{
			create: func(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSession, error) {
				return &CheckoutSession{}, nil
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", strings.NewReader("{"))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 got %d", rec.Code)
		}
	})

	t.Run("service error surfaces", func(t *testing.T) {
		handler := NewCheckoutHandler(&stubService{
			get: func(ctx context.Context, id string) (*CheckoutSession, error) {
				return nil, NewHTTPError(http.StatusNotFound, InvalidRequest, ErrorCode("not_found"), "missing")
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/checkout_sessions/unknown", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 got %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "missing") {
			t.Fatalf("unexpected body %s", rec.Body.String())
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		handler := NewCheckoutHandler(&stubService{})
		req := httptest.NewRequest(http.MethodGet, "/checkout_sessions", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405 got %d", rec.Code)
		}
	})
}

type stubService struct {
	create   func(context.Context, CheckoutSessionCreateRequest) (*CheckoutSession, error)
	update   func(context.Context, string, CheckoutSessionUpdateRequest) (*CheckoutSession, error)
	get      func(context.Context, string) (*CheckoutSession, error)
	complete func(context.Context, string, CheckoutSessionCompleteRequest) (*CheckoutSessionWithOrder, error)
	cancel   func(context.Context, string) (*CheckoutSession, error)
}

func (s *stubService) CreateSession(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSession, error) {
	if s.create != nil {
		return s.create(ctx, req)
	}
	return nil, NewHTTPError(http.StatusNotImplemented, InvalidRequest, ErrorCode("not_implemented"), "create not implemented")
}

func (s *stubService) UpdateSession(ctx context.Context, id string, req CheckoutSessionUpdateRequest) (*CheckoutSession, error) {
	if s.update != nil {
		return s.update(ctx, id, req)
	}
	return nil, NewHTTPError(http.StatusNotImplemented, InvalidRequest, ErrorCode("not_implemented"), "update not implemented")
}

func (s *stubService) GetSession(ctx context.Context, id string) (*CheckoutSession, error) {
	if s.get != nil {
		return s.get(ctx, id)
	}
	return nil, NewHTTPError(http.StatusNotImplemented, InvalidRequest, ErrorCode("not_implemented"), "get not implemented")
}

func (s *stubService) CompleteSession(ctx context.Context, id string, req CheckoutSessionCompleteRequest) (*CheckoutSessionWithOrder, error) {
	if s.complete != nil {
		return s.complete(ctx, id, req)
	}
	return nil, NewHTTPError(http.StatusNotImplemented, InvalidRequest, ErrorCode("not_implemented"), "complete not implemented")
}

func (s *stubService) CancelSession(ctx context.Context, id string) (*CheckoutSession, error) {
	if s.cancel != nil {
		return s.cancel(ctx, id)
	}
	return nil, NewHTTPError(http.StatusNotImplemented, InvalidRequest, ErrorCode("not_implemented"), "cancel not implemented")
}
