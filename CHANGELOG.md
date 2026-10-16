# Changelog

All notable changes to the FlexOps Go SDK are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `WebhookEventType` with the event names FlexOps sends and matches subscriptions on: `WebhookEventLabelCreated`, `WebhookEventLabelCanceled`, `WebhookEventBatchCompleted`, and, for USPS push-tracking pilot workspaces, `WebhookEventTrackingUpdated` and `WebhookEventTrackingDelivered`. The README has a subscription sample and the event table.
- `UpdateWebhookRequest` (`URL`, `Description`, `Events`, `Status`) for `Webhooks.Update`.
- `HazmatClassification` (`HazmatNone` = 0 ... `HazmatLithiumBatteryUsedDevice` = 6, the Gateway numbers) and an optional `*HazmatClassification` on `CreateLabelRequest` (also batch items) and `RateRequest`. It is sent only when set; the SDK never defaults it. Labels without an `OrderID` must declare it (`flexops.HazmatNone.Ptr()` for an ordinary parcel).
- `Returns.GenerateLabelWithOptions` sends a declared HAZMAT class with a return label.
- `ReturnLabelOptions` gains `WeightOz`, `LengthIn`, `WidthIn`, `HeightIn`, `CarrierCode`, `ServiceType` and `PaymentConfirmed`. The Gateway needs the package: a workspace with no default shipping preset gets a 400 (`PackageMeasurementsRequired`) naming the missing fields when they are left out (Gateway #672).
- `Shipping.CreateLabelBatch` and `PreviewLabelBatch` take the new `CreateBatchLabelRequest` / `BatchLabelItem`, which match the Gateway batch body. `CreateBatch`, `PreviewBatch` and `BatchLabelRequest` are deprecated (they never matched). `BatchLabelJob` gains the Gateway job fields (`TotalItems`, `CompletedItems`, `FailedItems`, ...); `TrackingInfo` gains the Gateway tracking fields (`CarrierCode`, `EstimatedDeliveryDate`, ...).
- `Shipping.ValidateShippingAddress` returns the Gateway `AddressValidationResponse`. `ShippingAddress` gains optional `Company`, `Phone` and `Email`.

### Fixed
- **BREAKING**: webhook subscriptions use the event names and fields the Gateway reads and returns (flexops-sdk-node#43, flexops-sdk-python#8). `Webhooks.Create` sent `Events` as an `events` array, which the Gateway ignores, so every subscription received every event. `CreateWebhookRequest.Events` is now `[]WebhookEventType`, sent as the Gateway's comma-separated `eventTypes` with no spaces or duplicates; leave it nil to receive every event. `CreateWebhookRequest` gains `Description`. `Webhooks.Update(ctx, id, data any)` is now `Update(ctx, id, req UpdateWebhookRequest)`. Both return a `*FlexOpsError` with `Code` `"UnknownWebhookEvent"`, before any request, for an empty non-nil `Events` or a name FlexOps never sends (such as `label.created`). `WebhookSubscription` now describes the Gateway response: `EventTypes` (a comma-separated string, or `"*"`), `Status`, `Description`, `ConsecutiveFailures`, `CreatedAt`, `UpdatedAt` and `Secret`; the never-populated `Events` and `IsActive` are removed. Old calls no longer compile: pass `[]flexops.WebhookEventType{...}` instead of `[]string{...}`, pass an `UpdateWebhookRequest` to `Update`, and read `EventTypes`/`Status` instead of `Events`/`IsActive`.
- `Webhooks.Create` returns the new subscription in `Data`, including the `whsec_...` signing secret. The Gateway sends it without an envelope, so `Data` was always empty and the secret, shown only once, was lost.
- **BREAKING**: webhook signature verification now matches what FlexOps sends (flexops-sdk-node#41). `VerifySignature(payload, signature, secret string) bool` hashed the body alone and compared bare hex, so it rejected every real delivery. It is now `VerifySignature(payload []byte, signatureHeader, timestampHeader, secret string) bool`, and the new `VerifyWebhookSignature(payload []byte, signatureHeader, timestampHeader, secret string, opts ...WebhookVerifyOption) error` returns why a delivery failed (`ErrWebhookSignatureMismatch`, `ErrWebhookTimestampOutsideTolerance`, `ErrWebhookInvalidTimestamp`, `ErrWebhookMissingSignature`, `ErrWebhookMissingSecret`). Both sign `{X-FlexOps-Timestamp}.{raw body}` with HMAC-SHA256 keyed by the whole `whsec_...` secret, accept the `sha256=` prefix, compare in constant time, and reject timestamps more than 5 minutes from now (`WithWebhookTolerance`, `WithWebhookClock`). Old calls no longer compile: pass the raw body as `[]byte` and add the `X-FlexOps-Timestamp` header value. Header name constants (`WebhookSignatureHeader`, `WebhookTimestampHeader`, `WebhookEventHeader`, ...) are exported.
- The README webhook sample called a function that did not exist and is now a complete handler. Webhook bodies have no `type`/`data` envelope; the event name is in `X-FlexOps-Event`.
- `Returns.GenerateLabel` sends `{}`; it sent no body, which the Gateway return-label route refuses with 400.
- `Shipping.CreateLabel` posts to the Gateway route `/api/shipping/labels`; it used a workspace-scoped path the Gateway does not serve.
- `Shipping.Track`, `CancelLabel`, `GetCarriers`, `GetRecommendations`, `PredictDelivery` and `GetSavings` use the Gateway `/api/shipping/...` routes instead of workspace-scoped paths the Gateway does not serve, and wrap its raw responses in `ApiResponse` as before. `CancelLabel(ctx, labelID, carrierCode)` takes the carrier code the Gateway requires (variadic, so existing calls compile; without it the call returns an error before any request).
- `DeliveryPredictionRequest` omits an unset `ShipDate` and `ServiceCode`; an empty `shipDate` failed Gateway model binding.
- Zip-only rate requests (deprecated `FromZip`/`ToZip` without full addresses) return a `*FlexOpsError` with `Code` `UnsupportedRequestShape` before any request; the Gateway needs full addresses.
- `Shipping.ValidateAddress` (now deprecated) posts the Gateway `{"address": {...}}` request to `/api/shipping/addresses/validate`, mapping the legacy `Address` to the Gateway field names and the response back.

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
