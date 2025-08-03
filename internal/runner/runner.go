// internal/runner/runner.go - Enhanced version with better error handling and logging
package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/night-slayer18/goforge/internal/logger"
)

// CommandOptions provides configuration for command execution
type CommandOptions struct {
	Dir         string
	Env         []string
	Timeout     time.Duration
	ShowOutput  bool
	ShowCommand bool
}

// DefaultOptions returns sensible default options
func DefaultOptions() *CommandOptions {
	return &CommandOptions{
		Dir:         "",
		Env:         os.Environ(),
		Timeout:     5 * time.Minute,
		ShowOutput:  true,
		ShowCommand: true,
	}
}

// ExecuteCommand runs an external command with enhanced error handling and logging
func ExecuteCommand(dir, name string, args ...string) error {
	opts := DefaultOptions()
	opts.Dir = dir
	return ExecuteCommandWithOptions(name, args, opts)
}

// ExecuteCommandWithOptions runs a command with custom options
func ExecuteCommandWithOptions(name string, args []string, opts *CommandOptions) error {
	start := time.Now()
	
	if opts.ShowCommand {
		logger.CommandStart(name, args...)
	}
	
	// Create command with context for timeout support
	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}
	
	cmd := exec.CommandContext(ctx, name, args...)
	
	// Set working directory
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}
	
	// Set environment
	if len(opts.Env) > 0 {
		cmd.Env = opts.Env
	}
	
	// Configure output handling
	if opts.ShowOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin
	
	// Execute command
	err := cmd.Run()
	duration := time.Since(start)
	
	// Log result
	if err != nil {
		// Check for timeout
		if ctx.Err() == context.DeadlineExceeded {
			logger.CommandError(name, fmt.Errorf("command timed out after %v", opts.Timeout), duration)
			return fmt.Errorf("command '%s' timed out after %v", name, opts.Timeout)
		}
		
		// Check for exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			logger.CommandError(name, fmt.Errorf("exit code %d", exitCode), duration)
			return fmt.Errorf("command '%s' failed with exit code %d", name, exitCode)
		}
		
		logger.CommandError(name, err, duration)
		return fmt.Errorf("command '%s' failed: %w", name, err)
	}
	
	logger.CommandSuccess(name, duration)
	return nil
}

// ExecuteScript runs a shell script with enhanced error handling
func ExecuteScript(dir, script string) error {
	return ExecuteScriptWithOptions(dir, script, DefaultOptions())
}

// ExecuteScriptWithOptions runs a shell script with custom options
func ExecuteScriptWithOptions(dir, script string, opts *CommandOptions) error {
	opts.Dir = dir
	return ExecuteCommandWithOptions("sh", []string{"-c", script}, opts)
}

// ExecuteCommandWithOutput runs a command and captures its output
func ExecuteCommandWithOutput(dir, name string, args ...string) (string, error) {
	start := time.Now()
	logger.CommandStart(name, args...)
	
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	
	output, err := cmd.Output()
	duration := time.Since(start)
	
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			stderr := string(exitError.Stderr)
			logger.CommandError(name, fmt.Errorf("exit code %d: %s", exitError.ExitCode(), stderr), duration)
			return "", fmt.Errorf("command '%s' failed: %s", name, stderr)
		}
		logger.CommandError(name, err, duration)
		return "", err
	}
	
	logger.CommandSuccess(name, duration)
	return strings.TrimSpace(string(output)), nil
}

// StreamingExecutor provides real-time output streaming for long-running commands
type StreamingExecutor struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

// NewStreamingExecutor creates a new streaming executor
func NewStreamingExecutor(dir, name string, args ...string) (*StreamingExecutor, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdout.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	return &StreamingExecutor{
		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

// Start begins command execution with streaming output
func (se *StreamingExecutor) Start() error {
	logger.CommandStart(se.cmd.Path, se.cmd.Args[1:]...)
	
	// Start the command
	if err := se.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}
	
	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(se.stdout)
		for scanner.Scan() {
			logger.Info("📤 %s", scanner.Text())
		}
	}()
	
	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(se.stderr)
		for scanner.Scan() {
			logger.Error("📥 %s", scanner.Text())
		}
	}()
	
	return nil
}

