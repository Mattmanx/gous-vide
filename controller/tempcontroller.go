package controller
import (
	"net/http"
	"time"
	"github.com/mattmanx/gous-vide/db"
	"github.com/mattmanx/gous-vide/model"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/mattmanx/gous-vide/hardware"
	"github.com/mattmanx/gous-vide/routine"
	"io/ioutil"
	"io"
)

const (
	NARROW_RFC3339 = "2006-01-02T15:04"
)

type PollRequest struct {
	Action string	`json:"action"`
	IntervalMilliseconds int	`json:"intervalMilliseconds"`
}

type TempController struct {
	poller *routine.TempPoller
}

// Creates a new TempController to handle requests related to temperature capture and recording.
// TODO: Probably need to inject the temp poller in the future so we can ensure the same poller is used for programmed recipes
func NewTempController() *TempController {
	return &TempController{poller: routine.NewTempPoller()}
}

func (c *TempController) GetCurrentTemperature(w http.ResponseWriter, req *http.Request) {
	save := false

	var err error
	var currentTemp float64

	//query params: save?
	q := req.URL.Query()
	saveText := q.Get("save")

	if len(saveText) > 0 {
		save, err = strconv.ParseBool(saveText)
	}

	if err != nil {
		respondMessage(http.StatusBadRequest, w, fmt.Sprintf("Error parsing query param 'save'. Expected 'true' or 'false'. Error: %v", err))
		return
	}

	currentTemp, err = hardware.CurrentTempCelsius()

	if err != nil {
		respondMessage(http.StatusInternalServerError, w, fmt.Sprintf("Error retrieving current temperature. Error: %v", err))
		return
	}

	if save {
		err = db.SaveTemperature(currentTemp, model.CELSIUS)
	}

	if err != nil {
		respondMessage(http.StatusInternalServerError, w, fmt.Sprintf("Error saving current temperature. Error: %v", err))
		return
	}

	reading := model.TemperatureReading{time.Now(), currentTemp}  //TODO: get our record from the db.save call so timestamp matches what actually gets saved
	summary := model.TemperatureSummary{[]model.TemperatureReading{reading}, model.CELSIUS}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)


	if err := json.NewEncoder(w).Encode(summary); err != nil {
		panic(err)
	}

}

// Gets historical temperatures from the database, filtering using the provided earliest and latest dates in the query
// parameters. Temperatures will be returned in the provided scale units - CELSIUS or FAHRENHEIT
func (c *TempController) GetHistoricalTemperatures(w http.ResponseWriter, req *http.Request) {
	earliest := time.Unix(0, 0)
	latest := time.Now()
	scale := model.CELSIUS

	//reusable error
	var err error

	//query params: earliest, latest
	q := req.URL.Query()
	e := q.Get("earliest")
	l := q.Get("latest")
	scaleText := q.Get("scale")

	if len(e) > 0 {
		earliest, err = time.Parse(NARROW_RFC3339, e) //use a narrowed RFC3339 for usability, precision not necessary
	}

	if len(l) > 0 {
		latest, err = time.Parse(NARROW_RFC3339, l)
	}

	if len(scaleText) > 0 {
		scale = model.StringToScale(scaleText)
	}

	//error checks before db lookup
	if (err != nil) {
		respondMessage(http.StatusBadRequest, w, fmt.Sprintf("Error parsing 'earliest' or 'latest' query parameters. Expected date format %v. Error: %v", NARROW_RFC3339, err))
		return
	}

	tempSummary := db.GetTempHistForDateRange(scale, earliest, latest)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tempSummary); err != nil {
		panic(err)
	}
}

// Handles updates to temperature polling. The command to stop or start polling should be provided within the JSON body
// of the incoming request.  When a request to start polling is submitted, a positive millisecond poll interval should
// also be provided.
func (c *TempController) Poll(w http.ResponseWriter, req *http.Request) {
	var poll PollRequest

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 500000))

	if err != nil {
		panic(err)
	}

	if err := req.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &poll); err != nil {
		respondMessage(http.StatusNotAcceptable, w, err.Error())
	}

	switch poll.Action {
	case "START":
		//validate interval > 0
		if poll.IntervalMilliseconds <= 0 {
			respondMessage(http.StatusBadRequest, w, fmt.Sprintf("Expected intervalMilliseconds field > 0. Got %v", poll.IntervalMilliseconds))
			return
		}

		//call poller.start
		err = c.poller.Start(poll.IntervalMilliseconds)

		if err != nil {
			respondMessage(http.StatusInternalServerError, w, fmt.Sprintf("Error starting temperature poller. Error %v", err))
			return
		}

		respondMessage(http.StatusOK, w, fmt.Sprintf("Started polling temperature with interval %v", poll.IntervalMilliseconds))
	case "STOP":
		//call poller.stop
		err = c.poller.Stop()

		if err != nil {
			respondMessage(http.StatusInternalServerError, w, fmt.Sprintf("Error stopping temperature poller. Error %v", err))
			return
		}

		respondMessage(http.StatusOK, w, "Stopped polling temperature.")
	default:
		respondMessage(http.StatusBadRequest, w, fmt.Sprintf("Invalid 'action' value. Expected 'START' or 'STOP'. Got %v", poll.Action))
	}
}

// Controller interface function, returns a list of routes handled by this controller
func (c *TempController) GetRoutes() Routes {
	routes := Routes{
		Route{Name: "GetCurrentTemperature", Method: "GET", Pattern: "/thermometer/temperature/now", HandlerFunc: c.GetCurrentTemperature},
		Route{Name: "GetTemperatureHistory", Method: "GET", Pattern: "/thermometer/temperature", HandlerFunc: c.GetHistoricalTemperatures},
		Route{Name: "Poll", Method: "PUT", Pattern: "/thermometer/poll", HandlerFunc: c.Poll},
	}

	return routes
}
