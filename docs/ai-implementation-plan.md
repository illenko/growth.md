# AI Integration Implementation Plan

**Created**: 2025-12-23
**Status**: In Progress - Phase 1 & 2 Complete ✅, Ready for Testing
**Goal**: Add AI-powered learning path generation and MCP server integration

## Overview

This plan focuses on making growth.md AI-native by:
1. **AI Integration**: Generate personalized learning paths using Google Gemini (free API)
2. **MCP Server**: Expose growth.md data to Claude Desktop and other AI assistants
3. **Provider-Agnostic Design**: Support multiple AI providers (Gemini, OpenAI, Anthropic, etc.)

## Why This Architecture?

- **Gemini First**: Free API with excellent quality (Google's latest models)
- **MCP Integration**: Make growth.md part of your AI workflow
- **Future-Proof**: Easy to add OpenAI, Anthropic, or local models later
- **Real Value**: AI-powered insights vs. just another CLI tool

---

## Phase 1: AI Client Architecture

**Goal**: Build provider-agnostic AI client with Gemini as the primary implementation

### 1.1 Create AI Client Interface

Create a clean abstraction that works with any AI provider.

**Files to create**:
- `internal/ai/client.go` - Core interfaces
- `internal/ai/types.go` - Shared types and responses

**Interface Design**:
```go
package ai

// AIClient is the main interface for AI providers
type AIClient interface {
    // GenerateLearningPath creates a personalized learning path from a goal
    GenerateLearningPath(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error)

    // SuggestResources recommends learning resources for a skill
    SuggestResources(ctx context.Context, req ResourceSuggestionRequest) (*ResourceSuggestionResponse, error)

    // AnalyzeProgress provides insights on progress and next steps
    AnalyzeProgress(ctx context.Context, req ProgressAnalysisRequest) (*ProgressAnalysisResponse, error)

    // Provider returns the name of the AI provider
    Provider() string
}

// PathGenerationRequest contains context for path generation
type PathGenerationRequest struct {
    Goal            *core.Goal
    CurrentSkills   []*core.Skill
    Background      string
    LearningStyle   string // e.g., "top-down", "bottom-up", "project-based"
    TimeCommitment  string // e.g., "10 hours/week"
    TargetDate      *time.Time
}

// PathGenerationResponse contains the generated path structure
type PathGenerationResponse struct {
    Path        *core.LearningPath
    Phases      []*core.Phase
    Resources   []*core.Resource
    Milestones  []*core.Milestone
    Reasoning   string
}

// ResourceSuggestionRequest for resource recommendations
type ResourceSuggestionRequest struct {
    Skill           *core.Skill
    CurrentLevel    core.ProficiencyLevel
    TargetLevel     core.ProficiencyLevel
    LearningStyle   string
    Budget          string // e.g., "free", "paid", "any"
}

// ResourceSuggestionResponse contains recommended resources
type ResourceSuggestionResponse struct {
    Resources   []*core.Resource
    Reasoning   string
}

// ProgressAnalysisRequest for analyzing user progress
type ProgressAnalysisRequest struct {
    Goal            *core.Goal
    Path            *core.LearningPath
    ProgressLogs    []*core.ProgressLog
    CurrentSkills   []*core.Skill
}

// ProgressAnalysisResponse contains insights and recommendations
type ProgressAnalysisResponse struct {
    Summary         string
    Insights        []string
    Recommendations []string
    IsOnTrack       bool
    SuggestedFocus  []string
}
```

**Config Structure**:
```go
// Config holds AI provider configuration
type Config struct {
    Provider    string // "gemini", "openai", "anthropic", "local"
    APIKey      string
    Model       string
    Temperature float32
    MaxTokens   int
    BaseURL     string // For custom endpoints
}
```

**Factory Pattern**:
```go
// NewClient creates an AI client based on config
func NewClient(cfg Config) (AIClient, error) {
    switch cfg.Provider {
    case "gemini":
        return NewGeminiClient(cfg)
    case "openai":
        return NewOpenAIClient(cfg)
    case "anthropic":
        return NewAnthropicClient(cfg)
    case "local":
        return NewLocalClient(cfg)
    default:
        return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
    }
}
```

**Tasks**:
- [ ] Create `internal/ai/client.go` with interfaces
- [ ] Create `internal/ai/types.go` with request/response types
- [ ] Create `internal/ai/config.go` with configuration
- [ ] Create `internal/ai/factory.go` with provider factory
- [ ] Add error types: `ErrAPIKeyMissing`, `ErrRateLimitExceeded`, `ErrInvalidResponse`

---

### 1.2 Implement Gemini Client

**Files to create**:
- `internal/ai/gemini/client.go` - Main client
- `internal/ai/gemini/prompts.go` - Prompt templates
- `internal/ai/gemini/parser.go` - Response parsing

**Gemini Client Structure**:
```go
package gemini

import (
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

type Client struct {
    client      *genai.Client
    model       *genai.GenerativeModel
    config      ai.Config
}

func NewClient(cfg ai.Config) (*Client, error) {
    ctx := context.Background()

    // Initialize Gemini client
    client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
    if err != nil {
        return nil, fmt.Errorf("failed to create Gemini client: %w", err)
    }

    // Select model (default: gemini-3-flash-preview)
    modelName := cfg.Model
    if modelName == "" {
        modelName = "gemini-3-flash-preview"
    }

    model := client.GenerativeModel(modelName)
    model.Temperature = &cfg.Temperature
    model.SetMaxOutputTokens(int32(cfg.MaxTokens))

    // Set structured output mode
    model.ResponseMIMEType = "application/json"

    return &Client{
        client: client,
        model:  model,
        config: cfg,
    }, nil
}

func (c *Client) Provider() string {
    return "gemini"
}
```

**Prompt Templates** (`internal/ai/gemini/prompts.go`):
```go
const PathGenerationPrompt = `You are an expert career coach for software engineers. Generate a personalized learning path.

GOAL: {{.Goal.Title}}
GOAL DESCRIPTION: {{.Goal.Body}}
PRIORITY: {{.Goal.Priority}}
{{if .Goal.TargetDate}}TARGET DATE: {{.Goal.TargetDate}}{{end}}

CURRENT SKILLS:
{{range .CurrentSkills}}
- {{.Title}} ({{.Level}}) - {{.Category}}
{{end}}

BACKGROUND:
{{.Background}}

LEARNING PREFERENCES:
- Learning Style: {{.LearningStyle}}
- Time Commitment: {{.TimeCommitment}}

TASK:
Create a structured learning path with:
1. Path Overview (title, description, type: manual/ai-generated)
2. Phases (3-6 phases, ordered by learning progression)
3. For each phase:
   - Title and description
   - Duration estimate
   - Skill requirements (prerequisite proficiency levels)
   - Milestones (concrete achievements)
   - Recommended resources (books, courses, projects)

OUTPUT FORMAT (JSON):
{
  "path": {
    "title": "string",
    "description": "string",
    "estimated_duration_weeks": "number"
  },
  "phases": [
    {
      "title": "string",
      "description": "string",
      "duration_weeks": "number",
      "skill_requirements": [
        {
          "skill_title": "string",
          "category": "string",
          "required_level": "beginner|intermediate|advanced|expert"
        }
      ],
      "milestones": [
        {
          "title": "string",
          "description": "string",
          "type": "goal-level|path-level|skill-level"
        }
      ],
      "resources": [
        {
          "title": "string",
          "type": "book|course|video|article|project|documentation",
          "author": "string",
          "url": "string",
          "estimated_hours": "number",
          "description": "string"
        }
      ]
    }
  ],
  "reasoning": "string - explain the learning path design rationale"
}

IMPORTANT:
- Make the path practical and achievable
- Consider the user's current skill level
- Prioritize hands-on projects and real-world application
- Include both foundational and advanced resources
- Suggest free resources when possible
- Provide clear milestones for tracking progress
`

const ResourceSuggestionPrompt = `You are an expert at recommending technical learning resources.

SKILL: {{.Skill.Title}}
CATEGORY: {{.Skill.Category}}
CURRENT LEVEL: {{.CurrentLevel}}
TARGET LEVEL: {{.TargetLevel}}
LEARNING STYLE: {{.LearningStyle}}
BUDGET: {{.Budget}}

TASK:
Recommend 5-10 high-quality learning resources to progress from {{.CurrentLevel}} to {{.TargetLevel}}.

OUTPUT FORMAT (JSON):
{
  "resources": [
    {
      "title": "string",
      "type": "book|course|video|article|project|documentation",
      "author": "string",
      "url": "string",
      "estimated_hours": "number",
      "cost": "free|paid",
      "description": "string",
      "why_recommended": "string"
    }
  ],
  "reasoning": "string - explain the resource selection rationale"
}

GUIDELINES:
- Prioritize {{.Budget}} resources
- Match {{.LearningStyle}} (e.g., top-down = projects first, bottom-up = theory first)
- Include diverse formats (books, courses, projects)
- Prefer well-reviewed, current resources (2023+)
- Start with foundational resources, progress to advanced
`

const ProgressAnalysisPrompt = `You are an expert career coach analyzing learning progress.

GOAL: {{.Goal.Title}}
LEARNING PATH: {{.Path.Title}}

PROGRESS LOGS (Last 30 days):
{{range .ProgressLogs}}
- {{.Date}}: {{.HoursInvested}} hours, Mood: {{.Mood}}
  Skills worked: {{.SkillsWorked}}
  Summary: {{.Body}}
{{end}}

CURRENT SKILLS:
{{range .CurrentSkills}}
- {{.Title}} ({{.Level}}, Status: {{.Status}})
{{end}}

TASK:
Analyze the user's progress and provide actionable insights.

OUTPUT FORMAT (JSON):
{
  "summary": "string - 2-3 sentence progress overview",
  "insights": [
    "string - key observation about progress patterns"
  ],
  "recommendations": [
    "string - specific actionable next step"
  ],
  "is_on_track": "boolean",
  "suggested_focus": [
    "string - skill or area to focus on next"
  ]
}

ANALYSIS GUIDELINES:
- Look for consistency patterns (regular vs sporadic)
- Identify skills with momentum vs stagnation
- Consider mood trends and energy levels
- Provide encouraging but honest assessment
- Suggest specific next actions, not generic advice
`
```

**Response Parser** (`internal/ai/gemini/parser.go`):
```go
package gemini

import (
    "encoding/json"
    "fmt"
    "github.com/illenko/growth.md/internal/core"
    "github.com/illenko/growth.md/internal/ai"
)

// PathGenerationOutput matches the JSON schema from the prompt
type PathGenerationOutput struct {
    Path      PathOutput      `json:"path"`
    Phases    []PhaseOutput   `json:"phases"`
    Reasoning string          `json:"reasoning"`
}

type PathOutput struct {
    Title               string `json:"title"`
    Description         string `json:"description"`
    EstimatedDurationWeeks int `json:"estimated_duration_weeks"`
}

type PhaseOutput struct {
    Title             string                  `json:"title"`
    Description       string                  `json:"description"`
    DurationWeeks     int                     `json:"duration_weeks"`
    SkillRequirements []SkillRequirementOutput `json:"skill_requirements"`
    Milestones        []MilestoneOutput       `json:"milestones"`
    Resources         []ResourceOutput        `json:"resources"`
}

// ... other output types

func ParsePathGeneration(responseText string) (*ai.PathGenerationResponse, error) {
    var output PathGenerationOutput

    if err := json.Unmarshal([]byte(responseText), &output); err != nil {
        return nil, fmt.Errorf("failed to parse AI response: %w", err)
    }

    // Convert to core entities
    path := &core.LearningPath{
        // ... map fields
    }

    phases := make([]*core.Phase, len(output.Phases))
    // ... convert phases

    return &ai.PathGenerationResponse{
        Path:      path,
        Phases:    phases,
        Resources: resources,
        Milestones: milestones,
        Reasoning: output.Reasoning,
    }, nil
}
```

**Tasks**:
- [ ] Add Gemini SDK dependency: `go get github.com/google/generative-ai-go`
- [ ] Create `internal/ai/gemini/client.go`
- [ ] Create `internal/ai/gemini/prompts.go` with all templates
- [ ] Create `internal/ai/gemini/parser.go` for response parsing
- [ ] Implement `GenerateLearningPath()`
- [ ] Implement `SuggestResources()`
- [ ] Implement `AnalyzeProgress()`
- [ ] Add retry logic with exponential backoff
- [ ] Add response validation and error handling
- [ ] Write unit tests with mocked responses

---

### 1.3 Add OpenAI Client (Future)

**Files to create** (stub for now):
- `internal/ai/openai/client.go`
- `internal/ai/openai/prompts.go`

**Implementation Notes**:
- Use same prompts as Gemini with minor adjustments
- Use `gpt-4` or `gpt-4-turbo` models
- Similar structured output approach
- Can be implemented later when users request it

**Tasks**:
- [ ] Create stub OpenAI client
- [ ] Return "not implemented" error
- [ ] Document how to implement when needed

---

### 1.4 Test AI Client

**Files to create**:
- `internal/ai/mock_client.go` - Mock for testing
- `internal/ai/gemini/client_test.go` - Unit tests

**Mock Client**:
```go
package ai

type MockClient struct {
    GenerateLearningPathFunc func(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error)
    // ... other funcs
}

func (m *MockClient) GenerateLearningPath(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error) {
    if m.GenerateLearningPathFunc != nil {
        return m.GenerateLearningPathFunc(ctx, req)
    }
    // Return default mock response
    return &PathGenerationResponse{
        Path: &core.LearningPath{
            Title: "Mock Learning Path",
            // ...
        },
    }, nil
}
```

**Test Coverage**:
- [ ] Test prompt template rendering
- [ ] Test response parsing (valid JSON)
- [ ] Test error handling (invalid JSON, API errors)
- [ ] Test retry logic
- [ ] Integration test with real Gemini API (manual, documented)

---

## Phase 2: CLI Commands for AI Features

**Goal**: Add user-facing commands to generate paths and get AI suggestions

### 2.1 Path Generation Command

**Update**: `internal/cli/path.go`

**Command**:
```bash
growth path generate <goal-id> [flags]
```

**Flags**:
- `--style` - Learning style: `top-down`, `bottom-up`, `project-based` (default: `project-based`)
- `--time` - Time commitment: e.g., `10 hours/week` (default: `5 hours/week`)
- `--background` - Additional background context
- `--provider` - AI provider: `gemini`, `openai` (default: from config)
- `--model` - Model override

**Implementation**:
```go
var pathGenerateCmd = &cobra.Command{
    Use:   "generate <goal-id>",
    Short: "Generate a learning path using AI",
    Long: `Generate a personalized learning path for a goal using AI.

The AI will analyze your goal, current skills, and preferences to create
a structured learning path with phases, milestones, and resource recommendations.

Examples:
  growth path generate goal-001
  growth path generate goal-001 --style top-down --time "10 hours/week"
  growth path generate goal-001 --background "I have 5 years Python experience"`,
    Args: cobra.ExactArgs(1),
    RunE: runPathGenerate,
}

func runPathGenerate(cmd *cobra.Command, args []string) error {
    goalID := core.EntityID(args[0])

    // Load goal
    goal, err := goalRepo.GetByIDWithBody(goalID)
    if err != nil {
        return fmt.Errorf("goal '%s' not found", goalID)
    }

    // Load current skills
    skills, err := skillRepo.GetAll()
    if err != nil {
        return fmt.Errorf("failed to load skills: %w", err)
    }

    // Initialize AI client
    aiConfig := ai.Config{
        Provider:    pathGenerateProvider,
        APIKey:      config.AI.APIKey,
        Model:       pathGenerateModel,
        Temperature: 0.7,
        MaxTokens:   8000,
    }

    client, err := ai.NewClient(aiConfig)
    if err != nil {
        return fmt.Errorf("failed to initialize AI client: %w", err)
    }

    // Show progress
    fmt.Printf("🤖 Generating learning path for: %s\n", goal.Title)
    fmt.Printf("   Provider: %s\n", client.Provider())
    fmt.Println()

    spinner := NewSpinner("Analyzing your goal and skills...")
    spinner.Start()

    // Generate path
    req := ai.PathGenerationRequest{
        Goal:           goal,
        CurrentSkills:  skills,
        Background:     pathGenerateBackground,
        LearningStyle:  pathGenerateStyle,
        TimeCommitment: pathGenerateTime,
        TargetDate:     goal.TargetDate,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    resp, err := client.GenerateLearningPath(ctx, req)
    if err != nil {
        spinner.Stop()
        return fmt.Errorf("failed to generate path: %w", err)
    }

    spinner.Stop()

    // Save path and related entities
    if err := saveGeneratedPath(resp, goalID); err != nil {
        return fmt.Errorf("failed to save path: %w", err)
    }

    // Display summary
    displayPathSummary(resp)

    return nil
}
```

**Tasks**:
- [ ] Add `pathGenerateCmd` to `path.go`
- [ ] Implement `runPathGenerate()`
- [ ] Add progress spinner during generation
- [ ] Implement `saveGeneratedPath()` helper
- [ ] Implement `displayPathSummary()` helper
- [ ] Link generated path to goal automatically
- [ ] Set path type as `ai-generated`
- [ ] Store AI provider and model metadata

---

### 2.2 Resource Suggestion Command

**Update**: `internal/cli/skill.go`

**Command**:
```bash
growth skill suggest-resources <skill-id> [flags]
```

**Flags**:
- `--target-level` - Target proficiency level (default: next level up)
- `--style` - Learning style (default: from config)
- `--budget` - Resource budget: `free`, `paid`, `any` (default: `any`)

**Implementation**:
```go
var skillSuggestResourcesCmd = &cobra.Command{
    Use:   "suggest-resources <skill-id>",
    Short: "Get AI-powered resource recommendations",
    Long: `Get personalized learning resource recommendations for a skill.

The AI will suggest books, courses, videos, and projects based on your
current level, target level, and learning preferences.

Examples:
  growth skill suggest-resources skill-001
  growth skill suggest-resources skill-001 --target-level advanced --budget free`,
    Args: cobra.ExactArgs(1),
    RunE: runSkillSuggestResources,
}
```

**Tasks**:
- [ ] Add `skillSuggestResourcesCmd` to `skill.go`
- [ ] Implement resource suggestion logic
- [ ] Display suggestions in table format
- [ ] Option to save suggestions: `--save` flag
- [ ] Show cost/time estimates

---

### 2.3 Progress Analysis Command

**Create**: `internal/cli/analyze.go`

**Command**:
```bash
growth analyze [goal-id]
```

**Description**: Get AI insights on your progress

**Implementation**:
```go
var analyzeCmd = &cobra.Command{
    Use:   "analyze [goal-id]",
    Short: "Get AI-powered progress insights",
    Long: `Analyze your learning progress and get personalized recommendations.

If a goal-id is provided, analyzes progress for that specific goal.
Otherwise, provides overall progress analysis.

Examples:
  growth analyze                  # Overall analysis
  growth analyze goal-001         # Goal-specific analysis`,
    Args: cobra.MaximumNArgs(1),
    RunE: runAnalyze,
}
```

**Output Example**:
```
🤖 Progress Analysis

