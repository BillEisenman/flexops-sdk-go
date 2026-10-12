package flexops_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v2"
)

type recordedCall struct {
	path string
	body map[string]any
}

func recordingClient(t *testing.T, response map[string]any) (*flexops.Client, *[]recordedCall) {
	t.Helper()
	calls := &[]recordedCall{}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		call := recordedCall{path: r.URL.Path}
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &call.body); err != nil {
				t.Fatalf("body is not a JSON object: %s", raw)
			}
		}
		*calls = append(*calls, call)
		writeJSON(w, 200, response)
	}))
	return client, calls
}

func sentHazmat(t *testing.T, body map[string]any) any {
	t.Helper()
	value, ok := body["hazmatClassification"]
	if !ok {
		return "absent"
	}
	return value
}

func TestHazmatClassification_ValuesMatchGateway(t *testing.T) {
	values := []flexops.HazmatClassification{
		flexops.HazmatNone, flexops.HazmatLimitedQuantity, flexops.HazmatLithiumBatteryContainedInEquipment,
		flexops.HazmatLithiumBatteryPackedWithEquipment, flexops.HazmatLithiumBatteryStandalone,
		flexops.HazmatOtherRestricted, flexops.HazmatLithiumBatteryUsedDevice,
	}
	for want, got := range values {
		if int(got) != want {
			t.Errorf("value %d = %d", want, got)
		}
	}
}

func TestHazmatClassification_LabelSentOnlyWhenSet(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"status": "Preview", "confirmationToken": "t"})
	request := flexops.CreateLabelRequest{CarrierCode: "USPS", ServiceCode: "PRIORITY",
		Origin: rateRequest().Origin, Destination: rateRequest().Destination, Package: rateRequest().Package}
	ctx := context.Background()

	if _, err := client.Shipping.CreateLabel(ctx, request); err != nil {
		t.Fatal(err)
	}
	request.HazmatClassification = flexops.HazmatNone.Ptr()
	if _, err := client.Shipping.CreateLabel(ctx, request); err != nil {
		t.Fatal(err)
	}
	request.HazmatClassification = flexops.HazmatLimitedQuantity.Ptr()
	if _, err := client.Shipping.CreateLabel(ctx, request); err != nil {
		t.Fatal(err)
	}

	for i, want := range []any{"absent", 0.0, 1.0} {
		call := (*calls)[i]
		if call.path != "/api/shipping/labels" {
			t.Errorf("call %d posted to %s, want /api/shipping/labels", i, call.path)
		}
		if got := sentHazmat(t, call.body); got != want {
			t.Errorf("call %d sent hazmatClassification %v, want %v", i, got, want)
		}
	}
}

func TestHazmatClassification_RatesSentOnlyWhenSet(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"rates": []any{}})
	ctx := context.Background()
	request := rateRequest()

	if _, err := client.Shipping.GetRates(ctx, request); err != nil {
		t.Fatal(err)
	}
	request.HazmatClassification = flexops.HazmatNone.Ptr()
	if _, err := client.Shipping.GetCheapestRate(ctx, request); err != nil {
		t.Fatal(err)
	}
	request.HazmatClassification = flexops.HazmatLithiumBatteryStandalone.Ptr()
	if _, err := client.Shipping.GetFastestRate(ctx, request); err != nil {
		t.Fatal(err)
	}

	for i, want := range []any{"absent", 0.0, 4.0} {
		if got := sentHazmat(t, (*calls)[i].body); got != want {
			t.Errorf("call %d sent hazmatClassification %v, want %v", i, got, want)
		}
	}
}

func TestHazmatClassification_BatchItemsAndReturnLabels(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	ctx := context.Background()
	item := flexops.CreateLabelRequest{OrderID: 42, HazmatClassification: flexops.HazmatNone.Ptr()}

	if _, err := client.Shipping.CreateBatch(ctx, flexops.BatchLabelRequest{Items: []flexops.CreateLabelRequest{item, {OrderID: 43}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Returns.GenerateLabel(ctx, "RMA-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Returns.GenerateLabelWithOptions(ctx, "RMA-1", flexops.ReturnLabelOptions{HazmatClassification: flexops.HazmatNone.Ptr()}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Returns.GenerateLabelWithOptions(ctx, "RMA-1", flexops.ReturnLabelOptions{}); err != nil {
		t.Fatal(err)
	}

	items := (*calls)[0].body["items"].([]any)
	if got := sentHazmat(t, items[0].(map[string]any)); got != 0.0 {
		t.Errorf("batch item sent hazmatClassification %v, want 0", got)
	}
	if got := sentHazmat(t, items[1].(map[string]any)); got != "absent" {
		t.Errorf("batch item without a declaration sent %v", got)
	}
	// The Gateway binds a required GenerateReturnLabelRequest body; {} takes its defaults.
	if (*calls)[1].body == nil || len((*calls)[1].body) != 0 {
		t.Errorf("GenerateLabel must send {}, sent %v", (*calls)[1].body)
	}
	if (*calls)[2].path != "/api/workspaces/ws-test-123/returns/RMA-1/label" || sentHazmat(t, (*calls)[2].body) != 0.0 {
		t.Errorf("return label with HAZMAT sent %s %v", (*calls)[2].path, (*calls)[2].body)
	}
	if got := sentHazmat(t, (*calls)[3].body); got != "absent" {
		t.Errorf("return label without a declaration sent %v", got)
	}
}

func TestShipping_ValidateAddressSendsGatewayRequest(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{
		"isValid": true, "provider": "USPS", "messages": []string{"Corrected ZIP+4"},
		"validatedAddress": map[string]any{"addressLine1": "123 MAIN ST", "city": "NEW YORK", "stateProvince": "NY", "postalCode": "10001-1234", "countryCode": "US"},
	})
	ctx := context.Background()
	canonical := flexops.ShippingAddress{AddressLine1: "123 Main St", City: "New York", StateProvince: "NY", PostalCode: "10001", CountryCode: "US"}

	response, err := client.Shipping.ValidateShippingAddress(ctx, canonical)
	if err != nil || !response.IsValid || response.ValidatedAddress == nil || response.ValidatedAddress.PostalCode != "10001-1234" {
		t.Fatalf("bad response: %+v %v", response, err)
	}
	legacy, err := client.Shipping.ValidateAddress(ctx, flexops.Address{Name: "J", Street1: "123 Main St", City: "New York", State: "NY", Zip: "10001", Country: "US", Phone: "555"})
	if err != nil || !legacy.Success || !legacy.Data.IsValid || legacy.Data.Normalized == nil || legacy.Data.Normalized.Zip != "10001-1234" {
		t.Fatalf("bad legacy response: %+v %v", legacy, err)
	}

	want := []map[string]any{
		{"addressLine1": "123 Main St", "city": "New York", "stateProvince": "NY", "postalCode": "10001", "countryCode": "US"},
		{"name": "J", "addressLine1": "123 Main St", "city": "New York", "stateProvince": "NY", "postalCode": "10001", "countryCode": "US", "phone": "555"},
	}
	for i, call := range *calls {
		if call.path != "/api/shipping/addresses/validate" {
			t.Errorf("call %d posted to %s", i, call.path)
		}
		got, _ := json.Marshal(call.body["address"])
		expected, _ := json.Marshal(want[i])
		if string(got) != string(expected) || len(call.body) != 1 {
			t.Errorf("call %d sent %v, want {address: %s}", i, call.body, expected)
		}
	}
}
