package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean build artifacts and temporary files",
	Long: `Remove build artifacts, temporary files, and caches from your project.

This includes:
  • dist/ directory (build outputs)
  • coverage files (*.out, coverage.html)
  • test cache
  • go module cache (with --all flag)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger.SetVerbose(verbose)

		_, projectRoot, err := project.LoadConfig()
		if err != nil {
			return fmt.Errorf("command must be run from a goforge project: %w", err)
		}

		all, _ := cmd.Flags().GetBool("all")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return cleanProject(projectRoot, all, dryRun)
	},
}

func cleanProject(projectRoot string, all, dryRun bool) error {
	logger.Info("🧹 Cleaning project...")

	filesToRemove := []string{
		"dist",
		"coverage.out",
		"coverage.html",
		"*.test",
	}

	var removed []string

	for _, pattern := range filesToRemove {
		matches, err := filepath.Glob(filepath.Join(projectRoot, pattern))
		if err != nil {
			logger.Debug("Error globbing %s: %v", pattern, err)
			continue
		}

		for _, match := range matches {
			relPath, _ := filepath.Rel(projectRoot, match)
			
			if dryRun {
				logger.Info("Would remove: %s", relPath)
				continue
			}

			if err := os.RemoveAll(match); err != nil {
				logger.Error("Failed to remove %s: %v", relPath, err)
				continue
			}

			removed = append(removed, relPath)
			logger.Debug("Removed: %s", relPath)
		}
	}

	if dryRun {
		logger.Info("🔍 Dry run completed - no files were actually removed")
		return nil
	}

	if len(removed) > 0 {
		logger.Success("✅ Removed %d items:", len(removed))
		for _, item := range removed {
			logger.Info("  • %s", item)
		}
	} else {
		logger.Info("✅ Project is already clean")
	}

	if all {
		logger.Info("🗑️  Cleaning Go module cache...")
		opts := runner.DefaultOptions()
		opts.Dir = projectRoot
		opts.ShowOutput = false

		if err := runner.ExecuteCommandWithOptions("go", []string{"clean", "-modcache"}, opts); err != nil {
			logger.Warn("Failed to clean module cache: %v", err)
		} else {
			logger.Success("✅ Module cache cleaned")
		}
	}

	return nil
}

func init() {
	cleanCmd.Flags().BoolP("all", "a", false, "Also clean Go module cache")
	cleanCmd.Flags().BoolP("dry-run", "n", false, "Show what would be removed without actually removing")
	cleanCmd.Flags().BoolP("verbose", "v", false, "Enable verbose logging")
}