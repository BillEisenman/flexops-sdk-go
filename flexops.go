// ***********************************************************************
// Package          : flexops-sdk-go
// Author           : FlexOps, LLC
// Created          : 2026-03-08
//
// Copyright (c) 2021-2026 by FlexOps, LLC. All rights reserved.
// ***********************************************************************

// Package flexops provides the official Go SDK for the FlexOps multi-carrier
// shipping platform API.
package flexops

import "time"

// Config configures the FlexOps client.
type Config struct {
	BaseURL     string
	APIKey      string
	AccessToken string
	WorkspaceID string
	Timeout     time.Duration
	// MaxRetries is how many times a retryable call is sent again (default 3).
	MaxRetries int
}

// Client is the main entry point for the FlexOps SDK.
type Client struct {
	WorkspaceID string
	Auth        *AuthService
	Workspaces  *WorkspacesService
	Shipping    *ShippingService
	Carriers    *CarriersService
	// Webhooks needs a signed-in user's access token; API keys get 403 SIGNED_IN_USER_REQUIRED.
	Webhooks      *WebhooksService
	Wallet        *WalletService
	Insurance     *InsuranceService
	Returns       *ReturnsService
	Analytics     *AnalyticsService
	Orders        *OrdersService
	Inventory     *InventoryService
	Pickups       *PickupsService
	ScanForms     *ScanFormsService
	Rules         *RulesService
	CarbonOffsets *CarbonOffsetsService
	HsCodes       *HsCodesService
	// Writes on these three need a user token or an API key with the resource's write scope.
	RecurringShipments *RecurringShipmentsService
	ReportSchedules    *ReportSchedulesService
	EmailTemplates     *EmailTemplatesService
	http               *httpClient
}

// NewClient creates a new FlexOps API client.
func NewClient(cfg Config) *Client {
	c := &Client{
		WorkspaceID: cfg.WorkspaceID,
		http:        newHTTPClient(cfg),
	}
	c.Auth = &AuthService{c: c}
	c.Workspaces = &WorkspacesService{c: c}
	c.Shipping = &ShippingService{c: c}
	c.Carriers = newCarriersService(c)
	c.Webhooks = &WebhooksService{c: c}
	c.Wallet = &WalletService{c: c}
	c.Insurance = &InsuranceService{c: c}
	c.Returns = &ReturnsService{c: c}
	c.Analytics = &AnalyticsService{c: c}
	c.Orders = &OrdersService{c: c}
	c.Inventory = &InventoryService{c: c}
	c.Pickups = &PickupsService{c: c}
	c.ScanForms = &ScanFormsService{c: c}
	c.Rules = &RulesService{c: c}
	c.CarbonOffsets = &CarbonOffsetsService{c: c}
	c.HsCodes = &HsCodesService{c: c}
	c.RecurringShipments = &RecurringShipmentsService{c: c}
	c.ReportSchedules = &ReportSchedulesService{c: c}
	c.EmailTemplates = &EmailTemplatesService{c: c}
	return c
}

// SetAccessToken sets the JWT access token.
func (c *Client) SetAccessToken(token string) { c.http.setAccessToken(token) }

// SetAPIKey sets the API key.
func (c *Client) SetAPIKey(key string) { c.http.setAPIKey(key) }

func (c *Client) wsPath(suffix string) string {
	return "/api/workspaces/" + c.WorkspaceID + "/" + suffix
}
