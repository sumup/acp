package acpauthentication_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpauthentication"
)

type authenticationProvider struct{}

func (authenticationProvider) CreateAuthenticationSession(context.Context, acpauthentication.DelegateAuthenticationCreateRequest) (*acpauthentication.DelegateAuthenticationSession, error) {
	return nil, errors.ErrUnsupported
}

func (authenticationProvider) AuthenticateSession(context.Context, string, acpauthentication.DelegateAuthenticationAuthenticateRequest) (*acpauthentication.DelegateAuthenticationSession, error) {
	return nil, errors.ErrUnsupported
}

func (authenticationProvider) GetAuthenticationSession(context.Context, string) (*acpauthentication.DelegateAuthenticationSessionWithResult, error) {
	return nil, errors.ErrUnsupported
}

func ExampleNewHandler() {
	handler := acpauthentication.NewHandler(
		authenticationProvider{},
		acpauth.StaticTokenAuthorizer("api_key_123"),
	)

	fmt.Printf("%T\n", handler)
	// Output: *acpauthentication.Handler
}
