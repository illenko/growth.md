package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/illenko/growth.md/internal/core"
	"gopkg.in/yaml.v3"
)

var _ Repository[any] = (*FilesystemRepository[any])(nil)

// FilesystemRepository implements the Repository interface using the local filesystem.
// Entities are stored as markdown files with YAML frontmatter.
type FilesystemRepository[T any] struct {
	basePath   string // Base directory for this repository
	entityType string // Entity type name (e.g., "skill", "goal")
}

// NewFilesystemRepository creates a new filesystem-based repository.
// basePath is the directory where entity files will be stored.
// entityType is used for file naming (e.g., "skill" -> "skill-001-python.md").
func NewFilesystemRepository[T any](basePath, entityType string) (*FilesystemRepository[T], error) {
	if basePath == "" {
		return nil, errors.New("basePath cannot be empty")
	}
	if entityType == "" {
		return nil, errors.New("entityType cannot be empty")
	}

	// Ensure directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", basePath, err)
	}

	return &FilesystemRepository[T]{
		basePath:   basePath,
		entityType: entityType,
	}, nil
}

// Create persists a new entity to storage.
func (r *FilesystemRepository[T]) Create(entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	id, err := r.getEntityID(entity)
	if err != nil {
		return fmt.Errorf("failed to get entity ID: %w", err)
	}

	// Check if entity already exists
	exists, err := r.Exists(id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("entity with ID %s already exists", id)
	}

	// Generate filename
	title := r.getEntityTitle(entity)
	filename := r.generateFileName(id, title)
	filepath := filepath.Join(r.basePath, filename)

	// Serialize entity
	content, err := r.serializeEntity(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize entity: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filepath, err)
	}

	return nil
}

// GetByID retrieves an entity by its ID (metadata only, without body).
func (r *FilesystemRepository[T]) GetByID(id core.EntityID) (*T, error) {
	return r.getByID(id, false)
}

// GetByIDWithBody retrieves an entity by its ID including the markdown body.
func (r *FilesystemRepository[T]) GetByIDWithBody(id core.EntityID) (*T, error) {
	return r.getByID(id, true)
}

// getByID is the internal implementation for both GetByID methods.
func (r *FilesystemRepository[T]) getByID(id core.EntityID, includeBody bool) (*T, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	// Find file by ID
	filePath, err := r.findFileByID(id)
	if err != nil {
		return nil, err
	}

	// Read and parse file
	entity, err := r.parseEntityFromFile(filePath, includeBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entity from file: %w", err)
	}

	return entity, nil
}

