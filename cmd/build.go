package cmd

import (
	"fmt"
	"io"
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
		if cfg.Build != nil && len(cfg.Build.Assets) > 0 {
			fmt.Println("📦 Copying assets...")
			for _, assetPath := range cfg.Build.Assets {
				sourcePath := filepath.Join(projectRoot, assetPath)

				info, err := os.Stat(sourcePath)
				if os.IsNotExist(err) {
					fmt.Printf("  - Asset not found, skipping: %s\n", assetPath)
					continue
				}
				if err != nil {
					fmt.Printf("  - Error accessing asset %s: %v\n", assetPath, err)
					continue
				}

				destPath := filepath.Join(outputDir, assetPath)

				if info.IsDir() {
					err = copyDir(sourcePath, destPath)
				} else {
					err = copyFile(sourcePath, destPath)
				}

				if err != nil {
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

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
}
