package flexops_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestClient(t *testing.T, handler http.Handler) (*flexops.Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := flexops.NewClient(flexops.Config{
		APIKey:      "sk_test_key",
		WorkspaceID: "ws-test-123",
		BaseURL:     srv.URL,
	})
	return client, srv
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func rateRequest() flexops.RateRequest {
	return flexops.RateRequest{
		Origin:      flexops.ShippingAddress{AddressLine1: "123 Main St", City: "New York", StateProvince: "NY", PostalCode: "10001"},
		Destination: flexops.ShippingAddress{AddressLine1: "456 Oak Ave", City: "Los Angeles", StateProvince: "CA", PostalCode: "90210"},
		Package:     flexops.ShippingPackage{Weight: 16, WeightUnit: "oz"},
	}
}

// ---------------------------------------------------------------------------
// 1. Client initialisation — API key auth
// ---------------------------------------------------------------------------

func TestClient_APIKeyAuth(t *testing.T) {
	var gotHeader string
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Api-Key")
		writeJSON(w, 200, map[string]any{"success": true, "data": []any{}})
	}))

	_, _ = client.Shipping.GetRates(context.Background(), rateRequest())

	if gotHeader != "sk_test_key" {
		t.Errorf("expected X-Api-Key header = sk_test_key, got %q", gotHeader)
	}
}

// ---------------------------------------------------------------------------
// 2. Client initialisation — Bearer token auth
// ---------------------------------------------------------------------------

func TestClient_BearerTokenAuth(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("Authorization")
		writeJSON(w, 200, map[string]any{"success": true, "data": []any{}})
	}))
	defer srv.Close()

	client := flexops.NewClient(flexops.Config{
		AccessToken: "jwt_test_token",
		WorkspaceID: "ws-test-123",
		BaseURL:     srv.URL,
	})

	_, _ = client.Shipping.GetRates(context.Background(), rateRequest())

	if gotHeader != "Bearer jwt_test_token" {
		t.Errorf("expected Authorization = 'Bearer jwt_test_token', got %q", gotHeader)
	}
}

// ---------------------------------------------------------------------------
// 3. SetAccessToken switches auth mode
// ---------------------------------------------------------------------------

func TestClient_SetAccessToken(t *testing.T) {
	var gotAuth string
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		writeJSON(w, 200, map[string]any{"success": true, "data": []any{}})
	}))

	client.SetAccessToken("new_jwt_token")
	_, _ = client.Shipping.GetRates(context.Background(), rateRequest())

	if gotAuth != "Bearer new_jwt_token" {
		t.Errorf("expected Bearer new_jwt_token, got %q", gotAuth)
	}
}

// ---------------------------------------------------------------------------
// 4-6. Webhook signature verification: see webhooks_test.go
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// 7. Shipping — GetRates parses response
// ---------------------------------------------------------------------------

func TestShipping_GetRates(t *testing.T) {
	rates := []flexops.ShippingRate{
		{Carrier: "ups", Service: "ground", Rate: 8.42, EstimatedDays: 3},
		{Carrier: "fedex", Service: "express", Rate: 18.90, EstimatedDays: 1},
	}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/shipping/rates" {
			t.Errorf("expected /api/shipping/rates, got %q", r.URL.Path)
		}
		var request map[string]any
		_ = json.NewDecoder(r.Body).Decode(&request)
		if request["origin"] == nil || request["destination"] == nil || request["package"] == nil {
			t.Errorf("expected canonical origin/destination/package body, got %#v", request)
		}
		writeJSON(w, 200, map[string]any{"currency": "USD", "rates": rates})
	}))

	resp, err := client.Shipping.GetRates(context.Background(), rateRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Rates) != 2 {
		t.Fatalf("expected 2 rates, got %d", len(resp.Rates))
	}
	if resp.Rates[0].Carrier != "ups" {
		t.Errorf("expected first carrier = ups, got %q", resp.Rates[0].Carrier)
	}
}

// ---------------------------------------------------------------------------
// 8. Shipping — CreateLabel parses label response
// ---------------------------------------------------------------------------

