package acpdiscovery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sumup/acp"
)

func TestHandlerServesDiscoveryDocument(t *testing.T) {
	t.Parallel()

	handler, err := NewHandler(Document{
		Protocol: Protocol{
			Name:              ProtocolName,
			Version:           acp.APIVersion,
			SupportedVersions: []string{acp.APIVersion},
		},
		APIBaseURL: "https://merchant.example/acp",
		Transports: []Transport{TransportREST},
		Capabilities: Capabilities{
			Services: []Service{ServiceCheckout, ServiceCarts},
		},
	})
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, Path, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "public, max-age=3600" {
		t.Fatalf("Cache-Control = %q", got)
	}
	var got Document
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Protocol.Version != acp.APIVersion {
		t.Fatalf("version = %q", got.Protocol.Version)
	}
}

func TestNewHandlerRejectsVersionOrderMismatch(t *testing.T) {
	t.Parallel()

	_, err := NewHandler(Document{
		Protocol: Protocol{
			Name:              ProtocolName,
			Version:           acp.APIVersion,
			SupportedVersions: []string{"2026-01-30"},
		},
		APIBaseURL: "https://merchant.example/acp",
		Capabilities: Capabilities{
			Services: []Service{ServiceCheckout},
		},
	})
	if err == nil {
		t.Fatal("NewHandler() error = nil")
	}
}
