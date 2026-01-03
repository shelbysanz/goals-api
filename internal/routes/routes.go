package routes

import (
	"goals-api/internal/handlers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(e *echo.Echo, db *gorm.DB) {
	mg := handlers.NewMonthGoalHandler(db)

	api := e.Group("/api")
	api.GET("/month-goals", mg.List)
	api.POST("/month-goals", mg.Create)
}
