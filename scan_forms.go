package flexops

import "context"

type ScanFormsService struct{ c *Client }

func (s *ScanFormsService) Create(ctx context.Context, req ScanFormRequest) (ApiResponse[ScanForm], error) { return decode[ApiResponse[ScanForm]](s.c.http.post(ctx, s.c.wsPath("scan-forms"), req)) }
func (s *ScanFormsService) List(ctx context.Context) (ApiResponse[[]ScanForm], error) { return decode[ApiResponse[[]ScanForm]](s.c.http.get(ctx, s.c.wsPath("scan-forms"), nil)) }
func (s *ScanFormsService) Get(ctx context.Context, scanFormID string) (ApiResponse[ScanForm], error) { return decode[ApiResponse[ScanForm]](s.c.http.get(ctx, s.c.wsPath("scan-forms/"+scanFormID), nil)) }
