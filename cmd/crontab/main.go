package main

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/utility"
)

func main() {
	// err := crontab.Load("test")
	// if err != nil {
	// 	panic(err)
	// }

	// err = crontab.AddEventJob("test", "*/5 * * * * *", "test", 0, true,
	// 	et.Json{
	// 		"test": "test",
	// 	},
	// 	func(msg event.EvenMessage) {
	// 		worker := msg.Data
	// 		console.Debug("test by event:", worker.ToString())
	// 	})
	// if err != nil {
	// 	panic(err)
	// }

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

	utility.AppWait()
}
