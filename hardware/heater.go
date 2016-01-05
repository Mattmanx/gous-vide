package hardware

import (
	"log"
	"os/exec"

)

type heater struct {
	on bool
}

func NewHeater() heater {
	h := heater {on: false}

	return h
}

func (h* heater) TurnOn() error {
	log.Printf("Turning heater on.")

	// Expect 'heater-on' to be on the system path
	cmd := exec.Command("heater-on")

	if e := cmd.Run(); e != nil {
		log.Printf("Error occurred while turning on heater.")
		return e
	} else {
		h.on = true
		log.Printf("Heater turned on.")
		return nil
	}
}

func (h* heater) TurnOff() error {
	log.Printf("Turning heater off.")

	// Expect 'heater-off' to be on the system path
	cmd := exec.Command("heater-off")

	if e := cmd.Run(); e != nil {
		log.Printf("Error occurred while turning off heater.")
		return e
	} else {
		h.on = false
		log.Printf("Heater turned off.")
		return nil
	}
}

func (h* heater) IsOn() bool {
	return h.on
}
