package acpauth_test

import (
	"context"
	"testing"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
)

func TestStaticTokenAuthorizer(t *testing.T) {
	t.Parallel()

	authorizer := acpauth.StaticTokenAuthorizer("test-key")

	if err := authorizer.Authorize(context.Background(), "test-key"); err != nil {
		t.Fatalf("Authorize() error = %v", err)
	}

	err := authorizer.Authorize(context.Background(), "wrong-key")
	if err == nil {
		t.Fatal("Authorize() error = nil, want error")
	}

	acpErr, ok := err.(*acp.Error)
	if !ok {
		t.Fatalf("Authorize() error type = %T, want *acp.Error", err)
	}
	if acpErr.Code != acp.InvalidAuthorization {
		t.Fatalf("Authorize() code = %q, want %q", acpErr.Code, acp.InvalidAuthorization)
	}
}
