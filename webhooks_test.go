package flexops_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	flexops "github.com/BillEisenman/flexops-sdk-go/v2"
)

// Fixed vector computed outside this SDK with the Gateway algorithm
// (WebhookDeliveryService): lowercase hex HMAC-SHA256 of the UTF-8 string
// "{timestamp}.{raw body}", keyed with the UTF-8 bytes of the whole secret
// string (the "whsec_" prefix included, nothing base64-decoded).
// Cross-checked with Python hmac and .NET HMACSHA256.
const (
	vectorSecret    = "whsec_AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8="
	vectorTimestamp = "1760600000"
	vectorBody      = `{"trackingNumber":"9400111899223456789012","carrierCode":"USPS","status":"Delivered","note":"café"}`
	vectorSignature = "sha256=47e071c489ace7e261e373daa4a358801f8d9b1d4c7d075d29e11c4fd1cd80f7"
)

// vectorNow is the delivery time of the vector, so tolerance checks are deterministic.
var vectorNow = time.Unix(1760600000, 0)

func verifyAt(now time.Time, body, signature, timestamp, secret string, opts ...flexops.WebhookVerifyOption) error {
	opts = append([]flexops.WebhookVerifyOption{flexops.WithWebhookClock(func() time.Time { return now })}, opts...)
	return flexops.VerifyWebhookSignature([]byte(body), signature, timestamp, secret, opts...)
}

func TestVerifyWebhookSignature_ServerVector(t *testing.T) {
	if err := verifyAt(vectorNow, vectorBody, vectorSignature, vectorTimestamp, vectorSecret); err != nil {
		t.Fatalf("expected the server vector to verify, got %v", err)
	}
}

func TestVerifySignature_UsesRealClock(t *testing.T) {
	// The bool helper uses time.Now, so the fixed 2025 vector is stale for it.
	if flexops.VerifySignature([]byte(vectorBody), vectorSignature, vectorTimestamp, vectorSecret) {
		t.Fatal("expected VerifySignature to reject a timestamp outside the default tolerance")
	}

	// A delivery signed now, the way the Gateway signs it, verifies.
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte(vectorSecret))
	mac.Write([]byte(ts + "." + vectorBody))
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if !flexops.VerifySignature([]byte(vectorBody), sig, ts, vectorSecret) {
		t.Fatal("expected VerifySignature to accept a fresh delivery")
	}
}

func TestVerifyWebhookSignature_PrefixHandling(t *testing.T) {
	bareHex := strings.TrimPrefix(vectorSignature, "sha256=")
	if err := verifyAt(vectorNow, vectorBody, bareHex, vectorTimestamp, vectorSecret); err != nil {
		t.Errorf("expected bare hex to verify, got %v", err)
	}
	if err := verifyAt(vectorNow, vectorBody, strings.ToUpper(bareHex), vectorTimestamp, vectorSecret); err != nil {
		t.Errorf("expected upper-case hex to verify, got %v", err)
	}
	if err := verifyAt(vectorNow, vectorBody, "sha1="+bareHex, vectorTimestamp, vectorSecret); !errors.Is(err, flexops.ErrWebhookSignatureMismatch) {
		t.Errorf("expected another scheme prefix to be rejected, got %v", err)
	}
}

