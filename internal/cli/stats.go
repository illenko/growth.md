package cli

import (
	"fmt"
	"sort"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display detailed statistics",
	Long: `Display detailed statistics about your growth journey.

Shows trends, top categories, progress over time, and more.

Examples:
  growth stats`,
	RunE: runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	fmt.Println("Growth Statistics")
	fmt.Println("=================")
	fmt.Println()

	// Skill categories
	skills, err := skillRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get skills: %w", err)
	}

	categoryCount := make(map[string]int)
	for _, skill := range skills {
		categoryCount[skill.Category]++
	}

	if len(categoryCount) > 0 {
		fmt.Println("Top Skill Categories:")
		type categoryStats struct {
			name  string
			count int
		}
		var categories []categoryStats
		for name, count := range categoryCount {
			categories = append(categories, categoryStats{name, count})
		}
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].count > categories[j].count
		})
		for i, cat := range categories {
			if i >= 5 {
				break
			}
			fmt.Printf("  %d. %s (%d skills)\n", i+1, cat.name, cat.count)
		}
		fmt.Println()
	}

	// Goals progress
	goals, err := goalRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get goals: %w", err)
	}

	completedGoals := 0
	upcomingTargets := 0
	now := time.Now()
	for _, goal := range goals {
		if goal.Status == core.StatusCompleted {
			completedGoals++
		}
		if goal.TargetDate != nil && goal.TargetDate.After(now) && goal.Status == core.StatusActive {
			upcomingTargets++
		}
	}

	if len(goals) > 0 {
		completionRate := float64(completedGoals) / float64(len(goals)) * 100
		fmt.Printf("Goal Completion: %d/%d (%.1f%%)\n", completedGoals, len(goals), completionRate)
		if upcomingTargets > 0 {
			fmt.Printf("  Upcoming targets: %d goals\n", upcomingTargets)
		}
		fmt.Println()
	}

	// Resources progress
	resources, err := resourceRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	completedResources := 0
	inProgressResources := 0
	totalHours := 0.0
	completedHours := 0.0
	for _, resource := range resources {
		totalHours += resource.EstimatedHours
		if resource.Status == core.ResourceCompleted {
			completedResources++
			completedHours += resource.EstimatedHours
		} else if resource.Status == core.ResourceInProgress {
			inProgressResources++
		}
	}

	if len(resources) > 0 {
		fmt.Printf("Learning Resources:\n")
		fmt.Printf("  Completed: %d/%d resources\n", completedResources, len(resources))
		if totalHours > 0 {
			fmt.Printf("  Hours completed: %.1f/%.1f (%.1f%%)\n", completedHours, totalHours, completedHours/totalHours*100)
		}
		if inProgressResources > 0 {
			fmt.Printf("  In progress: %d resources\n", inProgressResources)
		}
		fmt.Println()
	}

	// Milestones
	milestones, err := milestoneRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get milestones: %w", err)
	}

	achievedMilestones := 0
	recentAchievements := 0
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	for _, milestone := range milestones {
		if milestone.IsAchieved() {
			achievedMilestones++
			if milestone.AchievedDate != nil && milestone.AchievedDate.After(thirtyDaysAgo) {
				recentAchievements++
			}
		}
	}

	if len(milestones) > 0 {
		fmt.Printf("Milestones: %d/%d achieved\n", achievedMilestones, len(milestones))
		if recentAchievements > 0 {
			fmt.Printf("  Recent (last 30 days): %d milestones\n", recentAchievements)
		}
		fmt.Println()
	}

	// Progress logs
	progressLogs, err := progressRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get progress logs: %w", err)
	}

	if len(progressLogs) > 0 {
		totalProgressHours := 0.0
		recentWeeks := 0
		fourWeeksAgo := now.AddDate(0, 0, -28)
		recentHours := 0.0

		for _, log := range progressLogs {
			totalProgressHours += log.HoursInvested
			if log.Date.After(fourWeeksAgo) {
				recentWeeks++
				recentHours += log.HoursInvested
			}
		}

		fmt.Printf("Progress Tracking:\n")
		fmt.Printf("  Total logs: %d\n", len(progressLogs))
		fmt.Printf("  Total hours invested: %.1f\n", totalProgressHours)
		if len(progressLogs) > 0 {
			avgHours := totalProgressHours / float64(len(progressLogs))
			fmt.Printf("  Average per log: %.1f hours\n", avgHours)
		}
		if recentWeeks > 0 {
			avgRecentHours := recentHours / float64(recentWeeks)
			fmt.Printf("  Recent (last 4 weeks): %.1f hours/log\n", avgRecentHours)
		}
		fmt.Println()
	}

	// Learning velocity
	if len(progressLogs) > 0 && len(resources) > 0 {
		fmt.Println("Learning Velocity:")

		// Calculate skills worked on
		skillsWorked := make(map[core.EntityID]bool)
		for _, log := range progressLogs {
			for _, skillID := range log.SkillsWorked {
				skillsWorked[skillID] = true
			}
		}
		fmt.Printf("  Active skills: %d/%d\n", len(skillsWorked), len(skills))

		// Calculate resources completion rate
		thirtyDaysAgo := now.AddDate(0, 0, -30)
		recentCompletions := 0
		for _, resource := range resources {
			if resource.Status == core.ResourceCompleted && resource.Updated.After(thirtyDaysAgo) {
				recentCompletions++
			}
		}
		if recentCompletions > 0 {
			fmt.Printf("  Resources completed (last 30 days): %d\n", recentCompletions)
		}
		fmt.Println()
	}

	return nil
}
