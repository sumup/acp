package acpauthentication

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
)

type stubProvider struct {
	create func(context.Context, DelegateAuthenticationCreateRequest) (*DelegateAuthenticationSession, error)
}

func (s stubProvider) CreateAuthenticationSession(ctx context.Context, req DelegateAuthenticationCreateRequest) (*DelegateAuthenticationSession, error) {
	return s.create(ctx, req)
}

func (stubProvider) AuthenticateSession(context.Context, string, DelegateAuthenticationAuthenticateRequest) (*DelegateAuthenticationSession, error) {
	panic("unexpected AuthenticateSession call")
}

func (stubProvider) GetAuthenticationSession(context.Context, string) (*DelegateAuthenticationSessionWithResult, error) {
	panic("unexpected GetAuthenticationSession call")
}

func TestHandlerCreateAuthenticationSession(t *testing.T) {
	t.Parallel()

	provider := stubProvider{create: func(ctx context.Context, req DelegateAuthenticationCreateRequest) (*DelegateAuthenticationSession, error) {
		if acp.RequestContextFromContext(ctx) == nil {
			t.Fatal("request context missing")
		}
		if req.MerchantId != "merchant_123" || req.Amount.Value != 1000 {
			t.Fatalf("unexpected request: %#v", req)
		}
		return &DelegateAuthenticationSessionBase{
			AuthenticationSessionId: "auth_1",
			Status:                  DelegateAuthenticationSessionBaseStatusPending,
		}, nil
	}}
	handler := NewHandler(provider, acpauth.StaticTokenAuthorizer("test-key"))

	body := `{
		"merchant_id":"merchant_123",
		"payment_method":{"type":"card","number":"4917610000000000","exp_month":"03","exp_year":"2030","name":"Jane Doe"},
		"amount":{"value":1000,"currency":"EUR"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/delegate_authentication", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("API-Version", acp.APIVersion)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
