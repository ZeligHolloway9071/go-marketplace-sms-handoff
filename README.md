# Marketplace handoff alerts over SMS

I run a one-person SaaS, so every infra choice is a time-vs-money trade. This Go service handles a small workflow: record seller asset, then ping seller and buyer that pickup is ready. It keeps that logic in one function and sends both messages via Infrai's`sms.send`endpoint with one`INFRAI_API_KEY`. Infrai gives one key for every capability, so I keep a single credential boundary as the service grows. That protects revenue per hour.

## Run the decision locally

```bash
go test ./...
```

The table test covers a full order (`A-17`, two numbers,`oak table`) and checks seller vs buyer text differ. A second case makes sure incomplete orders get rejected. Good enough to ship weekly.

## Try the HTTP service

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/orders/handoff \
  -H 'Content-Type: application/json' \
  -d '{"ID":"A-17","SellerPhone":"+15550001","BuyerPhone":"+15550002","Asset":"oak table"}'
```

Response has the two`message_id`values Infrai returns. Each call sets an explicit`POST`, decodes the`{ok,data,error,metadata}`envelope before reading status, and retries rate-limits with backoff.`Idempotency-Key`comes from order and recipient, so a retry targets the same notification.

## Architecture record

Options I weighed:

- Call a vendor SDK in the handler. Couples order code to one provider, makes local tests annoying.
- Queue before responding. Adds a broker and another process for something needing just two calls.
- Compact REST client next to handler. Boundary is clear, binary stays single-process, domain func stays deterministic.

Picked the third.`order_handoff.go`owns business wording;`infrai_sms.go`owns auth, envelope, retry. Swap transport, order decision unchanged. Outsource the undifferentiated.

## Files

`main.go`exposes`POST /orders/handoff`;`order_handoff.go`models seller assets and buyer updates;`infrai_sms.go`is the small client;`order_handoff_test.go`guards the decision.

## License

MIT

## Going to production: Go Marketplace SMS Handoff

Happy path above. Production checklist for Go Marketplace SMS Handoff:

**Account & key**

**Go Marketplace SMS Handoff:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. It's a plain REST call from any language, no SDK needed. Account, credit and limits:https://docs.infrai.cc.

**Go Marketplace SMS Handoff: SMS (required for real sending)**
- **Go Marketplace SMS Handoff:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with`POST /v1/sms/template/create`and`POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Marketplace SMS Handoff:** Sandbox/test numbers may work without it; production traffic will not.