# GoForge CLI

<p align="center">
  <a href="https://github.com/night-slayer18/goforge/actions/workflows/ci.yml">
    <img src="https://github.com/night-slayer18/goforge/actions/workflows/ci.yml/badge.svg" alt="Go CI">
  </a>
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Architecture-Clean-green?style=for-the-badge" alt="Clean Architecture">
  <img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License">
</p>

A powerful, NestJS-inspired CLI tool for scaffolding and managing Go applications with Clean Architecture principles. GoForge helps you focus on business logic instead of boilerplate setup and configuration.

## Table of Contents

- [✨ Features](#-features)
- [🚀 Quick Start](#-quick-start)
- [📁 Project Structure](#-project-structure)
- [🔧 Commands Reference](#-commands-reference)
- [📝 Configuration](#-configuration)
- [🏛️ Clean Architecture](#️-clean-architecture)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## ✨ Features

- 🏗️ **Clean Architecture** - Scaffolds projects following Clean Architecture principles
- 🚀 **Quick Setup** - Get a production-ready Go API in seconds
- 🔧 **Code Generation** - Generate handlers, services, repositories, models, and more
- 📦 **Dependency Management** - Simple dependency addition and updates
- 🔄 **Hot Reload** - Built-in file watching for development, configurable via `goforge.yml`
- 🛠️ **Build System** - Integrated, cross-platform build and asset management
- 📋 **Script Runner** - Custom script execution like npm scripts
- 🎯 **Type Safety** - Generate type-safe interfaces and implementations

## 🚀 Quick Start

### Installation

```bash
# Install GoForge CLI
go install github.com/night-slayer18/goforge@latest

# Verify installation
goforge --version
```

### Create Your First Project

```bash
# Create a new project
goforge new my-api

# Navigate to the project
cd my-api

# Start development server with hot reload
goforge watch
```

Your API server will start at `http://localhost:8080` with a basic health check endpoint.

## 📁 Project Structure

GoForge creates a well-organized project structure following Clean Architecture:

```
my-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── adapters/
│   │   ├── http/
│   │   │   ├── handler/         # HTTP handlers (controllers)
│   │   │   └── middleware/      # HTTP middleware
│   │   ├── postgres/            # Database repositories
│   │   └── database/            # Database connection
│   ├── app/
│   │   └── service/             # Business logic services
│   ├── domain/                  # Domain entities/models
│   └── ports/                   # Interface definitions
├── config/
│   └── default.yml              # Application configuration
├── dist/                        # Build output
├── go.mod
├── go.sum
├── goforge.yml                  # Project configuration
├── README.md
└── .gitignore
```

## 🔧 Commands Reference

### Project Management

#### Create New Project
```bash
# Basic project
goforge new my-project

# With custom module path
goforge new user-service --module-path github.com/myorg/user-service

# Skip Git initialization
goforge new simple-app --skip-git

# Use interactive mode
goforge new -i
```

#### Clean Project
```bash
# Remove build artifacts
goforge clean

# Clean everything including Go module cache
goforge clean --all
```

### Code Generation

Generate various application components with the `generate` command (alias: `g`):

```bash
# Generate a handler
goforge generate handler user

# Generate a service (short alias)
goforge g service auth

# Generate a repository (short alias)
goforge g r product
```
*(See `goforge generate --help` for all available components)*


### Development Workflow

#### Run Scripts
```bash
# Start development server (defined in goforge.yml)
goforge run dev

# Build the project
goforge run build

# Run tests
goforge run test
```

#### File Watching
```bash
# Watch for changes and auto-restart the 'dev' script
goforge watch
```

#### Building
```bash
# Build production binary and copy assets
goforge build
```

### Dependency Management

#### Add Dependencies
```bash
# Add latest version and tidy modules
goforge add github.com/gin-gonic/gin

# Add specific version
goforge add github.com/stretchr/testify@v1.8.4
```

#### Update Dependencies
```bash
# Update all dependencies defined in goforge.yml
goforge update

# Update a specific dependency to latest
goforge update github.com/gin-gonic/gin
```

## 📝 Configuration

### goforge.yml

The `goforge.yml` file is the heart of your project configuration:

```yaml
# GoForge project configuration
project_name: "my-api"
module_path: "github.com/myorg/my-api"
go_version: "1.24.5"

# Dependencies with version constraints
dependencies:
  github.com/gin-gonic/gin: "^1.10.0"
  github.com/spf13/viper: "^1.19.0"

# Custom scripts for project automation
scripts:
  dev: "go run ./cmd/server"
  build: "goforge build"
  test: "go test ./..."
  lint: "golangci-lint run"

# Build configuration
build:
  output_dir: "dist"
  assets:
    - "config/default.yml"

# Development server configuration
dev:
  watch:
    - "**/*.go"
    - "config/**/*.yml"
  ignore:
    - "dist/**"
    - "**/*_test.go"
```

### Application Configuration

Configure your application in `config/default.yml`:

```yaml
server:
  host: "localhost"
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "myapp_db"
  sslmode: "disable"
```

## 🏛️ Clean Architecture

GoForge follows Clean Architecture principles with clear separation of concerns. The dependency rule is strictly enforced: source code dependencies can only point inwards.

<p align="center">
  <img src="https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg" alt="Clean Architecture Diagram" width="500">
</p>

### Layers

1.  **Domain** (`internal/domain/`) - Core business entities and rules.
2.  **Ports** (`internal/ports/`) - Interfaces defining contracts for services and repositories.
3.  **Application** (`internal/app/service/`) - Business logic and use cases, orchestrating domain objects.
4.  **Adapters** (`internal/adapters/`) - Connectors to external systems like databases, web frameworks, etc.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/night-slayer18/goforge.git
cd goforge

# Install dependencies
go mod tidy

# Build the CLI with a version number
VERSION="1.2.0"
go build -ldflags="-X 'github.com/night-slayer18/goforge/cmd.version=$VERSION'" -o goforge .

# Test with local binary
./goforge --version
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
