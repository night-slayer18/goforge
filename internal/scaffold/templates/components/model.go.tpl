// internal/scaffold/templates/components/model.go.tpl
package domain

import (
	"time"
)

// {{.NameTitle}} represents a {{.Name}} entity in the domain.
type {{.NameTitle}} struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// TODO: Add your domain-specific fields here
	// Example:
	// Name        string `json:"name" db:"name" validate:"required,min=1,max=100"`
	// Email       string `json:"email" db:"email" validate:"required,email"`
	// IsActive    bool   `json:"is_active" db:"is_active"`
}

// TableName returns the database table name for this entity.
func ({{.Name}} *{{.NameTitle}}) TableName() string {
	return "{{.Name | pluralize}}"
}

// Validate performs domain-level validation on the {{.Name}} entity.
func ({{.Name}} *{{.NameTitle}}) Validate() error {
	// TODO: Implement domain validation logic
	// Example:
	// if {{.Name}}.Name == "" {
	//     return errors.New("name is required")
	// }
	return nil
}

// IsNew returns true if this is a new entity (not persisted yet).
func ({{.Name}} *{{.NameTitle}}) IsNew() bool {
	return {{.Name}}.ID == 0
}

// Update updates the entity's timestamp.
func ({{.Name}} *{{.NameTitle}}) Update() {
	{{.Name}}.UpdatedAt = time.Now()
}