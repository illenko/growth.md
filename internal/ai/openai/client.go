package openai

import (
	"context"
	"fmt"

	"github.com/illenko/growth.md/internal/ai"
)

// Client is a stub implementation for OpenAI (not yet implemented)
type Client struct {
	config ai.Config
}

// NewClient creates a new OpenAI client (stub)
func NewClient(cfg ai.Config) (*Client, error) {
	return &Client{
		config: cfg,
	}, nil
}

// Provider returns the provider name
func (c *Client) Provider() string {
	return "openai"
}

// GenerateLearningPath is not yet implemented
func (c *Client) GenerateLearningPath(ctx context.Context, req ai.PathGenerationRequest) (*ai.PathGenerationResponse, error) {
	return nil, fmt.Errorf("OpenAI provider: %w (coming soon)", ai.ErrProviderNotSupported)
}

// SuggestResources is not yet implemented
func (c *Client) SuggestResources(ctx context.Context, req ai.ResourceSuggestionRequest) (*ai.ResourceSuggestionResponse, error) {
	return nil, fmt.Errorf("OpenAI provider: %w (coming soon)", ai.ErrProviderNotSupported)
}

// AnalyzeProgress is not yet implemented
func (c *Client) AnalyzeProgress(ctx context.Context, req ai.ProgressAnalysisRequest) (*ai.ProgressAnalysisResponse, error) {
	return nil, fmt.Errorf("OpenAI provider: %w (coming soon)", ai.ErrProviderNotSupported)
}
