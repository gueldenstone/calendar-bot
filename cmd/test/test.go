package main

import (
	"fmt"
	"time"
)

func main() {
	loc := time.Local
	fmt.Println(loc.String())
	// cal, err := calendar.NewCalendar("https://files.x-hain.de/remote.php/dav/public-calendars/Yi63cicwgDnjaBHR/?export")
	// if err != nil {
	// 	panic(err)
	// }
	// evts, err := cal.GetEventsOn(time.Now().Add(24 * time.Hour))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(evts)
	// msg, err := message.NewTemplatedMessage("./templates/event.html", "./templates/event.txt", evts)
	// if err != nil {
	// 	panic(err)
	// }
	// html, txt, err := msg.Render()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(html)
	// fmt.Println(txt)
}
