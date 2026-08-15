// Package acp provides shared Agentic Commerce Protocol types and HTTP helpers.
//
// The protocol-specific packages are:
//
//   - [github.com/sumup/acp/acpcheckout] for seller-hosted checkout sessions;
//   - [github.com/sumup/acp/acpcart] for seller-hosted carts;
//   - [github.com/sumup/acp/acppayment] for delegated payment;
//   - [github.com/sumup/acp/acpauthentication] for delegated 3DS authentication;
//   - [github.com/sumup/acp/acpfeed] for the agent-hosted Feed API;
//   - [github.com/sumup/acp/acpdiscovery] for the public discovery document; and
//   - [github.com/sumup/acp/acpwebhook] for order events and webhook signing.
//
// [APIVersion] identifies the stable ACP specification implemented by this
// module. HTTP handlers validate that version and make the supported version
// available in protocol-level version errors.
package acp
