package flexops

import (
	"context"
	"net/url"
	"strconv"
)

type ReturnsService struct{ c *Client }

// List returns the workspace's returns, filtered by return date. Unset options are not
// sent (the Gateway defaults to page 1 of 25).
func (s *ReturnsService) List(ctx context.Context, opts ReturnListOptions) (ApiResponse[[]ReturnAuthorization], error) {
	params := url.Values{}
	if opts.StartDate != "" {
		params.Set("startDate", opts.StartDate)
	}
	if opts.EndDate != "" {
		params.Set("endDate", opts.EndDate)
	}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PageSize > 0 {
		params.Set("pageSize", strconv.Itoa(opts.PageSize))
	}
	return decode[ApiResponse[[]ReturnAuthorization]](s.c.http.get(ctx, s.c.wsPath("returns"), params))
}

// Get returns one return by its numeric return ID.
func (s *ReturnsService) Get(ctx context.Context, returnID string) (ApiResponse[ReturnAuthorization], error) {
	return decode[ApiResponse[ReturnAuthorization]](s.c.http.get(ctx, s.c.wsPath("returns/"+returnID), nil))
}

// Create opens an RMA for an existing order (the Gateway answers 201 with the RMA).
func (s *ReturnsService) Create(ctx context.Context, req ReturnRequest) (ApiResponse[Rma], error) {
	return decode[ApiResponse[Rma]](s.c.http.post(ctx, s.c.wsPath("returns/rma"), req))
}

// GenerateLabel sends {}: the Gateway binds a required GenerateReturnLabelRequest body and refuses an empty one with 400.
// An empty body carries no package, so it only works for a workspace with a default shipping preset; otherwise the
// Gateway answers 400 PackageMeasurementsRequired. Prefer GenerateLabelWithOptions and send the weight and size.
func (s *ReturnsService) GenerateLabel(ctx context.Context, rmaNumber string) (ApiResponse[any], error) {
	return s.GenerateLabelWithOptions(ctx, rmaNumber, ReturnLabelOptions{})
}

// GenerateLabelWithOptions generates a return label for an RMA with the package weight
// and size (WeightOz, LengthIn, WidthIn, HeightIn) and optional fields such as a declared
// HAZMAT class. Unset fields are not sent; CarrierCode defaults to "usps" on the Gateway.
func (s *ReturnsService) GenerateLabelWithOptions(ctx context.Context, rmaNumber string, opts ReturnLabelOptions) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("returns/"+rmaNumber+"/label"), opts))
}

// MarkReceived records the inspected items of a received RMA, identified by its numeric RMA ID.
func (s *ReturnsService) MarkReceived(ctx context.Context, rmaID string, req ReceiveReturnRequest) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("returns/"+rmaID+"/receive"), req))
}
