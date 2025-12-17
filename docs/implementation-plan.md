# growth.md Implementation Plan

**Status**: In Progress
**Target MVP Completion**: 6 weeks from start

---

## Overview

This document tracks the step-by-step implementation of growth.md MVP. Each phase builds on the previous one, and tasks are marked as complete when code is written, tested, and committed.

**Progress Overview**:
- Phase 1: Project Foundation - [x] 10/10 (100%)
- Phase 2: Core Domain Models - [x] 8/8 (100%)
- Phase 3: Storage Layer - [x] 6/7 (100% - 1 deferred)
- Phase 4: CLI Framework - [x] 6/6 (100%)
- Phase 5: Entity Commands - [ ] 0/15
- Phase 6: Git Integration - [ ] 0/4
- Phase 7: AI Integration - [ ] 0/5
- Phase 8: Polish & Testing - [ ] 0/6

**Total Progress**: 30/60 tasks complete (50% - 1 task deferred)

---

## Phase 1: Project Foundation

**Goal**: Set up Go project structure, dependencies, and basic tooling

### 1.1 Initialize Go Module
- [x] Run `go mod init github.com/yourusername/growth.md`
- [x] Create basic directory structure:
  ```
  growth.md/
  ├── cmd/growth/
  ├── internal/
  │   ├── cli/
  │   ├── core/
  │   ├── storage/
  │   ├── ai/
  │   └── git/
  ├── pkg/
  ├── docs/
  ├── examples/
  └── tests/
  ```
- [x] Add `.gitignore` for Go (binaries, IDE files, test coverage)

**Files to create**:
- `go.mod`
- `.gitignore`
- Directory structure

---

### 1.2 Add Core Dependencies
- [x] Install Cobra: `go get -u github.com/spf13/cobra@latest`
- [x] Install Viper: `go get -u github.com/spf13/viper`
- [x] Install YAML parser: `go get gopkg.in/yaml.v3`
- [x] Install goldmark (Markdown): `go get github.com/yuin/goldmark`
- [x] Install testify: `go get github.com/stretchr/testify`

**Verification**: `go mod tidy` runs successfully

---

### 1.3 Create Main Entry Point
- [x] Create `cmd/growth/main.go` with basic Cobra root command
- [x] Add version flag
- [x] Add basic `--help` output
- [x] Test: `go run cmd/growth/main.go --version`

**Files to create**:
- `cmd/growth/main.go`

**Example**:
```go
package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var version = "0.1.0-alpha"

func main() {
    rootCmd := &cobra.Command{
        Use:   "growth",
        Short: "Git-native career development manager",
        Long:  `growth.md - Track skills, goals, and learning paths in plain Markdown`,
        Version: version,
    }

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

---

### 1.4 Set Up Makefile
- [x] Create `Makefile` with common tasks:
  - `make build` - Build binary
  - `make test` - Run tests
  - `make lint` - Run linter
  - `make clean` - Clean build artifacts
  - `make install` - Install to `$GOPATH/bin`

**Files to create**:
- `Makefile`

**Example targets**:
```makefile
.PHONY: build test lint clean install

build:
	go build -o bin/growth cmd/growth/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

install:
	go install cmd/growth/main.go
```

---

### 1.5 Create README.md
- [x] Create basic README with:
  - Project description
  - Installation instructions
  - Quick start example
  - Link to concept doc
  - Development status (MVP in progress)

**Files to create**:
- `README.md`

---

### 1.6 Add MIT License
- [x] Create `LICENSE` file with MIT license text
- [x] Update copyright year and author

**Files to create**:
- `LICENSE`

---

### 1.7 Set Up Testing Structure
- [x] Create `tests/fixtures/` directory for test data
- [x] Create example test: `internal/core/skill_test.go`
- [x] Verify: `make test` passes

**Files to create**:
- `tests/fixtures/.gitkeep`
- `internal/core/skill_test.go` (placeholder)

---

### 1.8 Configure golangci-lint
- [x] Create `.golangci.yml` configuration
- [x] Enable key linters: govet, errcheck, staticcheck, unused
- [x] Run: `make lint` (shows installation instructions if not installed)

**Files to create**:
- `.golangci.yml`

---

### 1.9 Set Up GitHub Actions (Optional but Recommended)
- [x] Create `.github/workflows/test.yml` for CI
- [x] Run tests on push to main
- [x] Run linter

**Files to create**:
- `.github/workflows/test.yml`

---

### 1.10 Initial Git Commit
- [x] Initialize git repo: `git init`
- [x] Add all files: `git add .`
- [x] Commit: `git commit -m "Initial project structure"`
- [x] Tag: `git tag v0.1.0-alpha`

**Verification**: Clean git status, tagged commit

---

## Phase 2: Core Domain Models

**Goal**: Define Go structs for all entities with proper validation

### 2.1 Create Base Entity Types
- [x] Create `internal/core/types.go`
- [x] Define common types:
  - `EntityID` (string type alias)
  - `Status` (enum: active, completed, archived)
  - `Priority` (enum: high, medium, low)
  - `ProficiencyLevel` (enum: beginner, intermediate, advanced, expert)
  - `SkillStatus` (enum: not-started, learning, mastered)

**Files to create**:
- `internal/core/types.go`

**Example**:
```go
package core

