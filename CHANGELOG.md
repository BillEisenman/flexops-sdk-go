# Changelog

All notable changes to the FlexOps Go SDK are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `ScanFormRequest` has `ManifestDate` (`"YYYY-MM-DD"`), `RequestID`, `OriginPostalCode` and `WarehouseID`, the fields
  of the Gateway's SCAN form request. Live creation needs `ManifestDate`: without it the Gateway answers 400
  `Choose the label mailing date.` It stays optional because the sandbox defaults it to today.

### Changed
- `ScanFormRequest.Carrier` is deprecated and sent only when set (`omitempty`): the Gateway has no carrier field on this
  request and ignores it.
- Docs: `Pickups.Schedule`/`Cancel` and `ScanForms.Create` need `shipping:write` with an API key, and the public demo
  keys cannot use them. `Rules.Test` needs `shipping:read`. Shipping rule changes (`Rules.Create`, `Update`, `Delete`,
  `Reorder`) and member changes (`Workspaces.InviteMember`, `RemoveMember`, `UpdateMemberRole`) are for signed-in
  users only: API keys get 403 `SIGNED_IN_USER_REQUIRED`.
- Docs: `BatchLabelItem`, `CreateLabelBatch`, `PreviewLabelBatch`, `ReturnLabelOptions` and
  `GenerateLabelWithOptions` now say that `HazmatClassification` accepts only `HazmatNone` (0) or nil; the Gateway
  refuses any other class with 400 `HazmatNotSupported`. No behavior change.

### Fixed
- `ScanForms.Create` sends a `requestId`: yours, or a UUID generated for the call and sent on every retry of it.
  Before, it sent none, so the Gateway gave each attempt a new ID and a retry after a lost response could create a
  second form. Repeating a request with the same `requestId` returns the same form; reusing it with different
  inputs is 409.
- `ScanForm` now reads what the Gateway returns from `ScanForms.Create`, `List` and `Get`: `ScanFormID`, `State`
  (`ScanFormStateCreated` or `ScanFormStateOutcomeUnknown`), `OriginPostalCode`, `TrackingNumbers`, `CreatedAtUTC`,
  `ScanFormImage` (the base64 PDF, `""` for sandbox forms and unresolved operations; `PDF()` decodes it),
  `ManifestNumber`, `ManifestDate`, `PackageCount`, `StatusMessage`, `Success`, `ErrorMessage` and `IsSandbox`, with
  `Created()`. Before, it declared `ID`, `Carrier`, `FormURL` and `CreatedAt`, none of which the Gateway sends, so
  every field was empty. They stay, deprecated: `ID` and `CreatedAt` are now filled from `scanFormId` and the
  `createdAtUtc` text; `Carrier` and `FormURL` stay empty (SCAN forms are USPS, and the PDF is in `ScanFormImage`).
  No breaking change.
- Docs: `ScanForms.Create` returns a 409 `*FlexOpsError` when the `RequestID` was used with different inputs, when the
  packages already belong to a manifest, or when an operation needs review (its message names the operation). These
  are not retried.

## [3.2.0] - 2026-10-22

3.2.0 requires a maximum total price on multi-package labels: the Gateway refuses to buy a box that would take the shipment above it. It changes only the multi-package call released earlier today.

### Changed
- `MultiPackageShipmentRequest.MaximumPostageAmount` is required: the most the whole shipment may cost, in USD
  (positive, at most two decimals). The Gateway refuses to buy a box that would take the total above it, so a
  shipment never costs more. A box priced above what is left reports `PriceCeilingExceeded`; once the amount is
  used up the remaining boxes report `PriceCeilingReached` and are not bought. Boxes already bought stay bought.
  The field is always sent, so a zero value is refused by the Gateway (400 `ApprovalRequired`). The Gateway
  requires it from the release that ships this change, so this is breaking only for a `CreateMultiPackage` call
  written against 3.1.0, released 2026-10-22.

## [3.1.0] - 2026-10-22

3.1.0 adds multi-package labels and RMA reads, and corrects the returns `get` docs. Existing code keeps working.

### Added
- `Shipping.CreateMultiPackage(ctx, MultiPackageShipmentRequest, idempotencyKey...)` buys labels for one order shipped in 1 to 99 boxes (`POST /api/shipping/multi-package`), with `MultiPackageShipmentRequest`, `MultiPackageDetail`, `MultiPackageShipmentResponse` and `MultiPackageLabelResult`. `HazmatClassification` accepts only `HazmatNone` or nil; the Gateway refuses any other class with 400 `HazmatNotSupported`.
- `Returns.GetRma(ctx, rmaNumber)` and `Returns.ListRmas(ctx, RmaListOptions)` read return authorizations (`GET .../returns/rma/{rmaNumber}` and `GET .../returns/rma`). There was no way to read an RMA back after creating it.
- `FlexOpsError.Data` (`json.RawMessage`) carries the envelope's `data` on an error reply. A multi-package 422 `AllPackagesFailed` or 409 `OutcomeUnknown` puts the per-box results there.

