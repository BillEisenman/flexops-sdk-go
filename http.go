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
	defaultBaseURL = "https://gateway.flexops.io"
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
)

var retryableStatuses = map[int]bool{429: true, 500: true, 502: true, 503: true, 504: true}

type httpClient struct {
	baseURL     string
	timeout     time.Duration
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
	return &httpClient{
		baseURL:     strings.TrimRight(baseURL, "/"),
		timeout:     timeout,
		apiKey:      cfg.APIKey,
		accessToken: cfg.AccessToken,
		client:      &http.Client{Timeout: timeout},
	}
}

func (h *httpClient) setAccessToken(token string) { h.accessToken = token; h.apiKey = "" }
func (h *httpClient) setAPIKey(key string)         { h.apiKey = key; h.accessToken = "" }

func (h *httpClient) get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	return h.do(ctx, http.MethodGet, path, query, nil)
}

func (h *httpClient) post(ctx context.Context, path string, body any) ([]byte, error) {
	return h.do(ctx, http.MethodPost, path, nil, body)
}

func (h *httpClient) put(ctx context.Context, path string, body any) ([]byte, error) {
	return h.do(ctx, http.MethodPut, path, nil, body)
}

func (h *httpClient) patch(ctx context.Context, path string, body any) ([]byte, error) {
	return h.do(ctx, http.MethodPatch, path, nil, body)
}

func (h *httpClient) del(ctx context.Context, path string) ([]byte, error) {
	return h.do(ctx, http.MethodDelete, path, nil, nil)
}

func (h *httpClient) do(ctx context.Context, method, path string, query url.Values, body any) ([]byte, error) {
	u := h.buildURL(path, query)
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(h.backoff(attempt))
		}

		var bodyReader io.Reader
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(b)
		}

		req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if h.apiKey != "" {
			req.Header.Set("X-Api-Key", h.apiKey)
		} else if h.accessToken != "" {
			req.Header.Set("Authorization", "Bearer "+h.accessToken)
		}

		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		if resp.StatusCode == 401 {
			return nil, &AuthError{FlexOpsError{StatusCode: 401, Code: "UNAUTHORIZED", Message: "Authentication required"}}
		}
		if resp.StatusCode == 403 {
			return nil, &FlexOpsError{StatusCode: 403, Code: "FORBIDDEN", Message: "Access denied"}
		}
		if resp.StatusCode == 429 {
			ra, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
			lastErr = &RateLimitError{FlexOpsError: FlexOpsError{StatusCode: 429, Message: fmt.Sprintf("Rate limited, retry after %ds", ra)}, RetryAfter: ra}
			if retryableStatuses[429] {
				continue
			}
			return nil, lastErr
		}

		var errBody struct {
			Message string   `json:"message"`
			Errors  []string `json:"errors"`
		}
		_ = json.Unmarshal(respBody, &errBody)
		msg := errBody.Message
		if msg == "" {
			msg = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		apiErr := &FlexOpsError{StatusCode: resp.StatusCode, Message: msg, Errors: errBody.Errors}

		if retryableStatuses[resp.StatusCode] {
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
	ms := math.Min(1000*math.Pow(2, float64(attempt-1))*jitter, 30000)
	return time.Duration(ms) * time.Millisecond
}

// decode is a helper that unmarshals JSON response into target.
func decode[T any](data []byte, err error) (T, error) {
	var result T
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}
	return result, nil
}