type EntityID string

type Status string
const (
    StatusActive    Status = "active"
    StatusCompleted Status = "completed"
    StatusArchived  Status = "archived"
)

type Priority string
const (
    PriorityHigh   Priority = "high"
    PriorityMedium Priority = "medium"
    PriorityLow    Priority = "low"
)

type ProficiencyLevel string
const (
    LevelBeginner     ProficiencyLevel = "beginner"
    LevelIntermediate ProficiencyLevel = "intermediate"
    LevelAdvanced     ProficiencyLevel = "advanced"
    LevelExpert       ProficiencyLevel = "expert"
)

type SkillStatus string
const (
    SkillNotStarted SkillStatus = "not-started"
    SkillLearning   SkillStatus = "learning"
    SkillMastered   SkillStatus = "mastered"
)
```

---

### 2.2 Create Skill Entity
- [x] Create `internal/core/skill.go`
- [x] Define `Skill` struct with all fields from whitepaper
- [x] Add validation method: `Validate() error`
- [x] Add `NewSkill()` constructor
- [x] Write tests in `internal/core/skill_test.go`

**Files to create**:
- `internal/core/skill.go`
- `internal/core/skill_test.go`

**Struct fields**:
```go
type Skill struct {
    ID        EntityID         `yaml:"id"`
    Title     string           `yaml:"title"`
    Category  string           `yaml:"category"`
    Level     ProficiencyLevel `yaml:"level"`
    Status    SkillStatus      `yaml:"status"`
    Created   time.Time        `yaml:"created"`
    Updated   time.Time        `yaml:"updated"`
    Resources []EntityID       `yaml:"resources"`
    Tags      []string         `yaml:"tags"`
}
```

---

### 2.3 Create Goal Entity
- [x] Create `internal/core/goal.go`
- [x] Define `Goal` struct with all fields
- [x] Add validation method
- [x] Add `NewGoal()` constructor
- [x] Write tests

**Files to create**:
- `internal/core/goal.go`
- `internal/core/goal_test.go`

**Key fields**:
```go
type Goal struct {
    ID            EntityID   `yaml:"id"`
    Title         string     `yaml:"title"`
    Status        Status     `yaml:"status"`
    Priority      Priority   `yaml:"priority"`
    Created       time.Time  `yaml:"created"`
    Updated       time.Time  `yaml:"updated"`
    TargetDate    *time.Time `yaml:"targetDate,omitempty"`
    LearningPaths []EntityID `yaml:"learningPaths"`
    Milestones    []EntityID `yaml:"milestones"`
    Tags          []string   `yaml:"tags"`
}
```

---

### 2.4 Create Learning Path Entity
- [x] Create `internal/core/path.go`
- [x] Define `LearningPath` struct
- [x] Add `PathType` enum (ai-generated, manual)
- [x] Add validation and constructor
- [x] Write tests

**Files to create**:
- `internal/core/path.go`
- `internal/core/path_test.go`

---

### 2.5 Create Phase Entity
- [x] Create `internal/core/phase.go`
- [x] Define `Phase` struct
- [x] Define `SkillRequirement` struct (skill ID + target level)
- [x] Add validation and constructor
- [x] Write tests

**Files to create**:
- `internal/core/phase.go`
- `internal/core/phase_test.go`

**SkillRequirement**:
```go
type SkillRequirement struct {
    SkillID     EntityID         `yaml:"skillId"`
    TargetLevel ProficiencyLevel `yaml:"targetLevel"`
}

