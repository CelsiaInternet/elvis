package main

import (
	"errors"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/resilience"
)

var (
	atems = map[string]int{}
	total = 6
)

func main() {
	resilience.Load("test")

	err := test("test")
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

	logs.Ping(atems[name])
	if atems[name] == total {
		return nil
	}

	return errors.New("error " + name)
}
