# Default application configuration.
# These values can be overridden by environment variables.
server:
  host: "localhost"
  port: 8080

logging:
  level: "debug" # Options: debug, info, warn, error

# Database connection settings for PostgreSQL.
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "{{.ProjectName}}_db"
  sslmode: "disable" # Use "require" or "verify-full" in production.