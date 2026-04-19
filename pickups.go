package flexops

import "context"

type PickupsService struct{ c *Client }

func (s *PickupsService) Schedule(ctx context.Context, req PickupRequest) (ApiResponse[PickupConfirmation], error) { return decode[ApiResponse[PickupConfirmation]](s.c.http.post(ctx, s.c.wsPath("pickups"), req)) }
func (s *PickupsService) List(ctx context.Context) (ApiResponse[[]PickupConfirmation], error) { return decode[ApiResponse[[]PickupConfirmation]](s.c.http.get(ctx, s.c.wsPath("pickups"), nil)) }
func (s *PickupsService) Get(ctx context.Context, pickupID string) (ApiResponse[PickupConfirmation], error) { return decode[ApiResponse[PickupConfirmation]](s.c.http.get(ctx, s.c.wsPath("pickups/"+pickupID), nil)) }
func (s *PickupsService) Cancel(ctx context.Context, pickupID string) error { _, err := s.c.http.del(ctx, s.c.wsPath("pickups/"+pickupID)); return err }
