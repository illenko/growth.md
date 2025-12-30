package service

import (
	"context"
	"fmt"
	"time"

	"github.com/illenko/growth.md/internal/ai"
	"github.com/illenko/growth.md/internal/aifactory"
	"github.com/illenko/growth.md/internal/core"
	"github.com/illenko/growth.md/internal/storage"
)

type AIService struct {
	config        *storage.Config
	skillRepo     *storage.SkillRepository
	goalRepo      *storage.GoalRepository
	pathRepo      *storage.PathRepository
	phaseRepo     *storage.PhaseRepository
	resourceRepo  *storage.ResourceRepository
	milestoneRepo *storage.MilestoneRepository
	progressRepo  *storage.ProgressLogRepository
}

func NewAIService(
	config *storage.Config,
	skillRepo *storage.SkillRepository,
	goalRepo *storage.GoalRepository,
	pathRepo *storage.PathRepository,
	phaseRepo *storage.PhaseRepository,
	resourceRepo *storage.ResourceRepository,
	milestoneRepo *storage.MilestoneRepository,
	progressRepo *storage.ProgressLogRepository,
) *AIService {
	return &AIService{
		config:        config,
		skillRepo:     skillRepo,
		goalRepo:      goalRepo,
		pathRepo:      pathRepo,
		phaseRepo:     phaseRepo,
		resourceRepo:  resourceRepo,
		milestoneRepo: milestoneRepo,
		progressRepo:  progressRepo,
	}
}

type PathGenerationOptions struct {
	GoalID         core.EntityID
	Style          string
	TimeCommitment string
	Background     string
	Provider       string
	Model          string
}

type PathGenerationResult struct {
	Path       *core.LearningPath
	Phases     []*core.Phase
	Resources  []*core.Resource
	Milestones []*core.Milestone
	Reasoning  string
}

func (s *AIService) GenerateLearningPath(ctx context.Context, opts PathGenerationOptions) (*PathGenerationResult, error) {
	goal, err := s.goalRepo.GetByIDWithBody(opts.GoalID)
	if err != nil {
		return nil, fmt.Errorf("goal '%s' not found: %w", opts.GoalID, err)
	}

	skills, err := s.skillRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load skills: %w", err)
	}

	provider := s.config.AI.Provider
	if opts.Provider != "" {
		provider = opts.Provider
	}

	model := s.config.AI.Model
	if opts.Model != "" {
		model = opts.Model
	}

	style := s.config.AI.DefaultStyle
	if opts.Style != "" {
		style = opts.Style
	}

	aiConfig := ai.Config{
		Provider:    provider,
		Model:       model,
		Temperature: s.config.AI.Temperature,
		MaxTokens:   s.config.AI.MaxTokens,
	}

	if err := aiConfig.Validate(); err != nil {
		return nil, fmt.Errorf("AI configuration error: %w", err)
	}

	client, err := aifactory.NewClient(aiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI client: %w", err)
	}

	req := ai.PathGenerationRequest{
		Goal:           goal,
		CurrentSkills:  skills,
		Background:     opts.Background,
		LearningStyle:  style,
		TimeCommitment: opts.TimeCommitment,
		TargetDate:     goal.TargetDate,
	}

	resp, err := client.GenerateLearningPath(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate path: %w", err)
	}

	resp.Path.GenerationContext = fmt.Sprintf("Goal: %s | Style: %s | Time: %s",
		goal.Title, style, opts.TimeCommitment)

	return &PathGenerationResult{
		Path:       resp.Path,
		Phases:     resp.Phases,
		Resources:  resp.Resources,
		Milestones: resp.Milestones,
		Reasoning:  resp.Reasoning,
	}, nil
}

func (s *AIService) SaveGeneratedPath(result *PathGenerationResult, goalID core.EntityID) error {
	if err := s.pathRepo.Create(result.Path); err != nil {
		return fmt.Errorf("failed to save path: %w", err)
	}

	for _, phase := range result.Phases {
		if err := s.phaseRepo.Create(phase); err != nil {
			return fmt.Errorf("failed to save phase: %w", err)
		}
	}

	for _, resource := range result.Resources {
		if err := s.resourceRepo.Create(resource); err != nil {
			return fmt.Errorf("failed to save resource: %w", err)
		}
	}

	for _, milestone := range result.Milestones {
		if err := s.milestoneRepo.Create(milestone); err != nil {
			return fmt.Errorf("failed to save milestone: %w", err)
		}
	}

	goal, err := s.goalRepo.GetByIDWithBody(goalID)
	if err != nil {
		return fmt.Errorf("failed to load goal: %w", err)
	}

	goal.AddLearningPath(result.Path.ID)
	if err := s.goalRepo.Update(goal); err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}

	return nil
}

type ResourceSuggestionOptions struct {
	SkillID     core.EntityID
	TargetLevel core.ProficiencyLevel
	Style       string
	Budget      string
	Provider    string
	Model       string
}

