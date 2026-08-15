// Package acpfeed provides models and an HTTP client for the ACP Feed API.
//
// The Feed API is hosted by agents and called by merchants. Merchants create
// feed metadata and push product upserts to the agent. Products omitted from an
// upsert remain unchanged.
//
// [NewClientWithResponses] returns a client that decodes each documented JSON
// response into the corresponding JSON200, JSON201, JSON400, or JSON404 field.
package acpfeed
