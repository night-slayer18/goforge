package cmd

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// serviceCmd represents the command to generate an application service.
var serviceCmd = &cobra.Command{
	Use:     "service <name>",
	Short:   "Generate a new application service",
	Aliases: []string{"s"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		err := scaffold.GenerateComponent("service", name)
		if err!= nil {
			return err
		}
		fileName := strcase.ToSnake(name)
		fmt.Printf("✅ Generated service: internal/app/service/%s_service.go\n", fileName)
		fmt.Println("Remember to wire it up in your dependency injection setup in cmd/server/main.go.")
		return nil
	},
}
