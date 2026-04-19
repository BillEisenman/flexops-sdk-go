package flexops

import "context"

const proxyBase = "/api/ApiProxy"

type CarriersService struct {
	USPS  *UspsService
	UPS   *UpsService
	FedEx *FedExService
	DHL   *DhlService
	c     *Client
}

func newCarriersService(c *Client) *CarriersService {
	return &CarriersService{
		USPS:  &UspsService{c: c},
		UPS:   &UpsService{c: c},
		FedEx: &FedExService{c: c},
		DHL:   &DhlService{c: c},
		c:     c,
	}
}

// --- USPS ---
type UspsService struct{ c *Client }

func (s *UspsService) p(path string) string { return proxyBase + path }

func (s *UspsService) ValidateAddress(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/AddressValidation/getUspsValidateAndCorrectAddress"), toValues(params))) }
func (s *UspsService) CityStateLookup(ctx context.Context, zipCode string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/AddressValidation/getUspsCityStateLookupByZipCode"), toValues(map[string]string{"zipCode": zipCode}))) }
func (s *UspsService) ZipCodeLookup(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/AddressValidation/getUspsZipCodeLookupByAddress"), toValues(params))) }
func (s *UspsService) GetDomesticRates(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/RateCalculator/postUspsSearchDomesticBaseRates"), body)) }
func (s *UspsService) GetDomesticProducts(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/RateCalculator/postUspsSearchEligibleDomesticProducts"), body)) }
func (s *UspsService) GetDomesticPrices(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/RateCalculator/postUspsSearchEligibleDomesticPrices"), body)) }
func (s *UspsService) GetInternationalRates(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/RateCalculator/postUspsSearchInternationalBaseRates"), body)) }
func (s *UspsService) GetInternationalPrices(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/RateCalculator/postUspsSearchEligibleInternationalPrices"), body)) }
func (s *UspsService) CreateDomesticLabel(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Shipping/postUspsGenerateDomesticShippingLabel"), body)) }
func (s *UspsService) CreateReturnLabel(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Shipping/postUspsGenerateDomesticReturnsShippingLabel"), body)) }
func (s *UspsService) CreateInternationalLabel(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/Shipping/postUspsGenerateInternationalShippingLabel"), body)) }
func (s *UspsService) CancelDomesticLabel(ctx context.Context) error { _, err := s.c.http.del(ctx, s.p("/api/v3/Shipping/cancelUspsDomesticShipmentLabel")); return err }
func (s *UspsService) CancelInternationalLabel(ctx context.Context) error { _, err := s.c.http.del(ctx, s.p("/api/v3/Shipping/cancelUspsInternationalShipmentLabel")); return err }
func (s *UspsService) TrackSummary(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/Tracking/getUspsTrackingSummaryInformation"), toValues(params))) }
func (s *UspsService) TrackDetail(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/Tracking/getUspsTrackingDetailInformation"), toValues(params))) }
func (s *UspsService) CreatePickup(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/CarrierPickup/postUspsCreateCarrierPickupSchedule"), body)) }
func (s *UspsService) CancelPickup(ctx context.Context) error { _, err := s.c.http.del(ctx, s.p("/api/v3/CarrierPickup/cancelUspsCarrierPickupSchedule")); return err }
func (s *UspsService) CreateScanForm(ctx context.Context, body any) (any, error) { return decode[any](s.c.http.post(ctx, s.p("/api/v3/ScanForm/postUspsCreateScanFormLabelShipment"), body)) }
func (s *UspsService) DeliveryStandards(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/ServiceStandards/getUspsGetDeliveryStandardsEstimates"), toValues(params))) }
func (s *UspsService) FindDropOffLocations(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/LocationSearch/getUspsFindValidDropOffLocations"), toValues(params))) }
func (s *UspsService) FindPostOffices(ctx context.Context, params map[string]string) (any, error) { return decode[any](s.c.http.get(ctx, s.p("/api/v3/LocationSearch/getUspsFindValidPostOfficeLocations"), toValues(params))) }
