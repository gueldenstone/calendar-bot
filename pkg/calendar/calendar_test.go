package calendar

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/emersion/go-ical"
)

const (
	xHainDump_1 = "testData.txt"
)

func NewWantedCalendarEventWithUIDAndSummary(uid string, summary string) ical.Event {
	return ical.Event{
		Component: &ical.Component{
			Name: ical.CompEvent,
			Props: ical.Props{
				ical.PropSummary: []ical.Prop{
					{
						Name:   ical.PropSummary,
						Params: make(ical.Params),
						Value:  summary,
					},
				},
				ical.PropUID: []ical.Prop{
					{
						Name:   ical.PropUID,
						Params: make(ical.Params),
						Value:  uid,
					},
				},
			},
		},
	}
}

func getCalendarDataFromFile(file string) *os.File {
	f, _ := os.Open(file)
	return f
}

func TestCalendar_GetEventsOn(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name         string
		testDataFile io.ReadCloser
		args         args
		want         []ical.Event
		wantErr      bool
	}{
		{
			name:         "LED Workshop recurring",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 23, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDAndSummary("dae1e4eb-7213-4620-bc9f-1bdb8a023af9", "Workshop - Learn PCB design with KiCad"),
			},
		},
		{
			name:         "LED Workshop recurring exception",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 3, 2, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDAndSummary("66adcfc4-6827-45a2-a5b4-655923d5dd62", "How to use the latest AIs in your daily workflow - for... everything?"),
			},
		},
		{
			name:         "Gespräch unter Bäumen instances",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 3, 21, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDAndSummary("5fb7f276-54d6-4c30-a993-92cfe962e41b", "Gespräch unter Bäumen (mit Elisa Filevich)"),
			},
		},
		{
			name:         "Drones' night",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 24, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDAndSummary("1ec26b84-60e1-437d-a455-db6404dff879", "Drones' night"),
			},
		},
		{
			name:         "offener Montag",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 27, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDAndSummary("c9158eec-083a-4798-9860-99c4a83cce0f", "offener Montag"),
			},
		},
		{
			name:         "Multiple events: XMPP/Wednesday Meeting",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 8, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDAndSummary("3591c731-0e27-4902-9ae0-8748d46841f3", "XMPP-Meetup"),
				NewWantedCalendarEventWithUIDAndSummary("1695d4c9-aa37-4b52-bda9-2132ac92e3a2", "xHain Wednesday meeting"),
			},
		},
		{
			name:         "Multiple events: Kindernachmittag/Camp",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 1, 22, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDAndSummary("290b69b7-aaaf-47d4-88d1-d42366e36163", "Kindernachmittag"),
				NewWantedCalendarEventWithUIDAndSummary("77b564a3-4e8c-4e3b-b9db-990e622d04ff", "xHain @ Camp23 Brainstorming"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.testDataFile.Close()
			parser := ical.NewDecoder(tt.testDataFile)
			calendar, err := parser.Decode()
			if err != nil {
				t.Fatalf("could not open test data file: %s", tt.testDataFile)
			}

			cal := Calendar{
				url:      "file",
				tz:       time.Local,
				Logger:   log.New(os.Stdout, "[TEST] ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile),
				Calendar: calendar,
			}
			got, err := cal.GetEventsOn(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calendar.GetEventsOn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Calendar.GetEventsOn() not the right amount of events: got = %v, want %v", got, tt.want)
			}
			for i := range got {
				gotUID := got[i].Props.Get(ical.PropUID).Value
				wantUID := tt.want[i].Props.Get(ical.PropUID).Value
				if gotUID != wantUID {
					t.Errorf("Calendar.GetEventsOn() got = %v, want %v", gotUID, wantUID)
				}
			}
		})
	}
}
