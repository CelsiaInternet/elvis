package timezone

import (
	"time"
)

/**
* NowTime
* @return time.Time
* Remember to this function use ZONEINFO variable
**/
func NowTime() time.Time {
	return time.Now()
}

/**
* Now
* @return string
**/
func Now() string {
	return NowTime().Format("2006/01/02 15:04:05")
}
