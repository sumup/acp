// Package acpfeed provides models and an HTTP client for the ACP Feed API.
//
// The Feed API is hosted by agents and called by merchants. It complements the
// legacy pull-feed encoders in [github.com/sumup/acp/feed] with the current ACP
// push model for feed metadata and product upserts. Products omitted from an
// upsert remain unchanged.
//
// [NewClientWithResponses] returns a client that decodes each documented JSON
// response into the corresponding JSON200, JSON201, JSON400, or JSON404 field.
package acpfeed
