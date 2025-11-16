package acp

import (
	"fmt"
	"net/http"
	"net/url"
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
	endpoint        *url.URL
	signatureHeader string
	secret          []byte
	client          *http.Client
}

// Middleware is an HTTP middleware applied to the ACP handlers.
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
	// MerchantName controls the signature header name (the header name is Merchant_Name-Signature).
	MerchantName string
	// SecretKey is the HMAC secret provided by OpenAI for signing webhook payloads.
	SecretKey []byte
	// Client allows overriding the HTTP client used for delivering webhook events.
	Client *http.Client
}

// WithWebhookOptions configures webhook delivery for [CheckoutHandler.SendWebhook].
func WithWebhookOptions(opts WebhookOptions) Option {
	endpoint, err := url.Parse(opts.Endpoint)
	if err != nil {
		panic(fmt.Errorf("with webhook: parse endpoint url: %w", err))
	}

	merchantName := strings.TrimSpace(opts.MerchantName)
	if merchantName == "" {
		panic("with webhook: webhook header name is required")
	}
	header := fmt.Sprintf("%s-Signature", strings.ReplaceAll(merchantName, " ", "_"))

	if len(opts.SecretKey) == 0 {
		panic("with webhook: webhook secret key is required")
	}

	secret := append([]byte(nil), opts.SecretKey...)
	client := opts.Client
	if client == nil {
		client = http.DefaultClient
	}

	return func(cfg *config) {
		cfg.webhook = &webhookConfig{
			endpoint:        endpoint,
			signatureHeader: header,
			secret:          secret,
			client:          client,
		}
	}
}
