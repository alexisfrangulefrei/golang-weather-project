package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type jsonMetadata struct {
	Version      string `json:"version"`
	Source       string `json:"source"`
	GeneratedAt  string `json:"generated_at"`
	StationCount int    `json:"station_count"`
	License      string `json:"license"`
}

type jsonStations struct {
	JsonStations []jsonStation `json:"stations"`
	JsonMetadata jsonMetadata  `json:"metadata"`
}

type jsonStation struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Country          string            `json:"country"`
	Altitude         int16             `json:"altitude_m"` // Unit : Meter
	JsonLocation     jsonLocation      `json:"location"`
	JsonDevice       jsonDevice        `json:"device"`
	JsonObservations []jsonObservation `json:"observations"`
}

type jsonLocation struct {
	Latitude  float64 `json:"latitude"`  // Unit : Degree
	Longitude float64 `json:"longitude"` // Unit : Degree
}

type jsonDevice struct {
	DeviceType   string `json:"type"`
	Manufacturer string `json:"manufacturer"`
	InstalledOn  string `json:"installed_on"` // Date format : YYYY-MM-DD
}

type jsonObservation struct {
	Timestamp      string         `json:"timestamp"`
	Temperature    float64        `json:"temperature_celsius"` // Unit : Celsius
	Humidity       uint8          `json:"humidity_percent"`    // Unit : Percent
	Pressure       float64        `json:"pressure_hpa"`        // Unit : hPa
	JsonWind       jsonWind       `json:"wind"`
	Precipitation  float64        `json:"precipitation_mm"` // Unit : mm
	JsonAirQuality jsonAirQuality `json:"air_quality"`
	Conditions     string         `json:"conditions"`
	Notes          *string        `json:"notes,omitempty"`
}

type jsonWind struct {
	Speed     float64 `json:"speed_kmh"`     // Unit : km/h
	Direction uint16  `json:"direction_deg"` // Unit : Degree
}

type jsonAirQuality struct {
	Pm25 float64 `json:"pm25"`
	Pm10 float64 `json:"pm10"`
	No2  float64 `json:"no2"`
}

func convertBytesToJsonStations(bytes []byte) (jsonStations jsonStations) {
	if err := json.Unmarshal(bytes, &jsonStations); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	return
}

func convertCountryToCountryISOCode(country string) string {
	countriesISOCode := map[string]string{
		"France":             "FR",
		"États-Unis":         "US",
		"Allemagne":          "DE",
		"Espagne":            "ES",
		"Italie":             "IT",
		"Japon":              "JP",
		"Chine":              "CN",
		"Brésil":             "BR",
		"Inde":               "IN",
		"Canada":             "CA",
		"Portugal":           "PT",
		"Belgique":           "BE",
		"Pays-Bas":           "NL",
		"Autriche":           "AT",
		"Suisse":             "CH",
		"Royaume-Uni":        "GB",
		"Irlande":            "IE",
		"Danemark":           "DK",
		"Suède":              "SE",
		"Norvège":            "NO",
		"Finlande":           "FI",
		"Pologne":            "PL",
		"République tchèque": "CZ",
		"Slovaquie":          "SK",
		"Hongrie":            "HU",
		"Roumanie":           "RO",
		"Bulgarie":           "BG",
		"Grèce":              "GR",
		"Croatie":            "HR",
		"Slovénie":           "SI",
		"Serbie":             "RS",
		"Bosnie-Herzégovine": "BA",
		"Monténégro":         "ME",
		"Albanie":            "AL",
		"Macédoine du Nord":  "MK",
		"Estonie":            "EE",
		"Lettonie":           "LV",
		"Lituanie":           "LT",
		"Luxembourg":         "LU",
		"Liechtenstein":      "LI",
		"Monaco":             "MC",
		"Andorre":            "AD",
		"Saint-Marin":        "SM",
		"Vatican":            "VA",
		"Ukraine":            "UA",
		"Moldavie":           "MD",
		"Biélorussie":        "BY",
		"Turquie":            "TR",
		"Tchéquie":           "CZ",
	}

	countryISOCode, ok := countriesISOCode[country]

	if !ok {
		errorMessage := fmt.Sprintf("le code ISO n'est pas défini pour le pays suivant : %s", country)
		fmt.Println("Error:", errorMessage)
		panic(errorMessage)
	}

	return countryISOCode
}

func convertJsonStationsToStations(rawJsonStationsFirstLevel jsonStations) (stations []Station) {
	for _, rawJsonStation := range rawJsonStationsFirstLevel.JsonStations {
		var observations []Observation

		for _, rawJsonStationObservation := range rawJsonStation.JsonObservations {
			observations = append(observations, Observation{
				Timestamp:    ParseStringToTime(rawJsonStationObservation.Timestamp, time.RFC3339),
				Temperature:  rawJsonStationObservation.Temperature,
				SkyCondition: rawJsonStationObservation.Conditions,
				Wind: Wind{
					Speed:     rawJsonStationObservation.JsonWind.Speed,
					Direction: rawJsonStationObservation.JsonWind.Direction,
				},
				Note: rawJsonStationObservation.Notes,
			})
		}

		stations = append(stations, Station{
			Id:       rawJsonStation.Id,
			Name:     rawJsonStation.Name,
			Country:  convertCountryToCountryISOCode(rawJsonStation.Country),
			Altitude: rawJsonStation.Altitude,
			Device: Device{
				SensorType:       rawJsonStation.JsonDevice.DeviceType,
				InstallationDate: ParseStringToTime(rawJsonStation.JsonDevice.InstalledOn, "2006-01-02"),
			},
			Coordinates: Coordinates{
				Latitude:  rawJsonStation.JsonLocation.Latitude,
				Longitude: rawJsonStation.JsonLocation.Longitude,
			},
			Observations: observations,
		})
	}

	return stations
}

func LoadFromJSON(path string) ([]Station, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	rawJsonStationsFirstLevel := convertBytesToJsonStations(bytes)

	stations := convertJsonStationsToStations(rawJsonStationsFirstLevel)

	return stations, nil
}
