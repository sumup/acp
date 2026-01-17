package acp

import "testing"

func TestCheckoutSessionCompleteRequestValidateAuthenticationResult(t *testing.T) {
	t.Parallel()

	base := CheckoutSessionCompleteRequest{
		PaymentData: PaymentData{
			Token:    "tok",
			Provider: "stripe",
		},
	}

	t.Run("no authentication result", func(t *testing.T) {
		t.Parallel()
		if err := base.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("requires outcome details", func(t *testing.T) {
		t.Parallel()
		req := base
		req.AuthenticationResult = &AuthenticationResult{
			Outcome: AuthenticationOutcomeAuthenticated,
		}
		if err := req.Validate(); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("missing outcome details field", func(t *testing.T) {
		t.Parallel()
		req := base
		req.AuthenticationResult = &AuthenticationResult{
			Outcome: AuthenticationOutcomeAuthenticated,
			OutcomeDetails: &AuthenticationOutcomeDetails{
				ElectronicCommerceIndicator: AuthenticationECI05,
				TransactionId:               "dsTransId",
				Version:                     "2.2.0",
			},
		}
		if err := req.Validate(); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("valid outcome details", func(t *testing.T) {
		t.Parallel()
		req := base
		req.AuthenticationResult = &AuthenticationResult{
			Outcome: AuthenticationOutcomeAuthenticated,
			OutcomeDetails: &AuthenticationOutcomeDetails{
				ThreeDsCryptogram:           "AbCdEfGhIjKlMnOpQrStUvWxY0=",
				ElectronicCommerceIndicator: AuthenticationECI05,
				TransactionId:               "dsTransId",
				Version:                     "2.2.0",
			},
		}
		if err := req.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
