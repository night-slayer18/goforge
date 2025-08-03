package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch [script-name]",
	Short: "Watch for file changes and restart the application",
	Long: `The watch command monitors your project files for changes and automatically
restarts the specified script when changes are detected. If no script is specified,
it defaults to the 'dev' script.

This is useful for development workflows where you want automatic reloading.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up logging
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger.SetVerbose(verbose)

		cfg, projectRoot, err := project.LoadConfig()
		if err != nil {
			return fmt.Errorf("command must be run from a goforge project: %w", err)
		}

		// Determine script to run
		scriptName := "dev"
		if len(args) > 0 {
			scriptName = args[0]
		}

		script, exists := cfg.Scripts[scriptName]
		if !exists {
			return fmt.Errorf("script '%s' not found in goforge.yml\n\nAvailable scripts:\n%s", 
				scriptName, formatAvailableScripts(cfg.Scripts))
		}

		logger.Info("👀 Starting watch mode for script: %s", scriptName)
		logger.Info("📝 Command: %s", script)

		// Set up graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start watch mode
		watchMode := runner.NewWatchMode(projectRoot, "sh", "-c", script)
		if err := watchMode.Start(); err != nil {
			return fmt.Errorf("failed to start watch mode: %w", err)
		}

		// Wait for shutdown signal
		<-sigChan
		logger.Info("\n🛑 Shutting down...")

		if err := watchMode.Stop(); err != nil {
			logger.Error("Error stopping watch mode: %v", err)
			return err
		}

		logger.Info("✅ Watch mode stopped")
		return nil
	},
}

func formatAvailableScripts(scripts map[string]string) string {
	if len(scripts) == 0 {
		return "  No scripts defined"
	}

	result := ""
	for name, command := range scripts {
		result += fmt.Sprintf("  %s: %s\n", name, command)
	}
	return result
}

func init() {
	watchCmd.Flags().BoolP("verbose", "v", false, "Enable verbose logging")
}