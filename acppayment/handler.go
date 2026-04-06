package acppayment

import (
	"context"
	_ "embed"
	"errors"
	"io"
	"net/http"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/internal/openapi"
	"github.com/sumup/acp/internal/srv"
)

//go:embed spec/openapi.acppayment.yaml
var openAPISpec []byte

var requestValidator = openapi.MustNewRequestValidator(openAPISpec)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=handler_mock_test.go -package=acppayment_test

// Provider owns delegated payment tokenization.
type Provider interface {
	// DelegatePayment tokenizes the supplied payment method for later checkout completion.
	DelegatePayment(ctx context.Context, req DelegatePaymentRequest) (*DelegatePaymentResponse, error)
}

// Option configures a [Handler] during construction.
type Option func(*config)

type config struct {
	mux *http.ServeMux
}

// WithServeMux registers ACP payment routes on mux instead of creating a new [http.ServeMux].
func WithServeMux(mux *http.ServeMux) Option {
	return func(c *config) {
		c.mux = mux
	}
}

// Handler exposes the delegate payment API over net/http.
type Handler struct {
	service Provider
	mux     *srv.Mux
	auth    acpauth.Authorizer
}

// NewHandler returns a [Handler] that serves the ACP delegated payment API.
func NewHandler(service Provider, authorizer acpauth.Authorizer, opts ...Option) *Handler {
	cfg := new(config)

	for _, opt := range opts {
		opt(cfg)
	}

	h := &Handler{
		service: service,
		auth:    authorizer,
	}
	h.mux = srv.NewMux(cfg.mux)
	h.registerRoutes()
	return h
}

// ServeHTTP dispatches delegated payment requests to the configured ACP routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	h.mux.HandleFunc("POST /agentic_commerce/delegate_payment", h.handleDelegatePayment, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
}

func (h *Handler) handleDelegatePayment(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	var req DelegatePaymentRequest
	if err := srv.DecodeJSON(r.Body, &req); errors.Is(err, io.EOF) {
		return acp.NewInvalidRequestError("request body required")
	} else if err != nil {
		return acp.NewInvalidRequestError(err.Error())
	}

	resp, err := h.service.DelegatePayment(r.Context(), req)
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusCreated, resp)
}
