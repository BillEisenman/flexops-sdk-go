package flexops

import (
	"context"
	"net/url"
)

type RecurringShipmentsService struct{ c *Client }

func (s *RecurringShipmentsService) List(ctx context.Context, params url.Values) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("recurring-shipments"), params))
}
func (s *RecurringShipmentsService) Get(ctx context.Context, id string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("recurring-shipments/"+id), nil))
}
func (s *RecurringShipmentsService) Create(ctx context.Context, req any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("recurring-shipments"), req))
}
func (s *RecurringShipmentsService) Update(ctx context.Context, id string, req any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.put(ctx, s.c.wsPath("recurring-shipments/"+id), req))
}
func (s *RecurringShipmentsService) Delete(ctx context.Context, id string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("recurring-shipments/"+id))
	return err
}
func (s *RecurringShipmentsService) Pause(ctx context.Context, id string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("recurring-shipments/"+id+"/pause"), nil))
}
func (s *RecurringShipmentsService) Resume(ctx context.Context, id string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("recurring-shipments/"+id+"/resume"), nil))
}
func (s *RecurringShipmentsService) Trigger(ctx context.Context, id string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("recurring-shipments/"+id+"/trigger"), nil))
}
