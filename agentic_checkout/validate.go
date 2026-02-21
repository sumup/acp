package agentic_checkout

import (
	"errors"
	"fmt"
)

// Validate performs additional runtime checks when requests are built in Go
// (JSON unmarshalling already enforces schema constraints).
func (r CheckoutSessionCreateRequest) Validate() error {
	if len(r.LineItems) == 0 {
		return errors.New("line_items must contain at least one entry")
	}
	for i, item := range r.LineItems {
		if item.Id == "" {
			return fmt.Errorf("line_items[%d]: id is required", i)
		}
	}
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	if r.Buyer != nil && r.Buyer.Email == "" {
		return errors.New("buyer.email is required")
	}
	return nil
}

// Validate performs additional runtime checks when requests are built in Go.
func (r CheckoutSessionUpdateRequest) Validate() error {
	if r.LineItems != nil {
		for i, item := range *r.LineItems {
			if item.Id == "" {
				return fmt.Errorf("line_items[%d]: id is required", i)
			}
		}
	}
	if r.Buyer != nil && r.Buyer.Email == "" {
		return errors.New("buyer.email is required")
	}
	return nil
}

// Validate performs additional runtime checks when requests are built in Go.
func (r CheckoutSessionCompleteRequest) Validate() error {
	if len(r.PaymentData.union) == 0 &&
		r.PaymentData.HandlerId == nil &&
		r.PaymentData.Instrument == nil &&
		r.PaymentData.PurchaseOrderNumber == nil {
		return errors.New("payment_data is required")
	}
	return nil
}

// Validate performs additional runtime checks when requests are built in Go.
func (r CancelSessionRequest) Validate() error {
	// OpenAPI-based unmarshalling already constrains metadata item values.
	return nil
}
