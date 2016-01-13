package model

import (
	"time"
)

type TemperatureScale int

const (
	CELSIUS TemperatureScale = iota
	FAHRENHEIT
)

// Represents a temperature reading at a specific point in time. This struct is unit-agnostic; it is up to the caller to
// determine the units (celsius, fahrenheit... kelvin?). If addtional metadata on the reading is required, it should be
// wrapped in a TemperatureSummary.
type TemperatureReading struct {
	Timestamp   time.Time    `json:"timestamp"`
	Temperature float64      `json:"temperature"`
}

// Wraps a slice of temperature readings and metadata describing those temperatures.
type TemperatureSummary struct {
	Readings []TemperatureReading 	`json:"readings"`
	Scale    TemperatureScale 		`json:"scale"`
}
