// internal/scaffold/templates/components/repository.go.tpl
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"{{.ModulePath}}/internal/domain"
	"{{.ModulePath}}/internal/ports"
)

// ensure {{.NameTitle}}Repository implements the port at compile time.
var _ ports.{{.NameTitle}}Repository = (*{{.NameTitle}}Repository)(nil)

// {{.NameTitle}}Repository handles database operations for {{.Name}} entities.
type {{.NameTitle}}Repository struct {
	pool *pgxpool.Pool
}

// New{{.NameTitle}}Repository creates a new {{.NameTitle}}Repository.
func New{{.NameTitle}}Repository(pool *pgxpool.Pool) *{{.NameTitle}}Repository {
	return &{{.NameTitle}}Repository{pool: pool}
}

// FindByID retrieves a {{.Name}} by ID.
func (r *{{.NameTitle}}Repository) FindByID(ctx context.Context, id int64) (*domain.{{.NameTitle}}, error) {
	{{.Name}} := &domain.{{.NameTitle}}{}
	query := "SELECT id, created_at, updated_at FROM {{.Name | pluralize}} WHERE id = $1"

	err := r.pool.QueryRow(ctx, query, id).Scan(&{{.Name}}.ID, &{{.Name}}.CreatedAt, &{{.Name}}.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("{{.Name}} not found")
		}
		return nil, err
	}
	return {{.Name}}, nil
}

// Create inserts a new {{.Name}} into the database.
func (r *{{.NameTitle}}Repository) Create(ctx context.Context, {{.Name}} *domain.{{.NameTitle}}) error {
	query := `INSERT INTO {{.Name | pluralize}} (created_at, updated_at) 
			  VALUES (NOW(), NOW()) 
			  RETURNING id, created_at, updated_at`
	
	err := r.pool.QueryRow(ctx, query).Scan(&{{.Name}}.ID, &{{.Name}}.CreatedAt, &{{.Name}}.UpdatedAt)
	return err
}

// Update modifies an existing {{.Name}} in the database.
func (r *{{.NameTitle}}Repository) Update(ctx context.Context, {{.Name}} *domain.{{.NameTitle}}) error {
	query := `UPDATE {{.Name | pluralize}} 
			  SET updated_at = NOW() 
			  WHERE id = $1 
			  RETURNING updated_at`
	
	err := r.pool.QueryRow(ctx, query, {{.Name}}.ID).Scan(&{{.Name}}.UpdatedAt)
	return err
}

// Delete removes a {{.Name}} from the database.
func (r *{{.NameTitle}}Repository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM {{.Name | pluralize}} WHERE id = $1"
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// List retrieves multiple {{.Name | pluralize}} with pagination.
func (r *{{.NameTitle}}Repository) List(ctx context.Context, limit, offset int) ([]*domain.{{.NameTitle}}, error) {
	query := `SELECT id, created_at, updated_at 
			  FROM {{.Name | pluralize}} 
			  ORDER BY created_at DESC 
			  LIMIT $1 OFFSET $2`
	
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var {{.Name | pluralize}} []*domain.{{.NameTitle}}
	for rows.Next() {
		{{.Name}} := &domain.{{.NameTitle}}{}
		err := rows.Scan(&{{.Name}}.ID, &{{.Name}}.CreatedAt, &{{.Name}}.UpdatedAt)
		if err != nil {
			return nil, err
		}
		{{.Name | pluralize}} = append({{.Name | pluralize}}, {{.Name}})
	}

	return {{.Name | pluralize}}, rows.Err()
}