// Package signature provides low-level helpers for optional signed ACP request
// headers.
//
// [HMACVerifier] verifies a base64url HMAC-SHA256 signature over an RFC 3339
// timestamp and canonical JSON body. Merchant-Signature on order webhooks uses
// a different wire format; use [github.com/sumup/acp/acpwebhook] to send those
// events.
package signature
