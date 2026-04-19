package flexops

import "context"

type ShippingService struct{ c *Client }

func (s *ShippingService) GetRates(ctx context.Context, req RateRequest) (ApiResponse[[]ShippingRate], error) {
	return decode[ApiResponse[[]ShippingRate]](s.c.http.post(ctx, s.c.wsPath("shipping/rates"), req))
}
func (s *ShippingService) GetCheapestRate(ctx context.Context, req RateRequest) (ApiResponse[ShippingRate], error) {
	return decode[ApiResponse[ShippingRate]](s.c.http.post(ctx, s.c.wsPath("shipping/rates/cheapest"), req))
}
func (s *ShippingService) GetFastestRate(ctx context.Context, req RateRequest) (ApiResponse[ShippingRate], error) {
	return decode[ApiResponse[ShippingRate]](s.c.http.post(ctx, s.c.wsPath("shipping/rates/fastest"), req))
}
func (s *ShippingService) CreateLabel(ctx context.Context, req CreateLabelRequest) (ApiResponse[Label], error) {
	return decode[ApiResponse[Label]](s.c.http.post(ctx, s.c.wsPath("shipping/labels"), req))
}
func (s *ShippingService) CancelLabel(ctx context.Context, labelID string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("shipping/labels/"+labelID))
	return err
}
func (s *ShippingService) Track(ctx context.Context, trackingNumber string) (ApiResponse[TrackingInfo], error) {
	return decode[ApiResponse[TrackingInfo]](s.c.http.get(ctx, s.c.wsPath("shipping/track/"+trackingNumber), nil))
}
func (s *ShippingService) ValidateAddress(ctx context.Context, address Address) (ApiResponse[AddressValidationResult], error) {
	return decode[ApiResponse[AddressValidationResult]](s.c.http.post(ctx, s.c.wsPath("shipping/addresses/validate"), address))
}
func (s *ShippingService) CreateBatch(ctx context.Context, req BatchLabelRequest) (ApiResponse[BatchLabelJob], error) {
	return decode[ApiResponse[BatchLabelJob]](s.c.http.post(ctx, s.c.wsPath("labels/batch"), req))
}
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
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("shipping/carriers"), nil))
}
func (s *ShippingService) GetRecommendations(ctx context.Context, req CarrierRecommendationRequest) (ApiResponse[CarrierRecommendationResponse], error) {
	return decode[ApiResponse[CarrierRecommendationResponse]](s.c.http.post(ctx, s.c.wsPath("shipping/recommendations"), req))
}
func (s *ShippingService) PredictDelivery(ctx context.Context, req DeliveryPredictionRequest) (ApiResponse[DeliveryPredictionResponse], error) {
	return decode[ApiResponse[DeliveryPredictionResponse]](s.c.http.post(ctx, s.c.wsPath("shipping/predictions/delivery"), req))
}
func (s *ShippingService) GetSavings(ctx context.Context) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("shipping/savings"), nil))
}
