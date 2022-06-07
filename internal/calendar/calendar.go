package calendar

import (
	"net/http"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

type Calendar struct {
	url string
	*ical.Calendar
}

func NewCalendar(url string) (Calendar, error) {
	calendar := Calendar{url: url}
	resp, err := http.Get(url)
	if err != nil {
		return calendar, err
	}
	defer resp.Body.Close()

	parser := ical.NewDecoder(resp.Body)

	cal, err := parser.Decode()
	if err != nil {
		return calendar, err
	}
	calendar.Calendar = cal
	return calendar, nil
}

func (cal Calendar) GetEventsOn(date time.Time) ([]ical.Event, error) {
	events := make([]ical.Event, 0)
	todayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	todayEnd := todayStart.Add(24 * time.Hour)
	for _, e := range cal.Events() {
		start, err := e.DateTimeStart(time.Local)
		if err != nil {
			return []ical.Event{}, err
		}
		// regular event
		if start.After(todayStart) && start.Before(todayEnd) {
			events = append(events, e)
			continue
		}
		// recurring event?
		roption, err := e.Props.RecurrenceRule()
		if err != nil {
			return []ical.Event{}, err
		}
		if roption != nil {
			roption.Dtstart = start
			rule, err := rrule.NewRRule(*roption)
			if err != nil {
				return []ical.Event{}, err
			}
			times := rule.Between(todayStart, todayEnd, true)
			for _, t := range times {
				copyEvent := e
				copyEvent.Props.SetDateTime(ical.PropDateTimeStart, t)
				events = append(events, copyEvent)
			}
		}
	}
	return events, nil
}
