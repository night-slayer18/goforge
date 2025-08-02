package cmd

import (
	"github.com/spf13/cobra"
)

// generateCmd represents the base for all 'generate' subcommands (e.g., 'goforge generate handler').
// It doesn't have its own Run function because it only serves as a parent for other commands.
var generateCmd = &cobra.Command{
	Use:     "generate <component> <name>",
	Short:   "Generate a new application component",
	Long:    `The 'generate' command (alias 'g') creates boilerplate files for various components of your application.`,
	Aliases: []string{"g"},
}

func init() {
	// Register component-specific generation commands as subcommands of 'generate'.
	generateCmd.AddCommand(handlerCmd)
	generateCmd.AddCommand(serviceCmd)
}
