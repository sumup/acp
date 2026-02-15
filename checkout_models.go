package acp

import (
	"encoding/json"
	"time"

	"github.com/oapi-codegen/runtime"
)

// CheckoutSessionStatus defines model for CheckoutSessionBase.Status.
type CheckoutSessionStatus string

// Defines values for CheckoutSessionBaseStatus.
const (
	CheckoutSessionStatusIncomplete             CheckoutSessionStatus = "incomplete"
	CheckoutSessionStatusCanceled               CheckoutSessionStatus = "canceled"
	CheckoutSessionStatusCompleted              CheckoutSessionStatus = "completed"
	CheckoutSessionStatusCompleteInProgress     CheckoutSessionStatus = "complete_in_progress"
	CheckoutSessionStatusInProgress             CheckoutSessionStatus = "in_progress"
	CheckoutSessionStatusNotReadyForPayment     CheckoutSessionStatus = "not_ready_for_payment"
	CheckoutSessionStatusRequiresEscalation     CheckoutSessionStatus = "requires_escalation"
	CheckoutSessionStatusAuthenticationRequired CheckoutSessionStatus = "authentication_required"
	CheckoutSessionStatusReadyForPayment        CheckoutSessionStatus = "ready_for_payment"
	CheckoutSessionStatusPendingApproval        CheckoutSessionStatus = "pending_approval"
	CheckoutSessionStatusExpired                CheckoutSessionStatus = "expired"
)

// LinkType defines model for Link.Type.
type LinkType string

// Defines values for LinkType.
const (
	PrivacyPolicy LinkType = "privacy_policy"
	ReturnPolicy  LinkType = "return_policy"
	// Deprecated: removed in ACP 2025-12-11; attach policies to line items in marketplace scenarios.
	SellerShopPolicies LinkType = "seller_shop_policies"
	ShippingPolicy     LinkType = "shipping_policy"
	ContactUs          LinkType = "contact_us"
	AboutUs            LinkType = "about_us"
	FAQ                LinkType = "faq"
	Support            LinkType = "support"
	TermsOfUse         LinkType = "terms_of_use"
)

// DisclosureType defines model for Disclosure.Type.
type DisclosureType string

// Defines values for DisclosureType.
const (
	DisclosureTypeDisclaimer DisclosureType = "disclaimer"
)

// DisclosureContentType defines model for Disclosure.ContentType.
type DisclosureContentType string

// Defines values for DisclosureContentType.
const (
	DisclosureContentTypeMarkdown DisclosureContentType = "markdown"
	DisclosureContentTypePlain    DisclosureContentType = "plain"
)

// MessageErrorCode defines model for MessageError.Code.
type MessageErrorCode string

// Defines values for MessageErrorCode.
const (
	Invalid                 MessageErrorCode = "invalid"
	Missing                 MessageErrorCode = "missing"
	OutOfStock              MessageErrorCode = "out_of_stock"
	PaymentDeclined         MessageErrorCode = "payment_declined"
	Requires3ds             MessageErrorCode = "requires_3ds"
	RequiresSignIn          MessageErrorCode = "requires_sign_in"
	LowStock                MessageErrorCode = "low_stock"
	QuantityExceeded        MessageErrorCode = "quantity_exceeded"
	CouponInvalid           MessageErrorCode = "coupon_invalid"
	CouponExpired           MessageErrorCode = "coupon_expired"
	MinimumNotMet           MessageErrorCode = "minimum_not_met"
	MaximumExceeded         MessageErrorCode = "maximum_exceeded"
	RegionRestricted        MessageErrorCode = "region_restricted"
	AgeVerificationRequired MessageErrorCode = "age_verification_required"
	ApprovalRequired        MessageErrorCode = "approval_required"
	Unsupported             MessageErrorCode = "unsupported"
	NotFound                MessageErrorCode = "not_found"
	Conflict                MessageErrorCode = "conflict"
	RateLimited             MessageErrorCode = "rate_limited"
	Expired                 MessageErrorCode = "expired"
	InterventionRequired    MessageErrorCode = "intervention_required"
)

// MessageErrorContentType defines model for MessageError.ContentType.
type MessageErrorContentType string

// Defines values for MessageErrorContentType.
const (
	MessageErrorContentTypeMarkdown MessageErrorContentType = "markdown"
	MessageErrorContentTypePlain    MessageErrorContentType = "plain"
)

// MessageSeverity defines model for MessageInfo/MessageWarning/MessageError.Severity.
type MessageSeverity string

// Defines values for MessageSeverity.
const (
	MessageSeverityInfo     MessageSeverity = "info"
	MessageSeverityLow      MessageSeverity = "low"
	MessageSeverityMedium   MessageSeverity = "medium"
	MessageSeverityHigh     MessageSeverity = "high"
	MessageSeverityCritical MessageSeverity = "critical"
)

// MessageResolution defines who is responsible for resolving a message.
type MessageResolution string

// Defines values for MessageResolution.
const (
	MessageResolutionRecoverable         MessageResolution = "recoverable"
	MessageResolutionRequiresBuyerInput  MessageResolution = "requires_buyer_input"
	MessageResolutionRequiresBuyerReview MessageResolution = "requires_buyer_review"
)

// MessageWarningCode defines model for MessageWarning.Code.
type MessageWarningCode string

// Defines values for MessageWarningCode.
const (
	MessageWarningLowStock            MessageWarningCode = "low_stock"
	MessageWarningHighDemand          MessageWarningCode = "high_demand"
	MessageWarningShippingDelay       MessageWarningCode = "shipping_delay"
	MessageWarningPriceChange         MessageWarningCode = "price_change"
	MessageWarningExpiringPromotion   MessageWarningCode = "expiring_promotion"
	MessageWarningLimitedAvailability MessageWarningCode = "limited_availability"
)

// MessageInfoContentType defines model for MessageInfo.ContentType.
type MessageInfoContentType string

// Defines values for MessageInfoContentType.
const (
	MessageInfoContentTypeMarkdown MessageInfoContentType = "markdown"
	MessageInfoContentTypePlain    MessageInfoContentType = "plain"
)

// IntentTraceReasonCode defines model for IntentTrace.ReasonCode.
type IntentTraceReasonCode string

// Defines values for IntentTraceReasonCode.
const (
	IntentTraceReasonCodePriceSensitivity IntentTraceReasonCode = "price_sensitivity"
	IntentTraceReasonCodeShippingCost     IntentTraceReasonCode = "shipping_cost"
	IntentTraceReasonCodeShippingSpeed    IntentTraceReasonCode = "shipping_speed"
	IntentTraceReasonCodeProductFit       IntentTraceReasonCode = "product_fit"
	IntentTraceReasonCodeTrustSecurity    IntentTraceReasonCode = "trust_security"
	IntentTraceReasonCodeReturnsPolicy    IntentTraceReasonCode = "returns_policy"
	IntentTraceReasonCodePaymentOptions   IntentTraceReasonCode = "payment_options"
	IntentTraceReasonCodeComparison       IntentTraceReasonCode = "comparison"
	IntentTraceReasonCodeTimingDeferred   IntentTraceReasonCode = "timing_deferred"
	IntentTraceReasonCodeOther            IntentTraceReasonCode = "other"
)

// BuyerAccountType defines model for Buyer.AccountType.
type BuyerAccountType string

// Defines values for BuyerAccountType.
const (
	BuyerAccountTypeGuest      BuyerAccountType = "guest"
	BuyerAccountTypeRegistered BuyerAccountType = "registered"
	BuyerAccountTypeBusiness   BuyerAccountType = "business"
)

// BuyerAuthenticationStatus defines model for Buyer.AuthenticationStatus.
type BuyerAuthenticationStatus string

