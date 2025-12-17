package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	goalPriority   string
	goalStatus     string
	goalTags       string
	goalTargetDate string
	goalTitle      string
)

var goalCmd = &cobra.Command{
	Use:   "goal",
	Short: "Manage goals",
	Long:  `Create, list, view, edit, and delete career goals.`,
}

var goalCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new goal",
	Long: `Create a new career goal with the specified title.

You can provide the title as an argument or be prompted for it.
Optionally specify priority, target date, and tags using flags.

Examples:
  growth goal create "Senior Engineer by 2025" --priority high --target 2025-12-31
  growth goal create "Learn Cloud Architecture" --tags cloud,aws,architecture
  growth goal create`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGoalCreate,
}

var goalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all goals",
	Long: `List all goals in the repository.

Optionally filter by status or priority using flags.

Examples:
  growth goal list
  growth goal list --status active
  growth goal list --priority high`,
	Aliases: []string{"ls"},
	RunE:    runGoalList,
}

var goalViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View goal details",
	Long: `View detailed information about a specific goal.

The output format can be controlled with the --format flag (table, json, yaml).

Examples:
  growth goal view goal-001
  growth goal view goal-042 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runGoalView,
}

var goalEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing goal",
	Long: `Edit an existing goal by ID.

You can update any field using flags. If no flags are provided, you'll be prompted
to update each field interactively (press Enter to keep current value).

Examples:
  growth goal edit goal-001 --priority high
  growth goal edit goal-042 --status completed --target 2025-06-30
  growth goal edit goal-001`,
	Args: cobra.ExactArgs(1),
	RunE: runGoalEdit,
}

var goalDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a goal",
	Long: `Delete a goal by ID.

This will permanently remove the goal file. You'll be prompted for confirmation
before deletion.

Examples:
  growth goal delete goal-001
  growth goal delete goal-042`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runGoalDelete,
}

var goalAddPathCmd = &cobra.Command{
	Use:   "add-path <goal-id> <path-id>",
	Short: "Add a learning path to a goal",
	Long: `Associate a learning path with a goal.

Examples:
  growth goal add-path goal-001 path-001`,
	Args: cobra.ExactArgs(2),
	RunE: runGoalAddPath,
}

var goalRemovePathCmd = &cobra.Command{
	Use:   "remove-path <goal-id> <path-id>",
	Short: "Remove a learning path from a goal",
	Long: `Disassociate a learning path from a goal.

