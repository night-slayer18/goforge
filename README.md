# GoForge CLI

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Architecture-Clean-green?style=for-the-badge" alt="Clean Architecture">
  <img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License">
</p>

A powerful, NestJS-inspired CLI tool for scaffolding and managing Go applications with Clean Architecture principles. GoForge helps you focus on business logic instead of boilerplate setup and configuration.

## ✨ Features

- 🏗️ **Clean Architecture** - Scaffolds projects following Clean Architecture principles
- 🚀 **Quick Setup** - Get a production-ready Go API in seconds
- 🔧 **Code Generation** - Generate handlers, services, repositories, models, and more
- 📦 **Dependency Management** - Simple dependency addition and updates
- 🔄 **Hot Reload** - Built-in file watching for development
- 🛠️ **Build System** - Integrated build and asset management
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
goforge run dev
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

# With verbose output
goforge new my-project --verbose
```

#### Clean Project
```bash
# Remove build artifacts
goforge clean

# Clean everything including Go module cache
goforge clean --all

# Dry run (see what would be removed)
goforge clean --dry-run
```

### Code Generation

Generate various application components with the `generate` command (alias: `g`):

#### HTTP Handlers
```bash
goforge generate handler user
goforge g handler auth
goforge g h product          # Short alias
```
*Creates: `internal/adapters/http/handler/user_handler.go`*

#### Business Services
```bash
goforge generate service user
goforge g service email
goforge g s notification     # Short alias
```
*Creates: `internal/app/service/user_service.go`*

#### Data Repositories
```bash
goforge generate repository user
goforge g repository product
goforge g repo order         # Alias
goforge g r customer         # Short alias
```
*Creates: `internal/adapters/postgres/user_repo.go`*

#### Domain Models
```bash
goforge generate model user
goforge g model product
goforge g mod order          # Short alias
```
*Creates: `internal/domain/user.go`*

#### HTTP Middleware
```bash
goforge generate middleware auth
goforge g middleware cors
goforge g m logging          # Short alias
```
*Creates: `internal/adapters/http/middleware/auth.go`*

#### Port Interfaces
```bash
goforge generate port user
goforge g port notification
goforge g p email            # Short alias
```
*Creates: `internal/ports/user_port.go`*

### Development Workflow

#### Run Scripts
```bash
# Start development server
goforge run dev

# Build the project
goforge run build

# Run tests
goforge run test

# Custom scripts (defined in goforge.yml)
goforge run lint
goforge run db:migrate
```

#### File Watching
```bash
# Watch for changes and auto-restart
goforge watch

# Watch specific script
goforge watch dev
goforge watch test
```

#### Building
```bash
# Build production binary
goforge build

# Binary will be created in dist/ directory
```

### Dependency Management

#### Add Dependencies
```bash
# Add latest version
goforge add github.com/gin-gonic/gin

# Add specific version
goforge add github.com/stretchr/testify@v1.8.0
```

#### Update Dependencies
```bash
# Update all dependencies
goforge update

# Update specific dependency
goforge update github.com/gin-gonic/gin
```

## 📝 Configuration

### goforge.yml

The `goforge.yml` file is the heart of your project configuration:

```yaml
# Project metadata
project_name: "my-api"
module_path: "github.com/myorg/my-api"
go_version: "1.24.5"

# Dependencies
dependencies:
  github.com/gin-gonic/gin: "^1.10.0"
  github.com/spf13/viper: "^1.19.0"
  github.com/jackc/pgx/v5: "^5.6.0"

# Custom scripts
scripts:
  dev: "go run ./cmd/server"
  build: "goforge build"
  test: "go test ./..."
  lint: "golangci-lint run"
  fmt: "go fmt ./..."
  db:migrate: "migrate -path ./migrations -database postgres://localhost/mydb up"

# Build configuration
build:
  output_dir: "dist"
  binary_name: "my-api"
  assets:
    - "config/default.yml"
    - "web/static"
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

logging:
  level: "debug"
```

## 🏛️ Clean Architecture

GoForge follows Clean Architecture principles with clear separation of concerns:

### Layers

1. **Domain** (`internal/domain/`) - Core business entities
2. **Ports** (`internal/ports/`) - Interface definitions
3. **Application** (`internal/app/service/`) - Business logic
4. **Adapters** (`internal/adapters/`) - External interfaces

### Dependency Flow

```
HTTP Handler → Service → Repository Interface ← Repository Implementation
      ↓            ↓              ↑                        ↓
   Adapters → Application ← Ports Domain ←          Adapters
```

### Example Component Flow

```bash
# 1. Generate domain model
goforge g model user

# 2. Generate port interface
goforge g port user

# 3. Generate repository implementation
goforge g repository user

# 4. Generate business service
goforge g service user

# 5. Generate HTTP handler
goforge g handler user
```

## 🚀 Complete Example

Here's a complete example of building an e-commerce API:

```bash
# 1. Create the project
goforge new ecommerce-api -m github.com/mycompany/ecommerce-api
cd ecommerce-api

# 2. Generate domain models
goforge g model product
goforge g model user
goforge g model order

# 3. Generate port interfaces
goforge g port product
goforge g port user
goforge g port order

# 4. Generate repositories
goforge g repository product
goforge g repository user
goforge g repository order

# 5. Generate services
goforge g service product
goforge g service user
goforge g service order

# 6. Generate handlers
goforge g handler product
goforge g handler user
goforge g handler order

# 7. Generate middleware
goforge g middleware auth
goforge g middleware cors

# 8. Add dependencies
goforge add github.com/golang-jwt/jwt/v4
goforge add github.com/go-playground/validator/v10

# 9. Start development
goforge watch
```

## 🛠️ Advanced Usage

### Custom Scripts

Add custom scripts to your `goforge.yml`:

```yaml
scripts:
  # Development
  dev: "go run ./cmd/server"
  dev:debug: "dlv debug ./cmd/server"
  
  # Testing
  test: "go test ./..."
  test:coverage: "go test -coverprofile=coverage.out ./..."
  test:race: "go test -race ./..."
  
  # Database
  db:up: "migrate -path ./migrations -database $DATABASE_URL up"
  db:down: "migrate -path ./migrations -database $DATABASE_URL down"
  db:seed: "go run ./scripts/seed.go"
  
  # Deployment
  docker:build: "docker build -t myapp ."
  deploy:staging: "./scripts/deploy.sh staging"
  deploy:prod: "./scripts/deploy.sh production"
```

### Multiple Environments

Organize configuration by environment:

```
config/
├── default.yml      # Base configuration
├── development.yml  # Development overrides
├── staging.yml      # Staging overrides
└── production.yml   # Production overrides
```

### Database Integration

GoForge generates PostgreSQL-ready repositories. To integrate:

1. **Add database dependencies:**
```bash
goforge add github.com/jackc/pgx/v5
goforge add github.com/golang-migrate/migrate/v4
```

2. **Update your main.go:**
```go
// Uncomment database connection in cmd/server/main.go
dbPool := database.Connect()
defer dbPool.Close()

userRepo := postgres.NewPostgresUserRepository(dbPool)
```

3. **Create migrations:**
```sql
-- migrations/001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/night-slayer18/goforge.git
cd goforge

# Install dependencies
go mod tidy

# Build the CLI
go build -o bin/goforge .

# Test with local binary
./bin/goforge new test-project
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [NestJS CLI](https://nestjs.com/) for Node.js
- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [Gin](https://github.com/gin-gonic/gin) for HTTP routing
- Follows [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) principles

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/night-slayer18">night-slayer18</a>
</p>