// Defines values for BuyerAuthenticationStatus.
const (
	BuyerAuthenticationStatusAuthenticated  BuyerAuthenticationStatus = "authenticated"
	BuyerAuthenticationStatusGuest          BuyerAuthenticationStatus = "guest"
	BuyerAuthenticationStatusRequiresSignIn BuyerAuthenticationStatus = "requires_signin"
)

// PaymentMethodType defines model for PaymentMethod.Type.
type PaymentMethodType string

// Defines values for PaymentMethodType.
const (
	PaymentMethodTypeCard PaymentMethodType = "card"
)

// PaymentCardNetwork defines model for PaymentMethod.SupportedCardNetworks items.
type PaymentCardNetwork string

// Defines values for PaymentCardNetwork.
const (
	PaymentCardNetworkAmex       PaymentCardNetwork = "amex"
	PaymentCardNetworkDiscover   PaymentCardNetwork = "discover"
	PaymentCardNetworkMastercard PaymentCardNetwork = "mastercard"
	PaymentCardNetworkVisa       PaymentCardNetwork = "visa"
)

// WeightUnit defines model for WeightInfo.Unit.
type WeightUnit string

// Defines values for WeightUnit.
const (
	WeightUnitGrams     WeightUnit = "g"
	WeightUnitKilograms WeightUnit = "kg"
	WeightUnitOunces    WeightUnit = "oz"
	WeightUnitPounds    WeightUnit = "lb"
)

// DimensionsUnit defines model for DimensionsInfo.Unit.
type DimensionsUnit string

// Defines values for DimensionsUnit.
const (
	DimensionsUnitCentimeters DimensionsUnit = "cm"
	DimensionsUnitInches      DimensionsUnit = "in"
)

// DiscountDetailType defines model for DiscountDetail.Type.
type DiscountDetailType string

// Defines values for DiscountDetailType.
const (
	DiscountDetailTypePercentage DiscountDetailType = "percentage"
	DiscountDetailTypeFixed      DiscountDetailType = "fixed"
	DiscountDetailTypeBogo       DiscountDetailType = "bogo"
	DiscountDetailTypeVolume     DiscountDetailType = "volume"
)

// DiscountDetailSource defines model for DiscountDetail.Source.
type DiscountDetailSource string

// Defines values for DiscountDetailSource.
const (
	DiscountDetailSourceCoupon    DiscountDetailSource = "coupon"
	DiscountDetailSourceAutomatic DiscountDetailSource = "automatic"
	DiscountDetailSourceLoyalty   DiscountDetailSource = "loyalty"
)

// AvailabilityStatus defines model for LineItem.AvailabilityStatus.
type AvailabilityStatus string

// Defines values for AvailabilityStatus.
const (
	AvailabilityStatusInStock    AvailabilityStatus = "in_stock"
	AvailabilityStatusLowStock   AvailabilityStatus = "low_stock"
	AvailabilityStatusOutOfStock AvailabilityStatus = "out_of_stock"
	AvailabilityStatusBackorder  AvailabilityStatus = "backorder"
	AvailabilityStatusPreOrder   AvailabilityStatus = "pre_order"
)

// FulfillmentOptionType defines model for SelectedFulfillmentOptions.Type.
type FulfillmentOptionType string

// Defines values for FulfillmentOptionType.
const (
	FulfillmentOptionTypeShipping      FulfillmentOptionType = "shipping"
	FulfillmentOptionTypeDigital       FulfillmentOptionType = "digital"
	FulfillmentOptionTypePickup        FulfillmentOptionType = "pickup"
	FulfillmentOptionTypeLocalDelivery FulfillmentOptionType = "local_delivery"
)

// FulfillmentPickupType defines model for FulfillmentOptionPickup.PickupType.
type FulfillmentPickupType string

// Defines values for FulfillmentPickupType.
const (
	FulfillmentPickupTypeInStore  FulfillmentPickupType = "in_store"
	FulfillmentPickupTypeCurbside FulfillmentPickupType = "curbside"
	FulfillmentPickupTypeLocker   FulfillmentPickupType = "locker"
)

// PaymentTerms defines model for PaymentData.PaymentTerms.
type PaymentTerms string

// Defines values for PaymentTerms.
const (
	PaymentTermsImmediate PaymentTerms = "immediate"
	PaymentTermsNet15     PaymentTerms = "net_15"
	PaymentTermsNet30     PaymentTerms = "net_30"
	PaymentTermsNet60     PaymentTerms = "net_60"
	PaymentTermsNet90     PaymentTerms = "net_90"
)

// CheckoutOrderStatus is the order-level lifecycle status returned in Order.
// It aligns with the webhook/order superset used for post-purchase tracking.
type CheckoutOrderStatus string

// Defines values for CheckoutOrderStatus.
const (
	// CheckoutOrderStatusCreated indicates an order record has been created.
	CheckoutOrderStatusCreated CheckoutOrderStatus = "created"
	// CheckoutOrderStatusConfirmed indicates the order has been accepted/confirmed.
	CheckoutOrderStatusConfirmed CheckoutOrderStatus = "confirmed"
	// CheckoutOrderStatusManualReview indicates the order is waiting for manual review.
	CheckoutOrderStatusManualReview CheckoutOrderStatus = "manual_review"
	// CheckoutOrderStatusProcessing indicates fulfillment/processing has started.
	CheckoutOrderStatusProcessing CheckoutOrderStatus = "processing"
	// CheckoutOrderStatusShipped indicates at least one shipment has been handed off.
	CheckoutOrderStatusShipped CheckoutOrderStatus = "shipped"
	// CheckoutOrderStatusDelivered indicates fulfillment is complete and delivered.
	CheckoutOrderStatusDelivered CheckoutOrderStatus = "delivered"
	// CheckoutOrderStatusCanceled indicates the order has been canceled.
	CheckoutOrderStatusCanceled CheckoutOrderStatus = "canceled"
)

// TotalType identifies a monetary component in a totals breakdown.
type TotalType string

// Defines values for TotalType.
const (
	TotalTypeDiscount        TotalType = "discount"
	TotalTypeFee             TotalType = "fee"
	TotalTypeFulfillment     TotalType = "fulfillment"
	TotalTypeGiftWrap        TotalType = "gift_wrap"
	TotalTypeItemsBaseAmount TotalType = "items_base_amount"
	TotalTypeItemsDiscount   TotalType = "items_discount"
	TotalTypeStoreCredit     TotalType = "store_credit"
	TotalTypeSubtotal        TotalType = "subtotal"
	TotalTypeTax             TotalType = "tax"
	TotalTypeTip             TotalType = "tip"
	TotalTypeTotal           TotalType = "total"
	// TotalTypeAmountRefunded tracks cumulative refunded amount post-purchase.
	TotalTypeAmountRefunded TotalType = "amount_refunded"
)

// Address defines model for Address.
type Address struct {
	Name       string  `json:"name"`
	LineOne    string  `json:"line_one"`
	LineTwo    *string `json:"line_two,omitempty"`
	PostalCode string  `json:"postal_code"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
}

// CompanyInfo defines model for Buyer.Company.
type CompanyInfo struct {
	Name       string  `json:"name"`
	TaxID      *string `json:"tax_id,omitempty"`
	Department *string `json:"department,omitempty"`
	CostCenter *string `json:"cost_center,omitempty"`
}

// LoyaltyInfo defines model for Buyer.Loyalty.
type LoyaltyInfo struct {
	Tier          *string    `json:"tier,omitempty"`
	PointsBalance *int       `json:"points_balance,omitempty"`
	MemberSince   *time.Time `json:"member_since,omitempty"`
}

// TaxExemption defines model for Buyer.TaxExemption.
type TaxExemption struct {
	CertificateID   string     `json:"certificate_id"`
	CertificateType string     `json:"certificate_type"`
	ExemptRegions   []string   `json:"exempt_regions,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
}

