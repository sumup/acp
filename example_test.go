package acp_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/sumup/acp"
)

func ExampleRequestContextFromRequest() {
	req := httptest.NewRequest("POST", "/checkout_sessions", nil)
	req.Header.Set("API-Version", acp.APIVersion)
	req.Header.Set("Authorization", "Bearer api_key_123")
	req.Header.Set("Idempotency-Key", "idem_123")

	requestContext, err := acp.RequestContextFromRequest(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(requestContext.APIVersion, requestContext.IdempotencyKey)
	// Output: 2026-04-17 idem_123
}

func ExampleNewHTTPError() {
	err := acp.NewHTTPError(
		409,
		acp.InvalidRequest,
		acp.IdempotencyConflict,
		"idempotency key was already used with different parameters",
		acp.WithOffendingParam("$.line_items"),
	)

	fmt.Println(err.StatusCode(), err.Type, err.Code, *err.Param)
	// Output: 409 invalid_request idempotency_conflict $.line_items
}
