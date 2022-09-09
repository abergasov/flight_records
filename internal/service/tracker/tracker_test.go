package tracker_test

import (
	"flight_records/internal/entities"
	"flight_records/internal/service/tracker"
	"testing"
)

type testCase struct {
	name string
	src  []entities.FlightPair
	res  entities.FlightPair
}

func TestRouteTracker_Track(t *testing.T) {
	rt := tracker.InitRouteTracker()
	table := []testCase{
		{
			name: "empty pair",
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
		{
			name: "multiply pairs with dublicates",
			src: []entities.FlightPair{
				// c -> b, b -> d, d -> c, c -> a
				{"c", "b"},
				{"b", "d"},
				{"d", "c"},
				{"c", "a"},
			},
			res: entities.FlightPair{"C", "A"},
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

func TestRouteTracker_CalculateRouteOverMaps(t *testing.T) {
	rt := tracker.InitRouteTracker()
	table := []testCase{
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
			res := rt.CalculateRouteOverMaps(tc.src)
			if *res != tc.res {
				t.Errorf("expected %v, got %v", tc.res, res)
			}
		})
	}
}

func TestRouteTracker_CalculateRouterFor2Pairs(t *testing.T) {
	rt := tracker.InitRouteTracker()
	table := []testCase{
		{
			name: "correct order",
			src: []entities.FlightPair{
				{"SFO", "IND"},
				{"IND", "EWR"},
			},
			res: entities.FlightPair{"SFO", "EWR"},
		},
		{
			name: "wrong order",
			src: []entities.FlightPair{
				{"IND", "EWR"},
				{"SFO", "IND"},
			},
			res: entities.FlightPair{"SFO", "EWR"},
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			res := rt.CalculateRouterFor2Pairs(tc.src)
			if res != tc.res {
				t.Errorf("expected %v, got %v", tc.res, res)
			}
		})
	}
}

func TestRouteTracker_CalculateRouterFor3Pairs(t *testing.T) {
	rt := tracker.InitRouteTracker()
	table := []testCase{
		{
			name: "correct order",
			src: []entities.FlightPair{
				{"SFO", "IND"},
				{"IND", "EWR"},
				{"EWR", "SPB"},
			},
			res: entities.FlightPair{"SFO", "SPB"},
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			res := rt.CalculateRouterFor3Pairs(tc.src)
			if res != tc.res {
				t.Errorf("expected %v, got %v", tc.res, res)
			}
		})
	}
}

func TestRouteTracker_CalculateRouteWithDuplicates(t *testing.T) {
	rt := tracker.InitRouteTracker()
	table := []testCase{
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
		{
			name: "multiply pairs with dublicates",
			src: []entities.FlightPair{
				// c -> b, b -> d, d -> c, c -> a
				{"c", "b"},
				{"b", "d"},
				{"d", "c"},
				{"c", "a"},
			},
			res: entities.FlightPair{"C", "A"},
		},
		//{
		//	name: "multiply pairs circle",
		//	src: []entities.FlightPair{
		//		// a -> c, c -> e, e -> d, d -> a;
		//		// d -> a, a -> c, c -> e, e -> d;
		//		// c -> e, e -> d, d -> a, a -> c;
		//		// e -> d, d -> a, a -> c, c -> e;
		//		{"a", "c"},
		//		{"e", "d"},
		//		{"d", "a"},
		//		{"c", "e"},
		//	},
		//	res: entities.FlightPair{"A", "A"},
		//},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			res := rt.CalculateRouteWithDuplicates(tc.src)
			if res != tc.res {
				t.Errorf("expected %v, got %v", tc.res, res)
			}
		})
	}
}

func BenchmarkRouteTracker_CalculateRouteOverMaps(b *testing.B) {
	rt := tracker.InitRouteTracker()
	src := []entities.FlightPair{
		{"IND", "EWR"},
		{"SFO", "ATL"},
		{"GSO", "IND"},
		{"ATL", "GSO"},
	}
	for i := 0; i < b.N; i++ {
		rt.CalculateRouteOverMaps(src)
	}
}

func BenchmarkRouteTracker_CalculateRouteOverMaps_2Values(b *testing.B) {
	rt := tracker.InitRouteTracker()
	src := []entities.FlightPair{
		{"IND", "EWR"},
		{"GSO", "IND"},
	}
	for i := 0; i < b.N; i++ {
		rt.CalculateRouteOverMaps(src)
	}
}

func BenchmarkRouteTracker_CalculateRouterFor3Pairs(b *testing.B) {
	rt := tracker.InitRouteTracker()
	src := []entities.FlightPair{
		{"SFO", "IND"},
		{"IND", "EWR"},
		{"EWR", "SPB"},
	}
	for i := 0; i < b.N; i++ {
		rt.CalculateRouterFor3Pairs(src)
	}
}

func BenchmarkRouteTracker_CalculateRouterFor2Pairs(b *testing.B) {
	rt := tracker.InitRouteTracker()
	src := []entities.FlightPair{
		{"IND", "EWR"},
		{"GSO", "IND"},
	}
	for i := 0; i < b.N; i++ {
		rt.CalculateRouterFor2Pairs(src)
	}
}

func BenchmarkRouteTracker_CalculateRouteWithDuplicates(b *testing.B) {
	rt := tracker.InitRouteTracker()
	src := []entities.FlightPair{
		{"IND", "EWR"},
		{"SFO", "ATL"},
		{"GSO", "IND"},
		{"ATL", "GSO"},
	}
	for i := 0; i < b.N; i++ {
		rt.CalculateRouteWithDuplicates(src)
	}
}
