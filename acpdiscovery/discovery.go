package acpdiscovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// Path is the RFC 8615 well-known location for ACP discovery.
	Path = "/.well-known/acp.json"
	// ProtocolName is the required ACP protocol identifier.
	ProtocolName = "acp"
)

// Transport is a supported ACP transport binding.
type Transport string

const (
	TransportREST Transport = "rest"
	TransportMCP  Transport = "mcp"
)

// Service is an ACP service advertised by a seller.
type Service string

const (
	ServiceCheckout        Service = "checkout"
	ServiceOrders          Service = "orders"
	ServiceDelegatePayment Service = "delegate_payment"
	ServiceCarts           Service = "carts"
)

// InterventionType is a buyer-intervention capability.
type InterventionType string

const (
	Intervention3DS                 InterventionType = "3ds"
	InterventionBiometric           InterventionType = "biometric"
	InterventionAddressVerification InterventionType = "address_verification"
)

// Protocol identifies ACP versions supported by the seller.
type Protocol struct {
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	SupportedVersions []string `json:"supported_versions"`
	DocumentationURL  *string  `json:"documentation_url,omitempty"`
}

// Extension declares a seller-supported ACP extension.
type Extension struct {
	Name   string  `json:"name"`
	Spec   *string `json:"spec,omitempty"`
	Schema *string `json:"schema,omitempty"`
}

// Capabilities contains seller capabilities that are stable across sessions.
type Capabilities struct {
	Services            []Service          `json:"services"`
	Extensions          []Extension        `json:"extensions,omitempty"`
	InterventionTypes   []InterventionType `json:"intervention_types,omitempty"`
	SupportedCurrencies []string           `json:"supported_currencies,omitempty"`
	SupportedLocales    []string           `json:"supported_locales,omitempty"`
}

// Document is the response served from [Path].
type Document struct {
	Protocol     Protocol     `json:"protocol"`
	APIBaseURL   string       `json:"api_base_url"`
	Transports   []Transport  `json:"transports"`
	Capabilities Capabilities `json:"capabilities"`
}

// Handler serves an immutable discovery document.
type Handler struct {
	body []byte
}

// NewHandler validates and serializes document for concurrent serving.
func NewHandler(document Document) (*Handler, error) {
	if err := validate(document); err != nil {
		return nil, err
	}
	body, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("discovery: marshal document: %w", err)
	}
	return &Handler{body: body}, nil
}

// ServeHTTP serves the discovery document at [Path].
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || r.URL.Path != Path {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.body)
}

func validate(document Document) error {
	if document.Protocol.Name != ProtocolName {
		return fmt.Errorf("discovery: protocol name must be %q", ProtocolName)
	}
	if document.Protocol.Version == "" {
		return errors.New("discovery: protocol version is required")
	}
	if len(document.Protocol.SupportedVersions) == 0 {
		return errors.New("discovery: at least one supported version is required")
	}
	if latest := document.Protocol.SupportedVersions[len(document.Protocol.SupportedVersions)-1]; latest != document.Protocol.Version {
		return errors.New("discovery: protocol version must be the last supported version")
	}
	baseURL, err := url.Parse(document.APIBaseURL)
	if err != nil || !baseURL.IsAbs() || baseURL.Host == "" {
		return errors.New("discovery: API base URL must be absolute")
	}
	if len(document.Capabilities.Services) == 0 {
		return errors.New("discovery: at least one service is required")
	}
	return nil
}
