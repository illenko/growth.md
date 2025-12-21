package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// InitRepo initializes a new git repository at the specified path
func InitRepo(path string) error {
	if path == "" {
		return fmt.Errorf("repository path cannot be empty")
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// IsRepo checks if the specified path is inside a git repository
func IsRepo(path string) bool {
	if path == "" {
		return false
	}

	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	err := cmd.Run()
	return err == nil
}

// GetRepoRoot returns the root directory of the git repository
func GetRepoRoot(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository (or any of the parent directories)")
	}

	root := strings.TrimSpace(string(output))

	// Resolve symlinks (e.g., /var -> /private/var on macOS)
	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		return root, nil // Return original if symlink resolution fails
	}

	return resolved, nil
}

// Status returns the list of modified/untracked files in the repository
func Status(repoPath string) ([]string, error) {
	if !IsRepo(repoPath) {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

// Add stages files for commit
func Add(repoPath string, files []string) error {
	if !IsRepo(repoPath) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	if len(files) == 0 {
		return nil
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add files: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Commit creates a commit with the specified message and files
// If files is empty, commits all staged changes
func Commit(repoPath string, message string, files []string) error {
	if !IsRepo(repoPath) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	// Stage files if provided
	if len(files) > 0 {
		if err := Add(repoPath, files); err != nil {
			return err
		}
	}

	// Commit
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's just "nothing to commit"
		if strings.Contains(string(output), "nothing to commit") {
			return nil // Not an error, just nothing changed
		}
		return fmt.Errorf("failed to commit: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Log returns the last n commit messages
func Log(repoPath string, count int) ([]string, error) {
	if !IsRepo(repoPath) {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	if count <= 0 {
		count = 10
	}

	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", count), "--oneline")
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Empty repo has no commits - check both error and output
		outputStr := string(output)
		if strings.Contains(outputStr, "does not have any commits yet") ||
			strings.Contains(outputStr, "your current branch") ||
			strings.Contains(err.Error(), "does not have any commits yet") {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get git log: %w\nOutput: %s", err, outputStr)
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

// CommitFile is a convenience function that commits a single file with a message
func CommitFile(repoPath string, filePath string, message string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Convert to relative path if it's absolute
	if filepath.IsAbs(filePath) {
		repoRoot, err := GetRepoRoot(repoPath)
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(repoRoot, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		filePath = relPath
	}

	return Commit(repoPath, message, []string{filePath})
}

// HasUncommittedChanges checks if there are uncommitted changes in the repository
func HasUncommittedChanges(repoPath string) (bool, error) {
	status, err := Status(repoPath)
	if err != nil {
		return false, err
	}
	return len(status) > 0, nil
}

// GetCurrentBranch returns the name of the current branch
func GetCurrentBranch(repoPath string) (string, error) {
	if !IsRepo(repoPath) {
		return "", fmt.Errorf("not a git repository: %s", repoPath)
	}

	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// EnsureGitInstalled checks if git is installed and available
func EnsureGitInstalled() error {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git is not installed or not available in PATH")
	}
	return nil
}

// SetConfig sets a git config value (e.g., user.name, user.email)
func SetConfig(repoPath string, key string, value string, global bool) error {
	if key == "" || value == "" {
		return fmt.Errorf("config key and value cannot be empty")
	}

	args := []string{"config"}
	if global {
		args = append(args, "--global")
	}
	args = append(args, key, value)

	cmd := exec.Command("git", args...)
	if !global && repoPath != "" {
		cmd.Dir = repoPath
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set git config: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetConfig gets a git config value
func GetConfig(repoPath string, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("config key cannot be empty")
	}

	cmd := exec.Command("git", "config", "--get", key)
	if repoPath != "" {
		cmd.Dir = repoPath
	}
	output, err := cmd.Output()
	if err != nil {
		// Config key doesn't exist
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("failed to get git config: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
