<div align="center">

# Go SDK for the Agentic Commerce Protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/sumup/acp.svg)](https://pkg.go.dev/github.com/sumup/acp)
[![CI Status](https://github.com/sumup/acp/workflows/CI/badge.svg)](https://github.com/sumup/acp/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/sumup/acp)](./LICENSE)

</div>

Go SDK for the [Agentic Commerce Protocol](https://developers.openai.com/commerce) (ACP). `github.com/sumup/acr` supports [Agentic Checkout](https://developers.openai.com/commerce/specs/checkout), [Delegated Payment](https://developers.openai.com/commerce/specs/payment), and [Product Feeds](https://developers.openai.com/commerce/specs/feed).

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
  -d '{
        "items": [
          {"id": "latte", "quantity": 1},
          {"id": "mug", "quantity": 1}
        ],
        "buyer": {
          "first_name": "Ava",
          "last_name": "Agent",
          "email": "ava.agent@example.com"
        }
      }'

# Complete the session once you have the id from the response above
curl -sS -X POST http://localhost:8080/checkout_sessions/<session_id>/complete \
  -H 'Content-Type: application/json' \
  -d '{
        "payment_data": {
          "provider": "stripe",
          "token": "pm_sample_token"
        }
      }'

Feel free to copy this sample into your own project and swap the in-memory store for your real product catalog, fulfillment rules, and payment hooks.

To see webhook delivery end-to-end, export the environment variables below before starting the sample server. The handler will POST an `order_created` event every time a checkout session completes.

```bash
export ACP_WEBHOOK_ENDPOINT="https://webhook.site/your-endpoint"
export ACP_WEBHOOK_HEADER="Merchant_Name-Signature"
export ACP_WEBHOOK_SECRET="super-secret"
go run ./examples/checkout
```

### Delegated Payment Sample

```bash
go run ./examples/delegated_payment
```

Then call it with:

```bash
curl -sS -X POST http://localhost:8080/agentic_commerce/delegate_payment \
  -H 'Content-Type: application/json' \
  -d '{
        "payment_method": {
          "type": "card",
          "card_number_type": "fpan",
          "number": "4242424242424242",
          "exp_month": "11",
          "exp_year": "2026",
          "display_last4": "4242",
          "display_card_funding_type": "credit",
          "metadata": {"issuer": "demo-bank"}
        },
        "allowance": {
          "reason": "one_time",
          "max_amount": 2000,
          "currency": "usd",
          "checkout_session_id": "cs_000001",
          "merchant_id": "demo-merchant",
          "expires_at": "2025-12-31T23:59:59Z"
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

This writes compressed feed exports to `examples/feed/output/product_feed.jsonl.gz` and `examples/feed/output/product_feed.csv.gz`.

## License

[Apache 2.0](/LICENSE)
