package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	milestoneType       string
	milestoneRefType    string
	milestoneRefID      string
	milestoneTargetDate string
	milestoneProof      string
	milestoneStatus     string
	milestoneTitle      string
	milestoneFilterType string
)

var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Manage milestones",
	Long:  `Create, list, view, edit, and delete milestones.`,
}

var milestoneCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new milestone",
	Long: `Create a new milestone with the specified title.

Milestones must be associated with a goal, path, or skill using --ref-type and --ref-id.

Examples:
  growth milestone create "Deploy first app" --type skill-level --ref-type skill --ref-id skill-001
  growth milestone create "Complete course" --type goal-level --ref-type goal --ref-id goal-001 --target 2025-06-30
  growth milestone create`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMilestoneCreate,
}

var milestoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all milestones",
	Long: `List all milestones in the repository.

Optionally filter by type, status, or reference ID using flags.

Examples:
  growth milestone list
  growth milestone list --type goal-level
  growth milestone list --status completed
  growth milestone list --ref-id goal-001`,
	Aliases: []string{"ls"},
	RunE:    runMilestoneList,
}

var milestoneViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View milestone details",
	Long: `View detailed information about a specific milestone.

The output format can be controlled with the --format flag (table, json, yaml).

Examples:
  growth milestone view milestone-001
  growth milestone view milestone-042 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runMilestoneView,
}

var milestoneEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing milestone",
	Long: `Edit an existing milestone by ID.

You can update any field using flags.

Examples:
  growth milestone edit milestone-001 --status completed --proof https://github.com/user/repo
  growth milestone edit milestone-042 --target 2025-12-31
  growth milestone edit milestone-001 --title "New Title"`,
	Args: cobra.ExactArgs(1),
	RunE: runMilestoneEdit,
}

var milestoneDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a milestone",
	Long: `Delete a milestone by ID.

This will permanently remove the milestone file. You'll be prompted for confirmation
before deletion.

Examples:
  growth milestone delete milestone-001
  growth milestone delete milestone-042`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runMilestoneDelete,
}

var milestoneAchieveCmd = &cobra.Command{
	Use:   "achieve <id>",
	Short: "Mark milestone as achieved",
	Long: `Mark a milestone as achieved with optional proof URL.

Examples:
  growth milestone achieve milestone-001
  growth milestone achieve milestone-001 --proof https://github.com/user/repo`,
	Args: cobra.ExactArgs(1),
	RunE: runMilestoneAchieve,
}

func init() {
	rootCmd.AddCommand(milestoneCmd)
	milestoneCmd.AddCommand(milestoneCreateCmd)
	milestoneCmd.AddCommand(milestoneListCmd)
	milestoneCmd.AddCommand(milestoneViewCmd)
	milestoneCmd.AddCommand(milestoneEditCmd)
	milestoneCmd.AddCommand(milestoneDeleteCmd)
	milestoneCmd.AddCommand(milestoneAchieveCmd)

	milestoneCreateCmd.Flags().StringVarP(&milestoneType, "type", "t", "", "milestone type (goal-level, path-level, skill-level)")
	milestoneCreateCmd.Flags().StringVar(&milestoneRefType, "ref-type", "", "reference type (goal, path, skill)")
	milestoneCreateCmd.Flags().StringVar(&milestoneRefID, "ref-id", "", "reference ID (e.g., goal-001)")
	milestoneCreateCmd.Flags().StringVar(&milestoneTargetDate, "target", "", "target date (YYYY-MM-DD)")
	milestoneCreateCmd.MarkFlagRequired("ref-type")
	milestoneCreateCmd.MarkFlagRequired("ref-id")

	milestoneListCmd.Flags().StringVarP(&milestoneFilterType, "type", "t", "", "filter by type")
	milestoneListCmd.Flags().StringVarP(&milestoneStatus, "status", "s", "", "filter by status (active, completed)")
	milestoneListCmd.Flags().StringVar(&milestoneRefID, "ref-id", "", "filter by reference ID")

	milestoneEditCmd.Flags().StringVar(&milestoneTitle, "title", "", "milestone title")
	milestoneEditCmd.Flags().StringVarP(&milestoneStatus, "status", "s", "", "milestone status")
	milestoneEditCmd.Flags().StringVar(&milestoneTargetDate, "target", "", "target date (YYYY-MM-DD)")
	milestoneEditCmd.Flags().StringVar(&milestoneProof, "proof", "", "proof URL")

	milestoneAchieveCmd.Flags().StringVar(&milestoneProof, "proof", "", "proof URL")
}

