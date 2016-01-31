package main

import (
	"flag"
	"fmt"
	"github.com/mattmanx/gous-vide/hardware"
	"github.com/mattmanx/gous-vide/server"
	"github.com/mattmanx/gous-vide/db"
	"github.com/mattmanx/gous-vide/model"
)

var (
	serve = flag.Bool("serve", false, "Run in 'serve' mode, access the gous-vide temp control application over HTTP")
	port = flag.Int("port", 8080, "When in 'serve' mode, the HTTP port to listen on")
	command = flag.String("command", "get-temp", "When not in 'serve' mode, allows execution of: 'get-temp', 'heater-status', 'heater-on' or 'heater-off'");
)

func main() {
	flag.Parse()

	if *serve {
		server.Start(*port)
	}

	heater := hardware.NewHeater()

	switch(*command) {
	case "get-temp":
		temp, e := hardware.CurrentTempCelsius()
		if e != nil {
			fmt.Printf("Error when checking temperature %v", e)

			return
		}

		fmt.Printf("Current temp celsius: %v", temp)

		// save the temperature to the database
		e = db.Open()

		if e != nil {
			fmt.Printf("Error opening database: %v", e)
			return
		}

		defer db.Close()

		e = db.SaveTemperature(temp, model.CELSIUS)

		if e != nil {
			fmt.Printf("Error saving temperature to database: %v", e)
		}
	case "heater-on":
		heater.TurnOn()
	case "heater-off":
		heater.TurnOff()
	case "heater-status":
		fmt.Printf("heater is on: %v", heater.IsOn())
	default:
		fmt.Errorf("unknown command %s, must be one of 'get-temp', 'heater-status', 'heater-on', or 'heater-off'", *command)
	}
}
