package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sumup/acp/agentic_checkout"
)

type memoryService struct{}

func (memoryService) CreateSession(_ context.Context, req agentic_checkout.CheckoutSessionCreateRequest) (*agentic_checkout.CheckoutSessionBase, error) {
	return &agentic_checkout.CheckoutSessionBase{
		Id:                 "cs_demo_1",
		Currency:           req.Currency,
		Capabilities:       req.Capabilities,
		Status:             agentic_checkout.CheckoutSessionBaseStatusInProgress,
		FulfillmentOptions: []agentic_checkout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []agentic_checkout.LineItem{},
		Links:              []agentic_checkout.Link{},
		Messages:           []agentic_checkout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func (memoryService) UpdateSession(_ context.Context, id string, _ agentic_checkout.CheckoutSessionUpdateRequest) (*agentic_checkout.CheckoutSessionBase, error) {
	return &agentic_checkout.CheckoutSessionBase{
		Id:                 id,
		Currency:           "USD",
		Status:             agentic_checkout.CheckoutSessionBaseStatusInProgress,
		FulfillmentOptions: []agentic_checkout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []agentic_checkout.LineItem{},
		Links:              []agentic_checkout.Link{},
		Messages:           []agentic_checkout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func (memoryService) GetSession(_ context.Context, id string) (*agentic_checkout.CheckoutSessionBase, error) {
	return &agentic_checkout.CheckoutSessionBase{
		Id:                 id,
		Currency:           "USD",
		Status:             agentic_checkout.CheckoutSessionBaseStatusReadyForPayment,
		FulfillmentOptions: []agentic_checkout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []agentic_checkout.LineItem{},
		Links:              []agentic_checkout.Link{},
		Messages:           []agentic_checkout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func (memoryService) CompleteSession(_ context.Context, id string, _ agentic_checkout.CheckoutSessionCompleteRequest) (agentic_checkout.CheckoutSessionWithOrder, error) {
	return agentic_checkout.CheckoutSessionWithOrder{
		Id:                 id,
		Currency:           "USD",
		Capabilities:       agentic_checkout.Capabilities{},
		FulfillmentOptions: []agentic_checkout.CheckoutSessionWithOrder_FulfillmentOptions_Item{},
		LineItems:          []agentic_checkout.LineItem{},
		Links:              []agentic_checkout.Link{},
		Messages:           []agentic_checkout.CheckoutSessionWithOrder_Messages_Item{},
		Status:             agentic_checkout.CheckoutSessionWithOrderStatusCompleted,
		Totals:             []agentic_checkout.Total{},
		Order: agentic_checkout.Order{
			Id:                "ord_demo_1",
			CheckoutSessionId: id,
			PermalinkUrl:      "https://example.com/orders/ord_demo_1",
		},
	}, nil
}

func (memoryService) CancelSession(_ context.Context, id string, _ *agentic_checkout.CancelSessionRequest) (*agentic_checkout.CheckoutSessionBase, error) {
	return &agentic_checkout.CheckoutSessionBase{
		Id:                 id,
		Currency:           "USD",
		Status:             agentic_checkout.CheckoutSessionBaseStatusCanceled,
		FulfillmentOptions: []agentic_checkout.CheckoutSessionBase_FulfillmentOptions_Item{},
		LineItems:          []agentic_checkout.LineItem{},
		Links:              []agentic_checkout.Link{},
		Messages:           []agentic_checkout.CheckoutSessionBase_Messages_Item{},
	}, nil
}

func main() {
	handler := agentic_checkout.NewCheckoutHandler(memoryService{})
	log.Println("checkout example listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
