package acp

import (
	"context"
	"errors"
	"net/http"
	"strings"
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

func (h *DelegatedPaymentHandler) authenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.cfg.authenticator == nil {
			next(w, r)
			return
		}
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader == "" {
			writeJSONError(w, NewHTTPError(http.StatusUnauthorized, InvalidRequest, MissingAuthorization, "Authorization header is required"))
			return
		}
		schema, apiKey, ok := strings.Cut(authHeader, " ")
		if !ok || !strings.EqualFold(schema, "Bearer") {
			writeJSONError(w, NewHTTPError(http.StatusUnauthorized, InvalidRequest, InvalidAuthorization, "Authorization header must be in the format 'Bearer <api_key>'"))
			return
		}
		if apiKey == "" {
			writeJSONError(w, NewHTTPError(http.StatusUnauthorized, InvalidRequest, InvalidAuthorization, "API key is required"))
			return
		}
		if err := h.cfg.authenticator.Authenticate(r.Context(), apiKey); err != nil {
			var httpErr *Error
			if errors.As(err, &httpErr) {
				writeJSONError(w, httpErr)
				return
			}
			writeJSONError(w, NewHTTPError(http.StatusUnauthorized, InvalidRequest, InvalidAuthorization, "invalid API key"))
			return
		}
		next(w, r)
	}
}
