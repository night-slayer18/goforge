package cmd

import (
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// handlerCmd represents the command to generate an HTTP handler.
var handlerCmd = &cobra.Command{
	Use:     "handler <name>",
	Short:   "Generate a new HTTP handler",
	Aliases: []string{"h"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return scaffold.GenerateComponent("handler", name)
	},
}