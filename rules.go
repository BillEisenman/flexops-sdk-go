package flexops

import "context"

type RulesService struct{ c *Client }

func (s *RulesService) List(ctx context.Context) (ApiResponse[[]ShippingRule], error) {
	return decode[ApiResponse[[]ShippingRule]](s.c.http.get(ctx, s.c.wsPath("shipping-rules"), nil))
}
func (s *RulesService) Get(ctx context.Context, ruleID string) (ApiResponse[ShippingRule], error) {
	return decode[ApiResponse[ShippingRule]](s.c.http.get(ctx, s.c.wsPath("shipping-rules/"+ruleID), nil))
}

// Create creates a shipping rule.
//
// Signed-in users only: API keys get 403 SIGNED_IN_USER_REQUIRED, whatever their scopes.
func (s *RulesService) Create(ctx context.Context, rule CreateShippingRuleRequest) (ApiResponse[ShippingRule], error) {
	return decode[ApiResponse[ShippingRule]](s.c.http.post(ctx, s.c.wsPath("shipping-rules"), rule))
}

// Update changes a rule; unset (nil) fields are left as they are.
//
// Signed-in users only: API keys get 403 SIGNED_IN_USER_REQUIRED, whatever their scopes.
func (s *RulesService) Update(ctx context.Context, ruleID string, rule UpdateShippingRuleRequest) (ApiResponse[ShippingRule], error) {
	return decode[ApiResponse[ShippingRule]](s.c.http.put(ctx, s.c.wsPath("shipping-rules/"+ruleID), rule))
}

// Delete deletes a shipping rule.
//
// Signed-in users only: API keys get 403 SIGNED_IN_USER_REQUIRED, whatever their scopes.
func (s *RulesService) Delete(ctx context.Context, ruleID string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("shipping-rules/"+ruleID))
	return err
}

// Reorder sets rule priorities and returns the reordered rules.
//
// Signed-in users only: API keys get 403 SIGNED_IN_USER_REQUIRED, whatever their scopes.
func (s *RulesService) Reorder(ctx context.Context, rules []RulePriority) (ApiResponse[[]ShippingRule], error) {
	body := struct {
		Rules []RulePriority `json:"rules"`
	}{rules}
	return decode[ApiResponse[[]ShippingRule]](s.c.http.post(ctx, s.c.wsPath("shipping-rules/reorder"), body))
}

// Test evaluates the workspace's rules against a sample shipment without changing anything.
// API keys need shipping:read; the public demo keys are refused.
func (s *RulesService) Test(ctx context.Context, req TestShippingRuleRequest) (ApiResponse[TestShippingRuleResponse], error) {
	return decode[ApiResponse[TestShippingRuleResponse]](s.c.http.post(ctx, s.c.wsPath("shipping-rules/test"), req))
}
