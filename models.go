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
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2,omitempty"`
	City          string `json:"city"`
	StateProvince string `json:"stateProvince"`
	PostalCode    string `json:"postalCode"`
	CountryCode   string `json:"countryCode,omitempty"`
}

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
	Carrier       string  `json:"carrierCode"`
	CarrierName   string  `json:"carrierName"`
	Service       string  `json:"serviceCode"`
	ServiceName   string  `json:"serviceName"`
	Rate          float64 `json:"rate"`
	Currency      string  `json:"currency"`
	EstimatedDays int     `json:"estimatedDays"`
	DeliveryDate  string  `json:"deliveryDate,omitempty"`
}

type RateShoppingResponse struct {
	Rates    []ShippingRate `json:"rates"`
	Currency string         `json:"currency"`
}

type CreateLabelRequest struct {
	Carrier     string   `json:"carrier"`
	Service     string   `json:"service"`
	FromAddress *Address `json:"fromAddress"`
	ToAddress   *Address `json:"toAddress"`
	Parcel      *Parcel  `json:"parcel"`
}

type Label struct {
	LabelID        string  `json:"labelId"`
	TrackingNumber string  `json:"trackingNumber"`
	Carrier        string  `json:"carrier"`
	Service        string  `json:"service"`
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
}

type AddressValidationResult struct {
	IsValid     bool      `json:"isValid"`
	Normalized  *Address  `json:"normalized,omitempty"`
	Suggestions []Address `json:"suggestions,omitempty"`
	Errors      []string  `json:"errors,omitempty"`
}

type BatchLabelRequest struct {
	Items []CreateLabelRequest `json:"items"`
}

type BatchLabelJob struct {
	JobID     string `json:"jobId"`
	Status    string `json:"status"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
	Failed    int    `json:"failed"`
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

type WebhookSubscription struct {
	ID       string   `json:"id"`
	URL      string   `json:"url"`
	Events   []string `json:"events"`
	IsActive bool     `json:"isActive"`
}

type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
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
	ServiceCode           string `json:"serviceCode"`
	OriginPostalCode      string `json:"originPostalCode"`
	DestinationPostalCode string `json:"destinationPostalCode"`
	ShipDate              string `json:"shipDate"`
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
