package main

import (
	"flag"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gueldenstone/calendar-bot/internal/config"
	"github.com/gueldenstone/calendar-bot/internal/matrix"
	"github.com/gueldenstone/calendar-bot/internal/message"
	"github.com/gueldenstone/calendar-bot/pkg/calendar"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const flags = log.Ldate | log.Ltime | log.Lmsgprefix

var (
	errLog  = log.New(os.Stderr, "[ERROR] ", flags|log.Lshortfile)
	infoLog = log.New(os.Stdout, "[INFO] ", flags)
)

var (
	configFile   = flag.String("config", "", "Path to config file")
	htmlTmplPath = flag.String("html", "", "Path to html template file")
	txtTmplPath  = flag.String("txt", "", "Path to txt template file")
)

func main() {
	flag.Parse()
	// load config file
	if *configFile == "" || *htmlTmplPath == "" || *txtTmplPath == "" {
		errLog.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	conf, err := config.Parse(*configFile)
	if err != nil {
		errLog.Fatalf("could not read config file: %s", *configFile)
	}
	// validate roomID format for rooms publish in
	validRomms := make([]string, 0)
	for _, id := range conf.Rooms {
		if !matrix.IsValidRoomId(id) {
			errLog.Printf("RoomID '%s' is not valid, ignoring this room\n", id)
		} else {
			validRomms = append(validRomms, id)
		}
	}
	if len(validRomms) == 0 {
		errLog.Fatalf("No valid roomIDs have been provided, exiting...\n")
	}
	conf.Rooms = validRomms

	infoLog.Println("Logging into", conf.Homeserver, "as", conf.Username)
	client, err := matrix.LoginToHomeServer(conf.Homeserver, conf.Username, conf.Password)
	if err != nil {
		errLog.Fatal(err)
	}
	defer func() {
		if _, err := client.Logout(); err != nil {
			errLog.Println(err)
		}
	}()
	infoLog.Println("Login successful")

	// validate roomIDs
	rooms := make([]id.RoomID, 0)
	for _, rid := range conf.Rooms {
		var roomID id.RoomID
		if strings.HasPrefix(rid, "#") {
			resp, err := client.ResolveAlias(id.RoomAlias(rid))
			if err != nil {
				errLog.Printf("Error: Could not find the room: %s\n", err)
				continue
			}
			roomID = resp.RoomID
		} else {
			roomID = id.RoomID(rid)
		}
		if _, err := client.JoinRoomByID(roomID); err != nil {
			errLog.Printf("could not join room: %s\n", err)
		} else {
			rooms = append(rooms, roomID)
		}

	}
	if len(rooms) == 0 {
		errLog.Fatalf("Could not resolve or find any of the provided rooms!\n")
	}
	var notifyTime time.Time
	if conf.Testing {
		infoLog.Println("Running in test configuration")
		infoLog.Println("Using 5 seconds from now as notify time")
		notifyTime = time.Now().Add(5 * time.Second)
	} else {
		notifyTime, err = time.Parse("15:04", conf.NotifyTime)
		if err != nil {
			// default time
			errLog.Println("could not parse notification time, using 10:00")
			notifyTime, _ = time.Parse("15:04", "10:00")
		}
	}
	timezone := time.Local
	s := gocron.NewScheduler(timezone)
	infoLog.Printf("Scheduling notifications for %s", notifyTime.Format("15:04"))
	_, err = s.Every(1).Day().At(notifyTime).Do(func() {
		infoLog.Println("Start Notification")
		cal, err := calendar.NewNextcloudCalendar(conf.Calendar)
		if err != nil {
			errLog.Println(err)
			return
		}
		todayEvents, err := cal.GetEventsOn(time.Now())
		if err != nil {
			errLog.Println(err)
			return
		}

		// filter events which do not start today
		todayEvents = slices.Collect(func(yield func(calendar.Event) bool) {
			for _, event := range todayEvents {
				if event.Start.Local().Day() == time.Now().Local().Day() {
					if !yield(event) {
						return
					}
				}
			}
		})

		if len(todayEvents) == 0 {
			infoLog.Println("No events today!")
			return
		}
		infoLog.Printf("Sending notification with %d events\n", len(todayEvents))

		tmplMsg, err := message.NewTemplatedMessage(*htmlTmplPath, *txtTmplPath, todayEvents, timezone)
		if err != nil {
			errLog.Println(err)
			return
		}
		matrixMsg, err := tmplMsg.MatrixMessage()
		if err != nil {
			errLog.Println(err)
			return
		}
		for _, room := range rooms {
			if _, err := client.SendMessageEvent(room, event.EventMessage, matrixMsg); err != nil {
				errLog.Println(err)
				continue
			}
		}
	})
	if err != nil {
		errLog.Println(err)
		return
	}
	s.StartBlocking()
}
