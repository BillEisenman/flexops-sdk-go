package flexops_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go"
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

func computeHMAC(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
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

	_, _ = client.Shipping.GetRates(context.Background(), flexops.RateRequest{
		FromZip: "10001", ToZip: "90210", Weight: 16,
	})

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

	_, _ = client.Shipping.GetRates(context.Background(), flexops.RateRequest{
		FromZip: "10001", ToZip: "90210", Weight: 16,
	})

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
	_, _ = client.Shipping.GetRates(context.Background(), flexops.RateRequest{FromZip: "10001", ToZip: "90210"})

	if gotAuth != "Bearer new_jwt_token" {
		t.Errorf("expected Bearer new_jwt_token, got %q", gotAuth)
	}
}

// ---------------------------------------------------------------------------
// 4. Webhook signature verification — valid
// ---------------------------------------------------------------------------

func TestVerifySignature_Valid(t *testing.T) {
	payload := `{"event":"label.created","labelId":"lbl_123"}`
	secret := "whsec_test_secret"
	sig := computeHMAC(payload, secret)

	if !flexops.VerifySignature(payload, sig, secret) {
		t.Error("expected VerifySignature to return true for valid signature")
	}
}

// ---------------------------------------------------------------------------
// 5. Webhook signature verification — wrong signature
// ---------------------------------------------------------------------------

func TestVerifySignature_InvalidSignature(t *testing.T) {
	payload := `{"event":"label.created"}`
	secret := "whsec_test_secret"

	if flexops.VerifySignature(payload, "deadbeef", secret) {
		t.Error("expected VerifySignature to return false for wrong signature")
	}
}

// ---------------------------------------------------------------------------
// 6. Webhook signature verification — wrong secret
// ---------------------------------------------------------------------------

func TestVerifySignature_WrongSecret(t *testing.T) {
	payload := `{"event":"label.created"}`
	sig := computeHMAC(payload, "correct_secret")

	if flexops.VerifySignature(payload, sig, "wrong_secret") {
		t.Error("expected VerifySignature to return false for wrong secret")
	}
}

// ---------------------------------------------------------------------------
// 7. Shipping — GetRates parses response
// ---------------------------------------------------------------------------

func TestShipping_GetRates(t *testing.T) {
	rates := []flexops.ShippingRate{
		{Carrier: "ups", Service: "ground", Rate: 8.42, EstimatedDays: 3},
		{Carrier: "fedex", Service: "express", Rate: 18.90, EstimatedDays: 1},
	}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"success": true, "data": rates})
	}))

	resp, err := client.Shipping.GetRates(context.Background(), flexops.RateRequest{
		FromZip: "80202", ToZip: "10001", Weight: 32,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 rates, got %d", len(resp.Data))
	}
	if resp.Data[0].Carrier != "ups" {
		t.Errorf("expected first carrier = ups, got %q", resp.Data[0].Carrier)
	}
}

// ---------------------------------------------------------------------------
// 8. Shipping — CreateLabel parses label response
// ---------------------------------------------------------------------------

func TestShipping_CreateLabel(t *testing.T) {
	label := flexops.Label{
		LabelID: "lbl_abc123", TrackingNumber: "1Z999AA10123456784",
		Carrier: "ups", Rate: 8.42,
	}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"success": true, "data": label})
	}))

	resp, err := client.Shipping.CreateLabel(context.Background(), flexops.CreateLabelRequest{
		Carrier:     "ups",
		Service:     "ground",
		FromAddress: &flexops.Address{Name: "Sender", Street1: "123 Main St", City: "Denver", State: "CO", Zip: "80202", Country: "US"},
		ToAddress:   &flexops.Address{Name: "Recipient", Street1: "456 Park Ave", City: "New York", State: "NY", Zip: "10001", Country: "US"},
		Parcel:      &flexops.Parcel{Weight: 32},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.LabelID != "lbl_abc123" {
		t.Errorf("expected labelId = lbl_abc123, got %q", resp.Data.LabelID)
	}
	if resp.Data.TrackingNumber != "1Z999AA10123456784" {
		t.Errorf("unexpected tracking number: %q", resp.Data.TrackingNumber)
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
		writeJSON(w, 200, map[string]any{"success": true, "data": info})
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

	_, err := client.Shipping.GetRates(context.Background(), flexops.RateRequest{FromZip: "10001", ToZip: "90210"})
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

	_, err := client.Shipping.GetRates(context.Background(), flexops.RateRequest{FromZip: "10001", ToZip: "90210"})
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

func TestReturns_Create(t *testing.T) {
	rma := map[string]any{"rmaId": "rma_xyz", "status": "pending"}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(w, 200, map[string]any{"success": true, "data": rma})
	}))

	resp, err := client.Returns.Create(context.Background(), flexops.ReturnRequest{
		OrderNumber: "ord_123",
		Reason:      "damaged",
		Items:       []flexops.ReturnItem{{SKU: "SKU-001", Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Error("expected success = true")
	}
}

// ---------------------------------------------------------------------------
// 13. New services are present on the client
// ---------------------------------------------------------------------------

func TestClient_NewServicesPresent(t *testing.T) {
	client := flexops.NewClient(flexops.Config{APIKey: "key", WorkspaceID: "ws-123"})

	if client.Offsets == nil {
		t.Error("Offsets service is nil")
	}
	if client.HsCodes == nil {
		t.Error("HsCodes service is nil")
	}
	if client.RecurringShipments == nil {
		t.Error("RecurringShipments service is nil")
	}
	if client.EmailTemplates == nil {
		t.Error("EmailTemplates service is nil")
	}
	if client.Reports == nil {
		t.Error("Reports service is nil")
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
		writeJSON(w, 200, map[string]any{"success": true, "data": resp})
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

func TestOffsets_Offset(t *testing.T) {
	var gotPath string
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		writeJSON(w, 200, map[string]any{"success": true, "data": map[string]any{
			"labelId": "lbl_abc123", "offsetId": "off_xyz", "kgCo2e": 0.42,
		}})
	}))

	_, err := client.Offsets.Offset(context.Background(), "lbl_abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "/api/workspaces/ws-test-123/shipping/labels/lbl_abc123/offset"
	if gotPath != expected {
		t.Errorf("expected path %q, got %q", expected, gotPath)
	}
}
