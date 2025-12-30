package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/illenko/growth.md/internal/ai"
	"github.com/illenko/growth.md/internal/aifactory"
	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	pathType       string
	pathStatus     string
	pathTags       string
	pathTitle      string
	pathFilterType string

	// Path generate flags
	pathGenerateStyle      string
	pathGenerateTime       string
	pathGenerateBackground string
	pathGenerateProvider   string
	pathGenerateModel      string
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Manage learning paths",
	Long:  `Create, list, view, edit, and delete learning paths.`,
}

var pathCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new learning path",
	Long: `Create a new learning path with the specified title.

You can provide the title as an argument or be prompted for it.
Paths can be manual or AI-generated.

Examples:
  growth path create "Backend Development" --type manual
  growth path create "Full Stack Path" --type ai-generated --tags backend,frontend
  growth path create`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPathCreate,
}

var pathListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all learning paths",
	Long: `List all learning paths in the repository.

Optionally filter by type or status using flags.

Examples:
  growth path list
  growth path list --type manual
  growth path list --status active`,
	Aliases: []string{"ls"},
	RunE:    runPathList,
}

var pathViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View path details",
	Long: `View detailed information about a specific learning path.

The output format can be controlled with the --format flag (table, json, yaml).

Examples:
  growth path view path-001
  growth path view path-042 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runPathView,
}

var pathEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing path",
	Long: `Edit an existing learning path by ID.

You can update any field using flags.

Examples:
  growth path edit path-001 --status completed
  growth path edit path-042 --title "New Title"
  growth path edit path-001 --tags backend,devops`,
	Args: cobra.ExactArgs(1),
	RunE: runPathEdit,
}

var pathDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a path",
	Long: `Delete a learning path by ID.

This will permanently remove the path file. You'll be prompted for confirmation
before deletion.

Examples:
  growth path delete path-001
  growth path delete path-042`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runPathDelete,
}

var pathGenerateCmd = &cobra.Command{
	Use:   "generate <goal-id>",
	Short: "Generate a learning path using AI",
	Long: `Generate a personalized learning path for a goal using AI.

The AI will analyze your goal, current skills, and preferences to create
a structured learning path with phases, milestones, and resource recommendations.

