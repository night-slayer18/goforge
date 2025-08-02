package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"{{.ModulePath}}/internal/app/service"
)

// {{.NameTitle}}Handler handles HTTP requests related to the {{.Name}} resource.
type {{.NameTitle}}Handler struct {
	service *service.{{.NameTitle}}Service
}

// New{{.NameTitle}}Handler creates a new {{.NameTitle}}Handler.
func New{{.NameTitle}}Handler(s *service.{{.NameTitle}}Service) *{{.NameTitle}}Handler {
	return &{{.NameTitle}}Handler{
		service: s,
	}
}

// HandleSomething is an example handler method.
// TODO: Rename and implement your handler logic.
func (h *{{.NameTitle}}Handler) HandleSomething(c *gin.Context) {
	// 1. Parse request from c.Param, c.Query, or c.ShouldBindJSON.
	// 2. Call the service.
	// 3. Write response.
	c.JSON(http.StatusOK, gin.H{"message": "{{.NameTitle}} handler called"})
}
