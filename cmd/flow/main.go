package main

import (
	"fmt"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/workflow"
)

func main() {
	test := workflow.New("test", "1.0.0", "test", "test", func(flow *workflow.Flow, ctx et.Json) (et.Json, error) {
		console.Debug("Respuesta desde test start, contexto:", ctx.ToString())
		atrib := fmt.Sprintf("step_%d", flow.Current)
		ctx.Set(atrib, "start")

		return ctx, nil
	}, "test").
		Resilence(3, 5*time.Second, 10*time.Minute).
		Step("Step 1", "Step 1", func(flow *workflow.Flow, ctx et.Json) (et.Json, error) {
			console.Debug("Respuesta desde step 1, contexto:", ctx.ToString())
			atrib := fmt.Sprintf("step_%d", flow.Current)
			ctx.Set(atrib, "step1")

			return ctx, nil
		}, false).
		IfElse("ctx.test == 'test'", 3, 2).
		Step("Step 2", "Step 2", func(flow *workflow.Flow, ctx et.Json) (et.Json, error) {
			console.Debug("Respuesta desde step 2, con este contexto:", ctx.ToString())
			atrib := fmt.Sprintf("step_%d", flow.Current)
			ctx.Set(atrib, "step2")

			return ctx, nil
		}, false).
		Step("Step 3", "Step 3", func(flow *workflow.Flow, ctx et.Json) (et.Json, error) {
			console.Debug("Respuesta desde step 3, con este contexto:", ctx.ToString())
			atrib := fmt.Sprintf("step_%d", flow.Current)
			ctx.Set(atrib, "step3")

			return ctx, nil
		}, false)

	console.Debug("Flow:", test.ToJson().ToString())

	result, err := workflow.Run("", "test", 0, et.Json{
		"test": "test",
	})
	if err != nil {
		console.Error(err)
	} else {
		console.Debug("Result:", result.ToString())
	}

	utility.AppWait()

	console.Debug("Fin de flow")
}
