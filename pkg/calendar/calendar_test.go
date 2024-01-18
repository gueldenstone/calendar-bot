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

func NewWantedCalendarEvent(uid string, summary string, start time.Time, end time.Time) EventData {
	return EventData{
		UID:     uid,
		Summary: summary,
		Start:   start,
		End:     end,
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
		want         []EventData
		wantErr      bool
	}{
		{
			name:         "LED Workshop recurring",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 23, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []EventData{
				NewWantedCalendarEvent(
					"dae1e4eb-7213-4620-bc9f-1bdb8a023af9",
					"Workshop - Learn PCB design with KiCad",
					time.Date(2023, 2, 23, 18, 30, 0, 0, time.Local),
					time.Date(2023, 2, 23, 20, 30, 0, 0, time.Local),
				),
			},
		},
		{
			name:         "LED Workshop recurring exception",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 3, 2, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []EventData{
				NewWantedCalendarEvent(
					"66adcfc4-6827-45a2-a5b4-655923d5dd62",
					"How to use the latest AIs in your daily workflow - for... everything?",
					time.Date(2023, 3, 2, 18, 30, 0, 0, time.Local),
					time.Date(2023, 3, 2, 20, 30, 0, 0, time.Local),
				),
			},
		},
		{
			name:         "Gespräch unter Bäumen instances",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 3, 21, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []EventData{
				NewWantedCalendarEvent(
					"5fb7f276-54d6-4c30-a993-92cfe962e41b",
					"Gespräch unter Bäumen (mit Elisa Filevich)",
					time.Date(2023, 3, 21, 19, 0, 0, 0, time.Local),
					time.Date(2023, 3, 21, 20, 0, 0, 0, time.Local),
				),
			},
		},
		{
			name:         "Drones' night",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 24, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []EventData{
				NewWantedCalendarEvent(
					"1ec26b84-60e1-437d-a455-db6404dff879",
					"Drones' night",
					time.Date(2023, 2, 24, 18, 0, 0, 0, time.Local),
					time.Date(2023, 2, 24, 21, 0, 0, 0, time.Local),
				),
			},
		},
		{
			name:         "offener Montag",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 27, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []EventData{
				NewWantedCalendarEvent(
					"c9158eec-083a-4798-9860-99c4a83cce0f",
					"offener Montag",
					time.Date(2023, 2, 27, 18, 0, 0, 0, time.Local),
					time.Date(2023, 2, 27, 23, 59, 0, 0, time.Local),
				),
			},
		},
		{
			name:         "Multiple events: XMPP/Wednesday Meeting",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 2, 8, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []EventData{
				NewWantedCalendarEvent(
					"3591c731-0e27-4902-9ae0-8748d46841f3",
					"XMPP-Meetup",
					time.Date(2023, 2, 8, 18, 0, 0, 0, time.Local),
					time.Date(2023, 2, 8, 21, 0, 0, 0, time.Local),
				),
				NewWantedCalendarEvent(
					"1695d4c9-aa37-4b52-bda9-2132ac92e3a2",
					"xHain Wednesday meeting",
					time.Date(2023, 2, 8, 20, 30, 0, 0, time.Local),
					time.Date(2023, 2, 8, 22, 0, 0, 0, time.Local),
				),
			},
		},
		{
			name:         "Multiple events: Kindernachmittag/Camp",
			testDataFile: getCalendarDataFromFile(xHainDump_1),
			args: args{
				date: time.Date(2023, 1, 22, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
			want: []EventData{
				NewWantedCalendarEvent(
					"290b69b7-aaaf-47d4-88d1-d42366e36163",
					"Kindernachmittag",
					time.Date(2023, 1, 22, 13, 0, 0, 0, time.Local),
					time.Date(2023, 1, 22, 18, 0, 0, 0, time.Local),
				),
				NewWantedCalendarEvent(
					"77b564a3-4e8c-4e3b-b9db-990e622d04ff",
					"xHain @ Camp23 Brainstorming",
					time.Date(2023, 1, 22, 20, 0, 0, 0, time.Local),
					time.Date(2023, 1, 22, 21, 0, 0, 0, time.Local),
				),
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
				if got[i].UID != tt.want[i].UID {
					t.Errorf("Calendar.GetEventsOn() got UID = %v, want UID %v", got[i].UID, tt.want[i].UID)
				}
				if !got[i].Start.Equal(tt.want[i].Start) {
					t.Errorf("Calendar.GetEventsOn() got Start = %v, want Start %v", got[i].Start, tt.want[i].Start)
				}
				if !got[i].End.Equal(tt.want[i].End) {
					t.Errorf("Calendar.GetEventsOn() got End = %v, want End %v", got[i].End, tt.want[i].End)
				}
			}
		})
	}
}
