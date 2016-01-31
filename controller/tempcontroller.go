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
)

const (
	NARROW_RFC3339 = "2006-01-02T15:04"
)

type TempController struct {
	//nothing for now
}

func NewTempController() *TempController {
	return &TempController{}
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
		respondError(w, fmt.Sprintf("Error parsing query param 'save'. Expected 'true' or 'false'. Error: %v", err))
		return
	}

	currentTemp, err = hardware.CurrentTempCelsius()

	if err != nil {
		respondError(w, fmt.Sprintf("Error retrieving current temperature. Error: %v", err))
		return
	}

	if save {
		err = db.SaveTemperature(currentTemp, model.CELSIUS)
	}

	if err != nil {
		respondError(w, fmt.Sprintf("Error saving current temperature. Error: %v", err))
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
		respondError(w, fmt.Sprintf("Error parsing 'earliest' or 'latest' query parameters. Expected date format %v. Error: %v", NARROW_RFC3339, err))
		return
	}

	tempSummary := db.GetTempHistForDateRange(scale, earliest, latest)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tempSummary); err != nil {
		panic(err)
	}
}


// Controller interface function, returns a list of routes handled by this controller
func (c *TempController) GetRoutes() Routes {
	routes := Routes{
		Route{Name: "GetCurrentTemperature", Method: "GET", Pattern: "/thermometer/temperature/now", HandlerFunc: c.GetCurrentTemperature},
		Route{Name: "GetTemperatureHistory", Method: "GET", Pattern: "/thermometer/temperature", HandlerFunc: c.GetHistoricalTemperatures},
	}

	return routes
}
