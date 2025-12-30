package gemini

import (
	"context"
	"testing"

	"github.com/illenko/growth.md/internal/ai"
	"github.com/illenko/growth.md/internal/core"
)

func TestParsePathGeneration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		pathID      core.EntityID
		goalID      core.EntityID
		expectError bool
	}{
		{
			name: "valid path generation response",
			input: `{
				"path": {
					"title": "Test Learning Path",
					"description": "A test path",
					"estimated_duration_weeks": 12
				},
				"phases": [
					{
						"title": "Phase 1",
						"description": "First phase",
						"duration_weeks": 4,
						"skill_requirements": [],
						"milestones": [
							{
								"title": "Complete basics",
								"description": "Finish foundational work",
								"type": "path-level"
							}
						],
						"resources": [
							{
								"title": "Test Book",
								"type": "book",
								"author": "Test Author",
								"url": "https://example.com",
								"estimated_hours": 10,
								"description": "A test book"
							}
						]
					}
				],
				"reasoning": "This is a test path"
			}`,
			pathID:      "path-001",
			goalID:      "goal-001",
			expectError: false,
		},
		{
			name:        "invalid json",
			input:       `{invalid json`,
			pathID:      "path-001",
			goalID:      "goal-001",
			expectError: true,
		},
		{
			name: "missing path title",
			input: `{
				"path": {
					"title": "",
					"description": "No title",
					"estimated_duration_weeks": 12
				},
				"phases": [],
				"reasoning": "Test"
			}`,
			pathID:      "path-001",
			goalID:      "goal-001",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ParsePathGeneration(tt.input, tt.pathID, tt.goalID)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if resp.Path.ID != tt.pathID {
				t.Errorf("expected path ID %s, got %s", tt.pathID, resp.Path.ID)
			}

			if resp.Path.Type != core.PathTypeAIGenerated {
				t.Errorf("expected AI generated path type, got %s", resp.Path.Type)
			}
		})
	}
}

func TestParseResourceSuggestion(t *testing.T) {
	input := `{
		"resources": [
			{
				"title": "Test Resource",
				"type": "course",
				"author": "Test Author",
				"url": "https://example.com",
				"estimated_hours": 20,
				"cost": "free",
				"description": "A test resource",
				"why_recommended": "Good for beginners"
			}
		],
		"reasoning": "These resources are great"
	}`

	resp, err := ParseResourceSuggestion(input, "skill-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(resp.Resources))
	}

	if resp.Resources[0].SkillID != "skill-001" {
		t.Errorf("expected skill ID skill-001, got %s", resp.Resources[0].SkillID)
	}

	if resp.Resources[0].Type != core.ResourceCourse {
		t.Errorf("expected course type, got %s", resp.Resources[0].Type)
	}
}

func TestParseProgressAnalysis(t *testing.T) {
	input := `{
		"summary": "Great progress!",
		"insights": ["You're consistent", "Good momentum"],
		"recommendations": ["Focus on Docker", "Review Python"],
		"is_on_track": true,
		"suggested_focus": ["Docker", "Python testing"]
	}`

	resp, err := ParseProgressAnalysis(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsOnTrack {
		t.Error("expected IsOnTrack to be true")
	}

	if len(resp.Insights) != 2 {
		t.Errorf("expected 2 insights, got %d", len(resp.Insights))
	}

	if len(resp.Recommendations) != 2 {
		t.Errorf("expected 2 recommendations, got %d", len(resp.Recommendations))
	}
}

func TestMockClient(t *testing.T) {
	mockClient := &ai.MockClient{
		ProviderName: "test-mock",
	}

	if mockClient.Provider() != "test-mock" {
		t.Errorf("expected provider 'test-mock', got %s", mockClient.Provider())
	}

	// Test default response
	ctx := context.Background()
	req := ai.PathGenerationRequest{
		Goal: &core.Goal{
			ID:    "goal-001",
			Title: "Test Goal",
		},
	}

	resp, err := mockClient.GenerateLearningPath(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Path.Title != "Mock Learning Path" {
		t.Errorf("expected 'Mock Learning Path', got %s", resp.Path.Title)
	}
}
