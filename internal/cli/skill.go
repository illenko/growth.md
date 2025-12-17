package cli

import (
	"fmt"
	"strings"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	skillCategory    string
	skillLevel       string
	skillTags        string
	skillStatus      string
	skillFilterLevel string
	skillTitle       string
)

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage skills",
	Long:  `Create, list, view, edit, and delete technical skills.`,
}

var skillCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new skill",
	Long: `Create a new skill with the specified title.

You can provide the title as an argument or be prompted for it.
Optionally specify category, level, and tags using flags.

Examples:
  growth skill create "Python Programming" --category backend --level intermediate
  growth skill create "Docker" --tags containers,devops
  growth skill create`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSkillCreate,
}

var skillListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all skills",
	Long: `List all skills in the repository.

Optionally filter by category, level, or status using flags.

Examples:
  growth skill list
  growth skill list --category backend
  growth skill list --level intermediate
  growth skill list --status learning`,
	Aliases: []string{"ls"},
	RunE:    runSkillList,
}

var skillViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View skill details",
	Long: `View detailed information about a specific skill.

The output format can be controlled with the --format flag (table, json, yaml).

Examples:
  growth skill view skill-001
  growth skill view skill-042 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runSkillView,
}

var skillEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing skill",
	Long: `Edit an existing skill by ID.

You can update any field using flags. If no flags are provided, you'll be prompted
to update each field interactively (press Enter to keep current value).

Examples:
  growth skill edit skill-001 --level advanced
  growth skill edit skill-042 --category frontend --status learning
  growth skill edit skill-001`,
	Args: cobra.ExactArgs(1),
	RunE: runSkillEdit,
}

var skillDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a skill",
	Long: `Delete a skill by ID.

This will permanently remove the skill file. You'll be prompted for confirmation
before deletion unless --force is used.

