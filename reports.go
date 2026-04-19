package flexops

import "context"

type ReportsService struct{ c *Client }

func (s *ReportsService) List(ctx context.Context) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("report-schedules"), nil))
}
func (s *ReportsService) Get(ctx context.Context, id string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("report-schedules/"+id), nil))
}
func (s *ReportsService) Create(ctx context.Context, req any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("report-schedules"), req))
}
func (s *ReportsService) Update(ctx context.Context, id string, req any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.put(ctx, s.c.wsPath("report-schedules/"+id), req))
}
func (s *ReportsService) Delete(ctx context.Context, id string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("report-schedules/"+id))
	return err
}
