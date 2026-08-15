package acpfeed_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/sumup/acp/acpfeed"
)

func ExampleNewClientWithResponses() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"feed_8f3K2x","target_country":"US"}`))
	}))
	defer server.Close()

	client, err := acpfeed.NewClientWithResponses(server.URL)
	if err != nil {
		panic(err)
	}
	country := "US"
	response, err := client.CreateFeedWithResponse(context.Background(), acpfeed.CreateFeedRequest{
		TargetCountry: &country,
	})
	if err != nil {
		panic(err)
	}
	if response.JSON201 == nil {
		panic(fmt.Errorf("create feed: %s", response.Status()))
	}

	fmt.Println(response.JSON201.Id, *response.JSON201.TargetCountry)
	// Output: feed_8f3K2x US
}

func ExampleClientWithResponses_UpsertFeedProductsWithResponse() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"feed_8f3K2x","accepted":true}`))
	}))
	defer server.Close()

	client, err := acpfeed.NewClientWithResponses(server.URL)
	if err != nil {
		panic(err)
	}
	response, err := client.UpsertFeedProductsWithResponse(
		context.Background(),
		"feed_8f3K2x",
		acpfeed.UpsertProductsRequest{
			Products: []acpfeed.Product{
				{
					Id: "prod_classic_tee",
					Variants: []acpfeed.Variant{
						{Id: "sku-red-m", Title: "Classic Tee - Red / Medium"},
					},
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	if response.JSON200 == nil {
		panic(fmt.Errorf("upsert products: %s", response.Status()))
	}

	fmt.Println(response.JSON200.Id, response.JSON200.Accepted)
	// Output: feed_8f3K2x true
}
