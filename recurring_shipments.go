package flexops

import "context"

// RecurringShipmentsService schedules shipments that repeat on an interval. Each run
// needs a destination address (AddressID); a shipping preset supplies the package and
// service. The Gateway returns these results without an envelope.
//
// Reads work with any API key or a user token. Writes need a user token (Auth.Login)
// or an API key granted the recurring-shipments:write scope; other keys get 403
// API_KEY_SCOPE_REQUIRED.
type RecurringShipmentsService struct{ c *Client }

func (s *RecurringShipmentsService) List(ctx context.Context) ([]RecurringShipment, error) {
	return decode[[]RecurringShipment](s.c.http.get(ctx, s.c.wsPath("recurring-shipments"), nil))
}
func (s *RecurringShipmentsService) Get(ctx context.Context, id string) (RecurringShipment, error) {
	return decode[RecurringShipment](s.c.http.get(ctx, s.c.wsPath("recurring-shipments/"+id), nil))
}
func (s *RecurringShipmentsService) Create(ctx context.Context, req CreateRecurringShipmentRequest) (RecurringShipment, error) {
	return decode[RecurringShipment](s.c.http.post(ctx, s.c.wsPath("recurring-shipments"), req))
}

// Update changes a schedule; unset fields are left as they are.
func (s *RecurringShipmentsService) Update(ctx context.Context, id string, req UpdateRecurringShipmentRequest) (RecurringShipment, error) {
	return decode[RecurringShipment](s.c.http.put(ctx, s.c.wsPath("recurring-shipments/"+id), req))
}
func (s *RecurringShipmentsService) Delete(ctx context.Context, id string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("recurring-shipments/"+id))
	return err
}
func (s *RecurringShipmentsService) Pause(ctx context.Context, id string) error {
	_, err := s.c.http.post(ctx, s.c.wsPath("recurring-shipments/"+id+"/pause"), nil)
	return err
}
func (s *RecurringShipmentsService) Resume(ctx context.Context, id string) error {
	_, err := s.c.http.post(ctx, s.c.wsPath("recurring-shipments/"+id+"/resume"), nil)
	return err
}

// Trigger runs a schedule now, which creates an order. To reconcile an uncertain
// outcome, pass a stored key with WithIdempotencyKey and reuse it only for this run.
func (s *RecurringShipmentsService) Trigger(ctx context.Context, id string) error {
	_, err := s.c.http.post(ctx, s.c.wsPath("recurring-shipments/"+id+"/trigger"), nil)
	return err
}
