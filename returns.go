package flexops

import (
	"context"
	"net/url"
)

type ReturnsService struct{ c *Client }

func (s *ReturnsService) List(ctx context.Context, params url.Values) (ApiResponse[[]ReturnAuthorization], error) { return decode[ApiResponse[[]ReturnAuthorization]](s.c.http.get(ctx, s.c.wsPath("returns"), params)) }
func (s *ReturnsService) Get(ctx context.Context, returnID string) (ApiResponse[ReturnAuthorization], error) { return decode[ApiResponse[ReturnAuthorization]](s.c.http.get(ctx, s.c.wsPath("returns/"+returnID), nil)) }
func (s *ReturnsService) Create(ctx context.Context, req ReturnRequest) (ApiResponse[ReturnAuthorization], error) { return decode[ApiResponse[ReturnAuthorization]](s.c.http.post(ctx, s.c.wsPath("returns"), req)) }
func (s *ReturnsService) Approve(ctx context.Context, returnID string) (ApiResponse[ReturnAuthorization], error) { return decode[ApiResponse[ReturnAuthorization]](s.c.http.post(ctx, s.c.wsPath("returns/"+returnID+"/approve"), nil)) }
func (s *ReturnsService) Reject(ctx context.Context, returnID, reason string) (ApiResponse[ReturnAuthorization], error) { return decode[ApiResponse[ReturnAuthorization]](s.c.http.post(ctx, s.c.wsPath("returns/"+returnID+"/reject"), map[string]string{"reason": reason})) }
func (s *ReturnsService) Cancel(ctx context.Context, returnID string) error { _, err := s.c.http.post(ctx, s.c.wsPath("returns/"+returnID+"/cancel"), nil); return err }
func (s *ReturnsService) GenerateLabel(ctx context.Context, returnID string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("returns/"+returnID+"/label"), nil)) }
func (s *ReturnsService) MarkReceived(ctx context.Context, returnID string, items any) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("returns/"+returnID+"/receive"), map[string]any{"items": items})) }
func (s *ReturnsService) ProcessRefund(ctx context.Context, returnID string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("returns/"+returnID+"/refund"), nil)) }