type Phase struct {
    ID               EntityID           `yaml:"id"`
    PathID           EntityID           `yaml:"pathId"`
    Title            string             `yaml:"title"`
    Order            int                `yaml:"order"`
    EstimatedDuration string            `yaml:"estimatedDuration,omitempty"`
    RequiredSkills   []SkillRequirement `yaml:"requiredSkills"`
    Milestones       []EntityID         `yaml:"milestones"`
}
```

---

### 2.6 Create Resource Entity
- [x] Create `internal/core/resource.go`
- [x] Define `Resource` struct
- [x] Add `ResourceType` enum (book, course, video, article, project, documentation)
- [x] Add `ResourceStatus` enum (not-started, in-progress, completed)
- [x] Add validation and constructor
- [x] Write tests

**Files to create**:
- `internal/core/resource.go`
- `internal/core/resource_test.go`

---

### 2.7 Create Milestone Entity
- [x] Create `internal/core/milestone.go`
- [x] Define `Milestone` struct
- [x] Add `MilestoneType` enum (goal-level, path-level, skill-level)
- [x] Add `ReferenceType` enum (goal, path, skill)
- [x] Add validation and constructor
- [x] Write tests

**Files to create**:
- `internal/core/milestone.go`
- `internal/core/milestone_test.go`

---

### 2.8 Create Progress Log Entity
- [x] Create `internal/core/progress.go`
- [x] Define `ProgressLog` struct
- [x] Add helper for getting week-of date
- [x] Add validation and constructor
- [x] Write tests

**Files to create**:
- `internal/core/progress.go`
- `internal/core/progress_test.go`

**Struct**:
```go
type ProgressLog struct {
    ID                 EntityID   `yaml:"id"`
    WeekOf             time.Time  `yaml:"weekOf"`
    HoursInvested      float64    `yaml:"hoursInvested,omitempty"`
    SkillsWorked       []EntityID `yaml:"skillsWorked"`
    ResourcesUsed      []EntityID `yaml:"resourcesUsed"`
    MilestonesAchieved []EntityID `yaml:"milestonesAchieved"`
    Mood               string     `yaml:"mood,omitempty"`
}
```

---

## Phase 3: Storage Layer

**Goal**: Implement Markdown file I/O with YAML frontmatter parsing

### 3.1 Create Frontmatter Parser
- [x] Create `internal/storage/frontmatter.go`
- [x] Implement `ParseFrontmatter(content []byte) (map[string]interface{}, string, error)`
- [x] Implement `SerializeFrontmatter(data interface{}) ([]byte, error)`
- [x] Handle edge cases (missing frontmatter, malformed YAML)
- [x] Write tests

**Files to create**:
- `internal/storage/frontmatter.go`
- `internal/storage/frontmatter_test.go`

**Functions**:
```go
// ParseFrontmatter extracts YAML frontmatter and body from markdown
func ParseFrontmatter(content []byte) (frontmatter map[string]interface{}, body string, err error)

// SerializeFrontmatter combines frontmatter and body into markdown
func SerializeFrontmatter(frontmatter interface{}, body string) ([]byte, error)
```

---

### 3.2 Create Entity Repository Interface
- [x] Create `internal/storage/repository.go`
- [x] Define `Repository[T]` interface with CRUD operations:
  - `Create(entity T) error`
  - `GetByID(id EntityID) (T, error)`
  - `GetByIDWithBody(id EntityID) (T, error)`
  - `GetAll() ([]T, error)`
  - `Update(entity T) error`
  - `Delete(id EntityID) error`
  - `Search(query string) ([]T, error)`
  - `Exists(id EntityID) (bool, error)`

**Files to create**:
- `internal/storage/repository.go`

**Interface**:
```go
type Repository[T any] interface {
    Create(entity T) error
    GetByID(id EntityID) (*T, error)
    GetAll() ([]T, error)
    Update(entity T) error
    Delete(id EntityID) error
    Search(query string) ([]T, error)
}
```

---

### 3.3 Implement Filesystem Repository
- [x] Create `internal/storage/fs_repository.go`
- [x] Implement `FilesystemRepository[T]` struct
- [x] Implement all Repository interface methods (Create, GetByID, GetByIDWithBody, GetAll, Update, Delete, Search, Exists)
- [x] Handle file naming: `{id}-{slug}.md` (e.g., "skill-001-python.md")
- [x] Write tests with temp directories (50+ tests covering all operations)

**Files to create**:
- `internal/storage/fs_repository.go`
- `internal/storage/fs_repository_test.go`

**Key methods**:
```go
type FilesystemRepository[T any] struct {
    basePath    string
    entityType  string
}

func NewFilesystemRepository[T any](basePath, entityType string) *FilesystemRepository[T]