### Fixed
- `Returns.Get` and `Returns.List` say they read received returns: `Get` takes the numeric return ID, and an RMA number matches no Gateway route. Use `GetRma`.

## [3.0.0] - 2026-10-19

The SDK now calls only the documented public Gateway API (the Developer Docs OpenAPI). Calls with no Gateway route and portal-only features are removed; the remaining calls send the Gateway's routes, bodies and numeric enums. Carrier writes through the legacy `/api/ApiProxy` are deprecated and need a user token. A contract test checks every public method against the OpenAPI operation list.

### Removed
- **BREAKING**: `client.ApiKeys` (API keys are managed in the portal) and `client.Offsets` (its per-label routes never existed; use `client.CarbonOffsets`), with `ApiKeyInfo`, `CreateApiKeyRequest`, `CreateApiKeyResponse`, `RegisterRequest` and `CreateWorkspaceRequest`.
- **BREAKING**: `Auth.Register`, `GetProfile`, `UpdateProfile`, `ChangePassword`, `ForgotPassword`, `ResetPassword` and `VerifyEmail` (`Login`, `RefreshToken` and `Logout` stay); `Workspaces.Create` and `Update`; `Wallet.AddFunds`, `ListTransactions` and `ConfigureAutoReload`; `Returns.Approve`, `Reject`, `Cancel` and `ProcessRefund`; `Webhooks.Get`.
- **BREAKING**: proxy writes the Gateway refuses or closes: `Orders.Create` and `Cancel`, `Inventory.PostAsnReceipt` and `UpdateInventory` (proxy writes, refused for API keys); the proxy label purchases the Gateway closes with 410 in every environment: `Carriers.USPS.CreateDomesticLabel`, `CreateReturnLabel` and `CreateInternationalLabel`, `Carriers.UPS.CreateLabel`, `Carriers.DHL.CreateShipment`, `Carriers.FedEx.CreateShipment` and `CreateReturnShipment` (buy labels with `Shipping.CreateLabel`); and `Carriers.DHL.UpdatePickup` (PATCH, which the proxy does not route).
- **BREAKING**: `Shipping.ValidateShippingAddress`, the deprecated `Address`-based `Shipping.ValidateAddress`, `Address` and `AddressValidationResponse`: they called a route the Gateway hides from its API. See `Shipping.ValidateAddress` below.

### Changed
- **BREAKING**: the Go module/import path is now `github.com/BillEisenman/flexops-sdk-go/v3`.
- **BREAKING**: recurring shipments, report schedules and email templates call the Gateway routes; the 2.x paths and shapes never worked. Reads work with any API key or a user token; writes need a user token or an API key granted the new `recurring-shipments:write`, `report-schedules:write` or `email-templates:write` scope (other keys get 403 `API_KEY_SCOPE_REQUIRED`).
  - `RecurringShipments`: `List`, `Get`, `Create`, `Update`, `Delete`, `Pause`, `Resume` and `Trigger` on `recurring-shipments`, with raw `RecurringShipment` results, typed `CreateRecurringShipmentRequest` / `UpdateRecurringShipmentRequest` and `RecurrenceInterval` (Daily 0 ... Custom 4). `Trigger` creates an order now; pass `WithIdempotencyKey` to reconcile it.
  - `client.Reports` is renamed `client.ReportSchedules`: `List`, `Get`, `Create`, `Update`, `Delete` and `Toggle` on `report-schedules`, with enveloped `ReportSchedule` results, `ReportFormat` (Csv 0, JSON 1) and `ReportFrequency` (Daily 0, Weekly 1, Monthly 2).
  - `EmailTemplates` uses `email-templates` (2.x called `shipment-email-templates`): `List`, `Get`, `Upsert` (one template per `TrackingEmailEvent`, LabelCreated 0 ... ReturnToSender 5; it replaces `Create` and `Update`), `Delete` and `Preview(PreviewEmailTemplateRequest)`, with raw `EmailTemplate` and `EmailTemplatePreview` results.
