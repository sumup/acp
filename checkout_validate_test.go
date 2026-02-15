package acp

import "testing"

func TestCheckoutSessionCompleteRequestValidateAuthenticationResult(t *testing.T) {
	t.Parallel()

	base := CheckoutSessionCompleteRequest{
		PaymentData: PaymentData{
			HandlerID: ptr("handler_sumup_card"),
			Instrument: &PaymentInstrument{
				Type: "card",
				Credential: PaymentCredential{
					Type:  "spt",
					Token: "spt_123",
				},
			},
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

	t.Run("legacy token provider still accepted", func(t *testing.T) {
		t.Parallel()
		req := CheckoutSessionCompleteRequest{
			PaymentData: PaymentData{
				Token:    ptr("tok"),
				Provider: ptr(PaymentDataProvider("sumup")),
			},
		}
		if err := req.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("requires handler_id with instrument", func(t *testing.T) {
		t.Parallel()
		req := CheckoutSessionCompleteRequest{
			PaymentData: PaymentData{
				Instrument: &PaymentInstrument{
					Type: "card",
					Credential: PaymentCredential{
						Type:  "spt",
						Token: "spt_123",
					},
				},
			},
		}
		if err := req.Validate(); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("requires instrument credential token", func(t *testing.T) {
		t.Parallel()
		req := CheckoutSessionCompleteRequest{
			PaymentData: PaymentData{
				HandlerID: ptr("sumup"),
				Instrument: &PaymentInstrument{
					Type: "card",
					Credential: PaymentCredential{
						Type: "spt",
					},
				},
			},
		}
		if err := req.Validate(); err == nil {
			t.Fatalf("expected error")
		}
	})
}

// ptr returns to pointer of v.
// Can be replaced with `new` with Go 1.26+.
func ptr[V any](v V) *V {
	return &v
}
