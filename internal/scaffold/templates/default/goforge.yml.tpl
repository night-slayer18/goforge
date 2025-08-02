# The configuration file for your new Go project, managed by goforge.
project_name: "{{.ProjectName}}"
module_path: "{{.ModuleName}}"
go_version: "{{.GoVersion}}"

# Dependencies will be added here by 'goforge add <package>'
dependencies: {}

# Define custom commands that can be run with 'goforge run <script_name>'.
scripts:
  dev: go run ./cmd/server
  build: goforge build
  test: go test ./...

# Configure the build process.
build:
  # Specify non-Go assets to be copied to the 'dist' directory on build.
  assets:
    - "config/default.yml"
