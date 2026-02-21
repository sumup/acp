package delegate_payment

import "errors"

// Validate performs additional runtime checks when requests are built in Go
// (JSON unmarshalling already enforces schema constraints).
func (r DelegatePaymentRequest) Validate() error {
	if r.Allowance.Currency == "" {
		return errors.New("allowance.currency is required")
	}
	if len(r.RiskSignals) == 0 {
		return errors.New("risk_signals must contain at least one entry")
	}
	if r.Metadata == nil {
		return errors.New("metadata is required")
	}
	if r.PaymentMethod.Metadata == nil {
		return errors.New("payment_method.metadata is required")
	}
	return nil
}
