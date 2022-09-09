package tracker

import "flight_records/internal/entities"

type Tracker interface {
	Track(route []entities.FlightPair) entities.FlightPair
}
