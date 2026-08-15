package acpfeed

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientCreateFeed(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/feeds" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"feed_1","target_country":"US"}`))
	}))
	t.Cleanup(server.Close)

	client, err := NewClientWithResponses(server.URL)
	if err != nil {
		t.Fatalf("NewClientWithResponses() error = %v", err)
	}
	country := "US"
	response, err := client.CreateFeedWithResponse(context.Background(), CreateFeedRequest{TargetCountry: &country})
	if err != nil {
		t.Fatalf("CreateFeedWithResponse() error = %v", err)
	}
	if response.JSON201 == nil || response.JSON201.Id != "feed_1" {
		t.Fatalf("unexpected response: %#v", response)
	}
}
