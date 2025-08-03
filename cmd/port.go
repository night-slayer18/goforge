package cmd

import (
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// portCmd represents the command to generate a port interface.
var portCmd = &cobra.Command{
	Use:     "port <n>",
	Short:   "Generate a new port interface",
	Aliases: []string{"p"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return scaffold.GenerateComponent("port", name)
	},
}