package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/illenko/growth.md/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new growth repository",
	Long: `Initialize a new growth.md repository in the specified directory (or current directory).

This will:
- Create the directory structure (skills/, goals/, paths/, etc.)
- Initialize a Git repository
- Create a default config.yml
- Make an initial commit`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := createDirectoryStructure(absPath); err != nil {
		return err
	}

	config, err := promptForConfig()
	if err != nil {
		return err
	}

	configPath := filepath.Join(absPath, ".growth", "config.yml")
	if err := storage.SaveConfig(config, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := createGitignore(absPath); err != nil {
		return err
	}

	if err := createReadme(absPath); err != nil {
		return err
	}

	if err := initializeGit(absPath); err != nil {
		return err
	}

	fmt.Printf("\n✓ Initialized growth.md repository in %s\n", absPath)
	fmt.Println("\nNext steps:")
	fmt.Println("  cd", targetDir)
	fmt.Println("  growth skill create \"Your First Skill\" --category programming")
	fmt.Println("  growth goal create \"Your First Goal\" --priority high")
	fmt.Println("\nRun 'growth --help' to see all available commands.")

	return nil
}

func createDirectoryStructure(basePath string) error {
	dirs := []string{
		".growth",
		"skills",
		"goals",
		"paths",
		"phases",
		"resources",
		"milestones",
		"progress",
	}

	for _, dir := range dirs {
		path := filepath.Join(basePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return nil
}

func promptForConfig() (*storage.Config, error) {
	reader := bufio.NewReader(os.Stdin)
	config := storage.DefaultConfig()

	fmt.Println("\n📝 Let's set up your growth.md configuration")
	fmt.Println()

	fmt.Print("Your name (optional): ")
	name, _ := reader.ReadString('\n')
	config.User.Name = strings.TrimSpace(name)

	fmt.Print("Your email (optional): ")
	email, _ := reader.ReadString('\n')
	config.User.Email = strings.TrimSpace(email)

	fmt.Print("\nAI Provider (gemini/openai/anthropic/local) [gemini]: ")
	provider, _ := reader.ReadString('\n')
	provider = strings.TrimSpace(provider)
	if provider != "" {
		config.AI.Provider = provider
	}

	if config.AI.Provider == "gemini" {
		config.AI.Model = "gemini-3-flash-preview"
	} else if config.AI.Provider == "openai" {
		config.AI.Model = "gpt-4"
	} else if config.AI.Provider == "anthropic" {
		config.AI.Model = "claude-3-5-sonnet-20241022"
	}

	fmt.Print("\nEnable auto-commit to Git? (y/n) [n]: ")
	autoCommit, _ := reader.ReadString('\n')
	autoCommit = strings.TrimSpace(strings.ToLower(autoCommit))
	if autoCommit == "y" || autoCommit == "yes" {
		config.Git.AutoCommit = true
		config.Git.CommitOnUpdate = true
	}

	return config, nil
}

func createGitignore(basePath string) error {
	content := `# growth.md specific
.growth/cache/
.DS_Store

# Editor files
.vscode/
.idea/
*.swp
*.swo
*~
`

	path := filepath.Join(basePath, ".gitignore")
	return os.WriteFile(path, []byte(content), 0644)
}

func createReadme(basePath string) error {
	content := `# My Growth Journey

This repository tracks my career development using [growth.md](https://github.com/illenko/growth.md).

## Structure

- **skills/** - Technical skills I'm learning or have mastered
- **goals/** - Career objectives and targets
- **paths/** - Learning paths to achieve goals
- **phases/** - Phases within learning paths
- **resources/** - Books, courses, articles, and other learning materials
- **milestones/** - Achievement markers
- **progress/** - Weekly progress logs

## Quick Start

View your skills:
` + "```" + `bash
growth skill list
` + "```" + `

Create a new goal:
` + "```" + `bash
growth goal create "Senior Engineer by 2025" --priority high
` + "```" + `

Log your weekly progress:
` + "```" + `bash
growth progress log "Completed Python tutorial"
` + "```" + `

## Configuration

Configuration is stored in ` + "`" + `.growth/config.yml` + "`" + `. Edit this file to customize behavior.
`

	path := filepath.Join(basePath, "README.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func initializeGit(basePath string) error {
	if isGitRepo(basePath) {
		fmt.Println("\n⚠  Git repository already exists, skipping git init")
		return nil
	}

	commands := [][]string{
		{"git", "init"},
		{"git", "add", "."},
		{"git", "commit", "-m", "Initial commit: Initialize growth.md repository"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = basePath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git command failed (%s): %w", strings.Join(cmdArgs, " "), err)
		}
	}

	return nil
}

func isGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}
