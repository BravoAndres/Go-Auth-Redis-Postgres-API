package app

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BravoAndres/fiber-api/internal/middleware"
	"github.com/BravoAndres/fiber-api/internal/routes"
	"github.com/BravoAndres/fiber-api/pkg/httpserver"
	"github.com/BravoAndres/fiber-api/pkg/logger"
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload"
)

func Run() {
	// Define a new Fiber app with config.
	app := fiber.New(fiber.Config{
		ReadTimeout: time.Second * time.Duration(60),
	})

	middleware.FiberMiddleware(app)

	routes.UserRoutes(app)
	routes.AuthRoutes(app)
	routes.NotFoundRoute(app)

	handler := app.Handler()
	httpServer := httpserver.NewServer(handler)

	logger.Info("Server started")

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	err := errors.New("")
	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error(err, "app - Run - httpServer.Notify")
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - httpServer.Shutdown")
	}

	// Create channel for idle connections.
	// idleConnsClosed := make(chan struct{})

	// go func() {
	// 	sigint := make(chan os.Signal, 1)
	// 	signal.Notify(sigint, os.Interrupt) // Catch OS signals.
	// 	<-sigint

	// 	// Received an interrupt signal, shutdown.
	// 	if err := app.Shutdown(); err != nil {
	// 		// Error from closing listeners, or context timeout:
	// 		logger.Error(err)

	// 		return
	// 	}

	// 	close(idleConnsClosed)
	// }()

	// // Run server.
	// if err := app.Listen(os.Getenv("SERVER_URL")); err != nil {
	// 	logger.Error(err)

	// 	return
	// }

	// logger.Info("Server started")

	// <-idleConnsClosed
}
