# GoForge CLI

**GoForge** is a powerful, opinionated Command-Line Interface (CLI) for Go, heavily inspired by the development workflow of the NestJS CLI. It's designed to help you initialize, develop, and maintain scalable and maintainable backend applications with ease.

GoForge accelerates development by scaffolding your project with a production-ready, **Clean Architecture** (also known as Ports and Adapters), allowing you to focus on business logic instead of boilerplate and configuration.

## Core Features

* **Rapid Scaffolding:** Create a complete, feature-rich backend project with a single command (`goforge new`).
* **Best-Practice Architecture:** Generates a project pre-configured with a Clean/Hexagonal architecture, ensuring separation of concerns, testability, and maintainability from day one.
* **Code Generation:** Quickly scaffold new components like HTTP handlers and services with the `generate` command.
* **Integrated Tooling:** Comes with a built-in web server (Gin), configuration management (Viper), and a pre-configured database layer (pgx).
* **Dependency Management:** A simple `add` command to manage your project's Go module dependencies.
* **Custom Script Runner:** Define and run common tasks like `dev`, `test`, or `build` using a `goforge.yml` manifest, similar to `npm run`.

## Installation

Ensure you have Go (version 1.18 or newer) installed and your `GOPATH` is set up correctly. You can then install `goforge` globally with a single command:

```bash
go install github.com/night-slayer18/goforge@latest
```

Verify the installation by checking the version:

```bash
goforge --version
```

## Commands Usage

GoForge provides a suite of commands to assist you through the entire development lifecycle.

### `goforge new`

This command scaffolds a new, complete Go project in a new directory.

**Syntax:**

```bash
goforge new <project-name> [flags]
```

**Arguments:**

* `<project-name>`: The name of your new project (e.g., `my-api`).

**Flags:**

* `--module-path`, `-m`: Explicitly set the Go module path. This is **crucial** for ensuring correct import paths if you plan to host your code publicly (e.g., on GitHub). If omitted, it defaults to `<project-name>`.
* `--go`: Specify a Go version for the `go.mod` file. Defaults to the version of Go you are currently using.

**Examples:**

**1. Basic Usage (Default Module Path)**
This will create a project with the module path `github.com/night-slayer18/order-service`.

```bash
goforge new order-service
```

**2. Specifying a Custom GitHub Module Path (Recommended)**
Use the `-m` flag to set the exact import path for your repository.

```bash
goforge new my-blog-api -m github.com/user-name/my-blog-api
```

### `goforge generate` (or `g`)

This command generates new application components from templates.

**Syntax:**

```bash
goforge generate <component> <name>
```

**Arguments:**

* `<component>`: The type of component to generate (e.g., `handler`, `service`).
* `<name>`: The name of the resource (e.g., `product`, `user`).

**Example:**

```bash
# First, navigate into your project directory
cd my-blog-api

# Generate a new handler for a "product" resource
goforge generate handler product

# Generate a new service for an "invoice" resource using the alias
goforge g service invoice
```

### `goforge run`

Executes a custom script defined in the `scripts` section of your `goforge.yml` file.

**Syntax:**

```bash
goforge run <script-name>
```

**Example:**

```bash
# Run the development server
goforge run dev

# Run the project's tests
goforge run test
```

### `goforge add`

Adds a new dependency to your project. It automatically runs `go get` and updates your `goforge.yml` manifest.

**Syntax:**

```bash
goforge add <module-path>[@version]
```

**Example:**

```bash
# Add the popular 'go-playground/validator' package
goforge add github.com/go-playground/validator/v10
```

### `goforge build`

Compiles your application into a single, production-ready binary and copies any specified assets into a `dist/` directory.

**Syntax:**

```bash
goforge build
```

**Example:**

```bash
goforge build
```

## The Generated Project Architecture

When you run `goforge new`, you get more than just files—you get a well-defined architecture designed for scale and maintainability.

* **/cmd/server/main.go**: The application's entry point. Its only job is to wire up all the dependencies (database, repositories, services, handlers) and start the server.
* **/config/**: Holds your application's configuration files (e.g., `default.yml`).
* **/internal/domain/**: The core of your application. Contains your business models (e.g., the `User` struct). This layer has zero external dependencies.
* **/internal/app/**: Contains the application services that orchestrate your business logic (use cases).
* **/internal/ports/**: Defines the interfaces (the "Ports") that your application core uses to talk to the outside world (e.g., `UserRepository`).
* **/internal/adapters/**: Contains the concrete implementations (the "Adapters") for the ports. This is where specific technologies like Gin, PostgreSQL (`pgx`), etc., live.
* **goforge.yml**: The project manifest used by the `goforge` CLI to manage dependencies and scripts.

This structure ensures your core business logic is completely decoupled from your web framework and database, making it incredibly easy to test, maintain, and even migrate to new technologies in the future.

## Example Workflow: Building a Simple API

Here's how you can go from zero to a running API in just a few minutes.

1.  **Create the project with a specific module path:**
    ```bash
    goforge new my-blog-api -m github.com/user-name/my-blog-api
    ```

2.  **Navigate into the project:**
    ```bash
    cd my-blog-api
    ```

3.  **Run the development server:**
    ```bash
    goforge run dev
    # Output: 🚀 Server starting on http://localhost:8080
    ```
    You can now visit `http://localhost:8080/health` in your browser to see the health check endpoint.

4.  **Generate a new "post" resource:**
    ```bash
    goforge g service post
    goforge g handler post
    ```
    This creates `post_service.go` and `post_handler.go`.

5.  **Manually wire up the new route in `cmd/server/main.go`:**
    Open `cmd/server/main.go` and add the new handler and route. This is a manual step that ensures you have full control over your application's setup.

6.  **Build for production:**
    ```bash
    goforge build
    ```
    Your final executable is now ready in the `dist/` folder.

## Contributing

Contributions are welcome! If you have ideas for new features, find a bug, or want to improve the documentation, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
See the [LICENSE](LICENSE) file for details.