func runMilestoneCreate(cmd *cobra.Command, args []string) error {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = PromptStringRequired("Milestone title")
	}

	if milestoneRefType == "" {
		milestoneRefType = PromptSelectWithDefault(
			"Reference type",
			[]string{"goal", "path", "skill"},
			"goal",
		)
	}

	refType := core.ReferenceType(milestoneRefType)
	if !refType.IsValid() {
		return fmt.Errorf("invalid reference type: %s (must be goal, path, or skill)", milestoneRefType)
	}

	if milestoneRefID == "" {
		milestoneRefID = PromptStringRequired("Reference ID (e.g., goal-001)")
	}

	refID := core.EntityID(milestoneRefID)

	// Validate reference exists
	var exists bool
	var err error
	switch refType {
	case core.ReferenceGoal:
		exists, err = goalRepo.Exists(refID)
	case core.ReferencePath:
		exists, err = pathRepo.Exists(refID)
	case core.ReferenceSkill:
		exists, err = skillRepo.Exists(refID)
	}
	if err != nil {
		return fmt.Errorf("failed to check reference existence: %w", err)
	}
	if !exists {
		var listCmd string
		switch refType {
		case core.ReferenceGoal:
			listCmd = "growth goal list"
		case core.ReferencePath:
			listCmd = "growth path list"
		case core.ReferenceSkill:
			listCmd = "growth skill list"
		}
		return fmt.Errorf("%s '%s' not found. Use '%s' to see available %ss", refType, refID, listCmd, refType)
	}

	if milestoneType == "" {
		milestoneType = PromptSelectWithDefault(
			"Milestone type",
			[]string{"goal-level", "path-level", "skill-level"},
			string(refType)+"-level",
		)
	}

	mType := core.MilestoneType(milestoneType)
	if !mType.IsValid() {
		return fmt.Errorf("invalid milestone type '%s'. Valid options: goal-level, path-level, skill-level", milestoneType)
	}

	id, err := GenerateNextID("milestone")
	if err != nil {
		return fmt.Errorf("failed to generate milestone ID: %w", err)
	}

	milestone, err := core.NewMilestone(id, title, mType, refType, refID)
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}

	if milestoneTargetDate != "" {
		targetDate, err := time.Parse("2006-01-02", milestoneTargetDate)
		if err != nil {
			return fmt.Errorf("invalid target date format (use YYYY-MM-DD): %w", err)
		}
		milestone.SetTargetDate(targetDate)
	}

	description := PromptMultiline("Description (optional, press Ctrl+D or enter '.' to finish)")
	if description != "" {
		milestone.Body = description
	}

	if err := milestoneRepo.Create(milestone); err != nil {
		return fmt.Errorf("failed to save milestone: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Created milestone %s: %s", milestone.ID, milestone.Title))

	if verbose {
		fmt.Printf("\nMilestone details:\n")
		fmt.Printf("  ID: %s\n", milestone.ID)
		fmt.Printf("  Title: %s\n", milestone.Title)
		fmt.Printf("  Type: %s\n", milestone.Type)
		fmt.Printf("  Reference: %s (%s)\n", milestone.ReferenceID, milestone.ReferenceType)
		if milestone.TargetDate != nil {
			fmt.Printf("  Target: %s\n", milestone.TargetDate.Format("2006-01-02"))
		}
	}

	return nil
}

func runMilestoneList(cmd *cobra.Command, args []string) error {
	var milestones []*core.Milestone
	var err error

	if milestoneRefID != "" {
		refID := core.EntityID(milestoneRefID)
		// Need to determine refType from the ID prefix
		var refType core.ReferenceType
		if strings.HasPrefix(string(refID), "goal-") {
			refType = core.ReferenceGoal
		} else if strings.HasPrefix(string(refID), "path-") {
			refType = core.ReferencePath
		} else if strings.HasPrefix(string(refID), "skill-") {
			refType = core.ReferenceSkill
		} else {
			return fmt.Errorf("cannot determine reference type from ID: %s", refID)
		}
		milestones, err = milestoneRepo.FindByReferenceID(refType, refID)
	} else if milestoneFilterType != "" {
		mType := core.MilestoneType(milestoneFilterType)
		if !mType.IsValid() {
			return fmt.Errorf("invalid milestone type '%s'. Valid options: goal-level, path-level, skill-level", milestoneFilterType)
		}
		milestones, err = milestoneRepo.FindByType(mType)
	} else if milestoneStatus != "" {
		status := core.Status(milestoneStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status '%s'. Valid options: active, completed, archived", milestoneStatus)
		}
		milestones, err = milestoneRepo.FindByStatus(status)
	} else {
		milestones, err = milestoneRepo.GetAll()
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve milestones: %w\nTry running 'growth milestone list' without filters to see all milestones", err)
	}

	if len(milestones) == 0 {
		PrintInfo("No milestones found")
		return nil
	}

	return PrintOutputWithConfig(milestones)
}

