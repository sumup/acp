package acp

import (
	"errors"
	"fmt"
)

// Validate ensures CheckoutSessionCreateRequest satisfies required schema constraints.
func (r CheckoutSessionCreateRequest) Validate() error {
	if len(r.Items) == 0 {
		return errors.New("items must contain at least one entry")
	}
	for i, item := range r.Items {
		if item.ID == "" {
			return fmt.Errorf("items[%d]: id is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("items[%d]: quantity must be positive", i)
		}
	}
	if r.Buyer != nil {
		if r.Buyer.FirstName == "" || r.Buyer.LastName == "" || string(r.Buyer.Email) == "" {
			return errors.New("buyer requires first_name, last_name, and email")
		}
	}
	return nil
}

// Validate ensures CheckoutSessionUpdateRequest maintains schema constraints.
func (r CheckoutSessionUpdateRequest) Validate() error {
	if r.Items != nil {
		for i, item := range *r.Items {
			if item.ID == "" {
				return fmt.Errorf("items[%d]: id is required", i)
			}
			if item.Quantity <= 0 {
				return fmt.Errorf("items[%d]: quantity must be positive", i)
			}
		}
	}
	if r.Buyer != nil {
		if r.Buyer.FirstName == "" || r.Buyer.LastName == "" || string(r.Buyer.Email) == "" {
			return errors.New("buyer requires first_name, last_name, and email")
		}
	}
	return nil
}

// Validate ensures CheckoutSessionCompleteRequest satisfies payment requirements.
func (r CheckoutSessionCompleteRequest) Validate() error {
	if r.PaymentData.Token == "" {
		return errors.New("payment_data.token is required")
	}
	if r.PaymentData.Provider == "" {
		return errors.New("payment_data.provider is required")
	}
	return nil
}

// Validate ensures CancelSessionRequest conforms to the ACP schema.
func (r CancelSessionRequest) Validate() error {
	if r.IntentTrace == nil {
		return nil
	}
	if r.IntentTrace.ReasonCode == "" {
		return errors.New("intent_trace.reason_code is required")
	}
	if r.IntentTrace.TraceSummary != nil && len(*r.IntentTrace.TraceSummary) > 500 {
		return errors.New("intent_trace.trace_summary must be at most 500 characters")
	}
	if r.IntentTrace.Metadata != nil {
		for key, value := range r.IntentTrace.Metadata {
			switch value.(type) {
			case string, bool, float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				continue
			default:
				return fmt.Errorf("intent_trace.metadata[%s] must be string, number, or boolean", key)
			}
		}
	}
	return nil
}
