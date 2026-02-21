package agentic_checkout

import (
	"testing"
)

func TestCheckoutValidateRequiresPaymentData(t *testing.T) {
	req := CheckoutSessionCompleteRequest{}
	if err := req.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}
