package flexops

import "context"

type FedExService struct{ c *Client }

func (s *FedExService) p(path string) string { return proxyBase + path }

// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) ValidateAddress(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/AddressValidation/postFedExValidateAndCorrectDomesticAddress"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) ValidatePostalCode(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/AddressValidation/postFedExValidatePostalCode"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) GetRates(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/RateCalculator/postRetrieveFedExRateAndTransitTimesAsync"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) CancelShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.put(ctx, s.p("/api/v3/Shipping/putFedExCancelShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) ValidateShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Shipping/postFedExValidateShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) Track(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Tracking/postFedExRetrieveTrackingInfoByTrackingNumber"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) TrackMultiPiece(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Tracking/postFedExRetrieveTrackingInfoForMultiPieceShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) RegisterTrackingNotification(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Tracking/postFedExRegisterForTrackingNotification"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) CreatePickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/CarrierPickup/postFedExCreateCarrierPickupRequest"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) CancelPickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.put(ctx, s.p("/api/v3/CarrierPickup/putFedExCancelCarrierPickupRequest"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) SearchLocations(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/LocationSearch/postFedExSearchValidLocations"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) GetServiceStandards(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/ServiceStandards/postFedExRetrieveServicesAndTransitTimes"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) GetFreightRate(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Freight/postFedExGetFreightRateQuote"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) CreateFreightShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Freight/postFedExCreateFreightShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) GroundClose(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/GroundClose/postFedExCloseWithDocuments"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) UploadTradeDocuments(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Trade/postFedExUploadTradeDocuments"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) CreateOpenShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/OpenShip/postFedExCreateOpenShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) AddPackagesToOpenShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/OpenShip/postFedExAddPackagesToOpenShipment"), body)) }
// Deprecated: Legacy carrier proxy (sunset 2027-06-30). Needs a user token: API keys get 403 API_KEY_SCOPE_REQUIRED. Use client.Shipping / client.Pickups / client.ScanForms where they cover this.
func (s *FedExService) ConfirmOpenShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/OpenShip/postFedExConfirmOpenShipment"), body)) }
