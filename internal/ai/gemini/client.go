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

type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
	config ai.Config
}

func NewClient(cfg ai.Config) (*Client, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, &ai.APIError{
			Provider: "gemini",
			Message:  "failed to create client",
			Err:      err,
		}
	}

	modelName := cfg.Model
	if modelName == "" {
		modelName = "gemini-3-flash-preview"
	}

	model := client.GenerativeModel(modelName)

	temperature := cfg.Temperature
	model.Temperature = &temperature

	maxTokens := int32(cfg.MaxTokens)
	model.MaxOutputTokens = &maxTokens

	model.ResponseMIMEType = "application/json"

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

func (c *Client) Provider() string {
	return "gemini"
}

func (c *Client) GenerateLearningPath(ctx context.Context, req ai.PathGenerationRequest) (*ai.PathGenerationResponse, error) {
	prompt, err := c.renderPathPrompt(req)
	if err != nil {
		return nil, err
	}

	responseText, err := c.generateWithRetry(ctx, prompt, 3)
	if err != nil {
		return nil, err
	}

	pathID := core.EntityID(fmt.Sprintf("path-%03d", time.Now().Unix()%1000))
	goalID := req.Goal.ID

	resp, err := ParsePathGeneration(responseText, pathID, goalID)
	if err != nil {
		return nil, err
	}

	resp.Path.GenerationContext = fmt.Sprintf("Goal: %s | Style: %s | Time: %s",
		req.Goal.Title, req.LearningStyle, req.TimeCommitment)

	return resp, nil
}

func (c *Client) SuggestResources(ctx context.Context, req ai.ResourceSuggestionRequest) (*ai.ResourceSuggestionResponse, error) {
	prompt, err := c.renderResourcePrompt(req)
	if err != nil {
		return nil, err
	}

	responseText, err := c.generateWithRetry(ctx, prompt, 3)
	if err != nil {
		return nil, err
	}

	resp, err := ParseResourceSuggestion(responseText, req.Skill.ID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) AnalyzeProgress(ctx context.Context, req ai.ProgressAnalysisRequest) (*ai.ProgressAnalysisResponse, error) {
	prompt, err := c.renderProgressPrompt(req)
	if err != nil {
		return nil, err
	}

	responseText, err := c.generateWithRetry(ctx, prompt, 3)
	if err != nil {
		return nil, err
	}

	resp, err := ParseProgressAnalysis(responseText)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) generateWithRetry(ctx context.Context, prompt string, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
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

			if strings.Contains(err.Error(), "rate limit") {
				continue
			}
			if strings.Contains(err.Error(), "timeout") {
				continue
			}

			return "", lastErr
		}

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

		return text, nil
	}

	if lastErr != nil {
		return "", lastErr
	}
	return "", &ai.APIError{
		Provider: "gemini",
		Message:  "max retries exceeded",
	}
}

func (c *Client) renderPrompt(promptTemplate string, data interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

func (c *Client) renderPathPrompt(req ai.PathGenerationRequest) (string, error) {
	return c.renderPrompt(PathGenerationPrompt, req)
}

func (c *Client) renderResourcePrompt(req ai.ResourceSuggestionRequest) (string, error) {
	return c.renderPrompt(ResourceSuggestionPrompt, req)
}

func (c *Client) renderProgressPrompt(req ai.ProgressAnalysisRequest) (string, error) {
	return c.renderPrompt(ProgressAnalysisPrompt, req)
}

func (c *Client) Close() error {
	return c.client.Close()
}
