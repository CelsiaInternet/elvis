package timezone

import (
	"os"
	"time"
)

var loc *time.Location

/**
* NowTime
* @return time.Time
**/
func NowTime() time.Time {
	if loc == nil {
		loc = loadLocation()
	}

	return time.Now().In(loc)
}

/**
* loadLocation
* @return *time.Location
**/
func loadLocation() *time.Location {
	if loc != nil {
		return loc
	}

	timeZona := os.Getenv("TIME_ZONE")
	if timeZona == "" {
		timeZona = "America/Bogota"
	}

	var err error
	loc, err = time.LoadLocation(timeZona)
	if err != nil {
		loc = time.UTC
	}

	return loc
}

/**
* Now
* @return string
**/
func Now() string {
	return NowTime().Format("2006/01/02 15:04:05")
}

func init() {
	loadLocation()
}
