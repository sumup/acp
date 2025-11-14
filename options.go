package acp

import (
	"net/http"
	"time"

	"github.com/sumup/acp/signature"
)

type config struct {
	signatureVerifier     signature.Verifier
	maxClockSkew          time.Duration
	requireSignedRequests bool
	middleware            []Middleware
	authenticator         Authenticator
	clock                 func() time.Time
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func applyMiddleware(h http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// Option customizes the handler behavior.
type Option func(*config)

// WithSignatureVerifier enables canonical JSON signature enforcement.
func WithSignatureVerifier(verifier signature.Verifier) Option {
	return func(cfg *config) {
		cfg.signatureVerifier = verifier
	}
}

// WithMaxClockSkew sets the tolerated absolute difference between the
// Timestamp header and the server clock when verifying signed requests.
func WithMaxClockSkew(skew time.Duration) Option {
	if skew <= 0 {
		panic("checkout: max clock skew must be positive")
	}
	return func(cfg *config) {
		cfg.maxClockSkew = skew
	}
}

// WithRequireSignedRequests enforces that every request carries Signature and
// Timestamp headers when a verifier is configured.
func WithRequireSignedRequests() Option {
	return func(cfg *config) {
		cfg.requireSignedRequests = true
	}
}

// WithMiddleware appends custom middleware in the order provided.
func WithMiddleware(mw ...Middleware) Option {
	return func(cfg *config) {
		for _, m := range mw {
			if m == nil {
				continue
			}
			cfg.middleware = append(cfg.middleware, m)
		}
	}
}

// WithAuthenticator enables Authorization header API key validation.
func WithAuthenticator(auth Authenticator) Option {
	return func(cfg *config) {
		cfg.authenticator = auth
	}
}

// withClock provides deterministic time in tests.
func checkoutWithClock(fn func() time.Time) Option {
	return func(cfg *config) {
		cfg.clock = fn
	}
}
