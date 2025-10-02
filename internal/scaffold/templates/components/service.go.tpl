package service

import (
	// "{{.ModulePath}}/internal/ports" // TODO: Uncomment when you add repository dependencies.
)

// {{.NameTitle}}Service provides application logic for the {{.Name}} resource.
// It depends on interfaces (ports) defined in the ports package, not on
// concrete database implementations.
type {{.NameTitle}}Service struct {
	// TODO: Add dependencies like repositories here.
	// For example:
	// userRepo ports.UserRepository
}

// New{{.NameTitle}}Service is a factory function that creates a new {{.NameTitle}}Service.
func New{{.NameTitle}}Service(/* e.g., userRepo ports.UserRepository */) *{{.NameTitle}}Service {
	return &{{.NameTitle}}Service{
		// userRepo: userRepo,
	}
}

// DoSomething is an example service method.
// TODO: Rename and implement your business logic.
func (s *{{.NameTitle}}Service) DoSomething(/* parameters like id, or a data struct */) error {
	// 1. Perform validation on the input data.
	// 2. Interact with repositories to fetch or save data.
	// 3. Return the result or an error.
	return nil
}
