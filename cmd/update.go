package cmd

import (
	"fmt"

	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [module-path]",
	Short: "Update project dependencies",
	Long: `Update dependencies in your project. If a specific module is provided,
only that module will be updated. Otherwise, all dependencies will be updated.

Examples:
  goforge update                           # Update all dependencies
  goforge update github.com/gin-gonic/gin # Update specific dependency`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger.SetVerbose(verbose)

		cfg, projectRoot, err := project.LoadConfig()
		if err != nil {
			return fmt.Errorf("command must be run from a goforge project: %w", err)
		}

		if len(args) > 0 {
			// Update specific module
			modulePath := args[0]
			return updateSpecificModule(projectRoot, modulePath)
		}

		// Update all dependencies
		return updateAllDependencies(projectRoot, cfg)
	},
}

func updateSpecificModule(projectRoot, modulePath string) error {
	logger.Info("🔄 Updating dependency: %s", modulePath)

	if err := runner.InstallDependency(projectRoot, modulePath+"@latest"); err != nil {
		return fmt.Errorf("failed to update module: %w", err)
	}

	logger.Success("✅ Successfully updated: %s", modulePath)
	return nil
}

func updateAllDependencies(projectRoot string, cfg *project.Config) error {
	logger.Info("🔄 Updating all dependencies...")

	if len(cfg.Dependencies) == 0 {
		logger.Info("No dependencies to update")
		return nil
	}

	logger.Info("📦 Found %d dependencies to update", len(cfg.Dependencies))

	for module := range cfg.Dependencies {
		logger.Info("  Updating %s...", module)
		if err := runner.InstallDependency(projectRoot, module+"@latest"); err != nil {
			logger.Error("  ❌ Failed to update %s: %v", module, err)
			continue
		}
		logger.Success("  ✅ Updated %s", module)
	}

	// Run go mod tidy to clean up
	logger.Info("🧹 Cleaning up module files...")
	if err := runner.TidyGoModule(projectRoot); err != nil {
		return fmt.Errorf("failed to tidy module: %w", err)
	}

	logger.Success("✅ All dependencies updated successfully")
	return nil
}

func init() {
	updateCmd.Flags().BoolP("verbose", "v", false, "Enable verbose logging")
}