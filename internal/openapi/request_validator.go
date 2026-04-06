package openapi

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
	validatorconfig "github.com/pb33f/libopenapi-validator/config"
	validatorerrors "github.com/pb33f/libopenapi-validator/errors"
	validatorhelpers "github.com/pb33f/libopenapi-validator/helpers"
	"github.com/pb33f/libopenapi/datamodel"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/sumup/acp"
)

// RequestValidator validates incoming HTTP requests against a bundled OpenAPI document.
type RequestValidator struct {
	document  *v3high.Document
	validator validator.Validator
	mu        sync.Mutex
}

// NewRequestValidator builds a request validator from an OpenAPI v3 spec.
func NewRequestValidator(spec []byte) (*RequestValidator, error) {
	docCfg := datamodel.NewDocumentConfiguration()
	docCfg.SkipCircularReferenceCheck = true

	document, err := libopenapi.NewDocumentWithConfiguration(spec, docCfg)
	if err != nil {
		return nil, fmt.Errorf("parse OpenAPI document: %w", err)
	}

	v3Model, err := document.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("build OpenAPI model: %w", err)
	}

	requestValidator := validator.NewValidatorFromV3Model(
		&v3Model.Model,
		validatorconfig.WithoutSecurityValidation(),
		validatorconfig.WithContentAssertions(),
		validatorconfig.WithFormatAssertions(),
	)
	requestValidator.SetDocument(document)

	return &RequestValidator{
		document:  &v3Model.Model,
		validator: requestValidator,
	}, nil
}

// MustNewRequestValidator builds a validator and panics if the bundled spec is invalid.
func MustNewRequestValidator(spec []byte) *RequestValidator {
	v, err := NewRequestValidator(spec)
	if err != nil {
		panic(err)
	}
	return v
}

// Validate returns a normalized ACP error when the request does not match the spec.
func (v *RequestValidator) Validate(req *http.Request) *acp.Error {
	v.mu.Lock()
	defer v.mu.Unlock()

	valid, validationErrs := v.validator.ValidateHttpRequestSync(req)
	if valid {
		return nil
	}

	return convertValidationErrors(validationErrs)
}

func convertValidationErrors(errs []*validatorerrors.ValidationError) *acp.Error {
	errs = nonHeaderValidationErrors(errs)
	if len(errs) == 0 {
		return nil
	}

	code := acp.ErrorCode(acp.InvalidRequest)
	message := "request validation failed"
	status := http.StatusBadRequest

	var param string
	for _, err := range errs {
		if path, reason := firstSchemaDetail(err); path != "" || reason != "" {
			if path != "" {
				param = path
			}
			if reason != "" {
				message = reason
			}
			if isSemanticValidationError(err) {
				status = http.StatusUnprocessableEntity
			}
			break
		}

		if err.Reason != "" {
			message = err.Reason
			break
		}
		if err.Message != "" {
			message = err.Message
			break
		}
	}

	payload := acp.NewHTTPError(status, acp.InvalidRequest, code, message)
	if param != "" {
		payload.Param = &param
	}
	return payload
}

func nonHeaderValidationErrors(errs []*validatorerrors.ValidationError) []*validatorerrors.ValidationError {
	filtered := errs[:0]
	for _, err := range errs {
		if isHeaderValidationError(err) {
			continue
		}
		filtered = append(filtered, err)
	}
	return filtered
}

func isHeaderValidationError(err *validatorerrors.ValidationError) bool {
	return err.ValidationType == validatorhelpers.ParameterValidation &&
		err.ValidationSubType == validatorhelpers.ParameterValidationHeader
}

func firstSchemaDetail(err *validatorerrors.ValidationError) (string, string) {
	for _, schemaErr := range err.SchemaValidationErrors {
		if schemaErr == nil {
			continue
		}
		return schemaErr.FieldPath, schemaErr.Reason
	}
	return "", ""
}

func isSemanticValidationError(err *validatorerrors.ValidationError) bool {
	if err.ValidationType != validatorhelpers.RequestBodyValidation {
		return false
	}
	if strings.Contains(err.Reason, "cannot be decoded") {
		return false
	}
	if strings.Contains(err.Reason, "is empty but there is a schema defined") {
		return false
	}

	if len(err.SchemaValidationErrors) == 0 {
		return false
	}
	for _, schemaErr := range err.SchemaValidationErrors {
		if schemaErr == nil {
			continue
		}
		if strings.Contains(schemaErr.KeywordLocation, "/required") {
			return false
		}
		if strings.Contains(schemaErr.KeywordLocation, "/additionalProperties") {
			return false
		}
	}
	return true
}
