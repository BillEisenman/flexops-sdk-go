package flexops

import "context"

type InsuranceService struct{ c *Client }

func (s *InsuranceService) GetProviders(ctx context.Context) (ApiResponse[[]string], error) {
	return decode[ApiResponse[[]string]](s.c.http.get(ctx, s.c.wsPath("insurance/providers"), nil))
}
func (s *InsuranceService) GetQuote(ctx context.Context, req InsuranceQuoteRequest) (ApiResponse[InsuranceQuote], error) {
	return decode[ApiResponse[InsuranceQuote]](s.c.http.post(ctx, s.c.wsPath("insurance/quote"), req))
}

// Purchase buys a policy and charges the wallet. To reconcile an uncertain outcome,
// pass a stored key with WithIdempotencyKey and reuse it only for this purchase.
func (s *InsuranceService) Purchase(ctx context.Context, req PurchaseInsuranceRequest) (ApiResponse[InsurancePolicy], error) {
	return decode[ApiResponse[InsurancePolicy]](s.c.http.post(ctx, s.c.wsPath("insurance/policies"), req))
}

// Void cancels a policy and returns it.
func (s *InsuranceService) Void(ctx context.Context, policyID string) (ApiResponse[InsurancePolicy], error) {
	return decode[ApiResponse[InsurancePolicy]](s.c.http.del(ctx, s.c.wsPath("insurance/policies/"+policyID)))
}
func (s *InsuranceService) FileClaim(ctx context.Context, req FileClaimRequest) (ApiResponse[InsuranceClaim], error) {
	return decode[ApiResponse[InsuranceClaim]](s.c.http.post(ctx, s.c.wsPath("insurance/claims"), req))
}
