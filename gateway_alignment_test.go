package flexops_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// Routes, bodies and enum encodings checked against the Gateway source (2026-10-19).

func sent(c recordedCall) string {
	body, _ := json.Marshal(c.body)
	return c.method + " " + c.path + c.query + " " + string(body)
}

func expectCalls(t *testing.T, calls []recordedCall, want ...string) {
	t.Helper()
	got := make([]string, len(calls))
	for i, c := range calls {
		got[i] = sent(c)
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestWallet_GetBalanceUsesMyBalanceWithWorkspaceHeader(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"hasWallet": true, "balanceCents": 12345, "balance": 123.45})
	balance, err := client.Wallet.GetBalance(context.Background())
	if err != nil || !balance.HasWallet || balance.BalanceCents != 12345 || balance.Balance != 123.45 {
		t.Fatalf("balance: %+v %v", balance, err)
	}
	expectCalls(t, *calls, "GET /api/wallet/my-balance null")
	if got := (*calls)[0].header.Get("X-Workspace-Id"); got != "ws-test-123" {
		t.Errorf("X-Workspace-Id = %q", got)
	}
}

func TestHsCodes_UseShippingRoutesWithoutWorkspace(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"hsCode": "6109.10", "description": "T-shirts"})
	ctx := context.Background()
	if _, err := client.HsCodes.Lookup(ctx, "6109.10"); err != nil {
		t.Fatal(err)
	}
	_, _ = client.HsCodes.Search(ctx, flexops.HsCodeSearchRequest{Query: "cotton shirt", DestinationCountry: "CA", MaxResults: 5})
	_, _ = client.HsCodes.Search(ctx, flexops.HsCodeSearchRequest{Query: "mug"})
	estimate, err := client.HsCodes.EstimateLandedCost(ctx, flexops.LandedCostRequest{DestinationCountry: "CA", ShippingCost: 12.5,
		Items: []flexops.CustomsItem{{Description: "Shirt", Quantity: 2, Value: 10, WeightOz: 6, OriginCountry: "US", SkuOrProductID: "SKU-1"}}})
	if err != nil {
		t.Fatal(err)
	}
	_ = estimate
	expectCalls(t, *calls,
		"GET /api/shipping/hs-codes/6109.10 null",
		"GET /api/shipping/hs-codes/search?destinationCountry=CA&maxResults=5&query=cotton+shirt null",
		"GET /api/shipping/hs-codes/search?query=mug null",
		`POST /api/shipping/landed-cost {"destinationCountry":"CA","items":[{"description":"Shirt","originCountry":"US","quantity":2,"skuOrProductId":"SKU-1","value":10,"weightOz":6}],"shippingCost":12.5}`,
	)
}

func TestInsurance_UsesGatewayRoutesAndBodies(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true, "data": map[string]any{"policyId": "p1", "premium": 1.25}})
	ctx := context.Background()
	if _, err := client.Insurance.GetQuote(ctx, flexops.InsuranceQuoteRequest{DeclaredValue: 100, CarrierCode: "USPS"}); err != nil {
		t.Fatal(err)
	}
	policy, err := client.Insurance.Purchase(ctx, flexops.PurchaseInsuranceRequest{TrackingNumber: "9400", CarrierCode: "USPS", DeclaredValue: 100, CoverageTier: "Standard"})
	if err != nil || policy.Data.PolicyID != "p1" {
		t.Fatalf("purchase: %+v %v", policy, err)
	}
	if _, err := client.Insurance.FileClaim(ctx, flexops.FileClaimRequest{PolicyID: "p1", ClaimReason: "Lost", ClaimAmount: 100}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Insurance.Void(ctx, "p1"); err != nil {
		t.Fatal(err)
	}
	_, _ = client.Insurance.GetProviders(ctx)
	base := "/api/workspaces/ws-test-123/insurance/"
	expectCalls(t, *calls,
		"POST "+base+`quote {"carrierCode":"USPS","declaredValue":100}`,
		"POST "+base+`policies {"carrierCode":"USPS","coverageTier":"Standard","declaredValue":100,"trackingNumber":"9400"}`,
		"POST "+base+`claims {"claimAmount":100,"claimReason":"Lost","policyId":"p1"}`,
		"DELETE "+base+"policies/p1 null",
		"GET "+base+"providers null",
	)
}

