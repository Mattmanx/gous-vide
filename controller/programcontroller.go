package controller
import (
	"github.com/mattmanx/gous-vide/routine"
	"github.com/mattmanx/gous-vide/hardware"
	"net/http"
	"github.com/mattmanx/gous-vide/model"
	"io/ioutil"
	"io"
	"encoding/json"
	"sync"
	"fmt"
)

type ProgramController struct {
	program *routine.Program
	heater *hardware.Heater

	sync.RWMutex
}

func NewProgramController(heater *hardware.Heater) *ProgramController {
	return &ProgramController{heater: heater}
}

//TODO: Currently created in started state... update to check query param, allow start from id
func (c *ProgramController) NewProgram(w http.ResponseWriter, req *http.Request) {
	//expect a recipe here
	var recipe model.Recipe

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 500000))

	if err != nil {
		panic(err)
	}

	if err := req.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &recipe); err != nil {
		respondMessage(http.StatusNotAcceptable, w, err.Error())
	}

	//check status of existing program
	c.Lock()
	defer c.Unlock()

	if(c.program != nil && c.program.IsRunning()) {
		respondMessage(http.StatusPreconditionFailed, w, "Program is currently running. Unable to start new program.")
		return
	}

	c.program = routine.NewProgram(recipe, c.heater)
	err = c.program.Start(60000)	//TODO: hardcoded for now.. pull from optional query param in future

	if err != nil {
		respondMessage(http.StatusInternalServerError, w, fmt.Sprintf("Error starting program: %v", err))
		return
	}

	respondMessage(http.StatusOK, w, "Successfully started program.")	//TODO: return program run id
}

// Writes the current program status to the write stream, or returns a 404 if no current program
func (c *ProgramController) CurrentProgram(w http.ResponseWriter, req *http.Request) {
	//check status of existing program
	c.Lock()
	defer c.Unlock()

	if c.program == nil {
		respondMessage(http.StatusNotFound, w, "No current program")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)


	if err := json.NewEncoder(w).Encode(c.program.GetStatus()); err != nil {
		panic(err)
	}
}

// Stops the current program
func (c *ProgramController) StopCurrent(w http.ResponseWriter, req *http.Request) {
	//check status of existing program
	c.Lock()
	defer c.Unlock()

	if c.program == nil {
		respondMessage(http.StatusNotFound, w, "No current program")
		return
	}

	err := c.program.Stop()
	c.program = nil

	if err != nil {
		respondMessage(http.StatusInternalServerError, w, fmt.Sprintf("Error while stopping program: %v", err))
	}

	respondMessage(http.StatusOK, w, "Stopped program!")
}

// TODO: Add ability to lookup and view program history, lookup by ID, etc.!!!


// Controller interface function, returns a list of routes handled by this controller
func (c *ProgramController) GetRoutes() Routes {
	routes := Routes{
		Route{Name: "GetCurrentProgramStatus", Method: "GET", Pattern: "/program/current", HandlerFunc: c.CurrentProgram},
		Route{Name: "StopCurrentProgram", Method: "PUT", Pattern: "/program/current/stop", HandlerFunc: c.StopCurrent},
		Route{Name: "NewProgram", Method: "POST", Pattern: "/program", HandlerFunc: c.NewProgram},
	}

	return routes
}
