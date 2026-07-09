package main

import (
	"fmt"

	"github.com/celsiainternet/elvis/et"
)

func main() {

	v := et.Json{
		"document": []string{
			"https://objectstorage.sa-bogota-1.oraclecloud.com/n/axyxbjtesgqh/b/prueba/o//soportes-blacklist568d5278-b11c-4e81-96cf-788521baa317..pdf",
		},
	}

	fmt.Println(v.ToString())

	/*
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
	*/
}
