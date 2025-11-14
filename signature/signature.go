package signature

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	canonicaljson "github.com/gibson042/canonicaljson-go"
)

// Material captures the inputs needed to validate a signed request.
type Material struct {
	Signature     string
	Timestamp     time.Time
	CanonicalBody []byte
	Method        string
	Path          string
	RawQuery      string
	Headers       http.Header
}

// Verifier validates the authenticity of incoming requests.
type Verifier interface {
	Verify(ctx context.Context, material Material) error
}

// VerifierFunc lifts bare functions into [Verifier].
type VerifierFunc func(ctx context.Context, material Material) error

// Verify delegates to the wrapped function.
func (f VerifierFunc) Verify(ctx context.Context, material Material) error {
	return f(ctx, material)
}

// HMACVerifier validates signatures that were produced by taking the
// base64url-encoded HMAC-SHA256 of `RFC3339(timestamp) + "." + canonicalJSON`.
type HMACVerifier struct {
	Key []byte
}

// Verify implements [Verifier] by recomputing the expected HMAC signature.
func (v HMACVerifier) Verify(_ context.Context, material Material) error {
	if len(v.Key) == 0 {
		return errors.New("signature: HMACSignatureVerifier requires a non-empty key")
	}
	signingInput := BuildSigningPayload(material.Timestamp, material.CanonicalBody)
	mac := hmac.New(sha256.New, v.Key)
	if _, err := mac.Write(signingInput); err != nil {
		return fmt.Errorf("signature: compute signature: %w", err)
	}
	expected := mac.Sum(nil)
	decoded, err := base64.RawURLEncoding.DecodeString(material.Signature)
	if err != nil {
		return fmt.Errorf("signature: decode signature: %w", err)
	}
	if !hmac.Equal(decoded, expected) {
		return errors.New("signature: invalid signature")
	}
	return nil
}

// ReadAndBufferBody reads the request body while keeping it accessible for later handlers.
func ReadAndBufferBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		r.Body = io.NopCloser(bytes.NewReader(nil))
		return nil, nil
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(raw))
	return raw, nil
}

// CanonicalizeJSONBody normalizes arbitrary JSON into canonical form for signing.
func CanonicalizeJSONBody(raw []byte) ([]byte, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return []byte("null"), nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var payload any
	if err := dec.Decode(&payload); err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, errors.New("signature: multiple JSON documents in body")
	}
	return canonicaljson.Marshal(payload)
}

// ParseTimestamp accepts Timestamp header values in RFC3339 or RFC3339Nano format.
func ParseTimestamp(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("signature: empty timestamp")
	}
	if ts, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return ts, nil
	}
	return time.Parse(time.RFC3339, value)
}

// AbsDuration returns the absolute value of the supplied duration.
func AbsDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// BuildSigningPayload constructs the canonical string that is HMAC-signed.
func BuildSigningPayload(ts time.Time, canonicalBody []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString(ts.UTC().Format(time.RFC3339Nano))
	buf.WriteByte('.')
	buf.Write(canonicalBody)
	return buf.Bytes()
}
