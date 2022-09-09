package tracker

import (
	"flight_records/internal/entities"
	"strings"
)

type RouteTracker struct {
}

func InitRouteTracker() *RouteTracker {
	return &RouteTracker{}
}

func (r *RouteTracker) Track(route []entities.FlightPair) entities.FlightPair {
	return r.calculateRouteMap(route)
}

// calculateRouteMap calculates the route using 2 maps.
// Assuming that user not return back to airport that he already visited.
// In that key we loop over routes, and save soruce and destination airports to maps.
// If in any map this airport exist - this is transit, remove it from map.
// At the end we have 2 maps with source and destination airports.
func (r *RouteTracker) calculateRouteMap(route []entities.FlightPair) entities.FlightPair {
	if len(route) == 0 {
		return entities.FlightPair{}
	} else if len(route) == 1 {
		return route[0]
	}
	airportSource := make(map[string]struct{}, len(route))
	airportDestination := make(map[string]struct{}, len(route))
	for _, pair := range route {
		if _, ok := airportSource[pair.Destination]; ok {
			delete(airportSource, pair.Destination)
		} else {
			airportDestination[pair.Destination] = struct{}{}
		}
		if _, ok := airportDestination[pair.Source]; ok {
			delete(airportDestination, pair.Source)
		} else {
			airportSource[pair.Source] = struct{}{}
		}
	}
	return entities.FlightPair{
		Source:      strings.ToUpper(r.getMapKey(airportSource)),
		Destination: strings.ToUpper(r.getMapKey(airportDestination)),
	}
}

func (r *RouteTracker) getMapKey(dataMap map[string]struct{}) string {
	for key := range dataMap {
		return key
	}
	return ""
}
