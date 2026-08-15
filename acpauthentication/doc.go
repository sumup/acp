// Package acpauthentication serves the ACP Delegated Authentication API.
//
// The API provides a session-based 3DS2 lifecycle: create a session, submit
// browser fingerprint results, then retrieve the final authentication result.
// Implement [Provider] and pass it to [NewHandler].
package acpauthentication
