package main

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/crontab"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/utility"
)

func main() {
	err := crontab.Load("test", nil)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		timeStr := time.Now().Add(time.Second * 5).Format("2006-01-02T15:04:05")
		err = crontab.AddScheduleJob(fmt.Sprintf("test_%d", i), timeStr, true,
			et.Json{
				"test": fmt.Sprintf("test_%d", i),
			},
			func(msg event.EvenMessage) {
				console.Debug("Hola run test by event:", msg.ToString())
			})
		if err != nil {
			panic(err)
		}
	}

	utility.AppWait()
}
