package calendar

import (
	"fmt"
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

type EventData struct {
	UID         string
	Start       time.Time
	End         time.Time
	Summary     string
	Location    string
	Description string
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

func (cal Calendar) GetEventsOn(date time.Time) ([]EventData, error) {
	events := make([]ical.Event, 0)
	todayStart := GetDateWithoutTime(date)
	todayEnd := todayStart.Add(24 * time.Hour)
	for _, event := range cal.Events() {
		start, err := event.DateTimeStart(cal.tz)
		if err != nil {
			return []EventData{}, err
		}
		end, err := event.DateTimeEnd(cal.tz)
		if err != nil {
			return []EventData{}, err
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

	// Convert ical.Events to EventData
	eventDatas := make([]EventData, 0)
	for _, event := range dedupedEvents {
		eventData, err := cal.ConvertToEventData(event, date)
		if err != nil {
			return nil, err
		}
		eventDatas = append(eventDatas, eventData)
	}

	return eventDatas, nil
}

func GetDateWithoutTime(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}
func (cal Calendar) ConvertToEventData(icalEvent ical.Event, d time.Time) (EventData, error) {
	eventData := EventData{}

	// Handle UID
	uidProp := icalEvent.Props.Get(ical.PropUID)
	if uidProp != nil {
		eventData.UID = uidProp.Value
	} else {
		return eventData, fmt.Errorf("UID is missing for event %s", icalEvent.Name)
	}

	// Handle DTSTART
	startProp := icalEvent.Props.Get(ical.PropDateTimeStart)
	if startProp == nil {
		return eventData, fmt.Errorf("DTSTART is missing for event %s", icalEvent.Name)
	}
	eventStart, err := startProp.DateTime(cal.tz)
	if err != nil {
		return eventData, err
	}
	eventData.Start = time.Date(d.Year(), d.Month(), d.Day(), eventStart.Hour(), eventStart.Minute(), 0, 0, d.Location())

	// Handle DTEND
	endProp := icalEvent.Props.Get(ical.PropDateTimeEnd)
	if endProp != nil {
		eventEnd, err := endProp.DateTime(cal.tz)
		if err != nil {
			return eventData, err
		}
		// Calculate the difference in days and adjust eventData.End
		daysDiff := int(eventEnd.Sub(eventStart).Hours() / 24)
		eventData.End = time.Date(d.Year(), d.Month(), d.Day()+daysDiff, eventEnd.Hour(), eventEnd.Minute(), 0, 0, d.Location())
	}

	// Handle SUMMARY
	summaryProp := icalEvent.Props.Get(ical.PropSummary)
	if summaryProp != nil {
		eventData.Summary = summaryProp.Value
	}

	// Handle LOCATION
	locationProp := icalEvent.Props.Get(ical.PropLocation)
	if locationProp != nil {
		eventData.Location = locationProp.Value
	}

	// Handle DESCRIPTION
	descriptionProp := icalEvent.Props.Get(ical.PropDescription)
	if descriptionProp != nil {
		eventData.Description = descriptionProp.Value
	} else {
		eventData.Description = ""
	}

	return eventData, nil
}
