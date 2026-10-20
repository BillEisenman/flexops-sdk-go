package flexops_test

import (
	"context"
	"encoding/json"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// CreateBatchLabelRequest mirrors the Gateway batch body: orders plus package measurements.
func TestShipping_LabelBatchMatchesGatewayShape(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true, "data": map[string]any{
		"jobId": "6f1c", "status": "Pending", "totalItems": 2, "completedItems": 0, "failedItems": 0, "isDryRun": true}})
	ctx := context.Background()
	request := flexops.CreateBatchLabelRequest{
		Items: []flexops.BatchLabelItem{
			{OrderID: 42, PackageWeightInOunces: 16, PackageLengthInInches: 10, PackageWidthInInches: 6, PackageHeightInInches: 4,
				CarrierCode: "USPS", HazmatClassification: flexops.HazmatNone.Ptr()},
			{OrderID: 43, PackageWeightInOunces: 8, PackageLengthInInches: 6, PackageWidthInInches: 4, PackageHeightInInches: 2},
		},
		CallbackURL: "https://example.test/hook",
	}

	job, err := client.Shipping.CreateLabelBatch(ctx, request)
	if err != nil || job.Data.JobID != "6f1c" || job.Data.TotalItems != 2 || !job.Data.IsDryRun {
		t.Fatalf("bad job: %+v %v", job, err)
	}
	if _, err := client.Shipping.PreviewLabelBatch(ctx, request); err != nil {
		t.Fatal(err)
	}

	want := `{"callbackUrl":"https://example.test/hook","items":[` +
		`{"carrierCode":"USPS","hazmatClassification":0,"orderId":42,"packageHeightInInches":4,"packageLengthInInches":10,"packageWeightInOunces":16,"packageWidthInInches":6},` +
		`{"orderId":43,"packageHeightInInches":2,"packageLengthInInches":6,"packageWeightInOunces":8,"packageWidthInInches":4}]}`
	for i, path := range []string{"/api/workspaces/ws-test-123/labels/batch", "/api/workspaces/ws-test-123/labels/batch/preview"} {
		got, _ := json.Marshal((*calls)[i].body)
		if (*calls)[i].path != path || string(got) != want {
			t.Errorf("call %d: %s %s\nwant %s %s", i, (*calls)[i].path, got, path, want)
		}
	}
}
