// Package ports defines the interfaces that connect the application core (app)
// to outside adapters (like databases or external APIs).
package ports

import "{{.ModuleName}}/internal/domain"

// UserRepository is an interface for database operations on User objects.
// The implementation will be in the 'adapters' layer.
type UserRepository interface {
	FindByID(id int64) (*domain.User, error)
	Create(user *domain.User) error
}
