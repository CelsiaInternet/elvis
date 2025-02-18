package jdb

import (
	"log"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/lib/pq"
)

type HandlerListend func(res et.Json)

func (db *DB) defineListend(channels map[string]HandlerListend) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	listener := pq.NewListener(db.Connection, 10*time.Second, time.Minute, reportProblem)
	for channel := range channels {
		err := listener.Listen(channel)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		select {
		case notification := <-listener.Notify:
			if notification != nil {
				if db.channels[notification.Channel] == nil {
					continue
				}

				result, err := et.ToJson(notification.Extra)
				if err != nil {
					logs.Alertm("defineListend: Not conver to Json")
				}

				result.Set("channel", notification.Channel)
				db.channels[notification.Channel](result)
			}
		case <-time.After(90 * time.Second):
			go listener.Ping()
		}
	}
}

/**
* SetListen
* @param channels map[string]HandlerListend
**/
func (db *DB) SetListen(channels map[string]HandlerListend) {
	go db.defineListend(channels)
}
