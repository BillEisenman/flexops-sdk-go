package flexops

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

type WebhooksService struct{ c *Client }

func (s *WebhooksService) List(ctx context.Context) (ApiResponse[[]WebhookSubscription], error) { return decode[ApiResponse[[]WebhookSubscription]](s.c.http.get(ctx, s.c.wsPath("webhooks"), nil)) }
func (s *WebhooksService) Get(ctx context.Context, webhookID string) (ApiResponse[WebhookSubscription], error) { return decode[ApiResponse[WebhookSubscription]](s.c.http.get(ctx, s.c.wsPath("webhooks/"+webhookID), nil)) }
func (s *WebhooksService) Create(ctx context.Context, req CreateWebhookRequest) (ApiResponse[WebhookSubscription], error) { return decode[ApiResponse[WebhookSubscription]](s.c.http.post(ctx, s.c.wsPath("webhooks"), req)) }
func (s *WebhooksService) Update(ctx context.Context, webhookID string, data any) (ApiResponse[WebhookSubscription], error) { return decode[ApiResponse[WebhookSubscription]](s.c.http.put(ctx, s.c.wsPath("webhooks/"+webhookID), data)) }
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error { _, err := s.c.http.del(ctx, s.c.wsPath("webhooks/"+webhookID)); return err }
func (s *WebhooksService) RotateSecret(ctx context.Context, webhookID string) (ApiResponse[map[string]string], error) { return decode[ApiResponse[map[string]string]](s.c.http.post(ctx, s.c.wsPath("webhooks/"+webhookID+"/rotate-secret"), nil)) }
func (s *WebhooksService) ListDeliveryLogs(ctx context.Context, webhookID string) (ApiResponse[[]any], error) { return decode[ApiResponse[[]any]](s.c.http.get(ctx, s.c.wsPath("webhooks/"+webhookID+"/deliveries"), nil)) }

// Headers FlexOps sends with every webhook delivery.
const (
	WebhookEventHeader     = "X-FlexOps-Event"     // event name, e.g. "TrackingUpdated"; "test" for a test event
	WebhookEventIDHeader   = "X-FlexOps-Event-Id"  // event id
	WebhookTimestampHeader = "X-FlexOps-Timestamp" // Unix seconds, part of the signed string
	WebhookSignatureHeader = "X-FlexOps-Signature" // "sha256=" + lowercase hex HMAC-SHA256
	WebhookDeliveryHeader  = "X-FlexOps-Delivery"  // delivery attempt id
)

// DefaultWebhookTolerance is how far the X-FlexOps-Timestamp may be from the
// current time, in either direction, before a delivery is rejected.
const DefaultWebhookTolerance = 5 * time.Minute

// Errors returned by VerifyWebhookSignature. Test them with errors.Is.
var (
	ErrWebhookMissingSecret             = errors.New("flexops: webhook secret is empty")
	ErrWebhookMissingSignature          = errors.New("flexops: webhook signature header is missing")
	ErrWebhookInvalidTimestamp          = errors.New("flexops: webhook timestamp header is missing or not Unix seconds")
	ErrWebhookTimestampOutsideTolerance = errors.New("flexops: webhook timestamp is outside the tolerance")
	ErrWebhookSignatureMismatch         = errors.New("flexops: webhook signature does not match")
)

type webhookVerifyConfig struct {
	tolerance time.Duration
	now       func() time.Time
}

// WebhookVerifyOption configures VerifyWebhookSignature.
type WebhookVerifyOption func(*webhookVerifyConfig)

// WithWebhookTolerance sets how far the delivery timestamp may be from now.
// A non-positive value keeps DefaultWebhookTolerance; the check cannot be disabled.
func WithWebhookTolerance(d time.Duration) WebhookVerifyOption {
	return func(c *webhookVerifyConfig) {
		if d > 0 {
			c.tolerance = d
		}
	}
}

// WithWebhookClock replaces time.Now, for tests.
func WithWebhookClock(now func() time.Time) WebhookVerifyOption {
	return func(c *webhookVerifyConfig) {
		if now != nil {
			c.now = now
		}
	}
}

// VerifyWebhookSignature checks a FlexOps webhook delivery and returns nil when
// it is authentic and fresh.
//
// payload is the raw request body exactly as received (verify before decoding
// it). signatureHeader is the X-FlexOps-Signature value ("sha256=<hex>"; bare hex
// is also accepted), timestampHeader is the X-FlexOps-Timestamp value, and
// secret is the subscription secret ("whsec_..."), used as-is.
//
// FlexOps signs "{timestamp}.{body}" with HMAC-SHA256. The comparison is constant
// time, and a timestamp more than DefaultWebhookTolerance (or the
// WithWebhookTolerance value) from now is rejected to limit replays.
func VerifyWebhookSignature(payload []byte, signatureHeader, timestampHeader, secret string, opts ...WebhookVerifyOption) error {
	cfg := webhookVerifyConfig{tolerance: DefaultWebhookTolerance, now: time.Now}
	for _, opt := range opts {
		opt(&cfg)
	}

	if secret == "" {
		return ErrWebhookMissingSecret
	}
	if signatureHeader == "" {
		return ErrWebhookMissingSignature
	}
	if !isUnixSeconds(timestampHeader) {
		return ErrWebhookInvalidTimestamp
	}
	ts, err := strconv.ParseInt(timestampHeader, 10, 64)
	if err != nil {
		return ErrWebhookInvalidTimestamp
	}
	age := cfg.now().Unix() - ts
	if age < 0 {
		age = -age
	}
	if age > int64(cfg.tolerance/time.Second) {
		return ErrWebhookTimestampOutsideTolerance
	}

	received, err := hex.DecodeString(strings.TrimPrefix(signatureHeader, "sha256="))
	if err != nil {
		return ErrWebhookSignatureMismatch
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestampHeader))
	mac.Write([]byte{'.'})
	mac.Write(payload)
	if !hmac.Equal(mac.Sum(nil), received) {
		return ErrWebhookSignatureMismatch
	}
	return nil
}

// VerifySignature reports whether a FlexOps webhook delivery is authentic and
// within DefaultWebhookTolerance. It is VerifyWebhookSignature without the
// reason; use that when you need the error or options.
func VerifySignature(payload []byte, signatureHeader, timestampHeader, secret string) bool {
	return VerifyWebhookSignature(payload, signatureHeader, timestampHeader, secret) == nil
}

func isUnixSeconds(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
