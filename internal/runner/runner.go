package runner

import (
	"os"
	"os/exec"
)

// ExecuteCommand runs an external command within a specified directory,
// connecting its output and error streams to the user's terminal.
func ExecuteCommand(dir, name string, args...string) error {
	cmd := exec.Command(name, args...)
	// Set the working directory for the command.
	cmd.Dir = dir
	// Connect the command's standard streams to the parent process's streams.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ExecuteScript runs a shell script string by invoking "sh -c".
// This allows for complex shell commands with pipes, environment variables, etc.
func ExecuteScript(dir, script string) error {
	return ExecuteCommand(dir, "sh", "-c", script)
}

// InitGoModule runs 'go mod init' in the specified directory.
func InitGoModule(dir, modulePath string) error {
	return ExecuteCommand(dir, "go", "mod", "init", modulePath)
}

// TidyGoModule runs 'go mod tidy' in the specified directory.
func TidyGoModule(dir string) error {
	return ExecuteCommand(dir, "go", "mod", "tidy")
}

// InitGitRepository runs 'git init' in the specified directory, creating a 'main' branch by default.
func InitGitRepository(dir string) error {
	return ExecuteCommand(dir, "git", "init", "-b", "main")
}