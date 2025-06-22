package main

import (
	"errors"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/resilience"
)

var (
	atems = map[string]int{}
	total = 4
)

func main() {
	err := resilience.Load()
	if err != nil {
		logs.Log("resilience", "error", err)
		return
	}

	err = test("Hola Cristian")
	if err != nil {
		go resilience.Add("test", "test", test, "test")
	}

	time.Sleep(30 * 4 * time.Second)
	logs.Log("resilience", "finish")
}

func test(name string) error {
	_, ok := atems[name]
	if !ok {
		atems[name] = 1
	} else {
		atems[name]++
	}

	logs.Ping(atems[name], name)
	if atems[name] == total {
		return nil
	}

	return errors.New("error " + name)
}
