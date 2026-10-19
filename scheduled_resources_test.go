package flexops_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// Recurring shipments, report schedules and email templates on the Gateway routes and
// shapes (RecurringShipmentsController, ReportSchedulesController,
// ShipmentEmailTemplatesController).

type cannedResponse struct {
	status int
	body   any // nil writes no body
}

// routedClient answers "METHOD /path" from responses (204 with no body when absent) and
// records each request as "METHOD /path body".
func routedClient(t *testing.T, responses map[string]cannedResponse) (*flexops.Client, *[]string, *[]http.Header) {
	t.Helper()
	var mu sync.Mutex
	calls, headers := &[]string{}, &[]http.Header{}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		body := "null"
		if len(raw) > 0 {
			var v any
			if err := json.Unmarshal(raw, &v); err != nil {
				t.Fatalf("body is not JSON: %s", raw)
			}
			var b strings.Builder
			enc := json.NewEncoder(&b)
			enc.SetEscapeHTML(false)
			_ = enc.Encode(v)
			body = strings.TrimSuffix(b.String(), "\n")
		}
		key := r.Method + " " + r.URL.Path
		mu.Lock()
		*calls = append(*calls, key+" "+body)
		*headers = append(*headers, r.Header.Clone())
		mu.Unlock()
		resp, ok := responses[key]
		if !ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if resp.body == nil {
			w.WriteHeader(resp.status)
			return
		}
		writeJSON(w, resp.status, resp.body)
	}))
	return client, calls, headers
}

func expectSent(t *testing.T, calls []string, want ...string) {
	t.Helper()
	if strings.Join(calls, "\n") != strings.Join(want, "\n") {
		t.Errorf("requests:\n%s\nwant:\n%s", strings.Join(calls, "\n"), strings.Join(want, "\n"))
	}
}

func TestRecurringShipments_UseGatewayRoutesAndRawShapes(t *testing.T) {
	base := "/api/workspaces/ws-test-123/recurring-shipments"
	shipment := map[string]any{"id": "rs-1", "name": "Weekly restock", "interval": 1, "intervalName": "Weekly", "isActive": true, "maxFailures": 3}
	client, calls, headers := routedClient(t, map[string]cannedResponse{
		"GET " + base:           {200, []any{shipment}},
		"GET " + base + "/rs-1": {200, shipment},
		"POST " + base:          {201, shipment},
		"PUT " + base + "/rs-1": {200, shipment},
	})
	ctx := context.Background()

	list, err := client.RecurringShipments.List(ctx)
	if err != nil || len(list) != 1 || list[0].ID != "rs-1" || list[0].Interval != flexops.RecurrenceWeekly {
		t.Fatalf("list: %+v %v", list, err)
	}
	got, err := client.RecurringShipments.Get(ctx, "rs-1")
	if err != nil || got.Name != "Weekly restock" || !got.IsActive || got.MaxFailures != 3 {
		t.Fatalf("get: %+v %v", got, err)
	}
	created, err := client.RecurringShipments.Create(ctx, flexops.CreateRecurringShipmentRequest{
		Name: "Weekly restock", AddressID: "addr-1", PresetID: "preset-1", Interval: flexops.RecurrenceDaily, FirstRunUtc: "2026-10-20T15:00:00Z"})
	if err != nil || created.ID != "rs-1" {
		t.Fatalf("create: %+v %v", created, err)
	}
	monthly := flexops.RecurrenceMonthly
	if _, err := client.RecurringShipments.Update(ctx, "rs-1", flexops.UpdateRecurringShipmentRequest{Interval: &monthly}); err != nil {
		t.Fatal(err)
	}
	for _, call := range []func(context.Context, string) error{
		client.RecurringShipments.Delete, client.RecurringShipments.Pause, client.RecurringShipments.Resume,
	} {
		if err := call(ctx, "rs-1"); err != nil {
			t.Fatalf("204 must not be an error: %v", err)
		}
	}
	if err := client.RecurringShipments.Trigger(flexops.WithIdempotencyKey(ctx, "run-2026-10-20"), "rs-1"); err != nil {
		t.Fatal(err)
	}

	expectSent(t, *calls,
		"GET "+base+" null",
		"GET "+base+"/rs-1 null",
		"POST "+base+` {"addressId":"addr-1","firstRunUtc":"2026-10-20T15:00:00Z","interval":0,"name":"Weekly restock","presetId":"preset-1"}`,
		"PUT "+base+`/rs-1 {"interval":3}`,
		"DELETE "+base+"/rs-1 null",
		"POST "+base+"/rs-1/pause null",
		"POST "+base+"/rs-1/resume null",
		"POST "+base+"/rs-1/trigger null",
	)
	if key := (*headers)[7].Get("Idempotency-Key"); key != "run-2026-10-20" {
		t.Errorf("trigger sent Idempotency-Key %q", key)
	}
}

