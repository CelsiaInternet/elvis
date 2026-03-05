package main

import (
	"fmt"

	"github.com/celsiainternet/elvis/et"
)

func main() {
	data := et.Json{
		"a": et.Json{
			"b": et.Json{
				"c": "old_value",
			},
		},
	}
	data = et.SetNested(data, []string{"a", "b", "c"}, "value")
	fmt.Println(data.ToString())
}
