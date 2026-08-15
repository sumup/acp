package acppayment_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acppayment"
)

type paymentProvider struct{}

func (paymentProvider) DelegatePayment(context.Context, acppayment.DelegatePaymentRequest) (*acppayment.DelegatePaymentResponse, error) {
	return nil, errors.ErrUnsupported
}

func ExampleNewHandler() {
	handler := acppayment.NewHandler(
		paymentProvider{},
		acpauth.StaticTokenAuthorizer("api_key_123"),
	)

	fmt.Printf("%T\n", handler)
	// Output: *acppayment.Handler
}
