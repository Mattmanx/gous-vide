package server
import (
	"github.com/mattmanx/gous-vide/controller"
	"github.com/gorilla/mux"
	"github.com/mattmanx/gous-vide/middleware"
	"strconv"
	"net/http"
	"log"
	"github.com/mattmanx/gous-vide/hardware"
)

func Start(port int) {
	router := newRouter()

	// Create dependencies for controllers
	heater := hardware.NewHeater()

	// Create controllers
	// -- tempController := controller.NewTempController
	heaterController := controller.NewHeaterController(heater)

	controllers := []controller.Controller{heaterController}

	// Register with router
	register(controllers, router)

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), router))
}

// Helper method to create a new router and set basic configurations
func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	return router
}

// Helper method to register a set of controllers with our router
func register(controllers []controller.Controller, router *mux.Router) {
	// Completely unnecessary, just playing around with Go's first-class func
	addRoute := func(route controller.Route) {
		//wrap each route in a logger
		routeHandler := middleware.LoggerFilter(route.HandlerFunc, route.Name)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(routeHandler)
	}

	// a controllers.flatmap(GetRoutes).apply(addRoute) would be nice here ;)
	for _, controller := range controllers {
		for _, route := range controller.GetRoutes() {
			addRoute(route)
		}
	}
}
