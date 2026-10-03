package flexops_test

import (
	"context"
	"encoding/json"
	"errors"
	flexops "github.com/BillEisenman/flexops-sdk-go/v2"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestInternationalCustomsContract(t *testing.T) {
	data, err := os.ReadFile("examples/international-label.json")
	if err != nil {
		t.Fatal(err)
	}
	var request flexops.CreateLabelRequest
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatal(err)
	}
	calls := 0
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var sent flexops.CreateLabelRequest
		if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(sent.CustomsDeclaration, request.CustomsDeclaration) || sent.OrderID != 42 || sent.ShipDate != "2099-01-01" {
			t.Error("international contract lost fields")
		}
		if calls == 1 {
			if sent.ConfirmationToken != "" {
				t.Error("preview confirmed")
			}
			writeJSON(w, 200, map[string]any{"status": "Preview", "confirmationToken": "approved"})
		} else {
			if r.Header.Get("Idempotency-Key") != "international-1" || sent.ConfirmationToken != "approved" {
				t.Error("approval identity lost")
			}
			writeJSON(w, 201, map[string]any{"labelId": "intl-1", "currency": "USD"})
		}
	}))
	preview, err := client.Shipping.CreateLabel(context.Background(), request)
	if err != nil || preview.Status != "Preview" || calls != 1 {
		t.Fatalf("preview failed: %v", err)
	}
	request.ConfirmationToken = preview.ConfirmationToken
	label, err := client.Shipping.CreateLabel(context.Background(), request, "international-1")
	if err != nil || label.LabelID != "intl-1" || calls != 2 {
		t.Fatalf("purchase failed: %v", err)
	}
}

func TestInternationalDisabledNoRetry(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		writeJSON(w, 403, map[string]any{"errorCode": "FeatureDisabled", "message": "International disabled"})
	}))
	_, err := client.Shipping.GetRates(context.Background(), rateRequest())
	var apiErr *flexops.FlexOpsError
	if !errors.As(err, &apiErr) || apiErr.Code != "FeatureDisabled" || calls != 1 {
		t.Fatalf("gate refusal lost: %v", err)
	}
}
