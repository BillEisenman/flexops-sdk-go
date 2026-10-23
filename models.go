// ***********************************************************************
// Package          : flexops-sdk-go
// Author           : FlexOps, LLC
// Created          : 2026-03-08
//
// Copyright (c) 2021-2026 by FlexOps, LLC. All rights reserved.
// ***********************************************************************

package flexops

// ApiResponse is the standard API response wrapper.
type ApiResponse[T any] struct {
	Success bool     `json:"success"`
	Data    T        `json:"data"`
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// PaginatedResponse wraps paginated results.
type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"totalCount"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the Gateway login and refresh-token result.
type LoginResponse struct {
	Token             string   `json:"token"`
	RefreshToken      string   `json:"refreshToken"`
	ExpiresAt         string   `json:"expiresAt"`
	UserID            string   `json:"userId"`
	Email             string   `json:"email"`
	FirstName         string   `json:"firstName"`
	LastName          string   `json:"lastName"`
	Roles             []string `json:"roles"`
	RequiresTwoFactor bool     `json:"requiresTwoFactor"`
	TwoFactorToken    string   `json:"twoFactorToken,omitempty"`
	TwoFactorMethod   string   `json:"twoFactorMethod,omitempty"`
}

// WorkspaceRole is a member's role, sent and returned as the Gateway number.
type WorkspaceRole int

const (
	WorkspaceRoleOwner    WorkspaceRole = 0
	WorkspaceRoleAdmin    WorkspaceRole = 1
	WorkspaceRoleMember   WorkspaceRole = 2
	WorkspaceRoleMerchant WorkspaceRole = 3
)

// WorkspaceEntitlement is a product plan the workspace holds.
type WorkspaceEntitlement struct {
	ProductType        string `json:"productType"`
	PlanID             string `json:"planId"`
	IsPrimary          bool   `json:"isPrimary"`
	SubscriptionStatus string `json:"subscriptionStatus"`
}

// Workspace is a workspace as Workspaces.List returns it.
type Workspace struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Slug         string                 `json:"slug"`
	MyRole       WorkspaceRole          `json:"myRole"`
	PlanID       string                 `json:"planId"`
	MemberCount  int                    `json:"memberCount"`
	IsCreator    bool                   `json:"isCreator"`
	ProductType  string                 `json:"productType,omitempty"`
	Entitlements []WorkspaceEntitlement `json:"entitlements"`
}

// WorkspaceDetail is a workspace as Workspaces.Get returns it.
type WorkspaceDetail struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Slug                string                 `json:"slug"`
	Domain              string                 `json:"domain,omitempty"`
	PlanID              string                 `json:"planId"`
	Status              int                    `json:"status"`
	SubscriptionStatus  string                 `json:"subscriptionStatus"`
	UserLimit           int                    `json:"userLimit"`
	CurrentUserCount    int                    `json:"currentUserCount"`
	MonthlyLabelLimit   int                    `json:"monthlyLabelLimit"`
	LabelsUsedThisMonth int                    `json:"labelsUsedThisMonth"`
	IsTrialing          bool                   `json:"isTrialing"`
	TrialExpiresAt      string                 `json:"trialExpiresAt,omitempty"`
	MyRole              WorkspaceRole          `json:"myRole"`
	CanManageMembers    bool                   `json:"canManageMembers"`
	CanManageBilling    bool                   `json:"canManageBilling"`
	CreatedAt           string                 `json:"createdAt"`
	Entitlements        []WorkspaceEntitlement `json:"entitlements"`
}

// WorkspaceMember is one member in MemberList. Status and StoreScope are the Gateway
// numbers (Status 0 = active, 1 = pending; StoreScope 0 = all stores).
type WorkspaceMember struct {
	MembershipID       string        `json:"membershipId"`
	UserID             string        `json:"userId"`
	Email              string        `json:"email"`
	Name               string        `json:"name,omitempty"`
	Role               WorkspaceRole `json:"role"`
	Status             int           `json:"status"`
	JoinedAt           string        `json:"joinedAt"`
	IsCreator          bool          `json:"isCreator"`
	IsDeactivated      bool          `json:"isDeactivated"`
	StoreScope         int           `json:"storeScope"`
	StoreConnectionIDs []string      `json:"storeConnectionIds"`
}

// PendingInvite is an invitation not yet accepted.
type PendingInvite struct {
	InviteID           string        `json:"inviteId"`
	Email              string        `json:"email"`
	Role               WorkspaceRole `json:"role"`
	ExpiresAt          string        `json:"expiresAt"`
	InvitedByName      string        `json:"invitedByName"`
	StoreScope         int           `json:"storeScope"`
	StoreConnectionIDs []string      `json:"storeConnectionIds"`
}

// MemberList is what Workspaces.ListMembers returns.
type MemberList struct {
	Members        []WorkspaceMember `json:"members"`
	PendingInvites []PendingInvite   `json:"pendingInvites"`
	UserLimit      int               `json:"userLimit"`
	ActiveCount    int               `json:"activeCount"`
	CanInviteMore  bool              `json:"canInviteMore"`
}

type Parcel struct {
	Weight     float64 `json:"weight"`
	WeightUnit string  `json:"weightUnit"`
	Length     float64 `json:"length,omitempty"`
	Width      float64 `json:"width,omitempty"`
	Height     float64 `json:"height,omitempty"`
	DimUnit    string  `json:"dimUnit,omitempty"`
}

type RateRequest struct {
	Origin      ShippingAddress `json:"origin"`
	Destination ShippingAddress `json:"destination"`
	Package     ShippingPackage `json:"package"`
	Carriers    []string        `json:"carriers,omitempty"`
	// HazmatClassification is the declared HAZMAT class, so quotes match what the
	// label route will accept. Sent only when set.
	HazmatClassification *HazmatClassification `json:"hazmatClassification,omitempty"`
	// Deprecated: use Origin, Destination, and Package.
	FromZip string `json:"fromZip,omitempty"`
	// Deprecated: use Origin, Destination, and Package.
	ToZip string `json:"toZip,omitempty"`
	// Deprecated: use Package.Weight.
	Weight float64 `json:"weight,omitempty"`
	// Deprecated: use Package.WeightUnit.
	WeightUnit string `json:"weightUnit,omitempty"`
	// Deprecated: use Package.Length.
	Length float64 `json:"length,omitempty"`
	// Deprecated: use Package.Width.
	Width float64 `json:"width,omitempty"`
	// Deprecated: use Package.Height.
	Height float64 `json:"height,omitempty"`
	// Deprecated: use Package.DimensionUnit.
	DimensionUnit string `json:"dimensionUnit,omitempty"`
	// Deprecated: use Package.PredefinedPackage.
	PackageType string `json:"packageType,omitempty"`
}

