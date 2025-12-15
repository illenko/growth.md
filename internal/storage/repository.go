package storage

import "github.com/illenko/growth.md/internal/core"

// Repository defines the interface for entity persistence operations.
// It uses Go generics to work with any entity type.
type Repository[T any] interface {
	// Create persists a new entity to storage.
	// Returns an error if the entity already exists or if persistence fails.
	Create(entity *T) error

	// GetByID retrieves an entity by its ID (metadata only, without body).
	// Returns an error if the entity is not found.
	GetByID(id core.EntityID) (*T, error)

	// GetByIDWithBody retrieves an entity by its ID including the markdown body.
	// Returns an error if the entity is not found.
	GetByIDWithBody(id core.EntityID) (*T, error)

	// GetAll retrieves all entities of this type (metadata only, without bodies).
	// Returns an empty slice if no entities exist.
	GetAll() ([]*T, error)

	// Update persists changes to an existing entity.
	// Returns an error if the entity does not exist or if persistence fails.
	Update(entity *T) error

	// Delete removes an entity from storage by its ID.
	// Returns an error if the entity does not exist or if deletion fails.
	Delete(id core.EntityID) error

	// Search finds entities matching the given query string.
	// The query searches in titles, tags, and other text fields.
	// Returns entities with metadata only (no bodies).
	Search(query string) ([]*T, error)

	// Exists checks if an entity with the given ID exists.
	Exists(id core.EntityID) (bool, error)
}