// Buyer defines model for Buyer.
type Buyer struct {
	Email                string                     `json:"email"`
	FirstName            *string                    `json:"first_name,omitempty"`
	LastName             *string                    `json:"last_name,omitempty"`
	FullName             *string                    `json:"full_name,omitempty"`
	PhoneNumber          *string                    `json:"phone_number,omitempty"`
	CustomerID           *string                    `json:"customer_id,omitempty"`
	AccountType          *BuyerAccountType          `json:"account_type,omitempty"`
	AuthenticationStatus *BuyerAuthenticationStatus `json:"authentication_status,omitempty"`
	Company              *CompanyInfo               `json:"company,omitempty"`
	Loyalty              *LoyaltyInfo               `json:"loyalty,omitempty"`
	TaxExemption         *TaxExemption              `json:"tax_exemption,omitempty"`
}

// FulfillmentDetails defines model for FulfillmentDetails.
type FulfillmentDetails struct {
	Name        *string  `json:"name,omitempty"`
	PhoneNumber *string  `json:"phone_number,omitempty"`
	Email       *string  `json:"email,omitempty"`
	Address     *Address `json:"address,omitempty"`
}

// CheckoutSession defines model for CheckoutSession.
type CheckoutSession struct {
	ID              string           `json:"id"`
	Protocol        *ProtocolVersion `json:"protocol,omitempty"`
	Buyer           *Buyer           `json:"buyer,omitempty"`
	PaymentProvider *PaymentProvider `json:"payment_provider,omitempty"`
	// Payment contains payment provider response data (if available).
	Payment  *PaymentResponse      `json:"payment,omitempty"`
	Status   CheckoutSessionStatus `json:"status"`
	Currency string                `json:"currency"`
	// PresentmentCurrency is the currency used for buyer-facing pricing.
	PresentmentCurrency *string `json:"presentment_currency,omitempty"`
	// ExchangeRate is the conversion rate from Currency to PresentmentCurrency.
	ExchangeRate *float64 `json:"exchange_rate,omitempty"`
	// ExchangeRateTimestamp is the time the exchange rate was captured.
	ExchangeRateTimestamp *time.Time          `json:"exchange_rate_timestamp,omitempty"`
	Locale                *string             `json:"locale,omitempty"`
	Timezone              *string             `json:"timezone,omitempty"`
	LineItems             []LineItem          `json:"line_items"`
	FulfillmentDetails    *FulfillmentDetails `json:"fulfillment_details,omitempty"`
	FulfillmentOptions    []FulfillmentOption `json:"fulfillment_options"`
	// SelectedFulfillmentOptions lists fulfillment choices by item group.
	SelectedFulfillmentOptions []SelectedFulfillmentOptions `json:"selected_fulfillment_options,omitempty"`
	// FulfillmentGroups allow splitting items across multiple destinations.
	FulfillmentGroups []FulfillmentGroup `json:"fulfillment_groups,omitempty"`
	Totals            []Total            `json:"totals"`
	Messages          []Message          `json:"messages"`
	Links             []Link             `json:"links"`
	// AuthenticationMetadata is seller-provided authentication metadata for 3DS flows.
	AuthenticationMetadata *AuthenticationMetadata `json:"authentication_metadata,omitempty"`
	// AvailablePromotions lists promotions applicable to the current cart.
	AvailablePromotions []AvailablePromotion `json:"available_promotions,omitempty"`
	CreatedAt           *time.Time           `json:"created_at,omitempty"`
	UpdatedAt           *time.Time           `json:"updated_at,omitempty"`
	ExpiresAt           *time.Time           `json:"expires_at,omitempty"`
	ContinueURL         *string              `json:"continue_url,omitempty"`
	Metadata            map[string]any       `json:"metadata,omitempty"`
	QuoteID             *string              `json:"quote_id,omitempty"`
	QuoteExpiresAt      *time.Time           `json:"quote_expires_at,omitempty"`
}

// FulfillmentOption defines model for CheckoutSessionBase.fulfillment_options.Item.
type FulfillmentOption struct {
	union json.RawMessage
}

// Message defines model for CheckoutSessionBase.messages.Item.
type Message struct {
	union json.RawMessage
}

// CheckoutSessionCompleteRequest defines model for CheckoutSessionCompleteRequest.
type CheckoutSessionCompleteRequest struct {
	Buyer       *Buyer      `json:"buyer,omitempty"`
	PaymentData PaymentData `json:"payment_data"`
	// AuthenticationResult is agent-provided 3DS authentication results for card payments.
	AuthenticationResult *AuthenticationResult `json:"authentication_result,omitempty"`
	// AffiliateAttribution contains optional attribution data for crediting third-party publishers.
	AffiliateAttribution *AffiliateAttribution `json:"affiliate_attribution,omitempty"`
	// RiskSignals captures client-provided signals for fraud analysis.
	RiskSignals *RiskSignals `json:"risk_signals,omitempty"`
}

// CheckoutSessionCreateRequest defines model for CheckoutSessionCreateRequest.
type CheckoutSessionCreateRequest struct {
	Buyer *Buyer `json:"buyer,omitempty"`
	// LineItems is the list of requested items for the checkout session.
	LineItems []Item `json:"line_items"`
	// Currency is the ISO 4217 currency code for this checkout.
	Currency           string              `json:"currency"`
	FulfillmentDetails *FulfillmentDetails `json:"fulfillment_details,omitempty"`
	FulfillmentGroups  []FulfillmentGroup  `json:"fulfillment_groups,omitempty"`
	// AffiliateAttribution contains optional attribution data for crediting third-party publishers.
	AffiliateAttribution *AffiliateAttribution `json:"affiliate_attribution,omitempty"`
	Coupons              []string              `json:"coupons,omitempty"`
	Locale               *string               `json:"locale,omitempty"`
	Timezone             *string               `json:"timezone,omitempty"`
	QuoteID              *string               `json:"quote_id,omitempty"`
	Metadata             map[string]any        `json:"metadata,omitempty"`
}

// CheckoutSessionUpdateRequest defines model for CheckoutSessionUpdateRequest.
type CheckoutSessionUpdateRequest struct {
	Buyer                      *Buyer                       `json:"buyer,omitempty"`
	LineItems                  *[]Item                      `json:"line_items,omitempty"`
	FulfillmentDetails         *FulfillmentDetails          `json:"fulfillment_details,omitempty"`
	FulfillmentGroups          []FulfillmentGroup           `json:"fulfillment_groups,omitempty"`
	SelectedFulfillmentOptions []SelectedFulfillmentOptions `json:"selected_fulfillment_options,omitempty"`
	Coupons                    []string                     `json:"coupons,omitempty"`
}

// SessionWithOrder defines model for SessionWithOrder.
type SessionWithOrder struct {
	CheckoutSession
	Order Order `json:"order"`
}

// FulfillmentOptionDigital defines model for FulfillmentOptionDigital.
type FulfillmentOptionDigital struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Totals      []Total `json:"totals"`
	Type        string  `json:"type"`
}

// FulfillmentOptionShipping defines model for FulfillmentOptionShipping.
type FulfillmentOptionShipping struct {
	ID                   string     `json:"id"`
	Title                string     `json:"title"`
	Description          *string    `json:"description,omitempty"`
	Carrier              *string    `json:"carrier,omitempty"`
	EarliestDeliveryTime *time.Time `json:"earliest_delivery_time,omitempty"`
	LatestDeliveryTime   *time.Time `json:"latest_delivery_time,omitempty"`
	Totals               []Total    `json:"totals"`
	Type                 string     `json:"type"`
}

// FulfillmentPickupLocation describes pickup location metadata.
type FulfillmentPickupLocation struct {
	Name         string  `json:"name"`
	Address      Address `json:"address"`
	Phone        *string `json:"phone,omitempty"`
	Instructions *string `json:"instructions,omitempty"`
}

