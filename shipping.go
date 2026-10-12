package flexops

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type ShippingService struct{ c *Client }

func (s *ShippingService) GetRates(ctx context.Context, req RateRequest) (RateShoppingResponse, error) {
	return decode[RateShoppingResponse](s.c.http.post(ctx, "/api/shipping/rates", req))
}
func (s *ShippingService) GetCheapestRate(ctx context.Context, req RateRequest) (ShippingRate, error) {
	return decode[ShippingRate](s.c.http.post(ctx, "/api/shipping/rates/cheapest", req))
}
func (s *ShippingService) GetFastestRate(ctx context.Context, req RateRequest) (ShippingRate, error) {
	return decode[ShippingRate](s.c.http.post(ctx, "/api/shipping/rates/fastest", req))
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

// ValidateShippingAddress validates and corrects an address. It sends the Gateway
// AddressValidationRequest ({"address": {...}}) and returns the Gateway response.
func (s *ShippingService) ValidateShippingAddress(ctx context.Context, address ShippingAddress) (AddressValidationResponse, error) {
	return decode[AddressValidationResponse](s.c.http.post(ctx, "/api/shipping/addresses/validate", AddressValidationRequest{Address: address}))
}

// ValidateAddress validates a legacy Address.
//
// Deprecated: use ValidateShippingAddress. The address is mapped to the Gateway
// field names; the validated address and suggestions are mapped back, and the
// Gateway messages are returned as Errors.
func (s *ShippingService) ValidateAddress(ctx context.Context, address Address) (ApiResponse[AddressValidationResult], error) {
	response, err := s.ValidateShippingAddress(ctx, address.shippingAddress())
	if err != nil {
		return ApiResponse[AddressValidationResult]{}, err
	}
	result := AddressValidationResult{IsValid: response.IsValid, Errors: response.Messages}
	if response.ValidatedAddress != nil {
		normalized := response.ValidatedAddress.legacyAddress()
		result.Normalized = &normalized
	}
	for _, suggestion := range response.Suggestions {
		result.Suggestions = append(result.Suggestions, suggestion.legacyAddress())
	}
	return ApiResponse[AddressValidationResult]{Success: true, Data: result}, nil
}

func (a Address) shippingAddress() ShippingAddress {
	return ShippingAddress{Name: a.Name, AddressLine1: a.Street1, AddressLine2: a.Street2, City: a.City,
		StateProvince: a.State, PostalCode: a.Zip, CountryCode: a.Country, Phone: a.Phone}
}

func (a ShippingAddress) legacyAddress() Address {
	return Address{Name: a.Name, Street1: a.AddressLine1, Street2: a.AddressLine2, City: a.City,
		State: a.StateProvince, Zip: a.PostalCode, Country: a.CountryCode, Phone: a.Phone}
}

// CreateLabelBatch queues labels for a batch of orders (the Gateway CreateBatchLabelRequest).
func (s *ShippingService) CreateLabelBatch(ctx context.Context, req CreateBatchLabelRequest) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.post(ctx, s.c.wsPath("labels/batch"), req))
}

// PreviewLabelBatch estimates a batch without creating labels.
func (s *ShippingService) PreviewLabelBatch(ctx context.Context, req CreateBatchLabelRequest) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.post(ctx, s.c.wsPath("labels/batch/preview"), req))
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
