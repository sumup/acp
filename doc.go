// Package acp documents the Go SDK for the Agentic Commerce Protocol (ACP).
// It aggregates the checkout and delegated payment packages under a single module
// so merchants and PSPs can share common helpers, models, and documentation.
//
// # Checkout
//
// Use [NewCheckoutHandler] with your [CheckoutProvider] implementation to
// expose the ACP checkout contract over `net/http`. Handler options such as
// [WithSignatureVerifier] and [WithRequireSignedRequests] enforce the
// canonical JSON signatures and timestamp skew requirements spelled out in the
// spec.
//
// # Delegated Payment
//
// Payment service providers can call [NewDelegatedPaymentHandler] with their own
// [DelegatedPaymentProvider] to accept delegate payment payloads, validate them,
// and emit vault tokens scoped to a checkout session. Optional helpers such as
// [WithAuthenticator] and
// [DelegatedPaymentWithSignatureVerifier] keep API keys and signed requests in
// sync with ACP's security requirements.
//
// ## How it works
//
//   - Buyers check out using their preferred payment method and save it in ChatGPT.
//   - The delegated payment payload is sent to the merchant’s PSP or vault directly. The delegated payment is single-use and set with allowances.
//   - The PSP or vault returns a payment token scoped to the delegated payment outside of PCI scope.
//   - OpenAI forwards the token during the complete-checkout call to enable the merchant to complete the transaction.
package acp
