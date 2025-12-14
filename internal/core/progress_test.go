package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProgressLog(t *testing.T) {
	weekDate := time.Date(2025, 3, 19, 0, 0, 0, 0, time.UTC) // Wednesday

	t.Run("creates valid progress log", func(t *testing.T) {
		log, err := NewProgressLog("progress-001", weekDate)

		require.NoError(t, err)
		assert.Equal(t, EntityID("progress-001"), log.ID)
		// Should normalize to Monday of that week
		expectedMonday := time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, expectedMonday, log.WeekOf)
		assert.Equal(t, 0.0, log.HoursInvested)
		assert.Empty(t, log.SkillsWorked)
		assert.Empty(t, log.ResourcesUsed)
		assert.Empty(t, log.MilestonesAchieved)
		assert.Empty(t, log.Mood)
		assert.False(t, log.Created.IsZero())
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		_, err := NewProgressLog("", weekDate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("fails with zero time", func(t *testing.T) {
		_, err := NewProgressLog("progress-001", time.Time{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "weekOf is required")
	})
}

func TestProgressLog_Validate(t *testing.T) {
	tests := []struct {
		name    string
		log     *ProgressLog
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid progress log",
			log: &ProgressLog{
				ID:            "progress-001",
				WeekOf:        time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
				HoursInvested: 10.5,
				Timestamps:    NewTimestamps(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			log: &ProgressLog{
				WeekOf:     time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "ID is required",
		},
		{
			name: "zero weekOf",
			log: &ProgressLog{
				ID:         "progress-001",
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "weekOf is required",
		},
		{
			name: "negative hours invested",
			log: &ProgressLog{
				ID:            "progress-001",
				WeekOf:        time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
				HoursInvested: -5,
				Timestamps:    NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProgressLog_AddSkillWorked(t *testing.T) {
	log, _ := NewProgressLog("progress-001", time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC))

	t.Run("adds skill", func(t *testing.T) {
		log.AddSkillWorked("skill-001")
		assert.Len(t, log.SkillsWorked, 1)
		assert.Contains(t, log.SkillsWorked, EntityID("skill-001"))
	})

	t.Run("adds multiple skills", func(t *testing.T) {
		log.AddSkillWorked("skill-002")
		log.AddSkillWorked("skill-003")
		assert.Len(t, log.SkillsWorked, 3)
	})

	t.Run("does not add duplicate", func(t *testing.T) {
		log.AddSkillWorked("skill-001")
		assert.Len(t, log.SkillsWorked, 3)
	})
}

func TestProgressLog_AddResourceUsed(t *testing.T) {
	log, _ := NewProgressLog("progress-001", time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC))

	t.Run("adds resource", func(t *testing.T) {
		log.AddResourceUsed("resource-001")
		assert.Len(t, log.ResourcesUsed, 1)
		assert.Contains(t, log.ResourcesUsed, EntityID("resource-001"))
	})

	t.Run("adds multiple resources", func(t *testing.T) {
		log.AddResourceUsed("resource-002")
		log.AddResourceUsed("resource-003")
		assert.Len(t, log.ResourcesUsed, 3)
	})

	t.Run("does not add duplicate", func(t *testing.T) {
		log.AddResourceUsed("resource-001")
		assert.Len(t, log.ResourcesUsed, 3)
	})
}

func TestProgressLog_AddMilestoneAchieved(t *testing.T) {
	log, _ := NewProgressLog("progress-001", time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC))

	t.Run("adds milestone", func(t *testing.T) {
		log.AddMilestoneAchieved("milestone-001")
		assert.Len(t, log.MilestonesAchieved, 1)
		assert.Contains(t, log.MilestonesAchieved, EntityID("milestone-001"))
	})

	t.Run("adds multiple milestones", func(t *testing.T) {
		log.AddMilestoneAchieved("milestone-002")
		log.AddMilestoneAchieved("milestone-003")
		assert.Len(t, log.MilestonesAchieved, 3)
	})

	t.Run("does not add duplicate", func(t *testing.T) {
		log.AddMilestoneAchieved("milestone-001")
		assert.Len(t, log.MilestonesAchieved, 3)
	})
}

func TestProgressLog_SetHoursInvested(t *testing.T) {
	log, _ := NewProgressLog("progress-001", time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC))

	t.Run("sets valid hours", func(t *testing.T) {
		err := log.SetHoursInvested(15.5)
		assert.NoError(t, err)
		assert.Equal(t, 15.5, log.HoursInvested)
	})

	t.Run("sets zero hours", func(t *testing.T) {
		err := log.SetHoursInvested(0)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, log.HoursInvested)
	})

	t.Run("fails with negative hours", func(t *testing.T) {
		err := log.SetHoursInvested(-10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
	})
}

func TestProgressLog_SetMood(t *testing.T) {
	log, _ := NewProgressLog("progress-001", time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC))

	log.SetMood("motivated")

	assert.Equal(t, "motivated", log.Mood)
}

func TestGetStartOfWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Monday stays Monday",
			input:    time.Date(2025, 3, 17, 10, 30, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Wednesday goes to Monday",
			input:    time.Date(2025, 3, 19, 15, 45, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Friday goes to Monday",
			input:    time.Date(2025, 3, 21, 9, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Saturday goes to Monday",
			input:    time.Date(2025, 3, 22, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Sunday goes to Monday",
			input:    time.Date(2025, 3, 23, 18, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Tuesday goes to Monday",
			input:    time.Date(2025, 3, 18, 8, 15, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Thursday goes to Monday",
			input:    time.Date(2025, 3, 20, 14, 30, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStartOfWeek(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
