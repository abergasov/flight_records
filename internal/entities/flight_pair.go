package entities

import "encoding/json"

type FlightPair struct {
	Source      string
	Destination string
}

func (u FlightPair) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{u.Source, u.Destination})
}
