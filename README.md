<div align="center">

# Go SDK for the Agentic Commerce Protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/sumup/acp.svg)](https://pkg.go.dev/github.com/sumup/acp)
[![CI Status](https://github.com/sumup/acp/workflows/CI/badge.svg)](https://github.com/sumup/acp/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/sumup/acp)](./LICENSE)

</div>

Go SDK for the [Agentic Commerce Protocol](https://github.com/agentic-commerce-protocol/agentic-commerce-protocol) (ACP), targeting the stable API version [`2026-04-17`](https://github.com/agentic-commerce-protocol/agentic-commerce-protocol/tree/main/spec/2026-04-17).

```bash
go get github.com/sumup/acp
```

## Packages

| Package | Use |
| --- | --- |
| [`acpcheckout`](https://pkg.go.dev/github.com/sumup/acp/acpcheckout) | Serve seller-hosted checkout sessions. |
| [`acpcart`](https://pkg.go.dev/github.com/sumup/acp/acpcart) | Serve seller-hosted pre-checkout carts. |
| [`acppayment`](https://pkg.go.dev/github.com/sumup/acp/acppayment) | Serve delegated payment tokenization. |
| [`acpauthentication`](https://pkg.go.dev/github.com/sumup/acp/acpauthentication) | Serve delegated 3DS authentication. |
| [`acpfeed`](https://pkg.go.dev/github.com/sumup/acp/acpfeed) | Call an agent-hosted Feed API. |
| [`acpdiscovery`](https://pkg.go.dev/github.com/sumup/acp/acpdiscovery) | Serve `/.well-known/acp.json`. |
| [`acpwebhook`](https://pkg.go.dev/github.com/sumup/acp/acpwebhook) | Send signed order lifecycle events. |
| [`acpauth`](https://pkg.go.dev/github.com/sumup/acp/acpauth) | Implement bearer-token authorization. |
| [`signature`](https://pkg.go.dev/github.com/sumup/acp/signature) | Verify optional signed ACP requests. |
| [`discount`](https://pkg.go.dev/github.com/sumup/acp/discount), [`extension`](https://pkg.go.dev/github.com/sumup/acp/extension) | Use generated extension models. |

## Server handlers

The checkout, cart, payment, and authentication packages expose a `Provider` interface and a `net/http` handler. Implement the interface, supply an `acpauth.Authorizer`, and mount the handler directly or into an existing `http.ServeMux` with `WithServeMux`.

These handlers validate JSON payloads and ACP headers. Send `API-Version: 2026-04-17` and `Authorization: Bearer <token>` on authenticated service requests. Mutating checkout, cart, and delegated-payment requests also require `Idempotency-Key`; the OpenAPI definitions describe endpoint-specific exceptions. Discovery is public and does not use those headers.

### Checkout example

Run the in-memory example:

```bash
go run ./examples/checkout
```

Create and complete a session with spec-valid minimal payloads:

```bash
curl -sS -X POST http://localhost:8080/checkout_sessions \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-key' \
  -H 'API-Version: 2026-04-17' \
  -H 'Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000' \
  -d '{
        "line_items": [{"id": "latte"}],
        "currency": "eur",
        "capabilities": {}
      }'

curl -sS -X POST http://localhost:8080/checkout_sessions/<session_id>/complete \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-key' \
  -H 'API-Version: 2026-04-17' \
  -H 'Idempotency-Key: 550e8400-e29b-41d4-a716-446655440001' \
  -d '{
        "payment_data": {
          "handler_id": "card_tokenized",
          "instrument": {
            "type": "card",
            "credential": {"type": "spt", "token": "spt_demo"}
          }
        }
      }'
```

### Delegated payment example

```bash
go run ./examples/delegated_payment
```

In another terminal:

```bash
curl -sS -X POST http://localhost:8080/agentic_commerce/delegate_payment \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-key' \
  -H 'API-Version: 2026-04-17' \
  -H 'Idempotency-Key: 550e8400-e29b-41d4-a716-446655440002' \
  -d '{
        "payment_method": {
          "type": "card",
          "card_number_type": "fpan",
          "number": "4242424242424242",
          "exp_month": "11",
          "exp_year": "2030",
          "display_last4": "4242",
          "display_card_funding_type": "credit",
          "metadata": {"issuer": "demo-bank"}
        },
        "allowance": {
          "reason": "one_time",
          "max_amount": 2000,
          "currency": "eur",
          "checkout_session_id": "cs_000001",
          "merchant_id": "demo-merchant",
          "expires_at": "2027-12-31T23:59:59Z"
        },
        "risk_signals": [],
        "metadata": {"source": "sample"}
      }'
```

## Feed client

The stable Feed API is hosted by agents and called by merchants. It is a push model: merchants create feed metadata and partially upsert products by `Product.id`.

```go
country := "US"
client, err := acpfeed.NewClientWithResponses(agentURL)
if err != nil {
	return err
}

response, err := client.CreateFeedWithResponse(ctx, acpfeed.CreateFeedRequest{
	TargetCountry: &country,
})
if err != nil {
	return err
}
if response.JSON201 == nil {
	return fmt.Errorf("create feed: %s", response.Status())
}
```

See the executable examples on [pkg.go.dev](https://pkg.go.dev/github.com/sumup/acp) for discovery, webhooks, request signing, and feed upserts.

## Development

Generated files are built from the vendored stable ACP specifications. Edit the specs or generator configuration, then run:

```bash
make generate
make fmt
make lint
make test
```

## License

[Apache 2.0](./LICENSE)
