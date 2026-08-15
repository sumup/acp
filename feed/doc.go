// Package feed encodes and decodes the legacy pull-based product feed format
// as gzip-compressed CSV or JSON Lines.
//
// Deprecated: ACP 2026-04-17 uses a push model in which merchants create feed
// metadata and upsert products into an agent-hosted API. New ACP integrations
// should use [github.com/sumup/acp/acpfeed]. This package remains available for
// integrations that still consume the older flat [Product] record.
package feed
