package cmd

import (
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// modelCmd represents the command to generate a domain model.
var modelCmd = &cobra.Command{
	Use:     "model <n>",
	Short:   "Generate a new domain model",
	Aliases: []string{"mod"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return scaffold.GenerateComponent("model", name)
	},
}