package signature_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/sumup/acp/signature"
)

func ExampleHMACVerifier() {
	key := []byte("request-signing-secret")
	timestamp := time.Date(2026, 4, 17, 10, 30, 0, 0, time.UTC)
	body, err := signature.CanonicalizeJSONBody([]byte(`{"currency":"usd","line_items":[]}`))
	if err != nil {
		panic(err)
	}

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(signature.BuildSigningPayload(timestamp, body))
	encoded := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	err = (signature.HMACVerifier{Key: key}).Verify(context.Background(), signature.Material{
		Signature:     encoded,
		Timestamp:     timestamp,
		CanonicalBody: body,
	})
	fmt.Println(err)
	// Output: <nil>
}