// Wait waits for the command to complete
func (se *StreamingExecutor) Wait() error {
	defer se.stdout.Close()
	defer se.stderr.Close()
	
	err := se.cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("command failed with exit code %d", exitError.ExitCode())
		}
		return fmt.Errorf("command failed: %w", err)
	}
	
	return nil
}

// Kill terminates the running command
func (se *StreamingExecutor) Kill() error {
	if se.cmd.Process != nil {
		return se.cmd.Process.Kill()
	}
	return nil
}

// InitGoModule runs 'go mod init' with enhanced error handling
func InitGoModule(dir, modulePath string) error {
	logger.Debug("Initializing Go module: %s", modulePath)
	
	opts := DefaultOptions()
	opts.Dir = dir
	opts.ShowOutput = false // We'll handle output ourselves
	
	err := ExecuteCommandWithOptions("go", []string{"mod", "init", modulePath}, opts)
	if err != nil {
		return fmt.Errorf("failed to initialize Go module '%s': %w\n\nTroubleshooting:\n  • Ensure Go is installed and in PATH\n  • Check that the module path is valid\n  • Verify you have write permissions in the directory", modulePath, err)
	}
	
	logger.Debug("Go module initialized successfully")
	return nil
}

// TidyGoModule runs 'go mod tidy' with enhanced error handling
func TidyGoModule(dir string) error {
	logger.Debug("Tidying Go module dependencies...")
	
	opts := DefaultOptions()
	opts.Dir = dir
	opts.Timeout = 2 * time.Minute // Increase timeout for network operations
	
	err := ExecuteCommandWithOptions("go", []string{"mod", "tidy"}, opts)
	if err != nil {
		return fmt.Errorf("failed to tidy Go module: %w\n\nTroubleshooting:\n  • Check your internet connection\n  • Verify go.mod file is valid\n  • Ensure dependencies are accessible\n  • Try running 'go mod tidy' manually for more details", err)
	}
	
	logger.Debug("Go module dependencies tidied successfully")
	return nil
}

// InitGitRepository runs 'git init' with enhanced error handling
func InitGitRepository(dir string) error {
	logger.Debug("Initializing Git repository...")
	
	// Check if Git is available
	if !isCommandAvailable("git") {
		return fmt.Errorf("git is not installed or not available in PATH")
	}
	
	// Check if already a Git repository
	if isGitRepository(dir) {
		logger.Debug("Directory is already a Git repository")
		return nil
	}
	
	opts := DefaultOptions()
	opts.Dir = dir
	opts.ShowOutput = false
	
	err := ExecuteCommandWithOptions("git", []string{"init", "-b", "main"}, opts)
	if err != nil {
		// Try fallback for older Git versions
		logger.Debug("Trying fallback Git init for older versions...")
		err = ExecuteCommandWithOptions("git", []string{"init"}, opts)
		if err != nil {
			return fmt.Errorf("failed to initialize Git repository: %w", err)
		}
		
		// Set default branch to main for older Git versions
		_ = ExecuteCommandWithOptions("git", []string{"checkout", "-b", "main"}, opts)
	}
	
	// Create initial commit
	if err := createInitialCommit(dir); err != nil {
		logger.Warn("Failed to create initial commit: %v", err)
		// Don't fail the entire process for this
	}
	
	logger.Debug("Git repository initialized successfully")
	return nil
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// isGitRepository checks if a directory is already a Git repository
func isGitRepository(dir string) bool {
	gitDir := fmt.Sprintf("%s/.git", dir)
	if _, err := os.Stat(gitDir); err == nil {
		return true
	}
	return false
}

// createInitialCommit creates an initial commit in the Git repository
func createInitialCommit(dir string) error {
	opts := DefaultOptions()
	opts.Dir = dir
	opts.ShowOutput = false
	
	// Add all files
	if err := ExecuteCommandWithOptions("git", []string{"add", "."}, opts); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}
	
	// Set user config if not set (for CI environments)
	_ = ExecuteCommandWithOptions("git", []string{"config", "user.name", "GoForge"}, opts)
	_ = ExecuteCommandWithOptions("git", []string{"config", "user.email", "goforge@localhost"}, opts)
	
	// Create initial commit
	if err := ExecuteCommandWithOptions("git", []string{"commit", "-m", "Initial commit from GoForge"}, opts); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}
	
	return nil
}

