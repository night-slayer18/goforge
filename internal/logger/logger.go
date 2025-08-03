// internal/logger/logger.go
package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger provides structured logging with colors and levels
type Logger struct {
	level  LogLevel
	writer io.Writer
	
	// Color functions
	debugColor *color.Color
	infoColor  *color.Color
	warnColor  *color.Color
	errorColor *color.Color
	successColor *color.Color
}

// Global logger instance
var globalLogger *Logger

func init() {
	globalLogger = New(INFO, os.Stdout)
}

// New creates a new logger instance
func New(level LogLevel, writer io.Writer) *Logger {
	return &Logger{
		level:  level,
		writer: writer,
		
		debugColor:   color.New(color.FgCyan),
		infoColor:    color.New(color.FgBlue),
		warnColor:    color.New(color.FgYellow),
		errorColor:   color.New(color.FgRed),
		successColor: color.New(color.FgGreen),
	}
}

// SetLevel sets the logging level
func SetLevel(level LogLevel) {
	globalLogger.level = level
}

// SetVerbose enables debug logging
func SetVerbose(verbose bool) {
	if verbose {
		SetLevel(DEBUG)
	} else {
		SetLevel(INFO)
	}
}

// Debug logs debug messages (only shown in verbose mode)
func Debug(format string, args ...interface{}) {
	globalLogger.Debug(format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", l.debugColor, format, args...)
	}
}

// Info logs informational messages
func Info(format string, args ...interface{}) {
	globalLogger.Info(format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", l.infoColor, format, args...)
	}
}

// Warn logs warning messages
func Warn(format string, args ...interface{}) {
	globalLogger.Warn(format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", l.warnColor, format, args...)
	}
}

// Error logs error messages
func Error(format string, args ...interface{}) {
	globalLogger.Error(format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.log("ERROR", l.errorColor, format, args...)
	}
}

// Success logs success messages with green color
func Success(format string, args ...interface{}) {
	globalLogger.Success(format, args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	if l.level <= INFO {
		l.log("SUCCESS", l.successColor, format, args...)
	}
}

// Step logs a step in a process with emoji
func Step(step int, total int, message string, args ...interface{}) {
	prefix := fmt.Sprintf("[%d/%d]", step, total)
	Info("🔄 %s %s", prefix, fmt.Sprintf(message, args...))
}

// Progress shows progress without newline
func Progress(message string, args ...interface{}) {
	fmt.Fprintf(globalLogger.writer, "\r⏳ %s", fmt.Sprintf(message, args...))
}

// Complete completes a progress line
func Complete(message string, args ...interface{}) {
	fmt.Fprintf(globalLogger.writer, "\r✅ %s\n", fmt.Sprintf(message, args...))
}

func (l *Logger) log(level string, colorFunc *color.Color, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	
	if colorFunc != nil {
		levelStr := colorFunc.Sprintf("%-7s", level)
		fmt.Fprintf(l.writer, "%s %s %s\n", 
			color.New(color.Faint).Sprint(timestamp),
			levelStr,
			message,
		)
	} else {
		fmt.Fprintf(l.writer, "%s %-7s %s\n", timestamp, level, message)
	}
}

// Command execution logging helpers
func CommandStart(cmd string, args ...string) {
	Debug("Executing command: %s %s", cmd, strings.Join(args, " "))
}

func CommandSuccess(cmd string, duration time.Duration) {
	Debug("Command '%s' completed successfully in %v", cmd, duration)
}

func CommandError(cmd string, err error, duration time.Duration) {
	Error("Command '%s' failed after %v: %v", cmd, duration, err)
}

// Project lifecycle logging
func ProjectCreationStart(name string) {
	Info("🚀 Creating new project '%s'...", name)
}

func ProjectCreationComplete(name string, duration time.Duration) {
	Success("✨ Project '%s' created successfully in %v", name, duration.Round(time.Millisecond))
	Info("")
	Info("Next steps:")
	Info("  cd %s", name)
	Info("  goforge run dev")
	Info("")
	Info("Happy coding! 🎉")
}

// File operations logging
func FileCreated(path string) {
	Debug("📄 Created file: %s", path)
}

func DirectoryCreated(path string) {
	Debug("📁 Created directory: %s", path)
}

func AssetCopied(src, dst string) {
	Debug("📦 Copied asset: %s -> %s", src, dst)
}

// Component generation logging
func ComponentGenerationStart(componentType, name string) {
	Info("🔧 Generating %s: %s", componentType, name)
}

func ComponentGenerationComplete(componentType, name, path string) {
	Success("✅ Generated %s '%s' at: %s", componentType, name, path)
}

// Dependency management logging
func DependencyAdding(module string) {
	Info("📦 Adding dependency: %s", module)
}

func DependencyAdded(module string) {
	Success("✅ Successfully added dependency: %s", module)
}

// Build process logging
func BuildStart(projectName string) {
	Info("🏗️  Building project '%s'...", projectName)
}

func BuildComplete(outputPath string, duration time.Duration) {
	Success("✅ Build completed successfully in %v", duration.Round(time.Millisecond))
	Info("📦 Binary created at: %s", outputPath)
}

// Error helpers
func FatalError(err error, message string, args ...interface{}) {
	Error(message, args...)
	if err != nil {
		Error("Details: %v", err)
	}
	os.Exit(1)
}

func ValidationError(field, value, message string, suggestions []string) {
	Error("❌ Validation failed for %s: %s", field, message)
	if value != "" {
		Error("   Value: '%s'", value)
	}
	if len(suggestions) > 0 {
		Info("")
		Info("💡 Suggestions:")
		for _, suggestion := range suggestions {
			Info("   • %s", suggestion)
		}
	}
}

// Progress indicator for long-running operations
type ProgressIndicator struct {
	message string
	done    chan bool
}

func NewProgress(message string) *ProgressIndicator {
	p := &ProgressIndicator{
		message: message,
		done:    make(chan bool),
	}
	p.start()
	return p
}

func (p *ProgressIndicator) start() {
	go func() {
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		
		for {
			select {
			case <-p.done:
				return
			default:
				fmt.Fprintf(globalLogger.writer, "\r%s %s", chars[i%len(chars)], p.message)
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

func (p *ProgressIndicator) Complete(message string) {
	p.done <- true
	fmt.Fprintf(globalLogger.writer, "\r✅ %s\n", message)
}

func (p *ProgressIndicator) Stop() {
	p.done <- true
	fmt.Fprintf(globalLogger.writer, "\r")
}