SUMMARY
You've logged 45 hours over 12 sessions in the last 30 days. Your learning
momentum is strong with consistent 3-4 hour sessions.

INSIGHTS
✓ Strong progress on Python fundamentals (beginner → intermediate)
✓ Consistent learning schedule (3-4 sessions per week)
⚠ Docker skills haven't been practiced in 2 weeks
⚠ Less time invested when mood is "frustrated" - consider breaks

RECOMMENDATIONS
1. Focus on Docker this week - you're close to intermediate level
2. Start a small project combining Python + Docker
3. Review "Clean Code" chapters 4-6 based on your notes

ON TRACK: Yes - you're 65% through your path, on pace for target date

SUGGESTED FOCUS
- Docker containerization (skill-003)
- Python testing frameworks (skill-001)
```

**Tasks**:
- [ ] Create `internal/cli/analyze.go`
- [ ] Implement progress analysis
- [ ] Format insights nicely
- [ ] Add color coding for insights
- [ ] Cache analysis (avoid re-running frequently)

---

### 2.4 Update Config for AI

**Update**: `internal/storage/config.go`

**Add AI Configuration**:
```go
type Config struct {
    // ... existing fields

    AI AIConfig `yaml:"ai"`
}

type AIConfig struct {
    Provider         string  `yaml:"provider"`          // gemini, openai, anthropic
    APIKey          string  `yaml:"api_key"`            // or use env var
    Model           string  `yaml:"model"`              // model name
    Temperature     float32 `yaml:"temperature"`        // 0.0 - 1.0
    MaxTokens       int     `yaml:"max_tokens"`         // max output tokens
    DefaultStyle    string  `yaml:"default_style"`      // learning style preference
    DefaultBudget   string  `yaml:"default_budget"`     // resource budget preference
}
```

**Environment Variables**:
- `GEMINI_API_KEY` - Gemini API key
- `OPENAI_API_KEY` - OpenAI API key
- `ANTHROPIC_API_KEY` - Anthropic API key

**Default Config**:
```yaml
ai:
  provider: gemini
  api_key: ""  # Leave empty, use env var
  model: gemini-3-flash-preview
  temperature: 0.7
  max_tokens: 8000
  default_style: project-based
  default_budget: any
