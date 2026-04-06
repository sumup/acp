package srv

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
)

// HandlerFunc handles an HTTP request and returns any error that should be
// processed by an [ErrorHandler].
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ErrorHandler handles an error returned from a [HandlerFunc].
type ErrorHandler func(http.ResponseWriter, *http.Request, error)

// Middleware wraps a [HandlerFunc] with additional request processing.
type Middleware func(http.ResponseWriter, *http.Request, HandlerFunc) error

// Mux wraps [http.ServeMux] and routes handler errors through a shared
// [ErrorHandler].
type Mux struct {
	*http.ServeMux
	handleError ErrorHandler
}

// NewMux returns a [Mux] that dispatches handler errors to handleError.
//
// It panics if handleError is nil.
func NewMux(handleError ErrorHandler) *Mux {
	return &Mux{
		ServeMux:    http.NewServeMux(),
		handleError: handleError,
	}
}

// HandleFunc registers pattern with a handler that may return an error.
//
// If the handler returns a non-nil error, the mux forwards it to the
// [ErrorHandler] provided to [NewMux]. Middleware is applied in the order it is
// provided.
func (m *Mux) HandleFunc(pattern string, f HandlerFunc, middleware ...Middleware) {
	for i := len(middleware) - 1; i >= 0; i-- {
		next := f
		current := middleware[i]
		f = func(w http.ResponseWriter, r *http.Request) error {
			return current(w, r, next)
		}
	}

	m.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			m.handleError(w, r, err)
		}
	})
}

// WriteACPError writes err as a JSON response using build to construct the
// package-specific error payload.
//
// The build function receives err's type, code, message, and param values as
// strings that can be mapped into the caller's error model.
func WriteACPError[T any](w http.ResponseWriter, err *acp.Error, build func(errorType, code, message string, param *string) T) error {
	return WriteError(w, err.StatusCode(), build(string(err.Type), string(err.Code), err.Message, err.Param))
}

// WriteError writes payload as JSON with ACP error response headers and status.
func WriteError(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("API-Version", acp.APIVersion)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

// WriteJSON writes payload as JSON with ACP response headers and status.
//
// If the request includes an Idempotency-Key header, WriteJSON echoes it back
// in the response. A nil payload writes headers and status without a body.
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, payload any) error {
	w.Header().Set("API-Version", acp.APIVersion)
	w.Header().Set("Content-Type", "application/json")
	if key := r.Header.Get("Idempotency-Key"); key != "" {
		w.Header().Set("Idempotency-Key", key)
	}
	w.WriteHeader(status)
	if payload != nil {
		return json.NewEncoder(w).Encode(payload)
	}

	return nil
}

// DecodeJSON decodes exactly one JSON value from body into v and closes body.
//
// It returns an error if the body is empty, contains invalid JSON, or contains
// extra trailing JSON data after the first value.
func DecodeJSON(body io.ReadCloser, v any) error {
	defer func() { _ = body.Close() }()

	dec := json.NewDecoder(body)
	if err := dec.Decode(v); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("unexpected data after JSON body")
	}

	return nil
}

// AuthorizationMiddleware validates the Authorization header and authorizes its
// bearer token before calling next.
func AuthorizationMiddleware(authorizer acpauth.Authorizer) Middleware {
	if authorizer == nil {
		panic("srv: authorizer is required")
	}

	return func(w http.ResponseWriter, r *http.Request, next HandlerFunc) error {
		requestCtx := acp.RequestContextFromContext(r.Context())
		if requestCtx == nil {
			return acp.NewProcessingError("request context is missing")
		}

		bearer, err := bearerTokenFromAuthorization(requestCtx.Authorization)
		if err != nil {
			return err
		}
		if err := authorizer.Authorize(r.Context(), bearer); err != nil {
			return err
		}
		return next(w, r)
	}
}

// RequestContextMiddleware validates ACP request headers and stores the
// resulting [acp.RequestContext] on the request context passed to next.
func RequestContextMiddleware() Middleware {
	return func(w http.ResponseWriter, r *http.Request, next HandlerFunc) error {
		requestCtx, err := acp.RequestContextFromRequest(r)
		if err != nil {
			return err
		}

		r = r.WithContext(acp.ContextWithRequestContext(r.Context(), requestCtx))
		return next(w, r)
	}
}

func bearerTokenFromAuthorization(authorization string) (string, error) {
	if authorization == "" {
		return "", acp.NewHTTPError(http.StatusBadRequest, acp.InvalidRequest, acp.MissingAuthorization, "Authorization header is required")
	}

	scheme, bearer, ok := strings.Cut(authorization, " ")
	bearer = strings.TrimSpace(bearer)
	if !ok || !strings.EqualFold(scheme, "Bearer") || bearer == "" || strings.Contains(bearer, " ") {
		return "", acp.NewHTTPError(http.StatusBadRequest, acp.InvalidRequest, acp.InvalidAuthorization, "Authorization header must use Bearer authentication")
	}

	return bearer, nil
}
