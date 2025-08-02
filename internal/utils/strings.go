package utils

import (
	"github.com/iancoleman/strcase"
)

// ToPascalCase converts a string to PascalCase.
// Example: "user-details" -> "UserDetails"
func ToPascalCase(s string) string {
	return strcase.ToCamel(s)
}

// ToCamelCase converts a string to camelCase.
// Example: "user-details" -> "userDetails"
func ToCamelCase(s string) string {
	return strcase.ToLowerCamel(s)
}

// ToSnakeCase converts a string to snake_case.
// Example: "userDetails" -> "user_details"
func ToSnakeCase(s string) string {
	return strcase.ToSnake(s)
}
