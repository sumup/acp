package acpwebhook_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/sumup/acp/acpwebhook"
)

func ExampleNewSender() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var event acpwebhook.WebhookEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			panic(err)
		}
		fmt.Println(event.Type, strings.HasPrefix(r.Header.Get("Merchant-Signature"), "t="))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	sender, err := acpwebhook.NewSender(server.URL, []byte("webhook-secret"))
	if err != nil {
		panic(err)
	}
	err = sender.Send(context.Background(), acpwebhook.WebhookEvent{
		Type: "order_create",
		Data: acpwebhook.EventDataOrder{
			Type:              acpwebhook.EventDataOrderTypeOrder,
			Id:                "ord_123",
			CheckoutSessionId: "cs_123",
			PermalinkUrl:      "https://merchant.example/orders/ord_123",
		},
	})
	if err != nil {
		panic(err)
	}

	// Output: order_create true
}
