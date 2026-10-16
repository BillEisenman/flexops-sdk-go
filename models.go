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

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type Workspace struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	PlanID   string `json:"planId"`
	IsActive bool   `json:"isActive"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

type WorkspaceMember struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type Address struct {
	Name    string `json:"name,omitempty"`
	Street1 string `json:"street1"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
	Phone   string `json:"phone,omitempty"`
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

// AddressValidationRequest is the Gateway address-validation body.
type AddressValidationRequest struct {
	Address ShippingAddress `json:"address"`
}

// AddressValidationResponse is the Gateway address-validation result.
type AddressValidationResponse struct {
	ValidationUnavailable bool              `json:"validationUnavailable"`
	IsValid               bool              `json:"isValid"`
	ConfidenceScore       float64           `json:"confidenceScore"`
	ValidatedAddress      *ShippingAddress  `json:"validatedAddress,omitempty"`
	IsResidential         *bool             `json:"isResidential,omitempty"`
	AddressClassification string            `json:"addressClassification,omitempty"`
	Provider              string            `json:"provider"`
	Messages              []string          `json:"messages,omitempty"`
	Suggestions           []ShippingAddress `json:"suggestions,omitempty"`
}

// Deprecated: returned by the deprecated ValidateAddress; use AddressValidationResponse.
type AddressValidationResult struct {
	IsValid     bool      `json:"isValid"`
	Normalized  *Address  `json:"normalized,omitempty"`
	Suggestions []Address `json:"suggestions,omitempty"`
	Errors      []string  `json:"errors,omitempty"`
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

type WalletBalance struct {
	Balance           float64 `json:"balance"`
	Currency          string  `json:"currency"`
	AutoReloadEnabled bool    `json:"autoReloadEnabled"`
}

type InsuranceQuote struct {
	Provider string  `json:"provider"`
	Premium  float64 `json:"premium"`
	Coverage float64 `json:"coverage"`
}

type InsurancePolicy struct {
	PolicyID       string  `json:"policyId"`
	TrackingNumber string  `json:"trackingNumber"`
	Provider       string  `json:"provider"`
	Coverage       float64 `json:"coverage"`
	Premium        float64 `json:"premium"`
	Status         string  `json:"status"`
}

type ReturnRequest struct {
	OrderNumber string       `json:"orderNumber"`
	Reason      string       `json:"reason"`
	Items       []ReturnItem `json:"items"`
}

type ReturnItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Reason   string `json:"reason,omitempty"`
}

type ReturnAuthorization struct {
	ID          string       `json:"id"`
	OrderNumber string       `json:"orderNumber"`
	Status      string       `json:"status"`
	Items       []ReturnItem `json:"items"`
}

type ApiKeyInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Prefix    string `json:"prefix"`
	CreatedAt string `json:"createdAt"`
}

type CreateApiKeyRequest struct {
	Name string `json:"name"`
}

type CreateApiKeyResponse struct {
	ID  string `json:"id"`
	Key string `json:"key"`
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
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Priority   int             `json:"priority"`
	IsActive   bool            `json:"isActive"`
	Conditions []RuleCondition `json:"conditions"`
	Actions    []RuleAction    `json:"actions"`
}

type RuleCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Group    int    `json:"group,omitempty"`
}

type RuleAction struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PickupRequest struct {
	Carrier   string `json:"carrier"`
	Date      string `json:"date"`
	TimeStart string `json:"timeStart,omitempty"`
	TimeEnd   string `json:"timeEnd,omitempty"`
}

type PickupConfirmation struct {
	ID      string `json:"id"`
	Carrier string `json:"carrier"`
	Date    string `json:"date"`
	Status  string `json:"status"`
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
