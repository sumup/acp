// Package acppayment serves the ACP Delegated Payment API.
//
// Delegated Payment tokenizes a credential for controlled use by the
// merchant's payment service provider under the request's [Allowance]. The
// stable specification currently supports card credentials.
//
// Implement [Provider] and pass it to [NewHandler]. See the [standalone
// example] for an in-memory implementation.
//
// [standalone example]: https://github.com/sumup/acp/tree/main/examples/delegated_payment
package acppayment
