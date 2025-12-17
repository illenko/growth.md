package cli

import (
	"fmt"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Display repository overview",
	Long: `Display a high-level overview of your growth repository.

Shows counts and status of all entities: skills, goals, resources, paths, milestones, and progress logs.

Examples:
  growth overview`,
	RunE: runOverview,
}

func init() {
	rootCmd.AddCommand(overviewCmd)
}

func runOverview(cmd *cobra.Command, args []string) error {
	fmt.Println("Growth Repository Overview")
	fmt.Println("==========================")
	fmt.Println()

	// Skills
	skills, err := skillRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get skills: %w", err)
	}

	skillsByLevel := make(map[core.ProficiencyLevel]int)
	skillsByStatus := make(map[core.SkillStatus]int)
	for _, skill := range skills {
		skillsByLevel[skill.Level]++
		skillsByStatus[skill.Status]++
	}

	fmt.Printf("Skills: %d total\n", len(skills))
	if len(skills) > 0 {
		fmt.Printf("  Beginner: %d | Intermediate: %d | Advanced: %d | Expert: %d\n",
			skillsByLevel[core.LevelBeginner],
			skillsByLevel[core.LevelIntermediate],
			skillsByLevel[core.LevelAdvanced],
			skillsByLevel[core.LevelExpert])
		fmt.Printf("  Not Started: %d | Learning: %d | Mastered: %d\n",
			skillsByStatus[core.SkillNotStarted],
			skillsByStatus[core.SkillLearning],
			skillsByStatus[core.SkillMastered])
	}
	fmt.Println()

	// Goals
	goals, err := goalRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get goals: %w", err)
	}

	goalsByPriority := make(map[core.Priority]int)
	goalsByStatus := make(map[core.Status]int)
	for _, goal := range goals {
		goalsByPriority[goal.Priority]++
		goalsByStatus[goal.Status]++
	}

	fmt.Printf("Goals: %d total\n", len(goals))
	if len(goals) > 0 {
		fmt.Printf("  High: %d | Medium: %d | Low: %d\n",
			goalsByPriority[core.PriorityHigh],
			goalsByPriority[core.PriorityMedium],
			goalsByPriority[core.PriorityLow])
		fmt.Printf("  Active: %d | Completed: %d | Archived: %d\n",
			goalsByStatus[core.StatusActive],
			goalsByStatus[core.StatusCompleted],
			goalsByStatus[core.StatusArchived])
	}
	fmt.Println()

	// Resources
	resources, err := resourceRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	resourcesByType := make(map[core.ResourceType]int)
	resourcesByStatus := make(map[core.ResourceStatus]int)
	totalHours := 0.0
	for _, resource := range resources {
		resourcesByType[resource.Type]++
		resourcesByStatus[resource.Status]++
		totalHours += resource.EstimatedHours
	}

	fmt.Printf("Resources: %d total (%.1f hours estimated)\n", len(resources), totalHours)
	if len(resources) > 0 {
		fmt.Printf("  Books: %d | Courses: %d | Videos: %d | Articles: %d | Projects: %d | Docs: %d\n",
			resourcesByType[core.ResourceBook],
			resourcesByType[core.ResourceCourse],
			resourcesByType[core.ResourceVideo],
			resourcesByType[core.ResourceArticle],
			resourcesByType[core.ResourceProject],
			resourcesByType[core.ResourceDocumentation])
		fmt.Printf("  Not Started: %d | In Progress: %d | Completed: %d\n",
			resourcesByStatus[core.ResourceNotStarted],
			resourcesByStatus[core.ResourceInProgress],
			resourcesByStatus[core.ResourceCompleted])
	}
	fmt.Println()

	// Paths
	paths, err := pathRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get paths: %w", err)
	}

	pathsByType := make(map[core.PathType]int)
	pathsByStatus := make(map[core.Status]int)
	for _, path := range paths {
		pathsByType[path.Type]++
		pathsByStatus[path.Status]++
	}

	fmt.Printf("Learning Paths: %d total\n", len(paths))
	if len(paths) > 0 {
		fmt.Printf("  Manual: %d | AI-Generated: %d\n",
			pathsByType[core.PathTypeManual],
			pathsByType[core.PathTypeAIGenerated])
		fmt.Printf("  Active: %d | Completed: %d | Archived: %d\n",
			pathsByStatus[core.StatusActive],
			pathsByStatus[core.StatusCompleted],
			pathsByStatus[core.StatusArchived])
	}
	fmt.Println()

	// Milestones
	milestones, err := milestoneRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get milestones: %w", err)
	}

	milestonesByType := make(map[core.MilestoneType]int)
	milestonesAchieved := 0
	for _, milestone := range milestones {
		milestonesByType[milestone.Type]++
		if milestone.IsAchieved() {
			milestonesAchieved++
		}
	}

	fmt.Printf("Milestones: %d total (%d achieved)\n", len(milestones), milestonesAchieved)
	if len(milestones) > 0 {
		fmt.Printf("  Goal-level: %d | Path-level: %d | Skill-level: %d\n",
			milestonesByType[core.MilestoneGoalLevel],
			milestonesByType[core.MilestonePathLevel],
			milestonesByType[core.MilestoneSkillLevel])
	}
	fmt.Println()

	// Progress Logs
	progressLogs, err := progressRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get progress logs: %w", err)
	}

	totalProgressHours := 0.0
	for _, log := range progressLogs {
		totalProgressHours += log.HoursInvested
	}

	fmt.Printf("Progress Logs: %d total (%.1f hours logged)\n", len(progressLogs), totalProgressHours)
	fmt.Println()

	return nil
}
