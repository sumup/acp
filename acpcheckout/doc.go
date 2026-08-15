// Package acpcheckout serves the seller-hosted ACP Agentic Checkout API.
//
// Implement [Provider] with the business logic that owns checkout sessions,
// then pass it to [NewHandler]. The handler validates request payloads and the
// standard ACP headers before invoking the provider.
//
// The API creates, retrieves, updates, completes, and cancels checkout
// sessions. Every successful mutation returns the latest authoritative session
// state. See the [standalone example] for a complete provider skeleton.
//
// [standalone example]: https://github.com/sumup/acp/tree/main/examples/checkout
package acpcheckout
