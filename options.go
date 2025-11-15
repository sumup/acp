package acp

import (
	"net/http"
	"strings"
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
	webhook               *webhookConfig
}

type webhookConfig struct {
	endpoint string
	header   string
	secret   []byte
	client   *http.Client
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

// WebhookOptions configure how the checkout handler emits webhook events to OpenAI.
type WebhookOptions struct {
	// Endpoint is the absolute URL provided by OpenAI for receiving webhook events.
	Endpoint string
	// HeaderName controls the signature header name (for example Merchant_Name-Signature).
	HeaderName string
	// SecretKey is the HMAC secret provided by OpenAI for signing webhook payloads.
	SecretKey []byte
	// Client allows overriding the HTTP client used for delivering webhook events.
	Client *http.Client
}

// WithWebhookOptions configures webhook delivery for [CheckoutHandler.SendWebhook].
func WithWebhookOptions(opts WebhookOptions) Option {
	endpoint := strings.TrimSpace(opts.Endpoint)
	if endpoint == "" {
		panic("checkout: webhook endpoint is required")
	}
	header := strings.TrimSpace(opts.HeaderName)
	if header == "" {
		panic("checkout: webhook header name is required")
	}
	if len(opts.SecretKey) == 0 {
		panic("checkout: webhook secret key is required")
	}
	secret := append([]byte(nil), opts.SecretKey...)
	client := opts.Client
	if client == nil {
		client = http.DefaultClient
	}
	return func(cfg *config) {
		cfg.webhook = &webhookConfig{
			endpoint: endpoint,
			header:   header,
			secret:   secret,
			client:   client,
		}
	}
}
