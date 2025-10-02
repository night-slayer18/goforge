package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time
var version = "dev"
// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "GoForge is a powerful CLI for scaffolding and managing Go projects.",
	Long: `GoForge is a NestJS-inspired command-line interface tool that helps you
initialize, develop, and maintain your Go applications.
It provides a robust, scalable architecture out-of-the-box, allowing you
to focus on business logic instead of setup and configuration.

Interactive Mode:
  GoForge supports both traditional command-line and interactive modes.
  Use --interactive flag or run commands without arguments to access
  the guided interactive experience.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err!= nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{printf "GoForge CLI Version: %s\n" .Version}}`)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(watchCmd)    
	rootCmd.AddCommand(updateCmd) 
	rootCmd.AddCommand(cleanCmd) 
	
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
}