// FulfillmentOptionPickup defines model for FulfillmentOptionPickup.
type FulfillmentOptionPickup struct {
	Type        string                    `json:"type"`
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Description *string                   `json:"description,omitempty"`
	Location    FulfillmentPickupLocation `json:"location"`
	PickupType  *FulfillmentPickupType    `json:"pickup_type,omitempty"`
	ReadyBy     *time.Time                `json:"ready_by,omitempty"`
	PickupBy    *time.Time                `json:"pickup_by,omitempty"`
	Totals      []Total                   `json:"totals"`
}

// FulfillmentDeliveryWindow defines model for FulfillmentOptionLocalDelivery.DeliveryWindow.
type FulfillmentDeliveryWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// FulfillmentServiceArea defines model for FulfillmentOptionLocalDelivery.ServiceArea.
type FulfillmentServiceArea struct {
	RadiusMiles      *float64 `json:"radius_miles,omitempty"`
	CenterPostalCode *string  `json:"center_postal_code,omitempty"`
}

// FulfillmentOptionLocalDelivery defines model for FulfillmentOptionLocalDelivery.
type FulfillmentOptionLocalDelivery struct {
	Type           string                     `json:"type"`
	ID             string                     `json:"id"`
	Title          string                     `json:"title"`
	Description    *string                    `json:"description,omitempty"`
	DeliveryWindow *FulfillmentDeliveryWindow `json:"delivery_window,omitempty"`
	ServiceArea    *FulfillmentServiceArea    `json:"service_area,omitempty"`
	Totals         []Total                    `json:"totals"`
}

// SelectedFulfillmentOptions defines model for SelectedFulfillmentOptions.
type SelectedFulfillmentOptions struct {
	Type     FulfillmentOptionType `json:"type"`
	OptionID string                `json:"option_id"`
	ItemIDs  []string              `json:"item_ids"`
}

// AvailablePromotion defines model for AvailablePromotion.
type AvailablePromotion struct {
	Code         string  `json:"code"`
	Description  string  `json:"description"`
	Requirements *string `json:"requirements,omitempty"`
}

// GiftWrap defines model for GiftWrap.
type GiftWrap struct {
	Enabled bool    `json:"enabled"`
	Style   *string `json:"style,omitempty"`
	Charge  *int    `json:"charge,omitempty"`
}

// SplitPayment defines model for SplitPayment.
type SplitPayment struct {
	Amount int `json:"amount"`
}

// FulfillmentGroup defines model for FulfillmentGroup.
type FulfillmentGroup struct {
	ID                 string                `json:"id"`
	ItemIDs            []string              `json:"item_ids"`
	DestinationType    FulfillmentOptionType `json:"destination_type"`
	FulfillmentDetails *FulfillmentDetails   `json:"fulfillment_details,omitempty"`
	LocationID         *string               `json:"location_id,omitempty"`
	Instructions       *string               `json:"instructions,omitempty"`
}

// EstimatedDelivery defines model for EstimatedDelivery.
type EstimatedDelivery struct {
	Earliest time.Time `json:"earliest"`
	Latest   time.Time `json:"latest"`
}

// OrderConfirmation defines model for OrderConfirmation.
type OrderConfirmation struct {
	ConfirmationNumber    *string `json:"confirmation_number,omitempty"`
	ConfirmationEmailSent *bool   `json:"confirmation_email_sent,omitempty"`
	ReceiptURL            *string `json:"receipt_url,omitempty"`
	InvoiceNumber         *string `json:"invoice_number,omitempty"`
}

