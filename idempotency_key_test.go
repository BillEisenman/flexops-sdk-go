package flexops_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

func TestWithIdempotencyKey_IsSentOnEveryAttempt(t *testing.T) {
	client, attempts := retryServer(t, "", 503, 200)
	ctx := flexops.WithIdempotencyKey(context.Background(), "offset-42")
	if _, err := client.CarbonOffsets.Purchase(ctx, flexops.CarbonOffsetPurchaseRequest{TrackingNumber: "9400", WeightLbs: 1, DistanceMiles: 10}); err != nil {
		t.Fatal(err)
	}
	if len(*attempts) != 2 || (*attempts)[0].key != "offset-42" || (*attempts)[1].key != "offset-42" {
		t.Fatalf("attempts: %v", *attempts)
	}
}

func TestCreateLabel_ExplicitKeyWinsOverContextKey(t *testing.T) {
	client, attempts := retryServer(t, "", 201)
	ctx := flexops.WithIdempotencyKey(context.Background(), "from-context")
	if _, err := client.Shipping.CreateLabel(ctx, flexops.CreateLabelRequest{ConfirmationToken: "t"}, "purchase-1"); err != nil {
		t.Fatal(err)
	}
	if (*attempts)[0].key != "purchase-1" {
		t.Fatalf("sent key %q, want purchase-1", (*attempts)[0].key)
	}
}

// A 2xx without a body (204, or Content-Length 0) returns the zero value and no error.
func TestEmptySuccessBody_ReturnsZeroValue(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusOK} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(status) }))
		client := flexops.NewClient(flexops.Config{APIKey: "k", WorkspaceID: "ws", BaseURL: srv.URL})
		policy, err := client.Insurance.Void(context.Background(), "p1")
		if err != nil || policy.Success || policy.Data.PolicyID != "" {
			t.Errorf("status %d: %+v %v", status, policy, err)
		}
		if err := client.Rules.Delete(context.Background(), "r1"); err != nil {
			t.Errorf("status %d: delete: %v", status, err)
		}
		balance, err := client.Wallet.GetBalance(context.Background())
		if err != nil || balance.HasWallet {
			t.Errorf("status %d: balance %+v %v", status, balance, err)
		}
		srv.Close()
	}
}

func TestConfigMaxRetries_IsHonored(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		writeJSON(w, 503, map[string]any{"message": "down"})
	}))
	t.Cleanup(srv.Close)
	client := flexops.NewClient(flexops.Config{APIKey: "k", WorkspaceID: "ws", BaseURL: srv.URL, MaxRetries: 1})
	flexops.SetRetryBaseDelay(client, time.Millisecond)
	if _, err := client.Rules.List(context.Background()); err == nil || calls != 2 {
		t.Fatalf("MaxRetries 1: %v after %d calls, want an error after 2", err, calls)
	}
}
