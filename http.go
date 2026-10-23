// ***********************************************************************
// Package          : flexops-sdk-go
// Author           : FlexOps, LLC
// Created          : 2026-03-08
//
// Copyright (c) 2021-2026 by FlexOps, LLC. All rights reserved.
// ***********************************************************************

package flexops

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://gateway.flexops.io"
	defaultTimeout   = 30 * time.Second
	defaultBaseDelay = time.Second
	maxRetries       = 3

	idempotencyKeyHeader = "Idempotency-Key"
	// idempotencyConflict is the Gateway's 409 code while the first request with the
	// same Idempotency-Key is still running.
	idempotencyConflict = "IDEMPOTENCY_CONFLICT"
)

// retryableStatuses are the statuses a GET retries on.
var retryableStatuses = map[int]bool{429: true, 500: true, 502: true, 503: true, 504: true}

// writeRetryableStatuses are the statuses where a keyed POST, PUT or PATCH either never
// ran or is replayed by the Gateway. A 500 is not one: the Gateway does not cache a
// response >= 500, so a retry would run the write (and charge the wallet) again.
var writeRetryableStatuses = map[int]bool{429: true, 502: true, 503: true, 504: true}

type httpClient struct {
	baseURL     string
	timeout     time.Duration
	maxRetries  int
	baseDelay   time.Duration
	apiKey      string
	accessToken string
	client      *http.Client
}

func newHTTPClient(cfg Config) *httpClient {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	retries := cfg.MaxRetries
	if retries <= 0 {
		retries = maxRetries
	}
	return &httpClient{
		baseURL:     strings.TrimRight(baseURL, "/"),
		timeout:     timeout,
		maxRetries:  retries,
		baseDelay:   defaultBaseDelay,
		apiKey:      cfg.APIKey,
		accessToken: cfg.AccessToken,
		client:      &http.Client{Timeout: timeout},
	}
}

func (h *httpClient) setAccessToken(token string) { h.accessToken = token; h.apiKey = "" }
func (h *httpClient) setAPIKey(key string)        { h.apiKey = key; h.accessToken = "" }

type idempotencyKeyContextKey struct{}

// WithIdempotencyKey returns a context that makes the next POST, PUT or PATCH sent with
// it carry key as its Idempotency-Key header. Without one, the SDK generates a random
// key per call and reuses it on that call's retries. Supply your own key when you
// persist it to reconcile an uncertain purchase (an insurance policy, a carbon offset);
// reuse a key only to retry the same request.
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKeyContextKey{}, key)
}

func (h *httpClient) get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	return h.do(ctx, http.MethodGet, path, query, nil)
}

func (h *httpClient) post(ctx context.Context, path string, body any) ([]byte, error) {
	return h.do(ctx, http.MethodPost, path, nil, body)
}

func (h *httpClient) put(ctx context.Context, path string, body any) ([]byte, error) {
	return h.do(ctx, http.MethodPut, path, nil, body)
}

func (h *httpClient) del(ctx context.Context, path string) ([]byte, error) {
	return h.do(ctx, http.MethodDelete, path, nil, nil)
}

