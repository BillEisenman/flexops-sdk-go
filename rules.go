package flexops

import "context"

type RulesService struct{ c *Client }

func (s *RulesService) List(ctx context.Context) (ApiResponse[[]ShippingRule], error) { return decode[ApiResponse[[]ShippingRule]](s.c.http.get(ctx, s.c.wsPath("shipping-rules"), nil)) }
func (s *RulesService) Get(ctx context.Context, ruleID string) (ApiResponse[ShippingRule], error) { return decode[ApiResponse[ShippingRule]](s.c.http.get(ctx, s.c.wsPath("shipping-rules/"+ruleID), nil)) }
func (s *RulesService) Create(ctx context.Context, rule any) (ApiResponse[ShippingRule], error) { return decode[ApiResponse[ShippingRule]](s.c.http.post(ctx, s.c.wsPath("shipping-rules"), rule)) }
func (s *RulesService) Update(ctx context.Context, ruleID string, rule any) (ApiResponse[ShippingRule], error) { return decode[ApiResponse[ShippingRule]](s.c.http.put(ctx, s.c.wsPath("shipping-rules/"+ruleID), rule)) }
func (s *RulesService) Delete(ctx context.Context, ruleID string) error { _, err := s.c.http.del(ctx, s.c.wsPath("shipping-rules/"+ruleID)); return err }
func (s *RulesService) Reorder(ctx context.Context, ruleIDs []string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.put(ctx, s.c.wsPath("shipping-rules/reorder"), ruleIDs)) }