// InstallDependency adds a Go dependency with enhanced error handling
func InstallDependency(dir, module string) error {
	logger.DependencyAdding(module)
	
	opts := DefaultOptions()
	opts.Dir = dir
	opts.Timeout = 3 * time.Minute // Longer timeout for downloads
	
	err := ExecuteCommandWithOptions("go", []string{"get", module}, opts)
	if err != nil {
		return fmt.Errorf("failed to install dependency '%s': %w\n\nTroubleshooting:\n  • Check your internet connection\n  • Verify the module path is correct\n  • Ensure the module version exists\n  • Check if the module requires authentication", module, err)
	}
	
	logger.DependencyAdded(module)
	return nil
}

// BuildBinary builds a Go binary with enhanced error handling
func BuildBinary(dir, outputPath, entrypoint string) error {
	logger.BuildStart(fmt.Sprintf("binary at %s", outputPath))
	start := time.Now()
	
	// Ensure output directory exists
	if err := os.MkdirAll(fmt.Sprintf("%s/..", outputPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	opts := DefaultOptions()
	opts.Dir = dir
	opts.Timeout = 5 * time.Minute
	
	args := []string{"build", "-o", outputPath}
	if entrypoint != "" {
		args = append(args, entrypoint)
	}
	
	err := ExecuteCommandWithOptions("go", args, opts)
	if err != nil {
		return fmt.Errorf("failed to build binary: %w\n\nTroubleshooting:\n  • Check for compilation errors above\n  • Ensure all dependencies are available\n  • Verify the entry point path is correct", err)
	}
	
	duration := time.Since(start)
	logger.BuildComplete(outputPath, duration)
	return nil
}

// RunTests executes Go tests with enhanced output
func RunTests(dir string, packages ...string) error {
	logger.Info("🧪 Running tests...")
	start := time.Now()
	
	args := []string{"test", "-v"}
	if len(packages) > 0 {
		args = append(args, packages...)
	} else {
		args = append(args, "./...")
	}
	
	opts := DefaultOptions()
	opts.Dir = dir
	opts.Timeout = 10 * time.Minute
	
	err := ExecuteCommandWithOptions("go", args, opts)
	duration := time.Since(start)
	
	if err != nil {
		logger.Error("❌ Tests failed after %v", duration.Round(time.Millisecond))
		return fmt.Errorf("tests failed: %w", err)
	}
	
	logger.Success("✅ All tests passed in %v", duration.Round(time.Millisecond))
	return nil
}

// WatchMode provides file watching capabilities for development
type WatchMode struct {
	dir     string
	command string
	args    []string
	process *os.Process
}

// NewWatchMode creates a new watch mode instance
func NewWatchMode(dir, command string, args ...string) *WatchMode {
	return &WatchMode{
		dir:     dir,
		command: command,
		args:    args,
	}
}

// Start begins watching for file changes and restarting the command
func (w *WatchMode) Start() error {
	logger.Info("👀 Starting watch mode...")
	logger.Info("🔄 Watching directory: %s", w.dir)
	logger.Info("📝 Press Ctrl+C to stop")
	
	// Initial run
	if err := w.restart(); err != nil {
		return err
	}
	
	// TODO: Implement file watching logic
	// This would require a file watching library like fsnotify
	// For now, we'll just run once
	
	return nil
}

// restart stops the current process and starts a new one
func (w *WatchMode) restart() error {
	// Stop existing process
	if w.process != nil {
		logger.Debug("Stopping existing process...")
		if err := w.process.Signal(syscall.SIGTERM); err != nil {
			w.process.Kill()
		}
		w.process.Wait()
	}
	
	// Start new process
	logger.Info("🔄 Restarting: %s %s", w.command, strings.Join(w.args, " "))
	
	cmd := exec.Command(w.command, w.args...)
	cmd.Dir = w.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}
	
	w.process = cmd.Process
	logger.Info("✅ Process started (PID: %d)", w.process.Pid)
	
	return nil
}

// Stop terminates the watch mode
func (w *WatchMode) Stop() error {
	if w.process != nil {
		logger.Info("🛑 Stopping watch mode...")
		err := w.process.Signal(syscall.SIGTERM)
		if err != nil {
			return w.process.Kill()
		}
		_, err = w.process.Wait()
		return err
	}
	return nil
}