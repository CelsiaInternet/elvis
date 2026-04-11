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

	// err = crontab.AddEventJob("test", "0 30 11 * * *", 0, true,
	// 	et.Json{
	// 		"test": "test",
	// 	},
	// 	func(msg event.EvenMessage) {
	// 		// worker := msg.Data
	// 		console.Debug("Hol run test by event:", msg.ToString())
	// 	})
	// if err != nil {
	// 	panic(err)
	// }

	for i := 0; i < 5; i++ {
		timeStr := time.Now().Add(time.Second * 5).Format("2006-01-02T15:04:05")
		err = crontab.AddScheduleJob(fmt.Sprintf("test_%d", i), timeStr, true,
			et.Json{
				"test": fmt.Sprintf("test_%d", i),
			},
			func(msg event.EvenMessage) {
				// worker := msg.Data
				console.Debug("Hola run test by event:", msg.ToString())
			})
		if err != nil {
			panic(err)
		}
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
