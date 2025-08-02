package ports

import "{{.ModuleName}}/internal/domain"

// UserRepository defines the contract for database operations on User objects.
type UserRepository interface {
	FindByID(id int64) (*domain.User, error)
	Create(user *domain.User) error
}
