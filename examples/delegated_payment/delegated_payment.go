package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sumup/acp/acpauth"
	"github.com/sumup/acp/acppayment"
)

type delegatedMemoryService struct {
	mu     sync.Mutex
	tokens map[string]*acppayment.DelegatePaymentResponse
}

func newService() *delegatedMemoryService {
	return &delegatedMemoryService{tokens: map[string]*acppayment.DelegatePaymentResponse{}}
}

func (s *delegatedMemoryService) DelegatePayment(_ context.Context, req acppayment.DelegatePaymentRequest) (*acppayment.DelegatePaymentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := req.Allowance.CheckoutSessionId
	if tok, ok := s.tokens[key]; ok {
		return &acppayment.DelegatePaymentResponse{Id: tok.Id, Created: tok.Created, Metadata: cloneMap(tok.Metadata)}, nil
	}

	resp := &acppayment.DelegatePaymentResponse{
		Id:      "vt_demo_1",
		Created: time.Now().UTC(),
		Metadata: map[string]string{
			"merchant_id": req.Allowance.MerchantId,
			"source":      "example",
		},
	}
	s.tokens[key] = resp
	return &acppayment.DelegatePaymentResponse{Id: resp.Id, Created: resp.Created, Metadata: cloneMap(resp.Metadata)}, nil
}

func cloneMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func main() {
	handler := acppayment.NewHandler(newService(), acpauth.StaticTokenAuthorizer("demo-key"))
	log.Println("delegated payment example listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
