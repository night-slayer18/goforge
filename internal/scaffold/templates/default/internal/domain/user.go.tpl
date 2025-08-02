// Package domain contains the core business models and logic.
// These are the most central and stable parts of the application.
package domain

// User represents a user in the system.
// This is a core business entity.
type User struct {
	ID    int64
	Email string
	Name  string
}
