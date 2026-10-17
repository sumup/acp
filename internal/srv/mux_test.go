package srv

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"one object", `{"value":1}`, false},
		{"trailing whitespace", "{\"value\":1}\n\t ", false},
		{"empty body", "", true},
		{"two objects", `{} {}`, true},
		{"trailing closing object", `{"value":1}}`, true},
		{"trailing closing array", `{"value":1}]`, true},
		{"trailing garbage", `{"value":1} garbage`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got map[string]int
			err := DecodeJSON(io.NopCloser(strings.NewReader(tc.body)), &got)
			if (err != nil) != tc.wantErr {
				t.Fatalf("DecodeJSON() error = %v, want error = %v", err, tc.wantErr)
			}
			if !tc.wantErr && got["value"] != 1 {
				t.Fatalf("decoded body = %v", got)
			}
		})
	}

	t.Run("nil body", func(t *testing.T) {
		t.Parallel()

		var got any
		if err := DecodeJSON(nil, &got); !errors.Is(err, io.EOF) {
			t.Fatalf("DecodeJSON(nil) error = %v, want EOF", err)
		}
	})
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		id   string
		key  string
	}{
		{"echo correlation and idempotency", "req-123", "idem-123"},
		{"optional headers absent", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
			req.Header.Set("Request-Id", tc.id)
			req.Header.Set("Idempotency-Key", tc.key)
			rec := httptest.NewRecorder()
			if err := WriteJSON(rec, req, http.StatusCreated, map[string]string{"id": "123"}); err != nil {
				t.Fatal(err)
			}
			if rec.Code != http.StatusCreated {
				t.Fatalf("status = %d", rec.Code)
			}
			if rec.Header().Get("Request-Id") != tc.id || rec.Header().Get("Idempotency-Key") != tc.key || rec.Header().Get("API-Version") != acp.APIVersion {
				t.Fatalf("headers = %v", rec.Header())
			}
		})
	}
}
