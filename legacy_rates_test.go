package flexops_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	flexops "github.com/BillEisenman/flexops-sdk-go/v2"
)

// The Gateway RateShoppingRequest needs full origin and destination addresses, so a
// zip-only request is refused before it is sent (as in the Node SDK).
func TestShipping_ZipOnlyRateRequestIsRefused(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		writeJSON(w, 200, map[string]any{"rates": []any{}})
	}))
	ctx := context.Background()
	legacy := flexops.RateRequest{FromZip: "10001", ToZip: "90210", Weight: 16}

	_, errRates := client.Shipping.GetRates(ctx, legacy)
	_, errCheapest := client.Shipping.GetCheapestRate(ctx, legacy)
	_, errFastest := client.Shipping.GetFastestRate(ctx, legacy)
	for _, err := range []error{errRates, errCheapest, errFastest} {
		var apiErr *flexops.FlexOpsError
		if !errors.As(err, &apiErr) || apiErr.Code != "UnsupportedRequestShape" || apiErr.StatusCode != 400 {
			t.Errorf("expected UnsupportedRequestShape, got %v", err)
		}
	}
	if calls != 0 {
		t.Fatalf("sent %d zip-only requests", calls)
	}

	canonical := rateRequest()
	canonical.FromZip = "10001" // A stray deprecated field does not block a full request.
	if _, err := client.Shipping.GetRates(ctx, canonical); err != nil || calls != 1 {
		t.Fatalf("canonical request: %v (calls %d)", err, calls)
	}
}
