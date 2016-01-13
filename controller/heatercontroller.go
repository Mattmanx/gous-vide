package controller
import (
	"github.com/mattmanx/gous-vide/hardware"
	"net/http"
	"encoding/json"
)

type HeaterResponse struct {
	Message string	`json:"message"`
	Status string	`json:"status"`
}

type HeaterController struct {
	heater *hardware.Heater
}

func NewHeaterController(heater *hardware.Heater) *HeaterController {
	return &HeaterController{heater: heater}
}

func (c *HeaterController) TurnOn(w http.ResponseWriter, r *http.Request) {
	(*c.heater).TurnOn()

	isOn := (*c.heater).IsOn()

	respondSuccess(w, "Heater command processed", isOn)
}

func (c *HeaterController) TurnOff(w http.ResponseWriter, r *http.Request) {
	(*c.heater).TurnOff()

	isOn := (*c.heater).IsOn()

	respondSuccess(w, "Heater command processed", isOn)
}

func (c *HeaterController) GetStatus(w http.ResponseWriter, r *http.Request) {
	isOn := (*c.heater).IsOn()

	respondSuccess(w, "Heater status retrieved", isOn)
}

func (c* HeaterController) GetRoutes() Routes {
	routes := Routes {
		Route{Name: "HeaterOn", Method: "PUT", Pattern: "/heater/on", HandlerFunc: c.TurnOn},
		Route{Name: "HeaterOff", Method: "PUT", Pattern: "/heater/off", HandlerFunc: c.TurnOff},
		Route{Name: "HeaterStatus", Method: "GET", Pattern: "/heater", HandlerFunc: c.GetStatus},
	}

	return routes
}

func respondSuccess(w http.ResponseWriter, message string, isOn bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var status string

	if isOn {
		status = "ON"
	} else {
		status = "OFF"
	}

	if err := json.NewEncoder(w).Encode(HeaterResponse{Message: message, Status: status}); err != nil {
		panic(err)
	}
}

func respondError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(message); err != nil {
		panic(err)
	}
}