type ShippingAddress struct {
	Name          string `json:"name,omitempty"`
	Company       string `json:"company,omitempty"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2,omitempty"`
	City          string `json:"city"`
	StateProvince string `json:"stateProvince"`
	PostalCode    string `json:"postalCode"`
	CountryCode   string `json:"countryCode,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
}

// HazmatClassification is the USPS HAZMAT class, sent as the Gateway
// HazmatClassification number. There is no default: leave the request field nil
// when you have no answer, or send HazmatNone (0) to declare that the parcel holds
// no dangerous goods.
type HazmatClassification int

const (
	HazmatNone                               HazmatClassification = 0
	HazmatLimitedQuantity                    HazmatClassification = 1
	HazmatLithiumBatteryContainedInEquipment HazmatClassification = 2
	HazmatLithiumBatteryPackedWithEquipment  HazmatClassification = 3
	HazmatLithiumBatteryStandalone           HazmatClassification = 4
	HazmatOtherRestricted                    HazmatClassification = 5
	HazmatLithiumBatteryUsedDevice           HazmatClassification = 6
)

// Ptr returns a pointer to c for the optional HazmatClassification request fields.
func (c HazmatClassification) Ptr() *HazmatClassification { return &c }

type ShippingPackage struct {
	Weight            float64 `json:"weight"`
	WeightUnit        string  `json:"weightUnit,omitempty"`
	Length            float64 `json:"length,omitempty"`
	Width             float64 `json:"width,omitempty"`
	Height            float64 `json:"height,omitempty"`
	DimensionUnit     string  `json:"dimensionUnit,omitempty"`
	PredefinedPackage string  `json:"predefinedPackage,omitempty"`
}

type ShippingRate struct {
	RateIndicator   string  `json:"rateIndicator,omitempty"`
	UspsMailClass   string  `json:"uspsMailClass,omitempty"`
	IsInternational bool    `json:"isInternational,omitempty"`
	Carrier         string  `json:"carrierCode"`
	CarrierName     string  `json:"carrierName"`
	Service         string  `json:"serviceCode"`
	ServiceName     string  `json:"serviceName"`
	Rate            float64 `json:"rate"`
	Currency        string  `json:"currency"`
	EstimatedDays   int     `json:"estimatedDays"`
	DeliveryDate    string  `json:"deliveryDate,omitempty"`
}

type RateShoppingResponse struct {
	Rates    []ShippingRate `json:"rates"`
	Currency string         `json:"currency"`
}

// CustomsItem value and weightOz describe one unit, not the entire line.
type CustomsItem struct {
	Description   string  `json:"description"`
	Quantity      int     `json:"quantity"`
	Value         float64 `json:"value"`
	WeightOz      float64 `json:"weightOz"`
	HsCode        string  `json:"hsCode,omitempty"`
	OriginCountry string  `json:"originCountry"`
	// SkuOrProductID is an optional product reference.
	SkuOrProductID string `json:"skuOrProductId,omitempty"`
}
type CustomsDeclaration struct {
	ContentsType        string        `json:"contentsType"`
	ContentsExplanation string        `json:"contentsExplanation,omitempty"`
	NonDeliveryOption   string        `json:"nonDeliveryOption"`
	Currency            string        `json:"currency"`
	DeclaredValue       float64       `json:"declaredValue"`
	AesItn              string        `json:"aesItn"`
	InvoiceNumber       string        `json:"invoiceNumber,omitempty"`
	LicenseNumber       string        `json:"licenseNumber,omitempty"`
	CertificateNumber   string        `json:"certificateNumber,omitempty"`
	RestrictionType     string        `json:"restrictionType,omitempty"`
	RestrictionComments string        `json:"restrictionComments,omitempty"`
	Items               []CustomsItem `json:"items"`
}

type CreateLabelRequest struct {
	OrderID              int64               `json:"orderId,omitempty"`
	ShipDate             string              `json:"shipDate,omitempty"`
	DeclaredValue        *float64            `json:"declaredValue,omitempty"`
	CustomsDeclaration   *CustomsDeclaration `json:"customsDeclaration,omitempty"`
	RateIndicator        string              `json:"rateIndicator,omitempty"`
	CarrierCode          string              `json:"carrierCode"`
	ServiceCode          string              `json:"serviceCode"`
	Origin               ShippingAddress     `json:"origin"`
	Destination          ShippingAddress     `json:"destination"`
	Package              ShippingPackage     `json:"package"`
	MaximumPostageAmount float64             `json:"maximumPostageAmount,omitempty"`
	ConfirmationToken    string              `json:"confirmationToken,omitempty"`
	LabelFormat          string              `json:"labelFormat,omitempty"`
	// HazmatClassification is the declared HAZMAT class. A label without an
	// OrderID must declare one (HazmatNone.Ptr() for an ordinary parcel). Sent only
	// when set; the SDK never fills it in.
	HazmatClassification *HazmatClassification `json:"hazmatClassification,omitempty"`
}

// ReturnLabelOptions are the fields of a return-label request.
//
// The carrier needs the package weight and all three dimensions. Set WeightOz,
// LengthIn, WidthIn and HeightIn; any you leave out come from the workspace's
// default shipping preset. With no preset the Gateway refuses the call with a 400
// (code PackageMeasurementsRequired) that names the missing fields.
type ReturnLabelOptions struct {
	// CarrierCode defaults to "usps" on the Gateway (the only carrier supported).
	CarrierCode string `json:"carrierCode,omitempty"`
	// ServiceType is the carrier service for the return label.
	ServiceType string `json:"serviceType,omitempty"`
	// WeightOz is the package weight in ounces.
	WeightOz float64 `json:"weightOz,omitempty"`
	// LengthIn is the package length in inches.
	LengthIn float64 `json:"lengthIn,omitempty"`
	// WidthIn is the package width in inches.
	WidthIn float64 `json:"widthIn,omitempty"`
	// HeightIn is the package height in inches.
	HeightIn float64 `json:"heightIn,omitempty"`
	// PaymentConfirmed is set when the workspace return-cost policy is CustomerPays
	// and the customer has paid.
	PaymentConfirmed bool `json:"paymentConfirmed,omitempty"`
	// HazmatClassification is what the return holds. Sent only when set.
	HazmatClassification *HazmatClassification `json:"hazmatClassification,omitempty"`
	// RateIndicator is the USPS rate indicator for the return label.
	RateIndicator string `json:"rateIndicator,omitempty"`
}

// LabelPurchaseResult is the unwrapped Gateway response. Status == "Preview"
// means no purchase occurred; only a confirmed purchase contains LabelID.
type LabelPurchaseResult struct {
	Label
	Status               string  `json:"status,omitempty"`
	QuotedPostageAmount  float64 `json:"quotedPostageAmount,omitempty"`
	MaximumPostageAmount float64 `json:"maximumPostageAmount,omitempty"`
	ExpiresAt            string  `json:"expiresAt,omitempty"`
	ConfirmationToken    string  `json:"confirmationToken,omitempty"`
}

type Label struct {
	LabelID        string  `json:"labelId"`
	TrackingNumber string  `json:"trackingNumber"`
	CarrierCode    string  `json:"carrierCode"`
	Currency       string  `json:"currency"`
	LabelData      string  `json:"labelData"`
	LabelFormat    string  `json:"labelFormat"`
	Rate           float64 `json:"rate"`
	CreatedAt      string  `json:"createdAt"`
}

type TrackingEvent struct {
	Timestamp   string `json:"timestamp"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location,omitempty"`
}

