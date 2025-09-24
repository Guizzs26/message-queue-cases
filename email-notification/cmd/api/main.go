package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health-check", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK!")
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":8080"
	}

	log.Println("Starting API server on port", port)
	if err := e.Start(port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("shutting down the server: %v", err)
	}
}
