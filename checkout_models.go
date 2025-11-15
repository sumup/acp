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
	CheckoutSessionStatusCanceled           CheckoutSessionStatus = "canceled"
	CheckoutSessionStatusCompleted          CheckoutSessionStatus = "completed"
	CheckoutSessionStatusInProgress         CheckoutSessionStatus = "in_progress"
	CheckoutSessionStatusNotReadyForPayment CheckoutSessionStatus = "not_ready_for_payment"
	CheckoutSessionStatusReadyForPayment    CheckoutSessionStatus = "ready_for_payment"
)

// LinkType defines model for Link.Type.
type LinkType string

// Defines values for LinkType.
const (
	PrivacyPolicy      LinkType = "privacy_policy"
	SellerShopPolicies LinkType = "seller_shop_policies"
	TermsOfUse         LinkType = "terms_of_use"
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

// SupportedPaymentMethods defines model for PaymentProvider.SupportedPaymentMethods.
type SupportedPaymentMethods string

// Defines values for PaymentProviderSupportedPaymentMethods.
const (
	Card SupportedPaymentMethods = "card"
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
	ID                  string                `json:"id"`
	Buyer               *Buyer                `json:"buyer,omitempty"`
	Currency            string                `json:"currency"`
	FulfillmentAddress  *Address              `json:"fulfillment_address,omitempty"`
	FulfillmentOptionId *string               `json:"fulfillment_option_id,omitempty"`
	FulfillmentOptions  []FulfillmentOption   `json:"fulfillment_options"`
	LineItems           []LineItem            `json:"line_items"`
	Links               []Link                `json:"links"`
	Messages            []Message             `json:"messages"`
	PaymentProvider     *PaymentProvider      `json:"payment_provider,omitempty"`
	Status              CheckoutSessionStatus `json:"status"`
	Totals              []Total               `json:"totals"`
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
	Subtotal string  `json:"subtotal"`
	Tax      string  `json:"tax"`
	Title    string  `json:"title"`
	Total    string  `json:"total"`
	Type     string  `json:"type"`
}

// FulfillmentOptionShipping defines model for FulfillmentOptionShipping.
type FulfillmentOptionShipping struct {
	ID                   string     `json:"id"`
	Carrier              *string    `json:"carrier,omitempty"`
	EarliestDeliveryTime *time.Time `json:"earliest_delivery_time,omitempty"`
	LatestDeliveryTime   *time.Time `json:"latest_delivery_time,omitempty"`
	Subtitle             *string    `json:"subtitle,omitempty"`
	Subtotal             string     `json:"subtotal"`
	Tax                  string     `json:"tax"`
	Title                string     `json:"title"`
	Total                string     `json:"total"`
	Type                 string     `json:"type"`
}

// Item defines model for Item.
type Item struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

// LineItem defines model for LineItem.
type LineItem struct {
	ID         string `json:"id"`
	BaseAmount int    `json:"base_amount"`
	Discount   int    `json:"discount"`
	Item       Item   `json:"item"`
	Subtotal   int    `json:"subtotal"`
	Tax        int    `json:"tax"`
	Total      int    `json:"total"`
}

// Link defines model for Link.
type Link struct {
	Type LinkType `json:"type"`
	Url  string   `json:"url"`
}

// MessageInfo defines model for MessageInfo.
type MessageInfo struct {
	Content     string                 `json:"content"`
	ContentType MessageInfoContentType `json:"content_type"`

	// Param RFC 9535 JSONPath
	Param *string `json:"param,omitempty"`
	Type  string  `json:"type"`
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

// PaymentProvider defines model for PaymentProvider.
type PaymentProvider struct {
	Provider                PaymentProviderProvider   `json:"provider"`
	SupportedPaymentMethods []SupportedPaymentMethods `json:"supported_payment_methods"`
}

// PaymentProviderProvider defines model for PaymentProvider.Provider.
type PaymentProviderProvider string

// Total defines model for Total.
type Total struct {
	Amount      int       `json:"amount"`
	DisplayText string    `json:"display_text"`
	Type        TotalType `json:"type"`
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
