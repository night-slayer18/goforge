package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
)

// The crucial //go:embed directive. This tells the Go compiler to embed the
// entire 'templates' directory (located two levels up from this file)
// into the templatesFS variable.
//go:embed all:templates
var templatesFS embed.FS

// TemplateData holds the dynamic data passed to templates for rendering.
type TemplateData struct {
	ProjectName string
	ModuleName  string
	GoVersion   string
	Name        string // For component generation
	NameTitle   string // e.g., "User"
}

// CreateProject scaffolds a new project at the given destination path.
func CreateProject(projectName, moduleName, goVersion, destPath string) error {
	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
		GoVersion:   goVersion,
	}

	// Create the root project directory.
	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create project directory: %w", err)
	}

	templateRoot := "templates/default"

	// Walk through the embedded templates and generate files.
	err := fs.WalkDir(templatesFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path from the template root (e.g., "cmd/server/main.go.tpl").
		relativePath := strings.TrimPrefix(path, templateRoot+"/")
		if relativePath == "" || relativePath == "." {
			return nil // Skip the root directory itself.
		}

		targetPath := filepath.Join(destPath, relativePath)
		// Remove the .tpl extension for the final file name.
		targetPath = strings.TrimSuffix(targetPath, ".tpl")

		if d.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}

		return generateFile(path, targetPath, data)
	})

	if err != nil {
		return fmt.Errorf("error during file generation: %w", err)
	}

	// Run post-generation commands to finalize the project setup.
	fmt.Println("🔧 Initializing Go module...")
	if err := runner.InitGoModule(destPath, moduleName); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	fmt.Println("🧹 Tidying dependencies...")
	if err := runner.TidyGoModule(destPath); err != nil {
		return fmt.Errorf("failed to tidy go module: %w", err)
	}

	return nil
}

// GenerateComponent scaffolds a single architectural component.
func GenerateComponent(componentType, name string) error {
	// Load the project config to get the module path
	cfg, projectRoot, err := project.LoadConfig()
	if err != nil {
		return fmt.Errorf("command must be run from the root of a goforge project: %w", err)
	}

	data := TemplateData{
		Name:       name,
		NameTitle:  strcase.ToCamel(name), // Use a library for robust case conversion
		ModuleName: cfg.ModuleName,
	}

	var templateFile, targetFile string
	snakeName := strcase.ToSnake(name)

	switch componentType {
	case "handler":
		templateFile = "templates/components/handler.go.tpl"
		targetFile = filepath.Join(projectRoot, "internal/adapters/http/handler", fmt.Sprintf("%s_handler.go", snakeName))
	case "service":
		templateFile = "templates/components/service.go.tpl"
		targetFile = filepath.Join(projectRoot, "internal/app/service", fmt.Sprintf("%s_service.go", snakeName))
	default:
		return fmt.Errorf("unknown component type: %s", componentType)
	}

	return generateFile(templateFile, targetFile, data)
}

// generateFile is a helper to parse and execute a single template file.
func generateFile(templateFile, targetFile string, data TemplateData) error {
	// Read the template content from the embedded filesystem.
	tplContent, err := templatesFS.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("could not read template file %s: %w", templateFile, err)
	}

	// Parse the template.
	tmpl, err := template.New(filepath.Base(templateFile)).Parse(string(tplContent))
	if err != nil {
		return fmt.Errorf("could not parse template %s: %w", templateFile, err)
	}

	// Create the destination file.
	file, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("could not create target file %s: %w", targetFile, err)
	}
	defer file.Close()

	// Execute the template and write the output to the file.
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("could not execute template %s: %w", templateFile, err)
	}

	return nil
}
