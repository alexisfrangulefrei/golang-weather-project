package main

import (
	"fmt"

	"github.com/efrei/weather/internal"
)

func countObservations(stations []internal.Station) int {
	observationsCount := 0

	for _, station := range stations {
		observationsCount += len(station.Observations)
	}

	return observationsCount
}

func getBordeauxMerignacStationWithIndex(stations []internal.Station) (internal.Station, int) {
	var bordeauxMerignacStation internal.Station
	var bordeauxMerignacStationIndex int = -1

	for index, station := range stations {
		if station.Id == "FR-BOR-001" {
			bordeauxMerignacStationIndex = index
			bordeauxMerignacStation = station
		}
	}

	if bordeauxMerignacStationIndex == -1 {
		var messageError = "La station de Bordeaux Mérignac n'existe pas :"
		fmt.Println(messageError)
		panic(messageError)
	}

	return bordeauxMerignacStation, bordeauxMerignacStationIndex
}

func main() {
	stationsFromJson, err := internal.LoadFromJSON("sources/weather_data.json")

	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	stationsFromXml, err := internal.LoadFromXML("sources/weather_data.xml")

	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	jsonObservationCount := countObservations(stationsFromJson)
	xmlObservationCount := countObservations(stationsFromXml)

	fmt.Printf("JSON : %d stations, %d observations\n", len(stationsFromJson), jsonObservationCount)
	fmt.Printf("XML : %d stations, %d observations\n", len(stationsFromXml), xmlObservationCount)

	if len(stationsFromJson) == len(stationsFromXml) && jsonObservationCount == xmlObservationCount {
		fmt.Println("Cohérence : OK")
	} else {
		fmt.Println("Cohérence : KO")
	}

	jsonStationWithMaxWindGust, jsonMaxWindGust := internal.MaxWindGust(stationsFromJson)
	fmt.Printf("Station la plus ventée : %s (%.1f km/h)\n", jsonStationWithMaxWindGust.Id, jsonMaxWindGust)

	bordeauxMerignacStation, bordeauxMerignacStationIndex := getBordeauxMerignacStationWithIndex(stationsFromJson)
	jsonAvgTemperature := internal.AvgTemperature(stationsFromJson[bordeauxMerignacStationIndex])
	fmt.Printf("Temp. moyenne %s : %.1f °C\n", bordeauxMerignacStation.Name, jsonAvgTemperature)

	jsonCountByCountry := internal.CountByCountry(stationsFromJson)

	fmt.Print("Stations par pays : ")
	for key, value := range jsonCountByCountry {
		fmt.Printf("%s : %d | ", key, value)
	}

	/*

		// Tests de la partie C
		stationsFromJson, err := LoadFromJSON("weather_data.json")

		if err != nil {
			fmt.Println("Error:", err)
			panic(err)
		}

		fmt.Println(len(stationsFromJson))
		fmt.Println(len(stationsFromJson[0].Observations))

		// Tests de la partie D
		stationsFromXml, err := LoadFromXML("weather_data.xml")

		if err != nil {
			fmt.Println("Error:", err)
			panic(err)
		}

		fmt.Println(len(stationsFromXml))
		fmt.Println(len(stationsFromXml[0].Observations))

		// Tests de la partie E (from JSON)
		jsonFilteredStationsByCountry := FilterByCountry(stationsFromJson, "FR")
		for _, station := range jsonFilteredStationsByCountry {
			fmt.Println(station)
		}

		jsonAvgTemperature := AvgTemperature(stationsFromJson[0])
		fmt.Println(jsonAvgTemperature)

		jsonStationWithMaxWindGust, jsonMaxWindGust := MaxWindGust(stationsFromJson)
		fmt.Printf("Station: %v, Wind: %f\n", jsonStationWithMaxWindGust, jsonMaxWindGust)

		jsonCountByCountry := CountByCountry(stationsFromJson)
		for key, value := range jsonCountByCountry {
			fmt.Println(key, value)
		}

		// Tests de la partie E (from XML)
		xmlFilteredStationsByCountry := FilterByCountry(stationsFromXml, "FR")
		for _, station := range xmlFilteredStationsByCountry {
			fmt.Println(station)
		}

		xmlAvgTemperature := AvgTemperature(stationsFromXml[0])
		fmt.Println(xmlAvgTemperature)

		xmlStationWithMaxWindGust, xmlMaxWindGust := MaxWindGust(stationsFromXml)
		fmt.Printf("Station: %v, Wind: %f\n", xmlStationWithMaxWindGust, xmlMaxWindGust)

		xmlCountByCountry := CountByCountry(stationsFromXml)
		for key, value := range xmlCountByCountry {
			fmt.Println(key, value)
		}

	*/
}
