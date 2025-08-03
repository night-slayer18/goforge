// internal/scaffold/scaffold.go - Enhanced version with better error handling and progress
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/night-slayer18/goforge/internal/validation"
)

//go:embed all:templates
var templatesFS embed.FS

// Options contains all configuration for project creation
type Options struct {
	ProjectName string
	ModulePath  string
	GoVersion   string
	DestPath    string
	Template    string
	SkipGit     bool
}

// TemplateData holds all dynamic values needed for file generation
type TemplateData struct {
	ProjectName string
	ModuleName  string
	GoVersion   string
	Name        string // For component generation
	NameTitle   string // e.g., "User"
	ModulePath  string // For component generation
}

// FileGenerationTask represents a single file to be generated
type FileGenerationTask struct {
	TemplatePath string
	TargetPath   string
	Data         TemplateData
}

// Scaffolder handles project and component generation
type Scaffolder struct {
	validator *validation.ProjectValidator
}

// NewScaffolder creates a new scaffolder instance
func NewScaffolder() *Scaffolder {
	return &Scaffolder{
		validator: validation.NewProjectValidator(),
	}
}

// CreateProject scaffolds a new project at the given destination path
// Maintains backward compatibility
func CreateProject(projectName, moduleName, goVersion, destPath string) error {
	options := Options{
		ProjectName: projectName,
		ModulePath:  moduleName,
		GoVersion:   goVersion,
		DestPath:    destPath,
		Template:    "default",
		SkipGit:     false,
	}
	return CreateProjectWithOptions(options)
}

// CreateProjectWithOptions scaffolds a new project with enhanced options
func CreateProjectWithOptions(options Options) error {
	scaffolder := NewScaffolder()
	return scaffolder.CreateProject(options)
}

// CreateProject creates a new project with the given options
func (s *Scaffolder) CreateProject(options Options) error {
	data := TemplateData{
		ProjectName: options.ProjectName,
		ModuleName:  options.ModulePath,
		GoVersion:   options.GoVersion,
	}

	// Determine template root
	templateRoot := fmt.Sprintf("templates/%s", options.Template)
	
	// Check if template exists
	if !s.templateExists(templateRoot) {
		return fmt.Errorf("template '%s' not found. Available templates: default, minimal", options.Template)
	}

	// Collect all files to generate
	tasks, err := s.collectGenerationTasks(templateRoot, options.DestPath, data)
	if err != nil {
		return fmt.Errorf("failed to collect generation tasks: %w", err)
	}

	logger.Debug("Found %d files to generate", len(tasks))

	// Generate files with progress tracking
	if err := s.generateFiles(tasks); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Initialize the project (go mod, git, etc.)
	logger.Step(2, 4, "Initializing Go module...")
	if err := s.initializeProject(options); err != nil {
		return fmt.Errorf("failed to initialize project: %w", err)
	}

	return nil
}

// collectGenerationTasks walks the template directory and collects all files to generate
func (s *Scaffolder) collectGenerationTasks(templateRoot, destPath string, data TemplateData) ([]FileGenerationTask, error) {
	var tasks []FileGenerationTask

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

		// Skip directories - they'll be created automatically
		if d.IsDir() {
			return nil
		}

		// Calculate target path
		targetPath := filepath.Join(destPath, strings.TrimSuffix(relativePath, ".tpl"))

		tasks = append(tasks, FileGenerationTask{
			TemplatePath: path,
			TargetPath:   targetPath,
			Data:         data,
		})

		return nil
	})

	return tasks, err
}

// generateFiles generates all files, potentially in parallel
func (s *Scaffolder) generateFiles(tasks []FileGenerationTask) error {
	logger.Debug("Generating %d files...", len(tasks))
	
	// For small numbers of files, generate sequentially for better error reporting
	if len(tasks) <= 10 {
		return s.generateFilesSequential(tasks)
	}
	
	// For larger projects, use parallel generation
	return s.generateFilesParallel(tasks)
}

// generateFilesSequential generates files one by one
func (s *Scaffolder) generateFilesSequential(tasks []FileGenerationTask) error {
	for i, task := range tasks {
		logger.Debug("Generating file %d/%d: %s", i+1, len(tasks), task.TargetPath)
		
		if err := s.generateFile(task); err != nil {
			return fmt.Errorf("failed to generate %s: %w", task.TargetPath, err)
		}
	}
	return nil
}

// generateFilesParallel generates files concurrently for better performance
func (s *Scaffolder) generateFilesParallel(tasks []FileGenerationTask) error {
	const maxWorkers = 5
	workers := len(tasks)
	if workers > maxWorkers {
		workers = maxWorkers
	}

	taskChan := make(chan FileGenerationTask, len(tasks))
	errChan := make(chan error, len(tasks))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if err := s.generateFile(task); err != nil {
					errChan <- fmt.Errorf("failed to generate %s: %w", task.TargetPath, err)
					return
				}
			}
		}()
	}

	// Send tasks to workers
	go func() {
		defer close(taskChan)
		for _, task := range tasks {
			taskChan <- task
		}
	}()

	// Wait for completion
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		return err
	}

	return nil
}

