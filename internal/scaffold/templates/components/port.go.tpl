package ports

import (
	"context"
	"{{.ModulePath}}/internal/domain"
)

// {{.NameTitle}}Repository defines the contract for {{.Name}} data access.
type {{.NameTitle}}Repository interface {
	// FindByID retrieves a {{.Name}} by its ID.
	FindByID(ctx context.Context, id int64) (*domain.{{.NameTitle}}, error)
	
	// Create inserts a new {{.Name}} into the repository.
	Create(ctx context.Context, {{.Name}} *domain.{{.NameTitle}}) error
	
	// Update modifies an existing {{.Name}} in the repository.
	Update(ctx context.Context, {{.Name}} *domain.{{.NameTitle}}) error
	
	// Delete removes a {{.Name}} from the repository.
	Delete(ctx context.Context, id int64) error
	
	// List retrieves multiple {{.Name | pluralize}} with pagination.
	List(ctx context.Context, limit, offset int) ([]*domain.{{.NameTitle}}, error)
	
	// TODO: Add domain-specific repository methods
	// Example:
	// FindByEmail(ctx context.Context, email string) (*domain.{{.NameTitle}}, error)
	// FindActive(ctx context.Context) ([]*domain.{{.NameTitle}}, error)
}

// {{.NameTitle}}Service defines the contract for {{.Name}} business logic.
type {{.NameTitle}}Service interface {
	// Get{{.NameTitle}} retrieves a {{.Name}} by ID.
	Get{{.NameTitle}}(ctx context.Context, id int64) (*domain.{{.NameTitle}}, error)
	
	// Create{{.NameTitle}} creates a new {{.Name}}.
	Create{{.NameTitle}}(ctx context.Context, {{.Name}} *domain.{{.NameTitle}}) error
	
	// Update{{.NameTitle}} updates an existing {{.Name}}.
	Update{{.NameTitle}}(ctx context.Context, {{.Name}} *domain.{{.NameTitle}}) error
	
	// Delete{{.NameTitle}} removes a {{.Name}}.
	Delete{{.NameTitle}}(ctx context.Context, id int64) error
	
	// List{{.NameTitle | pluralize}} retrieves multiple {{.Name | pluralize}}.
	List{{.NameTitle | pluralize}}(ctx context.Context, limit, offset int) ([]*domain.{{.NameTitle}}, error)
	
	// TODO: Add business logic methods
	// Example:
	// Activate{{.NameTitle}}(ctx context.Context, id int64) error
	// Deactivate{{.NameTitle}}(ctx context.Context, id int64) error
}