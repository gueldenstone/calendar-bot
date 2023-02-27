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

func NewWantedCalendarEventWithUIDSet(uid string) ical.Event {
	return ical.Event{
		Component: &ical.Component{
			Name: ical.CompEvent,
			Props: ical.Props{
				ical.PropUID: []ical.Prop{
					{
						Name:   ical.PropUID,
						Params: make(ical.Params, 0),
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
				date: time.Date(2023, 02, 23, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDSet("dae1e4eb-7213-4620-bc9f-1bdb8a023af9"),
			},
		},
		{
			name:         "Drones' night",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 02, 24, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDSet("1ec26b84-60e1-437d-a455-db6404dff879"),
			},
		},
		{
			name:         "offener Montag",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 02, 27, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDSet("c9158eec-083a-4798-9860-99c4a83cce0f"),
			},
		},
		{
			name:         "Multiple events: XMPP/Wednesday Meeting",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 02, 8, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []ical.Event{
				NewWantedCalendarEventWithUIDSet("3591c731-0e27-4902-9ae0-8748d46841f3"),
				NewWantedCalendarEventWithUIDSet("1695d4c9-aa37-4b52-bda9-2132ac92e3a2"),
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
				NewWantedCalendarEventWithUIDSet("290b69b7-aaaf-47d4-88d1-d42366e36163"),
				NewWantedCalendarEventWithUIDSet("77b564a3-4e8c-4e3b-b9db-990e622d04ff"),
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
				Logger:   log.New(os.Stdout, "[TEST]", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile),
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
