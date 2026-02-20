package timezone

import (
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/celsiainternet/elvis/envar"
)

var loc *time.Location

func init() {
	timezone := envar.GetStr("America/Bogota", "TIMEZONE")
	var err error
	loc, err = time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
}

/**
* NowTime
* @return time.Time
* Remember to this function use ZONEINFO variable
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

/**
* Add
* @param d time.Duration
* @return time.Time
**/
func Add(d time.Duration) time.Time {
	return time.Now().In(loc).Add(d)
}

/**
* Location
* @return *time.Location
* Remember to this function use ZONEINFO variable
**/
func Location() *time.Location {
	return loc
}

/**
* Parse
* @param layout, value string
* @return time.Time, error
**/
func Parse(layout string, value string) (time.Time, error) {
	current, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		if strings.Count(value, "+") == 2 || strings.Count(value, "-") == 2 {
			layout = "2006-01-02 15:04:05 -0700 -0700"
		} else {
			layout = "2006-01-02 15:04:05 -0700"
		}

		return time.ParseInLocation(layout, value, loc)
	}

	return current, nil
}

/**
* FormatMDYYYY
* @param layout, value string
* @return string
**/
func FormatMDYYYY(value string) string {
	t, err := time.Parse("2006-01-02T15:04:05", value)
	if err != nil {
		return value
	}

	months := map[time.Month]string{
		time.January:   "Ene",
		time.February:  "Feb",
		time.March:     "Mar",
		time.April:     "Abr",
		time.May:       "May",
		time.June:      "Jun",
		time.July:      "Jul",
		time.August:    "Ago",
		time.September: "Sep",
		time.October:   "Oct",
		time.November:  "Nov",
		time.December:  "Dic",
	}

	return fmt.Sprintf("%s %02d %d", months[t.Month()], t.Day(), t.Year())
}
