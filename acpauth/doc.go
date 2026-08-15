// Package acpauth defines bearer-token authorization for ACP HTTP handlers.
//
// Implement [Authorizer] to connect a handler to production authentication.
// [StaticTokenAuthorizer] is intended for examples, tests, and other cases where
// one fixed token is sufficient.
package acpauth
