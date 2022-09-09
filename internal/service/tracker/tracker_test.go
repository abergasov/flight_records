package tracker

import (
	"flight_records/internal/entities"
	"testing"
)

func TestRouteTracker_Track(t *testing.T) {
	rt := InitRouteTracker()
	table := []struct {
		name string
		src  []entities.FlightPair
		res  entities.FlightPair
	}{
		{
			name: "empty pairs",
			src:  []entities.FlightPair{},
			res:  entities.FlightPair{},
		},
		{
			name: "one pair",
			src: []entities.FlightPair{
				{"SFO", "EWR"},
			},
			res: entities.FlightPair{"SFO", "EWR"},
		},
		{
			name: "two pairs",
			src: []entities.FlightPair{
				{"ATL", "EWR"},
				{"SFO", "ATL"},
			},
			res: entities.FlightPair{"SFO", "EWR"},
		},
		{
			name: "multiply pairs",
			src: []entities.FlightPair{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
			res: entities.FlightPair{"SFO", "EWR"},
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			res := rt.Track(tc.src)
			if res != tc.res {
				t.Errorf("expected %v, got %v", tc.res, res)
			}
		})
	}
}
