package acpauthentication_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpauthentication"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	t.Run("authentication lifecycle", func(t *testing.T) {
		t.Parallel()

		provider := NewMockProvider(gomock.NewController(t))
		session := &acpauthentication.DelegateAuthenticationSession{
			AuthenticationSessionId: "auth_1",
			Status:                  acpauthentication.DelegateAuthenticationSessionBaseStatusPending,
		}
		result := &acpauthentication.DelegateAuthenticationSessionWithResult{
			AuthenticationSessionId: session.AuthenticationSessionId,
			Status:                  acpauthentication.DelegateAuthenticationSessionWithResultStatusAuthenticated,
			AuthenticationResult: &acpauthentication.AuthenticationResult{
				TransStatus: "Y", TransactionId: "ds_123", ThreeDsServerTransId: "server_123", Version: "2.2.0",
				ThreeDsCryptogram: new("cryptogram"), ElectronicCommerceIndicator: new("05"),
			},
		}
		gomock.InOrder(
			provider.EXPECT().CreateAuthenticationSession(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, req acpauthentication.DelegateAuthenticationCreateRequest) (*acpauthentication.DelegateAuthenticationSession, error) {
					if acp.RequestContextFromContext(ctx) == nil {
						t.Fatal("request context missing")
					}
					if req.MerchantId != "merchant_123" || req.Amount.Value != 1000 || req.PaymentMethod.Number != "4917610000000000" {
						t.Fatalf("unexpected request: %#v", req)
					}
					return session, nil
				}),
			provider.EXPECT().AuthenticateSession(gomock.Any(), session.AuthenticationSessionId,
				acpauthentication.DelegateAuthenticationAuthenticateRequest{FingerprintCompletion: acpauthentication.Y}).
				Return(&acpauthentication.DelegateAuthenticationSession{
					AuthenticationSessionId: session.AuthenticationSessionId,
					Status:                  acpauthentication.DelegateAuthenticationSessionBaseStatusAuthenticated,
				}, nil),
			provider.EXPECT().GetAuthenticationSession(gomock.Any(), session.AuthenticationSessionId).Return(result, nil),
		)
		handler := acpauthentication.NewHandler(provider, acpauth.StaticTokenAuthorizer("test-key"))
		for _, step := range []struct {
			method string
			path   string
			body   string
			status int
			state  string
		}{
			{http.MethodPost, "/delegate_authentication", `{"merchant_id":"merchant_123","payment_method":{"type":"card","number":"4917610000000000","exp_month":"03","exp_year":"2030","name":"Jane Doe"},"amount":{"value":1000,"currency":"EUR"}}`, http.StatusCreated, "pending"},
			{http.MethodPost, "/delegate_authentication/" + session.AuthenticationSessionId + "/authenticate", `{"fingerprint_completion":"Y"}`, http.StatusOK, "authenticated"},
			{http.MethodGet, "/delegate_authentication/" + session.AuthenticationSessionId, "", http.StatusOK, "authenticated"},
		} {
			req := newRequest(t, step.method, step.path, step.body)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != step.status {
				t.Fatalf("%s %s: status = %d, body = %s", step.method, step.path, rec.Code, rec.Body.String())
			}
			var got acpauthentication.DelegateAuthenticationSessionWithResult
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got.AuthenticationSessionId != session.AuthenticationSessionId || string(got.Status) != step.state {
				t.Fatalf("unexpected session: %#v", got)
			}
			if step.method == http.MethodGet && (got.AuthenticationResult == nil || got.AuthenticationResult.ThreeDsCryptogram == nil || *got.AuthenticationResult.ThreeDsCryptogram != "cryptogram") {
				t.Fatalf("unexpected authentication result: %#v", got.AuthenticationResult)
			}
			if rec.Header().Get("Request-Id") != "req-1" {
				t.Fatalf("headers = %v", rec.Header())
			}
		}
	})

	t.Run("invalid fingerprint does not reach the provider", func(t *testing.T) {
		t.Parallel()

		provider := NewMockProvider(gomock.NewController(t))
		handler := acpauthentication.NewHandler(provider, acpauth.StaticTokenAuthorizer("test-key"))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, newRequest(t, http.MethodPost, "/delegate_authentication/auth_1/authenticate", `{"fingerprint_completion":"invalid"}`))
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
	})
}

func newRequest(t *testing.T, method, path, body string) *http.Request {
	t.Helper()

	req := httptest.NewRequestWithContext(t.Context(), method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("API-Version", acp.APIVersion)
	req.Header.Set("Request-Id", "req-1")
	return req
}
