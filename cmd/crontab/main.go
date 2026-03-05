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

	err = crontab.AddEventJob("test", "0 50 8 * * *", 0, true,
		et.Json{
			"test": "test",
		},
		func(msg event.EvenMessage) {
			// worker := msg.Data
			console.Debug("Hol run test by event:", msg.ToString())
		})
	if err != nil {
		panic(err)
	}

	err = crontab.AddScheduleJob("test2", "2026-03-05T08:51:00", true,
		et.Json{
			"test": "test2",
		},
		func(msg event.EvenMessage) {
			// worker := msg.Data
			console.Debug("Hol run test2 by event:", msg.ToString())
		})
	if err != nil {
		panic(err)
	}

	utility.AppWait()
}

func test() {
	loc, _ := time.LoadLocation("America/Bogota")
	fmt.Println("Sistema:", time.Now())
	fmt.Println("Bogotá :", time.Now().In(loc))

	i := 0
	var shot *time.Timer
	shot = time.AfterFunc(time.Second*3, func() {
		i++
		console.Debug("Running job:", i)
		if i < 5 {
			shot.Reset(time.Second * 3)
		}
	})

	console.Debug("Before stop:", shot)
	shot.Stop()
	console.Debug("After stop")
}
