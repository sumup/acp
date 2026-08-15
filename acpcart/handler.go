package acpcart

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

//go:embed spec/openapi.acpcart.yaml
var openAPISpec []byte

var requestValidator = openapi.MustNewRequestValidator(openAPISpec)

// Provider owns the seller-side cart lifecycle.
type Provider interface {
	// CreateCart creates a cart and returns its authoritative estimated state.
	CreateCart(context.Context, CartCreateRequest) (*Cart, error)
	// GetCart returns the current cart state.
	GetCart(context.Context, string) (*Cart, error)
	// UpdateCart replaces the mutable cart state and returns the updated cart.
	UpdateCart(context.Context, string, CartUpdateRequest) (*Cart, error)
	// CancelCart removes a cart and returns its final state.
	CancelCart(context.Context, string) (*Cart, error)
}

// Option configures a [Handler].
type Option func(*config)

type config struct {
	mux *http.ServeMux
}

// WithServeMux registers cart routes on mux instead of creating a new [http.ServeMux].
func WithServeMux(mux *http.ServeMux) Option {
	return func(c *config) {
		c.mux = mux
	}
}

// Handler exposes the ACP cart API over net/http.
type Handler struct {
	provider Provider
	auth     acpauth.Authorizer
	mux      *srv.Mux
}

// NewHandler returns a [Handler] that serves the ACP cart API.
func NewHandler(provider Provider, authorizer acpauth.Authorizer, opts ...Option) *Handler {
	cfg := new(config)
	for _, opt := range opts {
		opt(cfg)
	}

	h := &Handler{
		provider: provider,
		auth:     authorizer,
		mux:      srv.NewMux(cfg.mux),
	}
	h.registerRoutes()
	return h
}

// ServeHTTP dispatches cart requests to the configured ACP routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	authenticated := []srv.Middleware{srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth)}
	mutating := append(append([]srv.Middleware(nil), authenticated...), srv.IdempotencyMiddleware())

	h.mux.HandleFunc("POST /carts", h.handleCreate, mutating...)
	h.mux.HandleFunc("GET /carts/{id}", h.handleGet, authenticated...)
	h.mux.HandleFunc("PUT /carts/{id}", h.handleUpdate, mutating...)
	h.mux.HandleFunc("POST /carts/{id}/cancel", h.handleCancel, mutating...)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	var req CartCreateRequest
	if err := decodeRequiredBody(r, &req); err != nil {
		return err
	}

	cart, err := h.provider.CreateCart(r.Context(), req)
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusCreated, cart)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	cart, err := h.provider.GetCart(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusOK, cart)
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	var req CartUpdateRequest
	if err := decodeRequiredBody(r, &req); err != nil {
		return err
	}

	cart, err := h.provider.UpdateCart(r.Context(), r.PathValue("id"), req)
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusOK, cart)
}

func (h *Handler) handleCancel(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	cart, err := h.provider.CancelCart(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusOK, cart)
}

func decodeRequiredBody(r *http.Request, dst any) error {
	if err := srv.DecodeJSON(r.Body, dst); errors.Is(err, io.EOF) {
		return acp.NewInvalidRequestError("request body required")
	} else if err != nil {
		return acp.NewInvalidRequestError(err.Error())
	}
	return nil
}
