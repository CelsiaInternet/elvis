package middleware

import (
	"time"

	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/strs"
)

type Request struct {
	Tag     string
	Day     int
	Hour    int
	Minute  int
	Seccond int
	Limit   int
}

var items map[string]int = make(map[string]int)

func CallRequests(tag string) Request {
	return Request{
		Tag:     tag,
		Day:     more(strs.Format(`%s-%d`, tag, time.Now().Unix()/86400), 86400),
		Hour:    more(strs.Format(`%s-%d`, tag, time.Now().Unix()/3600), 3600),
		Minute:  more(strs.Format(`%s-%d`, tag, time.Now().Unix()/60), 60),
		Seccond: more(strs.Format(`%s-%d`, tag, time.Now().Unix()/1), 1),
		Limit:   envar.EnvarInt(400, "REQUESTS_LIMIT"),
	}
}

func more(tag string, expiration time.Duration) int {
	value, ok := items[tag]
	if ok {
		value++
	} else {
		value = 1
	}

	items[tag] = value

	clean := func() {
		delete(items, tag)
	}

	duration := expiration * time.Second
	if duration != 0 {
		go time.AfterFunc(duration, clean)
	}

	return 0
}
