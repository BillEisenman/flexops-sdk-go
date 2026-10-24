package flexops_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// SCAN form replies as the Gateway writes them (ScanFormsController, ScanFormDto), from
// testdata/scan-forms.json.

type scanFormFixture struct {
	Status int             `json:"status"`
	Body   json.RawMessage `json:"body"`
}

func scanFormFixtures(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	raw, err := os.ReadFile("testdata/scan-forms.json")
	if err != nil {
		t.Fatal(err)
	}
	var all map[string]json.RawMessage
	if err := json.Unmarshal(raw, &all); err != nil {
		t.Fatal(err)
	}
	return all
}

func replyingClient(t *testing.T, status int, body json.RawMessage, calls *int32) *flexops.Client {
	t.Helper()
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(calls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(body)
	}))
	flexops.SetRetryBaseDelay(client, time.Millisecond)
	return client
}

func TestScanForms_CreatedFormHasTheGatewayFields(t *testing.T) {
	var calls int32
	client := replyingClient(t, 200, scanFormFixtures(t)["created"], &calls)

	reply, err := client.ScanForms.Create(context.Background(), flexops.ScanFormRequest{
		TrackingNumbers: []string{"9400111899223456789014", "9400111899223456789021"}, ManifestDate: "2026-09-08",
	})
	if err != nil {
		t.Fatal(err)
	}

	form := reply.Data
	want := map[string]any{
		"ScanFormID": "3f2b8c1e9d4a4f6b8a7e2c5d1b0f9e88", "State": flexops.ScanFormStateCreated,
		"OriginPostalCode": "78701", "ManifestNumber": "9475711899560013735215", "ManifestDate": "2026-09-08T00:00:00",
		"PackageCount": 2, "Success": true, "IsSandbox": false, "ErrorMessage": "",
		"StatusMessage": "Created — USPS acceptance requires a carrier scan.",
	}
	got := reflect.ValueOf(form)
	for field, value := range want {
		if actual := got.FieldByName(field).Interface(); !reflect.DeepEqual(actual, value) {
			t.Errorf("%s = %v, want %v", field, actual, value)
		}
	}
	if !reflect.DeepEqual(form.TrackingNumbers, []string{"9400111899223456789014", "9400111899223456789021"}) {
		t.Errorf("TrackingNumbers = %v", form.TrackingNumbers)
	}
	if !form.CreatedAtUTC.Equal(time.Date(2026, 9, 8, 21, 14, 3, 512_000_000, time.UTC)) {
		t.Errorf("CreatedAtUTC = %v", form.CreatedAtUTC)
	}
	if !form.Created() {
		t.Error("Created() = false for a Created form")
	}
	pdf, err := form.PDF()
	if err != nil || !strings.HasPrefix(string(pdf), "%PDF-") {
		t.Errorf("PDF() = %q, %v", pdf, err)
	}
	// The deprecated names, which the Gateway never sent, now carry the real values.
	if form.ID != form.ScanFormID || form.CreatedAt != "2026-09-08T21:14:03.512Z" {
		t.Errorf("deprecated ID %q, CreatedAt %q", form.ID, form.CreatedAt)
	}
}

func TestScanForms_ListHasSandboxFormsWithoutAPDFAndUnresolvedOperations(t *testing.T) {
	var calls int32
	client := replyingClient(t, 200, scanFormFixtures(t)["list"], &calls)

	reply, err := client.ScanForms.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(reply.Data) != 2 {
		t.Fatalf("got %d forms, want 2", len(reply.Data))
	}
	sandbox, unknown := reply.Data[0], reply.Data[1]
	if !sandbox.IsSandbox || sandbox.ScanFormImage != "" || !strings.HasPrefix(sandbox.ManifestNumber, "SBOX") {
		t.Errorf("sandbox form = %+v", sandbox)
	}
	if pdf, err := sandbox.PDF(); pdf != nil || err != nil {
		t.Errorf("sandbox PDF() = %v, %v; want nil, nil", pdf, err)
	}
	if unknown.State != flexops.ScanFormStateOutcomeUnknown || unknown.Success || unknown.Created() ||
		unknown.ManifestNumber != "" || unknown.PackageCount != 2 {
		t.Errorf("unresolved form = %+v", unknown)
	}
}

func TestScanForms_GetReadsOneForm(t *testing.T) {
	var calls int32
	client := replyingClient(t, 200, scanFormFixtures(t)["created"], &calls)

	reply, err := client.ScanForms.Get(context.Background(), "3f2b8c1e9d4a4f6b8a7e2c5d1b0f9e88")
	if err != nil {
		t.Fatal(err)
	}
	if reply.Data.ScanFormID != "3f2b8c1e9d4a4f6b8a7e2c5d1b0f9e88" {
		t.Errorf("ScanFormID = %q", reply.Data.ScanFormID)
	}
}

func TestScanForms_Gateway409IsReturnedOnceWithoutARetry(t *testing.T) {
	for _, name := range []string{"requestIdReused", "needsReview"} {
		t.Run(name, func(t *testing.T) {
			var fixture scanFormFixture
			if err := json.Unmarshal(scanFormFixtures(t)[name], &fixture); err != nil {
				t.Fatal(err)
			}
			var body struct {
				Message string `json:"message"`
			}
			_ = json.Unmarshal(fixture.Body, &body)
			var calls int32
			client := replyingClient(t, fixture.Status, fixture.Body, &calls)

			_, err := client.ScanForms.Create(context.Background(), flexops.ScanFormRequest{
				TrackingNumbers: []string{"9400"}, ManifestDate: "2026-09-08", RequestID: "634648d2-7368-4317-a487-c7b71ba55bb1",
			})

			var apiErr *flexops.FlexOpsError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 409 || apiErr.Message != body.Message {
				t.Fatalf("err = %v, want a 409 FlexOpsError with %q", err, body.Message)
			}
			if calls != 1 {
				t.Errorf("%d calls, want 1 (a 409 is not retried)", calls)
			}
		})
	}
}