// GetAll retrieves all entities of this type (metadata only, without bodies).
func (r *FilesystemRepository[T]) GetAll() ([]*T, error) {
	pattern := filepath.Join(r.basePath, fmt.Sprintf("%s-*.md", r.entityType))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	entities := make([]*T, 0, len(matches))
	for _, filePath := range matches {
		entity, err := r.parseEntityFromFile(filePath, false)
		if err != nil {
			// Log error but continue with other files
			continue
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

// Update persists changes to an existing entity.
func (r *FilesystemRepository[T]) Update(entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	id, err := r.getEntityID(entity)
	if err != nil {
		return fmt.Errorf("failed to get entity ID: %w", err)
	}

	// Find existing file
	oldFilePath, err := r.findFileByID(id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	// Generate new filename (title might have changed)
	title := r.getEntityTitle(entity)
	newFilename := r.generateFileName(id, title)
	newFilePath := filepath.Join(r.basePath, newFilename)

	// Serialize entity
	content, err := r.serializeEntity(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize entity: %w", err)
	}

	// Write to file
	if err := os.WriteFile(newFilePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Remove old file if filename changed
	if oldFilePath != newFilePath {
		if err := os.Remove(oldFilePath); err != nil {
			// Try to clean up the new file
			os.Remove(newFilePath)
			return fmt.Errorf("failed to remove old file: %w", err)
		}
	}

	return nil
}

// Delete removes an entity from storage by its ID.
func (r *FilesystemRepository[T]) Delete(id core.EntityID) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	filePath, err := r.findFileByID(id)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Search finds entities matching the given query string.
func (r *FilesystemRepository[T]) Search(query string) ([]*T, error) {
	if query == "" {
		return r.GetAll()
	}

	allEntities, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)
	var results []*T

	for _, entity := range allEntities {
		// Search in title and tags
		title := strings.ToLower(r.getEntityTitle(entity))
		if strings.Contains(title, queryLower) {
			results = append(results, entity)
			continue
		}

		// Search in tags if entity has them
		if tags := r.getEntityTags(entity); len(tags) > 0 {
			for _, tag := range tags {
				if strings.Contains(strings.ToLower(tag), queryLower) {
					results = append(results, entity)
					break
				}
			}
		}
	}

	return results, nil
}

// Exists checks if an entity with the given ID exists.
func (r *FilesystemRepository[T]) Exists(id core.EntityID) (bool, error) {
	if id == "" {
		return false, errors.New("id cannot be empty")
	}

	pattern := filepath.Join(r.basePath, fmt.Sprintf("%s-*.md", id))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return len(matches) > 0, nil
}

// Helper functions

// findFileByID finds a file matching the given entity ID.
func (r *FilesystemRepository[T]) findFileByID(id core.EntityID) (string, error) {
	// Pattern matches: {id}-{slug}.md (e.g., "skill-001-python.md")
	pattern := filepath.Join(r.basePath, fmt.Sprintf("%s-*.md", id))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for file: %w", err)
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("entity with ID %s not found", id)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("multiple files found for ID %s", id)
	}

	return matches[0], nil
}

// parseEntityFromFile reads a file and parses it into an entity.
func (r *FilesystemRepository[T]) parseEntityFromFile(filePath string, includeBody bool) (*T, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	frontmatter, body, err := ParseFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Marshal frontmatter back to YAML then unmarshal to entity type
	yamlBytes, err := yaml.Marshal(frontmatter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	var entity T
	if err := yaml.Unmarshal(yamlBytes, &entity); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entity: %w", err)
	}

	// Set body if requested and entity has a Body field
	if includeBody && body != "" {
		r.setEntityBody(&entity, body)
	}

	return &entity, nil
}

// serializeEntity converts an entity to markdown with YAML frontmatter.
func (r *FilesystemRepository[T]) serializeEntity(entity *T) ([]byte, error) {
	// Extract body if present
	body := r.getEntityBody(entity)

	// Serialize entity to YAML frontmatter
	content, err := SerializeFrontmatter(entity, body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// generateFileName creates a filename for an entity.
// Format: {id}-{slug}.md (e.g., "skill-001-python.md")
// Note: The ID already contains the entity type prefix (e.g., "skill-001")
func (r *FilesystemRepository[T]) generateFileName(id core.EntityID, title string) string {
	slug := slugify(title)
	if slug == "" {
		slug = "untitled"
	}
	return fmt.Sprintf("%s-%s.md", id, slug)
}

// getEntityID extracts the ID field from an entity using reflection.
func (r *FilesystemRepository[T]) getEntityID(entity *T) (core.EntityID, error) {
	v := reflect.ValueOf(entity).Elem()
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", errors.New("entity does not have an ID field")
	}

	id, ok := idField.Interface().(core.EntityID)
	if !ok {
		return "", errors.New("ID field is not of type EntityID")
	}

	if id == "" {
		return "", errors.New("entity ID is empty")
	}

	return id, nil
}

// getEntityTitle extracts the Title field from an entity using reflection.
func (r *FilesystemRepository[T]) getEntityTitle(entity *T) string {
	v := reflect.ValueOf(entity).Elem()
	titleField := v.FieldByName("Title")
	if !titleField.IsValid() {
		return ""
	}

	title, ok := titleField.Interface().(string)
	if !ok {
		return ""
	}

	return title
}

// getEntityBody extracts the Body field from an entity using reflection.
func (r *FilesystemRepository[T]) getEntityBody(entity *T) string {
	v := reflect.ValueOf(entity).Elem()
	bodyField := v.FieldByName("Body")
	if !bodyField.IsValid() {
		return ""
	}

	body, ok := bodyField.Interface().(string)
	if !ok {
		return ""
	}

	return body
}

// setEntityBody sets the Body field of an entity using reflection.
func (r *FilesystemRepository[T]) setEntityBody(entity *T, body string) {
	v := reflect.ValueOf(entity).Elem()
	bodyField := v.FieldByName("Body")
	if !bodyField.IsValid() || !bodyField.CanSet() {
		return
	}

	bodyField.SetString(body)
}

// getEntityTags extracts the Tags field from an entity using reflection.
func (r *FilesystemRepository[T]) getEntityTags(entity *T) []string {
	v := reflect.ValueOf(entity).Elem()
	tagsField := v.FieldByName("Tags")
	if !tagsField.IsValid() {
		return nil
	}

	tags, ok := tagsField.Interface().([]string)
	if !ok {
		return nil
	}

	return tags
}

// slugify converts a string to a URL-safe slug.
func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove duplicate hyphens
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	// Limit length to 50 characters
	if len(s) > 50 {
		s = s[:50]
		s = strings.TrimRight(s, "-")
	}

	return s
}
