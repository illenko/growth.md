package cli

import (
	"fmt"

	"github.com/illenko/growth.md/internal/git"
	"github.com/spf13/cobra"
)

var (
	gitLogCount int
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git operations for the growth repository",
	Long:  `View git status and log for the growth repository.`,
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status",
	Long: `Show the current git status of the growth repository.

Examples:
  growth git status`,
	RunE: runGitStatus,
}

var gitLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show git commit log",
	Long: `Show recent git commits in the growth repository.

Examples:
  growth git log
  growth git log --count 20`,
	RunE: runGitLog,
}

func init() {
	rootCmd.AddCommand(gitCmd)
	gitCmd.AddCommand(gitStatusCmd)
	gitCmd.AddCommand(gitLogCmd)

	gitLogCmd.Flags().IntVarP(&gitLogCount, "count", "n", 10, "number of commits to show")
}

func runGitStatus(cmd *cobra.Command, args []string) error {
	// Check if git is installed
	if err := git.EnsureGitInstalled(); err != nil {
		return fmt.Errorf("git is not installed: %w", err)
	}

	// Check if we're in a git repository
	if !git.IsRepo(repoPath) {
		PrintWarning("Not a git repository")
		fmt.Println("\nTo initialize git for this repository, run:")
		fmt.Println("  cd", repoPath)
		fmt.Println("  git init")
		return nil
	}

	// Get git status
	status, err := git.Status(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	// Get current branch
	branch, err := git.GetCurrentBranch(repoPath)
	if err != nil {
		branch = "unknown"
	}

	// Display status
	fmt.Printf("On branch: %s\n\n", branch)

	if len(status) == 0 {
		PrintSuccess("Working tree is clean")
		return nil
	}

	fmt.Printf("Changes (%d files):\n", len(status))
	for _, line := range status {
		fmt.Printf("  %s\n", line)
	}

	return nil
}

func runGitLog(cmd *cobra.Command, args []string) error {
	// Check if git is installed
	if err := git.EnsureGitInstalled(); err != nil {
		return fmt.Errorf("git is not installed: %w", err)
	}

	// Check if we're in a git repository
	if !git.IsRepo(repoPath) {
		PrintWarning("Not a git repository")
		return nil
	}

	// Get git log
	commits, err := git.Log(repoPath, gitLogCount)
	if err != nil {
		return fmt.Errorf("failed to get git log: %w", err)
	}

	if len(commits) == 0 {
		PrintInfo("No commits yet")
		return nil
	}

	fmt.Printf("Recent commits (%d):\n", len(commits))
	for _, commit := range commits {
		fmt.Printf("  %s\n", commit)
	}

	return nil
}
