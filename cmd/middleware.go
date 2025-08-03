package cmd

import (
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// middlewareCmd represents the command to generate a middleware.
var middlewareCmd = &cobra.Command{
	Use:     "middleware <name>",
	Short:   "Generate a new HTTP middleware",
	Aliases: []string{"m"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return scaffold.GenerateComponent("middleware", name)
	},
}