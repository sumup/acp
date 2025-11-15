package acp

import (
	"context"
	"net/http"
	"strings"
)

type RequestContext struct {
	// API Key used to make requests
	//
	// Example: Bearer api_key_123
	Authorization string
	// The preferred locale for content like messages and errors
	//
	// Example: en-US
	AcceptLanguage string
	// Information about the client making this request
	//
	// Example: ChatGPT/2.0 (Mac OS X 15.0.1; arm64; build 0)
	UserAgent string
	// Key used to ensure requests are idempotent
	//
	// Example: idempotency_key_123
	IdempotencyKey string
	// Unique key for each request for tracing purposes
	//
	// Example: request_id_123
	RequestID string
	// Base64 encoded signature of the request body
	//
	// Example: eyJtZX...
	Signature string
	// Formatted as an RFC 3339 string.
	//
	// Example: 2025-09-25T10:30:00Z
	Timestamp string
	// API version
	//
	// Example: 2025-09-12
	APIVersion string
}

func requestContextFromRequest(r *http.Request) *RequestContext {
	return &RequestContext{
		Authorization:  strings.TrimSpace(r.Header.Get("Authorization")),
		AcceptLanguage: strings.TrimSpace(r.Header.Get("Accept-Language")),
		UserAgent:      strings.TrimSpace(r.Header.Get("User-Agent")),
		IdempotencyKey: strings.TrimSpace(r.Header.Get("Idempotency-Key")),
		RequestID:      strings.TrimSpace(r.Header.Get("Request-Id")),
		Signature:      strings.TrimSpace(r.Header.Get("Signature")),
		Timestamp:      strings.TrimSpace(r.Header.Get("Timestamp")),
		APIVersion:     strings.TrimSpace(r.Header.Get("API-Version")),
	}
}

type requestContextKey struct{}

func contextWithRequestContext(ctx context.Context, requestCtx *RequestContext) context.Context {
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
