<div align="center">

# Go SDK for the Agentic Commerce Protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/sumup/acp.svg)](https://pkg.go.dev/github.com/sumup/acp)
[![CI Status](https://github.com/sumup/acp/workflows/CI/badge.svg)](https://github.com/sumup/acp/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/sumup/acp)](./LICENSE)

</div>

This repo bootstraps Go support for the [Agentic Commerce Protocol](https://developers.openai.com/commerce) (ACP). Both the checkout and delegated payment SDKs now live under a single Go module so you can depend on `github.com/sumup/acp` and call the handlers, models, and helpers you need.

## Features

- **Checkout API** — plug your own business logic into `NewCheckoutHandler` by implementing `CheckoutSessionService`. The handler exposes the official ACP checkout contract over `net/http`, supports optional signature verification and timestamp skew enforcement, and emits typed responses generated from the OpenAPI spec.
- **Delegated Payment API** — payment service providers implement `DelegatedPaymentProvider` and wire it up via `NewDelegatedPaymentHandler` (optionally adding `DelegatedPaymentWithAuthenticator` and signature enforcement) to tokenize credentials and emit delegated vault tokens.

## Example Servers

Two runnable samples live under [`examples`](examples):

- [`examples/checkout`](examples/checkout) implements `CheckoutSessionService` with an in-memory catalog and session store.
- [`examples/delegated_payment`](examples/delegated_payment) implements `DelegatedPaymentProvider` with an in-memory vault token map.

### Checkout sample

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

### Delegated payment sample

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

## License

[Apache 2.0](/LICENSE)
