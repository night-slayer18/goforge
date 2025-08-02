package service

import (
	"{{.ModuleName}}/internal/domain"
	"{{.ModuleName}}/internal/ports"
)

// UserService implements the application logic for user management.
type UserService struct {
	userRepo ports.UserRepository
}

// NewUserService is a factory function that creates a new UserService.
func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUser is a use case that retrieves a user by their ID.
func (s *UserService) GetUser(id int64) (*domain.User, error) {
	// Business logic would go here.
	return s.userRepo.FindByID(id)
}
