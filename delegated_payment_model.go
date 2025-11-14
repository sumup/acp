package acp

import "time"

// PaymentRequest mirrors the ACP DelegatePaymentRequest payload described in the spec:
// https://developers.openai.com/commerce/specs/payment.
type PaymentRequest struct {
	// Type of credential. The only accepted value is "CARD".
	PaymentMethod PaymentMethodCard `json:"payment_method" validate:"required"`
	// Use cases that the stored credential can be applied to.
	Allowance Allowance `json:"allowance" validate:"required"`
	// Address associated with the payment method.
	BillingAddress *DelegatedPaymentAddress `json:"billing_address,omitempty" validate:"omitempty"`
	// Arbitrary key/value pairs.
	Metadata map[string]string `json:"metadata" validate:"required,map_present"`
	// List of risk signals.
	RiskSignals []RiskSignal `json:"risk_signals" validate:"required,min=1,dive"`
}

// VaultToken is emitted by PSPs after tokenizing the delegated payment payload.
type VaultToken struct {
	// Unique vault token identifier vt_….
	ID string `json:"id" validate:"required"`
	// Time formatted as an RFC 3339 string.
	Created time.Time `json:"created" validate:"required"`
	// Arbitrary key/value pairs for correlation (e.g., source, merchant_id, idempotency_key).
	Metadata map[string]string `json:"metadata" validate:"omitempty"`
}

// PaymentMethodCard captures the delegated card credential.
type PaymentMethodCard struct {
	// The type of payment method used. Currently only card.
	Type PaymentMethodCardType `json:"type" validate:"required,eq=card"`
	// The type of card number. Network tokens are preferred with fallback to FPAN. See [PCI Scope] for more details.
	//
	// [PCI Scope]: https://developers.openai.com/commerce/guides/production#security-and-compliance
	CardNumberType PaymentMethodCardCardNumberType `json:"card_number_type" validate:"required,oneof=fpan network_token"`
	// Card number.
	Number string `json:"number" validate:"required"`
	// Expiry month.
	ExpMonth *string `json:"exp_month,omitempty" validate:"omitempty,len=2,numeric"`
	// Expiry year.
	ExpYear *string `json:"exp_year,omitempty" validate:"omitempty,len=4,numeric"`
	// Cardholder name.
	Name *string `json:"name,omitempty"`
	// Card CVC number.
	CVC *string `json:"cvc,omitempty" validate:"omitempty,numeric"`
	// In case of non-PAN, this is the original last 4 digits of the card for customer display.
	DisplayLast4 *string `json:"display_last4,omitempty" validate:"omitempty,len=4"`
	// Funding type of the card to display.
	DisplayCardFundingType PaymentMethodCardDisplayCardFundingType `json:"display_card_funding_type" validate:"required,oneof=credit debit prepaid"`
	// Brand of the card to display.
	//
	// Exapmple: "Visa", "amex", "discover"
	DisplayBrand *string `json:"display_brand,omitempty"`
	// If the card came via a digital wallet, what type of wallet.
	DisplayWalletType *string `json:"display_wallet_type,omitempty"`
	// Institution Identification Number (aka BIN). The first 6 digits on a card identifying the issuer.
	IIN *string `json:"iin,omitempty" validate:"omitempty,max=6"`
	// Cryptogram provided with network tokens.
	Cryptogram *string `json:"cryptogram,omitempty"`
	// Electronic Commerce Indicator / Security Level Indicator provided with network tokens.
	ECIValue *string `json:"eci_value,omitempty"`
	// Checks already performed on the card.
	ChecksPerformed []PaymentMethodCardChecksPerformed `json:"checks_performed,omitempty" validate:"omitempty,dive,required"`
	// Arbitrary key/value pairs.
	Metadata map[string]string `json:"metadata" validate:"required,map_present"`
}

// Allowance scopes token use per the spec.
type Allowance struct {
	// Current possible values: "one_time".
	Reason AllowanceReason `json:"reason" validate:"required,eq=one_time"`
	// Max amount the payment method can be charged for.
	MaxAmount int `json:"max_amount" validate:"required,gt=0"`
	// Currency.
	Currency string `json:"currency" validate:"required,currency"`
	// Reference to checkout_session_id.
	CheckoutSessionID string `json:"checkout_session_id" validate:"required"`
	// Merchant identifying descriptor.
	MerchantID string `json:"merchant_id" validate:"required"`
	// Time formatted as an RFC 3339 string.
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
}

// RiskSignal provides PSPs with fraud intelligence references.
type RiskSignal struct {
	// The type of risk signal.
	Type RiskSignalType `json:"type" validate:"required,oneof=card_testing"`
	// Action taken.
	Action RiskSignalAction `json:"action" validate:"required,oneof=manual_review authorized blocked"`
	// Details of the risk signal.
	Score int `json:"score" validate:"gte=0"`
}

// DelegatedPaymentAddress corresponds to the billing address object.
type DelegatedPaymentAddress struct {
	Name       string  `json:"name" validate:"required"`
	LineOne    string  `json:"line_one" validate:"required"`
	LineTwo    *string `json:"line_two,omitempty"`
	City       string  `json:"city" validate:"required"`
	State      string  `json:"state" validate:"required"`
	Country    string  `json:"country" validate:"required,len=2,uppercase"`
	PostalCode string  `json:"postal_code" validate:"required"`
}

type PaymentMethodCardType string

const (
	PaymentMethodCardTypeCard PaymentMethodCardType = "card"
)

type PaymentMethodCardCardNumberType string

const (
	PaymentMethodCardCardNumberTypeFPAN         PaymentMethodCardCardNumberType = "fpan"
	PaymentMethodCardCardNumberTypeNetworkToken PaymentMethodCardCardNumberType = "network_token"
)

type PaymentMethodCardDisplayCardFundingType string

const (
	PaymentMethodCardDisplayCardFundingTypeCredit  PaymentMethodCardDisplayCardFundingType = "credit"
	PaymentMethodCardDisplayCardFundingTypeDebit   PaymentMethodCardDisplayCardFundingType = "debit"
	PaymentMethodCardDisplayCardFundingTypePrepaid PaymentMethodCardDisplayCardFundingType = "prepaid"
)

type PaymentMethodCardChecksPerformed string

const (
	PaymentMethodCardChecksPerformedAVS  PaymentMethodCardChecksPerformed = "avs"
	PaymentMethodCardChecksPerformedCVV  PaymentMethodCardChecksPerformed = "cvv"
	PaymentMethodCardChecksPerformedANI  PaymentMethodCardChecksPerformed = "ani"
	PaymentMethodCardChecksPerformedAUTH PaymentMethodCardChecksPerformed = "auth0"
)

type AllowanceReason string

const (
	AllowanceReasonOneTime AllowanceReason = "one_time"
)

type RiskSignalType string

const (
	RiskSignalTypeCardTesting RiskSignalType = "card_testing"
)

type RiskSignalAction string

const (
	RiskSignalActionManualReview RiskSignalAction = "manual_review"
	RiskSignalActionAuthorized   RiskSignalAction = "authorized"
	RiskSignalActionBlocked      RiskSignalAction = "blocked"
)
