package tracker

import "flight_records/internal/entities"

type RouteTracker struct {
}

func InitRouteTracker() *RouteTracker {
	return &RouteTracker{}
}

func (rt *RouteTracker) Track(route []entities.FlightPair) entities.FlightPair {
	panic("implement me")
}
