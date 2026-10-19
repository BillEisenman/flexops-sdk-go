package flexops

import (
	"context"
	"net/url"
	"strconv"
)

// HsCodesService looks up HS codes and estimates landed cost. The Gateway resolves the
// workspace from the credentials and returns these results without an envelope.
type HsCodesService struct{ c *Client }

// Search finds HS codes for a product description. DestinationCountry and MaxResults
// are optional (the Gateway returns 10 results by default).
func (s *HsCodesService) Search(ctx context.Context, req HsCodeSearchRequest) ([]HsCodeResult, error) {
	params := url.Values{"query": {req.Query}}
	if req.DestinationCountry != "" {
		params.Set("destinationCountry", req.DestinationCountry)
	}
	if req.MaxResults > 0 {
		params.Set("maxResults", strconv.Itoa(req.MaxResults))
	}
	return decode[[]HsCodeResult](s.c.http.get(ctx, "/api/shipping/hs-codes/search", params))
}
func (s *HsCodesService) Lookup(ctx context.Context, code string) (HsCodeResult, error) {
	return decode[HsCodeResult](s.c.http.get(ctx, "/api/shipping/hs-codes/"+url.PathEscape(code), nil))
}
func (s *HsCodesService) EstimateLandedCost(ctx context.Context, req LandedCostRequest) (LandedCostEstimate, error) {
	return decode[LandedCostEstimate](s.c.http.post(ctx, "/api/shipping/landed-cost", req))
}
