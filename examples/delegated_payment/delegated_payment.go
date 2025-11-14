package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sumup/acp"
)

func main() {
	service := newDelegatedMemoryService()
	addr := ":8080"

	log.Printf("ACP delegated payment sample listening on %s", addr)
	log.Printf("Try: curl -XPOST %s/agentic_commerce/delegate_payment -d @- <<'JSON' ...", "http://localhost:8080")

	handler := acp.NewDelegatedPaymentHandler(service)
	log.Fatal(http.ListenAndServe(addr, cors(logging(handler))))
}

// logging adds basic request logs without external dependencies.
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start).Truncate(time.Millisecond))
	})
}

// cors allows the browser-based testbed to call the sample server directly.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept,API-Version")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code before forwarding to the real writer.
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

type delegatedMemoryService struct {
	mu      sync.Mutex
	tokens  map[string]*acp.VaultToken
	tokenID uint64
}

func newDelegatedMemoryService() *delegatedMemoryService {
	return &delegatedMemoryService{
		tokens: make(map[string]*acp.VaultToken),
	}
}

// DelegatePayment issues idempotent tokens keyed by checkout_session_id.
func (s *delegatedMemoryService) DelegatePayment(_ context.Context, req acp.PaymentRequest) (*acp.VaultToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := req.Allowance.CheckoutSessionID
	if token, ok := s.tokens[key]; ok {
		return cloneVaultToken(token), nil
	}

	metadata := cloneStringMap(req.Metadata)
	if metadata == nil {
		metadata = make(map[string]string, 2)
	}
	metadata["merchant_id"] = req.Allowance.MerchantID
	metadata["checkout_session_id"] = key

	token := &acp.VaultToken{
		ID:       s.nextTokenID(),
		Created:  time.Now().UTC(),
		Metadata: metadata,
	}

	s.tokens[key] = token
	return cloneVaultToken(token), nil
}

func (s *delegatedMemoryService) nextTokenID() string {
	id := atomic.AddUint64(&s.tokenID, 1)
	return fmt.Sprintf("vt_%06d", id)
}

func cloneVaultToken(src *acp.VaultToken) *acp.VaultToken {
	if src == nil {
		return nil
	}
	dst := *src
	dst.Metadata = cloneStringMap(src.Metadata)
	return &dst
}

func cloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