// do sends one call. Every POST, PUT and PATCH carries an Idempotency-Key (the caller's,
// from extraHeaders or WithIdempotencyKey, else one random UUID for the call) that is
// reused on each retry, so the Gateway replays a write it already finished instead of
// running it again. GETs retry on a network error, 429 and 500/502/503/504; writes on a
// network error, 429, 502/503/504 and a 409 IDEMPOTENCY_CONFLICT, never on 500; a
// DELETE, which the Gateway does not deduplicate, only on 429. A 2xx with no body
// returns nil data.
func (h *httpClient) do(ctx context.Context, method, path string, query url.Values, body any, extraHeaders ...http.Header) ([]byte, error) {
	u := h.buildURL(path, query)
	write := method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
	idempotencyKey := ""
	if write {
		idempotencyKey = callerIdempotencyKey(ctx, extraHeaders)
		if idempotencyKey == "" {
			key, err := newIdempotencyKey()
			if err != nil {
				return nil, err
			}
			idempotencyKey = key
		}
	}
	var payload []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = b
	}
	var lastErr error

	for attempt := 0; attempt <= h.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(h.backoff(attempt))
		}

		var bodyReader io.Reader
		if payload != nil {
			bodyReader = bytes.NewReader(payload)
		}

		req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
		if err != nil {
			return nil, err
		}
		for _, headers := range extraHeaders {
			for name, values := range headers {
				req.Header[name] = append([]string(nil), values...)
			}
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if h.apiKey != "" {
			req.Header.Set("X-Api-Key", h.apiKey)
		} else if h.accessToken != "" {
			req.Header.Set("Authorization", "Bearer "+h.accessToken)
		}
		if write {
			req.Header.Set(idempotencyKeyHeader, idempotencyKey)
		}

		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = err
			// A canceled call is not retried, and a DELETE that may have run is not resent.
			if ctx.Err() != nil || method == http.MethodDelete {
				return nil, err
			}
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if len(respBody) == 0 {
				return nil, nil
			}
			return respBody, nil
		}

		if resp.StatusCode == 401 {
			return nil, &AuthError{FlexOpsError{StatusCode: 401, Code: "UNAUTHORIZED", Message: "Authentication required"}}
		}
		if resp.StatusCode == 429 {
			ra, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
			// Rate limited: the request was refused before it ran, so every method retries.
			lastErr = &RateLimitError{FlexOpsError: FlexOpsError{StatusCode: 429, Message: fmt.Sprintf("Rate limited, retry after %ds", ra)}, RetryAfter: ra}
			continue
		}

		var errBody struct {
			Code      string          `json:"code"`
			ErrorCode string          `json:"errorCode"`
			Message   string          `json:"message"`
			Errors    []string        `json:"errors"`
			Data      json.RawMessage `json:"data"`
		}
		_ = json.Unmarshal(respBody, &errBody)
		if errBody.ErrorCode != "" {
			errBody.Code = errBody.ErrorCode
		}
		msg := errBody.Message
		if msg == "" {
			msg = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		if resp.StatusCode == 403 && errBody.Code == "" {
			errBody.Code = "FORBIDDEN"
		}
		apiErr := &FlexOpsError{StatusCode: resp.StatusCode, Code: errBody.Code, Message: msg, Errors: errBody.Errors}
		if len(errBody.Data) > 0 && string(errBody.Data) != "null" {
			apiErr.Data = errBody.Data
		}

		if retryable(method, resp.StatusCode, errBody.Code) {
			lastErr = apiErr
			continue
		}
		return nil, apiErr
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &FlexOpsError{Message: "Request failed after retries", Code: "RETRY_EXHAUSTED"}
}

// retryable reports whether a non-2xx, non-401, non-429 response may be sent again.
func retryable(method string, status int, code string) bool {
	switch method {
	case http.MethodGet, http.MethodHead:
		return retryableStatuses[status]
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if status == http.StatusConflict {
			return code == idempotencyConflict
		}
		return writeRetryableStatuses[status]
	default:
		return false
	}
}

// callerIdempotencyKey returns the key the caller supplied, if any: an explicit header
// wins over WithIdempotencyKey.
func callerIdempotencyKey(ctx context.Context, extraHeaders []http.Header) string {
	for _, headers := range extraHeaders {
		if key := headers.Get(idempotencyKeyHeader); key != "" {
			return key
		}
	}
	if key, ok := ctx.Value(idempotencyKeyContextKey{}).(string); ok {
		return key
	}
	return ""
}

// newIdempotencyKey returns a random UUID version 4.
func newIdempotencyKey() (string, error) {
	var b [16]byte
	if _, err := crand.Read(b[:]); err != nil {
		return "", fmt.Errorf("flexops: generate idempotency key: %w", err)
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func (h *httpClient) buildURL(path string, query url.Values) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := h.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	return u
}

func (h *httpClient) backoff(attempt int) time.Duration {
	jitter := 0.85 + rand.Float64()*0.3
	d := math.Min(float64(h.baseDelay)*math.Pow(2, float64(attempt-1))*jitter, float64(30*time.Second))
	return time.Duration(d)
}

// decode is a helper that unmarshals JSON response into target. A response with no
// body (204, or an empty 2xx) returns the zero value.
func decode[T any](data []byte, err error) (T, error) {
	var result T
	if err != nil {
		return result, err
	}
	if len(data) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}
	return result, nil
}
