# FlexOps Go SDK

Official Go SDK for the [FlexOps](https://flexops.io) multi-carrier shipping platform. Supports USPS, UPS, FedEx, DHL, OnTrac, Australia Post, Canada Post, Royal Mail, and LSO with rate shopping, label generation, tracking, webhooks, wallet, insurance, returns, and more.

## Installation

```bash
go get github.com/BillEisenman/flexops-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    flexops "github.com/BillEisenman/flexops-sdk-go"
)

func main() {
    // API key authentication (recommended for server-to-server)
    client := flexops.NewClient(flexops.Config{
        APIKey:      "fxk_live_...",
        WorkspaceID: "ws_abc123",
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

### Live label approval (unreleased SDK changes)

Migration: `CreateLabel` now returns `LabelPurchaseResult`, with purchased fields
directly on the result instead of under `Data`. Request fields now match Gateway:
`CarrierCode`, `ServiceCode`, `Origin`, `Destination`, and `Package`.

The example below requires this source revision; the published 1.0.2 packages do
not include the new per-call idempotency argument. Release these SDK changes before
using that argument from a package registry.

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

## Authentication

### API key (recommended)

```go
client := flexops.NewClient(flexops.Config{
    APIKey:      "fxk_live_...",
    WorkspaceID: "ws_abc123",
})
```

### Email / password

```go
client := flexops.NewClient(flexops.Config{BaseURL: "https://gateway.flexops.io"})
if err := client.Auth.Login(ctx, "user@example.com", "password"); err != nil {
    log.Fatal(err)
}
client.WorkspaceID = "ws_abc123"
```

## Sandbox / test keys

Use `fxk_test_...` (instead of `fxk_live_...`) to route to the sandbox environment. Mock carriers respond, nothing hits real carrier APIs, no charges, no real labels. Perfect for CI and integration tests.

```go
client := flexops.NewClient(flexops.Config{
    APIKey:      "fxk_test_...",
    WorkspaceID: "ws_abc123",
})
```

## Direct carrier operations

Access carrier-specific endpoints when you need full control. The carrier-specific services are typed:

```go
// USPS domestic label
label, err := client.Carriers.USPS.CreateDomesticLabel(ctx, flexops.UspsLabelRequest{
    ImageType:      "PDF",
    MailClass:      "PRIORITY_MAIL",
    WeightInOunces: 16,
})

// FedEx rate quote
rates, err := client.Carriers.FedEx.GetRates(ctx, flexops.FedExRateRequest{...})

// UPS tracking
info, err := client.Carriers.UPS.Track(ctx, "1Z999AA10123456784")

// DHL shipment
shipment, err := client.Carriers.DHL.CreateShipment(ctx, flexops.DhlShipmentRequest{...})
```

## Webhook verification

```go
import flexops "github.com/BillEisenman/flexops-sdk-go"

valid := flexops.VerifyWebhookSignature(
    payload,   // []byte of the raw request body
    signature, // value of the X-FlexOps-Signature header
    "whsec_...",
)
```

## Curl quickstart

Every SDK method is a thin wrapper around the FlexOps REST API. If you want to verify the API before committing to the SDK — or you're integrating from a language we don't ship a SDK for — these curl invocations hit the same endpoints:

```bash
# Shop rates across all connected carriers
curl -X POST https://gateway.flexops.io/api/shipping/rates \
  -H "X-API-Key: fxk_live_..." \
  -H "Content-Type: application/json" \
  -d '{
    "origin": {"addressLine1": "123 Main St", "city": "New York", "stateProvince": "NY", "postalCode": "10001", "countryCode": "US"},
    "destination": {"addressLine1": "456 Oak Ave", "city": "Los Angeles", "stateProvince": "CA", "postalCode": "90210", "countryCode": "US"},
    "package": {"weight": 16, "weightUnit": "oz"}
  }'

# Preview a live label; this does not purchase it.
curl -X POST https://gateway.flexops.io/api/workspaces/ws_abc123/shipping/labels \
  -H "X-API-Key: fxk_live_..." \
  -H "Content-Type: application/json" \
  -d '{
    "carrierCode": "USPS", "serviceCode": "GROUND_ADVANTAGE",
    "origin": {"name": "Warehouse", "addressLine1": "123 Main St", "city": "New York", "stateProvince": "NY", "postalCode": "10001", "countryCode": "US"},
    "destination": {"name": "Customer", "addressLine1": "456 Oak Ave", "city": "Los Angeles", "stateProvince": "CA", "postalCode": "90210", "countryCode": "US"},
    "package": {"weight": 16, "weightUnit": "oz"},
    "maximumPostageAmount": 10.25
  }'

# After explicit approval, resend the same body with confirmationToken from
# the preview and an Idempotency-Key header unique to this purchase.
# See the SDK example above; retain that key for retries.

# Track a shipment
curl https://gateway.flexops.io/api/workspaces/ws_abc123/shipping/track/9400111899223456789012 \
  -H "X-API-Key: fxk_live_..."

# Cancel a label (via the unified carrier-agnostic endpoint)
curl -X DELETE https://gateway.flexops.io/api/v3.0/shipping/Usps/cancel/9400111899223456789012 \
  -H "X-API-Key: fxk_live_..."
```

Use an `fxk_test_...` key instead of `fxk_live_...` to hit the sandbox environment; mock carriers respond, no real charges, no real labels.

## Services

| Service | Description |
|---------|-------------|
| `client.Auth` | Login, register, password management |
| `client.Workspaces` | Workspace CRUD, membership, branding |
| `client.Shipping` | Rate shopping, labels, tracking, batch, cancel |
| `client.Carriers` | USPS, UPS, FedEx, DHL direct endpoints |
| `client.Webhooks` | Subscription CRUD, signature verification, delivery logs |
| `client.Wallet` | Balance, transactions, auto-reload |
| `client.Insurance` | Quotes, purchase, claims (first-party + U-PIC) |
| `client.Returns` | RMA lifecycle: create, batch, QR codes, photo upload, cost recovery |
| `client.ApiKeys` | Key creation, rotation, revocation |
| `client.Analytics` | Shipments, orders, carrier performance |
| `client.Orders` | Order management |
| `client.Inventory` | Warehouse inventory |
| `client.Pickups` | Carrier pickup scheduling |
| `client.ScanForms` | USPS scan forms |
| `client.Rules` | Shipping automation rules |
| `client.Offsets` | Carbon offset purchases |
| `client.HsCodes` | HS code lookup for international customs |
| `client.RecurringShipments` | Scheduled recurring shipments |
| `client.EmailTemplates` | Branded post-purchase email templates |
| `client.Reports` | Report generation and scheduled delivery |

## Configuration

```go
client := flexops.NewClient(flexops.Config{
    BaseURL:     "https://gateway.flexops.io", // API base URL
    APIKey:      "fxk_live_...",           // API key auth
    WorkspaceID: "ws_abc123",              // Default workspace
    Timeout:     30 * time.Second,         // Request timeout
    MaxRetries:  3,                        // Retry on transient failures
})
```

## Requirements

- Go 1.22+

## License

MIT © FlexOps, LLC. See [LICENSE](LICENSE) for full text.
