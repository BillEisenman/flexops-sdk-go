package flexops

import "context"

type FedExService struct{ c *Client }

func (s *FedExService) p(path string) string { return proxyBase + path }

func (s *FedExService) ValidateAddress(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/AddressValidation/postFedExValidateAndCorrectDomesticAddress"), body)) }
func (s *FedExService) ValidatePostalCode(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/AddressValidation/postFedExValidatePostalCode"), body)) }
func (s *FedExService) GetRates(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/RateCalculator/postRetrieveFedExRateAndTransitTimesAsync"), body)) }
func (s *FedExService) CreateShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Shipping/postFedExCreateNewShipment"), body)) }
func (s *FedExService) CancelShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.put(ctx, s.p("/api/v3/Shipping/putFedExCancelShipment"), body)) }
func (s *FedExService) ValidateShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Shipping/postFedExValidateShipment"), body)) }
func (s *FedExService) CreateReturnShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Shipping/postFedExCreateNewReturnShipment"), body)) }
func (s *FedExService) Track(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Tracking/postFedExRetrieveTrackingInfoByTrackingNumber"), body)) }
func (s *FedExService) TrackMultiPiece(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Tracking/postFedExRetrieveTrackingInfoForMultiPieceShipment"), body)) }
func (s *FedExService) RegisterTrackingNotification(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Tracking/postFedExRegisterForTrackingNotification"), body)) }
func (s *FedExService) CreatePickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/CarrierPickup/postFedExCreateCarrierPickupRequest"), body)) }
func (s *FedExService) CancelPickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.put(ctx, s.p("/api/v3/CarrierPickup/putFedExCancelCarrierPickupRequest"), body)) }
func (s *FedExService) SearchLocations(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/LocationSearch/postFedExSearchValidLocations"), body)) }
func (s *FedExService) GetServiceStandards(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/ServiceStandards/postFedExRetrieveServicesAndTransitTimes"), body)) }
func (s *FedExService) GetFreightRate(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Freight/postFedExGetFreightRateQuote"), body)) }
func (s *FedExService) CreateFreightShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Freight/postFedExCreateFreightShipment"), body)) }
func (s *FedExService) GroundClose(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/GroundClose/postFedExCloseWithDocuments"), body)) }
func (s *FedExService) UploadTradeDocuments(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Trade/postFedExUploadTradeDocuments"), body)) }
func (s *FedExService) CreateOpenShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/OpenShip/postFedExCreateOpenShipment"), body)) }
func (s *FedExService) AddPackagesToOpenShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/OpenShip/postFedExAddPackagesToOpenShipment"), body)) }
func (s *FedExService) ConfirmOpenShipment(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/OpenShip/postFedExConfirmOpenShipment"), body)) }
