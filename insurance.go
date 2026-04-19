package flexops

import "context"

type InsuranceService struct{ c *Client }

func (s *InsuranceService) GetProviders(ctx context.Context) (ApiResponse[[]string], error) { return decode[ApiResponse[[]string]](s.c.http.get(ctx, s.c.wsPath("insurance/providers"), nil)) }
func (s *InsuranceService) GetQuote(ctx context.Context, req any) (ApiResponse[InsuranceQuote], error) { return decode[ApiResponse[InsuranceQuote]](s.c.http.post(ctx, s.c.wsPath("insurance/quote"), req)) }
func (s *InsuranceService) Purchase(ctx context.Context, req any) (ApiResponse[InsurancePolicy], error) { return decode[ApiResponse[InsurancePolicy]](s.c.http.post(ctx, s.c.wsPath("insurance/purchase"), req)) }
func (s *InsuranceService) Void(ctx context.Context, policyID string) error { _, err := s.c.http.del(ctx, s.c.wsPath("insurance/policies/"+policyID)); return err }
func (s *InsuranceService) FileClaim(ctx context.Context, policyID string, claim any) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("insurance/policies/"+policyID+"/claims"), claim)) }