Examples:
  growth skill delete skill-001
  growth skill delete skill-042`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runSkillDelete,
}

func init() {
	rootCmd.AddCommand(skillCmd)
	skillCmd.AddCommand(skillCreateCmd)
	skillCmd.AddCommand(skillListCmd)
	skillCmd.AddCommand(skillViewCmd)
	skillCmd.AddCommand(skillEditCmd)
	skillCmd.AddCommand(skillDeleteCmd)

	skillCreateCmd.Flags().StringVarP(&skillCategory, "category", "c", "", "skill category")
	skillCreateCmd.Flags().StringVarP(&skillLevel, "level", "l", "", "proficiency level (beginner, intermediate, advanced, expert)")
	skillCreateCmd.Flags().StringVarP(&skillTags, "tags", "t", "", "comma-separated tags")

	skillListCmd.Flags().StringVarP(&skillCategory, "category", "c", "", "filter by category")
	skillListCmd.Flags().StringVarP(&skillFilterLevel, "level", "l", "", "filter by level")
	skillListCmd.Flags().StringVarP(&skillStatus, "status", "s", "", "filter by status")

	skillEditCmd.Flags().StringVar(&skillTitle, "title", "", "skill title")
	skillEditCmd.Flags().StringVarP(&skillCategory, "category", "c", "", "skill category")
	skillEditCmd.Flags().StringVarP(&skillLevel, "level", "l", "", "proficiency level")
	skillEditCmd.Flags().StringVarP(&skillStatus, "status", "s", "", "skill status")
	skillEditCmd.Flags().StringVarP(&skillTags, "tags", "t", "", "comma-separated tags")
}

func runSkillCreate(cmd *cobra.Command, args []string) error {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = PromptStringRequired("Skill title")
	}

	if skillCategory == "" {
		skillCategory = PromptStringRequired("Category (e.g., backend, frontend, devops, data)")
	}

	if skillLevel == "" {
		skillLevel = PromptSelectWithDefault(
			"Proficiency level",
			[]string{"beginner", "intermediate", "advanced", "expert"},
			"beginner",
		)
	}

	level := core.ProficiencyLevel(skillLevel)
	if !level.IsValid() {
		return fmt.Errorf("invalid proficiency level: %s (must be beginner, intermediate, advanced, or expert)", skillLevel)
	}

	id, err := GenerateNextID("skill")
	if err != nil {
		return fmt.Errorf("failed to generate skill ID: %w", err)
	}

	skill, err := core.NewSkill(id, title, skillCategory, level)
	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}

	if skillTags != "" {
		tags := strings.Split(skillTags, ",")
		for _, tag := range tags {
			skill.AddTag(strings.TrimSpace(tag))
		}
	}

	description := PromptMultiline("Description (optional, press Ctrl+D or enter '.' to finish)")
	if description != "" {
		skill.Body = description
	}

	if err := skillRepo.Create(skill); err != nil {
		return fmt.Errorf("failed to save skill: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Created skill %s: %s", skill.ID, skill.Title))

	if verbose {
		fmt.Printf("\nSkill details:\n")
		fmt.Printf("  ID: %s\n", skill.ID)
		fmt.Printf("  Title: %s\n", skill.Title)
		fmt.Printf("  Category: %s\n", skill.Category)
		fmt.Printf("  Level: %s\n", skill.Level)
		if len(skill.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(skill.Tags, ", "))
		}
	}

	return nil
}

func runSkillList(cmd *cobra.Command, args []string) error {
	var skills []*core.Skill
	var err error

	if skillCategory != "" && skillFilterLevel != "" {
		level := core.ProficiencyLevel(skillFilterLevel)
		if !level.IsValid() {
			return fmt.Errorf("invalid proficiency level: %s", skillFilterLevel)
		}
		skills, err = skillRepo.FindByCategoryAndLevel(skillCategory, level)
	} else if skillCategory != "" {
		skills, err = skillRepo.FindByCategory(skillCategory)
	} else if skillFilterLevel != "" {
		level := core.ProficiencyLevel(skillFilterLevel)
		if !level.IsValid() {
			return fmt.Errorf("invalid proficiency level: %s", skillFilterLevel)
		}
		skills, err = skillRepo.FindByLevel(level)
	} else if skillStatus != "" {
		status := core.SkillStatus(skillStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid skill status: %s", skillStatus)
		}
		skills, err = skillRepo.FindByStatus(status)
	} else {
		skills, err = skillRepo.GetAll()
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve skills: %w", err)
	}

	if len(skills) == 0 {
		PrintInfo("No skills found")
		return nil
	}

	return PrintOutputWithConfig(skills)
}

func runSkillView(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	skill, err := skillRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve skill: %w", err)
	}

	if config.Display.OutputFormat == "table" {
		fmt.Printf("ID:       %s\n", skill.ID)
		fmt.Printf("Title:    %s\n", skill.Title)
		fmt.Printf("Category: %s\n", skill.Category)
		fmt.Printf("Level:    %s\n", skill.Level)
		fmt.Printf("Status:   %s\n", skill.Status)
		if len(skill.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(skill.Tags, ", "))
		}
		if len(skill.Resources) > 0 {
			fmt.Printf("Resources: %v\n", skill.Resources)
		}
		fmt.Printf("Created:  %s\n", skill.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:  %s\n", skill.Updated.Format("2006-01-02 15:04:05"))

		if skill.Body != "" {
			fmt.Printf("\nDescription:\n%s\n", skill.Body)
		}

		return nil
	}

	return PrintOutputWithConfig(skill)
}

func runSkillEdit(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	skill, err := skillRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve skill: %w", err)
	}

	updated := false

	if cmd.Flags().Changed("title") {
		skill.Title = skillTitle
		updated = true
	}

	if cmd.Flags().Changed("category") {
		skill.Category = skillCategory
		updated = true
	}

	if cmd.Flags().Changed("level") {
		level := core.ProficiencyLevel(skillLevel)
		if !level.IsValid() {
			return fmt.Errorf("invalid proficiency level: %s", skillLevel)
		}
		if err := skill.UpdateLevel(level); err != nil {
			return fmt.Errorf("failed to update level: %w", err)
		}
		updated = true
	}

	if cmd.Flags().Changed("status") {
		status := core.SkillStatus(skillStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid skill status: %s", skillStatus)
		}
		if err := skill.UpdateStatus(status); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
		updated = true
	}

	if cmd.Flags().Changed("tags") {
		skill.Tags = []string{}
		if skillTags != "" {
			tags := strings.Split(skillTags, ",")
			for _, tag := range tags {
				skill.AddTag(strings.TrimSpace(tag))
			}
		}
		updated = true
	}

	if !updated {
		PrintInfo("No changes specified. Use flags to update fields or run interactively.")

		if PromptConfirm("Update title?") {
			skill.Title = PromptString("New title", skill.Title)
			updated = true
		}

		if PromptConfirm("Update category?") {
			skill.Category = PromptString("New category", skill.Category)
			updated = true
		}

		if PromptConfirm("Update level?") {
			newLevel := PromptSelectWithDefault(
				"Proficiency level",
				[]string{"beginner", "intermediate", "advanced", "expert"},
				string(skill.Level),
			)
			level := core.ProficiencyLevel(newLevel)
			if err := skill.UpdateLevel(level); err != nil {
				return fmt.Errorf("failed to update level: %w", err)
			}
			updated = true
		}

		if PromptConfirm("Update status?") {
			newStatus := PromptSelectWithDefault(
				"Skill status",
				[]string{"not-started", "learning", "mastered"},
				string(skill.Status),
			)
			status := core.SkillStatus(newStatus)
			if err := skill.UpdateStatus(status); err != nil {
				return fmt.Errorf("failed to update status: %w", err)
			}
			updated = true
		}

		if PromptConfirm("Update description?") {
			description := PromptMultiline("Description (press Ctrl+D or enter '.' to finish)")
			skill.Body = description
			updated = true
		}
	}

	if !updated {
		PrintInfo("No changes made")
		return nil
	}

	if err := skillRepo.Update(skill); err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Updated skill %s: %s", skill.ID, skill.Title))
	return nil
}

func runSkillDelete(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	skill, err := skillRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve skill: %w", err)
	}

	fmt.Printf("You are about to delete:\n")
	fmt.Printf("  ID: %s\n", skill.ID)
	fmt.Printf("  Title: %s\n", skill.Title)
	fmt.Printf("  Category: %s\n", skill.Category)
	fmt.Println()

	if !PromptConfirm("Are you sure you want to delete this skill?") {
		PrintInfo("Deletion cancelled")
		return nil
	}

	if err := skillRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Deleted skill %s", id))
	return nil
}
