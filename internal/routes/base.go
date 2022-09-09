package routes

import (
	"flight_records/internal/service/tracker"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type AppRouter struct {
	appPort    string
	appService tracker.Tracker
	fiberApp   *fiber.App
}

// InitAppRouter initializes the app router.
func InitAppRouter(appPort string, appService tracker.Tracker) *AppRouter {
	fiberApp := fiber.New(
		fiber.Config{
			//ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			//	log.Error("error get info. path: `%s`, err: %s", ctx.Request().URI().PathOriginal(), err)
			//	if errors.Is(err, service.ErrNotFound) {
			//		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			//	} else if errors.Is(err, fiber.ErrBadRequest) {
			//		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
			//	}
			//	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
			//},
			//DisableStartupMessage: true,
			//ReadBufferSize:        4096 * 16,
			//WriteBufferSize:       4096 * 16,
			//BodyLimit:             100 * 1024 * 1024,
		},
	)

	fiberApp.Use(recover.New())
	fiberApp.Use(logger.New())

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
