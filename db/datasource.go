package db

import (
	"github.com/boltdb/bolt"
	"time"
	"log"
	"sync"
	m "github.com/mattmanx/gous-vide/model"
	"fmt"
	"encoding/binary"
	"bytes"
)

const (
	DEFAULT_DB = "gousvide.db"
	TEMP_HIST_BUCKET = "TempHist"
)

var (
	db *bolt.DB
	once sync.Once
)

// Opens the datasource encapsulated by this package, returning an error if the database cannot be opened
func Open() error {
	var err error

	// Only open once... all db operations will be routed through this package.
	once.Do(func() {
		// There should be no contention for access to the database, so fail fast instead of blocking until available
		db, err = bolt.Open(DEFAULT_DB, 0600, &bolt.Options{Timeout: 1 * time.Second})
	})

	if (err != nil) {
		log.Printf("Error attempting to open bolt database: %v. Error message: %v", db, err)
		return err
	}

	// Create all of our buckets
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TEMP_HIST_BUCKET))
		if err != nil {
			return fmt.Errorf("can't create bucket TempHist: %s", err)
		}

		return nil
	})

	return nil
}

// Closes the datasource encapsulated by this package, returning an error if the database cannot be opened
func Close() error {
	db.Close()
	return nil
}

// Keeping all DAO-esque functions within this source file. Will move into separate modules (same package) if they grow
// too large

// Saves a temperature record to the database, returning an error if the value was unable to be saved. All temperatures
// are saved in celcius, and can be converted when retrieved if necessary.
func SaveTemperature(temperature float64, scale m.TemperatureScale) error {
	// Get current timestamp as our key on the temperature record
	t := time.Now()

	// use a sortable time format, as suggested by bolt, so we can do range queries
	key := []byte(t.Format(time.RFC3339))

	//convert our temperature to an integer w/ 2 digit scale for easier storage / retrieval from db
	celsius := m.ConvertTemperature(temperature, scale, m.CELSIUS)

	saveTemp := int(celsius * 100)

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TEMP_HIST_BUCKET))

		err := b.Put(key, itob(saveTemp))
		return err
	})
}

// Gets all historical temperatures in units per the provided scale
func GetTempHist(scale m.TemperatureScale) m.TemperatureSummary {
	var temps []m.TemperatureReading;

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(TEMP_HIST_BUCKET))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			temps = append(temps, *createTempFromRaw(k, v, scale))
		}

		return nil
	})

	return m.TemperatureSummary{temps, scale}
}

// Gets historical temperatures within the defined date range in units per the provided scale
func GetTempHistForDateRange(scale m.TemperatureScale, earliest time.Time, latest time.Time) m.TemperatureSummary {
	var temps []m.TemperatureReading;

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(TEMP_HIST_BUCKET))

		min := []byte(earliest.Format(time.RFC3339))
		max := []byte(latest.Format(time.RFC3339))

		c := b.Cursor()

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			temps = append(temps, *createTempFromRaw(k, v, scale))
		}

		return nil
	})

	return m.TemperatureSummary{temps, scale}
}

// Converts the binary representations of the timestamp and temperature value, constructs a new TemperatureReading to
// encapsulate these values, and returns a pointer to that temperature reading.
func createTempFromRaw(timestamp []byte, value []byte, scale m.TemperatureScale) *m.TemperatureReading {
	t, err := time.Parse(time.RFC3339, string(timestamp))

	//The timestamp should NEVER be saved in the wrong format. If so, db should be reset.
	if(err != nil) {
		panic(err)
	}

	temperature := float64(btoi(value)) / 100

	temperature = m.ConvertTemperature(temperature, m.CELSIUS, scale)

	return &m.TemperatureReading{t, temperature}
}

// itob returns an 8-byte varint representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.PutVarint(b, int64(v))
	return b
}

// btoi returns an int representation of an assumed binary varint
func btoi(v []byte) int {
	i, _ := binary.Varint(v)

	return int(i)
}