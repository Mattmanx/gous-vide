package model

import (
	"time"
	"fmt"
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

// Simple convenience method to convert from one temperature scale to another
func ConvertTemperature(temperature float64, from TemperatureScale, to TemperatureScale) float64 {
	if from == to {
		return temperature
	}

	switch {
	case from == to:
		return temperature
	case from == CELSIUS && to == FAHRENHEIT:
		return temperature * 9 / 5 + 32
	case from == FAHRENHEIT && to == CELSIUS:
		return (temperature - 32) * 5 / 9
	default:
		panic(fmt.Errorf("Unknown temperature scale combination, %v to %v", from, to))
	}
}
