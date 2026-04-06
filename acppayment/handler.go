package acppayment

import (
	"context"
	_ "embed"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/internal/openapi"
	"github.com/sumup/acp/internal/srv"
)

//go:embed spec/openapi.acppayment.yaml
var openAPISpec []byte

var requestValidator = openapi.MustNewRequestValidator(openAPISpec)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=handler_mock_test.go -package=acppayment_test

// DelegatedPaymentProvider owns delegated payment tokenization.
type DelegatedPaymentProvider interface {
	DelegatePayment(ctx context.Context, req DelegatePaymentRequest) (*DelegatePaymentResponse, error)
}

// DelegatedPaymentHandler exposes the delegate payment API over net/http.
type DelegatedPaymentHandler struct {
	service DelegatedPaymentProvider
	mux     *srv.Mux
	auth    acpauth.Authorizer
}

func NewDelegatedPaymentHandler(service DelegatedPaymentProvider, authorizer acpauth.Authorizer) *DelegatedPaymentHandler {
	h := &DelegatedPaymentHandler{
		service: service,
		auth:    authorizer,
	}
	h.mux = srv.NewMux(h.handleError)
	h.registerRoutes()
	return h
}

func (h *DelegatedPaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *DelegatedPaymentHandler) registerRoutes() {
	h.mux.HandleFunc("POST /agentic_commerce/delegate_payment", h.handleDelegatePayment, srv.RequestContextMiddleware(), srv.AuthorizationMiddleware(h.auth))
}

func (h *DelegatedPaymentHandler) handleError(w http.ResponseWriter, _ *http.Request, err error) {
	if acpErr := new(acp.Error); errors.As(err, &acpErr) {
		_ = srv.WriteACPError(w, acpErr, func(errorType, code, message string, param *string) Error {
			mappedCode := ErrorCode(code)
			if mappedCode == ErrorCode(acp.InvalidRequest) && param != nil && strings.HasPrefix(*param, "$.payment_method") {
				mappedCode = InvalidCard
			}

			return Error{
				Type:    ErrorType(errorType),
				Code:    mappedCode,
				Message: message,
				Param:   param,
			}
		})
		return
	}

	_ = srv.WriteError(w, http.StatusInternalServerError, Error{
		Type:    ProcessingError,
		Code:    ErrorCode(acp.ProcessingError),
		Message: "internal server error",
	})
}

func (h *DelegatedPaymentHandler) handleDelegatePayment(w http.ResponseWriter, r *http.Request) error {
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
