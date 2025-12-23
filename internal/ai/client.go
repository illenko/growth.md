package ai

import (
	"context"
)

// AIClient is the main interface for AI providers
type AIClient interface {
	// GenerateLearningPath creates a personalized learning path from a goal
	GenerateLearningPath(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error)

	// SuggestResources recommends learning resources for a skill
	SuggestResources(ctx context.Context, req ResourceSuggestionRequest) (*ResourceSuggestionResponse, error)

	// AnalyzeProgress provides insights on progress and next steps
	AnalyzeProgress(ctx context.Context, req ProgressAnalysisRequest) (*ProgressAnalysisResponse, error)

	// Provider returns the name of the AI provider
	Provider() string
}
