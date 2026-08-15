package acpcheckout_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpcheckout"
)

type checkoutProvider struct{}

func (checkoutProvider) CreateSession(context.Context, acpcheckout.CheckoutSessionCreateRequest) (*acpcheckout.CheckoutSessionBase, error) {
	return nil, errors.ErrUnsupported
}

func (checkoutProvider) UpdateSession(context.Context, string, acpcheckout.CheckoutSessionUpdateRequest) (*acpcheckout.CheckoutSessionBase, error) {
	return nil, errors.ErrUnsupported
}

func (checkoutProvider) GetSession(context.Context, string) (*acpcheckout.CheckoutSessionBase, error) {
	return nil, errors.ErrUnsupported
}

func (checkoutProvider) CompleteSession(context.Context, string, acpcheckout.CheckoutSessionCompleteRequest) (acpcheckout.CheckoutSessionWithOrder, error) {
	return acpcheckout.CheckoutSessionWithOrder{}, errors.ErrUnsupported
}

func (checkoutProvider) CancelSession(context.Context, string, *acpcheckout.CancelSessionRequest) (*acpcheckout.CheckoutSessionBase, error) {
	return nil, errors.ErrUnsupported
}

func ExampleNewHandler() {
	handler := acpcheckout.NewHandler(
		checkoutProvider{},
		acpauth.StaticTokenAuthorizer("api_key_123"),
	)

	fmt.Printf("%T\n", handler)
	// Output: *acpcheckout.Handler
}