func TestShipping_CreateLabel(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["carrierCode"] != "USPS" || body["maximumPostageAmount"] != 10.25 || body["origin"] == nil || body["package"] == nil {
			t.Errorf("bad body: %v", body)
		}
		if calls == 1 || calls == 4 {
			// Every write carries a key; a preview gets a generated one, never the purchase key.
			if key := r.Header.Get("Idempotency-Key"); body["confirmationToken"] != nil || key == "" || key == "purchase-001" {
				t.Errorf("preview contains confirmation (key %q)", key)
			}
			writeJSON(w, 200, map[string]any{"status": "Preview", "quotedPostageAmount": 8.5, "maximumPostageAmount": 10.25, "confirmationToken": "approval", "currency": "USD", "expiresAt": "2026-09-19T00:05:00Z"})
		} else {
			if r.Header.Get("Idempotency-Key") != "purchase-001" || body["confirmationToken"] != "approval" {
				t.Error("missing purchase approval")
			}
			if calls == 2 {
				writeJSON(w, 503, map[string]string{"message": "temporary"})
				return
			}
			writeJSON(w, 201, map[string]any{"labelId": "lbl-001", "carrierCode": "USPS", "trackingNumber": "tracking", "rate": 8.5, "currency": "USD"})
		}
	}))
	request := flexops.CreateLabelRequest{CarrierCode: "USPS", ServiceCode: "GROUND_ADVANTAGE", MaximumPostageAmount: 10.25,
		Origin: flexops.ShippingAddress{AddressLine1: "1 St", CountryCode: "US"}, Destination: flexops.ShippingAddress{AddressLine1: "2 St", CountryCode: "US"}, Package: flexops.ShippingPackage{Weight: 16}}
	preview, err := client.Shipping.CreateLabel(context.Background(), request)
	if err != nil || preview.Status != "Preview" || preview.ConfirmationToken != "approval" || preview.QuotedPostageAmount != 8.5 || calls != 1 {
		t.Fatalf("bad preview: %+v %v", preview, err)
	}
	request.ConfirmationToken = preview.ConfirmationToken
	label, err := client.Shipping.CreateLabel(context.Background(), request, "purchase-001")
	if err != nil || label.LabelID != "lbl-001" || label.CarrierCode != "USPS" || label.Status == "Preview" || calls != 3 {
		t.Fatalf("bad purchase: %+v %v", label, err)
	}
	request.ConfirmationToken = ""
	preview, err = client.Shipping.CreateLabel(context.Background(), request)
	if err != nil || preview.Status != "Preview" || calls != 4 {
		t.Fatalf("bad subsequent preview: %+v %v", preview, err)
	}
}

func TestShipping_LabelApprovalErrors(t *testing.T) {
	for _, tc := range []struct {
		status int
		code   string
	}{{400, "ApprovalRequired"}, {409, "ApprovalExpired"}} {
		t.Run(tc.code, func(t *testing.T) {
			calls := 0
			client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				writeJSON(w, tc.status, map[string]string{"errorCode": tc.code, "message": tc.code})
			}))
			_, err := client.Shipping.CreateLabel(context.Background(), flexops.CreateLabelRequest{})
			apiErr, ok := err.(*flexops.FlexOpsError)
			if !ok || apiErr.Code != tc.code || apiErr.StatusCode != tc.status || calls != 1 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 9. Shipping — Track parses tracking response
// ---------------------------------------------------------------------------

func TestShipping_Track(t *testing.T) {
	info := flexops.TrackingInfo{
		TrackingNumber: "1Z999AA10123456784",
		Carrier:        "ups",
		Status:         "delivered",
		Events: []flexops.TrackingEvent{
			{Status: "delivered", Description: "Package delivered"},
		},
	}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, info) // The Gateway returns the tracking DTO without an envelope.
	}))

	resp, err := client.Shipping.Track(context.Background(), "1Z999AA10123456784")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Status != "delivered" {
		t.Errorf("expected status = delivered, got %q", resp.Data.Status)
	}
}

// ---------------------------------------------------------------------------
// 10. Retry — retries on 429 and succeeds on final attempt
// ---------------------------------------------------------------------------

func TestRetry_429_EventualSuccess(t *testing.T) {
	var callCount int32
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			return
		}
		writeJSON(w, 200, map[string]any{"success": true, "data": []any{}})
	}))

	_, err := client.Shipping.GetRates(context.Background(), rateRequest())
	if err != nil {
		t.Fatalf("expected eventual success after retries, got: %v", err)
	}
	if callCount < 3 {
		t.Errorf("expected at least 3 calls (2 retries), got %d", callCount)
	}
}

// ---------------------------------------------------------------------------
// 11. Error — 401 returns AuthError
// ---------------------------------------------------------------------------

func TestError_401_ReturnsAuthError(t *testing.T) {
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	}))

	_, err := client.Shipping.GetRates(context.Background(), rateRequest())
	if err == nil {
		t.Fatal("expected AuthError, got nil")
	}
	var authErr *flexops.AuthError
	if !isAuthError(err, &authErr) {
		t.Errorf("expected *flexops.AuthError, got %T: %v", err, err)
	}
}

// isAuthError checks if err is an *AuthError via type assertion.
func isAuthError(err error, target **flexops.AuthError) bool {
	e, ok := err.(*flexops.AuthError)
	if ok {
		*target = e
	}
	return ok
}

// ---------------------------------------------------------------------------
// 12. Returns — Create CRUD flow
// ---------------------------------------------------------------------------

