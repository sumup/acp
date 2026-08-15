// Package acpdiscovery serves the public ACP discovery document.
//
// Mount [Handler] at the origin root. It responds only to GET
// /.well-known/acp.json, requires no authentication, and advertises stable
// seller capabilities before a checkout session is created.
package acpdiscovery
