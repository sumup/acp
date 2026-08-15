<div align="center">

# Go SDK for the Agentic Commerce Protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/sumup/acp.svg)](https://pkg.go.dev/github.com/sumup/acp)
[![CI Status](https://github.com/sumup/acp/workflows/CI/badge.svg)](https://github.com/sumup/acp/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/sumup/acp)](./LICENSE)

</div>

Go SDK for the [Agentic Commerce Protocol](https://github.com/agentic-commerce-protocol/agentic-commerce-protocol) (ACP), targeting the latest stable API version, `2026-04-17`.

The SDK provides:

- seller-side HTTP handlers and generated models for Agentic Checkout and Carts;
- provider-side HTTP handlers and models for Delegated Payment and Delegated Authentication;
- a typed client and models for the agent-hosted Feed API;
- a public `/.well-known/acp.json` discovery handler;
- current order-webhook models and Stripe-style HMAC signing; and
- legacy CSV and JSONL product-feed encoders for pull-based integrations.

## Examples

- [`examples/checkout`](examples/checkout) sample checkout provider implementation.
- [`examples/delegated_payment`](examples/delegated_payment) sample PSP (payments service provider) implementation for Delegated Payment.
- [`examples/feed`](examples/feed) sample Product Feed that for exporting feeds in JSONL and CSV formats.

### Checkout Sample

```bash
go run ./examples/checkout
```

Once the server is up, try exercising the flow with `curl`:

```bash
# Create a checkout session with two SKUs
curl -sS -X POST http://localhost:8080/checkout_sessions \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-key' \
  -H 'API-Version: 2026-04-17' \
  -H 'Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000' \
  -d '{
        "line_items": [
          {"id": "latte"},
          {"id": "mug"}
        ],
        "currency": "EUR",
        "buyer": {
          "first_name": "Ava",
          "last_name": "Agent",
          "email": "ava.agent@example.com"
        }
      }'

# Complete the session once you have the id from the response above
curl -sS -X POST http://localhost:8080/checkout_sessions/<session_id>/complete \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-key' \
  -H 'API-Version: 2026-04-17' \
  -H 'Idempotency-Key: 550e8400-e29b-41d4-a716-446655440001' \
  -d '{
        "payment_data": {
          "provider": "sumup",
          "token": "pm_sample_token"
        }
      }'
```

Feel free to copy this sample into your own project and swap the in-memory store for your real product catalog, fulfillment rules, and payment hooks. Use `acpwebhook.NewSender` to deliver full order events with the `Merchant-Signature` format defined by ACP `2026-04-17`.

### Delegated Payment Sample

```bash
go run ./examples/delegated_payment
```

Then call it with:

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
          "currency": "EUR",
          "checkout_session_id": "cs_000001",
          "merchant_id": "demo-merchant",
          "expires_at": "2027-12-31T23:59:59Z"
        },
        "risk_signals": [
          {"type": "card_testing", "action": "manual_review", "score": 10}
        ],
        "metadata": {"source": "sample"}
      }'
```

### Product Feed Sample

```bash
go run ./examples/feed
```

This writes compressed legacy feed exports to `examples/feed/output/product_feed.jsonl.gz` and `examples/feed/output/product_feed.csv.gz`. For the stable push-based Feed API, use `acpfeed.NewClientWithResponses`.

## License

[Apache 2.0](/LICENSE)
