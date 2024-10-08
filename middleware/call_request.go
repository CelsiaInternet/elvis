package middleware

import (
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/timezone"
)

type Request struct {
	Tag     string
	Day     int
	Hour    int
	Minute  int
	Seccond int
	Limit   int
}

func callRequests(tag string) Request {
	now := timezone.NowTime().Unix()
	return Request{
		Tag:     tag,
		Day:     cache.More(strs.Format(`%s-%d`, tag, now/86400), 86400),
		Hour:    cache.More(strs.Format(`%s-%d`, tag, now/3600), 3600),
		Minute:  cache.More(strs.Format(`%s-%d`, tag, now/60), 60),
		Seccond: cache.More(strs.Format(`%s-%d`, tag, now/1), 1),
		Limit:   envar.GetInt(400, "REQUESTS_LIMIT"),
	}
}

var items map[string]int = make(map[string]int)

func localRequests(tag string) Request {
	now := timezone.NowTime().Unix()
	return Request{
		Tag:     tag,
		Day:     more(strs.Format(`%s-%d`, tag, now/86400), 86400),
		Hour:    more(strs.Format(`%s-%d`, tag, now/3600), 3600),
		Minute:  more(strs.Format(`%s-%d`, tag, now/60), 60),
		Seccond: more(strs.Format(`%s-%d`, tag, now/1), 1),
		Limit:   envar.GetInt(400, "REQUESTS_LIMIT"),
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
