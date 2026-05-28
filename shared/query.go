package shared

func FilterByCountry(stations []Station, ISOCode string) []Station {
	var filteredStationsByCountry []Station

	for _, station := range stations {
		if station.Country == ISOCode {
			filteredStationsByCountry = append(filteredStationsByCountry, station)
		}
	}

	return filteredStationsByCountry
}

func AvgTemperature(station Station) float64 {
	if len(station.Observations) == 0 {
		return 0
	}

	var total float64

	for _, observation := range station.Observations {
		total += observation.Temperature
	}

	return total / float64(len(station.Observations))
}

func MaxWindGust(stations []Station) (Station, float64) {
	var maxWindGustStation Station
	var maxWindGust float64

	for _, station := range stations {
		for _, observation := range station.Observations {
			if observation.Wind.Speed > maxWindGust {
				maxWindGustStation = station
				maxWindGust = observation.Wind.Speed
			}
		}
	}

	return maxWindGustStation, maxWindGust
}

func CountByCountry(stations []Station) map[string]int {
	stationsCountByCountry := make(map[string]int)

	for _, station := range stations {
		stationsCountByCountry[station.Country]++
	}

	return stationsCountByCountry
}
