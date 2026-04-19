package flexops

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type WebhooksService struct{ c *Client }

func (s *WebhooksService) List(ctx context.Context) (ApiResponse[[]WebhookSubscription], error) { return decode[ApiResponse[[]WebhookSubscription]](s.c.http.get(ctx, s.c.wsPath("webhooks"), nil)) }
func (s *WebhooksService) Get(ctx context.Context, webhookID string) (ApiResponse[WebhookSubscription], error) { return decode[ApiResponse[WebhookSubscription]](s.c.http.get(ctx, s.c.wsPath("webhooks/"+webhookID), nil)) }
func (s *WebhooksService) Create(ctx context.Context, req CreateWebhookRequest) (ApiResponse[WebhookSubscription], error) { return decode[ApiResponse[WebhookSubscription]](s.c.http.post(ctx, s.c.wsPath("webhooks"), req)) }
func (s *WebhooksService) Update(ctx context.Context, webhookID string, data any) (ApiResponse[WebhookSubscription], error) { return decode[ApiResponse[WebhookSubscription]](s.c.http.put(ctx, s.c.wsPath("webhooks/"+webhookID), data)) }
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error { _, err := s.c.http.del(ctx, s.c.wsPath("webhooks/"+webhookID)); return err }
func (s *WebhooksService) RotateSecret(ctx context.Context, webhookID string) (ApiResponse[map[string]string], error) { return decode[ApiResponse[map[string]string]](s.c.http.post(ctx, s.c.wsPath("webhooks/"+webhookID+"/rotate-secret"), nil)) }
func (s *WebhooksService) ListDeliveryLogs(ctx context.Context, webhookID string) (ApiResponse[[]any], error) { return decode[ApiResponse[[]any]](s.c.http.get(ctx, s.c.wsPath("webhooks/"+webhookID+"/deliveries"), nil)) }

// VerifySignature verifies an HMAC-SHA256 webhook signature.
func VerifySignature(payload, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
