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
	CheckoutSessionStatusCanceled               CheckoutSessionStatus = "canceled"
	CheckoutSessionStatusCompleted              CheckoutSessionStatus = "completed"
	CheckoutSessionStatusInProgress             CheckoutSessionStatus = "in_progress"
	CheckoutSessionStatusNotReadyForPayment     CheckoutSessionStatus = "not_ready_for_payment"
	CheckoutSessionStatusAuthenticationRequired CheckoutSessionStatus = "authentication_required"
	CheckoutSessionStatusReadyForPayment        CheckoutSessionStatus = "ready_for_payment"
)

// LinkType defines model for Link.Type.
type LinkType string

// Defines values for LinkType.
const (
	PrivacyPolicy LinkType = "privacy_policy"
	ReturnPolicy  LinkType = "return_policy"
	// Deprecated: removed in ACP 2025-12-11; attach policies to line items in marketplace scenarios.
	SellerShopPolicies LinkType = "seller_shop_policies"
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
	Invalid         MessageErrorCode = "invalid"
	Missing         MessageErrorCode = "missing"
	OutOfStock      MessageErrorCode = "out_of_stock"
	PaymentDeclined MessageErrorCode = "payment_declined"
	Requires3ds     MessageErrorCode = "requires_3ds"
	RequiresSignIn  MessageErrorCode = "requires_sign_in"
)

// MessageErrorContentType defines model for MessageError.ContentType.
type MessageErrorContentType string

