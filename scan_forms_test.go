package flexops_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// SCAN form creation, checked against the Gateway's CreateScanFormRequest (2026-10-24):
// live creation needs manifestDate, and a retry must resend the same requestId so the
// Gateway returns the first attempt's form instead of creating a second one.

func TestScanForms_CreateSendsTheGatewayFieldNames(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	warehouse := 7
	_, err := client.ScanForms.Create(context.Background(), flexops.ScanFormRequest{
		TrackingNumbers:  []string{"9400111899223456789014"},
		ManifestDate:     "2026-09-08",
		RequestID:        "634648d2-7368-4317-a487-c7b71ba55bb1",
		OriginPostalCode: "10001",
		WarehouseID:      &warehouse,
	})
	if err != nil {
		t.Fatal(err)
	}
	expectCalls(t, *calls, `POST /api/workspaces/ws-test-123/scan-forms {"manifestDate":"2026-09-08","originPostalCode":"10001","requestId":"634648d2-7368-4317-a487-c7b71ba55bb1","trackingNumbers":["9400111899223456789014"],"warehouseId":7}`)
}

func TestScanForms_CreateGeneratesARequestIDPerCall(t *testing.T) {
	client, calls := recordingClient(t, map[string]any{"success": true})
	req := flexops.ScanFormRequest{TrackingNumbers: []string{"9400"}, ManifestDate: "2026-09-08"}
	for range 2 {
		if _, err := client.ScanForms.Create(context.Background(), req); err != nil {
			t.Fatal(err)
		}
	}
	first, _ := (*calls)[0].body["requestId"].(string)
	second, _ := (*calls)[1].body["requestId"].(string)
	if !uuidV4.MatchString(first) || !uuidV4.MatchString(second) || first == second {
		t.Fatalf("requestIds %q and %q: want two different UUIDv4s", first, second)
	}
	if len((*calls)[0].body) != 3 {
		t.Errorf("body %v: want only trackingNumbers, manifestDate and requestId", (*calls)[0].body)
	}
	if req.RequestID != "" {
		t.Errorf("the caller's request was changed: %q", req.RequestID)
	}
}

func TestScanForms_CreateResendsItsRequestIDOnEveryAttempt(t *testing.T) {
	var mu sync.Mutex
	var ids []string
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		mu.Lock()
		id, _ := body["requestId"].(string)
		ids = append(ids, id)
		n := len(ids)
		mu.Unlock()
		if n == 1 {
			writeJSON(w, 503, map[string]any{"errorCode": "SERVICE_UNAVAILABLE"})
			return
		}
		writeJSON(w, 200, map[string]any{"success": true})
	}))
	flexops.SetRetryBaseDelay(client, time.Millisecond)

	if _, err := client.ScanForms.Create(context.Background(), flexops.ScanFormRequest{TrackingNumbers: []string{"9400"}, ManifestDate: "2026-09-08"}); err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || !uuidV4.MatchString(ids[0]) || ids[1] != ids[0] {
		t.Fatalf("requestIds %v: want the same UUIDv4 on both attempts", ids)
	}
}
