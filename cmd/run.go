package cmd

import (
	"fmt"

	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/spf13/cobra"
)

// runCmd represents the command to execute a custom script from goforge.yml.
var runCmd = &cobra.Command{
	Use:   "run <script-name>",
	Short: "Run a custom script defined in goforge.yml",
	Long: `Executes a command from the 'scripts' section of your project's goforge.yml file.
This is analogous to 'npm run <script-name>' in the Node.js ecosystem.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptName := args[0]

		// Load the project configuration to find the scripts.
		cfg, projectRoot, err := project.LoadConfig()
		if err!= nil {
			return err
		}

		scriptCommand, exists := cfg.Scripts[scriptName]
		if!exists {
			return fmt.Errorf("script '%s' not found in goforge.yml", scriptName)
		}

		fmt.Printf("▶️  Running script '%s': %s\n\n", scriptName, scriptCommand)
		// Delegate execution to the runner package.
		return runner.ExecuteScript(projectRoot, scriptCommand)
	},
}
