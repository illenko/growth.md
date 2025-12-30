package service

import (
	"testing"

	"github.com/illenko/growth.md/internal/core"
)

func TestGetNextLevel(t *testing.T) {
	tests := []struct {
		current  core.ProficiencyLevel
		expected core.ProficiencyLevel
	}{
		{core.LevelBeginner, core.LevelIntermediate},
		{core.LevelIntermediate, core.LevelAdvanced},
		{core.LevelAdvanced, core.LevelExpert},
		{core.LevelExpert, core.LevelExpert},
	}

	for _, tt := range tests {
		t.Run(string(tt.current), func(t *testing.T) {
			result := getNextLevel(tt.current)
			if result != tt.expected {
				t.Errorf("getNextLevel(%s) = %s, want %s", tt.current, result, tt.expected)
			}
		})
	}
}
