package main

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/crontab"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/utility"
)

func main() {
	err := crontab.Load("test")
	if err != nil {
		panic(err)
	}

	err = crontab.AddEventJob("test", "*/5 * * * * *", "test", 0, true,
		et.Json{
			"test": "test",
		},
		func(msg event.EvenMessage) {
			worker := msg.Data
			console.Debug("test by event:", worker.ToString())
		})
	if err != nil {
		panic(err)
	}

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