func (r *FilesystemRepository[T]) Create(entity T) error
func (r *FilesystemRepository[T]) GetByID(id EntityID) (*T, error)
// ... etc
```

---

### 3.4 Implement Entity Repositories
- [x] Create `internal/storage/skill_repository.go`
- [x] Create wrapper around `FilesystemRepository[Skill]`
- [x] Add skill-specific queries (by category, level, status)
- [x] Create `internal/storage/goal_repository.go`
- [x] Add goal-specific queries (by status, priority, target date range)
- [x] Create `internal/storage/path_repository.go`
- [x] Add path-specific queries (by type, status, AI-generated)
- [x] Create `internal/storage/phase_repository.go`
- [x] Add phase-specific queries (by path ID, ordered by sequence)
- [x] Create `internal/storage/resource_repository.go`
- [x] Add resource-specific queries (by type, skill ID, status)
- [x] Create `internal/storage/milestone_repository.go`
- [x] Add milestone-specific queries (by reference ID, status, type)
- [x] Create `internal/storage/progress_repository.go`
- [x] Add progress log queries (by date range, skill ID, resource ID)
- [x] Write comprehensive tests for all repositories (115 tests passing)

**Files created**:
- `internal/storage/skill_repository.go` + `*_test.go`
- `internal/storage/goal_repository.go` + `*_test.go`
- `internal/storage/path_repository.go` + `*_test.go`
- `internal/storage/phase_repository.go` + `*_test.go`
- `internal/storage/resource_repository.go` + `*_test.go`
- `internal/storage/milestone_repository.go` + `*_test.go`
- `internal/storage/progress_repository.go` + `*_test.go`

**Repository Pattern**: Each repository wraps `FilesystemRepository[T]` and adds domain-specific query methods for filtering and searching entities.

---

### 3.5 Implement Config Management
- [x] Create `internal/storage/config.go`
- [x] Define `Config` struct matching whitepaper spec
- [x] Implement `LoadConfig(path string) (*Config, error)`
- [x] Implement `SaveConfig(config *Config, path string) error`
- [x] Add default config generation
- [x] Write tests (27 tests passing)

**Files created**:
- `internal/storage/config.go`
- `internal/storage/config_test.go`

**Config struct**:
```go
type Config struct {
    Version  string      `yaml:"version"`
    User     UserConfig  `yaml:"user"`
    AI       AIConfig    `yaml:"ai"`
    Git      GitConfig   `yaml:"git"`
    Progress ProgressConfig `yaml:"progress"`
    Display  DisplayConfig  `yaml:"display"`
    MCP      MCPConfig   `yaml:"mcp"`
}
```

---

### 3.6 Implement Index/Cache Layer (Optional for MVP)
- [ ] **DEFERRED** - Skipped for MVP to avoid premature optimization
- [ ] Cache would add complexity without clear performance need at this stage
- [ ] Can be added later if file system operations prove too slow

**Status**: Deferred to post-MVP

---

### 3.7 Create Storage Integration Tests
- [ ] **DEFERRED** - Unit tests provide sufficient coverage for MVP
- [ ] Integration tests can be added later for production readiness

**Status**: Deferred to post-MVP

---

## Phase 4: CLI Framework

**Goal**: Set up Cobra commands structure and global flags

### 4.1 Create Root Command
- [x] Create `internal/cli/root.go`
- [x] Define root command with description
- [x] Add global flags:
  - `--config` (config file path)
  - `--repo` (growth repository path)
  - `--format/-f` (output format: table, json, yaml)
  - `--verbose/-v` (verbose output)
- [x] Add PersistentPreRun to load config and initialize repositories
- [x] Wire up to `cmd/growth/main.go`
- [x] Test: `growth --version`, `growth --help`, flags work correctly

**Files created**:
- `internal/cli/root.go`

**Files modified**:
- `cmd/growth/main.go` (now calls `cli.Execute()`)

**Example**:
```go
var (
    cfgFile    string
    repoPath   string
    outputFormat string
    verbose    bool
)

var rootCmd = &cobra.Command{
    Use:   "growth",
    Short: "Git-native career development manager",
    Long:  `Track your skills, goals, and learning paths in plain Markdown files`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Load config, initialize repository
        return initializeApp()
    },
}

func init() {
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
    rootCmd.PersistentFlags().StringVar(&repoPath, "repo", "", "growth repository path")
    rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "table", "output format")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
