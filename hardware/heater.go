package hardware

import (
	"log"
	"os/exec"
	"sync"

)

// TODO: Use a package-level channel to synchronize calls to turn on, turn off, and read temperature. No struct, no lock.
type Heater struct {
	on   bool
	lock *sync.RWMutex
}

func NewHeater() *Heater {
	return &Heater{on: false,
		lock: &sync.RWMutex{}}
}

func (h* Heater) TurnOn() error {
	log.Printf("Turning heater on.")

	h.lock.Lock()
	defer h.lock.Unlock()

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

func (h* Heater) TurnOff() error {
	log.Printf("Turning heater off.")

	h.lock.Lock()
	defer h.lock.Unlock()

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

func (h* Heater) IsOn() bool {
	h.lock.RLock()
	defer h.lock.RUnlock()

	return h.on
}