```

**Tasks**:
- [ ] Update `Config` struct
- [ ] Add AI config validation
- [ ] Add env var loading
- [ ] Update `growth init` to prompt for AI config
- [ ] Document API key setup in README

---

## Phase 3: MCP Server Integration

**Goal**: Expose growth.md to Claude Desktop and other MCP clients

### 3.1 Understand MCP Protocol

**Research**:
- Read [MCP Documentation](https://modelcontextprotocol.io/)
- Understand server/client model
- Review Go implementation examples
- Study tool schema definitions

**Key Concepts**:
- **Resources**: Read-only data (skills, goals, paths)
- **Tools**: Actions the AI can take (create skill, update progress)
- **Prompts**: Pre-defined workflows
- **Sampling**: AI-driven interactions

**Tasks**:
- [ ] Review MCP specification
- [ ] Find/create Go MCP SDK
- [ ] Understand Claude Desktop integration
- [ ] Plan resource and tool schemas

---

### 3.2 Implement MCP Server

**Files to create**:
- `cmd/growth-mcp/main.go` - MCP server entry point
- `internal/mcp/server.go` - MCP server implementation
- `internal/mcp/resources.go` - Resource handlers
- `internal/mcp/tools.go` - Tool handlers
- `internal/mcp/prompts.go` - Prompt definitions

**Server Structure**:
```go
package mcp