// SupportInfo defines model for SupportInfo.
type SupportInfo struct {
	Email         *string `json:"email,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Hours         *string `json:"hours,omitempty"`
	HelpCenterURL *string `json:"help_center_url,omitempty"`
}

// Item defines model for Item.
type Item struct {
	ID         string  `json:"id"`
	Name       *string `json:"name,omitempty"`
	UnitAmount *int    `json:"unit_amount,omitempty"`
}

// VariantOption defines model for VariantOption.
type VariantOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// WeightInfo defines model for WeightInfo.
type WeightInfo struct {
	Value float64    `json:"value"`
	Unit  WeightUnit `json:"unit"`
}

// DimensionsInfo defines model for DimensionsInfo.
type DimensionsInfo struct {
	Length float64        `json:"length"`
	Width  float64        `json:"width"`
	Height float64        `json:"height"`
	Unit   DimensionsUnit `json:"unit"`
}

// DiscountDetail defines model for DiscountDetail.
type DiscountDetail struct {
	Code        *string               `json:"code,omitempty"`
	Type        DiscountDetailType    `json:"type"`
	Amount      int                   `json:"amount"`
	Description *string               `json:"description,omitempty"`
	Source      *DiscountDetailSource `json:"source,omitempty"`
}

// LineItem defines model for LineItem.
type LineItem struct {
	ID                       string                    `json:"id"`
	Item                     Item                      `json:"item"`
	Quantity                 int                       `json:"quantity"`
	Name                     *string                   `json:"name,omitempty"`
	Description              *string                   `json:"description,omitempty"`
	Images                   []string                  `json:"images,omitempty"`
	UnitAmount               *int                      `json:"unit_amount,omitempty"`
	Disclosures              []Disclosure              `json:"disclosures,omitempty"`
	CustomAttributes         []CustomAttribute         `json:"custom_attributes,omitempty"`
	MarketplaceSellerDetails *MarketplaceSellerDetails `json:"marketplace_seller_details,omitempty"`
	ProductID                *string                   `json:"product_id,omitempty"`
	SKU                      *string                   `json:"sku,omitempty"`
	VariantID                *string                   `json:"variant_id,omitempty"`
	Category                 *string                   `json:"category,omitempty"`
	Tags                     []string                  `json:"tags,omitempty"`
	Weight                   *WeightInfo               `json:"weight,omitempty"`
	Dimensions               *DimensionsInfo           `json:"dimensions,omitempty"`
	AvailabilityStatus       *AvailabilityStatus       `json:"availability_status,omitempty"`
	AvailableQuantity        *int                      `json:"available_quantity,omitempty"`
	MaxQuantityPerOrder      *int                      `json:"max_quantity_per_order,omitempty"`
	FulfillableOn            *time.Time                `json:"fulfillable_on,omitempty"`
	VariantOptions           []VariantOption           `json:"variant_options,omitempty"`
	DiscountDetails          []DiscountDetail          `json:"discount_details,omitempty"`
	TaxExempt                *bool                     `json:"tax_exempt,omitempty"`
	TaxExemptionReason       *string                   `json:"tax_exemption_reason,omitempty"`
	ParentID                 *string                   `json:"parent_id,omitempty"`
	// Totals contains the line-item totals breakdown including base_amount, discount, subtotal, tax, and total.
	Totals []Total `json:"totals"`
}

// CustomAttribute defines model for CustomAttribute.
type CustomAttribute struct {
	DisplayName string `json:"display_name"`
	Value       string `json:"value"`
}

// Disclosure defines model for Disclosure.
type Disclosure struct {
	Type        DisclosureType        `json:"type"`
	ContentType DisclosureContentType `json:"content_type"`
	Content     string                `json:"content"`
}

// Link defines model for Link.
type Link struct {
	Type  LinkType `json:"type"`
	Title *string  `json:"title,omitempty"`
	Url   string   `json:"url"`
}

// MarketplaceSellerDetails defines model for MarketplaceSellerDetails.
type MarketplaceSellerDetails struct {
	Name string `json:"name"`
}

// MessageInfo defines model for MessageInfo.
type MessageInfo struct {
	Content     string                 `json:"content"`
	ContentType MessageInfoContentType `json:"content_type"`
	Severity    *MessageSeverity       `json:"severity,omitempty"`
	Resolution  *MessageResolution     `json:"resolution,omitempty"`

	// Param RFC 9535 JSONPath
	Param *string `json:"param,omitempty"`
	Type  string  `json:"type"`
}

// MessageWarning defines model for MessageWarning.
type MessageWarning struct {
	Code        MessageWarningCode     `json:"code"`
	Content     string                 `json:"content"`
	ContentType MessageInfoContentType `json:"content_type"`
	Severity    *MessageSeverity       `json:"severity,omitempty"`
	Resolution  *MessageResolution     `json:"resolution,omitempty"`
	// Param RFC 9535 JSONPath
	Param *string `json:"param,omitempty"`
	Type  string  `json:"type"`
}

// MessageError defines model for MessageError.
type MessageError struct {
	Code        MessageErrorCode        `json:"code"`
	Content     string                  `json:"content"`
	ContentType MessageErrorContentType `json:"content_type"`
	Severity    *MessageSeverity        `json:"severity,omitempty"`
	Resolution  *MessageResolution      `json:"resolution,omitempty"`
	Param       *string                 `json:"param,omitempty"`
	Type        string                  `json:"type"`
}

// Order defines model for Order.
type Order struct {
	// Type is the discriminator for webhook payloads when present.
	Type *EventDataType `json:"type,omitempty"`
	// ID is the order identifier.
	ID string `json:"id"`
	// CheckoutSessionId identifies the checkout session used to create this order.
	CheckoutSessionId string `json:"checkout_session_id"`
	// OrderNumber is a human-readable order reference.
	OrderNumber *string `json:"order_number,omitempty"`
	// PermalinkUrl is the buyer-facing URL to view order details.
	PermalinkUrl string `json:"permalink_url"`
	// Status is the order-level lifecycle state.
	Status *CheckoutOrderStatus `json:"status,omitempty"`
	// EstimatedDelivery is the expected delivery window for this order.
	EstimatedDelivery *EstimatedDelivery `json:"estimated_delivery,omitempty"`
	// Confirmation contains confirmation metadata such as receipt URL.
	Confirmation *OrderConfirmation `json:"confirmation,omitempty"`
	// Support contains merchant support contact details for this order.
	Support *SupportInfo `json:"support,omitempty"`
	// LineItems tracks ordered items and fulfillment progress.
	LineItems []OrderLineItem `json:"line_items,omitempty"`
	// Fulfillments tracks shipping, pickup, and digital delivery.
	Fulfillments []Fulfillment `json:"fulfillments,omitempty"`
	// Adjustments tracks post-order changes like refunds/returns/disputes.
	Adjustments []Adjustment `json:"adjustments,omitempty"`
	// Totals captures order-level financial breakdown, including amount_refunded.
	Totals []Total `json:"totals,omitempty"`
}

// OrderLineItemStatus represents per-line fulfillment progress in an order.
type OrderLineItemStatus string

// Defines values for OrderLineItemStatus.
const (
	// OrderLineItemStatusProcessing means no units have shipped yet.
	OrderLineItemStatusProcessing OrderLineItemStatus = "processing"
	// OrderLineItemStatusPartial means some but not all units have shipped.
	OrderLineItemStatusPartial OrderLineItemStatus = "partial"
	// OrderLineItemStatusShipped means all ordered units have shipped.
	OrderLineItemStatusShipped OrderLineItemStatus = "shipped"
	// OrderLineItemStatusDelivered means the item has been delivered.
	OrderLineItemStatusDelivered OrderLineItemStatus = "delivered"
	// OrderLineItemStatusCanceled means fulfillment for this item was canceled.
	OrderLineItemStatusCanceled OrderLineItemStatus = "canceled"
)

// OrderLineItemQuantity captures ordered vs shipped quantity for a line item.
// Shipped is optional and defaults to zero when omitted.
type OrderLineItemQuantity struct {
	// Ordered is the quantity originally purchased.
	Ordered int `json:"ordered"`
	// Shipped is the quantity handed to carrier/fulfilled so far.
	Shipped *int `json:"shipped,omitempty"`
}

// OrderLineItem is a post-purchase line item with optional fulfillment and price details.
type OrderLineItem struct {
	// ID is the line item identifier for references in fulfillments/adjustments.
	ID string `json:"id"`
	// Title is the product name.
	Title string `json:"title"`
	// ProductID is the merchant catalog product identifier.
	ProductID *string `json:"product_id,omitempty"`
	// Description is the product description.
	Description *string `json:"description,omitempty"`
	// ImageURL is the product image URL.
	ImageURL *string `json:"image_url,omitempty"`
	// URL is the product page URL.
	URL *string `json:"url,omitempty"`
	// Quantity tracks ordered vs shipped units.
	Quantity OrderLineItemQuantity `json:"quantity"`
	// UnitPrice is the price per unit in minor currency units.
	UnitPrice *int `json:"unit_price,omitempty"`
	// Subtotal is quantity.ordered * unit_price in minor currency units.
	Subtotal *int `json:"subtotal,omitempty"`
	// Totals is the optional line-item totals breakdown.
	Totals []Total `json:"totals,omitempty"`
	// Status is the current fulfillment state for this line item.
	Status *OrderLineItemStatus `json:"status,omitempty"`
}

// LineItemReference points to a line item and quantity in fulfillments/adjustments.
type LineItemReference struct {
	// ID references an OrderLineItem.ID.
	ID string `json:"id"`
	// Quantity is the affected quantity for this reference.
	Quantity int `json:"quantity"`
}

// FulfillmentType identifies the delivery method for a fulfillment.
type FulfillmentType string

// Defines values for FulfillmentType.
const (
	// FulfillmentTypeShipping is carrier shipment.
	FulfillmentTypeShipping FulfillmentType = "shipping"
	// FulfillmentTypePickup is in-store or curbside pickup.
	FulfillmentTypePickup FulfillmentType = "pickup"
	// FulfillmentTypeDigital is non-physical digital delivery.
	FulfillmentTypeDigital FulfillmentType = "digital"
)

// FulfillmentStatus is the current state of a fulfillment.
// Not every state applies to every fulfillment type.
type FulfillmentStatus string

// Defines values for FulfillmentStatus.
const (
	// FulfillmentStatusPending means fulfillment is queued but not started.
	FulfillmentStatusPending FulfillmentStatus = "pending"
	// FulfillmentStatusProcessing means fulfillment is in progress.
	FulfillmentStatusProcessing FulfillmentStatus = "processing"
	// FulfillmentStatusShipped means shipment was handed to a carrier.
	FulfillmentStatusShipped FulfillmentStatus = "shipped"
	// FulfillmentStatusInTransit means shipment is moving through carrier network.
	FulfillmentStatusInTransit FulfillmentStatus = "in_transit"
	// FulfillmentStatusOutForDelivery means final-mile delivery is in progress.
	FulfillmentStatusOutForDelivery FulfillmentStatus = "out_for_delivery"
	// FulfillmentStatusReadyForPickup means buyer can collect the order.
	FulfillmentStatusReadyForPickup FulfillmentStatus = "ready_for_pickup"
	// FulfillmentStatusDelivered means fulfillment completed successfully.
	FulfillmentStatusDelivered FulfillmentStatus = "delivered"
	// FulfillmentStatusFailed means fulfillment could not be completed.
	FulfillmentStatusFailed FulfillmentStatus = "failed"
	// FulfillmentStatusCanceled means fulfillment was canceled.
	FulfillmentStatusCanceled FulfillmentStatus = "canceled"
)

// FulfillmentEventType is a point-in-time event in a fulfillment timeline.
type FulfillmentEventType string

// Defines values for FulfillmentEventType.
const (
	// FulfillmentEventTypeProcessing indicates processing activity.
	FulfillmentEventTypeProcessing FulfillmentEventType = "processing"
	// FulfillmentEventTypeShipped indicates carrier handoff.
	FulfillmentEventTypeShipped FulfillmentEventType = "shipped"
	// FulfillmentEventTypeInTransit indicates carrier network transit.
	FulfillmentEventTypeInTransit FulfillmentEventType = "in_transit"
	// FulfillmentEventTypeOutForDelivery indicates final-mile dispatch.
	FulfillmentEventTypeOutForDelivery FulfillmentEventType = "out_for_delivery"
	// FulfillmentEventTypeReadyForPickup indicates pickup availability.
	FulfillmentEventTypeReadyForPickup FulfillmentEventType = "ready_for_pickup"
	// FulfillmentEventTypeDelivered indicates successful delivery.
	FulfillmentEventTypeDelivered FulfillmentEventType = "delivered"
	// FulfillmentEventTypeFailedAttempt indicates an unsuccessful delivery attempt.
	FulfillmentEventTypeFailedAttempt FulfillmentEventType = "failed_attempt"
	// FulfillmentEventTypeReturned indicates returned/return-to-sender flow.
	FulfillmentEventTypeReturned FulfillmentEventType = "returned"
)

// FulfillmentDigitalDelivery holds digital access details for digital fulfillments.
type FulfillmentDigitalDelivery struct {
	// AccessURL is the URL where the buyer can access the digital content.
	AccessURL *string `json:"access_url,omitempty"`
	// LicenseKey is an activation/license key when applicable.
	LicenseKey *string `json:"license_key,omitempty"`
	// ExpiresAt is the RFC3339 timestamp after which access expires.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// FulfillmentEvent is an immutable event entry describing fulfillment progress.
type FulfillmentEvent struct {
	// ID is the event identifier.
	ID string `json:"id"`
	// Type is the event type.
	Type FulfillmentEventType `json:"type"`
	// OccurredAt is the RFC3339 timestamp when the event occurred.
	OccurredAt time.Time `json:"occurred_at"`
	// Description is optional human-readable event detail.
	Description *string `json:"description,omitempty"`
	// Location is the optional location of the event.
	Location *string `json:"location,omitempty"`
}

// Fulfillment describes how some order items are delivered (shipping, pickup, digital).
type Fulfillment struct {
	// ID is the fulfillment identifier.
	ID string `json:"id"`
	// Type selects shipping, pickup, or digital delivery.
	Type FulfillmentType `json:"type"`
	// Status is the current fulfillment state.
	Status *FulfillmentStatus `json:"status,omitempty"`
	// LineItems lists included line items and quantities.
	LineItems []LineItemReference `json:"line_items,omitempty"`
	// Carrier is the shipping carrier name for shipping fulfillments.
	Carrier *string `json:"carrier,omitempty"`
	// TrackingNumber is the carrier tracking number for shipping fulfillments.
	TrackingNumber *string `json:"tracking_number,omitempty"`
	// TrackingURL is the external tracking URL for shipping fulfillments.
	TrackingURL *string `json:"tracking_url,omitempty"`
	// Destination is the destination address when applicable.
	Destination *Address `json:"destination,omitempty"`
	// EstimatedDelivery is the expected delivery window.
	EstimatedDelivery *EstimatedDelivery `json:"estimated_delivery,omitempty"`
	// DigitalDelivery contains digital access credentials/links.
	DigitalDelivery *FulfillmentDigitalDelivery `json:"digital_delivery,omitempty"`
	// Description is optional human-readable fulfillment context.
	Description *string `json:"description,omitempty"`
	// Events is an append-only progress log.
	Events []FulfillmentEvent `json:"events,omitempty"`
}

// AdjustmentType classifies a post-order change such as refunds or disputes.
type AdjustmentType string

// Defines values for AdjustmentType.
const (
	// AdjustmentTypeRefund is a full refund.
	AdjustmentTypeRefund AdjustmentType = "refund"
	// AdjustmentTypePartialRefund is a partial refund.
	AdjustmentTypePartialRefund AdjustmentType = "partial_refund"
	// AdjustmentTypeStoreCredit grants store credit.
	AdjustmentTypeStoreCredit AdjustmentType = "store_credit"
	// AdjustmentTypeReturn records returned goods.
	AdjustmentTypeReturn AdjustmentType = "return"
	// AdjustmentTypeExchange records an item exchange.
	AdjustmentTypeExchange AdjustmentType = "exchange"
	// AdjustmentTypeCancellation records a cancellation adjustment.
	AdjustmentTypeCancellation AdjustmentType = "cancellation"
	// AdjustmentTypeDispute records a dispute.
	AdjustmentTypeDispute AdjustmentType = "dispute"
	// AdjustmentTypeChargeback records a chargeback.
	AdjustmentTypeChargeback AdjustmentType = "chargeback"
)

// AdjustmentStatus is the processing state of an order adjustment.
type AdjustmentStatus string

// Defines values for AdjustmentStatus.
const (
	// AdjustmentStatusPending means processing has not finished yet.
	AdjustmentStatusPending AdjustmentStatus = "pending"
	// AdjustmentStatusCompleted means the adjustment was applied.
	AdjustmentStatusCompleted AdjustmentStatus = "completed"
	// AdjustmentStatusFailed means the adjustment failed.
	AdjustmentStatusFailed AdjustmentStatus = "failed"
)

// Adjustment is a post-order financial/logistical change (refund, return, dispute, etc.).
type Adjustment struct {
	// ID is the adjustment identifier.
	ID string `json:"id"`
	// Type classifies the adjustment.
	Type AdjustmentType `json:"type"`
	// OccurredAt is the RFC3339 timestamp for the adjustment event.
	OccurredAt time.Time `json:"occurred_at"`
	// Status is the adjustment processing status.
	Status AdjustmentStatus `json:"status"`
	// LineItems lists affected line items and quantities.
	LineItems []LineItemReference `json:"line_items,omitempty"`
	// Amount is the credited amount in minor units, inclusive of tax.
	Amount *int `json:"amount,omitempty"`
	// Currency is the ISO 4217 currency code for Amount.
	Currency *string `json:"currency,omitempty"`
	// Description is an optional human-readable explanation.
	Description *string `json:"description,omitempty"`
	// Reason is an optional structured reason code.
	Reason *string `json:"reason,omitempty"`
}

// PaymentData defines model for PaymentData.
type PaymentData struct {
	// HandlerID identifies which advertised payment handler to use.
	HandlerID *string `json:"handler_id,omitempty"`
	// Instrument carries handler-specific payment instrument details.
	Instrument *PaymentInstrument `json:"instrument,omitempty"`
	// Deprecated: prefer handler_id + instrument for ACP payment handlers.
	// Token is a provider-issued payment token (required with Provider).
	Token *string `json:"token,omitempty"`
	// Deprecated: prefer handler_id + instrument for ACP payment handlers.
	// Provider identifies the payment provider that issued Token.
	Provider       *PaymentDataProvider `json:"provider,omitempty"`
	BillingAddress *Address             `json:"billing_address,omitempty"`
	// PurchaseOrderNumber enables non-tokenized payment flows (B2B).
	PurchaseOrderNumber *string `json:"purchase_order_number,omitempty"`
	// PaymentTerms specify invoice payment terms (for example "net_30").
	PaymentTerms *PaymentTerms `json:"payment_terms,omitempty"`
	// DueDate is the payment due date for invoice flows.
	DueDate *time.Time `json:"due_date,omitempty"`
	// ApprovalRequired signals whether buyer approval is required before payment capture.
	ApprovalRequired *bool `json:"approval_required,omitempty"`
}

// PaymentInstrument describes the selected payment instrument.
type PaymentInstrument struct {
	// Type is the instrument type (for example "card" or "wallet_token").
	Type string `json:"type"`
	// Credential carries the credential required by the selected handler.
	Credential PaymentCredential `json:"credential"`
}

// PaymentCredential holds a credential token used by a payment instrument.
type PaymentCredential struct {
	// Type is the credential type (for example "spt" or "wallet_token").
	Type string `json:"type"`
	// Token is the opaque credential/token value.
	Token string `json:"token"`
}

// PaymentDataProvider defines model for PaymentData.Provider.
type PaymentDataProvider string

// ProtocolVersion defines model for ProtocolVersion.
type ProtocolVersion struct {
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities,omitempty"`
}

