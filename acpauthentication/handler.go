package acpauthentication

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

//go:embed spec/openapi.acpauthentication.yaml
var openAPISpec []byte

var requestValidator = openapi.MustNewRequestValidator(openAPISpec)

// Provider owns delegated consumer-authentication sessions.
type Provider interface {
	// CreateAuthenticationSession initializes a 3DS authentication session.
	CreateAuthenticationSession(context.Context, DelegateAuthenticationCreateRequest) (*DelegateAuthenticationSession, error)
	// AuthenticateSession submits browser fingerprint results and starts authentication.
	AuthenticateSession(context.Context, string, DelegateAuthenticationAuthenticateRequest) (*DelegateAuthenticationSession, error)
	// GetAuthenticationSession returns the current session and its final result, when available.
	GetAuthenticationSession(context.Context, string) (*DelegateAuthenticationSessionWithResult, error)
}

// Option configures a [Handler].
type Option func(*config)

type config struct {
	mux *http.ServeMux
}

// WithServeMux registers delegated-authentication routes on mux instead of creating a new [http.ServeMux].
func WithServeMux(mux *http.ServeMux) Option {
	return func(c *config) {
		c.mux = mux
	}
}

// Handler exposes the ACP delegated authentication API over net/http.
type Handler struct {
	provider Provider
	auth     acpauth.Authorizer
	mux      *srv.Mux
}

// NewHandler returns a [Handler] that serves the ACP delegated authentication API.
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

// ServeHTTP dispatches delegated-authentication requests to the configured ACP routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	middleware := []srv.Middleware{srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth)}
	h.mux.HandleFunc("POST /delegate_authentication", h.handleCreate, middleware...)
	h.mux.HandleFunc("POST /delegate_authentication/{authentication_session_id}/authenticate", h.handleAuthenticate, middleware...)
	h.mux.HandleFunc("GET /delegate_authentication/{authentication_session_id}", h.handleGet, middleware...)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	var req DelegateAuthenticationCreateRequest
	if err := decodeRequiredBody(r, &req); err != nil {
		return err
	}

	session, err := h.provider.CreateAuthenticationSession(r.Context(), req)
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusCreated, session)
}

func (h *Handler) handleAuthenticate(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	var req DelegateAuthenticationAuthenticateRequest
	if err := decodeRequiredBody(r, &req); err != nil {
		return err
	}

	session, err := h.provider.AuthenticateSession(r.Context(), r.PathValue("authentication_session_id"), req)
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusOK, session)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	session, err := h.provider.GetAuthenticationSession(r.Context(), r.PathValue("authentication_session_id"))
	if err != nil {
		return err
	}
	return srv.WriteJSON(w, r, http.StatusOK, session)
}

func decodeRequiredBody(r *http.Request, dst any) error {
	if err := srv.DecodeJSON(r.Body, dst); errors.Is(err, io.EOF) {
		return acp.NewInvalidRequestError("request body required")
	} else if err != nil {
		return acp.NewInvalidRequestError(err.Error())
	}
	return nil
}
