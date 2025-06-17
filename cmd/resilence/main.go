package main

import (
	"errors"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/resilience"
)

var (
	atems = map[string]int{}
	total = 3
)

func main() {
	resilience.Load("test")

	time.Sleep(10 * time.Second)
	logs.Log("resilience", "finish")
}

func test2(name string) error {
	_, ok := atems["test2"]
	if !ok {
		atems["test2"] = 1
		return errors.New("test2")
	}
	atems["test2"]++

	if atems["test2"] == total {
		return nil
	}

	return errors.New("test2")
}

func test(name string) error {
	_, ok := atems["test"]
	if !ok {
		atems["test"] = 1
		return errors.New("test")
	}
	atems["test"]++

	if atems["test"] == total {
		return nil
	}

	return errors.New("test")
}
