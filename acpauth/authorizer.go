package acpauth

import (
	"context"

	"github.com/sumup/acp"
)

// Authorizer validates a bearer token extracted from the Authorization header.
// The token does not include the "Bearer" scheme prefix.
type Authorizer interface {
	Authorize(ctx context.Context, bearer string) error
}

// AuthorizerFunc lifts bare functions into [Authorizer].
type AuthorizerFunc func(ctx context.Context, bearer string) error

// Authorize validates bearer using the wrapped function.
func (f AuthorizerFunc) Authorize(ctx context.Context, bearer string) error {
	return f(ctx, bearer)
}

// StaticTokenAuthorizer returns an [Authorizer] that accepts exactly token.
func StaticTokenAuthorizer(token string) Authorizer {
	return AuthorizerFunc(func(_ context.Context, bearer string) error {
		if bearer != token {
			return acp.NewHTTPError(
				401,
				acp.InvalidRequest,
				acp.InvalidAuthorization,
				"Authorization bearer token is invalid",
			)
		}
		return nil
	})
}