// generateFile generates a single file from a template
func (s *Scaffolder) generateFile(task FileGenerationTask) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(task.TargetPath), os.ModePerm); err != nil {
		return fmt.Errorf("could not create parent directory for %s: %w", task.TargetPath, err)
	}

	// Read template content
	tplContent, err := templatesFS.ReadFile(task.TemplatePath)
	if err != nil {
		return fmt.Errorf("could not read template file %s: %w", task.TemplatePath, err)
	}

	// Create template with custom functions
	tmpl, err := template.New(filepath.Base(task.TemplatePath)).
		Funcs(s.getTemplateFunctions()).
		Parse(string(tplContent))
	if err != nil {
		return fmt.Errorf("could not parse template %s: %w", task.TemplatePath, err)
	}

	// Create target file
	file, err := os.Create(task.TargetPath)
	if err != nil {
		return fmt.Errorf("could not create target file %s: %w", task.TargetPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, task.Data); err != nil {
		return fmt.Errorf("could not execute template %s: %w", task.TemplatePath, err)
	}

	logger.FileCreated(task.TargetPath)
	return nil
}

// getTemplateFunctions returns custom template functions
func (s *Scaffolder) getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"toLower":    strings.ToLower,
		"toUpper":    strings.ToUpper,
		"toCamel":    strcase.ToCamel,
		"toSnake":    strcase.ToSnake,
		"toKebab":    strcase.ToKebab,
		"pluralize":  s.pluralize,
		"timestamp":  func() string { return time.Now().Format(time.RFC3339) },
	}
}

// pluralize is a simple pluralization function
func (s *Scaffolder) pluralize(word string) string {
	if strings.HasSuffix(word, "y") {
		return strings.TrimSuffix(word, "y") + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || 
	   strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") || 
	   strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	return word + "s"
}

// templateExists checks if a template directory exists
func (s *Scaffolder) templateExists(templateRoot string) bool {
	_, err := fs.Stat(templatesFS, templateRoot)
	return err == nil
}

// initializeProject runs post-scaffolding initialization commands
func (s *Scaffolder) initializeProject(options Options) error {
	// Initialize Go module
	logger.Debug("Initializing Go module: %s", options.ModulePath)
	if err := runner.InitGoModule(options.DestPath, options.ModulePath); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	logger.Step(3, 4, "Installing dependencies...")
	if err := runner.TidyGoModule(options.DestPath); err != nil {
		return fmt.Errorf("failed to tidy go module: %w", err)
	}

	// Initialize Git repository if not skipped
	if !options.SkipGit {
		logger.Step(4, 4, "Initializing Git repository...")
		if err := runner.InitGitRepository(options.DestPath); err != nil {
			logger.Warn("Failed to initialize Git repository: %v", err)
			logger.Info("💡 You can initialize Git manually later with: git init")
		} else {
			logger.Debug("Git repository initialized successfully")
		}
	} else {
		logger.Step(4, 4, "Skipping Git initialization...")
	}

	return nil
}

// GenerateComponent scaffolds a single architectural component
func GenerateComponent(componentType, name string) error {
	scaffolder := NewScaffolder()
	return scaffolder.GenerateComponent(componentType, name)
}

// GenerateComponent generates a single component with enhanced validation
func (s *Scaffolder) GenerateComponent(componentType, name string) error {
	// Validate component name
	if err := s.validator.ValidateComponentName(componentType, name); err != nil {
		if validationErr, ok := err.(*validation.ValidationError); ok {
			logger.ValidationError(validationErr.Field, validationErr.Value, validationErr.Message, validationErr.Suggestions)
			return fmt.Errorf("invalid component name")
		}
		return err
	}

	// Load project configuration
	cfg, projectRoot, err := project.LoadConfig()
	if err != nil {
		return fmt.Errorf("command must be run from the root of a goforge project: %w", err)
	}

	logger.ComponentGenerationStart(componentType, name)

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
	case "repository":
		templateFile = "templates/components/repository.go.tpl"
		targetFile = filepath.Join(projectRoot, "internal/adapters/postgres", fmt.Sprintf("%s_repo.go", snakeName))
	case "model":
		templateFile = "templates/components/model.go.tpl"
		targetFile = filepath.Join(projectRoot, "internal/domain", fmt.Sprintf("%s.go", snakeName))
	case "middleware":
		templateFile = "templates/components/middleware.go.tpl"
		targetFile = filepath.Join(projectRoot, "internal/adapters/http/middleware", fmt.Sprintf("%s.go", snakeName))
	default:
		return fmt.Errorf("unknown component type: %s\n\nAvailable types: handler, service, repository, model, middleware", componentType)
	}

	task := FileGenerationTask{
		TemplatePath: templateFile,
		TargetPath:   targetFile,
		Data:         data,
	}

	if err := s.generateFile(task); err != nil {
		return err
	}

	logger.ComponentGenerationComplete(componentType, name, targetFile)
	s.showComponentInstructions(componentType, name)

	return nil
}

// showComponentInstructions shows helpful instructions after component generation
func (s *Scaffolder) showComponentInstructions(componentType, name string) {
	switch componentType {
	case "handler":
		logger.Info("")
		logger.Info("📋 Next steps:")
		logger.Info("   1. Register the handler in cmd/server/main.go")
		logger.Info("   2. Define your routes and HTTP methods")
		logger.Info("   3. Implement your business logic")
		
	case "service":
		logger.Info("")
		logger.Info("📋 Next steps:")
		logger.Info("   1. Wire up dependencies in cmd/server/main.go")
		logger.Info("   2. Implement your business logic methods")
		logger.Info("   3. Add any required repository interfaces")
		
	case "repository":
		logger.Info("")
		logger.Info("📋 Next steps:")
		logger.Info("   1. Define the repository interface in internal/ports")
		logger.Info("   2. Implement the database operations")
		logger.Info("   3. Wire it up in your dependency injection")
		
	case "model":
		logger.Info("")
		logger.Info("📋 Next steps:")
		logger.Info("   1. Define your domain entity fields")
		logger.Info("   2. Add validation tags if needed")
		logger.Info("   3. Consider adding business logic methods")
	}
}