// The Gateway opens an RMA at POST returns/rma with a numeric orderId and answers 201
// with the enveloped RmaDto.
func TestReturns_Create(t *testing.T) {
	var gotPath string
	var body map[string]any
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&body)
		writeJSON(w, 201, map[string]any{"success": true, "data": map[string]any{"rmaNumber": "RMA-1", "rmaId": 7, "orderId": 42, "status": "Pending"}})
	}))

	resp, err := client.Returns.Create(context.Background(), flexops.ReturnRequest{
		OrderID: 42,
		Reason:  "damaged",
		Items:   []flexops.ReturnItem{{SKU: "SKU-001", Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success || resp.Data.RmaNumber != "RMA-1" || resp.Data.RmaID != 7 {
		t.Errorf("bad RMA: %+v", resp)
	}
	if gotPath != "/api/workspaces/ws-test-123/returns/rma" {
		t.Errorf("posted to %s", gotPath)
	}
	if body["orderId"] != 42.0 || body["reason"] != "damaged" {
		t.Errorf("bad body: %v", body)
	}
	for _, unset := range []string{"orderNumber", "notes", "customerEmail", "reasonId"} {
		if _, ok := body[unset]; ok {
			t.Errorf("sent unset %s: %v", unset, body)
		}
	}
}

// ---------------------------------------------------------------------------
// 13. New services are present on the client
// ---------------------------------------------------------------------------

func TestClient_NewServicesPresent(t *testing.T) {
	client := flexops.NewClient(flexops.Config{APIKey: "key", WorkspaceID: "ws-123"})

	if client.CarbonOffsets == nil {
		t.Error("CarbonOffsets service is nil")
	}
	if client.HsCodes == nil {
		t.Error("HsCodes service is nil")
	}
}

// ---------------------------------------------------------------------------
// 14. AI Recommendations endpoint
// ---------------------------------------------------------------------------

func TestShipping_GetRecommendations(t *testing.T) {
	resp := flexops.CarrierRecommendationResponse{
		Lane:       "80202→10001",
		SampleSize: 100,
		Recommendations: []flexops.CarrierRecommendation{
			{CarrierCode: "ups", Score: 0.91, OnTimePercent: 97.2},
		},
	}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, resp) // The Gateway returns the recommendation DTO without an envelope.
	}))

	result, err := client.Shipping.GetRecommendations(context.Background(), flexops.CarrierRecommendationRequest{
		OriginPostalCode:      "80202",
		DestinationPostalCode: "10001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data.Lane != "80202→10001" {
		t.Errorf("expected lane = 80202→10001, got %q", result.Data.Lane)
	}
	if len(result.Data.Recommendations) != 1 {
		t.Errorf("expected 1 recommendation, got %d", len(result.Data.Recommendations))
	}
}

// ---------------------------------------------------------------------------
// 15. Carbon Offsets endpoint
// ---------------------------------------------------------------------------

func TestCarbonOffsets_UseGatewayRoutes(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true, "data": map[string]any{"offsetId": "off-1", "carbonKg": 0.42}})
	ctx := context.Background()
	if _, err := client.CarbonOffsets.Quote(ctx, flexops.CarbonOffsetQuoteRequest{WeightLbs: 2, DistanceMiles: 300}); err != nil {
		t.Fatal(err)
	}
	purchase, err := client.CarbonOffsets.Purchase(ctx, flexops.CarbonOffsetPurchaseRequest{TrackingNumber: "9400", WeightLbs: 2, DistanceMiles: 300})
	if err != nil || purchase.Data.OffsetID != "off-1" || purchase.Data.CarbonKg != 0.42 {
		t.Fatalf("purchase: %+v %v", purchase, err)
	}
	// The fake answers every call with one object, so only the list request is checked here.
	_, _ = client.CarbonOffsets.List(ctx, 2, 10)
	if _, err := client.CarbonOffsets.Get(ctx, "off-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CarbonOffsets.Summary(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CarbonOffsets.GetSettings(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CarbonOffsets.UpdateSettings(ctx, flexops.CarbonOffsetSettings{AutoOffsetEnabled: false}); err != nil {
		t.Fatal(err)
	}

	base := "/api/workspaces/ws-test-123/carbon-offsets"
	want := []struct {
		call string
		body string
	}{
		{"POST " + base + "/quote", `{"distanceMiles":300,"weightLbs":2}`},
		{"POST " + base + "/purchase", `{"distanceMiles":300,"trackingNumber":"9400","weightLbs":2}`},
		{"GET " + base + "?page=2&pageSize=10", "null"},
		{"GET " + base + "/off-1", "null"},
		{"GET " + base + "/summary", "null"},
		{"GET " + base + "/settings", "null"},
		{"PUT " + base + "/settings", `{"autoOffsetEnabled":false}`},
	}
	if len(*calls) != len(want) {
		t.Fatalf("sent %d requests, want %d", len(*calls), len(want))
	}
	for i, w := range want {
		got, _ := json.Marshal((*calls)[i].body)
		if call := (*calls)[i].method + " " + (*calls)[i].path + (*calls)[i].query; call != w.call || string(got) != w.body {
			t.Errorf("call %d: %s %s, want %s %s", i, call, got, w.call, w.body)
		}
	}
}