func TestVerifyWebhookSignature_Rejects(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		signature string
		timestamp string
		secret    string
		now       time.Time
		want      error
	}{
		{"wrong secret", vectorBody, vectorSignature, vectorTimestamp, "whsec_wrong", vectorNow, flexops.ErrWebhookSignatureMismatch},
		{"tampered body", strings.Replace(vectorBody, "Delivered", "Returned", 1), vectorSignature, vectorTimestamp, vectorSecret, vectorNow, flexops.ErrWebhookSignatureMismatch},
		{"timestamp changed", vectorBody, vectorSignature, "1760600001", vectorSecret, vectorNow, flexops.ErrWebhookSignatureMismatch},
		{"body-only signature (old algorithm)", vectorBody, "sha256=4511bbc1123d6aaf6d2125b2043b8202385d5b35fdd4d180356499df7f18956d", vectorTimestamp, vectorSecret, vectorNow, flexops.ErrWebhookSignatureMismatch},
		{"truncated signature", vectorBody, vectorSignature[:len(vectorSignature)-2], vectorTimestamp, vectorSecret, vectorNow, flexops.ErrWebhookSignatureMismatch},
		{"extended signature", vectorBody, vectorSignature + "00", vectorTimestamp, vectorSecret, vectorNow, flexops.ErrWebhookSignatureMismatch},
		{"non-hex signature", vectorBody, "sha256=zz", vectorTimestamp, vectorSecret, vectorNow, flexops.ErrWebhookSignatureMismatch},
		{"missing signature", vectorBody, "", vectorTimestamp, vectorSecret, vectorNow, flexops.ErrWebhookMissingSignature},
		{"missing timestamp", vectorBody, vectorSignature, "", vectorSecret, vectorNow, flexops.ErrWebhookInvalidTimestamp},
		{"non-numeric timestamp", vectorBody, vectorSignature, "17606OOOOO", vectorSecret, vectorNow, flexops.ErrWebhookInvalidTimestamp},
		{"signed timestamp", vectorBody, vectorSignature, "+1760600000", vectorSecret, vectorNow, flexops.ErrWebhookInvalidTimestamp},
		{"stale timestamp", vectorBody, vectorSignature, vectorTimestamp, vectorSecret, vectorNow.Add(301 * time.Second), flexops.ErrWebhookTimestampOutsideTolerance},
		{"future timestamp", vectorBody, vectorSignature, vectorTimestamp, vectorSecret, vectorNow.Add(-301 * time.Second), flexops.ErrWebhookTimestampOutsideTolerance},
		{"empty secret", vectorBody, vectorSignature, vectorTimestamp, "", vectorNow, flexops.ErrWebhookMissingSecret},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyAt(tc.now, tc.body, tc.signature, tc.timestamp, tc.secret)
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}

func TestVerifyWebhookSignature_ToleranceBoundary(t *testing.T) {
	if err := verifyAt(vectorNow.Add(300*time.Second), vectorBody, vectorSignature, vectorTimestamp, vectorSecret); err != nil {
		t.Errorf("expected a 300 s old delivery to verify, got %v", err)
	}
	if err := verifyAt(vectorNow.Add(-300*time.Second), vectorBody, vectorSignature, vectorTimestamp, vectorSecret); err != nil {
		t.Errorf("expected a delivery 300 s in the future to verify, got %v", err)
	}
	if err := verifyAt(vectorNow.Add(10*time.Minute), vectorBody, vectorSignature, vectorTimestamp, vectorSecret,
		flexops.WithWebhookTolerance(15*time.Minute)); err != nil {
		t.Errorf("expected a custom tolerance to apply, got %v", err)
	}
	if err := verifyAt(vectorNow.Add(10*time.Minute), vectorBody, vectorSignature, vectorTimestamp, vectorSecret,
		flexops.WithWebhookTolerance(0)); !errors.Is(err, flexops.ErrWebhookTimestampOutsideTolerance) {
		t.Errorf("expected a non-positive tolerance to fall back to the default, got %v", err)
	}
}

// capturedRequest records one request the test server received.
type capturedRequest struct {
	method, path string
	body         map[string]any
}

// webhookServer answers like the Gateway WebhookSubscriptionsController: create
// returns the raw WebhookSubscriptionCreatedDto with 201, update an ApiResponse.
func webhookServer(t *testing.T) (*flexops.Client, *[]capturedRequest) {
	t.Helper()
	var seen []capturedRequest
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Errorf("request body is not JSON: %s", raw)
			}
		}
		seen = append(seen, capturedRequest{r.Method, r.URL.Path, body})
		subscription := map[string]any{
			"id": "0b6f2d4e-0000-0000-0000-000000000001", "url": "https://example.com/hook", "description": nil,
			"eventTypes": "LabelCreated,BatchCompleted", "status": "Active", "consecutiveFailures": 0,
			"createdAt": "2026-10-16T12:00:00Z", "updatedAt": nil,
		}
		if r.Method == http.MethodPost {
			subscription["secret"] = "whsec_test"
			writeJSON(w, http.StatusCreated, subscription)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": subscription})
	}))
	return client, &seen
}

func TestWebhookEventTypes_AreTheNamesTheGatewaySends(t *testing.T) {
	got := []flexops.WebhookEventType{
		flexops.WebhookEventLabelCreated, flexops.WebhookEventLabelCanceled, flexops.WebhookEventBatchCompleted,
		flexops.WebhookEventTrackingUpdated, flexops.WebhookEventTrackingDelivered,
	}
	want := []flexops.WebhookEventType{"LabelCreated", "LabelCanceled", "BatchCompleted", "TrackingUpdated", "TrackingDelivered"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("event names = %v, want %v", got, want)
	}
}