Examples:
  growth path generate goal-001
  growth path generate goal-001 --style top-down --time "10 hours/week"
  growth path generate goal-001 --background "I have 5 years Python experience"
  growth path generate goal-001 --provider gemini --model gemini-3-flash-preview`,
	Args: cobra.ExactArgs(1),
	RunE: runPathGenerate,
}

func init() {
	rootCmd.AddCommand(pathCmd)
	pathCmd.AddCommand(pathCreateCmd)
	pathCmd.AddCommand(pathListCmd)
	pathCmd.AddCommand(pathViewCmd)
	pathCmd.AddCommand(pathEditCmd)
	pathCmd.AddCommand(pathDeleteCmd)
	pathCmd.AddCommand(pathGenerateCmd)

	pathCreateCmd.Flags().StringVarP(&pathType, "type", "t", "", "path type (manual, ai-generated)")
	pathCreateCmd.Flags().StringVar(&pathTags, "tags", "", "comma-separated tags")

	pathListCmd.Flags().StringVarP(&pathFilterType, "type", "t", "", "filter by type")
	pathListCmd.Flags().StringVarP(&pathStatus, "status", "s", "", "filter by status")

	pathEditCmd.Flags().StringVar(&pathTitle, "title", "", "path title")
	pathEditCmd.Flags().StringVarP(&pathStatus, "status", "s", "", "path status")
	pathEditCmd.Flags().StringVar(&pathTags, "tags", "", "comma-separated tags")

	pathGenerateCmd.Flags().StringVar(&pathGenerateStyle, "style", "", "learning style (top-down, bottom-up, project-based) - defaults to config")
	pathGenerateCmd.Flags().StringVar(&pathGenerateTime, "time", "5 hours/week", "time commitment (e.g., '10 hours/week')")
	pathGenerateCmd.Flags().StringVar(&pathGenerateBackground, "background", "", "additional background context")
	pathGenerateCmd.Flags().StringVar(&pathGenerateProvider, "provider", "", "AI provider (gemini, openai) - defaults to config")
	pathGenerateCmd.Flags().StringVar(&pathGenerateModel, "model", "", "model override - defaults to config")
}

func runPathCreate(cmd *cobra.Command, args []string) error {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = PromptStringRequired("Path title")
	}

	if pathType == "" {
		pathType = PromptSelectWithDefault(
			"Path type",
			[]string{"manual", "ai-generated"},
			"manual",
		)
	}

	pType := core.PathType(pathType)
	if !pType.IsValid() {
		return fmt.Errorf("invalid path type: %s (must be manual or ai-generated)", pathType)
	}

	id, err := GenerateNextID("path")
	if err != nil {
		return fmt.Errorf("failed to generate path ID: %w", err)
	}

	path, err := core.NewLearningPath(id, title, pType)
	if err != nil {
		return fmt.Errorf("failed to create path: %w", err)
	}

	if pathTags != "" {
		tags := strings.Split(pathTags, ",")
		for _, tag := range tags {
			path.AddTag(strings.TrimSpace(tag))
		}
	}

	description := PromptMultiline("Description (optional, press Ctrl+D or enter '.' to finish)")
	if description != "" {
		path.Body = description
	}

	if err := pathRepo.Create(path); err != nil {
		return fmt.Errorf("failed to save path: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Created path %s: %s", path.ID, path.Title))

	if verbose {
		fmt.Printf("\nPath details:\n")
		fmt.Printf("  ID: %s\n", path.ID)
		fmt.Printf("  Title: %s\n", path.Title)
		fmt.Printf("  Type: %s\n", path.Type)
		fmt.Printf("  Status: %s\n", path.Status)
		if len(path.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(path.Tags, ", "))
		}
	}

	return nil
}

func runPathList(cmd *cobra.Command, args []string) error {
	var paths []*core.LearningPath
	var err error

	if pathFilterType != "" {
		pType := core.PathType(pathFilterType)
		if !pType.IsValid() {
			return fmt.Errorf("invalid path type '%s'. Valid options: manual, ai-generated", pathFilterType)
		}
		paths, err = pathRepo.FindByType(pType)
	} else if pathStatus != "" {
		status := core.Status(pathStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status '%s'. Valid options: active, completed, archived", pathStatus)
		}
		paths, err = pathRepo.FindByStatus(status)
	} else {
		paths, err = pathRepo.GetAll()
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve paths: %w\nTry running 'growth path list' without filters to see all paths", err)
	}

	if len(paths) == 0 {
		PrintInfo("No paths found")
		return nil
	}

	return PrintOutputWithConfig(paths)
}

func runPathView(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	path, err := pathRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("path '%s' not found. Use 'growth path list' to see available paths", id)
	}

	if config.Display.OutputFormat == "table" {
		fmt.Printf("ID:       %s\n", path.ID)
		fmt.Printf("Title:    %s\n", path.Title)
		fmt.Printf("Type:     %s\n", path.Type)
		fmt.Printf("Status:   %s\n", path.Status)
		if path.GeneratedBy != "" {
			fmt.Printf("Generated By: %s\n", path.GeneratedBy)
		}
		if len(path.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(path.Tags, ", "))
		}
		if len(path.Phases) > 0 {
			fmt.Printf("Phases:   %v\n", path.Phases)
		}
		fmt.Printf("Created:  %s\n", path.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:  %s\n", path.Updated.Format("2006-01-02 15:04:05"))

		if path.Body != "" {
			fmt.Printf("\nDescription:\n%s\n", path.Body)
		}

		return nil
	}

	return PrintOutputWithConfig(path)
}

func runPathEdit(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	path, err := pathRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("path '%s' not found. Use 'growth path list' to see available paths", id)
	}

	updated := false

	if cmd.Flags().Changed("title") {
		path.Title = pathTitle
		updated = true
	}

	if cmd.Flags().Changed("status") {
		status := core.Status(pathStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status '%s'. Valid options: active, completed, archived", pathStatus)
		}
		if err := path.UpdateStatus(status); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
		updated = true
	}

	if cmd.Flags().Changed("tags") {
		path.Tags = []string{}
		if pathTags != "" {
			tags := strings.Split(pathTags, ",")
			for _, tag := range tags {
				path.AddTag(strings.TrimSpace(tag))
			}
		}
		updated = true
	}

	if !updated {
		PrintInfo("No changes specified. Use flags to update fields.")
		return nil
	}

	if err := pathRepo.Update(path); err != nil {
		return fmt.Errorf("failed to update path: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Updated path %s: %s", path.ID, path.Title))
	return nil
}

func runPathDelete(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	path, err := pathRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("path '%s' not found. Use 'growth path list' to see available paths", id)
	}

	fmt.Printf("You are about to delete:\n")
	fmt.Printf("  ID: %s\n", path.ID)
	fmt.Printf("  Title: %s\n", path.Title)
	fmt.Printf("  Type: %s\n", path.Type)
	fmt.Println()

	if !PromptConfirm("Are you sure you want to delete this path?") {
		PrintInfo("Deletion cancelled")
		return nil
	}

	if err := pathRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete path: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Deleted path %s", id))
	return nil
}

func runPathGenerate(cmd *cobra.Command, args []string) error {
	goalID := core.EntityID(args[0])

	// Load goal
	goal, err := goalRepo.GetByIDWithBody(goalID)
	if err != nil {
		return fmt.Errorf("goal '%s' not found: %w", goalID, err)
	}

	// Load current skills
	skills, err := skillRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}

	// Initialize AI client - use config defaults, allow flags to override
	provider := config.AI.Provider
	if pathGenerateProvider != "" {
		provider = pathGenerateProvider
	}

	model := config.AI.Model
	if pathGenerateModel != "" {
		model = pathGenerateModel
	}

	style := config.AI.DefaultStyle
	if pathGenerateStyle != "" {
		style = pathGenerateStyle
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
	fmt.Printf("🤖 Generating learning path for: %s\n", goal.Title)
	fmt.Printf("   Provider: %s\n", client.Provider())
	if pathGenerateModel != "" {
		fmt.Printf("   Model: %s\n", pathGenerateModel)
	}
	fmt.Printf("   Style: %s\n", style)
	fmt.Printf("   Time Commitment: %s\n", pathGenerateTime)
	fmt.Println()
	fmt.Println("⏳ Analyzing your goal and skills...")

	// Generate path
	req := ai.PathGenerationRequest{
		Goal:           goal,
		CurrentSkills:  skills,
		Background:     pathGenerateBackground,
		LearningStyle:  style,
		TimeCommitment: pathGenerateTime,
		TargetDate:     goal.TargetDate,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.GenerateLearningPath(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate path: %w", err)
	}

	// Save path and related entities
	if err := saveGeneratedPath(resp, goalID); err != nil {
		return fmt.Errorf("failed to save path: %w", err)
	}

	// Display summary
	displayPathSummary(resp)

	return nil
}

func saveGeneratedPath(resp *ai.PathGenerationResponse, goalID core.EntityID) error {
	// Generate proper sequential IDs to avoid conflicts
	if err := reassignGeneratedIDs(resp); err != nil {
		return fmt.Errorf("failed to assign IDs: %w", err)
	}

	// Save the main path
	if err := pathRepo.Create(resp.Path); err != nil {
		return fmt.Errorf("failed to save path: %w", err)
	}

	// Save phases
	for _, phase := range resp.Phases {
		if err := phaseRepo.Create(phase); err != nil {
			return fmt.Errorf("failed to save phase %s: %w", phase.ID, err)
		}
	}

	// Save resources
	for _, resource := range resp.Resources {
		if err := resourceRepo.Create(resource); err != nil {
			return fmt.Errorf("failed to save resource %s: %w", resource.ID, err)
		}
	}

	// Save milestones
	for _, milestone := range resp.Milestones {
		if err := milestoneRepo.Create(milestone); err != nil {
			return fmt.Errorf("failed to save milestone %s: %w", milestone.ID, err)
		}
	}

	// Link path to goal
	goal, err := goalRepo.GetByIDWithBody(goalID)
	if err == nil {
		goal.LearningPaths = append(goal.LearningPaths, resp.Path.ID)
		if err := goalRepo.Update(goal); err != nil {
			// Non-fatal: path is already created
			PrintWarning(fmt.Sprintf("Failed to link path to goal: %v", err))
		}
	}

	return nil
}

func reassignGeneratedIDs(resp *ai.PathGenerationResponse) error {
	// Map old IDs to new IDs for references
	milestoneIDMap := make(map[core.EntityID]core.EntityID)
	phaseIDMap := make(map[core.EntityID]core.EntityID)

	// Generate new path ID
	newPathID, err := GenerateNextID("path")
	if err != nil {
		return fmt.Errorf("failed to generate path ID: %w", err)
	}
	resp.Path.ID = newPathID

	if len(resp.Phases) > 0 {
		startPhaseID, err := GenerateNextID("phase")
		if err != nil {
			return fmt.Errorf("failed to generate phase ID: %w", err)
		}
		phaseCounter := extractIDNumber(startPhaseID)

		for _, phase := range resp.Phases {
			oldPhaseID := phase.ID
			newPhaseID := core.EntityID(fmt.Sprintf("phase-%03d", phaseCounter))
			phaseIDMap[oldPhaseID] = newPhaseID
			phase.ID = newPhaseID
			phase.PathID = newPathID
			phaseCounter++
		}
	}

	if len(resp.Resources) > 0 {
		startResourceID, err := GenerateNextID("resource")
		if err != nil {
			return fmt.Errorf("failed to generate resource ID: %w", err)
		}
		resourceCounter := extractIDNumber(startResourceID)

		for _, resource := range resp.Resources {
			newResourceID := core.EntityID(fmt.Sprintf("resource-%03d", resourceCounter))
			resource.ID = newResourceID
			resourceCounter++
		}
	}

	if len(resp.Milestones) > 0 {
		startMilestoneID, err := GenerateNextID("milestone")
		if err != nil {
			return fmt.Errorf("failed to generate milestone ID: %w", err)
		}
		milestoneCounter := extractIDNumber(startMilestoneID)

		for _, milestone := range resp.Milestones {
			oldMilestoneID := milestone.ID
			newMilestoneID := core.EntityID(fmt.Sprintf("milestone-%03d", milestoneCounter))
			milestoneIDMap[oldMilestoneID] = newMilestoneID
			milestone.ID = newMilestoneID
			milestoneCounter++
		}
	}

	for _, phase := range resp.Phases {
		var newMilestones []core.EntityID
		for _, oldMilestoneID := range phase.Milestones {
			if newMilestoneID, ok := milestoneIDMap[oldMilestoneID]; ok {
				newMilestones = append(newMilestones, newMilestoneID)
			}
		}
		phase.Milestones = newMilestones
	}

	return nil
}

func extractIDNumber(id core.EntityID) int {
	// Extract number from ID like "phase-001" -> 1
	parts := strings.Split(string(id), "-")
	if len(parts) >= 2 {
		var result int
		if _, err := fmt.Sscanf(parts[len(parts)-1], "%d", &result); err == nil {
			return result
		}
	}
	return 1
}

func displayPathSummary(resp *ai.PathGenerationResponse) {
	fmt.Println()
	PrintSuccess("✨ Learning path generated successfully!")
	fmt.Println()

	fmt.Printf("📚 Path: %s (ID: %s)\n", resp.Path.Title, resp.Path.ID)
	fmt.Printf("   %s\n", resp.Path.Body)
	fmt.Println()

	fmt.Printf("📅 Phases: %d\n", len(resp.Phases))
	for i, phase := range resp.Phases {
		fmt.Printf("   %d. %s (%s)\n", i+1, phase.Title, phase.EstimatedDuration)
		fmt.Printf("      %s\n", phase.Body)
		if len(phase.RequiredSkills) > 0 {
			fmt.Printf("      Required: %d skills\n", len(phase.RequiredSkills))
		}
		if len(phase.Milestones) > 0 {
			fmt.Printf("      Milestones: %d\n", len(phase.Milestones))
		}
	}
	fmt.Println()

	fmt.Printf("📖 Resources: %d\n", len(resp.Resources))
	for i, resource := range resp.Resources {
		if i < 5 { // Show first 5
			fmt.Printf("   • %s (%s) - %.1f hours\n", resource.Title, resource.Type, resource.EstimatedHours)
		}
	}
	if len(resp.Resources) > 5 {
		fmt.Printf("   ... and %d more\n", len(resp.Resources)-5)
	}
	fmt.Println()

	fmt.Printf("🎯 Milestones: %d\n", len(resp.Milestones))
	for i, milestone := range resp.Milestones {
		if i < 3 { // Show first 3
			fmt.Printf("   • %s (%s)\n", milestone.Title, milestone.Type)
		}
	}
	if len(resp.Milestones) > 3 {
		fmt.Printf("   ... and %d more\n", len(resp.Milestones)-3)
	}
	fmt.Println()

	if resp.Reasoning != "" {
		fmt.Println("💡 AI Reasoning:")
		fmt.Printf("   %s\n", resp.Reasoning)
		fmt.Println()
	}

	fmt.Printf("View full path with: growth path view %s\n", resp.Path.ID)
}
