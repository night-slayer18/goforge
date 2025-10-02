package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/night-slayer18/goforge/internal/logger"
	"github.com/night-slayer18/goforge/internal/project"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watchCmd = &cobra.Command{
	Use:   "watch [script-name]",
	Short: "Watch for file changes and restart the application",
	Long: `The watch command monitors your project files for changes and automatically
restarts the specified script when changes are detected. If no script is specified,
it defaults to the 'dev' script.

GoForge handles all process management, port cleanup, and graceful restarts internally.
Your application code stays clean and simple.

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
		
		logger.Info("👀 Starting GoForge watch mode")
		logger.Info("📝 Script: %s → %s", scriptName, script)
		logger.Info("📁 Watching: %s", projectRoot)
		logger.Info("🔄 Press Ctrl+C to stop")
		logger.Info("")

		// Create the advanced watcher
		watcher := NewAdvancedWatcher(projectRoot, script, verbose, cfg)
		defer watcher.Close()

		// Set up graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start the watcher
		if err := watcher.Start(); err != nil {
			return fmt.Errorf("failed to start watcher: %w", err)
		}

		// Wait for shutdown signal
		<-sigChan
		logger.Info("\n🛑 Shutting down...")
		
		if err := watcher.Stop(); err != nil {
			logger.Error("Error during shutdown: %v", err)
		} else {
			logger.Info("✅ GoForge watch mode stopped")
		}

		return nil
	},
}

// AdvancedWatcher handles all the complexity of file watching and process management
type AdvancedWatcher struct {
	projectRoot    string
	script         string
	verbose        bool
	fileWatcher    *fsnotify.Watcher
	processManager *ProcessManager
	portManager    *PortManager
	debouncer      *Debouncer
	
	// Configuration from project
	projectPort    int
	watchPatterns  []string
	ignorePatterns []string
}

// NewAdvancedWatcher creates a new advanced watcher
func NewAdvancedWatcher(projectRoot, script string, verbose bool, cfg *project.Config) *AdvancedWatcher {
	watcher := &AdvancedWatcher{
		projectRoot: projectRoot,
		script:      script,
		verbose:     verbose,
		debouncer:   NewDebouncer(1500 * time.Millisecond), // Smart debouncing
	}
	
	watcher.loadProjectConfig(cfg)
	
	return watcher
}

// loadProjectConfig loads project-specific configuration
func (aw *AdvancedWatcher) loadProjectConfig(cfg *project.Config) {
	// Try to load viper config to detect port from config/default.yml
	viper.SetConfigName("default")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(aw.projectRoot, "config"))
	
	if err := viper.ReadInConfig(); err == nil {
		aw.projectPort = viper.GetInt("server.port")
	}
	
	if aw.projectPort == 0 {
		aw.projectPort = 8080 // Default
	}
	
	// Set up watch and ignore patterns
	// Use values from goforge.yml if available, otherwise use defaults.
	if cfg.Dev != nil && len(cfg.Dev.Watch) > 0 {
		aw.watchPatterns = cfg.Dev.Watch
		logger.Debug("Loaded %d watch patterns from goforge.yml", len(cfg.Dev.Watch))
	} else {
		aw.watchPatterns = []string{
			"**/*.go",
			"**/*.yml",
			"**/*.yaml",
			"**/*.json",
		}
		logger.Debug("Using default watch patterns")
	}
	
	if cfg.Dev != nil && len(cfg.Dev.Ignore) > 0 {
		aw.ignorePatterns = cfg.Dev.Ignore
		logger.Debug("Loaded %d ignore patterns from goforge.yml", len(cfg.Dev.Ignore))
	} else {
		aw.ignorePatterns = []string{
			"**/*_test.go",
			"dist/**",
			"vendor/**",
			".git/**",
			"node_modules/**",
			"**/*.tmp",
			"**/*.log",
		}
		logger.Debug("Using default ignore patterns")
	}
	
	logger.Debug("Detected project port: %d", aw.projectPort)
}

// Start begins watching and starts the initial process
func (aw *AdvancedWatcher) Start() error {
	var err error
	
	// Initialize file watcher
	aw.fileWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	
	// Initialize process manager
	aw.processManager = NewProcessManager(aw.projectRoot, aw.script, aw.verbose)
	
	// Initialize port manager
	aw.portManager = NewPortManager()
	
	// Add directories to watch
	if err := aw.addWatchPaths(); err != nil {
		return fmt.Errorf("failed to setup watch paths: %w", err)
	}
	
	// Start the initial process
	logger.Info("🚀 Starting initial process...")
	if err := aw.processManager.Start(); err != nil {
		return fmt.Errorf("failed to start initial process: %w", err)
	}
	
	// Start watching for file changes
	go aw.watchLoop()
	
	return nil
}

// watchLoop handles file change events
func (aw *AdvancedWatcher) watchLoop() {
	lastRestart := time.Now().Add(-10 * time.Second)
	
	for {
		select {
		case event, ok := <-aw.fileWatcher.Events:
			if !ok {
				return
			}
			
			if aw.shouldIgnoreEvent(event) {
				continue
			}
			
			logger.Debug("File changed: %s (%s)", event.Name, event.Op)
			
			// Prevent rapid restarts
			if time.Since(lastRestart) < 2*time.Second {
				logger.Debug("Ignoring change - too soon after last restart")
				continue
			}
			
			// Debounce the restart
			aw.debouncer.Debounce(func() {
				lastRestart = time.Now()
				logger.Info("🔄 Changes detected, restarting...")
				
				if err := aw.smartRestart(); err != nil {
					logger.Error("Failed to restart: %v", err)
				}
			})
			
		case err, ok := <-aw.fileWatcher.Errors:
			if !ok {
				return
			}
			logger.Warn("File watcher error: %v", err)
		}
	}
}

// smartRestart performs an intelligent restart with port management
func (aw *AdvancedWatcher) smartRestart() error {
	// Step 1: Stop the current process gracefully
	logger.Debug("Stopping current process...")
	if err := aw.processManager.Stop(); err != nil {
		logger.Warn("Error stopping process: %v", err)
	}
	
	// Step 2: Ensure port is available
	logger.Debug("Ensuring port %d is available...", aw.projectPort)
	if err := aw.portManager.EnsurePortAvailable(aw.projectPort, 8*time.Second); err != nil {
		logger.Warn("Port cleanup may have failed: %v", err)
		// Continue anyway - the process start might still work
	}
	
	// Step 3: Wait a moment for system cleanup
	time.Sleep(500 * time.Millisecond)
	
	// Step 4: Start new process
	logger.Debug("Starting new process...")
	if err := aw.processManager.Start(); err != nil {
		return fmt.Errorf("failed to start new process: %w", err)
	}
	
	logger.Success("✅ Process restarted successfully")
	return nil
}

// shouldIgnoreEvent determines if a file change event should be ignored
func (aw *AdvancedWatcher) shouldIgnoreEvent(event fsnotify.Event) bool {
	// Only care about write and create events
	if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
		return true
	}
	
	relPath, err := filepath.Rel(aw.projectRoot, event.Name)
	if err != nil {
		return true
	}
	
	// Check ignore patterns
	for _, pattern := range aw.ignorePatterns {
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return true
		}
		
		// Check if any directory in the path matches
		dirs := strings.Split(filepath.Dir(relPath), string(filepath.Separator))
		for _, dir := range dirs {
			if matched, _ := filepath.Match(strings.TrimSuffix(pattern, "/**"), dir); matched {
				return true
			}
		}
	}
	
	// Check if file matches watch patterns
	for _, pattern := range aw.watchPatterns {
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return false // Should watch this file
		}
		
		// Check extension-based patterns
		if strings.HasPrefix(pattern, "**/*") {
			ext := strings.TrimPrefix(pattern, "**/*")
			if strings.HasSuffix(relPath, ext) {
				return false // Should watch this file
			}
		}
	}
	
	return true // Ignore by default
}

// addWatchPaths recursively adds directories to the file watcher
func (aw *AdvancedWatcher) addWatchPaths() error {
	return filepath.Walk(aw.projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			return nil
		}
		
		relPath, err := filepath.Rel(aw.projectRoot, path)
		if err != nil {
			return err
		}
		
		// Check if directory should be ignored
		for _, pattern := range aw.ignorePatterns {
			dirPattern := strings.TrimSuffix(pattern, "/**")
			if matched, _ := filepath.Match(dirPattern, relPath); matched {
				logger.Debug("Ignoring directory: %s", relPath)
				return filepath.SkipDir
			}
		}
		
		logger.Debug("Watching directory: %s", relPath)
		return aw.fileWatcher.Add(path)
	})
}

// Stop stops the watcher and cleans up resources
func (aw *AdvancedWatcher) Stop() error {
	var errs []error
	
	if aw.processManager != nil {
		if err := aw.processManager.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("process manager: %w", err))
		}
	}
	
	if aw.fileWatcher != nil {
		if err := aw.fileWatcher.Close(); err != nil {
			errs = append(errs, fmt.Errorf("file watcher: %w", err))
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	
	return nil
}

// Close is an alias for Stop for consistency
func (aw *AdvancedWatcher) Close() error {
	return aw.Stop()
}

// ProcessManager handles process lifecycle with enhanced control
type ProcessManager struct {
	dir      string
	script   string
	verbose  bool
	cmd      *exec.Cmd
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewProcessManager creates a new process manager
func NewProcessManager(dir, script string, verbose bool) *ProcessManager {
	return &ProcessManager{
		dir:     dir,
		script:  script,
		verbose: verbose,
	}
}

// Start starts the process
func (pm *ProcessManager) Start() error {
	pm.ctx, pm.cancel = context.WithCancel(context.Background())
	
	pm.cmd = exec.CommandContext(pm.ctx, "sh", "-c", pm.script)
	pm.cmd.Dir = pm.dir
	
	// Set up process group for better control
	pm.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	
	if pm.verbose {
		pm.cmd.Stdout = os.Stdout
		pm.cmd.Stderr = os.Stderr
	} else {
		// Capture output for smart filtering
		stdout, err := pm.cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := pm.cmd.StderrPipe()
		if err != nil {
			return err
		}
		
		// Filter and display output
		go pm.handleOutput(stdout, false)
		go pm.handleOutput(stderr, true)
	}
	
	if err := pm.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}
	
	logger.Success("✅ Process started (PID: %d)", pm.cmd.Process.Pid)
	
	// Monitor process completion
	go func() {
		err := pm.cmd.Wait()
		if err != nil && pm.ctx.Err() == nil {
			// Process died unexpectedly (not due to cancellation)
			logger.Error("❌ Process exited unexpectedly: %v", err)
		}
	}()
	
	return nil
}

// handleOutput processes stdout/stderr with smart filtering
func (pm *ProcessManager) handleOutput(pipe io.Reader, isError bool) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Filter out noise
		if pm.shouldFilterLine(line) {
			continue
		}
		
		if isError {
			logger.Error("🔴 %s", line)
		} else {
			// Highlight important messages
			if pm.isImportantLine(line) {
				logger.Success("🟢 %s", line)
			} else {
				logger.Info("⚪ %s", line)
			}
		}
	}
}

// shouldFilterLine determines if a log line should be filtered out
func (pm *ProcessManager) shouldFilterLine(line string) bool {
	noisePatterns := []string{
		"[GIN-debug]",
		"Listening and serving HTTP",
	}
	
	for _, pattern := range noisePatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	
	return false
}

// isImportantLine determines if a log line is important
func (pm *ProcessManager) isImportantLine(line string) bool {
	importantPatterns := []string{
		"Server starting",
		"🚀",
		"✅",
		"Ready",
		"Started",
	}
	
	for _, pattern := range importantPatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	
	return false
}

// Stop stops the process gracefully
func (pm *ProcessManager) Stop() error {
	if pm.cmd == nil || pm.cmd.Process == nil {
		return nil
	}
	
	logger.Debug("Stopping process (PID: %d)...", pm.cmd.Process.Pid)
	
	// Cancel context first
	if pm.cancel != nil {
		pm.cancel()
	}
	
	// Get process group ID
	pgid, err := syscall.Getpgid(pm.cmd.Process.Pid)
	if err != nil {
		pgid = pm.cmd.Process.Pid
	}
	
	// Send SIGTERM to process group
	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
		// Fallback to killing just the process
		pm.cmd.Process.Signal(syscall.SIGTERM)
	}
	
	// Wait with timeout for graceful shutdown
	done := make(chan error, 1)
	go func() {
		_, err := pm.cmd.Process.Wait()
		done <- err
	}()
	
	select {
	case <-done:
		logger.Debug("Process stopped gracefully")
	case <-time.After(3 * time.Second):
		logger.Debug("Process didn't stop gracefully, force killing...")
		syscall.Kill(-pgid, syscall.SIGKILL)
		<-done // Wait for force kill to complete
	}
	
	pm.cmd = nil
	return nil
}

// PortManager handles port availability and cleanup
type PortManager struct{}

// NewPortManager creates a new port manager
func NewPortManager() *PortManager {
	return &PortManager{}
}

// EnsurePortAvailable ensures a port is available, with cleanup if necessary
func (pm *PortManager) EnsurePortAvailable(port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		if pm.isPortAvailable(port) {
			return nil
		}
		
		logger.Debug("Port %d still in use, attempting cleanup...", port)
		pm.attemptPortCleanup(port)
		
		time.Sleep(500 * time.Millisecond)
	}
	
	return fmt.Errorf("port %d is still not available after %v", port, timeout)
}

// isPortAvailable checks if a port is available
func (pm *PortManager) isPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// attemptPortCleanup tries to free up a port by calling a platform-specific implementation.
func (pm *PortManager) attemptPortCleanup(port int) {
	switch runtime.GOOS {
	case "linux", "darwin":
		pm.attemptPortCleanupUnix(port)
	case "windows":
		pm.attemptPortCleanupWindows(port)
	default:
		logger.Debug("Port cleanup not supported on this OS: %s", runtime.GOOS)
	}
}

// attemptPortCleanupUnix tries to free up a port on Unix-like systems.
func (pm *PortManager) attemptPortCleanupUnix(port int) {
	// Use lsof to find and kill processes using the port
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err != nil {
		return
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return
	}

	logger.Debug("Killing process %d using port %d", pid, port)
	if process, err := os.FindProcess(pid); err == nil {
		process.Signal(syscall.SIGTERM)

		// Wait a moment, then force kill if needed
		time.Sleep(1 * time.Second)
		if !pm.isPortAvailable(port) {
			process.Kill()
		}
	}
}

// attemptPortCleanupWindows tries to free up a port on Windows.
func (pm *PortManager) attemptPortCleanupWindows(port int) {
	portStr := fmt.Sprintf(":%d", port)
	findPidCmd := fmt.Sprintf("netstat -aon | findstr %s", portStr)

	out, err := exec.Command("cmd", "/C", findPidCmd).Output()
	if err != nil {
		logger.Debug("Could not find process for port %d: %v", port, err)
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 5 {
			continue
		}

		if strings.Contains(fields[1], portStr) {
			pidStr := fields[len(fields)-1]
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				logger.Debug("Could not parse PID from netstat output: %s", pidStr)
				continue
			}

			logger.Debug("Found process %d on port %d. Attempting to kill.", pid, port)

			killCmd := exec.Command("taskkill", "/F", "/PID", pidStr)
			if err := killCmd.Run(); err != nil {
				logger.Warn("Failed to kill process %d: %v", pid, err)
			} else {
				logger.Debug("Successfully sent kill signal to process %d", pid)
			}
			return
		}
	}
}

// Debouncer prevents rapid successive calls
type Debouncer struct {
	duration time.Duration
	timer    *time.Timer
}

// NewDebouncer creates a new debouncer
func NewDebouncer(duration time.Duration) *Debouncer {
	return &Debouncer{duration: duration}
}

// Debounce executes the function after the specified delay
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