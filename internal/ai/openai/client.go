package openai

import (
	"context"
	"fmt"

	"github.com/illenko/growth.md/internal/ai"
)

type Client struct {
	config ai.Config
}

func NewClient(cfg ai.Config) (*Client, error) {
	return &Client{
		config: cfg,
	}, nil
}

func (c *Client) Provider() string {
	return "openai"
}

func (c *Client) GenerateLearningPath(ctx context.Context, req ai.PathGenerationRequest) (*ai.PathGenerationResponse, error) {
	return nil, fmt.Errorf("OpenAI provider: %w (coming soon)", ai.ErrProviderNotSupported)
}

func (c *Client) SuggestResources(ctx context.Context, req ai.ResourceSuggestionRequest) (*ai.ResourceSuggestionResponse, error) {
	return nil, fmt.Errorf("OpenAI provider: %w (coming soon)", ai.ErrProviderNotSupported)
}

func (c *Client) AnalyzeProgress(ctx context.Context, req ai.ProgressAnalysisRequest) (*ai.ProgressAnalysisResponse, error) {
	return nil, fmt.Errorf("OpenAI provider: %w (coming soon)", ai.ErrProviderNotSupported)
}
