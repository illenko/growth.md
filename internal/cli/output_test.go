package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type testStruct struct {
	ID      string    `yaml:"id"`
	Name    string    `yaml:"name"`
	Status  string    `yaml:"status"`
	Created time.Time `yaml:"created"`
}

func TestPrintJSON(t *testing.T) {
	t.Run("prints valid JSON", func(t *testing.T) {
		data := testStruct{
			ID:      "test-001",
			Name:    "Test Item",
			Status:  "active",
			Created: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := PrintJSON(data)
		require.NoError(t, err)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)

		var result testStruct
		err = json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, "test-001", result.ID)
		assert.Equal(t, "Test Item", result.Name)
	})
}

func TestPrintYAML(t *testing.T) {
	t.Run("prints valid YAML", func(t *testing.T) {
		data := testStruct{
			ID:      "test-001",
			Name:    "Test Item",
			Status:  "active",
			Created: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := PrintYAML(data)
		require.NoError(t, err)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)

		var result testStruct
		err = yaml.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, "test-001", result.ID)
		assert.Equal(t, "Test Item", result.Name)
	})
}

func TestPrintTable(t *testing.T) {
	t.Run("prints single struct", func(t *testing.T) {
		data := testStruct{
			ID:     "test-001",
			Name:   "Test Item",
			Status: "active",
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := PrintTable(data)
		require.NoError(t, err)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.Contains(t, output, "ID")
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "test-001")
		assert.Contains(t, output, "Test Item")
	})

	t.Run("prints slice of structs", func(t *testing.T) {
		data := []testStruct{
			{ID: "test-001", Name: "Item 1", Status: "active"},
			{ID: "test-002", Name: "Item 2", Status: "completed"},
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := PrintTable(data)
		require.NoError(t, err)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.Contains(t, output, "test-001")
		assert.Contains(t, output, "test-002")
		assert.Contains(t, output, "Item 1")
		assert.Contains(t, output, "Item 2")
	})

	t.Run("handles empty slice", func(t *testing.T) {
		data := []testStruct{}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := PrintTable(data)
		require.NoError(t, err)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.Contains(t, output, "No results found")
	})

	t.Run("handles Skills", func(t *testing.T) {
		skill, _ := core.NewSkill("skill-001", "Go Programming", "programming", core.LevelIntermediate)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := PrintTable(skill)
		require.NoError(t, err)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.Contains(t, output, "skill-001")
		assert.Contains(t, output, "Go Prog")
		assert.Contains(t, output, "program")
		assert.Contains(t, output, "interme")
	})
}

func TestFormatFieldValue(t *testing.T) {
	t.Run("formats time correctly", func(t *testing.T) {
		now := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		v := reflect.ValueOf(now)

		result := formatFieldValue(v)

		assert.Equal(t, "2024-01-15", result)
	})

	t.Run("formats slice with max 3 items", func(t *testing.T) {
		slice := []string{"a", "b", "c", "d", "e"}
		v := reflect.ValueOf(slice)

		result := formatFieldValue(v)

		assert.Contains(t, result, "a")
		assert.Contains(t, result, "...")
	})

	t.Run("handles nil pointer", func(t *testing.T) {
		var ptr *string
		v := reflect.ValueOf(ptr)

		result := formatFieldValue(v)

		assert.Equal(t, "", result)
	})
}
