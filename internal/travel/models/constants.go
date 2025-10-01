package models

var (
	Hostels = map[string]struct{}{
		"Uniworld-1": {},
		"Uniworld-2": {},
	}

	AirportTerminals = map[string]struct{}{
		"Kempegowda International Airport Terminal-1": {},
		"Kempegowda International Airport Terminal-2": {},
	}

	RailwayStations = map[string]struct{}{
		"KSR SBC Bengaluru Junction Railway Station": {},
		"SMVT Bengaluru Railway station":             {},
		"Krishnarajapuram Railway Station":           {},
		"Yesvantpur Junction Railway station":        {},
		"Banglore Cantonment Railway Station":        {},
		"Bengaluru East Railway Station":             {},
	}
)

func IsHostel(loc string) bool {
	_, ok := Hostels[loc]
	return ok
}

func IsAirportTerminal(loc string) bool {
	_, ok := AirportTerminals[loc]
	return ok
}

func IsRailwayStation(loc string) bool {
	_, ok := RailwayStations[loc]
	return ok
}

// AreNearbyTerminals returns true when two airport terminals are considered close enough
// to be interchangeable for matching purposes (e.g. Terminal1 and Terminal2).
func AreNearbyTerminals(a, b string) bool {
	if a == b {
		return true
	}
	return IsAirportTerminal(a) && IsAirportTerminal(b)
}

// AreNearbyHostels returns true when two hostels are considered close enough
// to be interchangeable (e.g. Uniworld1 and Uniworld2).
func AreNearbyHostels(a, b string) bool {
	if a == b {
		return true
	}
	return IsHostel(a) && IsHostel(b)
}
