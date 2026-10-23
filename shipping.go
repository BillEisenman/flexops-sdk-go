package flexops

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type ShippingService struct{ c *Client }

func (s *ShippingService) GetRates(ctx context.Context, req RateRequest) (RateShoppingResponse, error) {
	if err := refuseZipOnly(req); err != nil {
		return RateShoppingResponse{}, err
	}
	return decode[RateShoppingResponse](s.c.http.post(ctx, "/api/shipping/rates", req))
}
func (s *ShippingService) GetCheapestRate(ctx context.Context, req RateRequest) (ShippingRate, error) {
	if err := refuseZipOnly(req); err != nil {
		return ShippingRate{}, err
	}
	return decode[ShippingRate](s.c.http.post(ctx, "/api/shipping/rates/cheapest", req))
}
func (s *ShippingService) GetFastestRate(ctx context.Context, req RateRequest) (ShippingRate, error) {
	if err := refuseZipOnly(req); err != nil {
		return ShippingRate{}, err
	}
	return decode[ShippingRate](s.c.http.post(ctx, "/api/shipping/rates/fastest", req))
}

// refuseZipOnly rejects a request that uses the deprecated FromZip/ToZip fields
// without full addresses: the Gateway requires Origin and Destination, so it would
// only fail there. It returns a *FlexOpsError with Code "UnsupportedRequestShape".
func refuseZipOnly(req RateRequest) error {
	if (req.FromZip != "" || req.ToZip != "") && (req.Origin.AddressLine1 == "" || req.Destination.AddressLine1 == "") {
		return &FlexOpsError{StatusCode: 400, Code: "UnsupportedRequestShape",
			Message: "The Gateway needs full origin and destination addresses. Use Origin, Destination and Package."}
	}
	return nil
}

// CreateLabel previews unless ConfirmationToken is supplied. For a purchase,
// supply one idempotency key and reuse it on retries; never automatically confirm.
func (s *ShippingService) CreateLabel(ctx context.Context, req CreateLabelRequest, idempotencyKey ...string) (LabelPurchaseResult, error) {
	if len(idempotencyKey) > 1 {
		return LabelPurchaseResult{}, fmt.Errorf("supply at most one idempotency key")
	}
	headers := http.Header{}
	if len(idempotencyKey) == 1 {
		headers.Set("Idempotency-Key", idempotencyKey[0])
	}
	return decode[LabelPurchaseResult](s.c.http.do(ctx, http.MethodPost, "/api/shipping/labels", nil, req, headers))
}

// CancelLabel voids a label. The Gateway requires the carrier code (for example
// "USPS"); it is variadic only so existing call sites still compile, and a call
// without it returns an error before any request.
func (s *ShippingService) CancelLabel(ctx context.Context, labelID string, carrierCode ...string) error {
	if len(carrierCode) != 1 || carrierCode[0] == "" {
		return fmt.Errorf("flexops: CancelLabel needs exactly one carrier code (for example \"USPS\")")
	}
	_, err := s.c.http.do(ctx, http.MethodDelete, "/api/shipping/labels/"+url.PathEscape(labelID),
		url.Values{"carrierCode": {carrierCode[0]}}, nil)
	return err
}

// Track returns the Gateway tracking result. The Gateway sends it without an
// envelope; the SDK wraps it in ApiResponse to keep this method's signature.
func (s *ShippingService) Track(ctx context.Context, trackingNumber string) (ApiResponse[TrackingInfo], error) {
	return wrap[TrackingInfo](s.c.http.get(ctx, "/api/shipping/track/"+url.PathEscape(trackingNumber), nil))
}

// ValidateAddress validates an address with the carrier (POST /api/addresses/validate)
// and returns the Gateway result without an envelope. Country defaults to "US" on the
// Gateway. ValidationUnavailable means the check could not complete, not that the
// address was rejected.
func (s *ShippingService) ValidateAddress(ctx context.Context, req AddressValidationRequest) (AddressValidationResult, error) {
	return decode[AddressValidationResult](s.c.http.post(ctx, "/api/addresses/validate", req))
}

