package calendar

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	jcal "github.com/xHain-hackspace/go-jcal"
)

type Event struct {
	Start       time.Time
	End         time.Time
	Summary     string
	Description string
}

type NextcloudCalendar struct {
	BaseUrl url.URL
}

// NewNextcloudCalendar creates a new nextcloud calendar instance
//
// Parameters:
//
//	baseUrlStr - the base url e.g. https://files.x-hain.de/remote.php/dav/public-calendars/Yi63cicwgDnjaBHR
func NewNextcloudCalendar(baseUrlStr string) (NextcloudCalendar, error) {
	calendar := NextcloudCalendar{}

	baseUrlStr = strings.TrimSpace(baseUrlStr)

	baseUrl, err := url.Parse(baseUrlStr)
	if err != nil {
		return calendar, err
	}
	calendar.BaseUrl = *baseUrl

	return calendar, nil
}

// Assembles a list of events on a given date
func (c NextcloudCalendar) GetEventsOn(date time.Time) ([]Event, error) {
	events := make([]Event, 0)

	// build url params
	params := url.Values{}
	params.Add("accept", "jcal")
	params.Add("expand", "1")
	params.Add("export", "")

	// Get the first and last date of the day
	dayBegin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)
	params.Add("start", fmt.Sprintf("%d", dayBegin.Unix()))
	params.Add("end", fmt.Sprintf("%d", dayEnd.Unix()))

	// add params to url
	requestUrl := c.BaseUrl
	requestUrl.RawQuery = params.Encode()

	// make request
	resp, err := http.Get(requestUrl.String())
	if err != nil {
		return events, err
	}

	// read data
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return events, err
	}
	resp.Body.Close()

	// parse jcal
	var jcalObj jcal.JCalObject
	if err := json.Unmarshal(respData, &jcalObj); err != nil {
		log.Fatalf("Error parsing jCal JSON: %v", err)
	}

	// convert from jcal
	for _, jEvent := range jcalObj.Events {
		events = append(events, fromJcal(jEvent))
	}

	// sort events by start time
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})
	return events, nil
}

func fromJcal(jEvent jcal.Event) Event {
	return Event{
		Summary:     jEvent.Summary,
		Description: jEvent.Description,
		Start:       jEvent.DtStart,
		End:         jEvent.DtEnd,
	}
}
