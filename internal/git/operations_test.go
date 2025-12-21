package git

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "growth-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Resolve symlinks (e.g., /var -> /private/var on macOS)
	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to resolve symlinks: %v", err)
	}

	// Initialize git repo
	if err := InitRepo(tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to init repo: %v", err)
	}

	// Configure git for testing
	SetConfig(tmpDir, "user.name", "Test User", false)
	SetConfig(tmpDir, "user.email", "test@example.com", false)

	return tmpDir
}

func TestInitRepo(t *testing.T) {
	t.Run("initializes repository successfully", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "growth-git-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		err = InitRepo(tmpDir)
		if err != nil {
			t.Errorf("InitRepo() error = %v", err)
		}

		// Check that .git directory exists
		gitDir := filepath.Join(tmpDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			t.Error(".git directory was not created")
		}
	})

	t.Run("fails with empty path", func(t *testing.T) {
		err := InitRepo("")
		if err == nil {
			t.Error("Expected error for empty path, got nil")
		}
	})
}

func TestIsRepo(t *testing.T) {
	t.Run("returns true for git repository", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		if !IsRepo(tmpDir) {
			t.Error("IsRepo() returned false for valid git repository")
		}
	})

	t.Run("returns false for non-repository", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "growth-git-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if IsRepo(tmpDir) {
			t.Error("IsRepo() returned true for non-repository")
		}
	})

	t.Run("returns false for empty path", func(t *testing.T) {
		if IsRepo("") {
			t.Error("IsRepo() returned true for empty path")
		}
	})
}

func TestGetRepoRoot(t *testing.T) {
	t.Run("returns repository root", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		root, err := GetRepoRoot(tmpDir)
		if err != nil {
			t.Errorf("GetRepoRoot() error = %v", err)
		}

		if root != tmpDir {
			t.Errorf("GetRepoRoot() = %v, want %v", root, tmpDir)
		}
	})

	t.Run("works from subdirectory", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		subDir := filepath.Join(tmpDir, "subdir")
		os.MkdirAll(subDir, 0755)

		root, err := GetRepoRoot(subDir)
		if err != nil {
			t.Errorf("GetRepoRoot() error = %v", err)
		}

		if root != tmpDir {
			t.Errorf("GetRepoRoot() = %v, want %v", root, tmpDir)
		}
	})

	t.Run("fails for non-repository", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "growth-git-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		_, err = GetRepoRoot(tmpDir)
		if err == nil {
			t.Error("Expected error for non-repository, got nil")
		}
	})
}

func TestStatus(t *testing.T) {
	t.Run("returns empty list for clean repository", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		status, err := Status(tmpDir)
		if err != nil {
			t.Errorf("Status() error = %v", err)
		}

		if len(status) != 0 {
			t.Errorf("Status() returned %d items, want 0", len(status))
		}
	})

	t.Run("detects untracked file", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create a file
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("test"), 0644)

		status, err := Status(tmpDir)
		if err != nil {
			t.Errorf("Status() error = %v", err)
		}

		if len(status) == 0 {
			t.Error("Status() returned empty list, expected untracked file")
		}
	})

	t.Run("fails for non-repository", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "growth-git-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		_, err = Status(tmpDir)
		if err == nil {
			t.Error("Expected error for non-repository, got nil")
		}
	})
}

func TestAddAndCommit(t *testing.T) {
	t.Run("adds and commits file successfully", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create a file
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("test content"), 0644)

		// Add file
		err := Add(tmpDir, []string{"test.txt"})
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// Commit
		err = Commit(tmpDir, "Test commit", []string{})
		if err != nil {
			t.Errorf("Commit() error = %v", err)
		}

		// Verify clean status
		status, _ := Status(tmpDir)
		if len(status) != 0 {
			t.Errorf("Expected clean status after commit, got %d changes", len(status))
		}
	})

	t.Run("commits with files in one step", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create files
		testFile1 := filepath.Join(tmpDir, "file1.txt")
		testFile2 := filepath.Join(tmpDir, "file2.txt")
		os.WriteFile(testFile1, []byte("content1"), 0644)
		os.WriteFile(testFile2, []byte("content2"), 0644)

		// Commit both files
		err := Commit(tmpDir, "Add two files", []string{"file1.txt", "file2.txt"})
		if err != nil {
			t.Errorf("Commit() error = %v", err)
		}

		// Verify clean status
		status, _ := Status(tmpDir)
		if len(status) != 0 {
			t.Errorf("Expected clean status after commit, got %d changes", len(status))
		}
	})

	t.Run("handles nothing to commit gracefully", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Try to commit with nothing staged
		err := Commit(tmpDir, "Empty commit", []string{})
		if err != nil {
			t.Errorf("Commit() with nothing to commit should not error, got: %v", err)
		}
	})

	t.Run("fails with empty commit message", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		err := Commit(tmpDir, "", []string{})
		if err == nil {
			t.Error("Expected error for empty commit message, got nil")
		}
	})
}

