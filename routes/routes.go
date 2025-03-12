package routes

import (
	"backend/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Health check
	e.GET("/health", handlers.HealthCheck)

	// API group
	api := e.Group("/api/v1")

	// Example routes
	api.GET("/tasks", handlers.GetAllTask)
	api.POST("/task", handlers.CreateTask)
	api.GET("/task/:id", handlers.GetTaskById)
	api.DELETE("/task/:id", handlers.DeleteTask)
	api.PUT("/task", handlers.EditTask)
}