func runMilestoneView(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	milestone, err := milestoneRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("milestone '%s' not found. Use 'growth milestone list' to see available milestones", id)
	}

	if config.Display.OutputFormat == "table" {
		fmt.Printf("ID:       %s\n", milestone.ID)
		fmt.Printf("Title:    %s\n", milestone.Title)
		fmt.Printf("Type:     %s\n", milestone.Type)
		fmt.Printf("Reference: %s (%s)\n", milestone.ReferenceID, milestone.ReferenceType)
		fmt.Printf("Status:   %s\n", milestone.Status)
		if milestone.TargetDate != nil {
			fmt.Printf("Target:   %s\n", milestone.TargetDate.Format("2006-01-02"))
		}
		if milestone.AchievedDate != nil {
			fmt.Printf("Achieved: %s\n", milestone.AchievedDate.Format("2006-01-02"))
		}
		if milestone.Proof != "" {
			fmt.Printf("Proof:    %s\n", milestone.Proof)
		}
		fmt.Printf("Created:  %s\n", milestone.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:  %s\n", milestone.Updated.Format("2006-01-02 15:04:05"))

		if milestone.Body != "" {
			fmt.Printf("\nDescription:\n%s\n", milestone.Body)
		}

		return nil
	}

	return PrintOutputWithConfig(milestone)
}

func runMilestoneEdit(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	milestone, err := milestoneRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("milestone '%s' not found. Use 'growth milestone list' to see available milestones", id)
	}

	updated := false

	if cmd.Flags().Changed("title") {
		milestone.Title = milestoneTitle
		updated = true
	}

	if cmd.Flags().Changed("status") {
		status := core.Status(milestoneStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status '%s'. Valid options: active, completed, archived", milestoneStatus)
		}
		milestone.Status = status
		if status == core.StatusCompleted && milestone.AchievedDate == nil {
			now := time.Now()
			milestone.AchievedDate = &now
		}
		updated = true
	}

	if cmd.Flags().Changed("target") {
		if milestoneTargetDate == "" {
			milestone.ClearTargetDate()
		} else {
			targetDate, err := time.Parse("2006-01-02", milestoneTargetDate)
			if err != nil {
				return fmt.Errorf("invalid target date format (use YYYY-MM-DD): %w", err)
			}
			milestone.SetTargetDate(targetDate)
		}
		updated = true
	}

	if cmd.Flags().Changed("proof") {
		milestone.SetProof(milestoneProof)
		updated = true
	}

	if !updated {
		PrintInfo("No changes specified. Use flags to update fields.")
		return nil
	}

	if err := milestoneRepo.Update(milestone); err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Updated milestone %s: %s", milestone.ID, milestone.Title))
	return nil
}

func runMilestoneDelete(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	milestone, err := milestoneRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("milestone '%s' not found. Use 'growth milestone list' to see available milestones", id)
	}

	fmt.Printf("You are about to delete:\n")
	fmt.Printf("  ID: %s\n", milestone.ID)
	fmt.Printf("  Title: %s\n", milestone.Title)
	fmt.Printf("  Type: %s\n", milestone.Type)
	fmt.Println()

	if !PromptConfirm("Are you sure you want to delete this milestone?") {
		PrintInfo("Deletion cancelled")
		return nil
	}

	if err := milestoneRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Deleted milestone %s", id))
	return nil
}

func runMilestoneAchieve(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	milestone, err := milestoneRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("milestone '%s' not found. Use 'growth milestone list' to see available milestones", id)
	}

	proof := milestoneProof
	if proof == "" && PromptConfirm("Add proof URL?") {
		proof = PromptString("Proof URL", "")
	}

	milestone.Achieve(proof)

	if err := milestoneRepo.Update(milestone); err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Achieved milestone %s: %s", milestone.ID, milestone.Title))
	if milestone.Proof != "" {
		fmt.Printf("Proof: %s\n", milestone.Proof)
	}

	return nil
}