type TrackingInfo struct {
	TrackingNumber    string          `json:"trackingNumber"`
	Carrier           string          `json:"carrier"`
	Status            string          `json:"status"`
	StatusDetail      string          `json:"statusDetail"`
	EstimatedDelivery string          `json:"estimatedDelivery,omitempty"`
	Events            []TrackingEvent `json:"events"`
	// Fields of the Gateway tracking result (api/shipping/track).
	CarrierCode           string `json:"carrierCode,omitempty"`
	Description           string `json:"description,omitempty"`
	Location              string `json:"location,omitempty"`
	EventDate             string `json:"eventDate,omitempty"`
	EstimatedDeliveryDate string `json:"estimatedDeliveryDate,omitempty"`
}

// AddressValidationRequest is the Gateway address-validation body. Line1, City and
// PostalCode are required; Country defaults to "US" on the Gateway.
type AddressValidationRequest struct {
	Name       string `json:"name,omitempty"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country,omitempty"`
}

// AddressValidationResult is the Gateway address-validation result.
type AddressValidationResult struct {
	// ValidationUnavailable is true when the carrier check could not complete; it is
	// not an address rejection.
	ValidationUnavailable bool              `json:"validationUnavailable"`
	CheckedAt             string            `json:"checkedAt,omitempty"`
	IsValid               bool              `json:"isValid"`
	SuggestedAddress      *SuggestedAddress `json:"suggestedAddress,omitempty"`
	IsResidential         *bool             `json:"isResidential,omitempty"`
	Warnings              []string          `json:"warnings"`
	ValidatedBy           string            `json:"validatedBy,omitempty"`
}

// SuggestedAddress is a carrier-suggested address correction.
type SuggestedAddress struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// Deprecated: the Gateway batch route takes orders with package measurements,
// not label requests. Use CreateBatchLabelRequest with CreateLabelBatch.
type BatchLabelRequest struct {
	Items []CreateLabelRequest `json:"items"`
}

// CreateBatchLabelRequest is the Gateway batch body. Addresses come from each order.
type CreateBatchLabelRequest struct {
	Items       []BatchLabelItem `json:"items"`
	CallbackURL string           `json:"callbackUrl,omitempty"`
	DryRun      bool             `json:"dryRun,omitempty"`
}

// BatchLabelItem is one order in a batch (the Gateway BatchLabelItemRequest).
// Package measurements are in ounces and inches; carrier and service are chosen by
// shipping rules when omitted.
type BatchLabelItem struct {
	OrderID                       int64    `json:"orderId"`
	OrderNumber                   string   `json:"orderNumber,omitempty"`
	PackageWeightInOunces         float64  `json:"packageWeightInOunces"`
	PackageLengthInInches         float64  `json:"packageLengthInInches"`
	PackageWidthInInches          float64  `json:"packageWidthInInches"`
	PackageHeightInInches         float64  `json:"packageHeightInInches"`
	CarrierCode                   string   `json:"carrierCode,omitempty"`
	ServiceCode                   string   `json:"serviceCode,omitempty"`
	DeliveryConfirmationRequired  bool     `json:"deliveryConfirmationRequired,omitempty"`
	SignatureConfirmationRequired bool     `json:"signatureConfirmationRequired,omitempty"`
	InsuranceCoverageRequested    bool     `json:"insuranceCoverageRequested,omitempty"`
	InsuredValue                  *float64 `json:"insuredValue,omitempty"`
	LabelFormat                   string   `json:"labelFormat,omitempty"`
	ShipDate                      string   `json:"shipDate,omitempty"`
	// HazmatClassification is what the parcel holds. Sent only when set.
	HazmatClassification *HazmatClassification `json:"hazmatClassification,omitempty"`
}

type BatchLabelJob struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
	// Deprecated: the Gateway sends TotalItems.
	Total int `json:"total"`
	// Deprecated: the Gateway sends CompletedItems.
	Completed int `json:"completed"`
	// Deprecated: the Gateway sends FailedItems.
	Failed int `json:"failed"`
	// Fields of the Gateway BatchLabelJobDto.
	TotalItems         int      `json:"totalItems"`
	CompletedItems     int      `json:"completedItems"`
	FailedItems        int      `json:"failedItems"`
	IsDryRun           bool     `json:"isDryRun"`
	TotalEstimatedCost *float64 `json:"totalEstimatedCost,omitempty"`
	CreatedAt          string   `json:"createdAt,omitempty"`
	CompletedAt        string   `json:"completedAt,omitempty"`
	StatusURL          string   `json:"statusUrl,omitempty"`
}

type Order struct {
	OrderNumber string      `json:"orderNumber"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items,omitempty"`
	CreatedAt   string      `json:"createdAt"`
}

type OrderItem struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// WebhookEventType is an event name FlexOps sends in the X-FlexOps-Event header
// and matches subscriptions against. Names are PascalCase and matched exactly: a
// dotted name such as "label.created" matches nothing.
type WebhookEventType string