import (
    "github.com/illenko/growth.md/internal/storage"
)

type Server struct {
    skillRepo     storage.SkillRepository
    goalRepo      storage.GoalRepository
    resourceRepo  storage.ResourceRepository
    pathRepo      storage.PathRepository
    progressRepo  storage.ProgressRepository
    milestoneRepo storage.MilestoneRepository
}

func NewServer(repoPath string) (*Server, error) {
    // Initialize repositories
    // ...

    return &Server{
        skillRepo:     skillRepo,
        goalRepo:      goalRepo,
        resourceRepo:  resourceRepo,
        pathRepo:      pathRepo,
        progressRepo:  progressRepo,
        milestoneRepo: milestoneRepo,
    }, nil
}
```

**Resources** (`resources.go`):
```go
// Resource: List all skills
// URI: growth://skills
func (s *Server) ListSkills() ([]Skill, error)

// Resource: Get specific skill
// URI: growth://skills/{id}
func (s *Server) GetSkill(id string) (*Skill, error)

// Resource: List all goals
// URI: growth://goals
func (s *Server) ListGoals() ([]Goal, error)

// Resource: Get learning path
// URI: growth://paths/{id}
func (s *Server) GetPath(id string) (*Path, error)

// Resource: Recent progress
// URI: growth://progress/recent?days=30
func (s *Server) RecentProgress(days int) ([]Progress, error)