Examples:
  growth goal remove-path goal-001 path-001`,
	Args: cobra.ExactArgs(2),
	RunE: runGoalRemovePath,
}

func init() {
	rootCmd.AddCommand(goalCmd)
	goalCmd.AddCommand(goalCreateCmd)
	goalCmd.AddCommand(goalListCmd)
	goalCmd.AddCommand(goalViewCmd)
	goalCmd.AddCommand(goalEditCmd)
	goalCmd.AddCommand(goalDeleteCmd)
	goalCmd.AddCommand(goalAddPathCmd)
	goalCmd.AddCommand(goalRemovePathCmd)

	goalCreateCmd.Flags().StringVarP(&goalPriority, "priority", "p", "", "goal priority (high, medium, low)")
	goalCreateCmd.Flags().StringVarP(&goalTargetDate, "target", "d", "", "target date (YYYY-MM-DD)")
	goalCreateCmd.Flags().StringVarP(&goalTags, "tags", "t", "", "comma-separated tags")

	goalListCmd.Flags().StringVarP(&goalStatus, "status", "s", "", "filter by status (active, completed, archived)")
	goalListCmd.Flags().StringVarP(&goalPriority, "priority", "p", "", "filter by priority (high, medium, low)")

	goalEditCmd.Flags().StringVar(&goalTitle, "title", "", "goal title")
	goalEditCmd.Flags().StringVarP(&goalPriority, "priority", "p", "", "goal priority")
	goalEditCmd.Flags().StringVarP(&goalStatus, "status", "s", "", "goal status")
	goalEditCmd.Flags().StringVarP(&goalTargetDate, "target", "d", "", "target date (YYYY-MM-DD)")
	goalEditCmd.Flags().StringVarP(&goalTags, "tags", "t", "", "comma-separated tags")
}

func runGoalCreate(cmd *cobra.Command, args []string) error {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = PromptStringRequired("Goal title")
	}

	if goalPriority == "" {
		goalPriority = PromptSelectWithDefault(
			"Priority",
			[]string{"high", "medium", "low"},
			"medium",
		)
	}

	priority := core.Priority(goalPriority)
	if !priority.IsValid() {
		return fmt.Errorf("invalid priority: %s (must be high, medium, or low)", goalPriority)
	}

	id, err := GenerateNextID("goal")
	if err != nil {
		return fmt.Errorf("failed to generate goal ID: %w", err)
	}

	goal, err := core.NewGoal(id, title, priority)
	if err != nil {
		return fmt.Errorf("failed to create goal: %w", err)
	}

	if goalTargetDate != "" {
		targetDate, err := time.Parse("2006-01-02", goalTargetDate)
		if err != nil {
			return fmt.Errorf("invalid target date format (use YYYY-MM-DD): %w", err)
		}
		goal.SetTargetDate(targetDate)
	}

	if goalTags != "" {
		tags := strings.Split(goalTags, ",")
		for _, tag := range tags {
			goal.AddTag(strings.TrimSpace(tag))
		}
	}

	description := PromptMultiline("Description (optional, press Ctrl+D or enter '.' to finish)")
	if description != "" {
		goal.Body = description
	}

	if err := goalRepo.Create(goal); err != nil {
		return fmt.Errorf("failed to save goal: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Created goal %s: %s", goal.ID, goal.Title))

	if verbose {
		fmt.Printf("\nGoal details:\n")
		fmt.Printf("  ID: %s\n", goal.ID)
		fmt.Printf("  Title: %s\n", goal.Title)
		fmt.Printf("  Priority: %s\n", goal.Priority)
		fmt.Printf("  Status: %s\n", goal.Status)
		if goal.TargetDate != nil {
			fmt.Printf("  Target: %s\n", goal.TargetDate.Format("2006-01-02"))
		}
		if len(goal.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(goal.Tags, ", "))
		}
	}

	return nil
}

func runGoalList(cmd *cobra.Command, args []string) error {
	var goals []*core.Goal
	var err error

	if goalStatus != "" {
		status := core.Status(goalStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status: %s", goalStatus)
		}
		goals, err = goalRepo.FindByStatus(status)
	} else if goalPriority != "" {
		priority := core.Priority(goalPriority)
		if !priority.IsValid() {
			return fmt.Errorf("invalid priority: %s", goalPriority)
		}
		goals, err = goalRepo.FindByPriority(priority)
	} else {
		goals, err = goalRepo.GetAll()
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve goals: %w", err)
	}

	if goalPriority != "" && goalStatus != "" {
		priority := core.Priority(goalPriority)
		status := core.Status(goalStatus)
		var filtered []*core.Goal
		for _, g := range goals {
			if g.Priority == priority && g.Status == status {
				filtered = append(filtered, g)
			}
		}
		goals = filtered
	}

	if len(goals) == 0 {
		PrintInfo("No goals found")
		return nil
	}

	return PrintOutputWithConfig(goals)
}

func runGoalView(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	goal, err := goalRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve goal: %w", err)
	}

	if config.Display.OutputFormat == "table" {
		fmt.Printf("ID:       %s\n", goal.ID)
		fmt.Printf("Title:    %s\n", goal.Title)
		fmt.Printf("Status:   %s\n", goal.Status)
		fmt.Printf("Priority: %s\n", goal.Priority)
		if goal.TargetDate != nil {
			fmt.Printf("Target:   %s\n", goal.TargetDate.Format("2006-01-02"))
		}
		if len(goal.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(goal.Tags, ", "))
		}
		if len(goal.LearningPaths) > 0 {
			fmt.Printf("Paths:    %v\n", goal.LearningPaths)
		}
		if len(goal.Milestones) > 0 {
			fmt.Printf("Milestones: %v\n", goal.Milestones)
		}
		fmt.Printf("Created:  %s\n", goal.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:  %s\n", goal.Updated.Format("2006-01-02 15:04:05"))

		if goal.Body != "" {
			fmt.Printf("\nDescription:\n%s\n", goal.Body)
		}

		return nil
	}

	return PrintOutputWithConfig(goal)
}

func runGoalEdit(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	goal, err := goalRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve goal: %w", err)
	}

	updated := false

	if cmd.Flags().Changed("title") {
		goal.Title = goalTitle
		updated = true
	}

	if cmd.Flags().Changed("priority") {
		priority := core.Priority(goalPriority)
		if !priority.IsValid() {
			return fmt.Errorf("invalid priority: %s", goalPriority)
		}
		if err := goal.UpdatePriority(priority); err != nil {
			return fmt.Errorf("failed to update priority: %w", err)
		}
		updated = true
	}

	if cmd.Flags().Changed("status") {
		status := core.Status(goalStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status: %s", goalStatus)
		}
		if err := goal.UpdateStatus(status); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
		updated = true
	}

	if cmd.Flags().Changed("target") {
		if goalTargetDate == "" {
			goal.ClearTargetDate()
		} else {
			targetDate, err := time.Parse("2006-01-02", goalTargetDate)
			if err != nil {
				return fmt.Errorf("invalid target date format (use YYYY-MM-DD): %w", err)
			}
			goal.SetTargetDate(targetDate)
		}
		updated = true
	}

	if cmd.Flags().Changed("tags") {
		goal.Tags = []string{}
		if goalTags != "" {
			tags := strings.Split(goalTags, ",")
			for _, tag := range tags {
				goal.AddTag(strings.TrimSpace(tag))
			}
		}
		updated = true
	}

	if !updated {
		PrintInfo("No changes specified. Use flags to update fields or run interactively.")

		if PromptConfirm("Update title?") {
			goal.Title = PromptString("New title", goal.Title)
			updated = true
		}

		if PromptConfirm("Update priority?") {
			newPriority := PromptSelectWithDefault(
				"Priority",
				[]string{"high", "medium", "low"},
				string(goal.Priority),
			)
			priority := core.Priority(newPriority)
			if err := goal.UpdatePriority(priority); err != nil {
				return fmt.Errorf("failed to update priority: %w", err)
			}
			updated = true
		}

		if PromptConfirm("Update status?") {
			newStatus := PromptSelectWithDefault(
				"Status",
				[]string{"active", "completed", "archived"},
				string(goal.Status),
			)
			status := core.Status(newStatus)
			if err := goal.UpdateStatus(status); err != nil {
				return fmt.Errorf("failed to update status: %w", err)
			}
			updated = true
		}

		if PromptConfirm("Update target date?") {
			defaultDate := ""
			if goal.TargetDate != nil {
				defaultDate = goal.TargetDate.Format("2006-01-02")
			}
			dateStr := PromptString("Target date (YYYY-MM-DD, empty to clear)", defaultDate)
			if dateStr == "" {
				goal.ClearTargetDate()
			} else {
				targetDate, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					return fmt.Errorf("invalid date format: %w", err)
				}
				goal.SetTargetDate(targetDate)
			}
			updated = true
		}

		if PromptConfirm("Update description?") {
			description := PromptMultiline("Description (press Ctrl+D or enter '.' to finish)")
			goal.Body = description
			updated = true
		}
	}

	if !updated {
		PrintInfo("No changes made")
		return nil
	}

	if err := goalRepo.Update(goal); err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Updated goal %s: %s", goal.ID, goal.Title))
	return nil
}

func runGoalDelete(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	goal, err := goalRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve goal: %w", err)
	}

	fmt.Printf("You are about to delete:\n")
	fmt.Printf("  ID: %s\n", goal.ID)
	fmt.Printf("  Title: %s\n", goal.Title)
	fmt.Printf("  Priority: %s\n", goal.Priority)
	fmt.Println()

	if !PromptConfirm("Are you sure you want to delete this goal?") {
		PrintInfo("Deletion cancelled")
		return nil
	}

	if err := goalRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Deleted goal %s", id))
	return nil
}

func runGoalAddPath(cmd *cobra.Command, args []string) error {
	goalID := core.EntityID(args[0])
	pathID := core.EntityID(args[1])

	goal, err := goalRepo.GetByIDWithBody(goalID)
	if err != nil {
		return fmt.Errorf("failed to retrieve goal: %w", err)
	}

	exists, err := pathRepo.Exists(pathID)
	if err != nil {
		return fmt.Errorf("failed to check path existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("path %s does not exist", pathID)
	}

	goal.AddLearningPath(pathID)

	if err := goalRepo.Update(goal); err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Added path %s to goal %s", pathID, goalID))
	return nil
}

func runGoalRemovePath(cmd *cobra.Command, args []string) error {
	goalID := core.EntityID(args[0])
	pathID := core.EntityID(args[1])

	goal, err := goalRepo.GetByIDWithBody(goalID)
	if err != nil {
		return fmt.Errorf("failed to retrieve goal: %w", err)
	}

	goal.RemoveLearningPath(pathID)

	if err := goalRepo.Update(goal); err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Removed path %s from goal %s", pathID, goalID))
	return nil
}
