package acp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

var errMissingIdempotencyKey = errors.New("Idempotency-Key header is required")

func decodeJSON(body io.ReadCloser, v any) error {
	defer func() { _ = body.Close() }()
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("unexpected data after JSON body")
	}
	return nil
}

func writeServiceError(ctx context.Context, w http.ResponseWriter, err error) {
	var httpErr *Error
	if errors.As(err, &httpErr) {
		writeJSONError(ctx, w, httpErr)
		return
	}
	writeJSONError(ctx, w, NewProcessingError("internal server error"))
}

func writeJSONError(ctx context.Context, w http.ResponseWriter, payload *Error) {
	if payload == nil {
		payload = NewProcessingError("internal server error")
	}
	setCommonResponseHeaders(w, ctx)
	w.Header().Set("Content-Type", "application/json")
	if seconds := retryAfterSeconds(payload.RetryAfter()); seconds > 0 {
		w.Header().Set("Retry-After", strconv.FormatInt(seconds, 10))
	}
	w.WriteHeader(payload.status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, payload any) {
	setCommonResponseHeaders(w, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func retryAfterSeconds(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	seconds := d / time.Second
	if d%time.Second != 0 {
		seconds++
	}
	if seconds <= 0 {
		return 1
	}
	return int64(seconds)
}

func setCommonResponseHeaders(w http.ResponseWriter, ctx context.Context) {
	w.Header().Set("API-Version", APIVersion)
	if requestCtx := RequestContextFromContext(ctx); requestCtx != nil && requestCtx.IdempotencyKey != "" {
		w.Header().Set("Idempotency-Key", requestCtx.IdempotencyKey)
	}
}

func requireIdempotencyKey(r *http.Request) error {
	requestCtx := RequestContextFromContext(r.Context())
	if requestCtx != nil && requestCtx.IdempotencyKey != "" {
		return nil
	}
	return errMissingIdempotencyKey
}