// The events FlexOps sends. WebhookEventTrackingUpdated and
// WebhookEventTrackingDelivered are sent only to workspaces in the USPS
// push-tracking pilot.
const (
	WebhookEventLabelCreated      WebhookEventType = "LabelCreated"
	WebhookEventLabelCanceled     WebhookEventType = "LabelCanceled"
	WebhookEventBatchCompleted    WebhookEventType = "BatchCompleted"
	WebhookEventTrackingUpdated   WebhookEventType = "TrackingUpdated"
	WebhookEventTrackingDelivered WebhookEventType = "TrackingDelivered"
)

// WebhookSubscription is a webhook subscription as the Gateway returns it.
type WebhookSubscription struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	// EventTypes is the subscribed event names, comma-separated, or "*" for every event.
	EventTypes string `json:"eventTypes"`
	// Status is "Active", "Paused" or "Disabled".
	Status              string `json:"status"`
	ConsecutiveFailures int    `json:"consecutiveFailures"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt,omitempty"`
	// Secret is the "whsec_..." signing secret. It is returned only by Create; store it.
	Secret string `json:"secret,omitempty"`
}

// CreateWebhookRequest creates a webhook subscription.
type CreateWebhookRequest struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	// Events lists the events to receive. Leave it nil to receive every event; an
	// empty, non-nil slice is refused. It is sent as the Gateway's comma-separated
	// eventTypes.
	Events []WebhookEventType `json:"-"`
}

// UpdateWebhookRequest changes a webhook subscription. Unset fields are left as they are.
type UpdateWebhookRequest struct {
	URL string `json:"url,omitempty"`
	// Description replaces the description when set; a pointer to "" clears it.
	Description *string `json:"description,omitempty"`
	// Events replaces the subscribed events when non-nil; an empty, non-nil slice
	// is refused. It is sent as the Gateway's comma-separated eventTypes.
	Events []WebhookEventType `json:"-"`
	// Status is "Active", "Paused" or "Disabled".
	Status string `json:"status,omitempty"`
}

// WalletBalance is the workspace wallet balance (GET /api/wallet/my-balance).
type WalletBalance struct {
	HasWallet    bool    `json:"hasWallet"`
	BalanceCents int64   `json:"balanceCents"`
	Balance      float64 `json:"balance"`
}

// InsuranceQuoteRequest prices coverage. DeclaredValue and CarrierCode are required.
type InsuranceQuoteRequest struct {
	DeclaredValue  float64 `json:"declaredValue"`
	CarrierCode    string  `json:"carrierCode"`
	TrackingNumber string  `json:"trackingNumber,omitempty"`
	Provider       string  `json:"provider,omitempty"`
}

type InsuranceQuote struct {
	DeclaredValue  float64 `json:"declaredValue"`
	CoverageAmount float64 `json:"coverageAmount"`
	Premium        float64 `json:"premium"`
	CoverageTier   string  `json:"coverageTier"`
	CarrierCode    string  `json:"carrierCode"`
	CoverageDays   int     `json:"coverageDays"`
	Provider       string  `json:"provider"`
}

// PurchaseInsuranceRequest buys a policy for a shipped label.
type PurchaseInsuranceRequest struct {
	TrackingNumber string  `json:"trackingNumber"`
	CarrierCode    string  `json:"carrierCode"`
	DeclaredValue  float64 `json:"declaredValue"`
	CoverageTier   string  `json:"coverageTier,omitempty"`
	Provider       string  `json:"provider,omitempty"`
}

type InsurancePolicy struct {
	PolicyID          string  `json:"policyId"`
	TrackingNumber    string  `json:"trackingNumber"`
	CarrierCode       string  `json:"carrierCode"`
	DeclaredValue     float64 `json:"declaredValue"`
	CoverageAmount    float64 `json:"coverageAmount"`
	Premium           float64 `json:"premium"`
	Status            string  `json:"status"`
	CoverageTier      string  `json:"coverageTier,omitempty"`
	PurchasedAt       string  `json:"purchasedAt"`
	ExpiresAt         string  `json:"expiresAt"`
	CancelledAt       string  `json:"cancelledAt,omitempty"`
	Provider          string  `json:"provider,omitempty"`
	ProviderReference string  `json:"providerReference,omitempty"`
	IsSandbox         bool    `json:"isSandbox"`
}

// FileClaimRequest files a claim against a policy. EvidenceJSON is optional.
type FileClaimRequest struct {
	PolicyID     string  `json:"policyId"`
	ClaimReason  string  `json:"claimReason"`
	ClaimAmount  float64 `json:"claimAmount"`
	EvidenceJSON string  `json:"evidenceJson,omitempty"`
}

type InsuranceClaim struct {
	ClaimID         string   `json:"claimId"`
	PolicyID        string   `json:"policyId"`
	ClaimReason     string   `json:"claimReason"`
	ClaimAmount     float64  `json:"claimAmount"`
	Status          string   `json:"status"`
	SubmittedAt     string   `json:"submittedAt"`
	ResolvedAt      string   `json:"resolvedAt,omitempty"`
	ApprovedAmount  *float64 `json:"approvedAmount,omitempty"`
	ResolutionNotes string   `json:"resolutionNotes,omitempty"`
	IsSandbox       bool     `json:"isSandbox"`
}

// ReturnRequest opens an RMA (the Gateway CreateRmaRequest). OrderID, Reason and at
// least one item are required.
type ReturnRequest struct {
	OrderID       int64        `json:"orderId"`
	OrderNumber   string       `json:"orderNumber,omitempty"`
	Reason        string       `json:"reason"`
	Items         []ReturnItem `json:"items"`
	Notes         string       `json:"notes,omitempty"`
	CustomerEmail string       `json:"customerEmail,omitempty"`
	ReasonID      *int         `json:"reasonId,omitempty"`
}

type ReturnItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Reason   string `json:"reason,omitempty"`
}

// RmaItem is a returned line as the Gateway reports it.
type RmaItem struct {
	SKU           string  `json:"sku"`
	Description   string  `json:"description,omitempty"`
	Quantity      int     `json:"quantity"`
	UnitPrice     float64 `json:"unitPrice"`
	ExtendedPrice float64 `json:"extendedPrice"`
	Reason        string  `json:"reason,omitempty"`
}

// Rma is a return merchandise authorization, as Returns.Create returns it.
type Rma struct {
	RmaNumber                 string    `json:"rmaNumber"`
	RmaID                     int64     `json:"rmaId"`
	OrderID                   int64     `json:"orderId"`
	OrderNumber               string    `json:"orderNumber,omitempty"`
	CustomerEmail             string    `json:"customerEmail,omitempty"`
	Status                    string    `json:"status"`
	Reason                    string    `json:"reason,omitempty"`
	CreatedAt                 string    `json:"createdAt"`
	Items                     []RmaItem `json:"items"`
	ItemCount                 int       `json:"itemCount"`
	ReturnLabelTrackingNumber string    `json:"returnLabelTrackingNumber,omitempty"`
	Notes                     string    `json:"notes,omitempty"`
	ReturnID                  *int64    `json:"returnId,omitempty"`
	DispositionStatus         string    `json:"dispositionStatus,omitempty"`
}

