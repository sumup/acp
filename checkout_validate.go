package acp

import (
	"errors"
	"fmt"
)

// Validate ensures CheckoutSessionCreateRequest satisfies required schema constraints.
func (r CheckoutSessionCreateRequest) Validate() error {
	if len(r.LineItems) == 0 {
		return errors.New("line_items must contain at least one entry")
	}
	for i, item := range r.LineItems {
		if item.ID == "" {
			return fmt.Errorf("line_items[%d]: id is required", i)
		}
	}
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	if r.Buyer != nil {
		if r.Buyer.Email == "" {
			return errors.New("buyer.email is required")
		}
	}
	if err := validateAffiliateAttribution(r.AffiliateAttribution); err != nil {
		return err
	}
	return nil
}

// Validate ensures CheckoutSessionUpdateRequest maintains schema constraints.
func (r CheckoutSessionUpdateRequest) Validate() error {
	if r.LineItems != nil {
		for i, item := range *r.LineItems {
			if item.ID == "" {
				return fmt.Errorf("line_items[%d]: id is required", i)
			}
		}
	}
	if r.Buyer != nil {
		if r.Buyer.Email == "" {
			return errors.New("buyer.email is required")
		}
	}
	return nil
}

// Validate ensures CheckoutSessionCompleteRequest satisfies payment requirements.
func (r CheckoutSessionCompleteRequest) Validate() error {
	hasTokenProvider := r.PaymentData.Token != nil && *r.PaymentData.Token != "" &&
		r.PaymentData.Provider != nil && *r.PaymentData.Provider != ""
	hasPurchaseOrder := r.PaymentData.PurchaseOrderNumber != nil && *r.PaymentData.PurchaseOrderNumber != ""
	if !hasTokenProvider && !hasPurchaseOrder {
		return errors.New("payment_data requires token+provider or purchase_order_number")
	}
	if r.PaymentData.Token != nil && *r.PaymentData.Token != "" && (r.PaymentData.Provider == nil || *r.PaymentData.Provider == "") {
		return errors.New("payment_data.provider is required when payment_data.token is set")
	}
	if r.PaymentData.Provider != nil && *r.PaymentData.Provider != "" && (r.PaymentData.Token == nil || *r.PaymentData.Token == "") {
		return errors.New("payment_data.token is required when payment_data.provider is set")
	}
	if r.AuthenticationResult != nil {
		if r.AuthenticationResult.Outcome == "" {
			return errors.New("authentication_result.outcome is required")
		}
		if requiresAuthenticationOutcomeDetails(r.AuthenticationResult.Outcome) {
			details := r.AuthenticationResult.OutcomeDetails
			if details == nil {
				return errors.New("authentication_result.outcome_details is required when outcome is authenticated, informational, or attempt_acknowledged")
			}
			if details.ThreeDsCryptogram == "" {
				return errors.New("authentication_result.outcome_details.three_ds_cryptogram is required")
			}
			if details.ElectronicCommerceIndicator == "" {
				return errors.New("authentication_result.outcome_details.electronic_commerce_indicator is required")
			}
			if details.TransactionId == "" {
				return errors.New("authentication_result.outcome_details.transaction_id is required")
			}
			if details.Version == "" {
				return errors.New("authentication_result.outcome_details.version is required")
			}
		}
	}
	if err := validateAffiliateAttribution(r.AffiliateAttribution); err != nil {
		return err
	}
	return nil
}

func requiresAuthenticationOutcomeDetails(outcome AuthenticationOutcome) bool {
	switch outcome {
	case AuthenticationOutcomeAuthenticated, AuthenticationOutcomeInformational, AuthenticationOutcomeAttemptAcknowledged:
		return true
	default:
		return false
	}
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

func validateAffiliateAttribution(attribution *AffiliateAttribution) error {
	if attribution == nil {
		return nil
	}
	if attribution.Provider == "" {
		return errors.New("affiliate_attribution.provider is required")
	}
	hasToken := attribution.Token != nil && *attribution.Token != ""
	hasPublisherID := attribution.PublisherID != nil && *attribution.PublisherID != ""
	if !hasToken && !hasPublisherID {
		return errors.New("affiliate_attribution requires token or publisher_id")
	}
	return nil
}
