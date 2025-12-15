package storage

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const frontmatterDelimiter = "---"

// ParseFrontmatter extracts YAML frontmatter and body from markdown content.
// Expected format:
//
//	---
//	yaml: content
//	---
//	markdown body
//
// Returns the frontmatter as a map, the body as a string, and any error.
func ParseFrontmatter(content []byte) (frontmatter map[string]interface{}, body string, err error) {
	if len(content) == 0 {
		return nil, "", errors.New("empty content")
	}

	contentStr := string(content)

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(contentStr, frontmatterDelimiter) {
		// No frontmatter, treat entire content as body
		return make(map[string]interface{}), contentStr, nil
	}

	// Find the second delimiter
	lines := strings.Split(contentStr, "\n")
	if len(lines) < 3 {
		return nil, "", errors.New("invalid frontmatter: too few lines")
	}

	// Skip first line (opening ---)
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == frontmatterDelimiter {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil, "", errors.New("invalid frontmatter: missing closing delimiter")
	}

	// Extract frontmatter YAML (between delimiters)
	frontmatterLines := lines[1:endIdx]
	frontmatterYAML := strings.Join(frontmatterLines, "\n")

	// Parse YAML
	frontmatter = make(map[string]interface{})
	if len(frontmatterYAML) > 0 {
		if err := yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter); err != nil {
			return nil, "", fmt.Errorf("failed to parse frontmatter YAML: %w", err)
		}
	}

	// Extract body (everything after closing delimiter)
	if endIdx+1 < len(lines) {
		body = strings.Join(lines[endIdx+1:], "\n")
		body = strings.TrimSpace(body)
	}

	return frontmatter, body, nil
}

// SerializeFrontmatter combines frontmatter and body into markdown with YAML frontmatter.
// The frontmatter parameter can be any struct or map that can be marshaled to YAML.
func SerializeFrontmatter(frontmatter interface{}, body string) ([]byte, error) {
	if frontmatter == nil {
		return nil, errors.New("frontmatter cannot be nil")
	}

	// Marshal frontmatter to YAML
	yamlBytes, err := yaml.Marshal(frontmatter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal frontmatter to YAML: %w", err)
	}

	var buf bytes.Buffer

	// Write opening delimiter
	buf.WriteString(frontmatterDelimiter)
	buf.WriteString("\n")

	// Write YAML content
	buf.Write(yamlBytes)

	// Write closing delimiter
	buf.WriteString(frontmatterDelimiter)
	buf.WriteString("\n")

	// Write body if present
	if body != "" {
		buf.WriteString("\n")
		buf.WriteString(body)
		buf.WriteString("\n")
	}

	return buf.Bytes(), nil
}