// PaymentResponse defines model for PaymentResponse.
type PaymentResponse struct {
	Provider    *string          `json:"provider,omitempty"`
	Instruments []map[string]any `json:"instruments,omitempty"`
	Handlers    []map[string]any `json:"handlers,omitempty"`
}

// RiskSignals defines model for RiskSignals.
type RiskSignals struct {
	IPAddress         *string `json:"ip_address,omitempty"`
	UserAgent         *string `json:"user_agent,omitempty"`
	AcceptLanguage    *string `json:"accept_language,omitempty"`
	SessionID         *string `json:"session_id,omitempty"`
	DeviceFingerprint *string `json:"device_fingerprint,omitempty"`
}

// AffiliateAttributionSource defines model for AffiliateAttribution.Source.
type AffiliateAttributionSource struct {
	Type string  `json:"type"`
	URL  *string `json:"url,omitempty"`
}

// AffiliateAttributionMetadata defines model for AffiliateAttribution.Metadata.
type AffiliateAttributionMetadata map[string]any

// AffiliateAttribution defines model for AffiliateAttribution.
type AffiliateAttribution struct {
	// Provider is the attribution provider / affiliate network namespace (for example "impact.com").
	Provider string `json:"provider"`
	// Token is an opaque provider-issued token for fraud-resistant validation.
	Token *string `json:"token,omitempty"`
	// PublisherID is the provider-scoped affiliate/publisher identifier.
	PublisherID *string                      `json:"publisher_id,omitempty"`
	CampaignID  *string                      `json:"campaign_id,omitempty"`
	CreativeID  *string                      `json:"creative_id,omitempty"`
	SubID       *string                      `json:"sub_id,omitempty"`
	Source      *AffiliateAttributionSource  `json:"source,omitempty"`
	IssuedAt    *time.Time                   `json:"issued_at,omitempty"`
	ExpiresAt   *time.Time                   `json:"expires_at,omitempty"`
	Metadata    AffiliateAttributionMetadata `json:"metadata,omitempty"`
	// Touchpoint is the attribution touchpoint type ("first" or "last").
	Touchpoint *string `json:"touchpoint,omitempty"`
}

