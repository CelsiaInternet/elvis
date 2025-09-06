package main

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/flow"
)

func main() {
	flow.Load()

	test, err := flow.NewFlow("test", "1.0.0", "test", "test", func(ctx et.Json) (et.Item, error) {
		console.Debug("Ejecutando desde test, con este contexto:", ctx.ToString())

		return et.Item{
			Ok:     true,
			Result: ctx,
		}, nil
	}, 0, 0, 0, "test")
	if err != nil {
		console.Error(err)
		return
	}

	result, err := test.Run("test", et.Json{
		"test": "test",
	})
	if err != nil {
		console.Error(err)
		return
	}

	console.Debug("Result:", result.ToString())
}
