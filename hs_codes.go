package flexops

import (
	"context"
	"net/url"
)

type HsCodesService struct{ c *Client }

func (s *HsCodesService) Search(ctx context.Context, query string, params url.Values) (ApiResponse[any], error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("q", query)
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("shipping/hs-codes/search"), params))
}
func (s *HsCodesService) Lookup(ctx context.Context, code string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("shipping/hs-codes/"+code), nil))
}
func (s *HsCodesService) EstimateLandedCost(ctx context.Context, req any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("shipping/landed-cost"), req))
}