// The Gateway reads a comma-separated eventTypes (no spaces) and ignores an
// events field, so an events array subscribed to every event.
func TestWebhooksCreate_SendsEventTypes(t *testing.T) {
	client, seen := webhookServer(t)

	res, err := client.Webhooks.Create(context.Background(), flexops.CreateWebhookRequest{
		URL:    "https://example.com/hook",
		Events: []flexops.WebhookEventType{flexops.WebhookEventLabelCreated, flexops.WebhookEventBatchCompleted, flexops.WebhookEventLabelCreated},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(*seen) != 1 || (*seen)[0].method != http.MethodPost || (*seen)[0].path != "/api/workspaces/ws-test-123/webhooks" {
		t.Fatalf("requests = %+v", *seen)
	}
	want := map[string]any{"url": "https://example.com/hook", "eventTypes": "LabelCreated,BatchCompleted"}
	if !reflect.DeepEqual((*seen)[0].body, want) {
		t.Fatalf("body = %v, want %v", (*seen)[0].body, want)
	}
	// Create returns the subscription without an envelope; the secret comes back only here.
	sub := res.Data
	if !res.Success || sub.ID == "" || sub.Secret != "whsec_test" || sub.EventTypes != "LabelCreated,BatchCompleted" || sub.Status != "Active" || sub.CreatedAt == "" {
		t.Fatalf("result = %+v", res)
	}
}

func TestWebhooksCreate_WithoutEventsReceivesEveryEvent(t *testing.T) {
	client, seen := webhookServer(t)

	if _, err := client.Webhooks.Create(context.Background(), flexops.CreateWebhookRequest{URL: "https://example.com/hook", Description: "orders"}); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"url": "https://example.com/hook", "description": "orders"}
	if len(*seen) != 1 || !reflect.DeepEqual((*seen)[0].body, want) {
		t.Fatalf("requests = %+v, want body %v", *seen, want)
	}
}

func TestWebhooksUpdate_SendsEventTypes(t *testing.T) {
	client, seen := webhookServer(t)

	res, err := client.Webhooks.Update(context.Background(), "wh-001", flexops.UpdateWebhookRequest{
		Events: []flexops.WebhookEventType{flexops.WebhookEventLabelCanceled},
		Status: "Active",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(*seen) != 1 || (*seen)[0].method != http.MethodPut || (*seen)[0].path != "/api/workspaces/ws-test-123/webhooks/wh-001" {
		t.Fatalf("requests = %+v", *seen)
	}
	want := map[string]any{"eventTypes": "LabelCanceled", "status": "Active"}
	if !reflect.DeepEqual((*seen)[0].body, want) {
		t.Fatalf("body = %v, want %v", (*seen)[0].body, want)
	}
	if res.Data.EventTypes == "" || res.Data.Status != "Active" {
		t.Fatalf("result = %+v", res)
	}
}

func TestWebhooksUpdate_ClearsDescription(t *testing.T) {
	client, seen := webhookServer(t)
	empty := ""

	if _, err := client.Webhooks.Update(context.Background(), "wh-001", flexops.UpdateWebhookRequest{Description: &empty}); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"description": ""}
	if len(*seen) != 1 || !reflect.DeepEqual((*seen)[0].body, want) {
		t.Fatalf("requests = %+v, want body %v", *seen, want)
	}
}

func TestWebhooks_RefuseEventsFlexOpsNeverSends(t *testing.T) {
	cases := map[string][]flexops.WebhookEventType{
		"dotted name":           {"label.created"},
		"enum name never sent":  {flexops.WebhookEventLabelCreated, "ShipmentCreated"},
		"wrong case":            {"labelcreated"},
		"space after comma":     {" LabelCanceled"},
		"empty non-nil list":    {},
		"comma inside one name": {"LabelCreated,LabelCanceled"},
	}
	for name, events := range cases {
		t.Run(name, func(t *testing.T) {
			client, seen := webhookServer(t)
			ctx := context.Background()

			_, createErr := client.Webhooks.Create(ctx, flexops.CreateWebhookRequest{URL: "https://example.com/hook", Events: events})
			_, updateErr := client.Webhooks.Update(ctx, "wh-001", flexops.UpdateWebhookRequest{Events: events})

			for _, err := range []error{createErr, updateErr} {
				var apiErr *flexops.FlexOpsError
				if !errors.As(err, &apiErr) || apiErr.Code != "UnknownWebhookEvent" {
					t.Errorf("expected an UnknownWebhookEvent error, got %v", err)
				}
			}
			if len(*seen) != 0 {
				t.Errorf("expected no request, got %+v", *seen)
			}
		})
	}
}