// Resource: Overview stats
// URI: growth://overview
func (s *Server) Overview() (*Overview, error)
```

**Tools** (`tools.go`):
```go
// Tool: Create a new skill
func (s *Server) CreateSkillTool() Tool {
    return Tool{
        Name: "create_skill",
        Description: "Create a new technical skill to track",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "title": map[string]string{
                    "type": "string",
                    "description": "Skill name (e.g., 'Python Programming')",
                },
                "category": map[string]string{
                    "type": "string",
                    "description": "Category (e.g., 'backend', 'frontend')",
                },
                "level": map[string]interface{}{
                    "type": "string",
                    "enum": []string{"beginner", "intermediate", "advanced", "expert"},
                },
            },
            "required": []string{"title", "category", "level"},
        },
    }
}

// Tool: Log progress
func (s *Server) LogProgressTool() Tool

// Tool: Update skill status
func (s *Server) UpdateSkillStatusTool() Tool

// Tool: Create goal
func (s *Server) CreateGoalTool() Tool

// Tool: Mark milestone achieved
func (s *Server) AchieveMilestoneTool() Tool

// Tool: Search across entities
func (s *Server) SearchTool() Tool
```

**Prompts** (`prompts.go`):
```go
// Prompt: Start tracking new goal
func (s *Server) TrackNewGoalPrompt() Prompt

