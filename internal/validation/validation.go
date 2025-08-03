// internal/validation/validator.go
package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Valid project name pattern: alphanumeric, hyphens, underscores
	projectNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	
	// Reserved Go keywords and common conflicts
	reservedNames = map[string]bool{
		"break": true, "case": true, "chan": true, "const": true, "continue": true,
		"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
		"func": true, "go": true, "goto": true, "if": true, "import": true,
		"interface": true, "map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true, "var": true,
		// Common directory names that might cause conflicts
		"vendor": true, "node_modules": true, "dist": true, "build": true,
		"main": true, "test": true, "tests": true,
	}
)

// ValidationError represents a validation error with suggestions
type ValidationError struct {
	Field       string
	Value       string
	Message     string
	Suggestions []string
}

func (e *ValidationError) Error() string {
	msg := fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
	if len(e.Suggestions) > 0 {
		msg += "\n\nSuggestions:\n"
		for _, suggestion := range e.Suggestions {
			msg += fmt.Sprintf("  - %s\n", suggestion)
		}
	}
	return msg
}

// ProjectValidator handles project-level validation
type ProjectValidator struct{}

func NewProjectValidator() *ProjectValidator {
	return &ProjectValidator{}
}

// ValidateProjectName validates the project name and provides suggestions
func (v *ProjectValidator) ValidateProjectName(name string) error {
	if name == "" {
		return &ValidationError{
			Field:   "project_name",
			Value:   name,
			Message: "project name cannot be empty",
			Suggestions: []string{
				"Use a descriptive name like 'my-api' or 'user-service'",
				"Project names should start with a letter",
			},
		}
	}

	if len(name) > 100 {
		return &ValidationError{
			Field:   "project_name",
			Value:   name,
			Message: "project name is too long (max 100 characters)",
			Suggestions: []string{
				"Use a shorter, more concise name",
				fmt.Sprintf("Consider: %s", suggestShorterName(name)),
			},
		}
	}

	if !projectNamePattern.MatchString(name) {
		return &ValidationError{
			Field:   "project_name",
			Value:   name,
			Message: "project name contains invalid characters",
			Suggestions: []string{
				"Use only letters, numbers, hyphens (-), and underscores (_)",
				"Start with a letter",
				fmt.Sprintf("Consider: %s", sanitizeProjectName(name)),
			},
		}
	}

	if reservedNames[strings.ToLower(name)] {
		return &ValidationError{
			Field:   "project_name",
			Value:   name,
			Message: "project name conflicts with reserved keywords or common directory names",
			Suggestions: []string{
				fmt.Sprintf("Try: %s-app, %s-service, or my-%s", name, name, name),
				"Avoid Go keywords and common directory names",
			},
		}
	}

	return nil
}

// ValidateModulePath validates Go module paths
func (v *ProjectValidator) ValidateModulePath(modulePath string) error {
	if modulePath == "" {
		return &ValidationError{
			Field:   "module_path",
			Value:   modulePath,
			Message: "module path cannot be empty",
			Suggestions: []string{
				"Use your domain: github.com/username/project",
				"Or a simple name: myproject",
			},
		}
	}

	// Check for valid module path format
	if strings.Contains(modulePath, " ") {
		return &ValidationError{
			Field:   "module_path",
			Value:   modulePath,
			Message: "module path cannot contain spaces",
			Suggestions: []string{
				fmt.Sprintf("Try: %s", strings.ReplaceAll(modulePath, " ", "-")),
			},
		}
	}

	return nil
}

// ValidateComponentName validates component names for generation
func (v *ProjectValidator) ValidateComponentName(componentType, name string) error {
	if name == "" {
		return &ValidationError{
			Field:   fmt.Sprintf("%s_name", componentType),
			Value:   name,
			Message: fmt.Sprintf("%s name cannot be empty", componentType),
			Suggestions: []string{
				fmt.Sprintf("Provide a descriptive name for your %s", componentType),
				"Use PascalCase or camelCase naming conventions",
			},
		}
	}

	// Check for valid Go identifier
	if !isValidGoIdentifier(name) {
		return &ValidationError{
			Field:   fmt.Sprintf("%s_name", componentType),
			Value:   name,
			Message: fmt.Sprintf("%s name is not a valid Go identifier", componentType),
			Suggestions: []string{
				"Use letters and numbers only",
				"Start with a letter",
				fmt.Sprintf("Consider: %s", sanitizeIdentifier(name)),
			},
		}
	}

	return nil
}

// Helper functions
func suggestShorterName(name string) string {
	if len(name) <= 20 {
		return name
	}
	
	// Extract meaningful parts
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	
	if len(parts) > 1 {
		// Use first two meaningful parts
		return strings.Join(parts[:2], "-")
	}
	
	// Truncate and add suffix
	return name[:15] + "-app"
}

func sanitizeProjectName(name string) string {
	var result strings.Builder
	
	for i, r := range name {
		if i == 0 && !unicode.IsLetter(r) {
			result.WriteRune('a')
		}
		
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else if r == ' ' || r == '.' {
			result.WriteRune('-')
		} else if r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}
	
	sanitized := result.String()
	if sanitized == "" {
		return "my-project"
	}
	
	return sanitized
}

func isValidGoIdentifier(name string) bool {
	if name == "" {
		return false
	}
	
	for i, r := range name {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	
	return true
}

func sanitizeIdentifier(name string) string {
	var result strings.Builder
	
	for i, r := range name {
		if i == 0 {
			if unicode.IsLetter(r) || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('A')
				if unicode.IsDigit(r) {
					result.WriteRune(r)
				}
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				result.WriteRune(r)
			}
		}
	}
	
	sanitized := result.String()
	if sanitized == "" {
		return "Component"
	}
	
	return sanitized
}