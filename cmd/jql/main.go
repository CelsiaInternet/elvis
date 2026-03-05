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
	data.SetNested([]string{"a", "b", "c"}, "value")
	data.SetNested([]string{"a", "b", "d"}, "value2")
	fmt.Println(data.ToString())
}
