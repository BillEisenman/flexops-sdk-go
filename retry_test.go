package flexops_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"
	"time"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// Write retries (founder 2026-10-19): a retried purchase must not charge twice. The
// Gateway deduplicates POST/PUT/PATCH by Idempotency-Key but does not cache a response
// >= 500, and does not deduplicate DELETE.

var uuidV4 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

type attempt struct {
	method string
	key    string
}

// retryServer answers each attempt with the next of statuses (the last repeats); a
// status of 0 drops the connection. The body is {"errorCode": code} for errors.
func retryServer(t *testing.T, code string, statuses ...int) (*flexops.Client, *[]attempt) {
	t.Helper()
	var mu sync.Mutex
	attempts := &[]attempt{}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		*attempts = append(*attempts, attempt{r.Method, r.Header.Get("Idempotency-Key")})
		n := len(*attempts)
		mu.Unlock()
		status := statuses[min(n, len(statuses))-1]
		switch {
		case status == 0:
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Fatal(err)
			}
			_ = conn.Close()
		case status == 429:
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
		case status >= 400:
			writeJSON(w, status, map[string]any{"errorCode": code, "message": "failed"})
		default:
			writeJSON(w, status, map[string]any{"success": true, "data": map[string]any{}})
		}
	}))
	flexops.SetRetryBaseDelay(client, time.Millisecond)
	return client, attempts
}

func purchasePolicy(client *flexops.Client) error {
	_, err := client.Insurance.Purchase(context.Background(), flexops.PurchaseInsuranceRequest{TrackingNumber: "9400", CarrierCode: "USPS", DeclaredValue: 50})
	return err
}

func TestWriteRetry_ReusesOneGeneratedKeyAcrossRetries(t *testing.T) {
	client, attempts := retryServer(t, "", 503, 0, 502, 200)
	if err := purchasePolicy(client); err != nil {
		t.Fatalf("expected success after retries: %v", err)
	}
	if len(*attempts) != 4 {
		t.Fatalf("sent %d attempts, want 4", len(*attempts))
	}
	key := (*attempts)[0].key
	if !uuidV4.MatchString(key) {
		t.Fatalf("Idempotency-Key %q is not a UUIDv4", key)
	}
	for i, a := range *attempts {
		if a.key != key {
			t.Errorf("attempt %d sent key %q, want %q", i, a.key, key)
		}
	}
}

func TestWriteRetry_NewKeyPerCall(t *testing.T) {
	client, attempts := retryServer(t, "", 200)
	_ = purchasePolicy(client)
	_ = purchasePolicy(client)
	if (*attempts)[0].key == "" || (*attempts)[0].key == (*attempts)[1].key {
		t.Errorf("two calls sent keys %q and %q, want two distinct keys", (*attempts)[0].key, (*attempts)[1].key)
	}
}

func TestWriteRetry_NeverOn500(t *testing.T) {
	client, attempts := retryServer(t, "InternalError", 500, 200)
	err := purchasePolicy(client)
	var apiErr *flexops.FlexOpsError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 500 {
		t.Fatalf("expected the 500, got %v", err)
	}
	if len(*attempts) != 1 {
		t.Errorf("a POST was sent %d times after a 500; the Gateway would run it again", len(*attempts))
	}
}

func TestWriteRetry_OnIdempotencyConflictOnly(t *testing.T) {
	client, attempts := retryServer(t, "IDEMPOTENCY_CONFLICT", 409, 200)
	if err := purchasePolicy(client); err != nil || len(*attempts) != 2 || (*attempts)[0].key != (*attempts)[1].key {
		t.Fatalf("IDEMPOTENCY_CONFLICT: %v after %d attempts %v", err, len(*attempts), *attempts)
	}

	// Codes are matched exactly: label purchases and managed returns answer a PascalCase
	// IdempotencyConflict for altered inputs, which must not be retried.
	for _, code := range []string{"ApprovalExpired", "IdempotencyConflict", "idempotency_conflict"} {
		client, attempts = retryServer(t, code, 409, 200)
		if err := purchasePolicy(client); err == nil || len(*attempts) != 1 {
			t.Fatalf("409 %s must not be retried: %v after %d attempts", code, err, len(*attempts))
		}
	}
}

