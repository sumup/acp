package acp

import (
	"context"
)

// Authenticator validates Authorization header API keys before the
// request reaches the provider.
type Authenticator interface {
	Authenticate(ctx context.Context, apiKey string) error
}

// AuthenticatorFunc lifts bare functions into [Authenticator].
type AuthenticatorFunc func(ctx context.Context, apiKey string) error

// Authenticate validates the API key using the wrapped function.
func (f AuthenticatorFunc) Authenticate(ctx context.Context, apiKey string) error {
	return f(ctx, apiKey)
}
