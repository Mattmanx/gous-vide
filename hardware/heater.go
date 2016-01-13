package hardware

import (
	"log"
	"os/exec"
	"sync"

)

type Heater struct {
	on   bool

	sync.RWMutex
}

func NewHeater() *Heater {
	return &Heater
}

func (h* Heater) TurnOn() error {
	log.Printf("Turning heater on.")

	h.Lock()
	defer h.Unlock()

	// Expect 'heater-on' to be on the system path
	cmd := exec.Command("heater-on")

	if e := cmd.Run(); e != nil {
		log.Printf("Error occurred while turning on heater.")
		return e
	}

	h.on = true
	log.Printf("Heater turned on.")
	return nil
}

func (h* Heater) TurnOff() error {
	log.Printf("Turning heater off.")

	h.Lock()
	defer h.Unlock()

	// Expect 'heater-off' to be on the system path
	cmd := exec.Command("heater-off")

	if e := cmd.Run(); e != nil {
		log.Printf("Error occurred while turning off heater.")
		return e
	}

	h.on = false
	log.Printf("Heater turned off.")

	return nil
}

func (h* Heater) IsOn() bool {
	h.RLock()
	defer h.RUnlock()

	return h.on
}
