package routes

import (
	"flight_records/internal/entities"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (a *AppRouter) track(ctx *fiber.Ctx) error {
	pairs, err := bindFlightPairs(ctx)
	if err != nil {
		return fiber.ErrBadRequest
	}
	return ctx.JSON(a.appService.Track(pairs))
}

func bindFlightPairs(ctx *fiber.Ctx) ([]entities.FlightPair, error) {
	var routes [][]string
	if err := ctx.BodyParser(&routes); err != nil {
		return nil, err
	}
	pairs := make([]entities.FlightPair, 0, len(routes))
	for i := range routes {
		pairs = append(pairs, entities.FlightPair{
			Source:      strings.ToLower(routes[i][0]),
			Destination: strings.ToLower(routes[i][1]),
		})
	}
	return pairs, nil
}
