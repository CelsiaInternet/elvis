package health

import (
	"fmt"

	"github.com/celsiainternet/elvis/et"
)

type Check struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}

/**
* ToJson
* @return Json
**/
func (s *Check) ToJson() et.Json {
	return et.Json{
		"ok":  s.Ok,
		"msg": s.Msg,
	}
}

type Services map[string]func() bool

/**
* Checked
* @param checks map[string]func() bool
* @return Check
**/
func Checked(checks Services) Check {
	for name, check := range checks {
		if !check() {
			return Check{Ok: false, Msg: fmt.Sprintf("%s is not ok", name)}
		}
	}

	return Check{Ok: true, Msg: "ok"}
}
