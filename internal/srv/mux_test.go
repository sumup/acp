package srv

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sumup/acp"
)

func TestWriteACPErrorSetsRetryAfter(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	err := acp.NewHTTPError(
		http.StatusConflict,
		acp.InvalidRequest,
		acp.IdempotencyInFlight,
		"request is still processing",
		acp.WithRetryAfter(1500*time.Millisecond),
	)
	if writeErr := WriteACPError(rec, err); writeErr != nil {
		t.Fatalf("WriteACPError() error = %v", writeErr)
	}
	if got := rec.Header().Get("Retry-After"); got != "2" {
		t.Fatalf("Retry-After = %q, want %q", got, "2")
	}
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d", rec.Code)
	}
}
