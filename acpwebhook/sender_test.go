package acpwebhook_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/sumup/acp/acpwebhook"
)

func TestSender_Send(t *testing.T) {
	t.Parallel()

	t.Run("signs the raw event with the current timestamp", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			secret := []byte("test-secret")
			key := bytes.Clone(secret)
			var received bool
			client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				received = true
				if r.Method != http.MethodPost || r.URL.String() != "https://agent.example/webhooks" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					return nil, err
				}
				mac := hmac.New(sha256.New, key)
				_, _ = fmt.Fprintf(mac, "%d.%s", time.Now().Unix(), body)
				want := fmt.Sprintf("t=%d,v1=%x", time.Now().Unix(), mac.Sum(nil))
				if got := r.Header.Get("Merchant-Signature"); got != want {
					t.Errorf("Merchant-Signature = %q, want %q", got, want)
				}
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("Content-Type = %q", got)
				}
				var event acpwebhook.WebhookEvent
				if err := json.Unmarshal(body, &event); err != nil {
					return nil, err
				}
				if event.Type != "order_create" || event.Data.Id != "order_1" {
					t.Errorf("event = %#v", event)
				}
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			})}
			sender, err := acpwebhook.NewSender("https://agent.example/webhooks", secret, acpwebhook.WithHTTPClient(client))
			if err != nil {
				t.Fatal(err)
			}
			secret[0] = 'X' // The sender must own a copy of the signing key.
			err = sender.Send(t.Context(), acpwebhook.WebhookEvent{
				Type: "order_create",
				Data: acpwebhook.EventDataOrder{
					Type: acpwebhook.EventDataOrderTypeOrder, Id: "order_1", CheckoutSessionId: "checkout_1",
					PermalinkUrl: "https://merchant.example/orders/order_1", Status: "created",
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			if !received {
				t.Fatal("webhook was not sent")
			}
		})
	})

	t.Run("respects the caller deadline", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				<-r.Context().Done()
				return nil, r.Context().Err()
			})}
			sender, err := acpwebhook.NewSender("https://agent.example/webhooks", []byte("secret"), acpwebhook.WithHTTPClient(client))
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()
			start := time.Now()
			err = sender.Send(ctx, acpwebhook.WebhookEvent{})
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("Send() error = %v, want deadline exceeded", err)
			}
			if elapsed := time.Since(start); elapsed != time.Second {
				t.Fatalf("Send() duration = %s", elapsed)
			}
		})
	})

	t.Run("reports non-success responses", func(t *testing.T) {
		t.Parallel()

		for _, status := range []int{http.StatusBadRequest, http.StatusTooManyRequests, http.StatusInternalServerError} {
			t.Run(http.StatusText(status), func(t *testing.T) {
				t.Parallel()

				client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Body: io.NopCloser(strings.NewReader("try later"))}, nil
				})}
				sender, err := acpwebhook.NewSender("https://agent.example/webhooks", []byte("secret"), acpwebhook.WithHTTPClient(client))
				if err != nil {
					t.Fatal(err)
				}
				if err := sender.Send(t.Context(), acpwebhook.WebhookEvent{}); err == nil || !strings.Contains(err.Error(), "try later") {
					t.Fatalf("Send() error = %v", err)
				}
			})
		}
	})
}

func TestNewSender(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		endpoint string
		secret   []byte
		opts     []acpwebhook.Option
	}{
		{"relative endpoint", "/webhooks", []byte("secret"), nil},
		{"unsupported scheme", "ftp://agent.example/webhooks", []byte("secret"), nil},
		{"empty key", "https://agent.example/webhooks", nil, nil},
		{"nil client", "https://agent.example/webhooks", []byte("secret"), []acpwebhook.Option{acpwebhook.WithHTTPClient(nil)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if _, err := acpwebhook.NewSender(tc.endpoint, tc.secret, tc.opts...); err == nil {
				t.Fatal("NewSender() error = nil")
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