// AuthenticationDirectoryServer defines supported 3DS directory servers.
type AuthenticationDirectoryServer string

// Defines values for AuthenticationDirectoryServer.
const (
	AuthenticationDirectoryServerAmericanExpress AuthenticationDirectoryServer = "american_express"
	AuthenticationDirectoryServerMastercard      AuthenticationDirectoryServer = "mastercard"
	AuthenticationDirectoryServerVisa            AuthenticationDirectoryServer = "visa"
)

// AuthenticationFlowPreferenceType defines flow preference type.
type AuthenticationFlowPreferenceType string

// Defines values for AuthenticationFlowPreferenceType.
const (
	AuthenticationFlowPreferenceTypeChallenge    AuthenticationFlowPreferenceType = "challenge"
	AuthenticationFlowPreferenceTypeFrictionless AuthenticationFlowPreferenceType = "frictionless"
)

// AuthenticationChallengePreferenceType defines challenge preference type.
type AuthenticationChallengePreferenceType string

// Defines values for AuthenticationChallengePreferenceType.
const (
	AuthenticationChallengePreferenceTypeMandated  AuthenticationChallengePreferenceType = "mandated"
	AuthenticationChallengePreferenceTypePreferred AuthenticationChallengePreferenceType = "preferred"
)

// AuthenticationOutcome defines 3DS authentication outcomes.
type AuthenticationOutcome string

// Defines values for AuthenticationOutcome.
const (
	AuthenticationOutcomeAbandoned           AuthenticationOutcome = "abandoned"
	AuthenticationOutcomeAttemptAcknowledged AuthenticationOutcome = "attempt_acknowledged"
	AuthenticationOutcomeAuthenticated       AuthenticationOutcome = "authenticated"
	AuthenticationOutcomeCanceled            AuthenticationOutcome = "canceled"
	AuthenticationOutcomeDenied              AuthenticationOutcome = "denied"
	AuthenticationOutcomeInformational       AuthenticationOutcome = "informational"
	AuthenticationOutcomeInternalError       AuthenticationOutcome = "internal_error"
	AuthenticationOutcomeNotSupported        AuthenticationOutcome = "not_supported"
	AuthenticationOutcomeProcessingError     AuthenticationOutcome = "processing_error"
	AuthenticationOutcomeRejected            AuthenticationOutcome = "rejected"
)

// AuthenticationECI defines accepted ECI values.
type AuthenticationECI string

// Defines values for AuthenticationECI.
const (
	AuthenticationECI01 AuthenticationECI = "01"
	AuthenticationECI02 AuthenticationECI = "02"
	AuthenticationECI05 AuthenticationECI = "05"
	AuthenticationECI06 AuthenticationECI = "06"
	AuthenticationECI07 AuthenticationECI = "07"
)

// AuthenticationMetadata captures seller-provided authentication metadata for 3DS flows.
type AuthenticationMetadata struct {
	// AcquirerDetails are details about the acquirer used for this 3DS Authentication.
	AcquirerDetails AuthenticationAcquirerDetails `json:"acquirer_details"`
	// DirectoryServer is the 3DS directory server used for this Authentication.
	DirectoryServer AuthenticationDirectoryServer `json:"directory_server"`
	// FlowPreference captures seller's preferred 3DS authentication flow, if any.
	FlowPreference *AuthenticationFlowPreference `json:"flow_preference,omitempty"`
}

// AuthenticationAcquirerDetails describes the acquirer details used for this 3DS Authentication.
type AuthenticationAcquirerDetails struct {
	// AcquirerBin is the acquirer BIN (directory-server specific).
	AcquirerBin string `json:"acquirer_bin"`
	// AcquirerCountry is the ISO 3166-1 alpha-2 country code.
	AcquirerCountry string `json:"acquirer_country"`
	// AcquirerMerchantId is the merchant ID assigned by the acquirer.
	AcquirerMerchantId string `json:"acquirer_merchant_id"`
	// MerchantName is the merchant name assigned by the acquirer.
	MerchantName string `json:"merchant_name"`
	// RequestorId is the requestor ID when required by the directory server.
	RequestorId *string `json:"requestor_id,omitempty"`
}

// AuthenticationFlowPreference captures seller preferences for the 3DS authentication flow.
type AuthenticationFlowPreference struct {
	// Type is the preferred flow type: "challenge" or "frictionless".
	Type AuthenticationFlowPreferenceType `json:"type"`
	// Challenge captures details about a requested challenge flow.
	Challenge *AuthenticationChallengePreference `json:"challenge,omitempty"`
	// Frictionless captures details about a requested frictionless flow.
	Frictionless *AuthenticationFrictionlessPreference `json:"frictionless,omitempty"`
}

// AuthenticationChallengePreference captures details about a requested challenge flow.
type AuthenticationChallengePreference struct {
	// Type is the challenge subtype preference.
	Type AuthenticationChallengePreferenceType `json:"type"`
}

// AuthenticationFrictionlessPreference captures details about a requested frictionless flow.
type AuthenticationFrictionlessPreference struct{}

// AuthenticationResult represents agent-provided 3DS authentication results.
type AuthenticationResult struct {
	// Outcome is the outcome of this 3DS Authentication.
	Outcome AuthenticationOutcome `json:"outcome"`
	// OutcomeDetails contains detailed authentication data.
	// Required when Outcome is authenticated, informational, or attempt_acknowledged.
	OutcomeDetails *AuthenticationOutcomeDetails `json:"outcome_details,omitempty"`
}

