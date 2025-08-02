package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"{{.ModuleName}}/internal/adapters/database"
	"{{.ModuleName}}/internal/adapters/http/handler"
	"{{.ModuleName}}/internal/adapters/postgres"
	"{{.ModuleName}}/internal/app/service"
	"{{.ModuleName}}/internal/ports"
)

func main() {
	// --- Configuration Setup ---
	viper.SetConfigName("default")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8080 // Default port
	}

	// --- Dependency Injection ---
	// This section sets up the dependencies for the application.
	// It includes the database connection, repositories, services, and handlers.
	// Each component is initialized and wired together to form the application.
	// The order of initialization is important to ensure that dependencies are available 
	// ---------------------------------------------------------------
	/* 1. Database Connection
	dbPool := database.Connect()
	defer dbPool.Close() // Ensure the connection pool is closed on exit

	// 2. Repositories (Adapters)
	userRepo := postgres.NewPostgresUserRepository(dbPool)
	*/

	// In a real app, initialize a database connection here.
	var userRepo ports.UserRepository // = nil for this example

	// 3. Services (Application Core)
	userService := service.NewUserService(userRepo)

	// 4. Handlers (Adapters)
	userHandler := handler.NewUserHandler(userService)

	// --- Gin Router Setup ---
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		userRoutes := api.Group("/users")
		{
			userRoutes.GET("/:id", userHandler.GetUser)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// --- Start Server ---
	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("🚀 Server starting on http://localhost%s\n", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Could not start server: %v", err)
	}
}