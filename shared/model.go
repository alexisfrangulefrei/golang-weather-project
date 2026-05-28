package shared

import "time"

type Station struct {
	Id           string
	Name         string
	Country      string // Format : ISO 3166-1 alpha-2
	Altitude     int16  // Unit : Meter
	Device       Device
	Coordinates  Coordinates
	Observations []Observation
}

type Observation struct {
	Timestamp    time.Time
	Temperature  float64 // Unit : Celsius
	SkyCondition string
	Wind         Wind
	Note         *string
}

type Device struct {
	SensorType       string
	InstallationDate time.Time
}

type Coordinates struct {
	Latitude  float64 // Unit : Degree
	Longitude float64 // Unit : Degree
}

type Wind struct {
	Speed     float64 // Unit : km/h
	Direction uint16  // Unit : Degree
}