- Carrier writes through the legacy proxy are deprecated (`// Deprecated:`) and need a user token (`Auth.Login`); API keys get 403 `API_KEY_SCOPE_REQUIRED`, and they stop working when the proxy is retired (sunset 2027-06-30). They keep their 2.x names, signatures and routes: `Carriers.USPS.GetDomesticRates`, `GetDomesticProducts`, `GetDomesticPrices`, `GetInternationalRates`, `GetInternationalPrices`, `CancelDomesticLabel`, `CancelInternationalLabel`, `CreatePickup`, `CancelPickup` and `CreateScanForm`; `Carriers.UPS.ValidateAddress`, `GetRates`, `CreatePickup`, `CancelPickup`, `GetTransitTimes`, `GetLandedCost`, `SearchLocations`, `UploadDocument`, `CreateFreightShipment` and `GetFreightRate`; `Carriers.DHL.GetMultiPieceRates`, `CreatePickup`, `CancelPickup`, `CalculateLandedCost`, `ScreenShipment` and `UploadInvoice`; and `Carriers.FedEx` without its label methods. Like any write they now carry an `Idempotency-Key` and follow the write retry rules. Use `Shipping`, `Pickups` and `ScanForms` where they cover the job.
- Retries are safe for purchases. Every POST, PUT and PATCH sends an `Idempotency-Key`: the caller's (the `CreateLabel` argument, or the new `WithIdempotencyKey(ctx, key)`), else one random UUID per call, reused on every retry of that call, so the Gateway replays a write it already finished instead of charging again. Writes retry only on a network error, a timeout, 429, 502, 503, 504 and a 409 `IDEMPOTENCY_CONFLICT`, never on 500 (the Gateway does not cache a 5xx, so a retry would run the write again). A DELETE, which the Gateway does not deduplicate, retries only on 429. A key reused with a different request gets 422 `IDEMPOTENCY_KEY_REUSED`, which is not retried; error codes are matched exactly, so a PascalCase `IdempotencyConflict` 409 (label purchases, managed returns) is not retried either. GETs retry as before. A canceled context is no longer retried.
- A 2xx response with no body (204, or an empty 200) returns the zero value and no error; it failed to decode.
- `Config.MaxRetries` is honored (default 3); it was ignored.