func TestCommitFile(t *testing.T) {
	t.Run("commits single file with relative path", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create file
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("content"), 0644)

		// Commit using relative path
		err := CommitFile(tmpDir, "test.txt", "Add test file")
		if err != nil {
			t.Errorf("CommitFile() error = %v", err)
		}

		// Verify clean status
		status, _ := Status(tmpDir)
		if len(status) != 0 {
			t.Error("Expected clean status after commit")
		}
	})

	t.Run("commits single file with absolute path", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create file
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("content"), 0644)

		// Commit using absolute path
		err := CommitFile(tmpDir, testFile, "Add test file")
		if err != nil {
			t.Errorf("CommitFile() error = %v", err)
		}

		// Verify clean status
		status, _ := Status(tmpDir)
		if len(status) != 0 {
			t.Error("Expected clean status after commit")
		}
	})
}

func TestLog(t *testing.T) {
	t.Run("returns commit history", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create and commit a file
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("content"), 0644)
		Commit(tmpDir, "Initial commit", []string{"test.txt"})

		// Get log
		commits, err := Log(tmpDir, 10)
		if err != nil {
			t.Errorf("Log() error = %v", err)
		}

		if len(commits) != 1 {
			t.Errorf("Log() returned %d commits, want 1", len(commits))
		}
	})

	t.Run("returns empty list for new repository", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		commits, err := Log(tmpDir, 10)
		if err != nil {
			t.Errorf("Log() error = %v", err)
		}

		if len(commits) != 0 {
			t.Errorf("Log() returned %d commits, want 0", len(commits))
		}
	})
}

func TestHasUncommittedChanges(t *testing.T) {
	t.Run("returns false for clean repository", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		hasChanges, err := HasUncommittedChanges(tmpDir)
		if err != nil {
			t.Errorf("HasUncommittedChanges() error = %v", err)
		}

		if hasChanges {
			t.Error("HasUncommittedChanges() = true, want false")
		}
	})

	t.Run("returns true for modified files", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create a file
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("content"), 0644)

		hasChanges, err := HasUncommittedChanges(tmpDir)
		if err != nil {
			t.Errorf("HasUncommittedChanges() error = %v", err)
		}

		if !hasChanges {
			t.Error("HasUncommittedChanges() = false, want true")
		}
	})
}

func TestGetCurrentBranch(t *testing.T) {
	t.Run("returns current branch name", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		// Create initial commit so branch is created
		testFile := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(testFile, []byte("content"), 0644)
		Commit(tmpDir, "Initial commit", []string{"test.txt"})

		branch, err := GetCurrentBranch(tmpDir)
		if err != nil {
			t.Errorf("GetCurrentBranch() error = %v", err)
		}

		// Default branch is usually "main" or "master"
		if branch != "main" && branch != "master" {
			t.Errorf("GetCurrentBranch() = %v, want main or master", branch)
		}
	})
}

func TestSetAndGetConfig(t *testing.T) {
	t.Run("sets and gets config value", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		err := SetConfig(tmpDir, "user.name", "Test User", false)
		if err != nil {
			t.Errorf("SetConfig() error = %v", err)
		}

		value, err := GetConfig(tmpDir, "user.name")
		if err != nil {
			t.Errorf("GetConfig() error = %v", err)
		}

		if value != "Test User" {
			t.Errorf("GetConfig() = %v, want 'Test User'", value)
		}
	})

	t.Run("returns empty string for non-existent config", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		defer os.RemoveAll(tmpDir)

		value, err := GetConfig(tmpDir, "nonexistent.key")
		if err != nil {
			t.Errorf("GetConfig() error = %v", err)
		}

		if value != "" {
			t.Errorf("GetConfig() = %v, want empty string", value)
		}
	})
}

func TestEnsureGitInstalled(t *testing.T) {
	t.Run("checks git is installed", func(t *testing.T) {
		err := EnsureGitInstalled()
		if err != nil {
			t.Errorf("EnsureGitInstalled() error = %v (git should be installed for tests)", err)
		}
	})
}
