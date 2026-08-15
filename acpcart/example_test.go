package acpcart_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpcart"
)

type cartProvider struct{}

func (cartProvider) CreateCart(context.Context, acpcart.CartCreateRequest) (*acpcart.Cart, error) {
	return nil, errors.ErrUnsupported
}

func (cartProvider) GetCart(context.Context, string) (*acpcart.Cart, error) {
	return nil, errors.ErrUnsupported
}

func (cartProvider) UpdateCart(context.Context, string, acpcart.CartUpdateRequest) (*acpcart.Cart, error) {
	return nil, errors.ErrUnsupported
}

func (cartProvider) CancelCart(context.Context, string) (*acpcart.Cart, error) {
	return nil, errors.ErrUnsupported
}

func ExampleNewHandler() {
	handler := acpcart.NewHandler(
		cartProvider{},
		acpauth.StaticTokenAuthorizer("api_key_123"),
	)

	fmt.Printf("%T\n", handler)
	// Output: *acpcart.Handler
}
