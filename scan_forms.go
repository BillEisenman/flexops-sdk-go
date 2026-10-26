package flexops

import "context"

type ScanFormsService struct{ c *Client }

// Create creates a USPS SCAN form. API keys need shipping:write; the public demo keys are
// refused. It sends req.RequestID, or a UUID generated for this call, on every attempt, so a
// retried call replays the same operation instead of creating a second form.
//
// It returns the created form, also when the same RequestID and inputs are sent again. It
// returns a *FlexOpsError with StatusCode 409 when the RequestID was used with different
// inputs, when the packages already belong to a manifest, or when this or an earlier operation
// needs review (its Message names the operation; do not resubmit those packages), and 400 when
// the request or package selection is refused. A 409 is not retried.
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

// List lists the workspace's saved scan forms, newest mailing date first, including any whose
// outcome is unknown.
func (s *ScanFormsService) List(ctx context.Context) (ApiResponse[[]ScanForm], error) {
	return decode[ApiResponse[[]ScanForm]](s.c.http.get(ctx, s.c.wsPath("scan-forms"), nil))
}

// Get gets a saved scan form by its ScanFormID. A form that does not exist is a *FlexOpsError
// with StatusCode 404.
func (s *ScanFormsService) Get(ctx context.Context, scanFormID string) (ApiResponse[ScanForm], error) {
	return decode[ApiResponse[ScanForm]](s.c.http.get(ctx, s.c.wsPath("scan-forms/"+scanFormID), nil))
}
