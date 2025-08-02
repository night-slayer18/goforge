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

// This directive embeds the 'templates' directory.
// For this to work, the 'templates' directory MUST be located inside
// this 'scaffold' package directory.
//go:embed all:templates
var templatesFS embed.FS

// TemplateData holds all dynamic values needed for file generation.
type TemplateData struct {
	ProjectName string
	ModuleName  string
	GoVersion   string
	Name        string // For component generation
	NameTitle   string // e.g., "User"
	ModulePath  string // For component generation
}

// CreateProject scaffolds a new project at the given destination path.
func CreateProject(projectName, moduleName, goVersion, destPath string) error {
	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
		GoVersion:   goVersion,
	}

	templateRoot := "templates/default"

	err := fs.WalkDir(templatesFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(templateRoot, path)
		if err != nil {
			return err
		}
		if relativePath == "." {
			return nil
		}
		targetPath := filepath.Join(destPath, strings.TrimSuffix(relativePath, ".tpl"))
		if d.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}
		return generateFile(path, targetPath, data)
	})

	if err != nil {
		return fmt.Errorf("error during file generation: %w", err)
	}

	// Call the initialization function to keep logic separate and clean.
	return InitializeProject(data)
}

// GenerateComponent scaffolds a single architectural component.
func GenerateComponent(componentType, name string) error {
	cfg, projectRoot, err := project.LoadConfig()
	if err != nil {
		return fmt.Errorf("command must be run from the root of a goforge project: %w", err)
	}

	data := TemplateData{
		Name:       name,
		NameTitle:  strcase.ToCamel(name),
		ModulePath: cfg.ModuleName,
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
func generateFile(templatePath, targetPath string, data interface{}) error {
	// --- THIS IS THE FIX ---
	// Ensure the parent directory of the target file exists before creating the file.
	if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
		return fmt.Errorf("could not create parent directory for %s: %w", targetPath, err)
	}

	tplContent, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("could not read template file %s: %w", templatePath, err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tplContent))
	if err != nil {
		return fmt.Errorf("could not parse template %s: %w", templatePath, err)
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("could not create target file %s: %w", targetPath, err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// InitializeProject runs post-scaffolding commands.
func InitializeProject(data TemplateData) error {
	projectPath := filepath.Join(".", data.ProjectName)
	fmt.Println("\n🔧 Initializing Go module...")
	if err := runner.InitGoModule(projectPath, data.ModuleName); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	fmt.Println("🧹 Tidying dependencies...")
	if err := runner.TidyGoModule(projectPath); err != nil {
		return fmt.Errorf("failed to tidy go module: %w", err)
	}

	return nil
}