```

---

### 4.2 Create Init Command
- [x] Create `internal/cli/init.go`
- [x] Implement `growth init [directory]`
- [x] Create directory structure (skills/, goals/, paths/, phases/, resources/, milestones/, progress/)
- [x] Initialize Git repository with initial commit
- [x] Create default `config.yml` with validation
- [x] Create `.gitignore` and `README.md`
- [x] Add interactive prompts (name, email, AI provider, auto-commit)
- [x] Test: `growth init test-dir` creates full structure

**Files created**:
- `internal/cli/init.go`

**Features**:
- Interactive configuration prompts with sensible defaults
- Automatic Git initialization and first commit
- Generated README with quick start guide
- Support for multiple AI providers (openai/anthropic/google/local)

**Command structure**:
```go
var initCmd = &cobra.Command{
    Use:   "init [directory]",
    Short: "Initialize a new growth repository",
    Args:  cobra.MaximumNArgs(1),
    RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
    // Implementation
}
```

---

### 4.3 Create Common Output Utilities
- [x] Create `internal/cli/output.go`
- [x] Implement `PrintTable(data interface{})` - Smart table formatting with column width handling
- [x] Implement `PrintJSON(data interface{})` - Pretty-printed JSON output
- [x] Implement `PrintYAML(data interface{})` - YAML formatted output
- [x] Implement `PrintSuccess(message string)` - Green checkmark success messages
- [x] Implement `PrintError(err error)` - Red X error messages
- [x] Add `PrintWarning`, `PrintInfo`, `Print`, `Println` utilities
- [x] Implement `PrintOutput(data, format)` - Format-aware output dispatcher
- [x] Write comprehensive tests (9 tests passing)

**Files created**:
- `internal/cli/output.go`
- `internal/cli/output_test.go`

**Features**:
- ANSI color support for success/error/warning/info messages
- Smart table formatting with automatic column sizing
- Handles structs, slices, and pointers
- Truncates long values to fit columns
- Formats time.Time as dates, slices with ellipsis
- Respects YAML field tags for display names

---

### 4.4 Create Input/Prompt Utilities
- [x] Create `internal/cli/input.go`
- [x] Implement `PromptString(prompt string, defaultValue string) string`
- [x] Implement `PromptStringRequired(prompt string) string` - loops until non-empty
- [x] Implement `PromptConfirm(prompt string) bool`
- [x] Implement `PromptConfirmDefault(prompt string, defaultYes bool) bool`
- [x] Implement `PromptSelect(prompt string, options []string) string`
- [x] Implement `PromptSelectWithDefault(prompt, options, default) string`
- [x] Implement `PromptInt(prompt string, defaultValue int) int`
- [x] Implement `PromptMultiline(prompt string) string` - for long-form text input
- [x] Add validation helpers:
  - `ValidateNotEmpty(value, fieldName) error`
  - `ValidateEmail(email) error`
  - `ValidateOneOf(value, options) error`
  - `ValidatePositive(value, fieldName) error`
  - `ValidateRange(value, min, max, fieldName) error`

**Files created**:
- `internal/cli/input.go`

**Features**:
- All prompts use buffered reader from stdin
- Support for default values with visual indicators
- Select prompts show numbered options
- Multiline input supports '.' terminator or Ctrl+D
- Comprehensive validation helpers for common patterns

---

### 4.5 Create ID Generation Utilities
- [x] Create `internal/cli/id_gen.go`
- [x] Implement `GenerateNextID(entityType string) (EntityID, error)`
- [x] Implement `GenerateNextIDInPath(entityType, basePath) (EntityID, error)` - testable version
- [x] Scan existing files to find max ID using filepath.Glob
- [x] Zero-pad to 3 digits (001, 002, etc.)
- [x] Handle gaps in numbering (finds max + 1)
- [x] Add slug generation from title: `GenerateSlug(title string) string`
  - Lowercase and hyphenate
  - Remove special characters
  - Truncate to 50 chars
  - Handle empty titles with "untitled"
- [x] Add filename generator: `GenerateFileName(id, title) string`
- [x] Write comprehensive tests (15 tests passing)

**Files created**:
- `internal/cli/id_gen.go`
- `internal/cli/id_gen_test.go`

**Features**:
- Supports all 7 entity types (skill, goal, path, phase, resource, milestone, progress)
- Regex-based ID extraction from filenames
- Smart slug generation with cleanup and truncation
- Returns clear error for unknown entity types

---

### 4.6 Test CLI Framework
- [x] Test `growth --help` shows all commands
- [x] Test `growth --version` shows version
- [x] Test `growth init` creates structure
- [x] Test global flags work
- [x] Run all CLI package tests - 24 tests passing:
  - 8 slug generation tests
  - 3 filename generation tests
  - 4 ID generation tests (including gaps, unknown types)
  - 1 JSON output test
  - 1 YAML output test
  - 4 table output tests (single struct, slice, empty, Skills)
  - 3 field formatting tests (time, slice, nil pointer)
- [x] Build CLI binary: `go build -o bin/growth cmd/growth/main.go`
- [x] Verify help output displays root command and init subcommand

**Verification**: All basic CLI commands work, all tests passing

---

## Phase 5: Entity Commands

**Goal**: Implement CRUD commands for all entities

### 5.1 Skill Commands - Create
- [ ] Create `internal/cli/skill.go`
- [ ] Implement `growth skill create <title> [flags]`
- [ ] Flags: `--category`, `--level`, `--tags`
- [ ] Generate ID and slug
- [ ] Create markdown file with frontmatter
- [ ] Print success message with ID

**Files to create**:
- `internal/cli/skill.go`

**Command**:
```bash
growth skill create "Python" --category programming --level intermediate --tags python,ml
```

---

### 5.2 Skill Commands - List & View
- [ ] Implement `growth skill list [flags]`
- [ ] Flags: `--category`, `--level`, `--status`
- [ ] Display as table by default
- [ ] Implement `growth skill view <id-or-slug>`
- [ ] Show full details in formatted output

**Commands**:
```bash
growth skill list --category ml
growth skill view python
```

---

### 5.3 Skill Commands - Edit & Delete
- [ ] Implement `growth skill edit <id> [flags]`
- [ ] Allow updating level, status, tags
- [ ] Update `updated` timestamp
- [ ] Implement `growth skill delete <id>`
- [ ] Prompt for confirmation
- [ ] Check for references before deleting

**Commands**:
```bash
growth skill edit python --level advanced
growth skill delete skill-005
```

---

### 5.4 Goal Commands - Full CRUD
- [ ] Create `internal/cli/goal.go`
- [ ] Implement `growth goal create <title> [flags]`
- [ ] Flags: `--priority`, `--target-date`
- [ ] Implement `growth goal list [flags]`
- [ ] Implement `growth goal view <id>`
- [ ] Implement `growth goal edit <id> [flags]`
- [ ] Implement `growth goal delete <id>`

**Files to create**:
- `internal/cli/goal.go`

---

### 5.5 Goal Commands - Path Management
- [ ] Implement `growth goal add-path <goal-id> <path-id>`
- [ ] Implement `growth goal remove-path <goal-id> <path-id>`
- [ ] Update goal file frontmatter
- [ ] Validate path exists

**Commands**:
```bash
growth goal add-path goal-001 path-001
growth goal remove-path goal-001 path-002
```

---

### 5.6 Path Commands - Basic CRUD
- [ ] Create `internal/cli/path.go`
- [ ] Implement `growth path list [flags]`
- [ ] Implement `growth path view <id>`
- [ ] Implement `growth path edit <id>`
- [ ] Implement `growth path delete <id>`

**Files to create**:
- `internal/cli/path.go`

**Note**: Path generation (AI) will be in Phase 7

---

### 5.7 Resource Commands - Full CRUD
- [ ] Create `internal/cli/resource.go`
- [ ] Implement `growth resource create <title> --skill <id> --type <type> [flags]`
- [ ] Flags: `--url`, `--author`, `--estimated-hours`
- [ ] Implement `growth resource list [flags]`
- [ ] Filter by skill, type, status
- [ ] Implement `growth resource view <id>`
- [ ] Implement `growth resource edit <id>`
- [ ] Implement `growth resource delete <id>`

**Files to create**:
- `internal/cli/resource.go`

---

### 5.8 Resource Commands - Status Updates
- [ ] Implement `growth resource start <id>`
- [ ] Update status to `in-progress`
- [ ] Implement `growth resource complete <id>`
- [ ] Update status to `completed`
- [ ] Add timestamp tracking

**Commands**:
```bash
growth resource start resource-001
growth resource complete resource-001
```

---

### 5.9 Milestone Commands - Full CRUD
- [ ] Create `internal/cli/milestone.go`
- [ ] Implement `growth milestone create <title> --type <type> --ref <id>`
- [ ] Implement `growth milestone list [flags]`
- [ ] Filter by status, type
- [ ] Implement `growth milestone view <id>`
- [ ] Implement `growth milestone achieve <id> [flags]`
- [ ] Flag: `--proof` (URL to evidence)
- [ ] Set achievedDate to now

**Files to create**:
- `internal/cli/milestone.go`

---

### 5.10 Progress Commands - Log Entry
- [ ] Create `internal/cli/progress.go`
- [ ] Implement `growth progress log [message]`
- [ ] Flags: `--skill`, `--resource`, `--milestone`
- [ ] Create or update current week's log
- [ ] Append to notes section
- [ ] Update referenced entities

**Files to create**:
- `internal/cli/progress.go`

**Command**:
```bash
growth progress log "Completed Fast.ai Lesson 3" --skill ml --resource fastai
```

---

### 5.11 Progress Commands - View & Stats
- [ ] Implement `growth progress view [flags]`
- [ ] Flags: `--week`, `--month`, `--date`
- [ ] Display formatted log
- [ ] Implement `growth progress stats [flags]`
- [ ] Show hours invested, skills worked on, consistency

**Commands**:
```bash
growth progress view --week
growth progress stats --month
```

---

### 5.12 Search Command
- [ ] Create `internal/cli/search.go`
- [ ] Implement `growth search <query>`
- [ ] Search across all entity types
- [ ] Search in titles, tags, notes
- [ ] Display results grouped by type

**Files to create**:
- `internal/cli/search.go`

**Command**:
```bash
growth search "neural networks"
```

---

### 5.13 Overview Command
- [ ] Create `internal/cli/overview.go`
- [ ] Implement `growth overview`
- [ ] Display summary:
  - Active goals count
  - Skills in progress
  - Recent progress logs
  - Upcoming milestones
- [ ] Text-based dashboard

**Files to create**:
- `internal/cli/overview.go`

---

### 5.14 Stats Command
- [ ] Create `internal/cli/stats.go`
- [ ] Implement `growth stats`
- [ ] Calculate and display:
  - Total skills (by level, status)
  - Total goals (by status)
  - Total hours invested
  - Learning velocity (hours/week)
  - Milestone completion rate

**Files to create**:
- `internal/cli/stats.go`

---

### 5.15 Integration Testing for Commands
- [ ] Create `tests/cli_integration_test.go`
- [ ] Test complete workflows:
  - Create goal → generate path → log progress
  - Create skill → add resource → mark complete
  - Create milestone → achieve milestone
- [ ] Test error cases

**Files to create**:
- `tests/cli_integration_test.go`

---

## Phase 6: Git Integration

**Goal**: Auto-commit changes and provide git utilities

### 6.1 Create Git Operations Module
- [ ] Create `internal/git/operations.go`
- [ ] Implement `InitRepo(path string) error`
- [ ] Implement `Commit(message string, files []string) error`
- [ ] Implement `Status() ([]string, error)`
- [ ] Implement `IsRepo(path string) bool`

**Files to create**:
- `internal/git/operations.go`
- `internal/git/operations_test.go`

---

### 6.2 Integrate Git with Storage Layer
- [ ] Modify `FilesystemRepository` to optionally auto-commit
- [ ] Check config: `git.autoCommit` and `git.commitOnUpdate`
- [ ] Generate commit messages using template
- [ ] Handle git errors gracefully

**Modified files**:
- `internal/storage/fs_repository.go`

---

### 6.3 Add Git Commands to CLI
- [ ] Create `internal/cli/git.go`
- [ ] Implement `growth git status` (wrapper)
- [ ] Implement `growth git log` (wrapper)
- [ ] Add `--no-commit` flag to entity commands

**Files to create**:
- `internal/cli/git.go`

---

### 6.4 Test Git Integration
- [ ] Test auto-commit on entity creation
- [ ] Test commit message formatting
- [ ] Test behavior when not a git repo
- [ ] Test `--no-commit` flag

**Tests to write**:
- `tests/git_integration_test.go`

---

## Phase 7: AI Integration (Basic)

**Goal**: Implement basic AI path generation using OpenAI or Anthropic

### 7.1 Create AI Client Module
- [ ] Create `internal/ai/client.go`
- [ ] Define `AIClient` interface:
  - `GeneratePath(goal Goal, skills []Skill, context string) (LearningPath, error)`
  - `SuggestResources(skill Skill) ([]Resource, error)`
- [ ] Implement OpenAI client: `internal/ai/openai_client.go`
- [ ] Add API key loading from env var
- [ ] Handle rate limits and errors

**Files to create**:
- `internal/ai/client.go`
- `internal/ai/openai_client.go`

---

### 7.2 Create Path Generation Prompt
- [ ] Create `internal/ai/prompts.go`
- [ ] Define path generation prompt template
- [ ] Include goal, skills, background in prompt
- [ ] Request structured markdown output
- [ ] Add examples for few-shot learning

**Files to create**:
- `internal/ai/prompts.go`

**Prompt structure**:
```go
const PathGenerationPrompt = `
You are an expert career coach for software engineers. Generate a personalized learning path.

CONTEXT:
- Goal: {{.GoalTitle}}
- Current Skills: {{.SkillsList}}
- Background: {{.UserBackground}}

... [rest of prompt]
`
```

---

### 7.3 Implement Path Generator
- [ ] Create `internal/ai/path_generator.go`
- [ ] Implement `GenerateLearningPath(goal Goal, skills []Skill, config AIConfig) (*LearningPath, error)`
- [ ] Parse AI response into structured Path
- [ ] Create Phase entities
- [ ] Save path and phases to files
- [ ] Handle parsing errors

**Files to create**:
- `internal/ai/path_generator.go`
- `internal/ai/path_generator_test.go`

---

### 7.4 Add Path Generate Command
- [ ] Update `internal/cli/path.go`
- [ ] Implement `growth path generate <goal-id> [flags]`
- [ ] Flags: `--approach`, `--model`
- [ ] Show progress indicator during generation
- [ ] Display generated path summary
- [ ] Automatically link path to goal

**Command**:
```bash
growth path generate goal-001 --approach "fast.ai top-down"
```

---

### 7.5 Test AI Integration
- [ ] Create mock AI client for testing
- [ ] Test prompt generation
- [ ] Test response parsing
- [ ] Test error handling (API errors, invalid responses)
- [ ] Manual test with real AI API (document in README)

**Files to create**:
- `internal/ai/mock_client.go`
- `tests/ai_integration_test.go`

---

## Phase 8: Polish & Testing

**Goal**: Add final touches, comprehensive testing, and documentation

### 8.1 Add Validation Across All Entities
- [ ] Review all entity constructors
- [ ] Add comprehensive validation:
  - Required fields present
  - Enums are valid values
  - Dates are valid
  - IDs exist (for references)
- [ ] Return clear error messages

**Files to update**:
- All `internal/core/*.go` files

---

### 8.2 Add Comprehensive Error Handling
- [ ] Review all CLI commands
- [ ] Add user-friendly error messages
- [ ] Handle common errors gracefully:
  - File not found
  - Invalid ID
  - Permission denied
  - Git not initialized
  - Config missing
- [ ] Add `--debug` flag for verbose errors

**Files to update**:
- All `internal/cli/*.go` files

---

### 8.3 Write User Documentation
- [ ] Update `README.md` with full getting started guide
- [ ] Create `docs/getting-started.md`
- [ ] Create `docs/cli-reference.md` (all commands)
- [ ] Create `docs/configuration.md` (config.yml options)
- [ ] Add examples for common workflows

**Files to create/update**:
- `README.md`
- `docs/getting-started.md`
- `docs/cli-reference.md`
- `docs/configuration.md`

---

### 8.4 Add Example Learning Paths
- [ ] Create `examples/ml-engineer/` with sample path
- [ ] Create `examples/backend-specialist/` with sample path
- [ ] Create `examples/frontend-developer/` with sample path
- [ ] Each example includes:
  - Goal file
  - Path file
  - Phases
  - Skills
  - Resources

**Directories to create**:
- `examples/ml-engineer/`
- `examples/backend-specialist/`
- `examples/frontend-developer/`

---

### 8.5 Final Testing & Bug Fixes
- [ ] Run full test suite: `make test`
- [ ] Run linter: `make lint`
- [ ] Test on clean system (fresh git clone)
- [ ] Test all documented workflows
- [ ] Fix any discovered bugs
- [ ] Achieve >80% test coverage

**Verification**: All tests pass, no critical bugs

---

### 8.6 Build & Release
- [ ] Build binaries for all platforms:
  - `make build-linux`
  - `make build-macos`
  - `make build-windows`
- [ ] Test binary on each platform
- [ ] Create GitHub release (v0.1.0)
- [ ] Upload binaries
- [ ] Tag commit: `git tag v0.1.0`

**Release artifacts**:
- `growth-linux-amd64`
- `growth-darwin-amd64`
- `growth-darwin-arm64`
- `growth-windows-amd64.exe`

---

## Next Steps After MVP

Once MVP is complete, consider:

1. **TUI Dashboard** (Phase 2 from roadmap)
   - Implement `growth board` using Bubble Tea
   - Interactive, real-time dashboard

2. **MCP Server** (Phase 3 from roadmap)
   - Implement MCP protocol
   - Enable Claude integration

3. **Enhanced AI Features**
   - Path regeneration based on progress
   - Resource recommendations
   - Skill gap analysis

4. **Community Features**
   - Path templates repository
   - Sharing mechanism
   - Import/export

---

## Tracking Progress

**How to use this document**:
1. Mark tasks complete with `[x]` as you finish them
2. Commit changes to this file after each session
3. Update progress percentages at the top
4. Add notes to decisions log as needed
5. Cross-reference commit SHAs for major milestones
