package cmd

import (
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// repositoryCmd represents the command to generate a repository.
var repositoryCmd = &cobra.Command{
	Use:     "repository <name>",
	Short:   "Generate a new repository",
	Aliases: []string{"repo", "r"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return scaffold.GenerateComponent("repository", name)
	},
}