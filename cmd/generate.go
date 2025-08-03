package cmd

import (
	"github.com/spf13/cobra"
)

// generateCmd represents the base for all 'generate' subcommands (e.g., 'goforge generate handler').
// It doesn't have its own Run function because it only serves as a parent for other commands.
var generateCmd = &cobra.Command{
	Use:     "generate <component> <n>",
	Short:   "Generate a new application component",
	Long: `The 'generate' command (alias 'g') creates boilerplate files for various components of your application.

Available components:
  handler     Generate HTTP handlers for API endpoints
  service     Generate application services for business logic  
  repository  Generate repository implementations for data access
  model       Generate domain models/entities
  middleware  Generate HTTP middleware components
  port        Generate port interfaces for clean architecture

Examples:
  goforge generate handler user
  goforge g service auth
  goforge g repository product
  goforge g model order
  goforge g middleware cors
  goforge g port notification`,
	Aliases: []string{"g"},
}

func init() {
	// Register all component-specific generation commands as subcommands of 'generate'.
	generateCmd.AddCommand(handlerCmd)
	generateCmd.AddCommand(serviceCmd)
	generateCmd.AddCommand(repositoryCmd)
	generateCmd.AddCommand(modelCmd)
	generateCmd.AddCommand(middlewareCmd)
	generateCmd.AddCommand(portCmd)
}