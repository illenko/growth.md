package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/illenko/growth.md/internal/ai"
	"github.com/illenko/growth.md/internal/aifactory"
	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	analyzeProvider string
	analyzeModel    string
	analyzeDays     int
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [goal-id]",
	Short: "Get AI-powered progress insights",
	Long: `Analyze your learning progress and get personalized recommendations.

If a goal-id is provided, analyzes progress for that specific goal.
Otherwise, provides overall progress analysis across all goals.

Examples:
  growth analyze                  # Overall analysis
  growth analyze goal-001         # Goal-specific analysis
  growth analyze --days 60        # Analyze last 60 days
  growth analyze goal-001 --provider gemini`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVar(&analyzeProvider, "provider", "", "AI provider (gemini, openai) - defaults to config")
	analyzeCmd.Flags().StringVar(&analyzeModel, "model", "", "model override - defaults to config")
	analyzeCmd.Flags().IntVar(&analyzeDays, "days", 30, "number of days to analyze")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	var goal *core.Goal
	var path *core.LearningPath
	var err error

	// Load goal if specified
	if len(args) > 0 {
		goalID := core.EntityID(args[0])
		goal, err = goalRepo.GetByIDWithBody(goalID)
		if err != nil {
			return fmt.Errorf("goal '%s' not found: %w", goalID, err)
		}

		// Load associated learning path if exists
		if len(goal.LearningPaths) > 0 {
			path, err = pathRepo.GetByIDWithBody(goal.LearningPaths[0])
			if err != nil {
				// Non-fatal: can analyze without path
				PrintWarning(fmt.Sprintf("Could not load learning path: %v", err))
			}
		}
	}

	// Load recent progress logs
	allProgress, err := progressRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load progress logs: %w", err)
	}

	// Filter progress logs by date
	cutoffDate := time.Now().AddDate(0, 0, -analyzeDays)
	var recentProgress []*core.ProgressLog
	for _, p := range allProgress {
		if p.Date.After(cutoffDate) {
			recentProgress = append(recentProgress, p)
		}
	}

	if len(recentProgress) == 0 {
		return fmt.Errorf("no progress logs found in the last %d days", analyzeDays)
	}

	// Load current skills
	skills, err := skillRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}

	// Initialize AI client - use config defaults, allow flags to override
	provider := config.AI.Provider
	if analyzeProvider != "" {
		provider = analyzeProvider
	}

	model := config.AI.Model
	if analyzeModel != "" {
		model = analyzeModel
	}

	aiConfig := ai.Config{
		Provider:    provider,
		Model:       model,
		Temperature: config.AI.Temperature,
		MaxTokens:   config.AI.MaxTokens,
	}

	if err := aiConfig.Validate(); err != nil {
		return fmt.Errorf("AI configuration error: %w", err)
	}

	client, err := aifactory.NewClient(aiConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize AI client: %w", err)
	}

	// Show progress
	fmt.Println("🤖 Progress Analysis")
	if goal != nil {
		fmt.Printf("   Goal: %s\n", goal.Title)
		if path != nil {
			fmt.Printf("   Path: %s\n", path.Title)
		}
	} else {
		fmt.Println("   Scope: Overall Progress")
	}
	fmt.Printf("   Period: Last %d days\n", analyzeDays)
	fmt.Printf("   Progress Logs: %d\n", len(recentProgress))
	fmt.Printf("   Provider: %s\n", client.Provider())
	fmt.Println()
	fmt.Println("⏳ Analyzing your learning journey...")

	// Create analysis request
	req := ai.ProgressAnalysisRequest{
		Goal:          goal,
		Path:          path,
		ProgressLogs:  recentProgress,
		CurrentSkills: skills,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.AnalyzeProgress(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to analyze progress: %w", err)
	}

	// Display analysis
	displayProgressAnalysis(resp, len(recentProgress))

	return nil
}

func displayProgressAnalysis(resp *ai.ProgressAnalysisResponse, logCount int) {
	fmt.Println()
	PrintSuccess("✨ Analysis Complete!")
	fmt.Println()

	// Summary
	fmt.Println("📊 SUMMARY")
	fmt.Printf("   %s\n", resp.Summary)
	fmt.Println()

	// On track status
	if resp.IsOnTrack {
		fmt.Println("✅ Status: On Track")
	} else {
		fmt.Println("⚠️  Status: Needs Attention")
	}
	fmt.Println()

	// Insights
	if len(resp.Insights) > 0 {
		fmt.Println("💡 KEY INSIGHTS")
		for i, insight := range resp.Insights {
			fmt.Printf("   %d. %s\n", i+1, insight)
		}
		fmt.Println()
	}

	// Recommendations
	if len(resp.Recommendations) > 0 {
		fmt.Println("🎯 RECOMMENDATIONS")
		for i, rec := range resp.Recommendations {
			fmt.Printf("   %d. %s\n", i+1, rec)
		}
		fmt.Println()
	}

	// Suggested focus
	if len(resp.SuggestedFocus) > 0 {
		fmt.Println("🔍 SUGGESTED FOCUS AREAS")
		for _, focus := range resp.SuggestedFocus {
			fmt.Printf("   • %s\n", focus)
		}
		fmt.Println()
	}

	fmt.Printf("💾 Based on %d progress log(s) from the last %d days\n", logCount, analyzeDays)
}
