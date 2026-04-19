package flexops

import "context"

type EmailTemplatesService struct{ c *Client }

func (s *EmailTemplatesService) List(ctx context.Context) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("shipment-email-templates"), nil))
}
func (s *EmailTemplatesService) Get(ctx context.Context, id string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("shipment-email-templates/"+id), nil))
}
func (s *EmailTemplatesService) Create(ctx context.Context, req any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("shipment-email-templates"), req))
}
func (s *EmailTemplatesService) Update(ctx context.Context, id string, req any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.put(ctx, s.c.wsPath("shipment-email-templates/"+id), req))
}
func (s *EmailTemplatesService) Delete(ctx context.Context, id string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("shipment-email-templates/"+id))
	return err
}
func (s *EmailTemplatesService) Preview(ctx context.Context, id string, context_ any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("shipment-email-templates/"+id+"/preview"), context_))
}
