package agentic_checkout

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const apiVersionHeaderValue = "2026-01-30"

// CheckoutProvider is implemented by business logic that owns checkout sessions.
type CheckoutProvider interface {
	CreateSession(ctx context.Context, req CheckoutSessionCreateRequest) (*CheckoutSessionBase, error)
	UpdateSession(ctx context.Context, id string, req CheckoutSessionUpdateRequest) (*CheckoutSessionBase, error)
	GetSession(ctx context.Context, id string) (*CheckoutSessionBase, error)
	CompleteSession(ctx context.Context, id string, req CheckoutSessionCompleteRequest) (CheckoutSessionWithOrder, error)
	CancelSession(ctx context.Context, id string, req *CancelSessionRequest) (*CheckoutSessionBase, error)
}

// CheckoutHandler wires ACP checkout routes to a CheckoutProvider.
type CheckoutHandler struct {
	service CheckoutProvider
	mux     *http.ServeMux
}

func NewCheckoutHandler(service CheckoutProvider) *CheckoutHandler {
	if service == nil {
		panic("checkout: service is required")
	}
	h := &CheckoutHandler{service: service, mux: http.NewServeMux()}
	h.registerRoutes()
	return h
}

func (h *CheckoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *CheckoutHandler) registerRoutes() {
	h.mux.HandleFunc("POST /checkout_sessions", h.handleCreate)
	h.mux.HandleFunc("GET /checkout_sessions/{id}", h.handleGet)
	h.mux.HandleFunc("POST /checkout_sessions/{id}", h.handleUpdate)
	h.mux.HandleFunc("POST /checkout_sessions/{id}/complete", h.handleComplete)
	h.mux.HandleFunc("POST /checkout_sessions/{id}/cancel", h.handleCancel)
}

func (h *CheckoutHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if err := requireIdempotencyKey(r); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "idempotency_key_required", err.Error())
		return
	}

	var req CheckoutSessionCreateRequest
	if err := decodeJSON(r.Body, &req); errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", "request body required")
		return
	} else if err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
		return
	}

	session, err := h.service.CreateSession(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, r, http.StatusCreated, session)
}

func (h *CheckoutHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", "checkout_session_id is required")
		return
	}
	session, err := h.service.GetSession(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, r, http.StatusOK, session)
}

func (h *CheckoutHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if err := requireIdempotencyKey(r); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "idempotency_key_required", err.Error())
		return
	}
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", "checkout_session_id is required")
		return
	}
	var req CheckoutSessionUpdateRequest
	if err := decodeJSON(r.Body, &req); errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", "request body required")
		return
	} else if err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
		return
	}
	session, err := h.service.UpdateSession(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, r, http.StatusOK, session)
}

func (h *CheckoutHandler) handleComplete(w http.ResponseWriter, r *http.Request) {
	if err := requireIdempotencyKey(r); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "idempotency_key_required", err.Error())
		return
	}
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", "checkout_session_id is required")
		return
	}
	var req CheckoutSessionCompleteRequest
	if err := decodeJSON(r.Body, &req); errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", "request body required")
		return
	} else if err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
		return
	}
	session, err := h.service.CompleteSession(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, r, http.StatusOK, session)
}

func (h *CheckoutHandler) handleCancel(w http.ResponseWriter, r *http.Request) {
	if err := requireIdempotencyKey(r); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, "idempotency_key_required", err.Error())
		return
	}
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", "checkout_session_id is required")
		return
	}
	var req CancelSessionRequest
	if r.Body != nil {
		if err := decodeJSON(r.Body, &req); err != nil && !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
			return
		}
		if err := req.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, InvalidRequest, "invalid_request", err.Error())
			return
		}
	}
	session, err := h.service.CancelSession(r.Context(), id, &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, r, http.StatusOK, session)
}

func decodeJSON(body io.ReadCloser, v any) error {
	defer func() { _ = body.Close() }()
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("unexpected data after JSON body")
	}
	return nil
}

func requireIdempotencyKey(r *http.Request) error {
	if r.Header.Get("Idempotency-Key") == "" {
		return errors.New("Idempotency-Key header is required")
	}
	return nil
}

type HandlerError struct {
	Status  int
	Type    ErrorType
	Code    string
	Message string
}

func (e *HandlerError) Error() string { return e.Message }

func writeServiceError(w http.ResponseWriter, err error) {
	var handlerErr *HandlerError
	if errors.As(err, &handlerErr) {
		writeError(w, handlerErr.Status, handlerErr.Type, handlerErr.Code, handlerErr.Message)
		return
	}
	writeError(w, http.StatusInternalServerError, ProcessingError, "processing_error", "internal server error")
}

func writeError(w http.ResponseWriter, status int, typ ErrorType, code, message string) {
	payload := Error{Type: typ, Code: code, Message: message}
	w.Header().Set("API-Version", apiVersionHeaderValue)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, payload any) {
	w.Header().Set("API-Version", apiVersionHeaderValue)
	w.Header().Set("Content-Type", "application/json")
	if key := r.Header.Get("Idempotency-Key"); key != "" {
		w.Header().Set("Idempotency-Key", key)
	}
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}
