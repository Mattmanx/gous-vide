package routine
import (
	"time"
	"github.com/mattmanx/gous-vide/model"
	"log"
	"fmt"
	"github.com/mattmanx/gous-vide/hardware"
)

//Not thread safe.. synchronize upstream
type Program struct {
	running bool
	ticker *time.Ticker
	heater *hardware.Heater
	recipeRun model.RecipeRun
}

func NewProgram(recipe model.Recipe, heater *hardware.Heater) *Program {
	return &Program{running: false, heater: heater, recipeRun: model.RecipeRun{Recipe: recipe}}
}

func (p *Program) IsRunning() bool {
	return p.running
}

func (p *Program) Start(pollIntervalMs int) error {
	if p.running {
		log.Print("Request to start running program, but already in started state.");
		return fmt.Errorf("Request to start program, but already in started state. Unable to process request.")
	}

	p.ticker = time.NewTicker(time.Duration(pollIntervalMs) * time.Millisecond)
	p.recipeRun.StartTime = time.Now()

	go func() {
		for t := range p.ticker.C {
			temp, err := hardware.CurrentTempCelsius()

			if err == nil {
				log.Printf("Program tick at %v. Operating on current temperature %vC, target temp %vC", t, temp, p.recipeRun.TargetTempCelcius)

				//TODO: should lock here... or make sure TargetTempCelcius is immutable
				if temp < p.recipeRun.TargetTempCelcius {
					err = p.heater.TurnOn()
				} else {
					err = p.heater.TurnOff()
				}

				if err != nil {
					log.Printf("ERROR: Program tick at %v, but cannot control heater.", t)
				}
			}

			if err != nil {
				log.Printf("ERROR: Poll tick at %v, but unable to get or save temp. Error: %v", t, err)
			}
		}
	}()

	p.running = true

	return nil
}

func (p *Program) Stop() error {
	if !p.running {
		log.Print("Request to stop program, but already in stopped state.");
		return fmt.Errorf("Request to stop program, but already in stopped state. Unable to process request.")
	}

	p.ticker.Stop()

	p.running = false

	e := p.heater.TurnOff()

	return e
}

func (p *Program) GetStatus() model.RecipeRun {
	return p.recipeRun
}