package delegate_payment

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type delegatedStub struct{}

func (delegatedStub) DelegatePayment(context.Context, DelegatePaymentRequest) (*DelegatePaymentResponse, error) {
	return &DelegatePaymentResponse{Id: "vt_1", Created: time.Now().UTC(), Metadata: map[string]string{"source": "test"}}, nil
}

func TestDelegatedPaymentHandlerCreate(t *testing.T) {
	h := NewDelegatedPaymentHandler(delegatedStub{})
	body := []byte(`{"payment_method":{"type":"card","card_number_type":"fpan","number":"4242424242424242","display_card_funding_type":"credit","metadata":{}},"allowance":{"reason":"one_time","max_amount":1000,"currency":"usd","checkout_session_id":"cs_1","merchant_id":"m_1","expires_at":"2026-12-31T00:00:00Z"},"risk_signals":[{"type":"card_testing","action":"authorized","score":0}],"metadata":{}}`)
	req := httptest.NewRequest(http.MethodPost, "/agentic_commerce/delegate_payment", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-1")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rec.Code, rec.Body.String())
	}
}
