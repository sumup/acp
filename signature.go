package acp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sumup/acp/signature"
)

type signatureMiddlewareConfig struct {
	Verifier      signature.Verifier
	RequireSigned bool
	MaxClockSkew  time.Duration
	Clock         func() time.Time
}

func newSignatureMiddleware(cfg signatureMiddlewareConfig) func(http.HandlerFunc) http.HandlerFunc {
	if cfg.Verifier == nil {
		return nil
	}
	if cfg.Clock == nil {
		cfg.Clock = time.Now
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			verifier := cfg.Verifier
			if verifier == nil {
				next(w, r)
				return
			}
			sig := strings.TrimSpace(r.Header.Get("Signature"))
			timestampHeader := strings.TrimSpace(r.Header.Get("Timestamp"))
			if sig == "" && timestampHeader == "" {
				if cfg.RequireSigned {
					writeJSONError(w, NewHTTPError(http.StatusUnauthorized, InvalidRequest, SignatureRequired, "Signature and Timestamp headers are required"))
					return
				}
				next(w, r)
				return
			}
			if sig == "" || timestampHeader == "" {
				writeJSONError(w, NewHTTPError(http.StatusBadRequest, InvalidRequest, InvalidSignature, "Signature and Timestamp headers must both be provided"))
				return
			}
			ts, err := signature.ParseTimestamp(timestampHeader)
			if err != nil {
				writeJSONError(w, NewHTTPError(http.StatusBadRequest, InvalidRequest, InvalidSignature, "Timestamp must be RFC3339"))
				return
			}
			ts = ts.UTC()
			if cfg.MaxClockSkew > 0 {
				skew := signature.AbsDuration(cfg.Clock().Sub(ts))
				if skew > cfg.MaxClockSkew {
					writeJSONError(w, NewHTTPError(http.StatusUnauthorized, InvalidRequest, StaleTimestamp, fmt.Sprintf("timestamp skew exceeds %s", cfg.MaxClockSkew)))
					return
				}
			}
			raw, err := signature.ReadAndBufferBody(r)
			if err != nil {
				writeJSONError(w, NewInvalidRequestError("unable to read request body"))
				return
			}
			canonicalBody, err := signature.CanonicalizeJSONBody(raw)
			if err != nil {
				writeJSONError(w, NewInvalidRequestError("request body must be valid JSON"))
				return
			}
			material := signature.Material{
				Signature:     sig,
				Timestamp:     ts,
				CanonicalBody: canonicalBody,
				Method:        r.Method,
				Path:          r.URL.Path,
				RawQuery:      r.URL.RawQuery,
				Headers:       r.Header.Clone(),
			}
			if err := verifier.Verify(r.Context(), material); err != nil {
				writeJSONError(w, NewHTTPError(http.StatusUnauthorized, InvalidRequest, InvalidSignature, "signature verification failed"))
				return
			}
			next(w, r)
		}
	}
}
