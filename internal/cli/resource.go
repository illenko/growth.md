package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	resourceType       string
	resourceSkillID    string
	resourceStatus     string
	resourceURL        string
	resourceAuthor     string
	resourceHours      string
	resourceTags       string
	resourceTitle      string
	resourceFilterType string
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Manage learning resources",
	Long:  `Create, list, view, edit, and delete learning resources (books, courses, videos, etc.).`,
}

var resourceCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new resource",
	Long: `Create a new learning resource with the specified title.

You can provide the title as an argument or be prompted for it.
A resource must be associated with a skill using --skill-id.

Examples:
  growth resource create "Clean Code" --skill-id skill-001 --type book --author "Robert Martin"
  growth resource create "Python Course" --skill-id skill-002 --type course --url https://example.com
  growth resource create`,
	Args: cobra.MaximumNArgs(1),
	RunE: runResourceCreate,
}

var resourceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all resources",
	Long: `List all resources in the repository.

Optionally filter by skill, type, or status using flags.

Examples:
  growth resource list
  growth resource list --skill-id skill-001
  growth resource list --type book
  growth resource list --status in-progress`,
	Aliases: []string{"ls"},
	RunE:    runResourceList,
}

var resourceViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View resource details",
	Long: `View detailed information about a specific resource.

The output format can be controlled with the --format flag (table, json, yaml).

Examples:
  growth resource view resource-001
  growth resource view resource-042 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runResourceView,
}

var resourceEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing resource",
	Long: `Edit an existing resource by ID.

You can update any field using flags. If no flags are provided, you'll be prompted
to update each field interactively (press Enter to keep current value).

Examples:
  growth resource edit resource-001 --status in-progress
  growth resource edit resource-042 --url https://example.com --hours 40
  growth resource edit resource-001`,
	Args: cobra.ExactArgs(1),
	RunE: runResourceEdit,
}

var resourceDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a resource",
	Long: `Delete a resource by ID.

This will permanently remove the resource file. You'll be prompted for confirmation
before deletion.

Examples:
  growth resource delete resource-001
  growth resource delete resource-042`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runResourceDelete,
}

var resourceStartCmd = &cobra.Command{
	Use:   "start <id>",
	Short: "Mark resource as in-progress",
	Long: `Update a resource status to in-progress.

Examples:
  growth resource start resource-001`,
	Args: cobra.ExactArgs(1),
	RunE: runResourceStart,
}

var resourceCompleteCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "Mark resource as completed",
	Long: `Update a resource status to completed.

