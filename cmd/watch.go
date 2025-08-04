// cmd/watch.go - Enhanced version with actual file watching
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/night-slayer18/goforge/internal/runner"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch [script-name]",
	Short: "Watch for file changes and restart the application",
	Long: `The watch command monitors your project files for changes and automatically
restarts the specified script when changes are detected. If no script is specified,
it defaults to the 'dev' script.

This is useful for development workflows where you want automatic reloading.

Examples:
  goforge watch           # Watch and run 'dev' script
  goforge watch dev       # Same as above
  goforge watch test      # Watch and run 'test' script`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up logging
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger.SetVerbose(verbose)

		cfg, projectRoot, err := project.LoadConfig()
		if err != nil {
			return fmt.Errorf("command must be run from a goforge project: %w", err)
		}

		// Determine script to run
		scriptName := "dev"
		if len(args) > 0 {
			scriptName = args[0]
		}

		script, exists := cfg.Scripts[scriptName]
		if !exists {
			return fmt.Errorf("script '%s' not found in goforge.yml\n\nAvailable scripts:\n%s", 
				scriptName, formatAvailableScripts(cfg.Scripts))
		}

		logger.Info("👀 Starting watch mode for script: %s", scriptName)
		logger.Info("📝 Command: %s", script)
		logger.Info("📁 Watching: %s", projectRoot)
		logger.Info("🔄 Press Ctrl+C to stop")
		logger.Info("")

		// Create file watcher
		watcher, err := NewFileWatcher(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to create file watcher: %w", err)
		}
		defer watcher.Close()

		// Set up graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start watch mode
		watchMode := runner.NewWatchMode(projectRoot, "sh", "-c", script)
		if err := watchMode.Start(); err != nil {
			return fmt.Errorf("failed to start watch mode: %w", err)
		}

		// File change handling
		debouncer := NewDebouncer(1 * time.Second) // Increase debounce time
		lastRestart := time.Now().Add(-10 * time.Second) // Allow immediate first restart
		
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				
				if watcher.ShouldIgnore(event.Name) {
					logger.Debug("Ignoring file change: %s", event.Name)
					continue
				}

				// Only care about write and create events for Go files
				if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
					continue
				}

				logger.Debug("File changed: %s (%s)", event.Name, event.Op)
				
				// Prevent too frequent restarts
				if time.Since(lastRestart) < 2*time.Second {
					logger.Debug("Restart too recent, debouncing...")
					continue
				}
				
				// Debounce rapid file changes
				debouncer.Debounce(func() {
					lastRestart = time.Now()
					logger.Info("🔄 Files changed, restarting...")
					if err := watchMode.Restart(); err != nil {
						logger.Error("Failed to restart: %v", err)
					}
				})

			case err, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
				logger.Warn("Watch error: %v", err)

			case <-sigChan:
				logger.Info("\n🛑 Shutting down...")
				if err := watchMode.Stop(); err != nil {
					logger.Error("Error stopping watch mode: %v", err)
					return err
				}
				logger.Info("✅ Watch mode stopped")
				return nil
			}
		}
	},
}

// FileWatcher wraps fsnotify.Watcher with project-specific logic
type FileWatcher struct {
	*fsnotify.Watcher
	projectRoot string
	ignorePaths []string
}

// NewFileWatcher creates a new file watcher for the project
func NewFileWatcher(projectRoot string) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		Watcher:     watcher,
		projectRoot: projectRoot,
		ignorePaths: []string{
			"dist",
			"vendor",
			"node_modules",
			".git",
			".idea",
			".vscode",
			"*.test",
			"*.out",
			"coverage.html",
		},
	}

	// Add project root and subdirectories to watch
	if err := fw.addWatchPaths(); err != nil {
		watcher.Close()
		return nil, err
	}

	return fw, nil
}

// addWatchPaths recursively adds directories to watch
func (fw *FileWatcher) addWatchPaths() error {
	return filepath.Walk(fw.projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-directories
		if !info.IsDir() {
			return nil
		}

		// Check if this path should be ignored
		relPath, err := filepath.Rel(fw.projectRoot, path)
		if err != nil {
			return err
		}

		if fw.shouldIgnorePath(relPath) {
			logger.Debug("Ignoring directory: %s", relPath)
			return filepath.SkipDir
		}

		logger.Debug("Watching directory: %s", relPath)
		return fw.Add(path)
	})
}

// shouldIgnorePath checks if a path should be ignored based on patterns
func (fw *FileWatcher) shouldIgnorePath(path string) bool {
	for _, pattern := range fw.ignorePaths {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
		// Also check if any part of the path matches
		parts := strings.Split(path, string(filepath.Separator))
		for _, part := range parts {
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
		}
	}
	return false
}

// ShouldIgnore checks if a file change event should be ignored
func (fw *FileWatcher) ShouldIgnore(filename string) bool {
	// Get relative path
	relPath, err := filepath.Rel(fw.projectRoot, filename)
	if err != nil {
		return true
	}

	// Ignore based on directory patterns
	if fw.shouldIgnorePath(relPath) {
		return true
	}

	// Only watch Go files and config files
	ext := filepath.Ext(filename)
	if ext != ".go" && ext != ".yml" && ext != ".yaml" && ext != ".json" && ext != ".mod" && ext != ".sum" {
		return true
	}

	// Ignore test files during development (they don't affect the running server)
	if strings.HasSuffix(filename, "_test.go") {
		return true
	}

	// Ignore temporary files and editor artifacts
	base := filepath.Base(filename)
	if strings.HasPrefix(base, ".") || 
	   strings.HasPrefix(base, "~") || 
	   strings.HasSuffix(base, ".tmp") ||
	   strings.HasSuffix(base, ".swp") ||
	   strings.HasSuffix(base, ".swo") ||
	   strings.Contains(base, "#") {
		return true
	}

	// Ignore generated files
	if strings.Contains(relPath, "/.git/") ||
	   strings.Contains(relPath, "/vendor/") ||
	   strings.Contains(relPath, "/node_modules/") {
		return true
	}

	logger.Debug("Will watch file: %s", relPath)
	return false
}

// Debouncer helps prevent rapid successive restarts
type Debouncer struct {
	duration time.Duration
	timer    *time.Timer
}

// NewDebouncer creates a new debouncer
func NewDebouncer(duration time.Duration) *Debouncer {
	return &Debouncer{duration: duration}
}

// Debounce executes the function after the specified delay, 
// canceling any previous pending execution
func (d *Debouncer) Debounce(fn func()) {
	if d.timer != nil {
		d.timer.Stop()
	}
	
	d.timer = time.AfterFunc(d.duration, fn)
}

func formatAvailableScripts(scripts map[string]string) string {
	if len(scripts) == 0 {
		return "  No scripts defined"
	}

	result := ""
	for name, command := range scripts {
		result += fmt.Sprintf("  %s: %s\n", name, command)
	}
	return result
}

func init() {
	watchCmd.Flags().BoolP("verbose", "v", false, "Enable verbose logging")
}