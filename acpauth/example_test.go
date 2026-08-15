package acpauth_test

import (
	"context"
	"fmt"

	"github.com/sumup/acp/acpauth"
)

func ExampleStaticTokenAuthorizer() {
	authorizer := acpauth.StaticTokenAuthorizer("api_key_123")

	err := authorizer.Authorize(context.Background(), "api_key_123")
	fmt.Println(err)
	// Output: <nil>
}
