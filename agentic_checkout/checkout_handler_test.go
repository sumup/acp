package agentic_checkout

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type checkoutStub struct{}

func (checkoutStub) CreateSession(context.Context, CheckoutSessionCreateRequest) (*CheckoutSessionBase, error) {
	return &CheckoutSessionBase{Id: "cs_1", Currency: "USD", Status: CheckoutSessionBaseStatusInProgress, FulfillmentOptions: []CheckoutSessionBase_FulfillmentOptions_Item{}, LineItems: []LineItem{}, Links: []Link{}, Messages: []CheckoutSessionBase_Messages_Item{}}, nil
}
func (checkoutStub) UpdateSession(context.Context, string, CheckoutSessionUpdateRequest) (*CheckoutSessionBase, error) {
	return &CheckoutSessionBase{Id: "cs_1", Currency: "USD", Status: CheckoutSessionBaseStatusInProgress, FulfillmentOptions: []CheckoutSessionBase_FulfillmentOptions_Item{}, LineItems: []LineItem{}, Links: []Link{}, Messages: []CheckoutSessionBase_Messages_Item{}}, nil
}
func (checkoutStub) GetSession(context.Context, string) (*CheckoutSessionBase, error) {
	return &CheckoutSessionBase{Id: "cs_1", Currency: "USD", Status: CheckoutSessionBaseStatusInProgress, FulfillmentOptions: []CheckoutSessionBase_FulfillmentOptions_Item{}, LineItems: []LineItem{}, Links: []Link{}, Messages: []CheckoutSessionBase_Messages_Item{}}, nil
}
func (checkoutStub) CompleteSession(context.Context, string, CheckoutSessionCompleteRequest) (CheckoutSessionWithOrder, error) {
	return CheckoutSessionWithOrder{
		Id:                 "cs_1",
		Currency:           "USD",
		Capabilities:       Capabilities{},
		FulfillmentOptions: []CheckoutSessionWithOrder_FulfillmentOptions_Item{},
		LineItems:          []LineItem{},
		Links:              []Link{},
		Messages:           []CheckoutSessionWithOrder_Messages_Item{},
		Status:             CheckoutSessionWithOrderStatusInProgress,
		Totals:             []Total{},
		Order: Order{
			Id:                "ord_1",
			CheckoutSessionId: "cs_1",
			PermalinkUrl:      "https://example.com/orders/ord_1",
		},
	}, nil
}
func (checkoutStub) CancelSession(context.Context, string, *CancelSessionRequest) (*CheckoutSessionBase, error) {
	return &CheckoutSessionBase{Id: "cs_1", Currency: "USD", Status: CheckoutSessionBaseStatusCanceled, FulfillmentOptions: []CheckoutSessionBase_FulfillmentOptions_Item{}, LineItems: []LineItem{}, Links: []Link{}, Messages: []CheckoutSessionBase_Messages_Item{}}, nil
}

func TestCheckoutHandlerCreate(t *testing.T) {
	h := NewCheckoutHandler(checkoutStub{})
	body := []byte(`{"capabilities":{},"currency":"USD","line_items":[{"id":"sku_1"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/checkout_sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-1")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rec.Code)
	}
}