Examples:
  growth resource complete resource-001`,
	Args: cobra.ExactArgs(1),
	RunE: runResourceComplete,
}

func init() {
	rootCmd.AddCommand(resourceCmd)
	resourceCmd.AddCommand(resourceCreateCmd)
	resourceCmd.AddCommand(resourceListCmd)
	resourceCmd.AddCommand(resourceViewCmd)
	resourceCmd.AddCommand(resourceEditCmd)
	resourceCmd.AddCommand(resourceDeleteCmd)
	resourceCmd.AddCommand(resourceStartCmd)
	resourceCmd.AddCommand(resourceCompleteCmd)

	resourceCreateCmd.Flags().StringVar(&resourceSkillID, "skill-id", "", "skill ID (required)")
	resourceCreateCmd.Flags().StringVarP(&resourceType, "type", "t", "", "resource type (book, course, video, article, project, documentation)")
	resourceCreateCmd.Flags().StringVar(&resourceURL, "url", "", "resource URL")
	resourceCreateCmd.Flags().StringVar(&resourceAuthor, "author", "", "resource author")
	resourceCreateCmd.Flags().StringVar(&resourceHours, "hours", "", "estimated hours")
	resourceCreateCmd.Flags().StringVar(&resourceTags, "tags", "", "comma-separated tags")
	resourceCreateCmd.MarkFlagRequired("skill-id")

	resourceListCmd.Flags().StringVar(&resourceSkillID, "skill-id", "", "filter by skill ID")
	resourceListCmd.Flags().StringVarP(&resourceFilterType, "type", "t", "", "filter by type")
	resourceListCmd.Flags().StringVarP(&resourceStatus, "status", "s", "", "filter by status")

	resourceEditCmd.Flags().StringVar(&resourceTitle, "title", "", "resource title")
	resourceEditCmd.Flags().StringVarP(&resourceType, "type", "t", "", "resource type")
	resourceEditCmd.Flags().StringVar(&resourceURL, "url", "", "resource URL")
	resourceEditCmd.Flags().StringVar(&resourceAuthor, "author", "", "resource author")
	resourceEditCmd.Flags().StringVar(&resourceHours, "hours", "", "estimated hours")
	resourceEditCmd.Flags().StringVarP(&resourceStatus, "status", "s", "", "resource status")
	resourceEditCmd.Flags().StringVar(&resourceTags, "tags", "", "comma-separated tags")
}

func runResourceCreate(cmd *cobra.Command, args []string) error {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = PromptStringRequired("Resource title")
	}

	if resourceSkillID == "" {
		resourceSkillID = PromptStringRequired("Skill ID (e.g., skill-001)")
	}

	skillID := core.EntityID(resourceSkillID)
	exists, err := skillRepo.Exists(skillID)
	if err != nil {
		return fmt.Errorf("failed to check skill existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("skill '%s' not found. Use 'growth skill list' to see available skills", skillID)
	}

	if resourceType == "" {
		resourceType = PromptSelectWithDefault(
			"Resource type",
			[]string{"book", "course", "video", "article", "project", "documentation"},
			"book",
		)
	}

	resType := core.ResourceType(resourceType)
	if !resType.IsValid() {
		return fmt.Errorf("invalid resource type '%s'. Valid options: book, course, video, article, project, documentation", resourceType)
	}

	id, err := GenerateNextID("resource")
	if err != nil {
		return fmt.Errorf("failed to generate resource ID: %w", err)
	}

	resource, err := core.NewResource(id, title, resType, skillID)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	if resourceURL != "" {
		resource.SetURL(resourceURL)
	}

	if resourceAuthor != "" {
		resource.SetAuthor(resourceAuthor)
	}

	if resourceHours != "" {
		hours, err := strconv.ParseFloat(resourceHours, 64)
		if err != nil {
			return fmt.Errorf("invalid hours value: %w", err)
		}
		if err := resource.SetEstimatedHours(hours); err != nil {
			return fmt.Errorf("failed to set estimated hours: %w", err)
		}
	}

	if resourceTags != "" {
		tags := strings.Split(resourceTags, ",")
		for _, tag := range tags {
			resource.AddTag(strings.TrimSpace(tag))
		}
	}

	notes := PromptMultiline("Notes (optional, press Ctrl+D or enter '.' to finish)")
	if notes != "" {
		resource.Body = notes
	}

	if err := resourceRepo.Create(resource); err != nil {
		return fmt.Errorf("failed to save resource: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Created resource %s: %s", resource.ID, resource.Title))

	if verbose {
		fmt.Printf("\nResource details:\n")
		fmt.Printf("  ID: %s\n", resource.ID)
		fmt.Printf("  Title: %s\n", resource.Title)
		fmt.Printf("  Type: %s\n", resource.Type)
		fmt.Printf("  Skill: %s\n", resource.SkillID)
		if resource.URL != "" {
			fmt.Printf("  URL: %s\n", resource.URL)
		}
		if resource.Author != "" {
			fmt.Printf("  Author: %s\n", resource.Author)
		}
	}

	return nil
}

func runResourceList(cmd *cobra.Command, args []string) error {
	var resources []*core.Resource
	var err error

	if resourceSkillID != "" {
		skillID := core.EntityID(resourceSkillID)
		resources, err = resourceRepo.FindBySkillID(skillID)
	} else if resourceFilterType != "" {
		resType := core.ResourceType(resourceFilterType)
		if !resType.IsValid() {
			return fmt.Errorf("invalid resource type '%s'. Valid options: book, course, video, article, project, documentation", resourceFilterType)
		}
		resources, err = resourceRepo.FindByType(resType)
	} else if resourceStatus != "" {
		status := core.ResourceStatus(resourceStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid resource status '%s'. Valid options: not-started, in-progress, completed", resourceStatus)
		}
		resources, err = resourceRepo.FindByStatus(status)
	} else {
		resources, err = resourceRepo.GetAll()
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve resources: %w\nTry running 'growth resource list' without filters to see all resources", err)
	}

	if len(resources) == 0 {
		PrintInfo("No resources found")
		return nil
	}

	return PrintOutputWithConfig(resources)
}

func runResourceView(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	resource, err := resourceRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("resource '%s' not found. Use 'growth resource list' to see available resources", id)
	}

	if config.Display.OutputFormat == "table" {
		fmt.Printf("ID:       %s\n", resource.ID)
		fmt.Printf("Title:    %s\n", resource.Title)
		fmt.Printf("Type:     %s\n", resource.Type)
		fmt.Printf("Skill:    %s\n", resource.SkillID)
		fmt.Printf("Status:   %s\n", resource.Status)
		if resource.URL != "" {
			fmt.Printf("URL:      %s\n", resource.URL)
		}
		if resource.Author != "" {
			fmt.Printf("Author:   %s\n", resource.Author)
		}
		if resource.EstimatedHours > 0 {
			fmt.Printf("Hours:    %.1f\n", resource.EstimatedHours)
		}
		if len(resource.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(resource.Tags, ", "))
		}
		fmt.Printf("Created:  %s\n", resource.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:  %s\n", resource.Updated.Format("2006-01-02 15:04:05"))

		if resource.Body != "" {
			fmt.Printf("\nNotes:\n%s\n", resource.Body)
		}

		return nil
	}

	return PrintOutputWithConfig(resource)
}

func runResourceEdit(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	resource, err := resourceRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("resource '%s' not found. Use 'growth resource list' to see available resources", id)
	}

	updated := false

	if cmd.Flags().Changed("title") {
		resource.Title = resourceTitle
		updated = true
	}

	if cmd.Flags().Changed("type") {
		resType := core.ResourceType(resourceType)
		if !resType.IsValid() {
			return fmt.Errorf("invalid resource type '%s'. Valid options: book, course, video, article, project, documentation", resourceType)
		}
		resource.Type = resType
		updated = true
	}

	if cmd.Flags().Changed("url") {
		resource.SetURL(resourceURL)
		updated = true
	}

	if cmd.Flags().Changed("author") {
		resource.SetAuthor(resourceAuthor)
		updated = true
	}

	if cmd.Flags().Changed("hours") {
		hours, err := strconv.ParseFloat(resourceHours, 64)
		if err != nil {
			return fmt.Errorf("invalid hours value: %w", err)
		}
		if err := resource.SetEstimatedHours(hours); err != nil {
			return fmt.Errorf("failed to set estimated hours: %w", err)
		}
		updated = true
	}

	if cmd.Flags().Changed("status") {
		status := core.ResourceStatus(resourceStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid resource status '%s'. Valid options: not-started, in-progress, completed", resourceStatus)
		}
		if err := resource.UpdateStatus(status); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
		updated = true
	}

	if cmd.Flags().Changed("tags") {
		resource.Tags = []string{}
		if resourceTags != "" {
			tags := strings.Split(resourceTags, ",")
			for _, tag := range tags {
				resource.AddTag(strings.TrimSpace(tag))
			}
		}
		updated = true
	}

	if !updated {
		PrintInfo("No changes specified. Use flags to update fields.")
		return nil
	}

	if err := resourceRepo.Update(resource); err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Updated resource %s: %s", resource.ID, resource.Title))
	return nil
}

func runResourceDelete(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	resource, err := resourceRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("resource '%s' not found. Use 'growth resource list' to see available resources", id)
	}

	fmt.Printf("You are about to delete:\n")
	fmt.Printf("  ID: %s\n", resource.ID)
	fmt.Printf("  Title: %s\n", resource.Title)
	fmt.Printf("  Type: %s\n", resource.Type)
	fmt.Println()

	if !PromptConfirm("Are you sure you want to delete this resource?") {
		PrintInfo("Deletion cancelled")
		return nil
	}

	if err := resourceRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Deleted resource %s", id))
	return nil
}

func runResourceStart(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	resource, err := resourceRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("resource '%s' not found. Use 'growth resource list' to see available resources", id)
	}

	resource.Start()

	if err := resourceRepo.Update(resource); err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Started resource %s: %s", resource.ID, resource.Title))
	return nil
}

func runResourceComplete(cmd *cobra.Command, args []string) error {
	id := core.EntityID(args[0])

	resource, err := resourceRepo.GetByIDWithBody(id)
	if err != nil {
		return fmt.Errorf("resource '%s' not found. Use 'growth resource list' to see available resources", id)
	}

	resource.Complete()

	if err := resourceRepo.Update(resource); err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Completed resource %s: %s", resource.ID, resource.Title))
	return nil
}