func TestReturns_LabelReceiveAndList(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	ctx := context.Background()
	if _, err := client.Returns.GenerateLabelWithOptions(ctx, "RMA-1", flexops.ReturnLabelOptions{CarrierCode: "usps", RateIndicator: "SP", WeightOz: 8, LengthIn: 6, WidthIn: 4, HeightIn: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Returns.MarkReceived(ctx, "7", flexops.ReceiveReturnRequest{RmaNumber: "RMA-1", WarehouseID: 2,
		InspectedItems: []flexops.InspectedItem{{PartNumber: "SKU-1", Quantity: 1, Disposition: flexops.ReturnDispositionDispose}, {PartNumber: "SKU-2", Quantity: 1, Disposition: flexops.ReturnDispositionRestock}}}); err != nil {
		t.Fatal(err)
	}
	_, _ = client.Returns.List(ctx, flexops.ReturnListOptions{StartDate: "2026-10-01", EndDate: "2026-10-19", Page: 2, PageSize: 50})
	_, _ = client.Returns.List(ctx, flexops.ReturnListOptions{})
	base := "/api/workspaces/ws-test-123/returns"
	expectCalls(t, *calls,
		"POST "+base+`/RMA-1/label {"carrierCode":"usps","heightIn":2,"lengthIn":6,"rateIndicator":"SP","weightOz":8,"widthIn":4}`,
		"POST "+base+`/7/receive {"inspectedItems":[{"disposition":1,"partNumber":"SKU-1","quantity":1},{"disposition":0,"partNumber":"SKU-2","quantity":1}],"rmaNumber":"RMA-1","warehouseId":2}`,
		"GET "+base+"?endDate=2026-10-19&page=2&pageSize=50&startDate=2026-10-01 null",
		"GET "+base+" null",
	)
}

func TestRules_UseGatewayBodies(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	ctx := context.Background()
	group := 1
	if _, err := client.Rules.Create(ctx, flexops.CreateShippingRuleRequest{Name: "Heavy", Priority: 5, IsEnabled: true,
		Conditions: []flexops.RuleCondition{{Field: "weightOz", Operator: "gt", Value: "160", Group: &group}},
		Actions:    []flexops.RuleAction{{Type: "setCarrier", Value: "UPS"}}}); err != nil {
		t.Fatal(err)
	}
	enabled := false
	if _, err := client.Rules.Update(ctx, "r1", flexops.UpdateShippingRuleRequest{IsEnabled: &enabled}); err != nil {
		t.Fatal(err)
	}
	_, _ = client.Rules.Reorder(ctx, []flexops.RulePriority{{RuleID: "r1", Priority: 1}, {RuleID: "r2", Priority: 2}})
	if _, err := client.Rules.Test(ctx, flexops.TestShippingRuleRequest{WeightOz: 200, DeclaredValue: 50, DestinationCountry: "US"}); err != nil {
		t.Fatal(err)
	}
	base := "/api/workspaces/ws-test-123/shipping-rules"
	expectCalls(t, *calls,
		"POST "+base+` {"actions":[{"type":"setCarrier","value":"UPS"}],"conditions":[{"field":"weightOz","group":1,"operator":"gt","value":"160"}],"isEnabled":true,"name":"Heavy","priority":5}`,
		"PUT "+base+`/r1 {"isEnabled":false}`,
		"POST "+base+`/reorder {"rules":[{"priority":1,"ruleId":"r1"},{"priority":2,"ruleId":"r2"}]}`,
		"POST "+base+`/test {"declaredValue":50,"destinationCountry":"US","weightOz":200}`,
	)
}

func TestPickups_UseGatewayRoutesAndBody(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	ctx := context.Background()
	if _, err := client.Pickups.Schedule(ctx, flexops.SchedulePickupRequest{CarrierCode: "usps", PickupDate: "2026-10-20",
		Street: "1 Main St", City: "Denver", State: "CO", PostalCode: "80202", PackageCount: 3, ContactPhone: "5551234567"}); err != nil {
		t.Fatal(err)
	}
	_, _ = client.Pickups.List(ctx)
	if _, err := client.Pickups.Get(ctx, "PU-1", "ups"); err != nil {
		t.Fatal(err)
	}
	if err := client.Pickups.Cancel(ctx, "PU-1", ""); err != nil {
		t.Fatal(err)
	}
	base := "/api/workspaces/ws-test-123/pickups"
	expectCalls(t, *calls,
		"POST "+base+` {"carrierCode":"usps","city":"Denver","contactPhone":"5551234567","packageCount":3,"pickupDate":"2026-10-20","postalCode":"80202","state":"CO","street":"1 Main St"}`,
		"GET "+base+"/scheduled null",
		"GET "+base+"/PU-1?carrierCode=ups null",
		"DELETE "+base+"/PU-1 null",
	)
}

func TestWorkspaces_MembersUseNumericRoles(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true, "data": map[string]any{
		"members":        []any{map[string]any{"userId": "u1", "email": "a@example.com", "role": 1, "status": 0}},
		"pendingInvites": []any{map[string]any{"email": "b@example.com", "role": 3}},
		"userLimit":      5, "activeCount": 1, "canInviteMore": true,
	}})
	ctx := context.Background()
	members, err := client.Workspaces.ListMembers(ctx)
	if err != nil || len(members.Data.Members) != 1 || members.Data.Members[0].Role != flexops.WorkspaceRoleAdmin ||
		members.Data.PendingInvites[0].Role != flexops.WorkspaceRoleMerchant || !members.Data.CanInviteMore {
		t.Fatalf("members: %+v %v", members, err)
	}
	if _, err := client.Workspaces.InviteMember(ctx, "c@example.com", flexops.WorkspaceRoleMember); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Workspaces.UpdateMemberRole(ctx, "u1", flexops.WorkspaceRoleOwner); err != nil {
		t.Fatal(err)
	}
	base := "/api/workspaces/ws-test-123/members"
	expectCalls(t, *calls,
		"GET "+base+" null",
		"POST "+base+`/invite {"email":"c@example.com","role":2}`,
		"PUT "+base+`/u1/role {"newRole":0}`,
	)
	for want, role := range []flexops.WorkspaceRole{flexops.WorkspaceRoleOwner, flexops.WorkspaceRoleAdmin, flexops.WorkspaceRoleMember, flexops.WorkspaceRoleMerchant} {
		if int(role) != want {
			t.Errorf("role %d = %d", want, role)
		}
	}
}

