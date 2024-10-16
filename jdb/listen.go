package jdb

import (
	"log"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/lib/pq"
)

type HandlerListend func(res et.Json)

var _channels map[string]bool = map[string]bool{}

func (db *DB) defineListend(channels []string, lited HandlerListend) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	listener := pq.NewListener(db.Connection, 10*time.Second, time.Minute, reportProblem)
	for _, channel := range channels {
		err := listener.Listen(channel)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		select {
		case notification := <-listener.Notify:
			if notification != nil {
				result, err := et.ToJson(notification.Extra)
				if err != nil {
					logs.Alertm("defineListend: Not conver to Json")
				}

				result.Set("channel", notification.Channel)
				lited(result)
			}
		case <-time.After(90 * time.Second):
			go listener.Ping()
		}
	}
}

/**
* SetListen
* @param channels []string
* @param listen HandlerListend
**/
func (db *DB) SetListen(channels []string, listen HandlerListend) {
	for _, channel := range channels {
		if !_channels[channel] {
			_channels[channel] = true
			go db.defineListend(channels, listen)
		}
	}
}