// Prompt: Weekly progress review
func (s *Server) WeeklyReviewPrompt() Prompt

// Prompt: Skill gap analysis
func (s *Server) SkillGapAnalysisPrompt() Prompt
```

**Tasks**:
- [ ] Create `cmd/growth-mcp/main.go`
- [ ] Implement MCP server with stdio transport
- [ ] Implement all resource handlers (6 resources)
- [ ] Implement all tool handlers (6 tools)
- [ ] Implement prompt templates (3 prompts)
- [ ] Add JSON schema validation
- [ ] Add error handling and logging
- [ ] Test with MCP inspector

---

### 3.3 Claude Desktop Integration

**Configuration File**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "growth": {
      "command": "/path/to/growth-mcp",
      "args": [],
      "env": {
        "GROWTH_REPO_PATH": "/Users/you/growth"
      }
    }
  }
}
```

**Installation Script**: `scripts/install-mcp.sh`
```bash
#!/bin/bash
# Install growth.md MCP server for Claude Desktop

set -e

echo "📦 Building MCP server..."
go build -o bin/growth-mcp cmd/growth-mcp/main.go

echo "📝 Configuring Claude Desktop..."
CONFIG_DIR="$HOME/Library/Application Support/Claude"
CONFIG_FILE="$CONFIG_DIR/claude_desktop_config.json"

mkdir -p "$CONFIG_DIR"

# Get absolute path to binary
BIN_PATH="$(pwd)/bin/growth-mcp"

# Read current config or create new one
if [ -f "$CONFIG_FILE" ]; then
    # Merge with existing config
    jq ".mcpServers.growth = {\"command\": \"$BIN_PATH\", \"args\": [], \"env\": {\"GROWTH_REPO_PATH\": \"$(pwd)\"}}" "$CONFIG_FILE" > "$CONFIG_FILE.tmp"
    mv "$CONFIG_FILE.tmp" "$CONFIG_FILE"
else
    # Create new config
    cat > "$CONFIG_FILE" <<EOF
{
  "mcpServers": {
    "growth": {
      "command": "$BIN_PATH",
      "args": [],
      "env": {
        "GROWTH_REPO_PATH": "$(pwd)"
      }
    }
  }
}
EOF
fi

echo "✅ MCP server installed!"
echo ""
echo "Next steps:"
echo "1. Restart Claude Desktop"
echo "2. Look for the 🔌 icon to see available tools"
echo "3. Try: 'Show me my current skills' or 'Create a new goal'"
```

