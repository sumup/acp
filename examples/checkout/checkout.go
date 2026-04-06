package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acpcheckout"
)

type memoryService struct{}

func (memoryService) CreateSession(_ context.Context, req acpcheckout.CheckoutSessionCreateRequest) (*acpcheckout.CheckoutSessionBase, error) {
	return &acpcheckout.CheckoutSessionBase{
		Id:                 "cs_demo_1",
		Currency:           req.Currency,
		Capabilities:       req.Capabilities,
		Status:             acpcheckout.CheckoutSessionBaseStatusInProgress,
		FulfillmentOptions: []acpcheckout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []acpcheckout.LineItem{},
		Links:              []acpcheckout.Link{},
		Messages:           []acpcheckout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func (memoryService) UpdateSession(_ context.Context, id string, _ acpcheckout.CheckoutSessionUpdateRequest) (*acpcheckout.CheckoutSessionBase, error) {
	return &acpcheckout.CheckoutSessionBase{
		Id:                 id,
		Currency:           "USD",
		Status:             acpcheckout.CheckoutSessionBaseStatusInProgress,
		FulfillmentOptions: []acpcheckout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []acpcheckout.LineItem{},
		Links:              []acpcheckout.Link{},
		Messages:           []acpcheckout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func (memoryService) GetSession(_ context.Context, id string) (*acpcheckout.CheckoutSessionBase, error) {
	return &acpcheckout.CheckoutSessionBase{
		Id:                 id,
		Currency:           "USD",
		Status:             acpcheckout.CheckoutSessionBaseStatusReadyForPayment,
		FulfillmentOptions: []acpcheckout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []acpcheckout.LineItem{},
		Links:              []acpcheckout.Link{},
		Messages:           []acpcheckout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func (memoryService) CompleteSession(_ context.Context, id string, _ acpcheckout.CheckoutSessionCompleteRequest) (acpcheckout.CheckoutSessionWithOrder, error) {
	return acpcheckout.CheckoutSessionWithOrder{
		Id:                 id,
		Currency:           "USD",
		Capabilities:       acpcheckout.Capabilities{},
		FulfillmentOptions: []acpcheckout.CheckoutSessionWithOrder_FulfillmentOptions_Item{},
		LineItems:          []acpcheckout.LineItem{},
		Links:              []acpcheckout.Link{},
		Messages:           []acpcheckout.CheckoutSessionWithOrder_Messages_Item{},
		Status:             acpcheckout.CheckoutSessionWithOrderStatusCompleted,
		Totals:             []acpcheckout.Total{},
		Order: acpcheckout.Order{
			Id:                "ord_demo_1",
			CheckoutSessionId: id,
			PermalinkUrl:      "https://example.com/orders/ord_demo_1",
		},
	}, nil
}

func (memoryService) CancelSession(_ context.Context, id string, _ *acpcheckout.CancelSessionRequest) (*acpcheckout.CheckoutSessionBase, error) {
	return &acpcheckout.CheckoutSessionBase{
		Id:                 id,
		Currency:           "USD",
		Status:             acpcheckout.CheckoutSessionBaseStatusCanceled,
		FulfillmentOptions: []acpcheckout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []acpcheckout.LineItem{},
		Links:              []acpcheckout.Link{},
		Messages:           []acpcheckout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func main() {
	handler := acpcheckout.NewCheckoutHandler(memoryService{}, acpauth.StaticTokenAuthorizer("demo-key"))
	log.Println("checkout example listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
