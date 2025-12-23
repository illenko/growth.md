package gemini

const PathGenerationPrompt = `You are an expert career coach for software engineers. Generate a personalized learning path.

GOAL: {{.Goal.Title}}
GOAL DESCRIPTION: {{.Goal.Body}}
PRIORITY: {{.Goal.Priority}}
{{if .Goal.TargetDate}}TARGET DATE: {{.Goal.TargetDate.Format "2006-01-02"}}{{end}}

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
1. Path Overview (title, description, estimated duration in weeks)
2. Phases (3-6 phases, ordered by learning progression)
3. For each phase:
   - Title and description
   - Duration estimate (in weeks)
   - Skill requirements (prerequisite proficiency levels)
   - Milestones (concrete achievements)
   - Recommended resources (books, courses, projects)

OUTPUT FORMAT (JSON):
{
  "path": {
    "title": "string",
    "description": "string",
    "estimated_duration_weeks": 12
  },
  "phases": [
    {
      "title": "string",
      "description": "string",
      "duration_weeks": 3,
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
          "estimated_hours": 10,
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
- Ensure all JSON fields use exact names as specified above
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
      "estimated_hours": 10,
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
- Ensure all JSON fields use exact names as specified above
`

const ProgressAnalysisPrompt = `You are an expert career coach analyzing learning progress.

GOAL: {{.Goal.Title}}
LEARNING PATH: {{.Path.Title}}

PROGRESS LOGS (Last 30 days):
{{range .ProgressLogs}}
- {{.Date.Format "2006-01-02"}}: {{.HoursInvested}} hours{{if .Mood}}, Mood: {{.Mood}}{{end}}
{{if .Body}}  Summary: {{.Body}}{{end}}
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
  "is_on_track": true,
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
- Ensure all JSON fields use exact names as specified above
`