func TestReportSchedules_UseGatewayRoutesAndEnvelope(t *testing.T) {
	base := "/api/workspaces/ws-test-123/report-schedules"
	schedule := map[string]any{"id": "sch-1", "reportType": "shipments", "format": 1, "frequencyType": 1, "dayOfWeek": 1,
		"timeOfDayUtc": "08:00:00", "isEnabled": true, "recipientEmails": "ops@example.com,cfo@example.com"}
	env := func(data any) map[string]any { return map[string]any{"success": true, "data": data} }
	client, calls, _ := routedClient(t, map[string]cannedResponse{
		"GET " + base:                    {200, env([]any{schedule})},
		"GET " + base + "/sch-1":         {200, env(schedule)},
		"POST " + base:                   {201, env(schedule)},
		"PUT " + base + "/sch-1":         {200, env(schedule)},
		"POST " + base + "/sch-1/toggle": {200, map[string]any{"success": true, "data": nil, "message": "Report schedule toggled"}},
	})
	ctx := context.Background()

	list, err := client.ReportSchedules.List(ctx)
	if err != nil || !list.Success || len(list.Data) != 1 || list.Data[0].Format != flexops.ReportFormatJSON {
		t.Fatalf("list: %+v %v", list, err)
	}
	got, err := client.ReportSchedules.Get(ctx, "sch-1")
	if err != nil || got.Data.FrequencyType != flexops.ReportFrequencyWeekly || got.Data.DayOfWeek == nil || *got.Data.DayOfWeek != 1 || got.Data.TimeOfDayUtc != "08:00:00" {
		t.Fatalf("get: %+v %v", got, err)
	}
	monday := 1
	if _, err := client.ReportSchedules.Create(ctx, flexops.CreateReportScheduleRequest{ReportType: "shipments", Format: flexops.ReportFormatCsv,
		FrequencyType: flexops.ReportFrequencyWeekly, DayOfWeek: &monday, TimeOfDayUtc: "08:00:00", RecipientEmails: "ops@example.com"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ReportSchedules.Update(ctx, "sch-1", flexops.UpdateReportScheduleRequest{RecipientEmails: "cfo@example.com"}); err != nil {
		t.Fatal(err)
	}
	toggled, err := client.ReportSchedules.Toggle(ctx, "sch-1")
	if err != nil || !toggled.Success || toggled.Data != nil || toggled.Message != "Report schedule toggled" {
		t.Fatalf("toggle: %+v %v", toggled, err)
	}
	if err := client.ReportSchedules.Delete(ctx, "sch-1"); err != nil {
		t.Fatalf("204 must not be an error: %v", err)
	}

	expectSent(t, *calls,
		"GET "+base+" null",
		"GET "+base+"/sch-1 null",
		"POST "+base+` {"dayOfWeek":1,"format":0,"frequencyType":1,"recipientEmails":"ops@example.com","reportType":"shipments","timeOfDayUtc":"08:00:00"}`,
		"PUT "+base+`/sch-1 {"recipientEmails":"cfo@example.com"}`,
		"POST "+base+"/sch-1/toggle null",
		"DELETE "+base+"/sch-1 null",
	)
}

func TestEmailTemplates_UseGatewayRoutesAndRawShapes(t *testing.T) {
	base := "/api/workspaces/ws-test-123/email-templates"
	template := map[string]any{"id": "tpl-1", "eventType": 3, "eventTypeName": "Delivered", "subjectTemplate": "Delivered: {{trackingNumber}}", "isEnabled": true}
	client, calls, _ := routedClient(t, map[string]cannedResponse{
		"GET " + base:               {200, []any{template}},
		"GET " + base + "/tpl-1":    {200, template},
		"PUT " + base:               {200, template},
		"POST " + base + "/preview": {200, map[string]any{"renderedSubject": "Delivered: 9400", "renderedHtmlBody": "<p>9400</p>"}},
	})
	ctx := context.Background()

	list, err := client.EmailTemplates.List(ctx)
	if err != nil || len(list) != 1 || list[0].EventType != flexops.TrackingEmailDelivered {
		t.Fatalf("list: %+v %v", list, err)
	}
	got, err := client.EmailTemplates.Get(ctx, "tpl-1")
	if err != nil || got.ID != "tpl-1" || got.EventTypeName != "Delivered" || !got.IsEnabled {
		t.Fatalf("get: %+v %v", got, err)
	}
	disabled := false
	if _, err := client.EmailTemplates.Upsert(ctx, flexops.UpsertEmailTemplateRequest{EventType: flexops.TrackingEmailLabelCreated,
		SubjectTemplate: "Shipped", HTMLBodyTemplate: "<p>{{trackingNumber}}</p>"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.EmailTemplates.Upsert(ctx, flexops.UpsertEmailTemplateRequest{EventType: flexops.TrackingEmailException,
		SubjectTemplate: "Problem", HTMLBodyTemplate: "<p>!</p>", IsEnabled: &disabled}); err != nil {
		t.Fatal(err)
	}
	preview, err := client.EmailTemplates.Preview(ctx, flexops.PreviewEmailTemplateRequest{SubjectTemplate: "Delivered: {{trackingNumber}}", HTMLBodyTemplate: "<p>{{trackingNumber}}</p>"})
	if err != nil || preview.RenderedSubject != "Delivered: 9400" || preview.RenderedHTMLBody != "<p>9400</p>" {
		t.Fatalf("preview: %+v %v", preview, err)
	}
	if err := client.EmailTemplates.Delete(ctx, "tpl-1"); err != nil {
		t.Fatalf("204 must not be an error: %v", err)
	}

	expectSent(t, *calls,
		"GET "+base+" null",
		"GET "+base+"/tpl-1 null",
		"PUT "+base+` {"eventType":0,"htmlBodyTemplate":"<p>{{trackingNumber}}</p>","subjectTemplate":"Shipped"}`,
		"PUT "+base+` {"eventType":4,"htmlBodyTemplate":"<p>!</p>","isEnabled":false,"subjectTemplate":"Problem"}`,
		"POST "+base+`/preview {"htmlBodyTemplate":"<p>{{trackingNumber}}</p>","subjectTemplate":"Delivered: {{trackingNumber}}"}`,
		"DELETE "+base+"/tpl-1 null",
	)
}

func TestScheduledResources_EnumNumbersMatchGateway(t *testing.T) {
	for want, got := range []int{int(flexops.RecurrenceDaily), int(flexops.RecurrenceWeekly), int(flexops.RecurrenceBiweekly), int(flexops.RecurrenceMonthly), int(flexops.RecurrenceCustom)} {
		if got != want {
			t.Errorf("RecurrenceInterval %d = %d", want, got)
		}
	}
	for want, got := range []int{int(flexops.TrackingEmailLabelCreated), int(flexops.TrackingEmailInTransit), int(flexops.TrackingEmailOutForDelivery),
		int(flexops.TrackingEmailDelivered), int(flexops.TrackingEmailException), int(flexops.TrackingEmailReturnToSender)} {
		if got != want {
			t.Errorf("TrackingEmailEvent %d = %d", want, got)
		}
	}
	if flexops.ReportFormatCsv != 0 || flexops.ReportFormatJSON != 1 || flexops.ReportFrequencyDaily != 0 || flexops.ReportFrequencyWeekly != 1 || flexops.ReportFrequencyMonthly != 2 {
		t.Error("report enums do not match the Gateway numbers")
	}
}
