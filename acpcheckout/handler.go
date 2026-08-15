package acpcheckout

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

//go:embed spec/openapi.acpcheckout.yaml
var openAPISpec []byte

var requestValidator = openapi.MustNewRequestValidator(openAPISpec)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=handler_mock_test.go -package=acpcheckout_test

// Provider is implemented by business logic that owns checkout sessions.
type Provider interface {
	// CreateSession initializes a new checkout session from items and (optionally) buyer and fulfillment info.
	// MUST return [CheckoutSessionBase] with a rich, authoritative cart state.
	CreateSession(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSessionBase, error)

	// UpdateSession applies changes (items, fulfillment address, fulfillment option) and returns an updated authoritative cart state.
	UpdateSession(ctx context.Context, id string, req CheckoutSessionUpdateRequest) (*CheckoutSessionBase, error)

	// GetSession returns the latest authoritative state for the checkout session.
	GetSession(ctx context.Context, id string) (*CheckoutSessionBase, error)

	// CompleteSession finalizes the checkout by applying a payment method. MUST create an order and return [CheckoutSessionWithOrder] on success.
	//
	// Agents MAY include [CheckoutSessionCompleteRequest.AffiliateAttribution] to ensure the merchant receives final attribution context.
	// Attribution is stored alongside the resulting order but not returned in the response.
	CompleteSession(ctx context.Context, id string, req CheckoutSessionCompleteRequest) (CheckoutSessionWithOrder, error)

	// CancelSession cancels a session if not already completed or canceled.
	// Agents MAY include an [CancelSessionRequest.IntentTrace] to communicate why the session is being abandoned.
	CancelSession(ctx context.Context, id string, req *CancelSessionRequest) (*CheckoutSessionBase, error)
}

// Option configures a [Handler] during construction.
type Option func(*config)

type config struct {
	mux *http.ServeMux
}

// WithServeMux registers ACP checkout routes on mux instead of creating a new [http.ServeMux].
func WithServeMux(mux *http.ServeMux) Option {
	return func(c *config) {
		c.mux = mux
	}
}

// Handler wires ACP checkout routes to a [Provider].
type Handler struct {
	service Provider
	mux     *srv.Mux
	auth    acpauth.Authorizer
}

// NewHandler returns a [Handler] that serves the ACP checkout API.
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

// ServeHTTP dispatches checkout requests to the configured ACP routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	mutating := []srv.Middleware{srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth), srv.IdempotencyMiddleware()}
	h.mux.HandleFunc("POST /checkout_sessions", h.handleCreate, mutating...)
	h.mux.HandleFunc("GET /checkout_sessions/{id}", h.handleGet, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
	h.mux.HandleFunc("POST /checkout_sessions/{id}", h.handleUpdate, mutating...)
	h.mux.HandleFunc("POST /checkout_sessions/{id}/complete", h.handleComplete, mutating...)
	h.mux.HandleFunc("POST /checkout_sessions/{id}/cancel", h.handleCancel, mutating...)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	var req CheckoutSessionCreateRequest
	if err := srv.DecodeJSON(r.Body, &req); errors.Is(err, io.EOF) {
		return acp.NewInvalidRequestError("request body required")
	} else if err != nil {
		return acp.NewInvalidRequestError(err.Error())
	}
	session, err := h.service.CreateSession(r.Context(), req)
	if acpErr, ok := errors.AsType[*acp.Error](err); ok {
		return acpErr
	} else if err != nil {
		return err
	}

	return srv.WriteJSON(w, r, http.StatusCreated, session)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	if id == "" {
		return acp.NewInvalidRequestError("checkout_session_id is required")
	}

	session, err := h.service.GetSession(r.Context(), id)
	if acpErr, ok := errors.AsType[*acp.Error](err); ok {
		return acpErr
	} else if err != nil {
		return err
	}

	return srv.WriteJSON(w, r, http.StatusOK, session)
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	if id == "" {
		return acp.NewInvalidRequestError("checkout_session_id is required")
	}

	var req CheckoutSessionUpdateRequest
	if err := srv.DecodeJSON(r.Body, &req); errors.Is(err, io.EOF) {
		return acp.NewInvalidRequestError("request body required")
	} else if err != nil {
		return acp.NewInvalidRequestError(err.Error())
	}

	session, err := h.service.UpdateSession(r.Context(), id, req)
	if acpErr, ok := errors.AsType[*acp.Error](err); ok {
		return acpErr
	} else if err != nil {
		return err
	}

	return srv.WriteJSON(w, r, http.StatusOK, session)
}

func (h *Handler) handleComplete(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	if id == "" {
		return acp.NewInvalidRequestError("checkout_session_id is required")
	}

	var req CheckoutSessionCompleteRequest
	if err := srv.DecodeJSON(r.Body, &req); errors.Is(err, io.EOF) {
		return acp.NewInvalidRequestError("request body required")
	} else if err != nil {
		return acp.NewInvalidRequestError(err.Error())
	}

	session, err := h.service.CompleteSession(r.Context(), id, req)
	if acpErr, ok := errors.AsType[*acp.Error](err); ok {
		return acpErr
	} else if err != nil {
		return err
	}

	return srv.WriteJSON(w, r, http.StatusOK, session)
}

func (h *Handler) handleCancel(w http.ResponseWriter, r *http.Request) error {
	if err := requestValidator.Validate(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	if id == "" {
		return acp.NewInvalidRequestError("checkout_session_id is required")
	}

	var req CancelSessionRequest
	if r.Body != nil {
		if err := srv.DecodeJSON(r.Body, &req); err != nil && !errors.Is(err, io.EOF) {
			return acp.NewInvalidRequestError(err.Error())
		}
	}

	session, err := h.service.CancelSession(r.Context(), id, &req)
	if acpErr, ok := errors.AsType[*acp.Error](err); ok {
		return acpErr
	} else if err != nil {
		return err
	}

	return srv.WriteJSON(w, r, http.StatusOK, session)
}
