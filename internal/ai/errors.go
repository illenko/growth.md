package ai

import (
	"errors"
	"fmt"
)

var (
	ErrAPIKeyMissing        = errors.New("API key is missing")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrInvalidResponse      = errors.New("invalid response from AI provider")
	ErrProviderNotSupported = errors.New("AI provider not supported")
	ErrTimeout              = errors.New("request timeout")
)

type APIError struct {
	Provider   string
	StatusCode int
	Message    string
	Err        error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s API error (%d): %s: %v", e.Provider, e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("%s API error (%d): %s", e.Provider, e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

type ParseError struct {
	Provider string
	Message  string
	Err      error
}

func (e *ParseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s parse error: %s: %v", e.Provider, e.Message, e.Err)
	}
	return fmt.Sprintf("%s parse error: %s", e.Provider, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}
