package acpdiscovery_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/sumup/acp"
	"github.com/sumup/acp/acpdiscovery"
)

func ExampleNewHandler() {
	handler, err := acpdiscovery.NewHandler(acpdiscovery.Document{
		Protocol: acpdiscovery.Protocol{
			Name:              acpdiscovery.ProtocolName,
			Version:           acp.APIVersion,
			SupportedVersions: []string{acp.APIVersion},
		},
		APIBaseURL: "https://merchant.example/acp",
		Transports: []acpdiscovery.Transport{acpdiscovery.TransportREST},
		Capabilities: acpdiscovery.Capabilities{
			Services: []acpdiscovery.Service{
				acpdiscovery.ServiceCheckout,
				acpdiscovery.ServiceCarts,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, acpdiscovery.Path, nil))

	var document acpdiscovery.Document
	if err := json.NewDecoder(recorder.Body).Decode(&document); err != nil {
		panic(err)
	}
	fmt.Println(recorder.Code, document.Protocol.Version, document.Capabilities.Services)
	// Output: 200 2026-04-17 [checkout carts]
}
