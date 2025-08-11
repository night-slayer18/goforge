package cmd

import (
	"fmt"

	"github.com/night-slayer18/goforge/internal/interactive"
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// generateCmd represents the base for all 'generate' subcommands (e.g., 'goforge generate handler').
var generateCmd = &cobra.Command{
	Use:     "generate [component] [name]",
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
  goforge g port notification
  
  # Interactive mode
  goforge generate --interactive
  goforge g -i`,
	Aliases: []string{"g"},
	Args:    cobra.MaximumNArgs(2), // Allow 0, 1, or 2 args for interactive mode
	RunE: func(cmd *cobra.Command, args []string) error {
		interactiveFlag, _ := cmd.Flags().GetBool("interactive")
		
		var componentType, name string
		
		// Determine if we should use interactive mode
		useInteractive := false
		
		if interactiveFlag {
			// Explicitly requested interactive mode
			if !interactive.IsInteractiveTerminal() {
				return fmt.Errorf("interactive mode requested but not running in an interactive terminal")
			}
			useInteractive = true
		} else if len(args) == 0 && interactive.IsInteractiveTerminal() {
			// No arguments provided and we're in a terminal - offer interactive mode
			fmt.Print("No component specified. Use interactive mode? (Y/n): ")
			useInteractive = interactive.PromptForInteractiveMode()
		} else if len(args) < 2 {
			// Not enough arguments and not interactive - show error
			return fmt.Errorf("component type and name are required when not using interactive mode\n\nUsage:\n  goforge generate <component> <name>\n  goforge generate --interactive")
		}
		
		if useInteractive {
			// Use interactive mode
			wizard := interactive.NewComponentWizard()
			options, err := wizard.RunComponentCreationWizard()
			if err != nil {
				return fmt.Errorf("interactive session failed: %w", err)
			}
			
			componentType = options.Type
			name = options.Name
		} else {
			// Use traditional command-line mode
			componentType = args[0]
			name = args[1]
		}
		
		// Generate the component (same for both modes)
		return scaffold.GenerateComponent(componentType, name)
	},
}

func init() {
	// Add interactive flag to generate command
	generateCmd.Flags().BoolP("interactive", "i", false, 
		"Use interactive mode for component generation")
	
	// Register all component-specific generation commands as subcommands of 'generate'.
	generateCmd.AddCommand(handlerCmd)
	generateCmd.AddCommand(serviceCmd)
	generateCmd.AddCommand(repositoryCmd)
	generateCmd.AddCommand(modelCmd)
	generateCmd.AddCommand(middlewareCmd)
	generateCmd.AddCommand(portCmd)
}