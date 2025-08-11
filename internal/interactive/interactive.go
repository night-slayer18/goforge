package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/night-slayer18/goforge/internal/validation"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// InteractiveSession manages the interactive CLI session
type InteractiveSession struct {
	scanner   *bufio.Scanner
	validator *validation.ProjectValidator
}

// NewInteractiveSession creates a new interactive session
func NewInteractiveSession() *InteractiveSession {
	return &InteractiveSession{
		scanner:   bufio.NewScanner(os.Stdin),
		validator: validation.NewProjectValidator(),
	}
}

// ProjectOptions holds all the options that can be configured interactively
type ProjectOptions struct {
	ProjectName string
	ModulePath  string
	Template    string
	SkipGit     bool
	Verbose     bool
}

// RunProjectCreationWizard runs the interactive project creation wizard
func (is *InteractiveSession) RunProjectCreationWizard() (*ProjectOptions, error) {
	options := &ProjectOptions{}
	
	// Welcome message
	is.printWelcome()
	
	// Step 1: Project name
	projectName, err := is.promptProjectName()
	if err != nil {
		return nil, err
	}
	options.ProjectName = projectName
	
	// Step 2: Module path
	modulePath, err := is.promptModulePath(projectName)
	if err != nil {
		return nil, err
	}
	options.ModulePath = modulePath
	
	// Step 3: Template selection
	template, err := is.promptTemplateSelection()
	if err != nil {
		return nil, err
	}
	options.Template = template
	
	// Step 4: Git initialization
	skipGit, err := is.promptGitInit()
	if err != nil {
		return nil, err
	}
	options.SkipGit = skipGit
	
	// Step 5: Verbose output
	verbose, err := is.promptVerboseOutput()
	if err != nil {
		return nil, err
	}
	options.Verbose = verbose
	
	// Summary
	is.showSummary(options)
	
	// Confirmation
	confirmed, err := is.promptConfirmation()
	if err != nil {
		return nil, err
	}
	
	if !confirmed {
		return nil, fmt.Errorf("project creation cancelled by user")
	}
	
	return options, nil
}

func (is *InteractiveSession) printWelcome() {
	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println("🚀 Welcome to GoForge Project Creator!")
	fmt.Println("Let's create your new Go project step by step.")
	fmt.Println("(Press Ctrl+C anytime to cancel)")
	fmt.Println()
}

func (is *InteractiveSession) promptProjectName() (string, error) {
	for {
		fmt.Print(color.New(color.FgYellow).Sprint("📝 Project name: "))
		
		if !is.scanner.Scan() {
			return "", fmt.Errorf("failed to read input")
		}
		
		name := strings.TrimSpace(is.scanner.Text())
		if name == "" {
			color.New(color.FgRed).Println("   ❌ Project name cannot be empty")
			continue
		}
		
		// Validate project name
		if err := is.validator.ValidateProjectName(name); err != nil {
			if validationErr, ok := err.(*validation.ValidationError); ok {
				color.New(color.FgRed).Printf("   ❌ %s\n", validationErr.Message)
				if len(validationErr.Suggestions) > 0 {
					color.New(color.FgBlue).Println("   💡 Suggestions:")
					for _, suggestion := range validationErr.Suggestions {
						color.New(color.FgBlue).Printf("      • %s\n", suggestion)
					}
				}
				continue
			}
			return "", err
		}
		
		// Check if directory exists
		if _, err := os.Stat(name); err == nil {
			color.New(color.FgRed).Printf("   ❌ Directory '%s' already exists\n", name)
			continue
		}
		
		color.New(color.FgGreen).Printf("   ✅ Project name: %s\n", name)
		return name, nil
	}
}

func (is *InteractiveSession) promptModulePath(defaultPath string) (string, error) {
	fmt.Printf("📦 Module path (press Enter for '%s'): ", defaultPath)
	
	if !is.scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}
	
	modulePath := strings.TrimSpace(is.scanner.Text())
	if modulePath == "" {
		modulePath = defaultPath
	}
	
	// Validate module path
	if err := is.validator.ValidateModulePath(modulePath); err != nil {
		if validationErr, ok := err.(*validation.ValidationError); ok {
			color.New(color.FgRed).Printf("   ❌ %s\n", validationErr.Message)
			// Ask again
			return is.promptModulePath(defaultPath)
		}
		return "", err
	}
	
	color.New(color.FgGreen).Printf("   ✅ Module path: %s\n", modulePath)
	return modulePath, nil
}

