package acp

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// RequestContext carries the standard ACP headers.
type RequestContext struct {
	// API Key used to make requests.
	//
	// Example: Bearer api_key_123
	Authorization string
	// The preferred locale for content like messages and errors.
	//
	// Example: en-US
	AcceptLanguage string
	// Information about the client making this request.
	//
	// Example: ChatGPT/2.0 (Mac OS X 15.0.1; arm64; build 0)
	UserAgent string
	// Key used to ensure requests are idempotent.
	//
	// Example: idempotency_key_123
	IdempotencyKey string
	// Unique key for each request for tracing purposes.
	//
	// Example: request_id_123
	RequestID string
	// Optional detached signature used to verify the request body.
	//
	// Example: eyJtZX...
	Signature string
	// Optional request-signing timestamp formatted as an RFC 3339 string.
	//
	// Example: 2025-09-25T10:30:00Z
	Timestamp string
	// API version.
	//
	// Example: 2026-04-17
	APIVersion string
}

// RequestContextFromRequest reads the standard ACP headers from r and validates
// that API-Version identifies the version implemented by this module.
func RequestContextFromRequest(r *http.Request) (*RequestContext, error) {
	requestCtx := &RequestContext{
		Authorization:  strings.TrimSpace(r.Header.Get("Authorization")),
		AcceptLanguage: strings.TrimSpace(r.Header.Get("Accept-Language")),
		UserAgent:      strings.TrimSpace(r.Header.Get("User-Agent")),
		IdempotencyKey: strings.TrimSpace(r.Header.Get("Idempotency-Key")),
		RequestID:      strings.TrimSpace(r.Header.Get("Request-Id")),
		Signature:      strings.TrimSpace(r.Header.Get("Signature")),
		Timestamp:      strings.TrimSpace(r.Header.Get("Timestamp")),
		APIVersion:     strings.TrimSpace(r.Header.Get("API-Version")),
	}
	if err := requestCtx.validate(); err != nil {
		return nil, err
	}
	return requestCtx, nil
}

// ContextWithRequestContextFromRequest reads and validates ACP headers from r,
// then returns a child of r.Context containing that metadata.
func ContextWithRequestContextFromRequest(r *http.Request) (context.Context, error) {
	requestCtx, err := RequestContextFromRequest(r)
	if err != nil {
		return nil, err
	}
	return ContextWithRequestContext(r.Context(), requestCtx), nil
}

func (r *RequestContext) validate() *Error {
	if r == nil {
		return NewInvalidRequestError("request context is required")
	}
	if r.APIVersion == "" {
		return NewHTTPError(
			http.StatusBadRequest,
			InvalidRequest,
			MissingAPIVersion,
			"API-Version header is required",
			WithSupportedVersions([]string{APIVersion}),
		)
	}
	if parsed, err := time.Parse("2006-01-02", r.APIVersion); err != nil || parsed.Format("2006-01-02") != r.APIVersion {
		return NewInvalidRequestError("API-Version header must be in YYYY-MM-DD format")
	}
	if r.APIVersion != APIVersion {
		return NewHTTPError(
			http.StatusBadRequest,
			InvalidRequest,
			UnsupportedAPIVersion,
			"API version is not supported",
			WithSupportedVersions([]string{APIVersion}),
		)
	}
	return nil
}

type requestContextKey struct{}

// ContextWithRequestContext returns a context containing requestCtx.
//
// A nil requestCtx leaves ctx unchanged. A nil ctx is replaced with
// [context.Background].
func ContextWithRequestContext(ctx context.Context, requestCtx *RequestContext) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if requestCtx == nil {
		return ctx
	}
	return context.WithValue(ctx, requestContextKey{}, requestCtx)
}

// RequestContextFromContext extracts the HTTP request metadata previously stored in the context.
func RequestContextFromContext(ctx context.Context) *RequestContext {
	if ctx == nil {
		return nil
	}
	if requestCtx, ok := ctx.Value(requestContextKey{}).(*RequestContext); ok {
		return requestCtx
	}
	return nil
}
