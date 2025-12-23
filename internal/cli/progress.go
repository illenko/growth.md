package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	progressDate   string
	progressHours  string
	progressMood   string
	progressSkills string
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Manage progress logs",
	Long:  `Log daily progress and view progress history.`,
}

var progressLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Log progress for a date",
	Long: `Create a progress log for a specific date.

Examples:
  growth progress log
  growth progress log --hours 15 --mood motivated
  growth progress log --date 2025-12-16`,
	RunE: runProgressLog,
}

var progressListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all progress logs",
	Long: `List all progress logs in chronological order.

Examples:
  growth progress list
  growth progress list --format json`,
	Aliases: []string{"ls"},
	RunE:    runProgressList,
}

var progressViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View progress log details",
	Long: `View detailed information about a specific progress log.

Examples:
  growth progress view progress-001
  growth progress view progress-042 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runProgressView,
}

func init() {
	rootCmd.AddCommand(progressCmd)
	progressCmd.AddCommand(progressLogCmd)
	progressCmd.AddCommand(progressListCmd)
	progressCmd.AddCommand(progressViewCmd)

	progressLogCmd.Flags().StringVar(&progressDate, "date", "", "date for progress log (YYYY-MM-DD), defaults to today")
	progressLogCmd.Flags().StringVar(&progressHours, "hours", "", "hours invested")
	progressLogCmd.Flags().StringVar(&progressMood, "mood", "", "mood (e.g., motivated, frustrated, focused)")
	progressLogCmd.Flags().StringVar(&progressSkills, "skills", "", "comma-separated skill IDs")
}

func runProgressLog(cmd *cobra.Command, args []string) error {
	var date time.Time
	var err error

	if progressDate != "" {
		date, err = time.Parse("2006-01-02", progressDate)
		if err != nil {
			return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
		}
	} else {
		date = time.Now()
	}

	id, err := GenerateNextID("progress")
	if err != nil {
		return fmt.Errorf("failed to generate progress ID: %w", err)
	}

	log, err := core.NewProgressLog(id, date)
	if err != nil {
		return fmt.Errorf("failed to create progress log: %w", err)
	}

	if progressHours != "" {
		hours, err := strconv.ParseFloat(progressHours, 64)
		if err != nil {
			return fmt.Errorf("invalid hours value: %w", err)
		}
		if err := log.SetHoursInvested(hours); err != nil {
			return fmt.Errorf("failed to set hours: %w", err)
		}
	} else {
		hours := PromptInt("Hours invested", 0)
		if hours > 0 {
			if err := log.SetHoursInvested(float64(hours)); err != nil {
				return fmt.Errorf("failed to set hours: %w", err)
			}
		}
	}

	if progressMood != "" {
		log.SetMood(progressMood)
	} else {
		mood := PromptString("Mood (optional)", "")
		if mood != "" {
			log.SetMood(mood)
		}
	}

	if progressSkills != "" {
		skills := strings.Split(progressSkills, ",")
		for _, skillStr := range skills {
			skillID := core.EntityID(strings.TrimSpace(skillStr))
			log.AddSkillWorked(skillID)
		}
	}

	summary := PromptMultiline("Daily summary (press Ctrl+D or enter '.' to finish)")
	if summary != "" {
		log.Body = summary
	}

	if err := progressRepo.Create(log); err != nil {
		return fmt.Errorf("failed to save progress log: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Logged progress %s for %s", log.ID, log.Date.Format("2006-01-02")))

	if verbose {
		fmt.Printf("\nProgress log details:\n")
		fmt.Printf("  ID: %s\n", log.ID)
		fmt.Printf("  Date: %s\n", log.Date.Format("2006-01-02"))
		if log.HoursInvested > 0 {
			fmt.Printf("  Hours: %.1f\n", log.HoursInvested)
		}
		if log.Mood != "" {
			fmt.Printf("  Mood: %s\n", log.Mood)
		}
		if len(log.SkillsWorked) > 0 {
			fmt.Printf("  Skills: %v\n", log.SkillsWorked)
		}
	}

	return nil
}

func runProgressList(cmd *cobra.Command, args []string) error {
	logs, err := progressRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to retrieve progress logs: %w\nTry running 'growth progress list' again or check your repository", err)
	}

	if len(logs) == 0 {
		PrintInfo("No progress logs found")
		return nil
	}

	return PrintOutputWithConfig(logs)
}

func runProgressView(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	log, err := progressRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("progress log '%s' not found. Use 'growth progress list' to see available logs", id)
	}

	if config.Display.OutputFormat == "table" {
		fmt.Printf("ID:       %s\n", log.ID)
		fmt.Printf("Date:     %s\n", log.Date.Format("2006-01-02"))
		if log.HoursInvested > 0 {
			fmt.Printf("Hours:    %.1f\n", log.HoursInvested)
		}
		if log.Mood != "" {
			fmt.Printf("Mood:     %s\n", log.Mood)
		}
		if len(log.SkillsWorked) > 0 {
			fmt.Printf("Skills:   %v\n", log.SkillsWorked)
		}
		if len(log.ResourcesUsed) > 0 {
			fmt.Printf("Resources: %v\n", log.ResourcesUsed)
		}
		if len(log.MilestonesAchieved) > 0 {
			fmt.Printf("Milestones: %v\n", log.MilestonesAchieved)
		}
		fmt.Printf("Created:  %s\n", log.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:  %s\n", log.Updated.Format("2006-01-02 15:04:05"))

		if log.Body != "" {
			fmt.Printf("\nSummary:\n%s\n", log.Body)
		}

		return nil
	}

	return PrintOutputWithConfig(log)
}