// CreateLabelBatch queues labels for a batch of orders (the Gateway CreateBatchLabelRequest).
func (s *ShippingService) CreateLabelBatch(ctx context.Context, req CreateBatchLabelRequest) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.post(ctx, s.c.wsPath("labels/batch"), req))
}

// PreviewLabelBatch estimates a batch without creating labels.
func (s *ShippingService) PreviewLabelBatch(ctx context.Context, req CreateBatchLabelRequest) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.post(ctx, s.c.wsPath("labels/batch/preview"), req))
}

// CreateMultiPackage buys labels for one order shipped in 1 to 99 boxes with the same
// carrier and service (POST /api/shipping/multi-package). Every box's label slot is
// reserved before any label is bought; a box that fails does not stop the others, and
// each box reports its own Success, ErrorMessage and ErrorCode.
//
// It returns the envelope when at least one label was bought (check FailedPackages).
// When none was, the Gateway answers 422 AllPackagesFailed, or 409 OutcomeUnknown when
// a box's carrier outcome is unknown and a label may exist: both return a
// *FlexOpsError whose Data holds the per-box results. After OutcomeUnknown, check the
// order before buying again. HazmatClassification accepts only HazmatNone or nil.
//
// Pass one idempotency key to retry safely across processes (WithIdempotencyKey also
// works); otherwise one is generated for the call.
func (s *ShippingService) CreateMultiPackage(ctx context.Context, req MultiPackageShipmentRequest, idempotencyKey ...string) (ApiResponse[MultiPackageShipmentResponse], error) {
	if len(idempotencyKey) > 1 {
		return ApiResponse[MultiPackageShipmentResponse]{}, fmt.Errorf("supply at most one idempotency key")
	}
	headers := http.Header{}
	if len(idempotencyKey) == 1 {
		headers.Set("Idempotency-Key", idempotencyKey[0])
	}
	return decode[ApiResponse[MultiPackageShipmentResponse]](s.c.http.do(ctx, http.MethodPost, "/api/shipping/multi-package", nil, req, headers))
}

// Deprecated: BatchLabelRequest does not match the Gateway batch body; use CreateLabelBatch.
func (s *ShippingService) CreateBatch(ctx context.Context, req BatchLabelRequest) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.post(ctx, s.c.wsPath("labels/batch"), req))
}

// Deprecated: BatchLabelRequest does not match the Gateway batch body; use PreviewLabelBatch.
func (s *ShippingService) PreviewBatch(ctx context.Context, req BatchLabelRequest) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.post(ctx, s.c.wsPath("labels/batch/preview"), req))
}
func (s *ShippingService) GetBatchStatus(ctx context.Context, jobID string) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.get(ctx, s.c.wsPath("labels/batch/"+jobID), nil))
}
func (s *ShippingService) DownloadBatchLabel(ctx context.Context, jobID, itemID string) ([]byte, error) {
	return s.c.http.get(ctx, s.c.wsPath("labels/batch/"+jobID+"/items/"+itemID+"/label"), nil)
}
func (s *ShippingService) GetCarriers(ctx context.Context) (ApiResponse[any], error) {
	return wrap[any](s.c.http.get(ctx, "/api/shipping/carriers", nil))
}
func (s *ShippingService) GetRecommendations(ctx context.Context, req CarrierRecommendationRequest) (ApiResponse[CarrierRecommendationResponse], error) {
	return wrap[CarrierRecommendationResponse](s.c.http.post(ctx, "/api/shipping/recommendations", req))
}
func (s *ShippingService) PredictDelivery(ctx context.Context, req DeliveryPredictionRequest) (ApiResponse[DeliveryPredictionResponse], error) {
	return wrap[DeliveryPredictionResponse](s.c.http.post(ctx, "/api/shipping/predictions/delivery", req))
}
func (s *ShippingService) GetSavings(ctx context.Context) (ApiResponse[any], error) {
	return wrap[any](s.c.http.get(ctx, "/api/shipping/savings", nil))
}

// wrap decodes a raw Gateway api/shipping response and wraps it in ApiResponse,
// the shape these methods have always returned.
func wrap[T any](data []byte, err error) (ApiResponse[T], error) {
	value, err := decode[T](data, err)
	if err != nil {
		return ApiResponse[T]{}, err
	}
	return ApiResponse[T]{Success: true, Data: value}, nil
}
