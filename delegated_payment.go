package acp

import (
	"context"
	"net/http"
	"time"
)

// DelegatedPaymentProvider owns the delegated payment tokenization lifecycle.
// To integrate your Payments Service Provider with the Delegate Payment Spec
// implement this interface.
type DelegatedPaymentProvider interface {
	DelegatePayment(ctx context.Context, req PaymentRequest) (*VaultToken, error)
}

// DelegatedPaymentHandler exposes the ACP delegate payment API over net/http.
type DelegatedPaymentHandler struct {
	service DelegatedPaymentProvider
	mux     *http.ServeMux
	cfg     config
}

// NewDelegatedPaymentHandler wires the delegate payment routes to the provided [DelegatedPaymentProvider].
func NewDelegatedPaymentHandler(service DelegatedPaymentProvider, opts ...Option) *DelegatedPaymentHandler {
	if service == nil {
		panic("delegatedpayment: service is required")
	}
	cfg := config{
		maxClockSkew: 5 * time.Minute,
		clock:        time.Now,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}
	if cfg.requireSignedRequests && cfg.signatureVerifier == nil {
		panic("delegatedpayment: signature verifier required when signed requests are enforced")
	}
	h := &DelegatedPaymentHandler{
		service: service,
		mux:     http.NewServeMux(),
		cfg:     cfg,
	}
	var middleware []Middleware
	if mw := newSignatureMiddleware(signatureMiddlewareConfig{
		Verifier:      cfg.signatureVerifier,
		RequireSigned: cfg.requireSignedRequests,
		MaxClockSkew:  cfg.maxClockSkew,
		Clock:         cfg.clock,
	}); mw != nil {
		middleware = append(middleware, Middleware(mw))
	}
	if cfg.authenticator != nil {
		middleware = append(middleware, h.authenticationMiddleware)
	}
	middleware = append(middleware, cfg.middleware...)
	h.registerRoutes(middleware...)
	return h
}

// ServeHTTP satisfies http.Handler.
func (h *DelegatedPaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *DelegatedPaymentHandler) registerRoutes(middleware ...Middleware) {
	h.mux.HandleFunc("POST /agentic_commerce/delegate_payment", applyMiddleware(h.handleDelegatePayment, middleware...))
}

func (h *DelegatedPaymentHandler) handleDelegatePayment(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	resp, err := h.service.DelegatePayment(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}
