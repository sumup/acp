package acp

import (
	"net/http"
	"time"
)

// ErrorType mirrors the ACP error.type field.
type ErrorType string

const (
	InvalidRequest ErrorType = "invalid_request" // Missing or malformed field.
	// Deprecated: use invalid_request + code idempotency_conflict/idempotency_key_required/idempotency_in_flight.
	RequestNotIdempotentType ErrorType = "request_not_idempotent" // Idempotency violation.
	ProcessingError          ErrorType = "processing_error"       // Downstream gateway or network failure.
	RateLimitExceeded        ErrorType = "rate_limit_exceeded"    // Too many requests.
	ServiceUnavailable       ErrorType = "service_unavailable"    // Temporary outage or maintenance.
)

// ErrorCode is a machine-readable identifier for the specific failure.
type ErrorCode string

const (
	DuplicateRequest       ErrorCode = "duplicate_request"        // Safe duplicate with the same idempotency key.
	IdempotencyConflict    ErrorCode = "idempotency_conflict"     // Same idempotency key but different parameters.
	IdempotencyKeyRequired ErrorCode = "idempotency_key_required" // Idempotency-Key header is missing.
	IdempotencyInFlight    ErrorCode = "idempotency_in_flight"    // Request with same idempotency key is still processing.
	InvalidCard            ErrorCode = "invalid_card"             // Credential failed basic validation (such as length or expiry).
	InvalidSignature       ErrorCode = "invalid_signature"        // Signature is missing or does not match the payload.
	SignatureRequired      ErrorCode = "signature_required"       // Signed requests are required but headers were missing.
	StaleTimestamp         ErrorCode = "stale_timestamp"          // Timestamp skew exceeded the allowed window.
	MissingAuthorization   ErrorCode = "missing_authorization"    // Authorization header missing.
	InvalidAuthorization   ErrorCode = "invalid_authorization"    // Authorization header malformed or API key invalid.
	RequestNotIdempotent   ErrorCode = "request_not_idempotent"
)

// Error represents a structured ACP error payload.
type Error struct {
	Type              ErrorType `json:"type"`
	Code              ErrorCode `json:"code"`
	Message           string    `json:"message"`
	Param             *string   `json:"param,omitempty"`
	SupportedVersions []string  `json:"supported_versions,omitempty"`

	Internal error `json:"-"`

	status     int           `json:"-"`
	retryAfter time.Duration `json:"-"`
}

// Error makes *Error satisfy the stdlib error interface.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// RetryAfter returns the duration clients should wait before retrying.
func (e *Error) RetryAfter() time.Duration {
	if e == nil {
		return 0
	}
	return e.retryAfter
}

// StatusCode returns the HTTP status associated with the error.
func (e *Error) StatusCode() int {
	if e == nil || e.status == 0 {
		return http.StatusInternalServerError
	}
	return e.status
}

type errorOption func(*Error)

// WithOffendingParam sets the JSON path for the field that triggered the error.
func WithOffendingParam(jsonPath string) errorOption {
	return func(er *Error) {
		er.Param = &jsonPath
	}
}

// WithStatusCode overrides the HTTP status code returned to the client.
func WithStatusCode(status int) errorOption {
	return func(er *Error) {
		er.status = status
	}
}

// WithRetryAfter specifies how long clients should wait before retrying.
func WithRetryAfter(d time.Duration) errorOption {
	return func(er *Error) {
		er.retryAfter = d
	}
}

// WithInternal sets an internal error that can be retrieved from [Error] for telemetry purposes.
func WithInternal(err error) errorOption {
	return func(er *Error) {
		er.Internal = err
	}
}

// WithSupportedVersions sets the supported protocol versions returned to clients.
func WithSupportedVersions(versions []string) errorOption {
	return func(er *Error) {
		if len(versions) == 0 {
			return
		}
		er.SupportedVersions = append([]string(nil), versions...)
	}
}

// NewRateLimitExceededError builds a Too Many Requests ACP error payload.
func NewRateLimitExceededError(message string, opts ...errorOption) *Error {
	return newError(RateLimitExceeded, ErrorCode(RateLimitExceeded), message, append([]errorOption{WithStatusCode(http.StatusTooManyRequests)}, opts...)...)
}

// NewServiceUnavailableError builds a Service Unavailable ACP error payload.
func NewServiceUnavailableError(message string, opts ...errorOption) *Error {
	return newError(ServiceUnavailable, ErrorCode(ServiceUnavailable), message, append([]errorOption{WithStatusCode(http.StatusServiceUnavailable)}, opts...)...)
}

// NewInvalidRequestError builds a Bad Request ACP error payload.
func NewInvalidRequestError(message string, opts ...errorOption) *Error {
	return newError(InvalidRequest, ErrorCode(InvalidRequest), message, append([]errorOption{WithStatusCode(http.StatusBadRequest)}, opts...)...)
}

// NewProcessingError builds an Internal Server Error ACP error payload.
func NewProcessingError(message string, opts ...errorOption) *Error {
	return newError(ProcessingError, ErrorCode(ProcessingError), message, append([]errorOption{WithStatusCode(http.StatusInternalServerError)}, opts...)...)
}

// NewHTTPError allows callers to control the status code explicitly.
func NewHTTPError(status int, typ ErrorType, code ErrorCode, message string, opts ...errorOption) *Error {
	return newError(typ, code, message, append(opts, WithStatusCode(status))...)
}

// newError builds a typed error payload matching the ACP schema.
func newError(typ ErrorType, code ErrorCode, message string, opts ...errorOption) *Error {
	errPayload := &Error{
		Type:    typ,
		Code:    code,
		Message: message,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(errPayload)
	}
	return errPayload
}
