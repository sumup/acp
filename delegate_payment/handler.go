package delegate_payment

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const apiVersionHeaderValue = "2026-01-30"

// DelegatedPaymentProvider owns delegated payment tokenization.
type DelegatedPaymentProvider interface {
	DelegatePayment(ctx context.Context, req DelegatePaymentRequest) (*DelegatePaymentResponse, error)
}

// DelegatedPaymentHandler exposes the delegate payment API over net/http.
type DelegatedPaymentHandler struct {
	service DelegatedPaymentProvider
	mux     *http.ServeMux
}

func NewDelegatedPaymentHandler(service DelegatedPaymentProvider) *DelegatedPaymentHandler {
	if service == nil {
		panic("delegate_payment: service is required")
	}
	h := &DelegatedPaymentHandler{service: service, mux: http.NewServeMux()}
	h.mux.HandleFunc("POST /agentic_commerce/delegate_payment", h.handleDelegatePayment)
	return h
}

func (h *DelegatedPaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *DelegatedPaymentHandler) handleDelegatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Idempotency-Key") == "" {
		writeError(w, http.StatusBadRequest, InvalidRequest, ErrorCode("idempotency_key_required"), "Idempotency-Key header is required")
		return
	}
	var req DelegatePaymentRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, ErrorCode("invalid_request"), err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, InvalidRequest, ErrorCode("invalid_request"), err.Error())
		return
	}
	resp, err := h.service.DelegatePayment(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, r, http.StatusCreated, resp)
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

type HandlerError struct {
	Status  int
	Type    ErrorType
	Code    ErrorCode
	Message string
}

func (e *HandlerError) Error() string { return e.Message }

func writeServiceError(w http.ResponseWriter, err error) {
	var handlerErr *HandlerError
	if errors.As(err, &handlerErr) {
		writeError(w, handlerErr.Status, handlerErr.Type, handlerErr.Code, handlerErr.Message)
		return
	}
	writeError(w, http.StatusInternalServerError, ProcessingError, ErrorCode("processing_error"), "internal server error")
}

func writeError(w http.ResponseWriter, status int, typ ErrorType, code ErrorCode, message string) {
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
