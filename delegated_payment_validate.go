package acp

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	currencyPattern = regexp.MustCompile(`^[a-z]{3}$`)
	validate        = newValidator()
)

// Validate ensures the request complies with the ACP Delegate Payment spec by
// running go-playground/validator rules plus custom constraints.
func (r PaymentRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return normalizeValidationError(err)
	}
	return nil
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})

	if err := v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return currencyPattern.MatchString(value)
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterValidation("map_present", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.Map {
			return false
		}
		return !fl.Field().IsNil()
	}); err != nil {
		panic(err)
	}

	return v
}

func normalizeValidationError(err error) error {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return err
	}
	first := validationErrs[0]
	fieldPath := jsonPath(first)
	message := validationMessage(first)
	return fmt.Errorf("%s %s", fieldPath, message)
}

func jsonPath(fe validator.FieldError) string {
	path := fe.Namespace()
	if idx := strings.Index(path, "."); idx >= 0 {
		path = path[idx+1:]
	}
	if path == "" {
		return fe.Field()
	}
	return path
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "map_present":
		return "must be provided"
	case "min":
		return fmt.Sprintf("must have at least %s entries", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("cannot exceed %s characters", fe.Param())
	case "numeric":
		return "must contain digits only"
	case "gt":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "eq":
		return fmt.Sprintf("must equal %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", strings.ReplaceAll(fe.Param(), " ", ", "))
	case "currency":
		return "must be a lowercase 3-letter ISO-4217 code"
	case "uppercase":
		return "must be uppercase"
	default:
		return fmt.Sprintf("failed validation: %s", fe.Tag())
	}
}
