# GoForge project configuration
project_name: "{{.ProjectName}}"
module_path: "{{.ModuleName}}"
go_version: "{{.GoVersion}}"

# Project metadata
description: "A Go application built with GoForge"
author: ""
license: "MIT"

# Dependencies with version constraints
dependencies:
  github.com/gin-gonic/gin: "^1.10.0"
  github.com/spf13/viper: "^1.19.0"
  github.com/jackc/pgx/v5: "^5.6.0"

# Development dependencies
dev_dependencies:
  github.com/stretchr/testify: "^1.10.0"
  github.com/golang/mock: "^1.6.0"

# Custom scripts for project automation
scripts:
  # Development
  dev: "go run ./cmd/server"
  dev:watch: "goforge watch dev"
  
  # Building
  build: "goforge build"
  build:prod: "go build -ldflags='-w -s' -o dist/{{.ProjectName}} ./cmd/server"
  
  # Testing
  test: "goforge test"
  test:race: "goforge test --race"
  test:all: "goforge test --coverage --race"
  
  # Code quality
  lint: "golangci-lint run"
  fmt: "go fmt ./..."
  vet: "go vet ./..."
  
  # Database
  db:migrate: "migrate -path ./migrations -database postgres://localhost/{{.ProjectName}}_db up"
  db:rollback: "migrate -path ./migrations -database postgres://localhost/{{.ProjectName}}_db down 1"
  
  # Deployment
  docker:build: "docker build -t {{.ProjectName}} ."
  docker:run: "docker run -p 8080:8080 {{.ProjectName}}"

# Build configuration
build:
  # Output directory for build artifacts
  output_dir: "dist"
  
  # Binary name (defaults to project name)
  binary_name: "{{.ProjectName}}"
  
  # Assets to copy to output directory
  assets:
    - "config/default.yml"
    - "web/static"
    - "templates"
    
  # Cross-compilation targets
  targets:
    - os: "linux"
      arch: "amd64"
    - os: "windows" 
      arch: "amd64"
    - os: "darwin"
      arch: "amd64"
    - os: "darwin"
      arch: "arm64"

# Development server configuration
dev:
  # Port for development server
  port: 8080
  
  # Files/directories to watch for changes
  watch:
    - "**/*.go"
    - "config/**/*.yml"
    - "web/templates/**/*"
  
  # Files/directories to ignore
  ignore:
    - "dist/**"
    - "**/*_test.go"
    - ".git/**"
    - "node_modules/**"
  
  # Commands to run on file changes
  on_change:
    - "go fmt ./..."
    - "go vet ./..."

# Testing configuration
test:
  # Test timeout
  timeout: "10m"
  
  # Coverage threshold (percentage)
  coverage_threshold: 80
  
  # Test database configuration
  database:
    driver: "postgres"
    dsn: "postgres://test:test@localhost/{{.ProjectName}}_test?sslmode=disable"

# Code generation settings
generate:
  # Default component templates
  templates:
    handler: "templates/components/handler.go.tpl"
    service: "templates/components/service.go.tpl"
    repository: "templates/components/repository.go.tpl"
    model: "templates/components/model.go.tpl"
    middleware: "templates/components/middleware.go.tpl"
  
  # Output directories for generated components
  output:
    handler: "internal/adapters/http/handler"
    service: "internal/app/service"
    repository: "internal/adapters/postgres"
    model: "internal/domain"
    middleware: "internal/adapters/http/middleware"

# Docker configuration
docker:
  # Base image for multi-stage build
  base_image: "golang:1.24-alpine"
  runtime_image: "alpine:latest"
  
  # Exposed port
  port: 8080
  
  # Environment variables
  env:
    GIN_MODE: "release"
    LOG_LEVEL: "info"

# Database migration settings
migrations:
  # Directory containing migration files
  dir: "migrations"
  
  # Migration table name
  table: "schema_migrations"
  
  # Database connection for migrations
  database_url: "postgres://localhost/{{.ProjectName}}_db?sslmode=disable"
