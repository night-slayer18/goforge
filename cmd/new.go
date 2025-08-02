package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/spf13/cobra"
)

// detectGoVersion runs "go version" and extracts the Go version (e.g., "1.24.5").
func detectGoVersion() (string, error) {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}
	// Expected output: "go version go1.24.5 darwin/amd64"
	parts := strings.Fields(string(out))
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected output from go version: %s", out)
	}
	return strings.TrimPrefix(parts[2], "go"), nil
}

// newCmd represents the 'new' command, responsible for creating new projects.
var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new Go project with a scalable architecture",
	Long: `The 'new' command creates a new directory with the specified project name,
and scaffolds a complete Go application based on Clean Architecture principles.

It sets up the entire project structure, including handlers, services, repositories,
a go.mod file, and a goforge.yml project manifest.`,
	// Enforce that exactly one argument (the project name) is provided.
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		modulePath, _ := cmd.Flags().GetString("module-path")

		if modulePath == "" {
			// Default module path if not provided by the user.
			modulePath = projectName
		}

		// Get absolute path for the new project directory.
		destPath, err := filepath.Abs(projectName)
		if err != nil {
			return fmt.Errorf("could not determine absolute path for project: %w", err)
		}

		goVersion, err := detectGoVersion()
		if err != nil {
			return fmt.Errorf("failed to detect go version: %w", err)
		}

		fmt.Printf("🚀 Creating new project '%s' at %s...\n", projectName, destPath)

		// Delegate the core scaffolding logic to the internal/scaffold package.
		err = scaffold.CreateProject(projectName, modulePath, goVersion, destPath)
		if err != nil {
			// Attempt to clean up the created directory on failure to avoid leaving partial projects.
			os.RemoveAll(destPath)
			return fmt.Errorf("failed to create project: %w", err)
		}

		fmt.Println("\n✅ Project created successfully!")
		fmt.Printf("\nNavigate to your new project:\n  cd %s\n", projectName)
		fmt.Println("\nTo run the development server:")
		fmt.Println("  goforge run dev")
		fmt.Println("\nHappy coding!")

		return nil
	},
}

func init() {
	// Define a flag for the 'new' command to allow users to specify a custom module path.
	newCmd.Flags().StringP("module-path", "m", "", "Explicitly set the Go module path (e.g., github.com/user/repo)")
}
