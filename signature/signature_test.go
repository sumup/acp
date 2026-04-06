package signature

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestReadAndBufferBody(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/checkout_sessions", strings.NewReader(`{"a":1}`))

	raw, err := ReadAndBufferBody(req)
	if err != nil {
		t.Fatalf("ReadAndBufferBody() error = %v", err)
	}
	if got, want := string(raw), `{"a":1}`; got != want {
		t.Fatalf("raw body = %q, want %q", got, want)
	}

	reread, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("reading buffered body: %v", err)
	}
	if got, want := string(reread), `{"a":1}`; got != want {
		t.Fatalf("buffered body = %q, want %q", got, want)
	}
}

func TestReadAndBufferBodyNilBody(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/checkout_sessions", nil)
	req.Body = nil

	raw, err := ReadAndBufferBody(req)
	if err != nil {
		t.Fatalf("ReadAndBufferBody() error = %v", err)
	}
	if raw != nil {
		t.Fatalf("raw body = %v, want nil", raw)
	}
	if req.Body == nil {
		t.Fatal("request body not restored")
	}
}

func TestCanonicalizeJSONBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr string
	}{
		{
			name: "canonicalizes object",
			raw:  "{\n  \"b\": 2,\n  \"a\": 1,\n  \"nested\": {\"d\": 4, \"c\": 3}\n}",
			want: `{"a":1,"b":2,"nested":{"c":3,"d":4}}`,
		},
		{
			name: "empty body becomes null",
			raw:  " \n\t ",
			want: "null",
		},
		{
			name:    "multiple documents rejected",
			raw:     "{}{}",
			wantErr: "multiple JSON documents",
		},
		{
			name:    "invalid json rejected",
			raw:     "{",
			wantErr: "unexpected EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := CanonicalizeJSONBody([]byte(tt.raw))
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("CanonicalizeJSONBody() error = %v, want substring %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("CanonicalizeJSONBody() error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("CanonicalizeJSONBody() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	t.Parallel()

	t.Run("accepts RFC3339", func(t *testing.T) {
		t.Parallel()

		got, err := ParseTimestamp("2025-09-25T10:30:00Z")
		if err != nil {
			t.Fatalf("ParseTimestamp() error = %v", err)
		}
		want := time.Date(2025, 9, 25, 10, 30, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Fatalf("ParseTimestamp() = %v, want %v", got, want)
		}
	})

	t.Run("accepts RFC3339Nano", func(t *testing.T) {
		t.Parallel()

		got, err := ParseTimestamp("2025-09-25T10:30:00.123456789Z")
		if err != nil {
			t.Fatalf("ParseTimestamp() error = %v", err)
		}
		want := time.Date(2025, 9, 25, 10, 30, 0, 123456789, time.UTC)
		if !got.Equal(want) {
			t.Fatalf("ParseTimestamp() = %v, want %v", got, want)
		}
	})

	t.Run("rejects empty value", func(t *testing.T) {
		t.Parallel()

		_, err := ParseTimestamp("")
		if err == nil || !strings.Contains(err.Error(), "empty timestamp") {
			t.Fatalf("ParseTimestamp() error = %v, want empty timestamp", err)
		}
	})
}

func TestAbsDuration(t *testing.T) {
	t.Parallel()

	if got, want := AbsDuration(-5*time.Second), 5*time.Second; got != want {
		t.Fatalf("AbsDuration(-5s) = %v, want %v", got, want)
	}
	if got, want := AbsDuration(3*time.Second), 3*time.Second; got != want {
		t.Fatalf("AbsDuration(3s) = %v, want %v", got, want)
	}
}

func TestBuildSigningPayload(t *testing.T) {
	t.Parallel()

	ts := time.Date(2025, 9, 25, 10, 30, 0, 123000000, time.UTC)
	got := BuildSigningPayload(ts, []byte(`{"a":1}`))
	want := `2025-09-25T10:30:00.123Z.{"a":1}`
	if string(got) != want {
		t.Fatalf("BuildSigningPayload() = %q, want %q", got, want)
	}
}

func TestHMACVerifierVerify(t *testing.T) {
	t.Parallel()

	key := []byte("top-secret")
	ts := time.Date(2025, 9, 25, 10, 30, 0, 0, time.UTC)
	body := []byte(`{"a":1}`)
	signingInput := BuildSigningPayload(ts, body)
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(signingInput); err != nil {
		t.Fatalf("mac.Write() error = %v", err)
	}
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	err := HMACVerifier{Key: key}.Verify(context.Background(), Material{
		Signature:     signature,
		Timestamp:     ts,
		CanonicalBody: body,
	})
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestHMACVerifierVerifyErrors(t *testing.T) {
	t.Parallel()

	ts := time.Date(2025, 9, 25, 10, 30, 0, 0, time.UTC)
	body := []byte(`{"a":1}`)

	tests := []struct {
		name     string
		verifier HMACVerifier
		material Material
		wantErr  string
	}{
		{
			name:     "empty key",
			verifier: HMACVerifier{},
			material: Material{Timestamp: ts, CanonicalBody: body},
			wantErr:  "requires a non-empty key",
		},
		{
			name:     "invalid base64 signature",
			verifier: HMACVerifier{Key: []byte("top-secret")},
			material: Material{Signature: "!!!", Timestamp: ts, CanonicalBody: body},
			wantErr:  "decode signature",
		},
		{
			name:     "signature mismatch",
			verifier: HMACVerifier{Key: []byte("top-secret")},
			material: Material{Signature: base64.RawURLEncoding.EncodeToString([]byte("wrong")), Timestamp: ts, CanonicalBody: body},
			wantErr:  "invalid signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.verifier.Verify(context.Background(), tt.material)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Verify() error = %v, want substring %q", err, tt.wantErr)
			}
		})
	}
}
