package main

import (
	"log"
	"net/http"

	"goals-api/internal/db"
	"goals-api/internal/models"
	"goals-api/internal/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	database, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running migrations...")
	err = database.AutoMigrate(
		&models.MonthGoal{},
		&models.MonthTodo{},
		&models.WeekGoal{},
		&models.WeekTodo{},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrations complete")

	// register routes
	routes.Register(e, database)

	e.Start(":8080")
}
