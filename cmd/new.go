package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/scaffold"
	"github.com/night-slayer18/goforge/internal/validation"
	"github.com/spf13/cobra"
)

// detectGoVersion runs "go version" and extracts the Go version (e.g., "1.24.5").
func detectGoVersion() (string, error) {
	logger.Debug("Detecting Go version...")
	
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", fmt.Errorf("failed to detect Go version: %w\n\nPlease ensure Go is installed and available in your PATH", err)
	}
	
	// Expected output: "go version go1.24.5 darwin/amd64"
	parts := strings.Fields(string(out))
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected output from go version: %s", out)
	}
	
	version := strings.TrimPrefix(parts[2], "go")
	logger.Debug("Detected Go version: %s", version)
	return version, nil
}

// checkPrerequisites ensures all required tools are available
func checkPrerequisites() error {
	logger.Debug("Checking prerequisites...")
	
	// Check Go installation
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("Go is not installed or not in PATH. Please install Go from https://golang.org/dl/")
	}
	
	// Check Git installation (optional but recommended)
	if _, err := exec.LookPath("git"); err != nil {
		logger.Warn("Git is not installed. Repository initialization will be skipped.")
		logger.Info("💡 Install Git to enable automatic repository initialization")
	}
	
	logger.Debug("Prerequisites check completed")
	return nil
}

// checkDirectoryExists checks if the target directory already exists
func checkDirectoryExists(projectName string) error {
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("directory '%s' already exists\n\nSuggestions:\n  • Choose a different project name\n  • Remove the existing directory: rm -rf %s\n  • Use a different location", projectName, projectName)
	}
	return nil
}

// newCmd represents the 'new' command, responsible for creating new projects.
var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new Go project with a scalable architecture",
	Long: `The 'new' command creates a new directory with the specified project name,
and scaffolds a complete Go application based on Clean Architecture principles.

It sets up the entire project structure, including handlers, services, repositories,
a go.mod file, and a goforge.yml project manifest.

Examples:
  goforge new my-api
  goforge new user-service --module-path github.com/myorg/user-service
  goforge new blog-app -m gitlab.com/company/blog-app`,
	
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Set up logging based on verbose flag
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger.SetVerbose(verbose)
		
		return checkPrerequisites()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()
		projectName := args[0]
		
		// Get flags
		modulePath, _ := cmd.Flags().GetString("module-path")
		skipGit, _ := cmd.Flags().GetBool("skip-git")
		template, _ := cmd.Flags().GetString("template")
		
		// Initialize validator
		validator := validation.NewProjectValidator()
		
		// Validate project name
		if err := validator.ValidateProjectName(projectName); err != nil {
			if validationErr, ok := err.(*validation.ValidationError); ok {
				logger.ValidationError(validationErr.Field, validationErr.Value, validationErr.Message, validationErr.Suggestions)
				return fmt.Errorf("invalid project name")
			}
			return err
		}
		
		// Set default module path if not provided
		if modulePath == "" {
			modulePath = projectName
			logger.Debug("Using default module path: %s", modulePath)
		}
		
		// Validate module path
		if err := validator.ValidateModulePath(modulePath); err != nil {
			if validationErr, ok := err.(*validation.ValidationError); ok {
				logger.ValidationError(validationErr.Field, validationErr.Value, validationErr.Message, validationErr.Suggestions)
				return fmt.Errorf("invalid module path")
			}
			return err
		}
		
		// Check if directory already exists
		if err := checkDirectoryExists(projectName); err != nil {
			logger.Error("❌ %v", err)
			return fmt.Errorf("directory conflict")
		}
		
		// Get absolute path for the new project directory
		destPath, err := filepath.Abs(projectName)
		if err != nil {
			logger.Error("Failed to determine absolute path for project")
			return fmt.Errorf("could not determine absolute path for project: %w", err)
		}
		
		// Detect Go version
		goVersion, err := detectGoVersion()
		if err != nil {
			logger.Error("Failed to detect Go version")
			return err
		}
		
		// Start project creation
		logger.ProjectCreationStart(projectName)
		logger.Info("📍 Location: %s", destPath)
		logger.Info("📦 Module: %s", modulePath)
		logger.Info("🐹 Go version: %s", goVersion)
		if template != "default" {
			logger.Info("📋 Template: %s", template)
		}
		logger.Info("")
		
		// Create project structure
		logger.Step(1, 4, "Setting up project structure...")
		
		scaffoldOptions := scaffold.Options{
			ProjectName: projectName,
			ModulePath:  modulePath,
			GoVersion:   goVersion,
			DestPath:    destPath,
			Template:    template,
			SkipGit:     skipGit,
		}
		
		if err := scaffold.CreateProjectWithOptions(scaffoldOptions); err != nil {
			// Clean up on failure
			if _, statErr := os.Stat(destPath); statErr == nil {
				logger.Debug("Cleaning up failed project creation...")
				os.RemoveAll(destPath)
			}
			
			logger.Error("Failed to create project: %v", err)
			return fmt.Errorf("failed to create project: %w", err)
		}
		
		// Calculate total time
		duration := time.Since(startTime)
		logger.ProjectCreationComplete(projectName, duration)
		
		// Show additional information
		showPostCreationInfo(projectName, modulePath)
		
		return nil
	},
}

// showPostCreationInfo displays helpful information after project creation
func showPostCreationInfo(projectName, modulePath string) {
	logger.Info("📋 Project Information:")
	logger.Info("   Name: %s", projectName)
	logger.Info("   Module: %s", modulePath)
	logger.Info("   Location: %s", filepath.Join(".", projectName))
	logger.Info("")
	
	logger.Info("🛠️  Available Commands:")
	logger.Info("   goforge run dev      # Start development server")
	logger.Info("   goforge build        # Build for production")
	logger.Info("   goforge add <pkg>    # Add dependencies")
	logger.Info("   goforge generate     # Generate components")
	logger.Info("")
	
	logger.Info("📚 Quick Start:")
	logger.Info("   cd %s                # Navigate to project", projectName)
	logger.Info("   goforge run dev      # Start coding!")
}

func init() {
	// Enhanced flags with better descriptions
	newCmd.Flags().StringP("module-path", "m", "", 
		"Go module path (e.g., github.com/user/repo)")
	
	newCmd.Flags().StringP("template", "t", "default", 
		"Project template to use (default, minimal, microservice)")
	
	newCmd.Flags().BoolP("skip-git", "", false, 
		"Skip Git repository initialization")
	
	newCmd.Flags().BoolP("verbose", "v", false, 
		"Enable verbose logging")
	
	// Add examples
	newCmd.Example = `  # Create a simple project
  goforge new my-api

  # Create with custom module path
  goforge new user-service -m github.com/myorg/user-service

  # Create with verbose output
  goforge new blog-app --verbose

  # Create without Git initialization
  goforge new simple-app --skip-git`
}