type ResourceSuggestionResult struct {
	Resources []*core.Resource
	Reasoning string
}

func (s *AIService) SuggestResources(ctx context.Context, opts ResourceSuggestionOptions) (*ResourceSuggestionResult, error) {
	skill, err := s.skillRepo.GetByIDWithBody(opts.SkillID)
	if err != nil {
		return nil, fmt.Errorf("skill '%s' not found: %w", opts.SkillID, err)
	}

	currentLevel := skill.Level
	targetLevel := opts.TargetLevel
	if targetLevel == "" {
		targetLevel = getNextLevel(currentLevel)
	}

	provider := s.config.AI.Provider
	if opts.Provider != "" {
		provider = opts.Provider
	}

	model := s.config.AI.Model
	if opts.Model != "" {
		model = opts.Model
	}

	style := s.config.AI.DefaultStyle
	if opts.Style != "" {
		style = opts.Style
	}

	budget := s.config.AI.DefaultBudget
	if opts.Budget != "" {
		budget = opts.Budget
	}

	aiConfig := ai.Config{
		Provider:    provider,
		Model:       model,
		Temperature: s.config.AI.Temperature,
		MaxTokens:   s.config.AI.MaxTokens,
	}

	if err := aiConfig.Validate(); err != nil {
		return nil, fmt.Errorf("AI configuration error: %w", err)
	}

	client, err := aifactory.NewClient(aiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI client: %w", err)
	}

	req := ai.ResourceSuggestionRequest{
		Skill:         skill,
		CurrentLevel:  currentLevel,
		TargetLevel:   targetLevel,
		LearningStyle: style,
		Budget:        budget,
	}

	resp, err := client.SuggestResources(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest resources: %w", err)
	}

	return &ResourceSuggestionResult{
		Resources: resp.Resources,
		Reasoning: resp.Reasoning,
	}, nil
}

type ProgressAnalysisOptions struct {
	GoalID   core.EntityID
	Days     int
	Provider string
	Model    string
}

type ProgressAnalysisResult struct {
	Summary         string
	Insights        []string
	Recommendations []string
	IsOnTrack       bool
	SuggestedFocus  []string
	LogCount        int
}

func (s *AIService) AnalyzeProgress(ctx context.Context, opts ProgressAnalysisOptions) (*ProgressAnalysisResult, error) {
	var goal *core.Goal
	var path *core.LearningPath
	var err error

	if opts.GoalID != "" {
		goal, err = s.goalRepo.GetByIDWithBody(opts.GoalID)
		if err != nil {
			return nil, fmt.Errorf("goal '%s' not found: %w", opts.GoalID, err)
		}

		if len(goal.LearningPaths) > 0 {
			path, _ = s.pathRepo.GetByIDWithBody(goal.LearningPaths[0])
		}
	}

	cutoffDate := time.Now().AddDate(0, 0, -opts.Days)
	logs, err := s.progressRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load progress logs: %w", err)
	}

	var recentLogs []*core.ProgressLog
	for _, log := range logs {
		if log.Date.After(cutoffDate) {
			recentLogs = append(recentLogs, log)
		}
	}

	skills, err := s.skillRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load skills: %w", err)
	}

	provider := s.config.AI.Provider
	if opts.Provider != "" {
		provider = opts.Provider
	}

	model := s.config.AI.Model
	if opts.Model != "" {
		model = opts.Model
	}

	aiConfig := ai.Config{
		Provider:    provider,
		Model:       model,
		Temperature: s.config.AI.Temperature,
		MaxTokens:   s.config.AI.MaxTokens,
	}

	if err := aiConfig.Validate(); err != nil {
		return nil, fmt.Errorf("AI configuration error: %w", err)
	}

	client, err := aifactory.NewClient(aiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI client: %w", err)
	}

	req := ai.ProgressAnalysisRequest{
		Goal:          goal,
		Path:          path,
		ProgressLogs:  recentLogs,
		CurrentSkills: skills,
	}

	resp, err := client.AnalyzeProgress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze progress: %w", err)
	}

	return &ProgressAnalysisResult{
		Summary:         resp.Summary,
		Insights:        resp.Insights,
		Recommendations: resp.Recommendations,
		IsOnTrack:       resp.IsOnTrack,
		SuggestedFocus:  resp.SuggestedFocus,
		LogCount:        len(recentLogs),
	}, nil
}

func getNextLevel(current core.ProficiencyLevel) core.ProficiencyLevel {
	switch current {
	case core.LevelBeginner:
		return core.LevelIntermediate
	case core.LevelIntermediate:
		return core.LevelAdvanced
	case core.LevelAdvanced:
		return core.LevelExpert
	case core.LevelExpert:
		return core.LevelExpert
	default:
		return core.LevelIntermediate
	}
}
