package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/night-slayer18/goforge/internal/logger"
)

// {{.NameTitle}}Middleware provides {{.Name}} functionality.
type {{.NameTitle}}Middleware struct {
	// TODO: Add any dependencies here
	// Example:
	// authService *service.AuthService
}

// New{{.NameTitle}}Middleware creates a new {{.NameTitle}}Middleware.
func New{{.NameTitle}}Middleware( /* dependencies */ ) *{{.NameTitle}}Middleware {
	return &{{.NameTitle}}Middleware{
		// Initialize dependencies
	}
}

// Handler returns the Gin middleware handler function.
func (m *{{.NameTitle}}Middleware) Handler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		
		// TODO: Implement your middleware logic here
		// Example for authentication middleware:
		// token := c.GetHeader("Authorization")
		// if token == "" {
		//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		//     c.Abort()
		//     return
		// }
		
		logger.Debug("{{.NameTitle}} middleware processing request: %s %s", c.Request.Method, c.Request.URL.Path)
		
		// Continue to next handler
		c.Next()
		
		// Post-processing logic (optional)
		duration := time.Since(start)
		logger.Debug("{{.NameTitle}} middleware completed in %v", duration)
	})
}

// Apply is a convenience method to apply this middleware to a router group.
func (m *{{.NameTitle}}Middleware) Apply(rg *gin.RouterGroup) {
	rg.Use(m.Handler())
}