**Tasks**:
- [ ] Create installation script
- [ ] Test with Claude Desktop
- [ ] Document setup process
- [ ] Add troubleshooting guide
- [ ] Create demo video/GIF

---

### 3.4 Test MCP Integration

**Manual Testing**:
1. Build MCP server: `go build -o bin/growth-mcp cmd/growth-mcp/main.go`
2. Test with MCP inspector: `npx @modelcontextprotocol/inspector bin/growth-mcp`
3. Install in Claude Desktop
4. Test each resource and tool
5. Test prompt workflows

**Test Scenarios**:
- [ ] Read all skills → shows table of skills
- [ ] Read specific goal → shows goal details
- [ ] Create new skill → skill file created
- [ ] Log progress → progress file created
- [ ] Search "python" → finds relevant entities
- [ ] Use "Track New Goal" prompt → guided workflow
- [ ] Use "Weekly Review" prompt → summary of progress

**Documentation**:
- [ ] Create `docs/mcp-setup.md` with setup instructions
- [ ] Add example interactions
- [ ] Document all resources and tools
- [ ] Add troubleshooting section

---

## Phase 4: Polish & Documentation

### 4.1 API Key Management

**Security Best Practices**:
- Never commit API keys to git
- Use environment variables
- Add `.env` support
- Document key rotation

**Files to update**:
- `.gitignore` - Add `.env`
- `README.md` - Add API key setup section

**Tasks**:
- [ ] Add dotenv library: `go get github.com/joho/godotenv`
- [ ] Load `.env` in root command
- [ ] Document key setup process
- [ ] Add validation for missing keys

---

### 4.2 Error Handling & Rate Limits

**Implement**:
- Exponential backoff for API calls
- Graceful degradation when API unavailable
- Clear error messages for common issues
- Rate limit detection and handling

**Tasks**:
- [ ] Add retry logic with backoff
- [ ] Detect rate limit errors
- [ ] Add friendly error messages
- [ ] Add `--debug` flag for verbose output

---

### 4.3 Documentation

**Files to create/update**:
- `docs/ai-features.md` - AI feature documentation
- `docs/mcp-setup.md` - MCP setup guide
- `docs/api-keys.md` - API key management guide
- `README.md` - Add AI features section

**README Updates**:
```markdown
## 🤖 AI-Powered Features

### Learning Path Generation

Generate personalized learning paths using AI:

```bash
# Get a free Gemini API key: https://aistudio.google.com/apikey
export GEMINI_API_KEY="your-key-here"

# Generate a path for your goal
growth path generate goal-001 --style project-based --time "10 hours/week"
```

### Resource Recommendations

Get AI-powered resource suggestions:

```bash
growth skill suggest-resources skill-001 --target-level advanced --budget free
```

### Progress Analysis

Get insights on your learning journey:

```bash
growth analyze goal-001
```

## 🔌 Claude Desktop Integration

Connect growth.md to Claude Desktop using MCP:

```bash
# Install MCP server
./scripts/install-mcp.sh

# Restart Claude Desktop
# Now you can ask Claude: "Show me my current skills"
```

See [MCP Setup Guide](docs/mcp-setup.md) for details.
```

