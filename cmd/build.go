package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/spf13/cobra"
)

// buildCmd represents the command to build the user's application.
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Go application binary",
	Long: `Compiles the Go application into an executable binary and copies any assets
specified in the 'build.assets' section of goforge.yml to the output directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, projectRoot, err := project.LoadConfig()
		if err!= nil {
			return err
		}

		outputDir := filepath.Join(projectRoot, "dist")
		binaryName := cfg.ProjectName
		if binaryName == "" {
			binaryName = filepath.Base(projectRoot)
		}
		outputPath := filepath.Join(outputDir, binaryName)

		fmt.Printf("🏗️  Building project '%s'...\n", cfg.ProjectName)

		// Ensure output directory exists.
		if err := os.MkdirAll(outputDir, os.ModePerm); err!= nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Build the binary.
		err = runner.ExecuteCommand(projectRoot, "go", "build", "-o", outputPath, "./cmd/server")
		if err!= nil {
			return fmt.Errorf("go build failed: %w", err)
		}
		fmt.Printf("✅ Binary created at: %s\n", outputPath)

		// Handle assets defined in goforge.yml.
		if cfg.Build!= nil && len(cfg.Build.Assets) > 0 {
			fmt.Println("📦 Copying assets...")
			// This is a simplified asset handler. A real implementation would
			// support glob patterns.
			for _, assetPath := range cfg.Build.Assets {
				sourcePath := filepath.Join(projectRoot, assetPath)
				destPath := filepath.Join(outputDir, assetPath)

				// Ensure destination directory exists.
				if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err!= nil {
					fmt.Printf("  - Could not create asset destination dir for %s: %v\n", assetPath, err)
					continue
				}

				// Using 'cp -r' is a simple approach for demonstration.
				err := runner.ExecuteCommand(projectRoot, "cp", "-r", sourcePath, destPath)
				if err!= nil {
					fmt.Printf("  - Failed to copy asset %s: %v\n", assetPath, err)
				} else {
					fmt.Printf("  - Copied: %s\n", assetPath)
				}
			}
		}

		fmt.Println("\n✨ Build complete.")
		return nil
	},
}