// Defines values for MessageErrorContentType.
const (
	MessageErrorContentTypeMarkdown MessageErrorContentType = "markdown"
	MessageErrorContentTypePlain    MessageErrorContentType = "plain"
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

// TotalType defines model for Total.Type.
type TotalType string

// Defines values for TotalType.
const (
	TotalTypeDiscount        TotalType = "discount"
	TotalTypeFee             TotalType = "fee"
	TotalTypeFulfillment     TotalType = "fulfillment"
	TotalTypeItemsBaseAmount TotalType = "items_base_amount"
	TotalTypeItemsDiscount   TotalType = "items_discount"
	TotalTypeSubtotal        TotalType = "subtotal"
	TotalTypeTax             TotalType = "tax"
	TotalTypeTotal           TotalType = "total"
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

// Buyer defines model for Buyer.
type Buyer struct {
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// CheckoutSession defines model for CheckoutSession.
type CheckoutSession struct {
	ID                  string              `json:"id"`
	Buyer               *Buyer              `json:"buyer,omitempty"`
	Currency            string              `json:"currency"`
	FulfillmentAddress  *Address            `json:"fulfillment_address,omitempty"`
	FulfillmentOptionId *string             `json:"fulfillment_option_id,omitempty"`
	FulfillmentOptions  []FulfillmentOption `json:"fulfillment_options"`
	LineItems           []LineItem          `json:"line_items"`
	Links               []Link              `json:"links"`
	Messages            []Message           `json:"messages"`
	// AuthenticationMetadata is seller-provided authentication metadata for 3DS flows.
	AuthenticationMetadata *AuthenticationMetadata `json:"authentication_metadata,omitempty"`
	PaymentProvider        *PaymentProvider        `json:"payment_provider,omitempty"`
	Status                 CheckoutSessionStatus   `json:"status"`
	Totals                 []Total                 `json:"totals"`
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
}

// CheckoutSessionCreateRequest defines model for CheckoutSessionCreateRequest.
type CheckoutSessionCreateRequest struct {
	Buyer              *Buyer   `json:"buyer,omitempty"`
	FulfillmentAddress *Address `json:"fulfillment_address,omitempty"`
	Items              []Item   `json:"items"`
}

// CheckoutSessionUpdateRequest defines model for CheckoutSessionUpdateRequest.
type CheckoutSessionUpdateRequest struct {
	Buyer               *Buyer   `json:"buyer,omitempty"`
	FulfillmentAddress  *Address `json:"fulfillment_address,omitempty"`
	FulfillmentOptionId *string  `json:"fulfillment_option_id,omitempty"`
	Items               *[]Item  `json:"items,omitempty"`
}

// SessionWithOrder defines model for SessionWithOrder.
type SessionWithOrder struct {
	CheckoutSession
	Order Order `json:"order"`
}

// FulfillmentOptionDigital defines model for FulfillmentOptionDigital.
type FulfillmentOptionDigital struct {
	ID       string  `json:"id"`
	Subtitle *string `json:"subtitle,omitempty"`
	Subtotal int     `json:"subtotal"`
	Tax      int     `json:"tax"`
	Title    string  `json:"title"`
	Total    int     `json:"total"`
	Type     string  `json:"type"`
}

// FulfillmentOptionShipping defines model for FulfillmentOptionShipping.
type FulfillmentOptionShipping struct {
	ID                   string     `json:"id"`
	Carrier              *string    `json:"carrier,omitempty"`
	EarliestDeliveryTime *time.Time `json:"earliest_delivery_time,omitempty"`
	LatestDeliveryTime   *time.Time `json:"latest_delivery_time,omitempty"`
	Subtitle             *string    `json:"subtitle,omitempty"`
	Subtotal             int        `json:"subtotal"`
	Tax                  int        `json:"tax"`
	Title                string     `json:"title"`
	Total                int        `json:"total"`
	Type                 string     `json:"type"`
}

// Item defines model for Item.
type Item struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

// LineItem defines model for LineItem.
type LineItem struct {
	ID                       string                    `json:"id"`
	BaseAmount               int                       `json:"base_amount"`
	Discount                 int                       `json:"discount"`
	Item                     Item                      `json:"item"`
	Subtotal                 int                       `json:"subtotal"`
	Tax                      int                       `json:"tax"`
	Total                    int                       `json:"total"`
	Name                     *string                   `json:"name,omitempty"`
	Description              *string                   `json:"description,omitempty"`
	Images                   []string                  `json:"images,omitempty"`
	UnitAmount               *int                      `json:"unit_amount,omitempty"`
	Disclosures              []Disclosure              `json:"disclosures,omitempty"`
	CustomAttributes         []CustomAttribute         `json:"custom_attributes,omitempty"`
	MarketplaceSellerDetails *MarketplaceSellerDetails `json:"marketplace_seller_details,omitempty"`
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
	Type LinkType `json:"type"`
	Url  string   `json:"url"`
}

// MarketplaceSellerDetails defines model for MarketplaceSellerDetails.
type MarketplaceSellerDetails struct {
	Name string `json:"name"`
}

// MessageInfo defines model for MessageInfo.
type MessageInfo struct {
	Content     string                 `json:"content"`
	ContentType MessageInfoContentType `json:"content_type"`

	// Param RFC 9535 JSONPath
	Param *string `json:"param,omitempty"`
	Type  string  `json:"type"`
}

// MessageError defines model for MessageError.
type MessageError struct {
	Code        MessageErrorCode        `json:"code"`
	Content     string                  `json:"content"`
	ContentType MessageErrorContentType `json:"content_type"`
	Param       *string                 `json:"param,omitempty"`
	Type        string                  `json:"type"`
}

// Order defines model for Order.
type Order struct {
	ID                string `json:"id"`
	CheckoutSessionId string `json:"checkout_session_id"`
	PermalinkUrl      string `json:"permalink_url"`
}

// PaymentData defines model for PaymentData.
type PaymentData struct {
	BillingAddress *Address            `json:"billing_address,omitempty"`
	Provider       PaymentDataProvider `json:"provider"`
	Token          string              `json:"token"`
}

// PaymentDataProvider defines model for PaymentData.Provider.
type PaymentDataProvider string

// AuthenticationChannelType defines the channel used for authentication.
type AuthenticationChannelType string

// Defines values for AuthenticationChannelType.
const (
	AuthenticationChannelTypeBrowser AuthenticationChannelType = "browser"
)

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
	// Channel captures details about the channel used for this 3DS Authentication.
	Channel AuthenticationChannel `json:"channel"`
	// AcquirerDetails are details about the acquirer used for this 3DS Authentication.
	AcquirerDetails AuthenticationAcquirerDetails `json:"acquirer_details"`
	// DirectoryServer is the 3DS directory server used for this Authentication.
	DirectoryServer AuthenticationDirectoryServer `json:"directory_server"`
	// FlowPreference captures seller's preferred 3DS authentication flow, if any.
	FlowPreference *AuthenticationFlowPreference `json:"flow_preference,omitempty"`
}

// AuthenticationChannel describes the channel used for a 3DS Authentication.
type AuthenticationChannel struct {
	// Type is the channel type. Use "browser" to indicate a browser-originated transaction.
	Type AuthenticationChannelType `json:"type"`
	// Browser contains browser details collected server- or client-side.
	Browser AuthenticationBrowser `json:"browser"`
}

// AuthenticationBrowser contains browser details collected server- or client-side for 3DS.
type AuthenticationBrowser struct {
	// AcceptHeader is the HTTP Accept header from the cardholder's browser.
	AcceptHeader string `json:"accept_header"`
	// IPAddress is the IP address of the browser.
	IPAddress string `json:"ip_address"`
	// JavascriptEnabled indicates whether the browser can execute JavaScript.
	JavascriptEnabled bool `json:"javascript_enabled"`
	// Language is an IETF BCP 47 language tag representing the browser language.
	Language string `json:"language"`
	// UserAgent is the browser user agent string.
	UserAgent string `json:"user_agent"`
	// ColorDepth is the screen color depth. Required if JavascriptEnabled is true.
	ColorDepth *int `json:"color_depth,omitempty"`
	// JavaEnabled indicates browser Java support. Required if JavascriptEnabled is true.
	JavaEnabled *bool `json:"java_enabled,omitempty"`
	// ScreenHeight is the screen height in pixels. Required if JavascriptEnabled is true.
	ScreenHeight *int `json:"screen_height,omitempty"`
	// ScreenWidth is the screen width in pixels. Required if JavascriptEnabled is true.
	ScreenWidth *int `json:"screen_width,omitempty"`
	// TimezoneOffset is the time difference in minutes between UTC and local time.
	// Required if JavascriptEnabled is true.
	TimezoneOffset *int `json:"timezone_offset,omitempty"`
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
	SupportedCardNetworks []PaymentCardNetwork `json:"supported_card_networks"`
	Type                  PaymentMethodType    `json:"type"`
}

// PaymentProvider defines model for PaymentProvider.
type PaymentProvider struct {
	Provider                PaymentProviderProvider `json:"provider"`
	SupportedPaymentMethods []PaymentMethod         `json:"supported_payment_methods"`
}

// PaymentProviderProvider defines model for PaymentProvider.Provider.
type PaymentProviderProvider string

// Total defines model for Total.
type Total struct {
	Amount      int       `json:"amount"`
	Description *string   `json:"description,omitempty"`
	DisplayText string    `json:"display_text"`
	Type        TotalType `json:"type"`
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
