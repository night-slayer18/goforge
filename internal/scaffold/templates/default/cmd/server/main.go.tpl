package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"{{.ModuleName}}/internal/adapters/http/handler"
	"{{.ModuleName}}/internal/app/service"
	"{{.ModuleName}}/internal/ports"
)

func main() {
	// --- Configuration Setup using Viper ---
	viper.SetConfigName("default")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8080 // Default port if not specified
	}

	// --- Dependency Injection ---
	// In a real app, initialize a database connection here.
	var userRepo ports.UserRepository // = nil for this example

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// --- Gin Router Setup ---
	router := gin.Default()

	// Group routes for better organization
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
