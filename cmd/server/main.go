package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	sqlConn, _ := database.DB()
	defer sqlConn.Close()

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

	// start echo in a goroutine
	server := &http.Server{
		Addr:    ":8080",
		Handler: e,
	}
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen %s\n", err)
		}
	}()

	// listening for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	// 10 second max wait if work is being done
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}
