package flexops

import "context"

// ReportSchedulesService emails reports on a schedule. Results use the ApiResponse
// envelope (Delete returns nothing).
//
// Reads work with any API key or a user token. Writes need a user token (Auth.Login)
// or an API key granted the report-schedules:write scope; other keys get 403
// API_KEY_SCOPE_REQUIRED.
type ReportSchedulesService struct{ c *Client }

func (s *ReportSchedulesService) List(ctx context.Context) (ApiResponse[[]ReportSchedule], error) {
	return decode[ApiResponse[[]ReportSchedule]](s.c.http.get(ctx, s.c.wsPath("report-schedules"), nil))
}
func (s *ReportSchedulesService) Get(ctx context.Context, id string) (ApiResponse[ReportSchedule], error) {
	return decode[ApiResponse[ReportSchedule]](s.c.http.get(ctx, s.c.wsPath("report-schedules/"+id), nil))
}
func (s *ReportSchedulesService) Create(ctx context.Context, req CreateReportScheduleRequest) (ApiResponse[ReportSchedule], error) {
	return decode[ApiResponse[ReportSchedule]](s.c.http.post(ctx, s.c.wsPath("report-schedules"), req))
}

// Update changes a schedule; unset fields are left as they are.
func (s *ReportSchedulesService) Update(ctx context.Context, id string, req UpdateReportScheduleRequest) (ApiResponse[ReportSchedule], error) {
	return decode[ApiResponse[ReportSchedule]](s.c.http.put(ctx, s.c.wsPath("report-schedules/"+id), req))
}
func (s *ReportSchedulesService) Delete(ctx context.Context, id string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("report-schedules/"+id))
	return err
}

// Toggle flips the schedule's IsEnabled. The response carries no data.
func (s *ReportSchedulesService) Toggle(ctx context.Context, id string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("report-schedules/"+id+"/toggle"), nil))
}
