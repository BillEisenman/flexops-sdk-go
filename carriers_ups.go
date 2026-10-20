package flexops

import (
	"context"
	"net/url"
)

type UpsService struct{ c *Client }

func (s *UpsService) p(path string) string { return proxyBase + path }

// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) ValidateAddress(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsVerifyAddress"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) GetRates(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsRateCheck"), body)) }
func (s *UpsService) Track(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getSingleUpsTrackingDetail"), toValues(params))) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) CreatePickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsCreatePickup"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) CancelPickup(ctx context.Context) error { _, err := s.c.http.del(ctx, s.p("/api/v2/ShippingLabel/deleteUpsPickup")); return err }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) GetTransitTimes(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsGetTransitTimes"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) GetLandedCost(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsGetLandedCostQuote"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) SearchLocations(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsSearchLocations"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) UploadDocument(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsUploadPaperlessDocument"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) CreateFreightShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsCreateFreightShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *UpsService) GetFreightRate(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postUpsGetFreightRate"), body)) }

func toValues(m map[string]string) url.Values {
	if m == nil { return nil }
	v := url.Values{}
	for k, val := range m { v.Set(k, val) }
	return v
}
