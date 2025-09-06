package main

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/flow"
)

func main() {
	flow.Load()

	test, err := flow.NewFlow("test", "1.0.0", "test", "test", func(ctx et.Json) (et.Item, error) {
		console.Debug("Respuesta desde test start, contexto:", ctx.ToString())

		return et.Item{
			Ok:     true,
			Result: ctx,
		}, nil
	}, 0, 0, 0, "test")
	if err != nil {
		console.Error(err)
		return
	}

	test.
		Step("Step 1", "Step 1", func(ctx et.Json) (et.Item, error) {
			console.Debug("Respuesta desde step 1, contexto:", ctx.ToString())

			return et.Item{
				Ok:     true,
				Result: ctx,
			}, nil
		}, false).
		IfElse("test == 'test'", 2, 3).
		Step("Step 2", "Step 2", func(ctx et.Json) (et.Item, error) {
			console.Debug("Respuesta desde step 2, con este contexto:", ctx.ToString())

			return et.Item{
				Ok:     true,
				Result: ctx,
			}, nil
		}, false).
		Step("Step 3", "Step 3", func(ctx et.Json) (et.Item, error) {
			console.Debug("Respuesta desde step 3, con este contexto:", ctx.ToString())

			return et.Item{
				Ok:     true,
				Result: ctx,
			}, nil
		}, false)

	console.Debug("Flow:", test.ToJson().ToString())

	result, err := flow.Run("", "test", 0, et.Json{
		"test": "test",
	})
	if err != nil {
		console.Error(err)
		return
	}

	console.Debug("Result:", result.ToString())
}
