package flexops

import (
	"context"
	"net/http"
	"net/url"
)

type PickupsService struct{ c *Client }

// Schedule schedules a carrier pickup. API keys need shipping:write; the public demo keys are refused.
func (s *PickupsService) Schedule(ctx context.Context, req SchedulePickupRequest) (ApiResponse[Pickup], error) {
	return decode[ApiResponse[Pickup]](s.c.http.post(ctx, s.c.wsPath("pickups"), req))
}

// List returns the workspace's scheduled pickups.
func (s *PickupsService) List(ctx context.Context) (ApiResponse[[]Pickup], error) {
	return decode[ApiResponse[[]Pickup]](s.c.http.get(ctx, s.c.wsPath("pickups/scheduled"), nil))
}

// Get returns a pickup. An empty carrierCode lets the Gateway default to "usps".
func (s *PickupsService) Get(ctx context.Context, pickupID, carrierCode string) (ApiResponse[Pickup], error) {
	return decode[ApiResponse[Pickup]](s.c.http.get(ctx, s.c.wsPath("pickups/"+url.PathEscape(pickupID)), carrierQuery(carrierCode)))
}

// Cancel cancels a pickup. An empty carrierCode lets the Gateway default to "usps". API keys
// need shipping:write; the public demo keys are refused.
func (s *PickupsService) Cancel(ctx context.Context, pickupID, carrierCode string) error {
	_, err := s.c.http.do(ctx, http.MethodDelete, s.c.wsPath("pickups/"+url.PathEscape(pickupID)), carrierQuery(carrierCode), nil)
	return err
}

func carrierQuery(carrierCode string) url.Values {
	if carrierCode == "" {
		return nil
	}
	return url.Values{"carrierCode": {carrierCode}}
}
