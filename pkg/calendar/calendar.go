package calendar

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/emersion/go-ical"
	// "github.com/teambition/rrule-go"
)

type Calendar struct {
	url string
	tz  *time.Location
	*ical.Calendar
	*log.Logger
}

func NewCalendar(url string, l *log.Logger) (Calendar, error) {
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
	calendar.Logger = l
	return calendar, nil
}

func (cal Calendar) GetEventsOn(date time.Time) ([]ical.Event, error) {
	events := make([]ical.Event, 0)
	todayStart := GetDateWithoutTime(date)
	todayEnd := todayStart.Add(24 * time.Hour)
	for _, event := range cal.Events() {
		start, err := event.DateTimeStart(cal.tz)
		if err != nil {
			return []ical.Event{}, err
		}
		end, err := event.DateTimeEnd(cal.tz)
		if err != nil {
			return []ical.Event{}, err
		}
		// regular event
		if (start.After(todayStart) || start.Local() == todayStart.Local()) && start.Before(todayEnd) || (start.Before(todayStart) && end.After(todayEnd)) {
			events = append(events, event)
			continue
		}
		// recurring event
		reccurenceSet, err := event.RecurrenceSet(cal.tz)
		if err != nil {
			cal.Printf("could not get recurrence set: %s\n", err)
			continue
		}
		if reccurenceSet == nil {
			// no recurrence
			continue
		}
		if GetDateWithoutTime(reccurenceSet.After(todayStart, true)).Local() == GetDateWithoutTime(date).Local() {
			events = append(events, event)
		}
	}
	// sort events
	sort.SliceStable(events, func(i, j int) bool {
		start1, _ := events[i].DateTimeStart(cal.tz)
		start2, _ := events[j].DateTimeStart(cal.tz)
		return start1.Before(start2)
	})

	// check for doubles via uid
	uids := make(map[string]struct{})
	dedupedEvents := make([]ical.Event, 0)
	for _, e := range events {
		uid := e.Props.Get(ical.PropUID).Value
		if _, ok := uids[uid]; !ok {
			dedupedEvents = append(dedupedEvents, e)
			uids[uid] = struct{}{}
		}
	}
	return dedupedEvents, nil
}

func GetDateWithoutTime(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}
