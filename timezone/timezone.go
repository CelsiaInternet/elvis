package timezone

import (
	"time"
)

var loc *time.Location

/**
* NowTime
* @return time.Time
* Remember to this function use ZONEINFO variable
**/
func NowTime() time.Time {
	if loc != nil {
		var err error
		loc, err = time.LoadLocation("America/Bogota")
		if err != nil {
			loc = time.Local
		}
	}

	return time.Now().In(loc)
}

/**
* Now
* @return string
**/
func Now() string {
	return NowTime().Format("2006/01/02 15:04:05")
}
