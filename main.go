package main

import (
	"flag"
	"fmt"
	"github.com/mattmanx/gous-vide/hardware"
	"github.com/mattmanx/gous-vide/server"
)

func main() {
	//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//	})

	//	router := NewRouter()
	//
	//	log.Fatal(http.ListenAndServe(":8080", router))

	serve := flag.Bool("serve", false, "Run in 'serve' mode, access the gous-vide temp control application over HTTP")
	port := flag.Int("port", 8080, "When in 'serve' mode, the HTTP port to listen on")

	command := flag.String("command", "get-temp", "When not in 'serve' mode, allows execution of: 'get-temp', 'heater-status', 'heater-on' or 'heater-off'");

	flag.Parse()

	if *serve {
		server.Start(*port)
	} else {
		heater := hardware.NewHeater()

		switch(*command) {
		case "get-temp":
			temp, e := hardware.CurrentTempCelsius()
			if e != nil {
				fmt.Errorf("Error when checking temperature %v", e)
			} else {
				fmt.Printf("Current temp celsius: %v", temp)
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

}
