package flexops_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

type recordedCall struct {
	method string
	path   string
	query  string // "?..." or ""
	header http.Header
	body   map[string]any
}

func recordingClient(t *testing.T, response map[string]any) (*flexops.Client, *[]recordedCall) {
	t.Helper()
	calls := &[]recordedCall{}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		call := recordedCall{method: r.Method, path: r.URL.Path, query: query(r), header: r.Header.Clone()}
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

// FlexOpsGateway#672: the package weight and size reach the return-label wire as the Gateway's fields.
func TestReturnLabel_SendsPackageWeightAndSize(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	ctx := context.Background()

	if _, err := client.Returns.GenerateLabelWithOptions(ctx, "RMA-1", flexops.ReturnLabelOptions{
		WeightOz: 16, LengthIn: 10, WidthIn: 8, HeightIn: 4,
	}); err != nil {
		t.Fatal(err)
	}

	body := (*calls)[0].body
	if (*calls)[0].path != "/api/workspaces/ws-test-123/returns/RMA-1/label" {
		t.Errorf("posted to %s", (*calls)[0].path)
	}
	for field, want := range map[string]float64{"weightOz": 16, "lengthIn": 10, "widthIn": 8, "heightIn": 4} {
		if got, ok := body[field].(float64); !ok || got != want {
			t.Errorf("%s = %v, want %v (body %v)", field, body[field], want, body)
		}
	}
	for _, unset := range []string{"carrierCode", "serviceType", "paymentConfirmed", "hazmatClassification"} {
		if _, present := body[unset]; present {
			t.Errorf("%s was sent although unset (body %v)", unset, body)
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

// Address validation uses the documented POST /api/addresses/validate with the flat
// Gateway body; the hidden /api/shipping/addresses/validate is not called.
func TestShipping_ValidateAddressSendsGatewayRequest(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{
		"isValid": true, "validatedBy": "USPS", "warnings": []string{"Corrected ZIP+4"}, "isResidential": true,
		"suggestedAddress": map[string]any{"line1": "123 MAIN ST", "city": "NEW YORK", "state": "NY", "postalCode": "10001-1234", "country": "US"},
	})
	response, err := client.Shipping.ValidateAddress(context.Background(), flexops.AddressValidationRequest{
		Line1: "123 Main St", City: "New York", State: "NY", PostalCode: "10001"})
	if err != nil || !response.IsValid || response.SuggestedAddress == nil || response.SuggestedAddress.PostalCode != "10001-1234" ||
		response.ValidatedBy != "USPS" || len(response.Warnings) != 1 || response.IsResidential == nil || !*response.IsResidential {
		t.Fatalf("bad response: %+v %v", response, err)
	}
	call := (*calls)[0]
	if call.method != "POST" || call.path != "/api/addresses/validate" {
		t.Errorf("sent %s %s", call.method, call.path)
	}
	got, _ := json.Marshal(call.body)
	if want := `{"city":"New York","line1":"123 Main St","postalCode":"10001","state":"NY"}`; string(got) != want {
		t.Errorf("sent %s, want %s", got, want)
	}
}
