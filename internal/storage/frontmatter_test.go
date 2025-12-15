package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFrontmatter(t *testing.T) {
	t.Run("parses valid frontmatter with body", func(t *testing.T) {
		content := []byte(`---
id: skill-001
title: Python
category: programming
level: intermediate
---

This is the markdown body content.

It can have multiple paragraphs.`)

		frontmatter, body, err := ParseFrontmatter(content)

		require.NoError(t, err)
		assert.Equal(t, "skill-001", frontmatter["id"])
		assert.Equal(t, "Python", frontmatter["title"])
		assert.Equal(t, "programming", frontmatter["category"])
		assert.Equal(t, "intermediate", frontmatter["level"])
		assert.Contains(t, body, "This is the markdown body content")
		assert.Contains(t, body, "It can have multiple paragraphs")
	})

	t.Run("parses frontmatter without body", func(t *testing.T) {
		content := []byte(`---
id: skill-001
title: Python
---`)

		frontmatter, body, err := ParseFrontmatter(content)

		require.NoError(t, err)
		assert.Equal(t, "skill-001", frontmatter["id"])
		assert.Equal(t, "Python", frontmatter["title"])
		assert.Empty(t, body)
	})

	t.Run("handles content without frontmatter", func(t *testing.T) {
		content := []byte(`This is just markdown content without frontmatter.`)

		frontmatter, body, err := ParseFrontmatter(content)

		require.NoError(t, err)
		assert.Empty(t, frontmatter)
		assert.Equal(t, "This is just markdown content without frontmatter.", body)
	})

	t.Run("handles empty frontmatter", func(t *testing.T) {
		content := []byte(`---
---

Body content here.`)

		frontmatter, body, err := ParseFrontmatter(content)

		require.NoError(t, err)
		assert.Empty(t, frontmatter)
		assert.Contains(t, body, "Body content here")
	})

	t.Run("fails with empty content", func(t *testing.T) {
		_, _, err := ParseFrontmatter([]byte{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty content")
	})

	t.Run("fails with missing closing delimiter", func(t *testing.T) {
		content := []byte(`---
id: skill-001
title: Python`)

		_, _, err := ParseFrontmatter(content)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing closing delimiter")
	})

	t.Run("fails with malformed YAML", func(t *testing.T) {
		content := []byte(`---
id: skill-001
title: [unclosed bracket
---`)

		_, _, err := ParseFrontmatter(content)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse frontmatter YAML")
	})

	t.Run("handles nested YAML structures", func(t *testing.T) {
		content := []byte(`---
id: path-001
skills:
  - skill-001
  - skill-002
metadata:
  generated: true
  version: 1
---

Path description.`)

		frontmatter, body, err := ParseFrontmatter(content)

		require.NoError(t, err)
		assert.Equal(t, "path-001", frontmatter["id"])
		assert.IsType(t, []interface{}{}, frontmatter["skills"])
		assert.IsType(t, map[string]interface{}{}, frontmatter["metadata"])
		assert.Contains(t, body, "Path description")
	})
}

func TestSerializeFrontmatter(t *testing.T) {
	t.Run("serializes map with body", func(t *testing.T) {
		frontmatter := map[string]interface{}{
			"id":       "skill-001",
			"title":    "Python",
			"category": "programming",
		}
		body := "This is the body content."

		result, err := SerializeFrontmatter(frontmatter, body)

		require.NoError(t, err)
		assert.Contains(t, string(result), "---")
		assert.Contains(t, string(result), "id: skill-001")
		assert.Contains(t, string(result), "title: Python")
		assert.Contains(t, string(result), "This is the body content")
	})

	t.Run("serializes map without body", func(t *testing.T) {
		frontmatter := map[string]interface{}{
			"id":    "skill-001",
			"title": "Python",
		}

		result, err := SerializeFrontmatter(frontmatter, "")

		require.NoError(t, err)
		assert.Contains(t, string(result), "---")
		assert.Contains(t, string(result), "id: skill-001")
		assert.NotContains(t, string(result), "\n\n\n") // No extra newlines
	})

	t.Run("serializes struct with tags", func(t *testing.T) {
		type TestStruct struct {
			ID    string `yaml:"id"`
			Title string `yaml:"title"`
		}

		frontmatter := TestStruct{
			ID:    "skill-001",
			Title: "Python",
		}
		body := "Body content."

		result, err := SerializeFrontmatter(frontmatter, body)

		require.NoError(t, err)
		assert.Contains(t, string(result), "id: skill-001")
		assert.Contains(t, string(result), "title: Python")
		assert.Contains(t, string(result), "Body content")
	})

	t.Run("fails with nil frontmatter", func(t *testing.T) {
		_, err := SerializeFrontmatter(nil, "body")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("handles empty map", func(t *testing.T) {
		frontmatter := map[string]interface{}{}
		body := "Just body content."

		result, err := SerializeFrontmatter(frontmatter, body)

		require.NoError(t, err)
		assert.Contains(t, string(result), "---")
		assert.Contains(t, string(result), "Just body content")
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("parse and serialize maintains data", func(t *testing.T) {
		original := map[string]interface{}{
			"id":       "skill-001",
			"title":    "Python",
			"category": "programming",
			"level":    "intermediate",
		}
		originalBody := "This is the markdown body.\n\nWith multiple paragraphs."

		// Serialize
		serialized, err := SerializeFrontmatter(original, originalBody)
		require.NoError(t, err)

		// Parse back
		parsed, parsedBody, err := ParseFrontmatter(serialized)
		require.NoError(t, err)

		// Verify
		assert.Equal(t, original["id"], parsed["id"])
		assert.Equal(t, original["title"], parsed["title"])
		assert.Equal(t, original["category"], parsed["category"])
		assert.Equal(t, original["level"], parsed["level"])
		assert.Contains(t, parsedBody, "This is the markdown body")
		assert.Contains(t, parsedBody, "With multiple paragraphs")
	})
}
