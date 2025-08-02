package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"{{.ModuleName}}/internal/domain"
	"{{.ModuleName}}/internal/ports"
)

// ensure PostgresUserRepository implements the port at compile time.
var _ ports.UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) FindByID(id int64) (*domain.User, error) {
	user := &domain.User{}
	query := "SELECT id, email, name FROM users WHERE id = $1"

	err := r.pool.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found") // Or a custom error type
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	query := "INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id"
	err := r.pool.QueryRow(context.Background(), query, user.Email, user.Name).Scan(&user.ID)
	return err
}
func (r *PostgresUserRepository) Update(user *domain.User) error {
    query := "UPDATE users SET email = $1, name = $2 WHERE id = $3"
    _, err := r.pool.Exec(context.Background(), query, user.Email, user.Name, user.ID)
    return err
}
func (r *PostgresUserRepository) Delete(id int64) error {
    query := "DELETE FROM users WHERE id = $1"
    _, err := r.pool.Exec(context.Background(), query, id)
    return err
}