// AuthenticationOutcomeDetails provides detailed 3DS authentication data.
type AuthenticationOutcomeDetails struct {
	// ThreeDsCryptogram is the 3DS cryptogram (authentication value / AAV/CAVV/AEVV).
	// This value is 20 bytes, base64-encoded into a 28-character string.
	ThreeDsCryptogram string `json:"three_ds_cryptogram"`
	// ElectronicCommerceIndicator is the ECI returned by the 3DS provider.
	ElectronicCommerceIndicator AuthenticationECI `json:"electronic_commerce_indicator"`
	// TransactionId is the 3DS transaction identifier (XID for 3DS1, dsTransID for 3DS2).
	TransactionId string `json:"transaction_id"`
	// Version is the 3D Secure version used for this authentication (for example "1.0.2" or "2.2.0").
	Version string `json:"version"`
}

// PaymentMethod defines model for PaymentProvider.SupportedPaymentMethods items.
type PaymentMethod struct {
	SupportedCardNetworks []PaymentCardNetwork `json:"supported_card_networks,omitempty"`
	Type                  PaymentMethodType    `json:"type"`
}

// PaymentProvider defines model for PaymentProvider.
type PaymentProvider struct {
	Provider                PaymentProviderProvider `json:"provider"`
	SupportedPaymentMethods []PaymentMethod         `json:"supported_payment_methods"`
}

// PaymentProviderProvider defines model for PaymentProvider.Provider.
type PaymentProviderProvider string

// TaxBreakdownItem defines model for TaxBreakdownItem.
type TaxBreakdownItem struct {
	Jurisdiction string  `json:"jurisdiction"`
	Rate         float64 `json:"rate"`
	Amount       int     `json:"amount"`
}

// Total defines model for Total.
type Total struct {
	Amount            int                `json:"amount"`
	PresentmentAmount *int               `json:"presentment_amount,omitempty"`
	Description       *string            `json:"description,omitempty"`
	DisplayText       string             `json:"display_text"`
	Type              TotalType          `json:"type"`
	Breakdown         []TaxBreakdownItem `json:"breakdown,omitempty"`
}

// IntentTrace defines model for IntentTrace.
type IntentTrace struct {
	ReasonCode   IntentTraceReasonCode `json:"reason_code"`
	TraceSummary *string               `json:"trace_summary,omitempty"`
	Metadata     map[string]any        `json:"metadata,omitempty"`
}

// CancelSessionRequest defines model for CancelSessionRequest.
type CancelSessionRequest struct {
	IntentTrace *IntentTrace `json:"intent_trace,omitempty"`
}

// AsFulfillmentOptionShipping returns the union data inside the CheckoutSessionBase_FulfillmentOptions_Item as a FulfillmentOptionShipping
func (t FulfillmentOption) AsFulfillmentOptionShipping() (FulfillmentOptionShipping, error) {
	var body FulfillmentOptionShipping
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromFulfillmentOptionShipping overwrites any union data inside the CheckoutSessionBase_FulfillmentOptions_Item as the provided FulfillmentOptionShipping
func (t *FulfillmentOption) FromFulfillmentOptionShipping(v FulfillmentOptionShipping) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeFulfillmentOptionShipping performs a merge with any union data inside the CheckoutSessionBase_FulfillmentOptions_Item, using the provided FulfillmentOptionShipping
func (t *FulfillmentOption) MergeFulfillmentOptionShipping(v FulfillmentOptionShipping) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsFulfillmentOptionDigital returns the union data inside the CheckoutSessionBase_FulfillmentOptions_Item as a FulfillmentOptionDigital
func (t FulfillmentOption) AsFulfillmentOptionDigital() (FulfillmentOptionDigital, error) {
	var body FulfillmentOptionDigital
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromFulfillmentOptionDigital overwrites any union data inside the CheckoutSessionBase_FulfillmentOptions_Item as the provided FulfillmentOptionDigital
func (t *FulfillmentOption) FromFulfillmentOptionDigital(v FulfillmentOptionDigital) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeFulfillmentOptionDigital performs a merge with any union data inside the CheckoutSessionBase_FulfillmentOptions_Item, using the provided FulfillmentOptionDigital
func (t *FulfillmentOption) MergeFulfillmentOptionDigital(v FulfillmentOptionDigital) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsFulfillmentOptionPickup returns the union data inside the CheckoutSessionBase_FulfillmentOptions_Item as a FulfillmentOptionPickup
func (t FulfillmentOption) AsFulfillmentOptionPickup() (FulfillmentOptionPickup, error) {
	var body FulfillmentOptionPickup
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromFulfillmentOptionPickup overwrites any union data inside the CheckoutSessionBase_FulfillmentOptions_Item as the provided FulfillmentOptionPickup
func (t *FulfillmentOption) FromFulfillmentOptionPickup(v FulfillmentOptionPickup) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeFulfillmentOptionPickup performs a merge with any union data inside the CheckoutSessionBase_FulfillmentOptions_Item, using the provided FulfillmentOptionPickup
func (t *FulfillmentOption) MergeFulfillmentOptionPickup(v FulfillmentOptionPickup) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsFulfillmentOptionLocalDelivery returns the union data inside the CheckoutSessionBase_FulfillmentOptions_Item as a FulfillmentOptionLocalDelivery
func (t FulfillmentOption) AsFulfillmentOptionLocalDelivery() (FulfillmentOptionLocalDelivery, error) {
	var body FulfillmentOptionLocalDelivery
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromFulfillmentOptionLocalDelivery overwrites any union data inside the CheckoutSessionBase_FulfillmentOptions_Item as the provided FulfillmentOptionLocalDelivery
func (t *FulfillmentOption) FromFulfillmentOptionLocalDelivery(v FulfillmentOptionLocalDelivery) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeFulfillmentOptionLocalDelivery performs a merge with any union data inside the CheckoutSessionBase_FulfillmentOptions_Item, using the provided FulfillmentOptionLocalDelivery
func (t *FulfillmentOption) MergeFulfillmentOptionLocalDelivery(v FulfillmentOptionLocalDelivery) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// MarshalJSON serializes the underlying union for CheckoutSessionBase_FulfillmentOptions_Item.
func (t FulfillmentOption) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

// UnmarshalJSON loads union data for CheckoutSessionBase_FulfillmentOptions_Item.
func (t *FulfillmentOption) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsMessageInfo returns the union data inside the CheckoutSessionBase_Messages_Item as a MessageInfo
func (t Message) AsMessageInfo() (MessageInfo, error) {
	var body MessageInfo
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromMessageInfo overwrites any union data inside the CheckoutSessionBase_Messages_Item as the provided MessageInfo
func (t *Message) FromMessageInfo(v MessageInfo) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeMessageInfo performs a merge with any union data inside the CheckoutSessionBase_Messages_Item, using the provided MessageInfo
func (t *Message) MergeMessageInfo(v MessageInfo) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsMessageWarning returns the union data inside the CheckoutSessionBase_Messages_Item as a MessageWarning
func (t Message) AsMessageWarning() (MessageWarning, error) {
	var body MessageWarning
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromMessageWarning overwrites any union data inside the CheckoutSessionBase_Messages_Item as the provided MessageWarning
func (t *Message) FromMessageWarning(v MessageWarning) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeMessageWarning performs a merge with any union data inside the CheckoutSessionBase_Messages_Item, using the provided MessageWarning
func (t *Message) MergeMessageWarning(v MessageWarning) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsMessageError returns the union data inside the CheckoutSessionBase_Messages_Item as a MessageError
func (t Message) AsMessageError() (MessageError, error) {
	var body MessageError
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromMessageError overwrites any union data inside the CheckoutSessionBase_Messages_Item as the provided MessageError
func (t *Message) FromMessageError(v MessageError) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeMessageError performs a merge with any union data inside the CheckoutSessionBase_Messages_Item, using the provided MessageError
func (t *Message) MergeMessageError(v MessageError) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// MarshalJSON serializes the underlying union for CheckoutSessionBase_Messages_Item.
func (t Message) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

// UnmarshalJSON loads union data for CheckoutSessionBase_Messages_Item.
func (t *Message) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
