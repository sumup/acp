// Package acpwebhook sends ACP order lifecycle events to agent endpoints.
//
// [Sender] signs each raw JSON body with HMAC-SHA256 and sets
// Merchant-Signature to "t=<unix_seconds>,v1=<hex_digest>". The signed payload
// is "timestamp.raw_body", as required by ACP 2026-04-17.
package acpwebhook
