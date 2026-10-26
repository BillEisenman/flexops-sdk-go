# FlexOps Go SDK

Official Go SDK for the [FlexOps](https://flexops.io) multi-carrier shipping platform. Supports USPS, UPS, FedEx, DHL, OnTrac, Australia Post, Canada Post, Royal Mail, and LSO with rate shopping, label generation, tracking, webhooks, wallet, insurance, returns, and more.

## Installation

```bash
go get github.com/BillEisenman/flexops-sdk-go/v3
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

func main() {
    // API key authentication (recommended for server-to-server)
    client := flexops.NewClient(flexops.Config{
        APIKey:      "live_...",
        WorkspaceID: "your-workspace-id",
    })

    ctx := context.Background()

    // Get shipping rates from all carriers
    rates, err := client.Shipping.GetRates(ctx, flexops.RateRequest{
        Origin: flexops.ShippingAddress{AddressLine1: "123 Main St", City: "New York", StateProvince: "NY", PostalCode: "10001"},
        Destination: flexops.ShippingAddress{AddressLine1: "456 Oak Ave", City: "Los Angeles", StateProvince: "CA", PostalCode: "90210"},
        Package: flexops.ShippingPackage{Weight: 16, WeightUnit: "oz"},
    })
    if err != nil {
        log.Fatal(err)
    }

    request := flexops.CreateLabelRequest{
        CarrierCode: "USPS", ServiceCode: "GROUND_ADVANTAGE",
        Origin: flexops.ShippingAddress{Name: "Warehouse", AddressLine1: "123 Main St", City: "New York", StateProvince: "NY", PostalCode: "10001", CountryCode: "US"},
        Destination: flexops.ShippingAddress{Name: "Customer", AddressLine1: "456 Oak Ave", City: "Los Angeles", StateProvince: "CA", PostalCode: "90210", CountryCode: "US"},
        Package: flexops.ShippingPackage{Weight: 16, WeightUnit: "oz"},
        MaximumPostageAmount: 10.25, // Caller-approved ceiling in USD.
        // Required without an OrderID. The SDK never fills it in; HazmatNone (0) = not hazardous.
        HazmatClassification: flexops.HazmatNone.Ptr(),
    }
    preview, err := client.Shipping.CreateLabel(ctx, request)
    if err != nil { log.Fatal(err) }
    fmt.Println(preview.QuotedPostageAmount, preview.ExpiresAt)

    // Invoke only after explicit caller approval, within five minutes.
    purchaseApprovedLabel := func(purchaseKey string) (flexops.LabelPurchaseResult, error) {
        request.ConfirmationToken = preview.ConfirmationToken
        return client.Shipping.CreateLabel(ctx, request, purchaseKey)
    }
    _ = purchaseApprovedLabel // Connect to your approval flow; persist/reuse the key.

    fmt.Println(rates.Rates)

    // Track a shipment
    info, err := client.Shipping.Track(ctx, "9400111899223456789012")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %s\n", info.Data.Status)
}
```

### Live label approval (2.0.0+)

Migration: `CreateLabel` now returns `LabelPurchaseResult`, with purchased fields
directly on the result instead of under `Data`. Request fields now match Gateway:
`CarrierCode`, `ServiceCode`, `Origin`, `Destination`, and `Package`.

The example below requires SDK 2.0.0 or later. Version 1.0.2 does not
include the per-call idempotency argument.

For live domestic single-label requests, `maximumPostageAmount` is required: positive
USD, at most two decimal places, up to 1,000,000. Missing or invalid values return
400 `ApprovalRequired`. Omitting `confirmationToken` returns a raw 200 preview
(`status`, `quotedPostageAmount`, `maximumPostageAmount`, `currency`, `expiresAt`,
`confirmationToken`). After explicit approval, resubmit the same shipment and ceiling
with the token and a unique per-purchase `Idempotency-Key` header. The token expires
after five minutes. A successful purchase returns the raw 201 label, including
`labelId`, `trackingNumber`, `carrierCode`, and `labelData`, without a `data` wrapper.

Keep the same key and request across retries. An uncertain outcome needs reconciliation;
do not start another purchase with a new key. An expired approval returns 409
`ApprovalExpired`; obtain and explicitly approve a fresh preview when no purchase is
unresolved. The ceiling bounds postage authorization, not later adjustments or separate
fees. Sandbox execution bypasses approval; batch, return, and raw carrier routes have
separate contracts. The SDK never automatically confirms a preview.

### Multi-package labels

One order shipped in several boxes with the same carrier and service; each box gets its own label and a box that
fails does not stop the others. Multi-package labels cannot carry HAZMAT: set `HazmatNone.Ptr()` or leave it nil.
`MaximumPostageAmount` is required: the most the whole shipment may cost in USD (positive, at most two decimals). The
Gateway refuses to buy a box that would take the total above it (`PriceCeilingExceeded`; once the amount is used up
the remaining boxes report `PriceCeilingReached`), and a request without a valid ceiling is 400 `ApprovalRequired`.

```go
shipment, err := client.Shipping.CreateMultiPackage(ctx, flexops.MultiPackageShipmentRequest{
    OrderID: 12345, CarrierCode: "USPS", ServiceCode: "PRIORITY",
    HazmatClassification: flexops.HazmatNone.Ptr(),
    MaximumPostageAmount: 25, // required: the most the whole shipment may cost, in USD
    Packages: []flexops.MultiPackageDetail{
        {WeightInOunces: 16, LengthInInches: 10, WidthInInches: 8, HeightInInches: 4, PackageReference: "Box 1 of 2"},
        {WeightInOunces: 8, LengthInInches: 6, WidthInInches: 6, HeightInInches: 6, PackageReference: "Box 2 of 2"},
    },
}, "order-12345-boxes") // optional Idempotency-Key; reuse it when you retry
var apiErr *flexops.FlexOpsError
if errors.As(err, &apiErr) && (apiErr.Code == "AllPackagesFailed" || apiErr.Code == "OutcomeUnknown") {
    // No label was bought. apiErr.Data holds the per-box results (a MultiPackageShipmentResponse);
    // after OutcomeUnknown a label may exist, so check the order before buying again.
} else if err == nil {
    fmt.Println(shipment.Data.SuccessfulPackages, shipment.Data.FailedPackages)
}
```

### Returns

`Returns.Create` opens an RMA; `GetRma` and `ListRmas` read RMAs by RMA number. `Returns.List` and `Returns.Get`
read *received* returns; `Get` takes the numeric return ID, not an RMA number.

```go
rma, _ := client.Returns.Create(ctx, flexops.ReturnRequest{OrderID: 12345, Reason: "Damaged", Items: []flexops.ReturnItem{{SKU: "SKU-1", Quantity: 1}}})
current, _ := client.Returns.GetRma(ctx, rma.Data.RmaNumber)
```

## Authentication

### API key (recommended)

```go
client := flexops.NewClient(flexops.Config{
    APIKey:      "live_...",
    WorkspaceID: "your-workspace-id",
})
```

### Email / password

```go
client := flexops.NewClient(flexops.Config{BaseURL: "https://gateway.flexops.io"})
login, err := client.Auth.Login(ctx, "user@example.com", "password")
if err != nil {
    log.Fatal(err)
}
// The client now sends login.Data.Token as its bearer token.
client.WorkspaceID = "your-workspace-id"

// Later, exchange the current token and its refresh token for new ones.
_, err = client.Auth.RefreshToken(ctx, login.Data.RefreshToken)
```

## Sandbox / test keys

Use `test_...` (instead of `live_...`) to route to the sandbox environment. Mock carriers respond, nothing hits real carrier APIs, no charges, no real labels. Perfect for CI and integration tests.

```go
client := flexops.NewClient(flexops.Config{
    APIKey:      "test_...",
    WorkspaceID: "your-workspace-id",
})
```

## Carrier-specific calls (legacy proxy)

`client.Carriers.USPS`, `UPS`, `FedEx` and `DHL` call carrier APIs through the deprecated legacy proxy (`/api/ApiProxy`, sunset 2027-06-30).

- **Lookups (GET)** work with API keys: `USPS.CityStateLookup`, `USPS.TrackSummary`, `UPS.Track`, `DHL.FindServicePoints`, ... They take the carrier route's query parameters as a `map[string]string`. The key needs the permission for the data: `shipping:read` for address, ZIP, city, location, service point and reference lookups; `tracking:read` for tracking and proof of delivery; `rates:read` for delivery standards, rates and products. Without it the Gateway returns 403 `API_KEY_SCOPE_REQUIRED`.
- **Writes** (POST/PUT/DELETE, such as `USPS.GetDomesticRates` or `FedEx.CreatePickup`) are deprecated. They need a user token (`client.Auth.Login`); API keys get 403 `API_KEY_SCOPE_REQUIRED`. They stop working when the proxy is retired.
- **Proxy label purchases are closed** (the Gateway answers 410) and are not in the SDK.

Native replacements: rates and labels with `client.Shipping` (`GetRates`, `CreateLabel`, `CancelLabel`), pickups with `client.Pickups`, scan forms with `client.ScanForms`, and landed cost with `client.HsCodes`.

```go
// USPS city/state for a ZIP code (works with an API key that has shipping:read)
place, err := client.Carriers.USPS.CityStateLookup(ctx, "10001")
```

## Retries and idempotency

GETs retry on network errors, timeouts, 429 and 500/502/503/504. Every POST, PUT and PATCH sends an `Idempotency-Key` header, reused on each retry of that call, so the Gateway replays a write it already finished instead of running it (and charging the wallet) again. Writes retry only on network errors, timeouts, 429, 502/503/504 and a 409 `IDEMPOTENCY_CONFLICT` (the first request with that key is still running), never on 500. A DELETE retries only on 429. A key reused with a different request body gets 422 `IDEMPOTENCY_KEY_REUSED`, which is not retried.

The SDK generates one key per call. To reconcile an uncertain purchase later, supply and store your own key:

```go
ctx := flexops.WithIdempotencyKey(context.Background(), "offset-order-1042")
offset, err := client.CarbonOffsets.Purchase(ctx, flexops.CarbonOffsetPurchaseRequest{
    TrackingNumber: "9400111899223456789012", WeightLbs: 2, DistanceMiles: 850,
})
```

Reuse a key only to retry the same request. `Shipping.CreateLabel` takes its purchase key as an argument.

## Webhook subscriptions

Managing webhooks needs a signed-in user's access token. The Gateway refuses API keys on every webhook route, reads included (`List` and `ListDeliveryLogs` as well as `Create`, `Update`, `Delete`, `RotateSecret` and `Test`), with 403 `SIGNED_IN_USER_REQUIRED`, whatever the key's scopes; the SDK returns a `*FlexOpsError` with that `Code`. Reads are refused too because webhook delivery logs carry each delivery's payload (order and customer data). A client with an API key sends the key even when an access token is also set, so use a client created without `APIKey`. Verifying deliveries (below) makes no request and needs no credentials.

```go
// Sign in as a user (or set Config.AccessToken); an API key gets 403 SIGNED_IN_USER_REQUIRED here.
client := flexops.NewClient(flexops.Config{WorkspaceID: "your-workspace-id"})
if _, err := client.Auth.Login(ctx, "user@example.com", "password"); err != nil {
    log.Fatal(err)
}

// Leave Events nil to receive every event.
webhook, err := client.Webhooks.Create(ctx, flexops.CreateWebhookRequest{
    URL: "https://your-app.com/webhooks/flexops",
    Events: []flexops.WebhookEventType{
        flexops.WebhookEventLabelCreated,
        flexops.WebhookEventLabelCanceled,
        flexops.WebhookEventBatchCompleted,
    },
})
if err != nil {
    log.Fatal(err)
}
signingSecret := webhook.Data.Secret // returned only by Create; store it

// Change the events later. Unset fields are left as they are.
_, err = client.Webhooks.Update(ctx, webhook.Data.ID, flexops.UpdateWebhookRequest{
    Events: []flexops.WebhookEventType{flexops.WebhookEventLabelCreated},
})
```

FlexOps sends these events, named exactly as below. Names are PascalCase and matched exactly, so a name like `label.created` would never match. `Create` and `Update` return a `*flexops.FlexOpsError` with `Code` `"UnknownWebhookEvent"`, without calling the API, for an unknown name or an empty (non-nil) `Events`.

| Event | Constant | Sent when |
|---|---|---|
| `LabelCreated` | `WebhookEventLabelCreated` | A live label is purchased (sandbox labels send nothing) |
| `LabelCanceled` | `WebhookEventLabelCanceled` | A live label is voided |
| `BatchCompleted` | `WebhookEventBatchCompleted` | A batch label job finishes |
| `TrackingUpdated` | `WebhookEventTrackingUpdated` | A USPS tracking scan arrives (USPS push-tracking pilot workspaces only) |
| `TrackingDelivered` | `WebhookEventTrackingDelivered` | A USPS package is delivered (USPS push-tracking pilot workspaces only) |

A subscription's `EventTypes` is the comma-separated list it receives, or `"*"` for every event.

## Webhook verification

FlexOps signs each delivery with HMAC-SHA256 over `{X-FlexOps-Timestamp}.{raw body}`, keyed with your subscription secret (`whsec_...`, used as-is), and sends it as `X-FlexOps-Signature: sha256=<hex>`. Verify the raw body before you decode it.

```go
import (
    "encoding/json"
    "io"
    "net/http"
    "os"

    flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    err = flexops.VerifyWebhookSignature(
        body,
        r.Header.Get(flexops.WebhookSignatureHeader), // "sha256=..."
        r.Header.Get(flexops.WebhookTimestampHeader), // Unix seconds
        os.Getenv("FLEXOPS_WEBHOOK_SECRET"),          // "whsec_..."
    )
    if err != nil {
        // errors.Is(err, flexops.ErrWebhookTimestampOutsideTolerance), ErrWebhookSignatureMismatch, ...
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    // The body is the event data itself (no type/data envelope).
    // The event name is in a header; test events send "test".
    switch r.Header.Get(flexops.WebhookEventHeader) {
    case "test":
        // body: {"type":"test","message":"...","timestamp":"...","subscriptionId":"..."}
    case string(flexops.WebhookEventTrackingUpdated):
        var data map[string]any
        _ = json.Unmarshal(body, &data)
        // handle data["trackingNumber"], data["status"], ...
    }
    w.WriteHeader(http.StatusOK)
}
```

Deliveries whose timestamp is more than 5 minutes from your clock (either way) are rejected to limit replays; pass `flexops.WithWebhookTolerance(d)` to change that. `flexops.VerifySignature(body, signature, timestamp, secret)` returns a `bool` with the default tolerance.

## Curl quickstart

Every SDK method is a thin wrapper around the FlexOps REST API. If you want to verify the API before committing to the SDK — or you're integrating from a language we don't ship a SDK for — these curl invocations hit the same endpoints:

```bash
# Shop rates across all connected carriers
curl -X POST https://gateway.flexops.io/api/shipping/rates \
  -H "X-API-Key: live_..." \
  -H "Content-Type: application/json" \
  -d '{
    "origin": {"addressLine1": "123 Main St", "city": "New York", "stateProvince": "NY", "postalCode": "10001", "countryCode": "US"},
    "destination": {"addressLine1": "456 Oak Ave", "city": "Los Angeles", "stateProvince": "CA", "postalCode": "90210", "countryCode": "US"},
    "package": {"weight": 16, "weightUnit": "oz"}
  }'

# Preview a live label; this does not purchase it.
# hazmatClassification 0 = not hazardous; required when there is no orderId.
curl -X POST https://gateway.flexops.io/api/shipping/labels \
  -H "X-API-Key: live_..." \
  -H "Content-Type: application/json" \
  -d '{
    "carrierCode": "USPS", "serviceCode": "GROUND_ADVANTAGE",
    "origin": {"name": "Warehouse", "addressLine1": "123 Main St", "city": "New York", "stateProvince": "NY", "postalCode": "10001", "countryCode": "US"},
    "destination": {"name": "Customer", "addressLine1": "456 Oak Ave", "city": "Los Angeles", "stateProvince": "CA", "postalCode": "90210", "countryCode": "US"},
    "package": {"weight": 16, "weightUnit": "oz"},
    "maximumPostageAmount": 10.25,
    "hazmatClassification": 0
  }'

# After explicit approval, resend the same body with confirmationToken from
# the preview and an Idempotency-Key header unique to this purchase.
# See the SDK example above; retain that key for retries.

# Track a shipment
curl https://gateway.flexops.io/api/shipping/track/9400111899223456789012 \
  -H "X-API-Key: live_..."

# Cancel (void) a label; the carrier code is required
curl -X DELETE "https://gateway.flexops.io/api/shipping/labels/{labelId}?carrierCode=USPS" \
  -H "X-API-Key: live_..."
```

Use a `test_...` key instead of `live_...` to hit the sandbox environment; mock carriers respond, no real charges, no real labels.

## Services

| Service | Description |
|---------|-------------|
| `client.Auth` | Login, token refresh, logout (user tokens) |
| `client.Workspaces` | Workspaces and membership (invite, roles, remove) |
| `client.Shipping` | Rate shopping, labels, multi-package labels, tracking, batch, cancel, address validation |
| `client.Carriers` | USPS, UPS, FedEx, DHL through the legacy proxy (reads; deprecated user-token writes) |
| `client.Webhooks` | Subscription CRUD, test deliveries, signature verification, delivery logs |
| `client.Wallet` | Workspace wallet balance |
| `client.Insurance` | Quotes, policies, claims |
| `client.Returns` | RMAs (create, read, list), return labels, receiving, received returns |
| `client.Analytics` | Shipments, orders, carrier performance (API keys need each method's read permission) |
| `client.Orders` | Order reads (API keys need `orders:read`) |
| `client.Inventory` | Warehouse inventory reads (API keys need `inventory:read`) |
| `client.Pickups` | Carrier pickup scheduling |
| `client.ScanForms` | USPS scan forms |
| `client.Rules` | Shipping automation rules, including dry-run tests |
| `client.CarbonOffsets` | Carbon offset quotes, purchases and settings |
| `client.HsCodes` | HS code lookup and landed cost for international customs |
| `client.RecurringShipments` | Shipments that repeat on a schedule (pause, resume, run now) |
| `client.ReportSchedules` | Reports emailed on a schedule |
| `client.EmailTemplates` | Tracking email templates per event, with preview |

`Analytics`, `Orders` and `Inventory` read through the legacy proxy. With an API key they need the permission for the data they read: `orders:read` for `Orders`, `inventory:read` for `Inventory`, and for `Analytics` the permission named on each method (`shipping:read`, `inventory:read`, `orders:read` or `returns:read`). A key without it gets 403 `API_KEY_SCOPE_REQUIRED`. Keys created with the default permissions have all of these.

`RecurringShipments`, `ReportSchedules` and `EmailTemplates` reads work with any API key. Their writes (including email-template preview) need a user token or an API key granted `recurring-shipments:write`, `report-schedules:write` or `email-templates:write`; other keys get 403 `API_KEY_SCOPE_REQUIRED`. API keys themselves are managed in the FlexOps portal.

## Configuration

```go
client := flexops.NewClient(flexops.Config{
    BaseURL:     "https://gateway.flexops.io", // API base URL
    APIKey:      "live_...",               // API key auth
    WorkspaceID: "your-workspace-id",      // Default workspace
    Timeout:     30 * time.Second,         // Request timeout
    MaxRetries:  3,                        // Retries per call (see Retries and idempotency)
})
```

## Requirements

- Go 1.22+

## License

MIT © FlexOps, LLC. See [LICENSE](LICENSE) for full text.
