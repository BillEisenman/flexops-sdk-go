package flexops

import "context"

type ScanFormsService struct{ c *Client }

// Create creates a USPS SCAN form. API keys need shipping:write; the public demo keys are
// refused. It sends req.RequestID, or a UUID generated for this call, on every attempt, so a
// retried call replays the same operation instead of creating a second form.
func (s *ScanFormsService) Create(ctx context.Context, req ScanFormRequest) (ApiResponse[ScanForm], error) {
	if req.RequestID == "" {
		id, err := newIdempotencyKey()
		if err != nil {
			return ApiResponse[ScanForm]{}, err
		}
		req.RequestID = id
	}
	return decode[ApiResponse[ScanForm]](s.c.http.post(ctx, s.c.wsPath("scan-forms"), req))
}

func (s *ScanFormsService) List(ctx context.Context) (ApiResponse[[]ScanForm], error) {
	return decode[ApiResponse[[]ScanForm]](s.c.http.get(ctx, s.c.wsPath("scan-forms"), nil))
}
func (s *ScanFormsService) Get(ctx context.Context, scanFormID string) (ApiResponse[ScanForm], error) {
	return decode[ApiResponse[ScanForm]](s.c.http.get(ctx, s.c.wsPath("scan-forms/"+scanFormID), nil))
}