### Added
- `client.CarbonOffsets`: `Quote`, `Purchase`, `List(page, pageSize)`, `Get`, `Summary`, `GetSettings` and `UpdateSettings` (`/api/workspaces/{ws}/carbon-offsets`).
- `Rules.Test` (`POST shipping-rules/test`, `TestShippingRuleRequest`) and `Webhooks.Test` (`POST webhooks/{id}/test`).
- `WithIdempotencyKey(ctx, key)` sets the `Idempotency-Key` of the next write, for purchases you reconcile with a stored key.
- `WebhookEventType` with the event names FlexOps sends and matches subscriptions on: `WebhookEventLabelCreated`, `WebhookEventLabelCanceled`, `WebhookEventBatchCompleted`, and, for USPS push-tracking pilot workspaces, `WebhookEventTrackingUpdated` and `WebhookEventTrackingDelivered`. The README has a subscription sample and the event table.
- `UpdateWebhookRequest` (`URL`, `Description`, `Events`, `Status`) for `Webhooks.Update`.
- `HazmatClassification` (`HazmatNone` = 0 ... `HazmatLithiumBatteryUsedDevice` = 6, the Gateway numbers) and an optional `*HazmatClassification` on `CreateLabelRequest` (also batch items) and `RateRequest`. It is sent only when set; the SDK never defaults it. Labels without an `OrderID` must declare it (`flexops.HazmatNone.Ptr()` for an ordinary parcel).
- `Returns.GenerateLabelWithOptions` sends a declared HAZMAT class with a return label.
- `ReturnLabelOptions` gains `WeightOz`, `LengthIn`, `WidthIn`, `HeightIn`, `CarrierCode`, `ServiceType` and `PaymentConfirmed`. The Gateway needs the package: a workspace with no default shipping preset gets a 400 (`PackageMeasurementsRequired`) naming the missing fields when they are left out (Gateway #672).
- `Shipping.CreateLabelBatch` and `PreviewLabelBatch` take the new `CreateBatchLabelRequest` / `BatchLabelItem`, which match the Gateway batch body. `CreateBatch`, `PreviewBatch` and `BatchLabelRequest` are deprecated (they never matched). `BatchLabelJob` gains the Gateway job fields (`TotalItems`, `CompletedItems`, `FailedItems`, ...); `TrackingInfo` gains the Gateway tracking fields (`CarrierCode`, `EstimatedDeliveryDate`, ...).
- `ShippingAddress` gains optional `Company`, `Phone` and `Email`.

### Fixed
- **BREAKING**: `Shipping.ValidateAddress(ctx, AddressValidationRequest)` posts the documented `/api/addresses/validate` with the Gateway's flat body (`Name`, `Line1`, `Line2`, `City`, `State`, `PostalCode`, `Country`, which defaults to `"US"`) and returns the raw `AddressValidationResult` (`IsValid`, `SuggestedAddress`, `IsResidential`, `Warnings`, `ValidatedBy`, `ValidationUnavailable`, `CheckedAt`).
- **BREAKING**: `Wallet.GetBalance` calls `GET /api/wallet/my-balance` with `X-Workspace-Id` and returns the raw `WalletBalance` (`HasWallet`, `BalanceCents`, `Balance`); it called a route that does not exist.
- **BREAKING**: `HsCodes` uses `/api/shipping/hs-codes/search` (`HsCodeSearchRequest`: `Query`, sent as `query`, `DestinationCountry`, `MaxResults`), `/api/shipping/hs-codes/{code}` and `POST /api/shipping/landed-cost` (`LandedCostRequest`) without a workspace in the path, and returns the raw `[]HsCodeResult`, `HsCodeResult` and `LandedCostEstimate`. `CustomsItem` gains `SkuOrProductID`.
- **BREAKING**: `Insurance.GetQuote` takes `InsuranceQuoteRequest`; `Purchase` posts `PurchaseInsuranceRequest` to `insurance/policies` (it posted to a missing `insurance/purchase`); `FileClaim(ctx, FileClaimRequest)` posts to `insurance/claims` with `PolicyID`, `ClaimReason`, `ClaimAmount` and `EvidenceJSON`; `Void` returns the cancelled policy. `InsuranceQuote` and `InsurancePolicy` match the Gateway DTOs; `InsuranceClaim` is new.
- **BREAKING**: `Returns.Create` posts to `returns/rma` with the Gateway `CreateRmaRequest`: `ReturnRequest` takes a numeric `OrderID` (required) and optional `OrderNumber`, `Notes`, `CustomerEmail` and `ReasonID`, and the result is the new `Rma`. `Returns.List` takes `ReturnListOptions` (`StartDate`, `EndDate`, `Page`, `PageSize`, the Gateway's filters). `ReturnAuthorization` matches the Gateway return (`ReturnID`, `RmaNumber`, `OrderID`, `Status`, `ReturnDate`, ...). `Returns.MarkReceived(ctx, rmaID, ReceiveReturnRequest)` sends `RmaNumber`, `WarehouseID` and `InspectedItems` with a numeric `ReturnDisposition` (Restock 0, Dispose 1, Forward 2). `ReturnLabelOptions` gains `RateIndicator`; `CarrierCode` stays optional (the Gateway defaults it to `"usps"`).
- **BREAKING**: `Rules.Create` and `Update` take `CreateShippingRuleRequest` and `UpdateShippingRuleRequest` (`Name`, `Priority`, `IsEnabled`, `WarehouseID`, `Conditions`, `Actions`); `Rules.Reorder(ctx, []RulePriority)` posts `{rules: [{ruleId, priority}]}` (it sent a bare ID array with PUT). `ShippingRule.IsActive` is now `IsEnabled`; `RuleCondition.Group` is `*int`.
- **BREAKING**: `Pickups.Schedule` takes `SchedulePickupRequest` (`CarrierCode`, `PickupDate`, `Street`, `City`, `State`, `PostalCode`, `Country`, `PackageCount`, ...); `List` reads `pickups/scheduled`; `Get` and `Cancel` take a `carrierCode` (empty means the Gateway default, `"usps"`). `PickupRequest` and `PickupConfirmation` are replaced by `SchedulePickupRequest` and `Pickup`.
- **BREAKING**: workspace roles are the Gateway numbers: `WorkspaceRole` (`WorkspaceRoleOwner` 0, `WorkspaceRoleAdmin` 1, `WorkspaceRoleMember` 2, `WorkspaceRoleMerchant` 3). `InviteMember(ctx, email, role WorkspaceRole)` sends `{email, role}`; `UpdateMemberRole(ctx, userID, role WorkspaceRole)` sends `{newRole}`. `ListMembers` returns `MemberList` (`Members`, `PendingInvites`, `UserLimit`, `ActiveCount`, `CanInviteMore`), not a member array, which never decoded. `Workspaces.List` returns `Workspace` (`MyRole`, `MemberCount`, `Entitlements`, ...) and `Get` returns the new `WorkspaceDetail`.
- **BREAKING**: `LoginResponse` matches the Gateway: the access token is `Token` (it was `AccessToken`, never populated), with `RefreshToken`, `ExpiresAt`, `UserID`, `Email`, `FirstName`, `LastName`, `Roles`, `RequiresTwoFactor`, `TwoFactorToken` and `TwoFactorMethod`. `Auth.Login` makes the client send the returned token. `Auth.RefreshToken(ctx, refreshToken, accessToken...)` sends the `token` the Gateway requires (default: the client's current token; with none it returns an error before any request) and switches the client to the new token.
- **BREAKING**: webhook subscriptions use the event names and fields the Gateway reads and returns (flexops-sdk-node#43, flexops-sdk-python#8). `Webhooks.Create` sent `Events` as an `events` array, which the Gateway ignores, so every subscription received every event. `CreateWebhookRequest.Events` is now `[]WebhookEventType`, sent as the Gateway's comma-separated `eventTypes` with no spaces or duplicates; leave it nil to receive every event. `CreateWebhookRequest` gains `Description`. `Webhooks.Update(ctx, id, data any)` is now `Update(ctx, id, req UpdateWebhookRequest)`. Both return a `*FlexOpsError` with `Code` `"UnknownWebhookEvent"`, before any request, for an empty non-nil `Events` or a name FlexOps never sends (such as `label.created`). `WebhookSubscription` now describes the Gateway response: `EventTypes` (a comma-separated string, or `"*"`), `Status`, `Description`, `ConsecutiveFailures`, `CreatedAt`, `UpdatedAt` and `Secret`; the never-populated `Events` and `IsActive` are removed. Old calls no longer compile: pass `[]flexops.WebhookEventType{...}` instead of `[]string{...}`, pass an `UpdateWebhookRequest` to `Update`, and read `EventTypes`/`Status` instead of `Events`/`IsActive`.
- `Webhooks.Create` returns the new subscription in `Data`, including the `whsec_...` signing secret. The Gateway answers `201` with the bare subscription, not the `ApiResponse` envelope, so `Data` was always empty and the secret, shown only once, was lost. Both the bare and an enveloped body are now read (as in flexops-sdk-dotnet#10).
- **BREAKING**: webhook signature verification now matches what FlexOps sends (flexops-sdk-node#41). `VerifySignature(payload, signature, secret string) bool` hashed the body alone and compared bare hex, so it rejected every real delivery. It is now `VerifySignature(payload []byte, signatureHeader, timestampHeader, secret string) bool`, and the new `VerifyWebhookSignature(payload []byte, signatureHeader, timestampHeader, secret string, opts ...WebhookVerifyOption) error` returns why a delivery failed (`ErrWebhookSignatureMismatch`, `ErrWebhookTimestampOutsideTolerance`, `ErrWebhookInvalidTimestamp`, `ErrWebhookMissingSignature`, `ErrWebhookMissingSecret`). Both sign `{X-FlexOps-Timestamp}.{raw body}` with HMAC-SHA256 keyed by the whole `whsec_...` secret, accept the `sha256=` prefix, compare in constant time, and reject timestamps more than 5 minutes from now (`WithWebhookTolerance`, `WithWebhookClock`). Old calls no longer compile: pass the raw body as `[]byte` and add the `X-FlexOps-Timestamp` header value. Header name constants (`WebhookSignatureHeader`, `WebhookTimestampHeader`, `WebhookEventHeader`, ...) are exported.
- The README webhook sample called a function that did not exist and is now a complete handler. Webhook bodies have no `type`/`data` envelope; the event name is in `X-FlexOps-Event`.
- `Returns.GenerateLabel` sends `{}`; it sent no body, which the Gateway return-label route refuses with 400.
- `Shipping.CreateLabel` posts to the Gateway route `/api/shipping/labels`; it used a workspace-scoped path the Gateway does not serve.
- `Shipping.Track`, `CancelLabel`, `GetCarriers`, `GetRecommendations`, `PredictDelivery` and `GetSavings` use the Gateway `/api/shipping/...` routes instead of workspace-scoped paths the Gateway does not serve, and wrap its raw responses in `ApiResponse` as before. `CancelLabel(ctx, labelID, carrierCode)` takes the carrier code the Gateway requires (variadic, so existing calls compile; without it the call returns an error before any request).
- `DeliveryPredictionRequest` omits an unset `ShipDate` and `ServiceCode`; an empty `shipDate` failed Gateway model binding.
- Zip-only rate requests (deprecated `FromZip`/`ToZip` without full addresses) return a `*FlexOpsError` with `Code` `UnsupportedRequestShape` before any request; the Gateway needs full addresses.

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
