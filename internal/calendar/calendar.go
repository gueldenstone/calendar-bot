package calendar

import (
	"net/http"
	"sort"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

type Calendar struct {
	url string
	tz  *time.Location
	*ical.Calendar
}

func NewCalendar(url string) (Calendar, error) {
	calendar := Calendar{url: url, tz: time.Local}
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

func (cal *Calendar) SetTimezone(tz *time.Location) {
	cal.tz = tz
}

func (cal Calendar) GetEventsOn(date time.Time) ([]ical.Event, error) {
	events := make([]ical.Event, 0)
	todayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, cal.tz)
	todayEnd := todayStart.Add(24 * time.Hour)
	for _, e := range cal.Events() {
		start, err := e.DateTimeStart(cal.tz)
		if err != nil {
			return []ical.Event{}, err
		}
		end, err := e.DateTimeEnd(cal.tz)
		if err != nil {
			return []ical.Event{}, err
		}
		// regular event
		if (start.After(todayStart) || start == todayStart) && start.Before(todayEnd) || (start.Before(todayStart) && end.After(todayEnd)) {
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
	// sort events
	sort.SliceStable(events, func(i, j int) bool {
		start1, _ := events[i].DateTimeStart(cal.tz)
		start2, _ := events[j].DateTimeStart(cal.tz)
		return start1.Before(start2)
	})
	return events, nil
}
