package gemini

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/illenko/growth.md/internal/ai"
	"github.com/illenko/growth.md/internal/core"
	"google.golang.org/api/option"
)

// Client implements the AI client interface for Google Gemini
type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
	config ai.Config
}

// NewClient creates a new Gemini client
func NewClient(cfg ai.Config) (*Client, error) {
	ctx := context.Background()

	// Initialize Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, &ai.APIError{
			Provider: "gemini",
			Message:  "failed to create client",
			Err:      err,
		}
	}

	// Select model (default: gemini-3-flash-preview)
	modelName := cfg.Model
	if modelName == "" {
		modelName = "gemini-3-flash-preview"
	}

	model := client.GenerativeModel(modelName)

	// Set generation config
	temperature := cfg.Temperature
	model.Temperature = &temperature

	maxTokens := int32(cfg.MaxTokens)
	model.MaxOutputTokens = &maxTokens

	// Set structured JSON output
	model.ResponseMIMEType = "application/json"

	// Set safety settings (allow all for technical content)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}

	return &Client{
		client: client,
		model:  model,
		config: cfg,
	}, nil
}

// Provider returns the provider name
func (c *Client) Provider() string {
	return "gemini"
}

// GenerateLearningPath creates a personalized learning path
func (c *Client) GenerateLearningPath(ctx context.Context, req ai.PathGenerationRequest) (*ai.PathGenerationResponse, error) {
	// Render prompt template
	prompt, err := c.renderPathPrompt(req)
	if err != nil {
		return nil, err
	}

	// Call Gemini API with retry
	responseText, err := c.generateWithRetry(ctx, prompt, 3)
	if err != nil {
		return nil, err
	}

	// Generate IDs for entities
	pathID := core.EntityID(fmt.Sprintf("path-%03d", time.Now().Unix()%1000))
	goalID := req.Goal.ID

	// Parse response
	resp, err := ParsePathGeneration(responseText, pathID, goalID)
	if err != nil {
		return nil, err
	}

	// Store generation context
	resp.Path.GenerationContext = fmt.Sprintf("Goal: %s | Style: %s | Time: %s",
		req.Goal.Title, req.LearningStyle, req.TimeCommitment)

	return resp, nil
}

// SuggestResources recommends learning resources
func (c *Client) SuggestResources(ctx context.Context, req ai.ResourceSuggestionRequest) (*ai.ResourceSuggestionResponse, error) {
	// Render prompt template
	prompt, err := c.renderResourcePrompt(req)
	if err != nil {
		return nil, err
	}

	// Call Gemini API with retry
	responseText, err := c.generateWithRetry(ctx, prompt, 3)
	if err != nil {
		return nil, err
	}

	// Parse response
	resp, err := ParseResourceSuggestion(responseText, req.Skill.ID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AnalyzeProgress provides insights on progress
func (c *Client) AnalyzeProgress(ctx context.Context, req ai.ProgressAnalysisRequest) (*ai.ProgressAnalysisResponse, error) {
	// Render prompt template
	prompt, err := c.renderProgressPrompt(req)
	if err != nil {
		return nil, err
	}

	// Call Gemini API with retry
	responseText, err := c.generateWithRetry(ctx, prompt, 3)
	if err != nil {
		return nil, err
	}

	// Parse response
	resp, err := ParseProgressAnalysis(responseText)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// generateWithRetry calls Gemini API with exponential backoff
func (c *Client) generateWithRetry(ctx context.Context, prompt string, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			lastErr = &ai.APIError{
				Provider: "gemini",
				Message:  "API call failed",
				Err:      err,
			}

			// Check if it's a retryable error
			if strings.Contains(err.Error(), "rate limit") {
				continue
			}
			if strings.Contains(err.Error(), "timeout") {
				continue
			}

			// Non-retryable error
			return "", lastErr
		}

		// Extract text from response
		if len(resp.Candidates) == 0 {
			lastErr = &ai.APIError{
				Provider: "gemini",
				Message:  "no response candidates returned",
			}
			continue
		}

		candidate := resp.Candidates[0]
		if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
			lastErr = &ai.APIError{
				Provider: "gemini",
				Message:  "empty response content",
			}
			continue
		}

		// Get text from first part
		var text string
		for _, part := range candidate.Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				text = string(txt)
				break
			}
		}

		if text == "" {
			lastErr = &ai.APIError{
				Provider: "gemini",
				Message:  "no text content in response",
			}
			continue
		}

		// Success!
		return text, nil
	}

	// All retries failed
	if lastErr != nil {
		return "", lastErr
	}
	return "", &ai.APIError{
		Provider: "gemini",
		Message:  "max retries exceeded",
	}
}

// renderPathPrompt renders the path generation prompt template
func (c *Client) renderPathPrompt(req ai.PathGenerationRequest) (string, error) {
	tmpl, err := template.New("path").Parse(PathGenerationPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// renderResourcePrompt renders the resource suggestion prompt template
func (c *Client) renderResourcePrompt(req ai.ResourceSuggestionRequest) (string, error) {
	tmpl, err := template.New("resource").Parse(ResourceSuggestionPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// renderProgressPrompt renders the progress analysis prompt template
func (c *Client) renderProgressPrompt(req ai.ProgressAnalysisRequest) (string, error) {
	tmpl, err := template.New("progress").Parse(ProgressAnalysisPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// Close closes the Gemini client
func (c *Client) Close() error {
	return c.client.Close()
}
