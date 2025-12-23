package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	searchType string
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search across all entities",
	Long: `Search for skills, goals, resources, paths, milestones, and progress logs.

The search looks through titles, descriptions, tags, and other text fields.

Examples:
  growth search python
  growth search "backend development"
  growth search docker --type skill
  growth search goal --type goal`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&searchType, "type", "t", "", "filter by entity type (skill, goal, resource, path, milestone, progress)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	if searchType != "" {
		return searchByType(query, searchType)
	}

	// Search all types
	PrintInfo(fmt.Sprintf("Searching for: %s", query))
	fmt.Println()

	hasResults := false

	// Search skills
	skills, err := skillRepo.Search(query)
	if err == nil && len(skills) > 0 {
		fmt.Printf("Skills (%d):\n", len(skills))
		for _, skill := range skills {
			fmt.Printf("  %s - %s (%s, %s)\n", skill.ID, skill.Title, skill.Category, skill.Level)
		}
		fmt.Println()
		hasResults = true
	}

	// Search goals
	goals, err := goalRepo.Search(query)
	if err == nil && len(goals) > 0 {
		fmt.Printf("Goals (%d):\n", len(goals))
		for _, goal := range goals {
			fmt.Printf("  %s - %s (%s priority, %s)\n", goal.ID, goal.Title, goal.Priority, goal.Status)
		}
		fmt.Println()
		hasResults = true
	}

	// Search resources
	resources, err := resourceRepo.Search(query)
	if err == nil && len(resources) > 0 {
		fmt.Printf("Resources (%d):\n", len(resources))
		for _, resource := range resources {
			fmt.Printf("  %s - %s (%s, skill: %s)\n", resource.ID, resource.Title, resource.Type, resource.SkillID)
		}
		fmt.Println()
		hasResults = true
	}

	// Search paths
	paths, err := pathRepo.Search(query)
	if err == nil && len(paths) > 0 {
		fmt.Printf("Paths (%d):\n", len(paths))
		for _, path := range paths {
			fmt.Printf("  %s - %s (%s, %s)\n", path.ID, path.Title, path.Type, path.Status)
		}
		fmt.Println()
		hasResults = true
	}

	// Search milestones
	milestones, err := milestoneRepo.Search(query)
	if err == nil && len(milestones) > 0 {
		fmt.Printf("Milestones (%d):\n", len(milestones))
		for _, milestone := range milestones {
			fmt.Printf("  %s - %s (%s, ref: %s)\n", milestone.ID, milestone.Title, milestone.Type, milestone.ReferenceID)
		}
		fmt.Println()
		hasResults = true
	}

	// Search progress logs
	progressLogs, err := progressRepo.Search(query)
	if err == nil && len(progressLogs) > 0 {
		fmt.Printf("Progress Logs (%d):\n", len(progressLogs))
		for _, log := range progressLogs {
			fmt.Printf("  %s - %s (%.1f hours)\n", log.ID, log.Date.Format("2006-01-02"), log.HoursInvested)
		}
		fmt.Println()
		hasResults = true
	}

	if !hasResults {
		PrintInfo("No results found")
	}

	return nil
}

func searchByType(query, entityType string) error {
	entityType = strings.ToLower(entityType)

	switch entityType {
	case "skill", "skills":
		skills, err := skillRepo.Search(query)
		if err != nil {
			return fmt.Errorf("search failed: %w\nTry running 'growth skill list' to see all skills", err)
		}
		if len(skills) == 0 {
			PrintInfo("No skills found")
			return nil
		}
		return PrintOutputWithConfig(skills)

	case "goal", "goals":
		goals, err := goalRepo.Search(query)
		if err != nil {
			return fmt.Errorf("search failed: %w\nTry running 'growth goal list' to see all goals", err)
		}
		if len(goals) == 0 {
			PrintInfo("No goals found")
			return nil
		}
		return PrintOutputWithConfig(goals)

	case "resource", "resources":
		resources, err := resourceRepo.Search(query)
		if err != nil {
			return fmt.Errorf("search failed: %w\nTry running 'growth resource list' to see all resources", err)
		}
		if len(resources) == 0 {
			PrintInfo("No resources found")
			return nil
		}
		return PrintOutputWithConfig(resources)

	case "path", "paths":
		paths, err := pathRepo.Search(query)
		if err != nil {
			return fmt.Errorf("search failed: %w\nTry running 'growth path list' to see all paths", err)
		}
		if len(paths) == 0 {
			PrintInfo("No paths found")
			return nil
		}
		return PrintOutputWithConfig(paths)

	case "milestone", "milestones":
		milestones, err := milestoneRepo.Search(query)
		if err != nil {
			return fmt.Errorf("search failed: %w\nTry running 'growth milestone list' to see all milestones", err)
		}
		if len(milestones) == 0 {
			PrintInfo("No milestones found")
			return nil
		}
		return PrintOutputWithConfig(milestones)

	case "progress":
		progressLogs, err := progressRepo.Search(query)
		if err != nil {
			return fmt.Errorf("search failed: %w\nTry running 'growth progress list' to see all progress logs", err)
		}
		if len(progressLogs) == 0 {
			PrintInfo("No progress logs found")
			return nil
		}
		return PrintOutputWithConfig(progressLogs)

	default:
		return fmt.Errorf("unknown entity type '%s'. Valid options: skill, goal, resource, path, milestone, progress", entityType)
	}
}
