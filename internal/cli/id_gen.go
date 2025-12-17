package cli

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/illenko/growth.md/internal/core"
)

func GenerateNextID(entityType string) (core.EntityID, error) {
	return GenerateNextIDInPath(entityType, repoPath)
}

func GenerateNextIDInPath(entityType string, basePath string) (core.EntityID, error) {
	var pattern string
	switch entityType {
	case "skill":
		pattern = filepath.Join(basePath, "skills", "skill-*.md")
	case "goal":
		pattern = filepath.Join(basePath, "goals", "goal-*.md")
	case "path":
		pattern = filepath.Join(basePath, "paths", "path-*.md")
	case "phase":
		pattern = filepath.Join(basePath, "phases", "phase-*.md")
	case "resource":
		pattern = filepath.Join(basePath, "resources", "resource-*.md")
	case "milestone":
		pattern = filepath.Join(basePath, "milestones", "milestone-*.md")
	case "progress":
		pattern = filepath.Join(basePath, "progress", "progress-*.md")
	default:
		return "", fmt.Errorf("unknown entity type: %s", entityType)
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to scan files: %w", err)
	}

	maxID := 0
	idPattern := regexp.MustCompile(fmt.Sprintf(`%s-(\d+)`, entityType))

	for _, match := range matches {
		basename := filepath.Base(match)
		if submatch := idPattern.FindStringSubmatch(basename); submatch != nil {
			id, err := strconv.Atoi(submatch[1])
			if err == nil && id > maxID {
				maxID = id
			}
		}
	}

	nextID := maxID + 1
	return core.EntityID(fmt.Sprintf("%s-%03d", entityType, nextID)), nil
}

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)

	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.TrimRight(slug, "-")
	}

	if slug == "" {
		slug = "untitled"
	}

	return slug
}

func GenerateFileName(id core.EntityID, title string) string {
	slug := GenerateSlug(title)
	return fmt.Sprintf("%s-%s.md", id, slug)
}
