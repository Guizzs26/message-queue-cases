package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Guizzs26/message-queue/internal/core/notification"
	"github.com/Guizzs26/message-queue/internal/platform/broker"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	APIPort     string `envconfig:"API_PORT" default:":9919" validate:"required"`
	RabbitMQURL string `envconfig:"RABBITMQ_URL" validate:"required,url"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
}

type App struct {
	config   *Config
	echo     *echo.Echo
	broker   broker.Broker
	validate *validator.Validate
}

func NewApp(cfg *Config) (*App, error) {
	app := &App{
		config:   cfg,
		echo:     echo.New(),
		validate: validator.New(),
	}

	rmqb, err := broker.NewRabbitMQBroker(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rabbitmq broker: %w", err)
	}
	app.broker = rmqb

	app.setupMiddleware()
	app.setupRoutes()

	return app, nil
}

func (a *App) setupMiddleware() {
	a.echo.Use(middleware.RequestID())
	a.echo.Use(middleware.Logger())
	a.echo.Use(middleware.Recover())
	a.echo.Use(middleware.CORS())
}

func (a *App) setupRoutes() {
	a.echo.GET("/health", a.healthCheckHandler)

	// API v1 routes
	v1 := a.echo.Group("/api/v1")
	a.setupNotificationRoutes(v1)
}

// setupNotificationRoutes configures notification module routes
func (a *App) setupNotificationRoutes(g *echo.Group) {
	ns := notification.NewNotificationService(a.broker)
	nh := notification.NewNotificationHTTPHandler(ns)

	notificationGroup := g.Group("/notifications")
	notificationGroup.POST("/email", nh.SendEmailHandler)
}

func (a *App) healthCheckHandler(c echo.Context) error {
	health := map[string]any{
		"status":      "OK",
		"version":     "1.0.0",
		"environment": a.config.Environment,
		"timestamp":   time.Now().UTC(),
	}
	return c.JSON(http.StatusOK, health)
}

func (a *App) Start() error {
	log.Printf("starting API server on port %s in %s environment",
		a.config.APIPort, a.config.Environment)

	return a.echo.Start(a.config.APIPort)
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("shutting down server...")

	defer a.broker.Close()
	return a.echo.Shutdown(ctx)
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found or error loading it:", err)
	}

	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment config: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}

	return &cfg, nil
}

func gracefulShutdown(app *App) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// start server in a goroutine
	go func() {
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	// wait for shutdown signal
	<-sigChan
	log.Println("shutdown signal received")

	// create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server shutdown completed")
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// composition root instance (centralized app structure)
	app, err := NewApp(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create app: %v", err)
	}

	gracefulShutdown(app)
}
