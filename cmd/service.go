package cmd

import (
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// serviceCmd represents the command to generate an application service.
var serviceCmd = &cobra.Command{
	Use:     "service <n>",
	Short:   "Generate a new application service",
	Aliases: []string{"s"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return scaffold.GenerateComponent("service", name)
	},
}