func TestWorkspaces_DecodeGatewayShapes(t *testing.T) {
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/workspaces" {
			writeJSON(w, 200, map[string]any{"success": true, "data": []any{map[string]any{"id": "ws1", "name": "Acme", "myRole": 1, "memberCount": 3,
				"entitlements": []any{map[string]any{"productType": "ship", "planId": "pro", "isPrimary": true, "subscriptionStatus": "active"}}}}})
			return
		}
		writeJSON(w, 200, map[string]any{"success": true, "data": map[string]any{"id": "ws-test-123", "myRole": 0, "userLimit": 5, "canManageBilling": true, "subscriptionStatus": "trialing"}})
	}))
	list, err := client.Workspaces.List(context.Background())
	if err != nil || list.Data[0].MyRole != flexops.WorkspaceRoleAdmin || list.Data[0].MemberCount != 3 || !list.Data[0].Entitlements[0].IsPrimary {
		t.Fatalf("list: %+v %v", list, err)
	}
	detail, err := client.Workspaces.Get(context.Background(), "")
	if err != nil || detail.Data.ID != "ws-test-123" || detail.Data.MyRole != flexops.WorkspaceRoleOwner || detail.Data.UserLimit != 5 || detail.Data.SubscriptionStatus != "trialing" {
		t.Fatalf("detail: %+v %v", detail, err)
	}
}

func TestWebhooks_TestPostsToSubscription(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true, "data": "Test event sent."})
	result, err := client.Webhooks.Test(context.Background(), "wh-1")
	if err != nil || result.Data != "Test event sent." {
		t.Fatalf("test: %+v %v", result, err)
	}
	expectCalls(t, *calls, "POST /api/workspaces/ws-test-123/webhooks/wh-1/test null")
}

func TestAuth_LoginAndRefreshUseGatewayTokens(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true, "data": map[string]any{"token": "jwt-1", "refreshToken": "rt-1", "userId": "u1"}})
	ctx := context.Background()
	login, err := client.Auth.Login(ctx, "a@example.com", "pw")
	if err != nil || login.Data.Token != "jwt-1" || login.Data.RefreshToken != "rt-1" {
		t.Fatalf("login: %+v %v", login, err)
	}
	if _, err := client.Auth.RefreshToken(ctx, "rt-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Auth.RefreshToken(ctx, "rt-2", "jwt-explicit"); err != nil {
		t.Fatal(err)
	}
	expectCalls(t, *calls,
		`POST /api/account/login {"email":"a@example.com","password":"pw"}`,
		`POST /api/account/refresh-token {"refreshToken":"rt-1","token":"jwt-1"}`,
		`POST /api/account/refresh-token {"refreshToken":"rt-2","token":"jwt-explicit"}`,
	)
	// Login switches the client from the API key to the returned bearer token.
	if got := (*calls)[1].header.Get("Authorization"); got != "Bearer jwt-1" || (*calls)[1].header.Get("X-Api-Key") != "" {
		t.Errorf("after login sent Authorization %q", got)
	}
}

func TestAuth_RefreshWithoutAccessTokenFailsBeforeRequest(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	if _, err := client.Auth.RefreshToken(context.Background(), "rt-1"); err == nil || !strings.Contains(err.Error(), "access token") {
		t.Fatalf("expected an access token error, got %v", err)
	}
	if len(*calls) != 0 {
		t.Errorf("sent %d requests", len(*calls))
	}
}
