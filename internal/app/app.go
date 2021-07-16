package app

import (
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

	var err error

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error(err, "app - Run - httpServer.Notify")
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - httpServer.Shutdown")
	}

}
