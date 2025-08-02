package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "GoForge is a powerful CLI for scaffolding and managing Go projects.",
	Long: `GoForge is a NestJS-inspired command-line interface tool that helps you
initialize, develop, and maintain your Go applications.

It provides a robust, scalable architecture out-of-the-box, allowing you
to focus on business logic instead of setup and configuration.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err!= nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init registers all subcommands to the root command. This function is called
// automatically by Go when the package is initialized.
func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(buildCmd)
}
