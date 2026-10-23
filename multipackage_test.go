package flexops_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// Multi-package labels (POST /api/shipping/multi-package) and the RMA reads, checked against
// the Gateway MultiPackageController, MultiPackageShipmentDTOs and ReturnsController.

func TestShipping_CreateMultiPackage_SendsGatewayBodyAndKey(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{
		"success": true,
		"data": map[string]any{"orderId": 15563627, "carrierCode": "USPS", "totalPackages": 2, "successfulPackages": 2,
			"failedPackages": 0, "totalCost": 18.94, "labels": []map[string]any{
				{"packageIndex": 0, "trackingNumber": "9400-1", "success": true},
				{"packageIndex": 1, "trackingNumber": "9400-2", "success": true},
			}},
	})
	insured := 50.0
	resp, err := client.Shipping.CreateMultiPackage(context.Background(), flexops.MultiPackageShipmentRequest{
		OrderID: 15563627, CarrierCode: "USPS", ServiceCode: "PRIORITY", HazmatClassification: flexops.HazmatNone.Ptr(),
		Packages: []flexops.MultiPackageDetail{
			{WeightInOunces: 16, LengthInInches: 10, WidthInInches: 8, HeightInInches: 4, PackageReference: "Box 1 of 2"},
			{WeightInOunces: 8, LengthInInches: 6, WidthInInches: 6, HeightInInches: 6, InsuranceCoverageRequested: true, InsuredValue: &insured},
		},
	}, "order-15563627")
	if err != nil {
		t.Fatal(err)
	}
	expectCalls(t, *calls,
		`POST /api/shipping/multi-package {"carrierCode":"USPS","hazmatClassification":0,"orderId":15563627,"packages":[`+
			`{"heightInInches":4,"lengthInInches":10,"packageReference":"Box 1 of 2","weightInOunces":16,"widthInInches":8},`+
			`{"heightInInches":6,"insuranceCoverageRequested":true,"insuredValue":50,"lengthInInches":6,"weightInOunces":8,"widthInInches":6}],"serviceCode":"PRIORITY"}`,
	)
	if got := (*calls)[0].header.Get("Idempotency-Key"); got != "order-15563627" {
		t.Errorf("Idempotency-Key = %q, want the caller's key", got)
	}
	if resp.Data.SuccessfulPackages != 2 || len(resp.Data.Labels) != 2 || resp.Data.Labels[1].TrackingNumber != "9400-2" {
		t.Errorf("response = %+v", resp.Data)
	}
}

func TestShipping_CreateMultiPackage_OmitsUnsetHazmat(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	_, _ = client.Shipping.CreateMultiPackage(context.Background(), flexops.MultiPackageShipmentRequest{
		OrderID: 1, CarrierCode: "USPS", Packages: []flexops.MultiPackageDetail{{WeightInOunces: 1, LengthInInches: 1, WidthInInches: 1, HeightInInches: 1}},
	})
	if got := sentHazmat(t, (*calls)[0].body); got != "absent" {
		t.Errorf("hazmatClassification = %v, want absent", got)
	}
	if (*calls)[0].header.Get("Idempotency-Key") == "" {
		t.Error("a generated Idempotency-Key was not sent")
	}
}

func TestShipping_CreateMultiPackage_UnknownOutcomeKeepsResultsAndIsNotRetried(t *testing.T) {
	var hits int32
	data := map[string]any{"orderId": 1, "carrierCode": "USPS", "totalPackages": 1, "successfulPackages": 0, "failedPackages": 1,
		"totalCost": 0, "labels": []map[string]any{{"packageIndex": 0, "success": false, "errorCode": "OutcomeUnknown"}}}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		writeJSON(w, http.StatusConflict, map[string]any{"success": false, "message": "No label was confirmed",
			"errorCode": "OutcomeUnknown", "errors": []string{"Package 0: unknown"}, "data": data})
	}))
	_, err := client.Shipping.CreateMultiPackage(context.Background(), flexops.MultiPackageShipmentRequest{
		OrderID: 1, CarrierCode: "USPS", Packages: []flexops.MultiPackageDetail{{WeightInOunces: 1, LengthInInches: 1, WidthInInches: 1, HeightInInches: 1}},
	})
	var apiErr *flexops.FlexOpsError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusConflict || apiErr.Code != "OutcomeUnknown" {
		t.Fatalf("err = %v, want 409 OutcomeUnknown", err)
	}
	var results flexops.MultiPackageShipmentResponse
	if err := json.Unmarshal(apiErr.Data, &results); err != nil {
		t.Fatalf("Data = %s: %v", apiErr.Data, err)
	}
	if len(results.Labels) != 1 || results.Labels[0].ErrorCode != "OutcomeUnknown" || results.FailedPackages != 1 {
		t.Errorf("per-box results = %+v", results)
	}
	if n := atomic.LoadInt32(&hits); n != 1 {
		t.Errorf("sent %d times, want 1 (an unknown outcome must not be retried)", n)
	}
}

func TestShipping_CreateMultiPackage_RefusesTwoKeys(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	if _, err := client.Shipping.CreateMultiPackage(context.Background(), flexops.MultiPackageShipmentRequest{}, "a", "b"); err == nil {
		t.Fatal("expected an error for two idempotency keys")
	}
	if len(*calls) != 0 {
		t.Errorf("sent %d requests, want none", len(*calls))
	}
}

func TestReturns_RmaReadsUseTheRmaRoutes(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	ctx := context.Background()
	if _, err := client.Returns.GetRma(ctx, "RMA-1001"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Returns.ListRmas(ctx, flexops.RmaListOptions{Page: 2, PageSize: 50}); err != nil {
		t.Fatal(err)
	}
	_, _ = client.Returns.ListRmas(ctx, flexops.RmaListOptions{})
	_, _ = client.Returns.Get(ctx, "42")
	base := "/api/workspaces/ws-test-123/returns"
	expectCalls(t, *calls,
		"GET "+base+"/rma/RMA-1001 null",
		"GET "+base+"/rma?page=2&pageSize=50 null",
		"GET "+base+"/rma null",
		"GET "+base+"/42 null",
	)
}
