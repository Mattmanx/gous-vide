package hardware

import (
	"os/exec"
	"strconv"
	"strings"
)

func CurrentTempCelsius() (float64, error) {
	return getTemp()
}

func CurrentTempFahrenheit() (float64, error) {
	tempCelsius, e := getTemp()

	if e != nil {
		return -1, e
	} else {
		return toFahrenheit(tempCelsius), nil
	}
}

// Gets the temp in celsius from an assumed 'get-temp' command on the system path.
func getTemp() (float64, error) {
	cmd := exec.Command("get-temp")

	output, e := cmd.Output()

	if e != nil {
		return -1, e
	} else {
		s := strings.TrimSpace(string(output[:]))

		f, err := strconv.ParseFloat(s, 64)

		return f, err
	}
}

// Helper function to convert celsius temp to fahrenheit
func toFahrenheit(tempCelcius float64) float64 {
	return tempCelcius * 1.8 + 32
}