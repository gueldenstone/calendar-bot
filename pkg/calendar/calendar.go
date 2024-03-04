package calendar

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	jcal "github.com/xHain-hackspace/go-jcal"
)

// Assembles a list of events on a given date
func GetEventsOn(date time.Time) ([]jcal.Event, error) {
	events := make([]jcal.Event, 0)

	// Get the first and last date of the day
	firstDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	lastDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)
	url := fmt.Sprintf("https://files.x-hain.de/remote.php/dav/public-calendars/Yi63cicwgDnjaBHR/?export&accept=jcal&start=%d&end=%d&expand=1", firstDate.Unix(), lastDate.Unix())
	resp, err := http.Get(url)
	if err != nil {
		return events, err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return events, err
	}
	resp.Body.Close()

	var jcalObj jcal.JCalObject
	if err := json.Unmarshal(respData, &jcalObj); err != nil {
		log.Fatalf("Error parsing jCal JSON: %v", err)
	}
	events = jcalObj.Events
	sortEvents(events)
	return events, nil
}

// Sorts events by start time
func sortEvents(events []jcal.Event) {
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].DtStart.Before(events[j].DtStart)
	})
}
