package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/envar"
)

var DefaultTelemetry func(next http.Handler) http.Handler

type Request struct {
	Tag     string
	Day     int
	Hour    int
	Minute  int
	Seccond int
	Limit   int
}

func CallRequests(tag string) Request {
	return Request{
		Tag:     tag,
		Day:     cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/86400), 86400),
		Hour:    cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/3600), 3600),
		Minute:  cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/60), 60),
		Seccond: cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/1), 1),
		Limit:   envar.EnvarInt(1000, "REQUESTS_LIMIT"),
	}
}
