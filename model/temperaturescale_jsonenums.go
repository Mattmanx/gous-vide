// generated by jsonenums -type=TemperatureScale; DO NOT EDIT

package model

import (
	"encoding/json"
	"fmt"
)

var (
	_TemperatureScaleNameToValue = map[string]TemperatureScale{
		"CELSIUS":    CELSIUS,
		"FAHRENHEIT": FAHRENHEIT,
	}

	_TemperatureScaleValueToName = map[TemperatureScale]string{
		CELSIUS:    "CELSIUS",
		FAHRENHEIT: "FAHRENHEIT",
	}
)

func init() {
	var v TemperatureScale
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_TemperatureScaleNameToValue = map[string]TemperatureScale{
			interface{}(CELSIUS).(fmt.Stringer).String():    CELSIUS,
			interface{}(FAHRENHEIT).(fmt.Stringer).String(): FAHRENHEIT,
		}
	}
}

// MarshalJSON is generated so TemperatureScale satisfies json.Marshaler.
func (r TemperatureScale) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _TemperatureScaleValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid TemperatureScale: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so TemperatureScale satisfies json.Unmarshaler.
func (r *TemperatureScale) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("TemperatureScale should be a string, got %s", data)
	}
	v, ok := _TemperatureScaleNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid TemperatureScale %q", s)
	}
	*r = v
	return nil
}

func ScaleToString(scale TemperatureScale) string {
	return _TemperatureScaleValueToName[scale]
}

func StringToScale(text string) TemperatureScale {
	return _TemperatureScaleNameToValue[text]
}
