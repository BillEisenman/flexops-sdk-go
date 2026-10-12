package flexops_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v2"
)

// The Gateway ShippingController lives at api/Shipping, resolves the workspace from the
// API key and returns raw DTOs (no success/data envelope).
func TestShipping_UsesGatewayRoutes(t *testing.T) {
	var seen []string
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.EscapedPath()+query(r))
		switch {
		case strings.HasPrefix(r.URL.Path, "/api/shipping/track/"):
			writeJSON(w, 200, map[string]any{"trackingNumber": "9400 1", "carrierCode": "USPS", "status": "Delivered",
				"estimatedDeliveryDate": "2026-10-14T00:00:00Z", "events": []any{map[string]any{"description": "Delivered"}}})
		case r.URL.Path == "/api/shipping/carriers":
			writeJSON(w, 200, []any{map[string]any{"carrierCode": "USPS"}})
		case r.URL.Path == "/api/shipping/recommendations":
			writeJSON(w, 200, map[string]any{"lane": "802->100", "recommendations": []any{map[string]any{"carrierCode": "UPS"}}})
		case r.URL.Path == "/api/shipping/predictions/delivery":
			writeJSON(w, 200, map[string]any{"carrierCode": "USPS", "predictedTransitDays": 3})
		case r.URL.Path == "/api/shipping/savings":
			writeJSON(w, 200, map[string]any{"lanesAnalyzed": 3})
		default:
			writeJSON(w, 200, map[string]any{"success": true})
		}
	}))
	ctx := context.Background()

	tracking, err := client.Shipping.Track(ctx, "9400 1")
	if err != nil || tracking.Data.Status != "Delivered" || tracking.Data.CarrierCode != "USPS" || len(tracking.Data.Events) != 1 {
		t.Fatalf("track: %+v %v", tracking, err)
	}
	if err := client.Shipping.CancelLabel(ctx, "lbl-1", "USPS"); err != nil {
		t.Fatal(err)
	}
	carriers, err := client.Shipping.GetCarriers(ctx)
	if err != nil || carriers.Data == nil {
		t.Fatalf("carriers: %+v %v", carriers, err)
	}
	recommendations, err := client.Shipping.GetRecommendations(ctx, flexops.CarrierRecommendationRequest{OriginPostalCode: "80202", DestinationPostalCode: "10001"})
	if err != nil || recommendations.Data.Lane != "802->100" || len(recommendations.Data.Recommendations) != 1 {
		t.Fatalf("recommendations: %+v %v", recommendations, err)
	}
	prediction, err := client.Shipping.PredictDelivery(ctx, flexops.DeliveryPredictionRequest{CarrierCode: "USPS", OriginPostalCode: "80202", DestinationPostalCode: "10001"})
	if err != nil || prediction.Data.PredictedTransitDays != 3 {
		t.Fatalf("prediction: %+v %v", prediction, err)
	}
	savings, err := client.Shipping.GetSavings(ctx)
	if err != nil || savings.Data == nil {
		t.Fatalf("savings: %+v %v", savings, err)
	}

	want := []string{
		"GET /api/shipping/track/9400%201",
		"DELETE /api/shipping/labels/lbl-1?carrierCode=USPS",
		"GET /api/shipping/carriers",
		"POST /api/shipping/recommendations",
		"POST /api/shipping/predictions/delivery",
		"GET /api/shipping/savings",
	}
	if strings.Join(seen, "\n") != strings.Join(want, "\n") {
		t.Errorf("requests:\n%s\nwant:\n%s", strings.Join(seen, "\n"), strings.Join(want, "\n"))
	}
}

func TestShipping_CancelLabelNeedsCarrierCode(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++ }))
	if err := client.Shipping.CancelLabel(context.Background(), "lbl-1"); err == nil || !strings.Contains(err.Error(), "carrier code") {
		t.Fatalf("expected a carrier code error, got %v", err)
	}
	if calls != 0 {
		t.Errorf("sent %d requests without a carrier code", calls)
	}
}

func TestShipping_PredictDeliveryOmitsUnsetOptionalFields(t *testing.T) {
	var body map[string]any
	client, calls := recordingClient(t, map[string]any{"carrierCode": "USPS"})
	if _, err := client.Shipping.PredictDelivery(context.Background(), flexops.DeliveryPredictionRequest{CarrierCode: "USPS", OriginPostalCode: "80202", DestinationPostalCode: "10001"}); err != nil {
		t.Fatal(err)
	}
	body = (*calls)[0].body
	for _, field := range []string{"shipDate", "serviceCode"} {
		if _, ok := body[field]; ok {
			t.Errorf("sent unset %s: %v (the Gateway binds shipDate as DateTime?, so \"\" is refused)", field, body)
		}
	}
}

func query(r *http.Request) string {
	if r.URL.RawQuery == "" {
		return ""
	}
	return "?" + r.URL.RawQuery
}
