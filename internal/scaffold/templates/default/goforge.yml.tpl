# The configuration file for your new Go project, managed by goforge.
project_name: "{{.ProjectName}}"
module_path: "{{.ModuleName}}"
go_version: "{{.GoVersion}}"

# Default dependencies for a robust backend.
dependencies:
  github.com/gin-gonic/gin: v1.10.0
  github.com/spf13/viper: v1.19.0

# Define custom commands that can be run with 'goforge run <script_name>'.
scripts:
  dev: go run ./cmd/server
  build: goforge build
  test: go test ./...

# Configuration for the 'goforge build' command.
build:
  assets:
    - "config/default.yml"
