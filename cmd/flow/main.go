package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/workflow"
)

func main() {
	workflow.New("ventas", "1.0.0", "Flujo de ventas", "flujo de ventas", func(flow *workflow.Instance, ctx et.Json) (et.Json, error) {
		console.Debug("Respuesta desde step 0, contexto:", ctx.ToString())
		atrib := fmt.Sprintf("step_%d", flow.Current)
		ctx.Set(atrib, "step0")

		return ctx, nil
	}, true, "test").
		// Debug().
		Retention(24*time.Hour).
		Resilence(3, 3*time.Second, "test", "1").
		Step("Step 1", "Step 1", func(flow *workflow.Instance, ctx et.Json) (et.Json, error) {
			console.Debug("Respuesta desde step 1, contexto:", ctx.ToString())
			atrib := fmt.Sprintf("step_%d", flow.Current)
			ctx.Set(atrib, "step1")

			// flow.Done()
			// flow.Stop()
			// flow.Goto(2)

			time.Sleep(3 * time.Second)

			return ctx, nil
		}, false).
		IfElse(`test == "test"`, 3, 2).
		Step("Step 2", "Step 2", func(flow *workflow.Instance, ctx et.Json) (et.Json, error) {
			console.Debug("Respuesta desde step 2, con este contexto:", ctx.ToString())
			atrib := fmt.Sprintf("step_%d", flow.Current)
			ctx.Set(atrib, "step2")

			// guardar en el Oss

			return ctx, nil
		}, true).
		Rollback(func(flow *workflow.Instance, ctx et.Json) (et.Json, error) {
			console.Debug("Respuesta desde rollback 2, con este contexto:", ctx.ToString())
			atrib := fmt.Sprintf("rollback_%d", flow.Current)
			ctx.Set(atrib, "step2")

			return ctx, nil
		}).
		Step("Step 3", "Step 3", func(flow *workflow.Instance, ctx et.Json) (et.Json, error) {
			console.Debug("Respuesta desde step 3, con este contexto:", ctx.ToString())
			atrib := fmt.Sprintf("step_%d", flow.Current)
			ctx.Set(atrib, "step3")

			return ctx, nil
		}, false)

	console.Debug("")
	console.Debug("")

	result, err := workflow.Run("1234", "ventas", 0, et.Json{
		"cedula": "91499023",
	}, et.Json{
		"test": "test",
	}, "test")
	if err != nil {
		console.Error(err)
	} else {
		console.Debug("Result 1:", result.ToString())
	}

	// go func() {
	// 	result, err := workflow.Run("", "ventas", 2, et.Json{
	// 		"cedula": "91499023",
	// 	}, et.Json{
	// 		"test": "test",
	// 	}, "test")
	// 	if err != nil {
	// 		console.Error(err)
	// 	} else {
	// 		console.Debug("Result 2:", result.ToString())
	// 	}
	// }()

	// result, err := workflow.Continue("", et.Json{
	// 	"cedula": "91499023",
	// }, et.Json{
	// 	"test": "test",
	// }, "test")
	// if err != nil {
	// 	console.Error(err)
	// } else {
	// 	console.Debug("Result 2:", result.ToString())
	// }

	// go func() {
	// 	result, err := workflow.Run("", "ventas", 2, et.Json{
	// 		"cedula": "91499023",
	// 	}, et.Json{
	// 		"test": "test",
	// 	}, "test")
	// 	if err != nil {
	// 		console.Error(err)
	// 	} else {
	// 		console.Debug("Result:", result.ToString())
	// 	}
	// }()

	utility.AppWait()

	console.Debug("Fin de flow")

}

func HttpVenta(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	tag := r.PathValue("tag")
	serviceId := body.Str("serviceId")
	tags := et.Json{
		"cedula": "91499023",
		"codigo": "112342",
	}
	step := body.Int("step")
	createdBy := claim.ClientName(r)
	result, err := workflow.Run(serviceId, tag, step, tags, body, createdBy)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok:     true,
		Result: result,
	})

}
