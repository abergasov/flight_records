package tracker

import (
	"flight_records/internal/entities"
	"flight_records/internal/utils"
	"strings"
)

type RouteTracker struct {
}

func InitRouteTracker() *RouteTracker {
	return &RouteTracker{}
}

func (r *RouteTracker) Track(route []entities.FlightPair) entities.FlightPair {
	switch len(route) {
	case 0:
		return entities.FlightPair{}
	case 1:
		return route[0]
	case 2:
		return r.CalculateRouterFor2Pairs(route)
	}
	if res := r.CalculateRouteOverMaps(route); res != nil {
		return *res
	}
	return r.CalculateRouteWithDuplicates(route)
}

// CalculateRouteOverMaps calculates the route using 2 maps.
// Assuming that user NOT return to airport that he already visited.
// In that key we loop over routes, and save source and destination airports to maps.
// If in any map this airport exist - this is transit, remove it from map (destination|source can be only in one map).
// At the end we have 2 maps with source and destination airports.
func (r *RouteTracker) CalculateRouteOverMaps(route []entities.FlightPair) *entities.FlightPair {
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
	return &entities.FlightPair{
		Source:      strings.ToUpper(r.getMapKey(airportSource)),
		Destination: strings.ToUpper(r.getMapKey(airportDestination)),
	}
}

func (r *RouteTracker) CalculateRouterFor2Pairs(route []entities.FlightPair) (result entities.FlightPair) {
	if route[0].Destination == route[1].Source {
		return entities.FlightPair{strings.ToUpper(route[0].Source), strings.ToUpper(route[1].Destination)}
	}
	return entities.FlightPair{strings.ToUpper(route[1].Source), strings.ToUpper(route[0].Destination)}
}

// CalculateRouterFor3Pairs works only if order not shuffled. Method for hypothesis check in benchmark.
func (r *RouteTracker) CalculateRouterFor3Pairs(route []entities.FlightPair) (result entities.FlightPair) {
	if route[0].Destination == route[1].Source && route[1].Destination == route[2].Source {
		return entities.FlightPair{strings.ToUpper(route[0].Source), strings.ToUpper(route[2].Destination)}
	}
	panic("not implemented")
}

type routeContainer struct {
	result    []string
	sourceMap map[string][]string
}

// CalculateRouteWithDuplicates calculates the route with assuming that user CAN return to airport that he already visited.
func (r *RouteTracker) CalculateRouteWithDuplicates(route []entities.FlightPair) entities.FlightPair {
	perAirportSource := make(map[string][]string, len(route))
	for _, pair := range route {
		perAirportSource[pair.Source] = append(perAirportSource[pair.Source], pair.Destination)
	}

	var resul []string
	for airport := range perAirportSource {
		container := &routeContainer{
			result:    make([]string, 0, len(route)),
			sourceMap: utils.CopyMap(perAirportSource),
		}
		if r.buildRoute(airport, container, len(route)) {
			resul = container.result
			break
		}
	}
	return entities.FlightPair{
		Source:      strings.ToUpper(resul[0]),
		Destination: strings.ToUpper(resul[len(resul)-1]),
	}
}

func (r *RouteTracker) buildRoute(start string, source *routeContainer, expectedRouteLen int) bool {
	source.result = append(source.result, start)
	counter := 0
	totalDestinations := len(source.sourceMap[start])
	for len(source.sourceMap[start]) > 0 && counter < totalDestinations {
		destination := source.sourceMap[start][0] // take first destination from possible pairs src->dest
		source.sourceMap[start] = source.sourceMap[start][1:]
		if !r.buildRoute(destination, source, expectedRouteLen) { // wrong, rollback changes
			counter += 1
			source.result = source.result[:len(source.result)-1] // remove appended destination
			source.sourceMap[start] = append(source.sourceMap[start], destination)
		} else {
			return true
		}
	}
	// if not all airports included - len will be mismatched with route and return false
	// len-1 cause it contain points, not a pairs.
	return len(source.result)-1 == expectedRouteLen
}

func (r *RouteTracker) getMapKey(dataMap map[string]struct{}) string {
	for key := range dataMap {
		return key
	}
	return ""
}
