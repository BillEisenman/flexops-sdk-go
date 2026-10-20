package flexops

import "context"

// EmailTemplatesService manages the tracking emails sent to customers. Templates are
// keyed by tracking event: Upsert creates or replaces the template for its EventType.
// The Gateway returns these results without an envelope.
//
// Reads work with any API key or a user token. Writes (including Preview) need a user
// token (Auth.Login) or an API key granted the email-templates:write scope; other keys
// get 403 API_KEY_SCOPE_REQUIRED.
type EmailTemplatesService struct{ c *Client }

func (s *EmailTemplatesService) List(ctx context.Context) ([]EmailTemplate, error) {
	return decode[[]EmailTemplate](s.c.http.get(ctx, s.c.wsPath("email-templates"), nil))
}
func (s *EmailTemplatesService) Get(ctx context.Context, templateID string) (EmailTemplate, error) {
	return decode[EmailTemplate](s.c.http.get(ctx, s.c.wsPath("email-templates/"+templateID), nil))
}
func (s *EmailTemplatesService) Upsert(ctx context.Context, req UpsertEmailTemplateRequest) (EmailTemplate, error) {
	return decode[EmailTemplate](s.c.http.put(ctx, s.c.wsPath("email-templates"), req))
}
func (s *EmailTemplatesService) Delete(ctx context.Context, templateID string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("email-templates/"+templateID))
	return err
}

// Preview renders a subject and HTML body with sample tracking data, without saving.
func (s *EmailTemplatesService) Preview(ctx context.Context, req PreviewEmailTemplateRequest) (EmailTemplatePreview, error) {
	return decode[EmailTemplatePreview](s.c.http.post(ctx, s.c.wsPath("email-templates/preview"), req))
}
