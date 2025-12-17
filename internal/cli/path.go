package cli

import (
	"fmt"
	"strings"

	"github.com/illenko/growth.md/internal/core"
	"github.com/spf13/cobra"
)

var (
	pathType       string
	pathStatus     string
	pathTags       string
	pathTitle      string
	pathFilterType string
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

func init() {
	rootCmd.AddCommand(pathCmd)
	pathCmd.AddCommand(pathCreateCmd)
	pathCmd.AddCommand(pathListCmd)
	pathCmd.AddCommand(pathViewCmd)
	pathCmd.AddCommand(pathEditCmd)
	pathCmd.AddCommand(pathDeleteCmd)

	pathCreateCmd.Flags().StringVarP(&pathType, "type", "t", "", "path type (manual, ai-generated)")
	pathCreateCmd.Flags().StringVar(&pathTags, "tags", "", "comma-separated tags")

	pathListCmd.Flags().StringVarP(&pathFilterType, "type", "t", "", "filter by type")
	pathListCmd.Flags().StringVarP(&pathStatus, "status", "s", "", "filter by status")

	pathEditCmd.Flags().StringVar(&pathTitle, "title", "", "path title")
	pathEditCmd.Flags().StringVarP(&pathStatus, "status", "s", "", "path status")
	pathEditCmd.Flags().StringVar(&pathTags, "tags", "", "comma-separated tags")
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
			return fmt.Errorf("invalid path type: %s", pathFilterType)
		}
		paths, err = pathRepo.FindByType(pType)
	} else if pathStatus != "" {
		status := core.Status(pathStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status: %s", pathStatus)
		}
		paths, err = pathRepo.FindByStatus(status)
	} else {
		paths, err = pathRepo.GetAll()
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve paths: %w", err)
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
		return fmt.Errorf("failed to retrieve path: %w", err)
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
		return fmt.Errorf("failed to retrieve path: %w", err)
	}

	updated := false

	if cmd.Flags().Changed("title") {
		path.Title = pathTitle
		updated = true
	}

	if cmd.Flags().Changed("status") {
		status := core.Status(pathStatus)
		if !status.IsValid() {
			return fmt.Errorf("invalid status: %s", pathStatus)
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
		return fmt.Errorf("failed to retrieve path: %w", err)
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
