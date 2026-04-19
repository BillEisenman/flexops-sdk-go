package flexops

import "context"

type OffsetService struct{ c *Client }

func (s *OffsetService) Offset(ctx context.Context, labelID string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("shipping/labels/"+labelID+"/offset"), nil))
}
func (s *OffsetService) GetEmissions(ctx context.Context, labelID string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("shipping/labels/"+labelID+"/emissions"), nil))
}
func (s *OffsetService) BatchOffset(ctx context.Context, labelIDs []string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("shipping/labels/offset/batch"), map[string][]string{"labelIds": labelIDs}))
}
