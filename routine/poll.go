package routine
import (
	"sync"
	"time"
	"fmt"
	"log"
	"github.com/mattmanx/gous-vide/hardware"
	"github.com/mattmanx/gous-vide/db"
	"github.com/mattmanx/gous-vide/model"
)

// TempPoller is a threadsafe, reusable structure that can be used to poll temperature from the thermometer at
// specified interval and save each temperature to the datasource for future retrieval. Retrieval and save is performed
// in a separate goroutine using a time.Ticker.  Calls to start an already started poller, or stop an already stopped
// poller, will result in an error and no change to the current state of the poller.
type TempPoller struct {
	polling bool
	ticker *time.Ticker
	sync.RWMutex	//lock to make sure we only ever have a single ticker running at a time
}

// Creates a new TempPoller in a stopped state.
func NewTempPoller() *TempPoller {
	return new(TempPoller)
}

// Whether this poller is currently in a started state and is polling for temperatures at interval, or not.
func (tp *TempPoller) IsPolling() bool {
	tp.RLock()
	defer tp.RUnlock()

	return tp.polling
}

// Starts polling at the specified interval. If the poller is already started, or otherwise an error occurs while
// starting the poller, an error will be returned and the state of this poller will not be changed.
func (tp *TempPoller) Start(intervalMs int) error {
	tp.Lock()
	defer tp.Unlock()

	if tp.polling {
		log.Print("Request to start poller, but already in started state.");
		return fmt.Errorf("Request to start poller, but already in started state. Unable to process request.")
	}

	tp.ticker = time.NewTicker(time.Duration(intervalMs) * time.Millisecond)

	go func() {
		for t := range tp.ticker.C {
			temp, err := hardware.CurrentTempCelsius()

			if err == nil {
				log.Printf("Poll tick at %v. Recording current temperature %vC", t, temp)

				err = db.SaveTemperature(temp, model.CELSIUS)
			}

			if err != nil {
				log.Printf("ERROR: Poll tick at %v, but unable to get or save temp. Error: %v", t, err)
			}
		}
	}()

	tp.polling = true

	return nil
}

// Stops polling, assuming the poller is currently in a started state.  If the poller is already stopped, an error will
// be returned and the state of this poller will not be changed.
func (tp *TempPoller) Stop() error {
	tp.Lock()
	defer tp.Unlock()

	if !tp.polling {
		log.Print("Request to stop poller, but already in stopped state.");
		return fmt.Errorf("Request to stop poller, but already in stopped state. Unable to process request.")
	}

	tp.ticker.Stop()

	tp.polling = false

	return nil
}

