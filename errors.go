// ***********************************************************************
// Package          : flexops-sdk-go
// Author           : FlexOps, LLC
// Created          : 2026-03-08
//
// Copyright (c) 2021-2026 by FlexOps, LLC. All rights reserved.
// ***********************************************************************

package flexops

import "fmt"

// FlexOpsError represents an API error response.
type FlexOpsError struct {
	StatusCode int
	Code       string
	Message    string
	Errors     []string
}

func (e *FlexOpsError) Error() string {
	return fmt.Sprintf("FlexOps API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
}

// AuthError is returned for 401 responses.
type AuthError struct {
	FlexOpsError
}

// RateLimitError is returned for 429 responses.
type RateLimitError struct {
	FlexOpsError
	RetryAfter int
}
