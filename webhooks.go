package flexops

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

type WebhooksService struct{ c *Client }

func (s *WebhooksService) List(ctx context.Context) (ApiResponse[[]WebhookSubscription], error) { return decode[ApiResponse[[]WebhookSubscription]](s.c.http.get(ctx, s.c.wsPath("webhooks"), nil)) }

// Create creates a webhook subscription. req.Events is sent as the Gateway's
// comma-separated eventTypes; leave it nil to receive every event. The result's
// Data.Secret is the signing secret, returned only here. The Gateway sends the
// new subscription without an envelope; the SDK wraps it in ApiResponse (an
// enveloped body is read as is).
//
// It returns a *FlexOpsError with Code "UnknownWebhookEvent", before any request,
// when req.Events is empty but non-nil or names an event FlexOps does not send.
func (s *WebhooksService) Create(ctx context.Context, req CreateWebhookRequest) (ApiResponse[WebhookSubscription], error) {
	eventTypes, err := webhookEventTypes(req.Events)
	if err != nil {
		return ApiResponse[WebhookSubscription]{}, err
	}
	body := struct {
		URL         string `json:"url"`
		Description string `json:"description,omitempty"`
		EventTypes  string `json:"eventTypes,omitempty"`
	}{req.URL, req.Description, eventTypes}
	return decodeCreatedWebhook(s.c.http.post(ctx, s.c.wsPath("webhooks"), body))
}

// decodeCreatedWebhook reads a create response. The Gateway answers 201 with the
// bare subscription (and its secret), not the ApiResponse envelope; an enveloped
// body ({"success":...,"data":...}) is read too, so the secret is never lost.
func decodeCreatedWebhook(data []byte, err error) (ApiResponse[WebhookSubscription], error) {
	if err != nil {
		return ApiResponse[WebhookSubscription]{}, err
	}
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return ApiResponse[WebhookSubscription]{}, err
	}
	_, hasSuccess := probe["success"]
	_, hasData := probe["data"]
	if hasSuccess && hasData {
		return decode[ApiResponse[WebhookSubscription]](data, nil)
	}
	return wrap[WebhookSubscription](data, nil)
}

// Update changes a webhook subscription. req.Events, when non-nil, replaces the
// subscribed events and is sent as the Gateway's comma-separated eventTypes.
//
// It returns a *FlexOpsError with Code "UnknownWebhookEvent", before any request,
// when req.Events is empty but non-nil or names an event FlexOps does not send.
func (s *WebhooksService) Update(ctx context.Context, webhookID string, req UpdateWebhookRequest) (ApiResponse[WebhookSubscription], error) {
	eventTypes, err := webhookEventTypes(req.Events)
	if err != nil {
		return ApiResponse[WebhookSubscription]{}, err
	}
	body := struct {
		URL         string  `json:"url,omitempty"`
		Description *string `json:"description,omitempty"`
		EventTypes  string  `json:"eventTypes,omitempty"`
		Status      string  `json:"status,omitempty"`
	}{req.URL, req.Description, eventTypes, req.Status}
	return decode[ApiResponse[WebhookSubscription]](s.c.http.put(ctx, s.c.wsPath("webhooks/"+webhookID), body))
}

func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error { _, err := s.c.http.del(ctx, s.c.wsPath("webhooks/"+webhookID)); return err }
func (s *WebhooksService) RotateSecret(ctx context.Context, webhookID string) (ApiResponse[map[string]string], error) { return decode[ApiResponse[map[string]string]](s.c.http.post(ctx, s.c.wsPath("webhooks/"+webhookID+"/rotate-secret"), nil)) }
func (s *WebhooksService) ListDeliveryLogs(ctx context.Context, webhookID string) (ApiResponse[[]any], error) { return decode[ApiResponse[[]any]](s.c.http.get(ctx, s.c.wsPath("webhooks/"+webhookID+"/deliveries"), nil)) }

// Test sends a test delivery (X-FlexOps-Event "test") to the subscription's URL.
func (s *WebhooksService) Test(ctx context.Context, webhookID string) (ApiResponse[string], error) {
	return decode[ApiResponse[string]](s.c.http.post(ctx, s.c.wsPath("webhooks/"+webhookID+"/test"), nil))
}

var knownWebhookEvents = []WebhookEventType{
	WebhookEventLabelCreated,
	WebhookEventLabelCanceled,
	WebhookEventBatchCompleted,
	WebhookEventTrackingUpdated,
	WebhookEventTrackingDelivered,
}

// webhookEventTypes joins events into the Gateway's comma-separated eventTypes,
// with no spaces and no duplicates. A nil slice returns "", which the request
// omits. The Gateway ignores an events field (subscribing to every event) and
// matches each name exactly, so a name it never sends would leave the
// subscription silent; it is refused here instead.
func webhookEventTypes(events []WebhookEventType) (string, error) {
	if events == nil {
		return "", nil
	}
	if len(events) == 0 {
		return "", &FlexOpsError{StatusCode: 400, Code: "UnknownWebhookEvent",
			Message: "Events must name at least one event; leave it nil to receive every event."}
	}
	names := make([]string, 0, len(events))
	var unknown []string
	for _, event := range events {
		if !slices.Contains(knownWebhookEvents, event) {
			unknown = append(unknown, strconv.Quote(string(event)))
			continue
		}
		if !slices.Contains(names, string(event)) {
			names = append(names, string(event))
		}
	}
	if len(unknown) > 0 {
		known := make([]string, len(knownWebhookEvents))
		for i, event := range knownWebhookEvents {
			known[i] = string(event)
		}
		return "", &FlexOpsError{StatusCode: 400, Code: "UnknownWebhookEvent",
			Message: "Unknown webhook event: " + strings.Join(unknown, ", ") + ". FlexOps sends " + strings.Join(known, ", ") + "."}
	}
	return strings.Join(names, ","), nil
}

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
