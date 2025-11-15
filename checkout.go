package acp

import (
	"context"
	"net/http"
	"time"
)

// CheckoutProvider is implemented by business logic that owns checkout sessions.
type CheckoutProvider interface {
	CreateSession(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSession, error)
	UpdateSession(ctx context.Context, id string, req CheckoutSessionUpdateRequest) (*CheckoutSession, error)
	GetSession(ctx context.Context, id string) (*CheckoutSession, error)
	CompleteSession(ctx context.Context, id string, req CheckoutSessionCompleteRequest) (*SessionWithOrder, error)
	CancelSession(ctx context.Context, id string) (*CheckoutSession, error)
}

// CheckoutHandler wires ACP checkout routes to a [CheckoutProvider].
type CheckoutHandler struct {
	service CheckoutProvider
	mux     *http.ServeMux
	cfg     config
}

// NewCheckoutHandler builds a [CheckoutHandler] backed by net/http's ServeMux.
func NewCheckoutHandler(service CheckoutProvider, opts ...Option) *CheckoutHandler {
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
		panic("checkout: signature verifier required when signed requests are enforced")
	}
	h := &CheckoutHandler{
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
	h.registerRoutes(middleware...)
	return h
}

// ServeHTTP satisfies http.Handler.
func (h *CheckoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestCtx := requestContextFromRequest(r)
	ctx := contextWithRequestContext(r.Context(), requestCtx)
	h.mux.ServeHTTP(w, r.WithContext(ctx))
}

func (h *CheckoutHandler) registerRoutes(middleware ...Middleware) {
	h.mux.HandleFunc("POST /checkout_sessions", applyMiddleware(h.handleCreate, middleware...))
	h.mux.HandleFunc("GET /checkout_sessions/{id}", applyMiddleware(h.handleGet, middleware...))
	h.mux.HandleFunc("POST /checkout_sessions/{id}", applyMiddleware(h.handleUpdate, middleware...))
	h.mux.HandleFunc("POST /checkout_sessions/{id}/complete", applyMiddleware(h.handleComplete, middleware...))
	h.mux.HandleFunc("POST /checkout_sessions/{id}/cancel", applyMiddleware(h.handleCancel, middleware...))
}

func (h *CheckoutHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CheckoutSessionCreateRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	session, err := h.service.CreateSession(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

func (h *CheckoutHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, NewInvalidRequestError("checkout_session_id is required"))
		return
	}
	session, err := h.service.GetSession(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *CheckoutHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, NewInvalidRequestError("checkout_session_id is required"))
		return
	}
	var req CheckoutSessionUpdateRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	session, err := h.service.UpdateSession(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *CheckoutHandler) handleComplete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, NewInvalidRequestError("checkout_session_id is required"))
		return
	}
	var req CheckoutSessionCompleteRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		writeJSONError(w, NewInvalidRequestError(err.Error()))
		return
	}
	session, err := h.service.CompleteSession(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *CheckoutHandler) handleCancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, NewInvalidRequestError("checkout_session_id is required"))
		return
	}
	session, err := h.service.CancelSession(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}
