package timezone

import (
	"time"

	"github.com/cgalvisleon/elvis/envar"
)

var loc *time.Location

/**
* NowTime
* @return time.Time
**/
func NowTime() time.Time {
	return time.Now().In(loc)
}

/**
* Now
* @return string
**/
func Now() string {
	return NowTime().Format("2006/01/02 15:04:05")
}

func init() {
	timeZona := envar.GetStr("America/Bogota", "TIME_ZONE")
	loc, _ = time.LoadLocation(timeZona)
}