// Gateway #736: a key reused with a different request gets 422 IDEMPOTENCY_KEY_REUSED,
// which is final.
func TestWriteRetry_NotOnIdempotencyKeyReused(t *testing.T) {
	client, attempts := retryServer(t, "IDEMPOTENCY_KEY_REUSED", 422, 200)
	err := purchasePolicy(client)
	var apiErr *flexops.FlexOpsError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 422 || apiErr.Code != "IDEMPOTENCY_KEY_REUSED" {
		t.Fatalf("expected the 422 IDEMPOTENCY_KEY_REUSED error, got %v", err)
	}
	if len(*attempts) != 1 {
		t.Errorf("a 422 IDEMPOTENCY_KEY_REUSED was sent %d times, want once", len(*attempts))
	}
}

func TestWriteRetry_On429(t *testing.T) {
	client, attempts := retryServer(t, "", 429, 200)
	if err := purchasePolicy(client); err != nil || len(*attempts) != 2 {
		t.Fatalf("429: %v after %d attempts", err, len(*attempts))
	}
}

func TestWriteRetry_PutCarriesKey(t *testing.T) {
	client, attempts := retryServer(t, "", 504, 200)
	if _, err := client.Rules.Update(context.Background(), "r1", flexops.UpdateShippingRuleRequest{}); err != nil {
		t.Fatal(err)
	}
	if len(*attempts) != 2 || (*attempts)[0].method != "PUT" || !uuidV4.MatchString((*attempts)[0].key) || (*attempts)[0].key != (*attempts)[1].key {
		t.Fatalf("PUT attempts: %v", *attempts)
	}
}

func TestWriteRetry_OnTimeout(t *testing.T) {
	var mu sync.Mutex
	var keys []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		first := len(keys) == 1
		mu.Unlock()
		if first {
			time.Sleep(300 * time.Millisecond)
		}
		writeJSON(w, 200, map[string]any{"success": true})
	}))
	t.Cleanup(srv.Close)
	client := flexops.NewClient(flexops.Config{APIKey: "k", WorkspaceID: "ws", BaseURL: srv.URL, Timeout: 100 * time.Millisecond})
	flexops.SetRetryBaseDelay(client, time.Millisecond)
	if err := purchasePolicy(client); err != nil {
		t.Fatalf("expected success after a timed-out attempt: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(keys) != 2 || keys[0] == "" || keys[0] != keys[1] {
		t.Fatalf("timeout retry keys: %v", keys)
	}
}

func TestDeleteRetry_OnlyOn429(t *testing.T) {
	for _, tc := range []struct {
		name     string
		statuses []int
		attempts int
		ok       bool
	}{
		{"503", []int{503, 204}, 1, false},
		{"500", []int{500, 204}, 1, false},
		{"network error", []int{0, 204}, 1, false},
		{"429", []int{429, 204}, 2, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, attempts := retryServer(t, "", tc.statuses...)
			err := client.Rules.Delete(context.Background(), "r1")
			if (err == nil) != tc.ok || len(*attempts) != tc.attempts {
				t.Fatalf("DELETE: %v after %d attempts, want ok=%v after %d", err, len(*attempts), tc.ok, tc.attempts)
			}
			if (*attempts)[0].key != "" {
				t.Errorf("DELETE sent Idempotency-Key %q", (*attempts)[0].key)
			}
		})
	}
}

func TestGetRetry_Unchanged(t *testing.T) {
	client, attempts := retryServer(t, "", 500, 0, 503, 200)
	if _, err := client.Rules.Get(context.Background(), "r1"); err != nil {
		t.Fatalf("GET: %v", err)
	}
	if len(*attempts) != 4 || (*attempts)[0].key != "" {
		t.Fatalf("GET attempts: %v", *attempts)
	}
}

func TestRetry_StopsAtMaxRetries(t *testing.T) {
	client, attempts := retryServer(t, "", 503)
	if err := purchasePolicy(client); err == nil || len(*attempts) != 4 {
		t.Fatalf("default: %v after %d attempts, want an error after 4", err, len(*attempts))
	}
}
