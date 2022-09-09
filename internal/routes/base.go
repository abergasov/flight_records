package routes

import (
	"errors"
	appLog "flight_records/internal/logger"
	"flight_records/internal/service/tracker"
	"net/http"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type AppRouter struct {
	appPort    string
	appService tracker.Tracker
	fiberApp   *fiber.App
}

// InitAppRouter initializes the app router.
func InitAppRouter(appPort string, logger appLog.AppLogger, appService tracker.Tracker) *AppRouter {
	log := logger.With(zap.String("module", "routes"))
	fiberApp := fiber.New(
		fiber.Config{
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				log.Error("error occurred", err, zap.String("path", ctx.Path()))
				if errors.Is(err, fiber.ErrBadRequest) {
					return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
				}
				return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
			},

			DisableStartupMessage: true,
		},
	)

	fiberApp.Use(recover.New())

	app := &AppRouter{
		appPort:    appPort,
		appService: appService,
		fiberApp:   fiberApp,
	}
	app.initRoutes()
	return app
}

func (a *AppRouter) initRoutes() {
	a.fiberApp.Get("/ping", a.pong)
	a.fiberApp.Post("/api/v1/track", a.track)
}

func (a *AppRouter) pong(ctx *fiber.Ctx) error {
	return ctx.SendString("pong")
}

// Run starts the server.
func (a *AppRouter) Run() error {
	return a.fiberApp.Listen(":" + a.appPort)
}

// Shutdown gracefully shuts down the server.
func (a *AppRouter) Shutdown() error {
	return a.fiberApp.Shutdown()
}
