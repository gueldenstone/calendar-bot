package calendar

import (
	"net/http"
	"time"

	"github.com/apognu/gocal"
)

type Calendar struct {
	url string
	tz  *time.Location
}

func NewCalendar(url string) Calendar {
	calendar := Calendar{url: url, tz: time.Local}
	return calendar
}

func (cal *Calendar) SetTimezone(tz *time.Location) {
	cal.tz = tz
}

func (cal Calendar) GetEventsOn(date time.Time) ([]gocal.Event, error) {
	events := make([]gocal.Event, 0)
	resp, err := http.Get(cal.url)
	if err != nil {
		return events, err
	}
	defer resp.Body.Close()

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, cal.tz)
	end := start.Add(24 * time.Hour)
	c := gocal.NewParser(resp.Body)
	c.Start, c.End = &start, &end
	c.Parse()
	return c.Events, nil
}
