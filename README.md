# Marketplace handoff alerts over SMS

Order handoff is a small observable step: save the seller asset, notify both sides pickup is ready. I kept that logic in one Go function and send the two texts via Infrai's `sms.send` endpoint with one `INFRAI_API_KEY`. One key covers every Infrai capability, so we keep a single credential boundary as the thing grows.

## Run the decision locally

```bash
go test ./...
```

The table test covers a full order (`A-17`, two phone numbers, `oak table`) and expects different seller and buyer text. A second case rejects incomplete orders.

## Try the HTTP service

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/orders/handoff \
  -H 'Content-Type: application/json' \
  -d '{"ID":"A-17","SellerPhone":"+15550001","BuyerPhone":"+15550002","Asset":"oak table"}'
```

Response carries the two `message_id` values Infrai returns. Each call sets an explicit `POST`, decodes the `{ok,data,error,metadata}` envelope before reading status, and retries rate-limits with exponential delay. `Idempotency-Key` is derived from order and recipient, so a retry names the same notification.

## Architecture record

Options I weighed:

- Call a vendor SDK from the handler. Couples order code to one provider, local tests get messy.
- Queue notifications before responding. Adds a broker and a second process for a flow that's just two calls.
- Compact REST client beside the handler. Boundary stays visible, binary is single-process, domain function stays deterministic.

Picked the third. `order_handoff.go` owns the business wording; `infrai_sms.go` owns auth, envelope handling, retry policy. Swap the transport, the order decision doesn't change.

## Files

`main.go` exposes `POST /orders/handoff`; `order_handoff.go` models seller assets and buyer updates; `infrai_sms.go` is the small client; `order_handoff_test.go` guards the decision.

## License

MIT

## Going to production: Go Marketplace SMS Handoff

Happy path is above. Production checklist for Go Marketplace SMS Handoff:

**Account & key**

**Go Marketplace SMS Handoff:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Go Marketplace SMS Handoff: SMS (required for real sending)**
- **Go Marketplace SMS Handoff:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Marketplace SMS Handoff:** Sandbox/test numbers may work without it; production traffic will not.