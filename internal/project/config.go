package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the goforge.yml file.
type Config struct {
	ProjectName  string            `yaml:"project_name"`
	ModuleName   string            `yaml:"module_path"`
	GoVersion    string            `yaml:"go_version"`
	Dependencies map[string]string `yaml:"dependencies"`
	Scripts      map[string]string `yaml:"scripts"`
	Build        *BuildConfig      `yaml:"build"`
}

// BuildConfig defines the build-specific configuration.
type BuildConfig struct {
	Assets []string `yaml:"assets"`
}

// LoadConfig finds and parses the goforge.yml file from the current directory
// or any parent directory. It returns the parsed config, the project root
// directory (where the config was found), and any error that occurred.
func LoadConfig() (*Config, string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get current directory: %w", err)
	}

	var configPath string
	dir := currentDir
	for {
		potentialPath := filepath.Join(dir, "goforge.yml")
		if _, err := os.Stat(potentialPath); err == nil {
			configPath = potentialPath
			break
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir { // Reached the root directory
			return nil, "", fmt.Errorf("goforge.yml not found in this directory or any parent")
		}
		dir = parentDir
	}

	projectRoot := filepath.Dir(configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read goforge.yml: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, "", fmt.Errorf("failed to parse goforge.yml: %w", err)
	}

	return &cfg, projectRoot, nil
}

// SaveConfig marshals the provided Config struct back to YAML and writes it
// to the goforge.yml file in the specified project root directory.
func SaveConfig(projectRoot string, cfg *Config) error {
	configPath := filepath.Join(projectRoot, "goforge.yml")

	// Marshal the struct into YAML format.
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write the new content back to the goforge.yml file.
	// 0644 is a standard file permission for text files.
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to goforge.yml: %w", err)
	}

	return nil
}
