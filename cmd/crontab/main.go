package main

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/crontab"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/utility"
)

func main() {
	err := crontab.Server()
	if err != nil {
		panic(err)
	}

	err = crontab.EventJob("", "test", "*/5 * * * * *", "test",
		et.Json{
			"test": "test",
		}, func(msg event.EvenMessage) {
			worker := msg.Data
			crontab.EventStatusRunning(worker)
			console.Debug("test by event:", worker.ToString())
			crontab.EventStatusDone(worker)
		})
	if err != nil {
		panic(err)
	}

	// job, err := crontab.AddEventJob("", "test", "*/5 * * * * *", "test",
	// 	et.Json{
	// 		"test": "test",
	// 	})
	// if err != nil {
	// 	panic(err)
	// }

	// time.AfterFunc(time.Second*9, func() {
	// 	job.Stop()
	// })

	// time.AfterFunc(time.Second*12, func() {
	// 	job.Start()
	// })

	// time.AfterFunc(time.Second*15, func() {
	// 	job.Remove()
	// })

	utility.AppWait()

	console.Debug("Fin de flow")
}