// ReturnListOptions filters Returns.List. StartDate and EndDate are ISO 8601 dates.
type ReturnListOptions struct {
	StartDate string
	EndDate   string
	Page      int
	PageSize  int
}

// ReturnAuthorization is a return as Returns.List (summary fields) and Returns.Get
// (with Items, Reason, Notes, TrackingNumber and CarrierCode) return it.
type ReturnAuthorization struct {
	ReturnID        int64     `json:"returnId"`
	RmaNumber       string    `json:"rmaNumber,omitempty"`
	OrderID         int64     `json:"orderId"`
	Status          string    `json:"status"`
	ReturnDate      string    `json:"returnDate"`
	DispositionType string    `json:"dispositionType,omitempty"`
	Items           []RmaItem `json:"items,omitempty"`
	Reason          string    `json:"reason,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	TrackingNumber  string    `json:"trackingNumber,omitempty"`
	CarrierCode     string    `json:"carrierCode,omitempty"`
}

// ReturnDisposition is what happens to an inspected item, sent as the Gateway number.
type ReturnDisposition int

const (
	ReturnDispositionRestock ReturnDisposition = 0
	ReturnDispositionDispose ReturnDisposition = 1
	ReturnDispositionForward ReturnDisposition = 2
)

// ReceiveReturnRequest records a received RMA (the Gateway ReceiveRmaRequest).
// WarehouseID defaults to 1 on the Gateway when left 0.
type ReceiveReturnRequest struct {
	RmaNumber      string          `json:"rmaNumber"`
	WarehouseID    int             `json:"warehouseId,omitempty"`
	InspectedItems []InspectedItem `json:"inspectedItems"`
}

// InspectedItem is one inspected line of a received RMA.
type InspectedItem struct {
	PartNumber      string            `json:"partNumber"`
	PartDescription string            `json:"partDescription,omitempty"`
	Quantity        int               `json:"quantity"`
	Disposition     ReturnDisposition `json:"disposition"`
	InspectionNotes string            `json:"inspectionNotes,omitempty"`
}

type ShipmentsTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type CarrierSummary struct {
	Carrier    string  `json:"carrier"`
	Shipments  int     `json:"shipments"`
	TotalSpend float64 `json:"totalSpend"`
}

type ShippingRule struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Priority    int             `json:"priority"`
	IsEnabled   bool            `json:"isEnabled"`
	WarehouseID string          `json:"warehouseId,omitempty"`
	Conditions  []RuleCondition `json:"conditions"`
	Actions     []RuleAction    `json:"actions"`
	CreatedAt   string          `json:"createdAt,omitempty"`
	UpdatedAt   string          `json:"updatedAt,omitempty"`
}

type RuleCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Group    *int   `json:"group,omitempty"`
}

type RuleAction struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// CreateShippingRuleRequest creates a rule. Name and at least one condition and action
// are required; IsEnabled is sent as given.
type CreateShippingRuleRequest struct {
	Name        string          `json:"name"`
	Priority    int             `json:"priority"`
	IsEnabled   bool            `json:"isEnabled"`
	WarehouseID string          `json:"warehouseId,omitempty"`
	Conditions  []RuleCondition `json:"conditions"`
	Actions     []RuleAction    `json:"actions"`
}

// UpdateShippingRuleRequest changes a rule. Nil fields are left as they are.
type UpdateShippingRuleRequest struct {
	Name       *string         `json:"name,omitempty"`
	Priority   *int            `json:"priority,omitempty"`
	IsEnabled  *bool           `json:"isEnabled,omitempty"`
	Conditions []RuleCondition `json:"conditions,omitempty"`
	Actions    []RuleAction    `json:"actions,omitempty"`
}

// RulePriority is one entry of a Rules.Reorder call.
type RulePriority struct {
	RuleID   string `json:"ruleId"`
	Priority int    `json:"priority"`
}

// TestShippingRuleRequest is a sample shipment for Rules.Test. DestinationCountry is
// required ("US" if unsure).
type TestShippingRuleRequest struct {
	WeightOz              float64  `json:"weightOz"`
	DeclaredValue         float64  `json:"declaredValue"`
	DestinationCountry    string   `json:"destinationCountry"`
	DestinationState      string   `json:"destinationState,omitempty"`
	DestinationPostalCode string   `json:"destinationPostalCode,omitempty"`
	OrderSource           string   `json:"orderSource,omitempty"`
	ItemCount             int      `json:"itemCount,omitempty"`
	ShipMethod            string   `json:"shipMethod,omitempty"`
	DestinationZone       string   `json:"destinationZone,omitempty"`
	Tags                  []string `json:"tags,omitempty"`
}

type TestShippingRuleResponse struct {
	Traces          []RuleEvaluationTrace `json:"traces"`
	WinningRuleID   string                `json:"winningRuleId,omitempty"`
	WinningRuleName string                `json:"winningRuleName,omitempty"`
	AppliedActions  []RuleAction          `json:"appliedActions"`
}

type RuleEvaluationTrace struct {
	RuleID          string           `json:"ruleId"`
	RuleName        string           `json:"ruleName"`
	Priority        int              `json:"priority"`
	Matched         bool             `json:"matched"`
	IsWinner        bool             `json:"isWinner"`
	ConditionTraces []ConditionTrace `json:"conditionTraces"`
}

type ConditionTrace struct {
	Field         string `json:"field"`
	Operator      string `json:"operator"`
	ExpectedValue string `json:"expectedValue"`
	ActualValue   string `json:"actualValue"`
	Matched       bool   `json:"matched"`
}

// SchedulePickupRequest schedules a carrier pickup. Dates and times are ISO 8601;
// Country defaults to "US" and PackageCount to 1 on the Gateway.
type SchedulePickupRequest struct {
	CarrierCode         string   `json:"carrierCode"`
	PickupDate          string   `json:"pickupDate"`
	Street              string   `json:"street"`
	City                string   `json:"city"`
	State               string   `json:"state"`
	PostalCode          string   `json:"postalCode"`
	Country             string   `json:"country,omitempty"`
	PackageCount        int      `json:"packageCount,omitempty"`
	TotalWeightOz       *float64 `json:"totalWeightOz,omitempty"`
	SpecialInstructions string   `json:"specialInstructions,omitempty"`
	ReadyTime           string   `json:"readyTime,omitempty"`
	CloseTime           string   `json:"closeTime,omitempty"`
	ContactName         string   `json:"contactName,omitempty"`
	ContactPhone        string   `json:"contactPhone,omitempty"`
}

// Pickup is a scheduled pickup as the Gateway returns it.
type Pickup struct {
	PickupID           string `json:"pickupId,omitempty"`
	CarrierCode        string `json:"carrierCode"`
	PickupDate         string `json:"pickupDate"`
	Status             string `json:"status"`
	ConfirmationNumber string `json:"confirmationNumber,omitempty"`
	WindowStart        string `json:"windowStart,omitempty"`
	WindowEnd          string `json:"windowEnd,omitempty"`
	Street             string `json:"street,omitempty"`
	City               string `json:"city,omitempty"`
	State              string `json:"state,omitempty"`
	PackageCount       int    `json:"packageCount"`
	CreatedAt          string `json:"createdAt,omitempty"`
	Message            string `json:"message,omitempty"`
	Success            bool   `json:"success"`
	ErrorMessage       string `json:"errorMessage,omitempty"`
	IsSandbox          bool   `json:"isSandbox"`
}

type ScanFormRequest struct {
	Carrier         string   `json:"carrier"`
	TrackingNumbers []string `json:"trackingNumbers"`
}

type ScanForm struct {
	ID        string `json:"id"`
	Carrier   string `json:"carrier"`
	FormURL   string `json:"formUrl"`
	CreatedAt string `json:"createdAt"`
}

type CarrierRecommendationRequest struct {
	OriginPostalCode      string  `json:"originPostalCode"`
	DestinationPostalCode string  `json:"destinationPostalCode"`
	WeightOz              float64 `json:"weightOz"`
	Priority              string  `json:"priority,omitempty"`
}

// CarrierRecommendation is one ranked carrier for a lane.
//
// AvgCost is the workspace's own average label cost in USD, and 0 when the
// carrier was scored on FlexOps network data: network cost is not reported, so
// 0 does not mean free. The API's dataSource field ("workspace" or "network")
// says which; this type does not include it yet. ShipmentCount is exact for
// workspace data and a band floor (5, 10, 25, 50, 100, 250, 500, 1000) for
// network data.
type CarrierRecommendation struct {
	CarrierCode    string  `json:"carrierCode"`
	Score          float64 `json:"score"`
	OnTimePercent  float64 `json:"onTimePercent"`
	AvgTransitDays float64 `json:"avgTransitDays"`
	AvgCost        float64 `json:"avgCost"`
	ShipmentCount  int     `json:"shipmentCount"`
	ExceptionRate  float64 `json:"exceptionRate"`
	Reason         string  `json:"reason,omitempty"`
}

// CarrierRecommendationResponse lists ranked carriers for a lane. SampleSize is
// exact for the workspace's own data and a band floor for each carrier scored
// on FlexOps network data.
type CarrierRecommendationResponse struct {
	Lane            string                  `json:"lane"`
	SampleSize      int                     `json:"sampleSize"`
	Recommendations []CarrierRecommendation `json:"recommendations"`
}

type DeliveryPredictionRequest struct {
	CarrierCode           string `json:"carrierCode"`
	ServiceCode           string `json:"serviceCode,omitempty"`
	OriginPostalCode      string `json:"originPostalCode"`
	DestinationPostalCode string `json:"destinationPostalCode"`
	ShipDate              string `json:"shipDate,omitempty"`
}

// DeliveryPredictionResponse is a delivery date range for a carrier and lane.
// SampleSize is exact for the workspace's own data and a band floor for
// FlexOps network data.
type DeliveryPredictionResponse struct {
	CarrierCode           string  `json:"carrierCode"`
	PredictedDeliveryDate string  `json:"predictedDeliveryDate"`
	EarliestDelivery      string  `json:"earliestDelivery"`
	LatestDelivery        string  `json:"latestDelivery"`
	WorstCaseDelivery     string  `json:"worstCaseDelivery"`
	PredictedTransitDays  int     `json:"predictedTransitDays"`
	Confidence            float64 `json:"confidence"`
	OnTimeRate            float64 `json:"onTimeRate"`
	SampleSize            int     `json:"sampleSize"`
}

// HsCodeSearchRequest searches HS codes. Query is required.
type HsCodeSearchRequest struct {
	Query              string
	DestinationCountry string
	MaxResults         int
}

type HsCodeResult struct {
	HsCode             string   `json:"hsCode"`
	Description        string   `json:"description"`
	Chapter            string   `json:"chapter"`
	ChapterDescription string   `json:"chapterDescription"`
	DutyRatePercent    *float64 `json:"dutyRatePercent,omitempty"`
	Notes              string   `json:"notes,omitempty"`
}

// LandedCostRequest estimates duty and tax. Items use per-unit value and weight.
type LandedCostRequest struct {
	DestinationCountry string        `json:"destinationCountry"`
	Items              []CustomsItem `json:"items"`
	ShippingCost       float64       `json:"shippingCost"`
}

type LandedCostEstimate struct {
	ShippingCost    float64 `json:"shippingCost"`
	EstimatedDuty   float64 `json:"estimatedDuty"`
	EstimatedTax    float64 `json:"estimatedTax"`
	TotalLandedCost float64 `json:"totalLandedCost"`
	Currency        string  `json:"currency"`
	DutyRateApplied string  `json:"dutyRateApplied,omitempty"`
	IsEstimate      bool    `json:"isEstimate"`
	Disclaimer      string  `json:"disclaimer,omitempty"`
}

type CarbonOffsetQuoteRequest struct {
	WeightLbs     float64 `json:"weightLbs"`
	DistanceMiles float64 `json:"distanceMiles"`
}

type CarbonOffsetQuote struct {
	CarbonKg               float64 `json:"carbonKg"`
	CostUsd                float64 `json:"costUsd"`
	EquivalentTreesPlanted float64 `json:"equivalentTreesPlanted"`
	Provider               string  `json:"provider"`
}

type CarbonOffsetPurchaseRequest struct {
	TrackingNumber string  `json:"trackingNumber"`
	WeightLbs      float64 `json:"weightLbs"`
	DistanceMiles  float64 `json:"distanceMiles"`
}

// CarbonOffset is a purchased offset.
type CarbonOffset struct {
	OffsetID          string  `json:"offsetId"`
	TrackingNumber    string  `json:"trackingNumber"`
	CarbonKg          float64 `json:"carbonKg"`
	CostUsd           float64 `json:"costUsd"`
	Status            string  `json:"status"`
	Provider          string  `json:"provider"`
	ProviderReference string  `json:"providerReference,omitempty"`
	IsAutoOffset      bool    `json:"isAutoOffset"`
	PurchasedAt       string  `json:"purchasedAt"`
}

type CarbonOffsetSummary struct {
	TotalCarbonKg          float64 `json:"totalCarbonKg"`
	TotalCostUsd           float64 `json:"totalCostUsd"`
	TotalOffsets           int     `json:"totalOffsets"`
	EquivalentTreesPlanted float64 `json:"equivalentTreesPlanted"`
}

type CarbonOffsetSettings struct {
	AutoOffsetEnabled bool `json:"autoOffsetEnabled"`
}

// RecurrenceInterval is how often a recurring shipment runs, sent as the Gateway number.
// RecurrenceCustom uses CronExpression.
type RecurrenceInterval int

const (
	RecurrenceDaily    RecurrenceInterval = 0
	RecurrenceWeekly   RecurrenceInterval = 1
	RecurrenceBiweekly RecurrenceInterval = 2
	RecurrenceMonthly  RecurrenceInterval = 3
	RecurrenceCustom   RecurrenceInterval = 4
)

// RecurringShipment is a recurring shipment schedule as the Gateway returns it.
type RecurringShipment struct {
	ID             string             `json:"id"`
	WorkspaceID    string             `json:"workspaceId"`
	PresetID       string             `json:"presetId,omitempty"`
	PresetName     string             `json:"presetName,omitempty"`
	AddressID      string             `json:"addressId,omitempty"`
	AddressSummary string             `json:"addressSummary,omitempty"`
	Name           string             `json:"name"`
	Interval       RecurrenceInterval `json:"interval"`
	IntervalName   string             `json:"intervalName"`
	CronExpression string             `json:"cronExpression,omitempty"`
	NextRunUtc     string             `json:"nextRunUtc"`
	LastRunUtc     string             `json:"lastRunUtc,omitempty"`
	IsActive       bool               `json:"isActive"`
	FailureCount   int                `json:"failureCount"`
	MaxFailures    int                `json:"maxFailures"`
	LastError      string             `json:"lastError,omitempty"`
	CreatedAt      string             `json:"createdAt"`
}

// CreateRecurringShipmentRequest creates a schedule. Name and FirstRunUtc (ISO 8601)
// are required; a MaxFailures of 0 uses the Gateway default (3).
type CreateRecurringShipmentRequest struct {
	Name           string             `json:"name"`
	PresetID       string             `json:"presetId,omitempty"`
	AddressID      string             `json:"addressId,omitempty"`
	Interval       RecurrenceInterval `json:"interval"`
	CronExpression string             `json:"cronExpression,omitempty"`
	FirstRunUtc    string             `json:"firstRunUtc"`
	MaxFailures    int                `json:"maxFailures,omitempty"`
}

// UpdateRecurringShipmentRequest changes a schedule. Unset fields are left as they are.
type UpdateRecurringShipmentRequest struct {
	Name           string              `json:"name,omitempty"`
	PresetID       string              `json:"presetId,omitempty"`
	AddressID      string              `json:"addressId,omitempty"`
	Interval       *RecurrenceInterval `json:"interval,omitempty"`
	CronExpression string              `json:"cronExpression,omitempty"`
	MaxFailures    *int                `json:"maxFailures,omitempty"`
}

// ReportFormat is a scheduled report's file format, sent as the Gateway number.
type ReportFormat int

const (
	ReportFormatCsv  ReportFormat = 0
	ReportFormatJSON ReportFormat = 1
)

// ReportFrequency is how often a report is sent, sent as the Gateway number.
type ReportFrequency int

const (
	ReportFrequencyDaily   ReportFrequency = 0
	ReportFrequencyWeekly  ReportFrequency = 1
	ReportFrequencyMonthly ReportFrequency = 2
)

// ReportSchedule is a report schedule as the Gateway returns it. TimeOfDayUtc is
// "HH:mm:ss"; RecipientEmails is comma-separated.
type ReportSchedule struct {
	ID                   string          `json:"id"`
	WorkspaceID          string          `json:"workspaceId"`
	ReportType           string          `json:"reportType"`
	Format               ReportFormat    `json:"format"`
	FilterParametersJSON string          `json:"filterParametersJson,omitempty"`
	FrequencyType        ReportFrequency `json:"frequencyType"`
	FrequencyName        string          `json:"frequencyName"`
	DayOfWeek            *int            `json:"dayOfWeek,omitempty"`
	DayOfMonth           *int            `json:"dayOfMonth,omitempty"`
	TimeOfDayUtc         string          `json:"timeOfDayUtc"`
	NextRunUtc           string          `json:"nextRunUtc"`
	LastRunUtc           string          `json:"lastRunUtc,omitempty"`
	IsEnabled            bool            `json:"isEnabled"`
	RecipientEmails      string          `json:"recipientEmails"`
	FailureCount         int             `json:"failureCount"`
	MaxFailures          int             `json:"maxFailures"`
	LastError            string          `json:"lastError,omitempty"`
	CreatedAt            string          `json:"createdAt"`
}

// CreateReportScheduleRequest creates a schedule. ReportType (free text, e.g.
// "shipments"), TimeOfDayUtc ("HH:mm:ss") and RecipientEmails (comma-separated) are
// required. DayOfWeek (0-6, Sunday 0) is for weekly reports, DayOfMonth (1-28) for
// monthly ones; a MaxFailures of 0 uses the Gateway default (3).
type CreateReportScheduleRequest struct {
	ReportType           string          `json:"reportType"`
	Format               ReportFormat    `json:"format"`
	FilterParametersJSON string          `json:"filterParametersJson,omitempty"`
	FrequencyType        ReportFrequency `json:"frequencyType"`
	DayOfWeek            *int            `json:"dayOfWeek,omitempty"`
	DayOfMonth           *int            `json:"dayOfMonth,omitempty"`
	TimeOfDayUtc         string          `json:"timeOfDayUtc"`
	RecipientEmails      string          `json:"recipientEmails"`
	MaxFailures          int             `json:"maxFailures,omitempty"`
}

// UpdateReportScheduleRequest changes a schedule. Unset fields are left as they are.
type UpdateReportScheduleRequest struct {
	ReportType           string           `json:"reportType,omitempty"`
	Format               *ReportFormat    `json:"format,omitempty"`
	FilterParametersJSON string           `json:"filterParametersJson,omitempty"`
	FrequencyType        *ReportFrequency `json:"frequencyType,omitempty"`
	DayOfWeek            *int             `json:"dayOfWeek,omitempty"`
	DayOfMonth           *int             `json:"dayOfMonth,omitempty"`
	TimeOfDayUtc         string           `json:"timeOfDayUtc,omitempty"`
	RecipientEmails      string           `json:"recipientEmails,omitempty"`
	MaxFailures          *int             `json:"maxFailures,omitempty"`
}

// TrackingEmailEvent is the tracking event a template is sent for, sent as the Gateway number.
type TrackingEmailEvent int

const (
	TrackingEmailLabelCreated   TrackingEmailEvent = 0
	TrackingEmailInTransit      TrackingEmailEvent = 1
	TrackingEmailOutForDelivery TrackingEmailEvent = 2
	TrackingEmailDelivered      TrackingEmailEvent = 3
	TrackingEmailException      TrackingEmailEvent = 4
	TrackingEmailReturnToSender TrackingEmailEvent = 5
)

// EmailTemplate is a tracking email template as the Gateway returns it.
type EmailTemplate struct {
	ID               string             `json:"id"`
	WorkspaceID      string             `json:"workspaceId"`
	EventType        TrackingEmailEvent `json:"eventType"`
	EventTypeName    string             `json:"eventTypeName"`
	SubjectTemplate  string             `json:"subjectTemplate"`
	HTMLBodyTemplate string             `json:"htmlBodyTemplate"`
	IsEnabled        bool               `json:"isEnabled"`
	CreatedAt        string             `json:"createdAt"`
	UpdatedAt        string             `json:"updatedAt"`
}

// UpsertEmailTemplateRequest creates or replaces the template for EventType.
// SubjectTemplate and HTMLBodyTemplate are required; a nil IsEnabled means enabled.
type UpsertEmailTemplateRequest struct {
	EventType        TrackingEmailEvent `json:"eventType"`
	SubjectTemplate  string             `json:"subjectTemplate"`
	HTMLBodyTemplate string             `json:"htmlBodyTemplate"`
	IsEnabled        *bool              `json:"isEnabled,omitempty"`
}

// PreviewEmailTemplateRequest is a subject and HTML body to render. Both are required.
type PreviewEmailTemplateRequest struct {
	SubjectTemplate  string `json:"subjectTemplate"`
	HTMLBodyTemplate string `json:"htmlBodyTemplate"`
}

// EmailTemplatePreview is a rendered template.
type EmailTemplatePreview struct {
	RenderedSubject  string `json:"renderedSubject"`
	RenderedHTMLBody string `json:"renderedHtmlBody"`
}

// MultiPackageShipmentRequest is the Gateway MultiPackageShipmentRequest for
// Shipping.CreateMultiPackage: one order shipped in 1 to 99 boxes with the same
// carrier and service. The addresses come from the order.
type MultiPackageShipmentRequest struct {
	// OrderID is the order to ship. Required and positive.
	OrderID     int64  `json:"orderId"`
	CarrierCode string `json:"carrierCode"`
	ServiceCode string `json:"serviceCode,omitempty"`
	// LabelFormat defaults to PDF on the Gateway.
	LabelFormat string `json:"labelFormat,omitempty"`
	// ShipDate is an ISO 8601 date; today when empty.
	ShipDate string `json:"shipDate,omitempty"`
	// HazmatClassification is what the shipment holds. Multi-package labels cannot
	// carry HAZMAT yet: set HazmatNone.Ptr() for no dangerous goods, or leave it nil
	// for no answer. The Gateway refuses any other class (400 HazmatNotSupported);
	// buy those labels one at a time with Shipping.CreateLabel.
	HazmatClassification *HazmatClassification `json:"hazmatClassification,omitempty"`
	Packages             []MultiPackageDetail  `json:"packages"`
}

// MultiPackageDetail is one box (the Gateway PackageDetail). Weight is in ounces
// (0.1 to 999999) and each dimension in inches (0.1 to 999); all four are required.
type MultiPackageDetail struct {
	WeightInOunces                float64 `json:"weightInOunces"`
	LengthInInches                float64 `json:"lengthInInches"`
	WidthInInches                 float64 `json:"widthInInches"`
	HeightInInches                float64 `json:"heightInInches"`
	DeliveryConfirmationRequired  bool    `json:"deliveryConfirmationRequired,omitempty"`
	SignatureConfirmationRequired bool    `json:"signatureConfirmationRequired,omitempty"`
	InsuranceCoverageRequested    bool    `json:"insuranceCoverageRequested,omitempty"`
	// InsuredValue is the coverage in USD. It needs InsuranceCoverageRequested; nil
	// insures the order value.
	InsuredValue *float64 `json:"insuredValue,omitempty"`
	// PackageReference is your label for the box, up to 100 characters (for example
	// "Box 1 of 3").
	PackageReference string `json:"packageReference,omitempty"`
}

// MultiPackageShipmentResponse is the Gateway MultiPackageShipmentResponse.
type MultiPackageShipmentResponse struct {
	OrderID            int64                     `json:"orderId"`
	CarrierCode        string                    `json:"carrierCode"`
	TotalPackages      int                       `json:"totalPackages"`
	SuccessfulPackages int                       `json:"successfulPackages"`
	FailedPackages     int                       `json:"failedPackages"`
	TotalCost          float64                   `json:"totalCost"`
	Labels             []MultiPackageLabelResult `json:"labels"`
}

// MultiPackageLabelResult is one box's outcome (the Gateway PackageLabelResult).
// ErrorCode OutcomeUnknown means a label may exist: check the order before retrying.
type MultiPackageLabelResult struct {
	PackageIndex     int      `json:"packageIndex"`
	PackageReference string   `json:"packageReference,omitempty"`
	TrackingNumber   string   `json:"trackingNumber,omitempty"`
	Rate             *float64 `json:"rate,omitempty"`
	Success          bool     `json:"success"`
	ErrorMessage     string   `json:"errorMessage,omitempty"`
	ErrorCode        string   `json:"errorCode,omitempty"`
}

// RmaListOptions pages Returns.ListRmas. Unset options are not sent (the Gateway
// defaults to page 1 of 25).
type RmaListOptions struct {
	Page     int
	PageSize int
}
