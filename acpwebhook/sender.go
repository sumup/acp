package acpwebhook

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
	"net/url"
	"strconv"
	"strings"
	"time"
)

const signatureHeader = "Merchant-Signature"

// Option configures a [Sender].
type Option func(*Sender)

// WithHTTPClient overrides the HTTP client used to deliver webhook events.
func WithHTTPClient(client *http.Client) Option {
	return func(s *Sender) {
		s.client = client
	}
}

func withClock(now func() time.Time) Option {
	return func(s *Sender) {
		s.now = now
	}
}

// Sender delivers signed ACP order webhook events to an agent endpoint.
type Sender struct {
	endpoint *url.URL
	secret   []byte
	client   *http.Client
	now      func() time.Time
}

// NewSender validates the endpoint and returns a webhook sender.
func NewSender(endpoint string, secret []byte, opts ...Option) (*Sender, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("webhook: parse endpoint: %w", err)
	}
	if !endpointURL.IsAbs() || endpointURL.Host == "" || (endpointURL.Scheme != "http" && endpointURL.Scheme != "https") {
		return nil, errors.New("webhook: endpoint must be an absolute HTTP(S) URL")
	}
	if len(secret) == 0 {
		return nil, errors.New("webhook: secret is required")
	}

	sender := &Sender{
		endpoint: endpointURL,
		secret:   append([]byte(nil), secret...),
		client:   http.DefaultClient,
		now:      time.Now,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(sender)
		}
	}
	if sender.client == nil {
		return nil, errors.New("webhook: HTTP client is required")
	}
	if sender.now == nil {
		return nil, errors.New("webhook: clock is required")
	}
	return sender, nil
}

// Send posts event using the ACP Merchant-Signature format.
func (s *Sender) Send(ctx context.Context, event WebhookEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("webhook: marshal event: %w", err)
	}

	timestamp := s.now().Unix()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(signatureHeader, sign(s.secret, timestamp, body))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send event: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("webhook: endpoint returned %s: %s", resp.Status, strings.TrimSpace(string(snippet)))
	}
	return nil
}

func sign(secret []byte, timestamp int64, body []byte) string {
	seconds := strconv.FormatInt(timestamp, 10)
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(seconds))
	_, _ = mac.Write([]byte("."))
	_, _ = mac.Write(body)
	return fmt.Sprintf("t=%s,v1=%x", seconds, mac.Sum(nil))
}
