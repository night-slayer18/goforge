package cmd

import (
	"fmt"

	"github.com/iancoleman/strcase"
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
		err := scaffold.GenerateComponent("handler", name)
		if err != nil {
			return err
		}
		fileName := strcase.ToSnake(name)
		fmt.Printf("✅ Generated handler: internal/adapters/http/handler/%s_handler.go\n", fileName)
		fmt.Println("Don't forget to register it in your router in cmd/server/main.go!")
		return nil
	},
}
