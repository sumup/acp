package agentic_checkout

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

//go:embed spec/openapi.agentic_checkout.yaml
var openAPISpec []byte

var requestValidator = openapi.MustNewRequestValidator(openAPISpec)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=handler_mock_test.go -package=agentic_checkout_test

// CheckoutProvider is implemented by business logic that owns checkout sessions.
type CheckoutProvider interface {
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

// CheckoutHandler wires ACP checkout routes to a CheckoutProvider.
type CheckoutHandler struct {
	service CheckoutProvider
	mux     *srv.Mux
	auth    acpauth.Authorizer
}

func NewCheckoutHandler(service CheckoutProvider, authorizer acpauth.Authorizer) *CheckoutHandler {
	h := &CheckoutHandler{
		service: service,
		auth:    authorizer,
	}
	h.mux = srv.NewMux(h.handleError)
	h.registerRoutes()
	return h
}

func (h *CheckoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *CheckoutHandler) registerRoutes() {
	h.mux.HandleFunc("POST /checkout_sessions", h.handleCreate, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
	h.mux.HandleFunc("GET /checkout_sessions/{id}", h.handleGet, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
	h.mux.HandleFunc("POST /checkout_sessions/{id}", h.handleUpdate, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
	h.mux.HandleFunc("POST /checkout_sessions/{id}/complete", h.handleComplete, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
	h.mux.HandleFunc("POST /checkout_sessions/{id}/cancel", h.handleCancel, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
}

func (h *CheckoutHandler) handleError(w http.ResponseWriter, _ *http.Request, err error) {
	if acpErr := new(acp.Error); errors.As(err, &acpErr) {
		_ = srv.WriteACPError(w, acpErr, func(errorType, code, message string, param *string) Error {
			return Error{
				Type:    ErrorType(errorType),
				Code:    code,
				Message: message,
				Param:   param,
			}
		})
		return
	}

	_ = srv.WriteError(w, http.StatusInternalServerError, Error{
		Type:    ProcessingError,
		Code:    string(acp.ProcessingError),
		Message: "internal server error",
	})
}

func (h *CheckoutHandler) handleCreate(w http.ResponseWriter, r *http.Request) error {
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

func (h *CheckoutHandler) handleGet(w http.ResponseWriter, r *http.Request) error {
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

func (h *CheckoutHandler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
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

func (h *CheckoutHandler) handleComplete(w http.ResponseWriter, r *http.Request) error {
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

func (h *CheckoutHandler) handleCancel(w http.ResponseWriter, r *http.Request) error {
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
