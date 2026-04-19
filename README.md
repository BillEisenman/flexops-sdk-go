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
        FromAddress: flexops.Address{Street1: "123 Main St", City: "New York",    State: "NY", Zip: "10001", Country: "US"},
        ToAddress:   flexops.Address{Street1: "456 Oak Ave", City: "Los Angeles", State: "CA", Zip: "90210", Country: "US"},
        Parcel:      flexops.Parcel{Weight: 16, WeightUnit: "oz"},
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a label with the cheapest rate
    cheapest := rates.Data[0] // rates are returned sorted by total cost
    label, err := client.Shipping.CreateLabel(ctx, flexops.CreateLabelRequest{
        Carrier:     cheapest.Carrier,
        Service:     cheapest.Service,
        FromAddress: flexops.Address{Name: "Warehouse", Street1: "123 Main St", City: "New York",    State: "NY", Zip: "10001", Country: "US"},
        ToAddress:   flexops.Address{Name: "Customer",  Street1: "456 Oak Ave", City: "Los Angeles", State: "CA", Zip: "90210", Country: "US"},
        Parcel:      flexops.Parcel{Weight: 16, WeightUnit: "oz"},
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Label URL: %s\n", label.Data.LabelURL)
    fmt.Printf("Tracking:  %s\n", label.Data.TrackingNumber)

    // Track a shipment
    info, err := client.Shipping.Track(ctx, "9400111899223456789012")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %s\n", info.Data.Status)
}
```

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
curl -X POST https://gateway.flexops.io/api/workspaces/ws_abc123/shipping/rates \
  -H "X-API-Key: fxk_live_..." \
  -H "Content-Type: application/json" \
  -d '{
    "fromAddress": {"street1": "123 Main St", "city": "New York", "state": "NY", "zip": "10001", "country": "US"},
    "toAddress":   {"street1": "456 Oak Ave", "city": "Los Angeles", "state": "CA", "zip": "90210", "country": "US"},
    "parcel":      {"weight": 16, "weightUnit": "oz"}
  }'

# Create a label
curl -X POST https://gateway.flexops.io/api/workspaces/ws_abc123/shipping/labels \
  -H "X-API-Key: fxk_live_..." \
  -H "Content-Type: application/json" \
  -d '{
    "carrier":  "USPS",
    "service":  "PRIORITY_MAIL",
    "fromAddress": {"name": "Warehouse", "street1": "123 Main St", "city": "New York", "state": "NY", "zip": "10001", "country": "US"},
    "toAddress":   {"name": "Customer",  "street1": "456 Oak Ave", "city": "Los Angeles", "state": "CA", "zip": "90210", "country": "US"},
    "parcel":   {"weight": 16, "weightUnit": "oz"}
  }'

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

Proprietary — FlexOps, LLC
