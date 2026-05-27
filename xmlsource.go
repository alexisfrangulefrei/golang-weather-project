package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"
)

type xmlWeatherDataset struct {
	Version     string       `xml:"version,attr"`
	Source      string       `xml:"source,attr"`
	XmlStations []xmlStation `xml:"station"`
}

type xmlStation struct {
	Id              string           `xml:"id,attr"`
	Country         string           `xml:"country,attr"`
	Name            string           `xml:"name"`
	XmlCoordinates  xmlCoordinates   `xml:"coordinates"`
	XmlHardware     xmlHardware      `xml:"hardware"`
	XmlObservations []xmlObservation `xml:"observations>observation"`
}

type xmlCoordinates struct {
	Latitude  float64 `xml:"lat,attr"`      // Unit : Degree
	Longitude float64 `xml:"lon,attr"`      // Unit : Degree
	Altitude  int16   `xml:"altitude,attr"` // Unit : Meter
}

type xmlHardware struct {
	Vendor string `xml:"vendor,attr"`
	Model  string `xml:"model,attr"`
	Since  string `xml:"since,attr"`
}

type xmlObservation struct {
	At            string        `xml:"at,attr"` // Timestamp
	Sky           string        `xml:"sky,attr"`
	XmlMeasures   []xmlMeasure  `xml:"measure"`
	XmlWind       xmlWind       `xml:"wind"`
	XmlAirQuality xmlAirQuality `xml:"air_quality"`
	Note          *string       `xml:"note"`
}

type xmlMeasure struct {
	Type  string `xml:"type,attr"` // Type : temperature | humidity | pressure | precipitation
	Unit  string `xml:"unit,attr"` // Unit : Celsius | Percent | hPa | mm
	Value string `xml:",chardata"`
}

type xmlWind struct {
	Speed     float64 `xml:"speed,attr"`     // Unit : km/h
	Direction uint16  `xml:"direction,attr"` // Unit : Degree
}

type xmlAirQuality struct {
	XmlPollutant []xmlPollutant `xml:"pollutant"`
}

type xmlPollutant struct {
	Name  string `xml:"name,attr"` // Name : PM2.5 | PM10 | NO2
	Value string `xml:",chardata"`
}

func convertBytesToXmlStations(bytes []byte) (xmlWeatherDataset xmlWeatherDataset) {
	if err := xml.Unmarshal(bytes, &xmlWeatherDataset); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	return
}

func getXmlMeasureValue(measures []xmlMeasure, measureType string) (float64, bool) {
	measureTypeFounded := false

	for _, measure := range measures {
		if measure.Type == measureType {
			measureTypeFounded = true
			measuredValue, err := strconv.ParseFloat(measure.Value, 64)

			if err != nil {
				fmt.Println("Error:", err)
				panic(err)
			}

			return measuredValue, measureTypeFounded
		}
	}

	return 0, measureTypeFounded
}

func convertXmlStationsToStations(rawXMLStations xmlWeatherDataset) (stations []Station) {
	for _, rawXMLStation := range rawXMLStations.XmlStations {
		var observations []Observation

		for _, rawXMLObservation := range rawXMLStation.XmlObservations {
			// Info : On a pas besoin de faire de switch ici car notre modèle unifié contient uniquement la température comme dit dans l'énoncé de la partie A. De plus, le nous n'avons pas besoin du reste pour la partie E.
			temperature, ok := getXmlMeasureValue(rawXMLObservation.XmlMeasures, "temperature")

			if !ok {
				errorMessage := fmt.Sprintf("la température n'a pas été trouvée.")
				fmt.Println("Error:", errorMessage)
				panic(errorMessage)
			}

			observations = append(observations, Observation{
				Timestamp:    ParseStringToTime(rawXMLObservation.At, time.RFC3339),
				Temperature:  temperature,
				SkyCondition: rawXMLObservation.Sky,
				Wind: Wind{
					Speed:     rawXMLObservation.XmlWind.Speed,
					Direction: rawXMLObservation.XmlWind.Direction,
				},
				Note: rawXMLObservation.Note,
			})
		}

		stations = append(stations, Station{
			Id:       rawXMLStation.Id,
			Name:     rawXMLStation.Name,
			Country:  rawXMLStation.Country,
			Altitude: rawXMLStation.XmlCoordinates.Altitude,
			Device: Device{
				SensorType:       rawXMLStation.XmlHardware.Model,
				InstallationDate: ParseStringToTime(rawXMLStation.XmlHardware.Since, "2006-01-02"),
			},
			Coordinates: Coordinates{
				Latitude:  rawXMLStation.XmlCoordinates.Latitude,
				Longitude: rawXMLStation.XmlCoordinates.Longitude,
			},
			Observations: observations,
		})
	}

	return stations
}

func LoadFromXML(path string) ([]Station, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	rawXMLStations := convertBytesToXmlStations(bytes)

	stations := convertXmlStationsToStations(rawXMLStations)

	return stations, nil
}