func (is *InteractiveSession) promptTemplateSelection() (string, error) {
	templates := []Template{
		{Name: "default", Description: "Full-featured web API with clean architecture"},
		{Name: "minimal", Description: "Lightweight template with basic structure"},
		{Name: "microservice", Description: "Microservice template with Docker and health checks"},
	}
	
	fmt.Println("📋 Available templates:")
	for i, template := range templates {
		fmt.Printf("   %d. %s - %s\n", i+1, 
			color.New(color.FgCyan).Sprint(template.Name), 
			template.Description)
	}
	
	for {
		fmt.Print("Select template (1-3, or press Enter for default): ")
		
		if !is.scanner.Scan() {
			return "", fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(is.scanner.Text())
		if input == "" {
			color.New(color.FgGreen).Println("   ✅ Template: default")
			return "default", nil
		}
		
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(templates) {
			color.New(color.FgRed).Println("   ❌ Invalid selection. Please choose 1-3.")
			continue
		}
		
		selected := templates[choice-1]
		color.New(color.FgGreen).Printf("   ✅ Template: %s\n", selected.Name)
		return selected.Name, nil
	}
}

func (is *InteractiveSession) promptGitInit() (bool, error) {
	for {
		fmt.Print("🔧 Initialize Git repository? (Y/n): ")
		
		if !is.scanner.Scan() {
			return false, fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(strings.ToLower(is.scanner.Text()))
		
		switch input {
		case "", "y", "yes":
			color.New(color.FgGreen).Println("   ✅ Git repository will be initialized")
			return false, nil // false means don't skip
		case "n", "no":
			color.New(color.FgYellow).Println("   ⚠️  Git repository will be skipped")
			return true, nil // true means skip
		default:
			color.New(color.FgRed).Println("   ❌ Please answer 'y' (yes) or 'n' (no)")
		}
	}
}

func (is *InteractiveSession) promptVerboseOutput() (bool, error) {
	for {
		fmt.Print("🔍 Enable verbose output? (y/N): ")
		
		if !is.scanner.Scan() {
			return false, fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(strings.ToLower(is.scanner.Text()))
		
		switch input {
		case "", "n", "no":
			color.New(color.FgGreen).Println("   ✅ Standard output mode")
			return false, nil
		case "y", "yes":
			color.New(color.FgGreen).Println("   ✅ Verbose output enabled")
			return true, nil
		default:
			color.New(color.FgRed).Println("   ❌ Please answer 'y' (yes) or 'n' (no)")
		}
	}
}

func (is *InteractiveSession) showSummary(options *ProjectOptions) {
	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println("📋 Project Summary:")
	fmt.Printf("   Project Name: %s\n", color.New(color.FgGreen).Sprint(options.ProjectName))
	fmt.Printf("   Module Path:  %s\n", color.New(color.FgGreen).Sprint(options.ModulePath))
	fmt.Printf("   Template:     %s\n", color.New(color.FgGreen).Sprint(options.Template))
	
	gitStatus := "Yes"
	if options.SkipGit {
		gitStatus = "No"
	}
	fmt.Printf("   Git Init:     %s\n", color.New(color.FgGreen).Sprint(gitStatus))
	
	verboseStatus := "No"
	if options.Verbose {
		verboseStatus = "Yes"
	}
	fmt.Printf("   Verbose:      %s\n", color.New(color.FgGreen).Sprint(verboseStatus))
	fmt.Println()
}

func (is *InteractiveSession) promptConfirmation() (bool, error) {
	for {
		fmt.Print("✨ Create this project? (Y/n): ")
		
		if !is.scanner.Scan() {
			return false, fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(strings.ToLower(is.scanner.Text()))
		
		switch input {
		case "", "y", "yes":
			return true, nil
		case "n", "no":
			color.New(color.FgYellow).Println("Project creation cancelled.")
			return false, nil
		default:
			color.New(color.FgRed).Println("Please answer 'y' (yes) or 'n' (no)")
		}
	}
}

// Template represents a project template
type Template struct {
	Name        string
	Description string
}

// ComponentWizard handles interactive component generation
type ComponentWizard struct {
	scanner   *bufio.Scanner
	validator *validation.ProjectValidator
}

// NewComponentWizard creates a new component wizard
func NewComponentWizard() *ComponentWizard {
	return &ComponentWizard{
		scanner:   bufio.NewScanner(os.Stdin),
		validator: validation.NewProjectValidator(),
	}
}

// ComponentOptions holds component generation options
type ComponentOptions struct {
	Type string
	Name string
}

// RunComponentCreationWizard runs the interactive component creation wizard
func (cw *ComponentWizard) RunComponentCreationWizard() (*ComponentOptions, error) {
	options := &ComponentOptions{}
	
	// Welcome
	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println("🔧 Component Generator")
	fmt.Println()
	
	// Select component type
	componentType, err := cw.promptComponentType()
	if err != nil {
		return nil, err
	}
	options.Type = componentType
	
	// Get component name
	componentName, err := cw.promptComponentName(componentType)
	if err != nil {
		return nil, err
	}
	options.Name = componentName
	
	return options, nil
}

func (cw *ComponentWizard) promptComponentType() (string, error) {
	components := []struct {
		Name        string
		Description string
	}{
		{"handler", "HTTP request handlers for API endpoints"},
		{"service", "Business logic services"},
		{"repository", "Data access layer implementations"},
		{"model", "Domain models and entities"},
		{"middleware", "HTTP middleware components"},
		{"port", "Interface definitions for clean architecture"},
	}
	
	fmt.Println("Available components:")
	for i, comp := range components {
		fmt.Printf("   %d. %s - %s\n", i+1, 
			color.New(color.FgCyan).Sprint(comp.Name), 
			comp.Description)
	}
	
	for {
		fmt.Print("Select component type (1-6): ")
		
		if !cw.scanner.Scan() {
			return "", fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(cw.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(components) {
			color.New(color.FgRed).Println("   ❌ Invalid selection. Please choose 1-6.")
			continue
		}
		
		selected := components[choice-1]
		color.New(color.FgGreen).Printf("   ✅ Component type: %s\n", selected.Name)
		return selected.Name, nil
	}
}

func (cw *ComponentWizard) promptComponentName(componentType string) (string, error) {
    titleCaser := cases.Title(language.English) // added
    for {
        // replaced strings.Title with titleCaser.String(componentType)
        fmt.Printf("📝 %s name: ", titleCaser.String(componentType))
        
        if !cw.scanner.Scan() {
            return "", fmt.Errorf("failed to read input")
        }
        
        name := strings.TrimSpace(cw.scanner.Text())
        if name == "" {
            color.New(color.FgRed).Printf("   ❌ %s name cannot be empty\n", componentType)
            continue
        }
        
        // Validate component name
        if err := cw.validator.ValidateComponentName(componentType, name); err != nil {
            if validationErr, ok := err.(*validation.ValidationError); ok {
                color.New(color.FgRed).Printf("   ❌ %s\n", validationErr.Message)
                if len(validationErr.Suggestions) > 0 {
                    color.New(color.FgBlue).Println("   💡 Suggestions:")
                    for _, suggestion := range validationErr.Suggestions {
                        color.New(color.FgBlue).Printf("      • %s\n", suggestion)
                    }
                }
                continue
            }
            return "", err
        }
        
        color.New(color.FgGreen).Printf("   ✅ %s name: %s\n", titleCaser.String(componentType), name)
        return name, nil
    }
}

// IsInteractiveTerminal checks if we're running in an interactive terminal
func IsInteractiveTerminal() bool {
	// Check if stdin is a terminal
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	
	// If it's a character device (terminal), it's interactive
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// ShouldUseInteractiveMode determines if interactive mode should be used
func ShouldUseInteractiveMode(args []string, interactiveFlag bool) bool {
	// If explicitly requested via flag
	if interactiveFlag {
		return IsInteractiveTerminal()
	}
	
	// If no arguments provided and we're in a terminal, suggest interactive mode
	if len(args) == 0 && IsInteractiveTerminal() {
		return true
	}
	
	return false
}

// PromptForInteractiveMode asks user if they want to use interactive mode
func PromptForInteractiveMode() bool {
	if !IsInteractiveTerminal() {
		return false
	}
	
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Print("No project name provided. Use interactive mode? (Y/n): ")
	
	if !scanner.Scan() {
		return false
	}
	
	input := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return input == "" || input == "y" || input == "yes"
}