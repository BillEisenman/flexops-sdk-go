# Changelog

All notable changes to the FlexOps Go SDK are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2026-09-19

### Bounded label purchase support

- **BREAKING**: the Go module/import path is now `github.com/BillEisenman/flexops-sdk-go/v2`.

- Single-label requests support maximumPostageAmount and confirmationToken.
- CreateLabel/create_label accepts a per-purchase idempotency key and preserves it on retries.
- Preview and purchase responses are raw objects, not success/data envelopes.
- Gateway errorCode is preserved on SDK errors.
- Read the README approval flow and migration notes before upgrading.
- Request models now use carrierCode/serviceCode/origin/destination/package and existing ShippingAddress/ShippingPackage types. Migrate the former carrier/service/fromAddress/toAddress/parcel fields. BatchLabelRequest also shares this corrected request model; approval fields remain optional.
- Label uses carrierCode and currency; service is not part of the Gateway label response.
- CreateLabel returns LabelPurchaseResult. Check Status == "Preview" before treating the result as a purchase; access purchased fields directly (LabelID, TrackingNumber), not through Data.

### Added
- Initial README with installation, quick start, authentication (API key and email/password), sandbox guidance, direct carrier operations, webhook verification, and a curl quickstart section.

### Changed
- **BREAKING**: module path changed from `github.com/FlexOps/flexops-sdk-go` to `github.com/BillEisenman/flexops-sdk-go`. Update your imports before upgrading.

## [1.0.0] - 2026-03-08

### Added
- Initial public release.
- `flexops.NewClient(flexops.Config{...})` entry point with API key and JWT authentication.
- 20 service pointers on the `Client` struct covering Auth, Workspaces, Shipping, Carriers, Webhooks, Wallet, Insurance, Returns, ApiKeys, Analytics, Orders, Inventory, Pickups, ScanForms, Rules, Offsets, HsCodes, RecurringShipments, EmailTemplates, and Reports.
- Direct carrier access for USPS, UPS, FedEx, and DHL with typed request/response structs.
- Context-based cancellation on every service method.
- Requires Go 1.22+.
