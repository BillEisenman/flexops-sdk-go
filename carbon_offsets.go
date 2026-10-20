package flexops

import (
	"context"
	"net/url"
	"strconv"
)

// CarbonOffsetsService quotes and buys carbon offsets for shipments.
type CarbonOffsetsService struct{ c *Client }

func (s *CarbonOffsetsService) Quote(ctx context.Context, req CarbonOffsetQuoteRequest) (ApiResponse[CarbonOffsetQuote], error) {
	return decode[ApiResponse[CarbonOffsetQuote]](s.c.http.post(ctx, s.c.wsPath("carbon-offsets/quote"), req))
}

// Purchase buys an offset and charges the wallet. To reconcile an uncertain outcome,
// pass a stored key with WithIdempotencyKey and reuse it only for this purchase.
func (s *CarbonOffsetsService) Purchase(ctx context.Context, req CarbonOffsetPurchaseRequest) (ApiResponse[CarbonOffset], error) {
	return decode[ApiResponse[CarbonOffset]](s.c.http.post(ctx, s.c.wsPath("carbon-offsets/purchase"), req))
}

// List returns purchased offsets. A page or pageSize of 0 uses the Gateway default (1, 25).
func (s *CarbonOffsetsService) List(ctx context.Context, page, pageSize int) (ApiResponse[[]CarbonOffset], error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.Itoa(pageSize))
	}
	return decode[ApiResponse[[]CarbonOffset]](s.c.http.get(ctx, s.c.wsPath("carbon-offsets"), params))
}
func (s *CarbonOffsetsService) Get(ctx context.Context, offsetID string) (ApiResponse[CarbonOffset], error) {
	return decode[ApiResponse[CarbonOffset]](s.c.http.get(ctx, s.c.wsPath("carbon-offsets/"+offsetID), nil))
}
func (s *CarbonOffsetsService) Summary(ctx context.Context) (ApiResponse[CarbonOffsetSummary], error) {
	return decode[ApiResponse[CarbonOffsetSummary]](s.c.http.get(ctx, s.c.wsPath("carbon-offsets/summary"), nil))
}
func (s *CarbonOffsetsService) GetSettings(ctx context.Context) (ApiResponse[CarbonOffsetSettings], error) {
	return decode[ApiResponse[CarbonOffsetSettings]](s.c.http.get(ctx, s.c.wsPath("carbon-offsets/settings"), nil))
}
func (s *CarbonOffsetsService) UpdateSettings(ctx context.Context, settings CarbonOffsetSettings) (ApiResponse[CarbonOffsetSettings], error) {
	return decode[ApiResponse[CarbonOffsetSettings]](s.c.http.put(ctx, s.c.wsPath("carbon-offsets/settings"), settings))
}