**Tasks**:
- [ ] Write comprehensive AI features guide
- [ ] Write MCP setup guide
- [ ] Add example workflows
- [ ] Add troubleshooting section
- [ ] Create demo GIFs/videos

---

## Phase 5: Testing & Validation

### 5.1 AI Integration Tests

**Files to create**:
- `tests/ai_integration_test.go`
- `internal/ai/gemini/client_test.go`

**Test Coverage**:
- [ ] Prompt rendering with different contexts
- [ ] Response parsing (valid responses)
- [ ] Error handling (invalid JSON, API errors)
- [ ] Retry logic
- [ ] Rate limit handling
- [ ] Integration test with real API (manual)

---

### 5.2 MCP Integration Tests

**Files to create**:
- `tests/mcp_integration_test.go`
- `internal/mcp/server_test.go`

**Test Coverage**:
- [ ] All resources return valid data
- [ ] All tools execute correctly
- [ ] Tool input validation
- [ ] Error handling
- [ ] Integration test with MCP inspector

---

### 5.3 End-to-End Testing

**Manual Test Scenarios**:
1. New user flow: init → create goal → generate path → view path
2. Resource suggestion: create skill → get suggestions → save resources
3. Progress tracking: log progress → analyze → get recommendations
4. MCP flow: Claude Desktop → read skills → create goal → log progress

**Tasks**:
- [ ] Create test plan document
- [ ] Execute all test scenarios
- [ ] Document results
- [ ] Fix any issues found

---

## Implementation Priority

### Week 1: AI Foundation
- [x] Phase 1.1: AI Client Interface ✅
- [x] Phase 1.2: Gemini Client Implementation ✅
- [x] Phase 1.3: OpenAI Stub (for future) ✅
- [x] Phase 1.4: Mock Client for Testing ✅
- [x] Code compiles successfully ✅

### Week 2: CLI Integration
- [x] Phase 2.1: Path Generation Command ✅
- [x] Phase 2.2: Resource Suggestion Command ✅
- [x] Phase 2.3: Progress Analysis Command ✅
- [x] Phase 2.4: AI Config ✅
- [ ] Manual testing with real API

### Week 3: MCP Server
- [ ] Phase 3.1: MCP Research ✅
- [ ] Phase 3.2: MCP Server Implementation ✅
- [ ] Phase 3.3: Claude Desktop Integration ✅
- [ ] Phase 3.4: MCP Testing ✅

### Week 4: Polish
- [ ] Phase 2.3: Progress Analysis ✅
- [ ] Phase 4: Documentation ✅
- [ ] Phase 5: Testing & Validation ✅
- [ ] Release! 🚀

---

## Success Metrics

**AI Integration**:
- ✅ Can generate learning paths in < 30 seconds
- ✅ Path quality is practical and achievable
- ✅ Resource suggestions are relevant and current
- ✅ Progress analysis provides actionable insights
- ✅ Supports multiple AI providers (at least Gemini)

**MCP Integration**:
- ✅ All resources accessible from Claude Desktop
- ✅ All tools work correctly
- ✅ Setup process takes < 5 minutes
- ✅ Prompts provide useful workflows
- ✅ Error handling is graceful

**Developer Experience**:
- ✅ Clear documentation
- ✅ Easy API key setup
- ✅ Good error messages
- ✅ Extensible architecture

---

## Future Enhancements (Post-MVP)

### Additional AI Providers
- [ ] OpenAI (GPT-4)
- [ ] Anthropic (Claude)
- [ ] Local models (Ollama)
- [ ] Azure OpenAI

### Advanced AI Features
- [ ] Skill gap analysis (what's missing?)
- [ ] Career path recommendations
- [ ] Interview prep suggestions
- [ ] Resume bullet point generator
- [ ] Learning style detection
- [ ] Adaptive difficulty adjustment

### MCP Enhancements
- [ ] Rich resource metadata
- [ ] Progress charts/visualizations
- [ ] Multi-repo support
- [ ] Team/collaboration features
- [ ] Export to other tools

### Intelligence Features
- [ ] Learning pattern detection
- [ ] Optimal study time recommendations
- [ ] Burnout detection and prevention
- [ ] Skill correlation analysis
- [ ] Market trend integration

---

## Notes

- **Gemini Free Tier**: 15 requests/minute, 1500 requests/day - more than enough for personal use
- **MCP is Experimental**: Protocol may change, stay updated
- **Privacy**: All data stays local, API calls are stateless
- **Costs**: Gemini is free, OpenAI/Anthropic are paid (add cost estimates in docs)

---

**Last Updated**: 2025-12-23
