package cmd

import (
	"fmt"
	"strings"

	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/spf13/cobra"
)

// addCmd represents the command to add a new Go module dependency.
var addCmd = &cobra.Command{
	Use:   "add <module-path>[@version]",
	Short: "Add a new dependency to the project",
	Long: `Downloads the specified module using 'go get' and adds it to the
'dependencies' section of your goforge.yml file for declarative dependency management.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modulePath := args[0]

		cfg, projectRoot, err := project.LoadConfig()
		if err!= nil {
			return err
		}

		fmt.Printf("📦 Adding dependency: %s\n", modulePath)
		// Execute 'go get' to download the dependency and update go.mod/go.sum.
		err = runner.ExecuteCommand(projectRoot, "go", "get", modulePath)
		if err!= nil {
			return fmt.Errorf("failed to 'go get' module: %w", err)
		}

		// Extract module base path and version for goforge.yml.
		var moduleName, version string
		if i := strings.LastIndex(modulePath, "@"); i != -1 {
			moduleName = modulePath[:i]
			version = modulePath[i+1:]
		} else {
			moduleName = modulePath
			version = "latest"
		}

		if cfg.Dependencies == nil {
			cfg.Dependencies = make(map[string]string)
		}
		cfg.Dependencies[moduleName] = version

		// Save the updated configuration back to goforge.yml.
		err = project.SaveConfig(projectRoot, cfg)
		if err!= nil {
			return fmt.Errorf("failed to update goforge.yml: %w", err)
		}

		// Tidy modules to ensure consistency
		if err := runner.TidyGoModule(projectRoot); err != nil {
			logger.Warn("Failed to tidy go modules: %v", err)
		}

		fmt.Printf("✅ Successfully added '%s' and updated goforge.yml.\n", modulePath)
		return nil
	},
}
