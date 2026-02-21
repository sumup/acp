package acp_test

import (
	"context"
	"testing"

	"github.com/sumup/acp"
)

func TestAuthenticatorFunc(t *testing.T) {
	fn := acp.AuthenticatorFunc(func(_ context.Context, _ string) error { return nil })
	if fn == nil {
		t.Fatalf("expected authenticator func")
	}
}
