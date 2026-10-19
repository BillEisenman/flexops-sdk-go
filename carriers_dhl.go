package flexops

import "context"

type DhlService struct{ c *Client }

func (s *DhlService) p(path string) string { return proxyBase + path }

func (s *DhlService) ValidateAddress(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlValidateAddress"), toValues(params))) }
func (s *DhlService) GetRates(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlRates"), toValues(params))) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *DhlService) GetMultiPieceRates(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postDhlMultiPieceRates"), body)) }
func (s *DhlService) GetProducts(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlProducts"), toValues(params))) }
func (s *DhlService) Track(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlTrackSingleShipment"), toValues(params))) }
func (s *DhlService) TrackMultiple(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlTrackMultipleShipments"), toValues(params))) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *DhlService) CreatePickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postDhlCreatePickup"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *DhlService) CancelPickup(ctx context.Context) error { _, err := s.c.http.del(ctx, s.p("/api/v2/ShippingLabel/deleteDhlPickup")); return err }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *DhlService) CalculateLandedCost(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postDhlCalculateLandedCost"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *DhlService) ScreenShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postDhlScreenShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *DhlService) UploadInvoice(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v2/ShippingLabel/postDhlUploadInvoice"), body)) }
func (s *DhlService) GetProofOfDelivery(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlElectronicProofOfDelivery"), toValues(params))) }
func (s *DhlService) GetReferenceData(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlReferenceData"), toValues(params))) }
func (s *DhlService) FindServicePoints(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v2/ShippingLabel/getDhlServicePoints"), toValues(params))) }
