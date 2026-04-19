package flexops

import (
	"context"
	"net/url"
)

type UpsService struct{ c *Client }

func (s *UpsService) p(path string) string { return proxyBase + path }

func (s *UpsService) ValidateAddress(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsVerifyAddress"), body)) }
func (s *UpsService) GetRates(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsRateCheck"), body)) }
func (s *UpsService) CreateLabel(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/generateNewUpsShipLabel"), body)) }
func (s *UpsService) Track(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getSingleUpsTrackingDetail"), toValues(params))) }
func (s *UpsService) CreatePickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsCreatePickup"), body)) }
func (s *UpsService) CancelPickup(ctx context.Context) error { _, err := s.c.http.del(ctx, s.p("/api/v2/ShippingLabel/deleteUpsPickup")); return err }
func (s *UpsService) GetTransitTimes(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsGetTransitTimes"), body)) }
func (s *UpsService) GetLandedCost(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsGetLandedCostQuote"), body)) }
func (s *UpsService) SearchLocations(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsSearchLocations"), body)) }
func (s *UpsService) UploadDocument(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsUploadPaperlessDocument"), body)) }
func (s *UpsService) CreateFreightShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsCreateFreightShipment"), body)) }
func (s *UpsService) GetFreightRate(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsGetFreightRate"), body)) }

func toValues(m map[string]string) url.Values {
	if m == nil { return nil }
	v := url.Values{}
	for k, val := range m { v.Set(k, val) }
	return v
}
