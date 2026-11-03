package create

import (
	"fmt"
	"strings"

	"github.com/celsiainternet/elvis/file"
	"github.com/celsiainternet/elvis/strs"
)

/**
* upsertModelInit: Writes pkg/<packageName>/model.go the first time a
* model is added, and appends a Define<modelo>(db) call to initModels
* on every subsequent call instead of leaving the file untouched, since
* file.MakeFile never overwrites a file that already exists.
* @param path, packageName, modelo string
* @return error
**/
func upsertModelInit(path, packageName, modelo string) error {
	modelPath := strs.Format(`%s/model.go`, path)
	if !file.ExistPath(modelPath) {
		_, err := file.MakeFile(path, "model.go", modelModel, packageName, modelo)
		return err
	}

	content, err := file.ReadFile(modelPath)
	if err != nil {
		return err
	}

	call := strs.Format(`Define%s(db)`, modelo)
	if strings.Contains(content, call) {
		return nil
	}

	marker := "\n\treturn nil\n}"
	idx := strings.LastIndex(content, marker)
	if idx == -1 {
		return fmt.Errorf("no se pudo actualizar initModels en %s", modelPath)
	}

	block := strs.Format("\n\tif err := %s; err != nil {\n\t\treturn err\n\t}\n", call)
	content = content[:idx] + block + content[idx:]

	return file.WriteFile(modelPath, content)
}

/**
* MakePkg: Genera pkg/<name>: config, eventos, mensajes, RPC, controller,
* router y, si hay schema, el primer modelo linq con su CRUD.
* @param name, schema string
* @return error
**/
/**
* routesAnchor: Linea que el generador escribe en Routes() justo despues de
* las rutas; las rutas de los modelos nuevos se insertan antes de ella.
**/
const routesAnchor = "\n\tctx := context.Background()\n"

/**
* upsertRoutes: Agrega a Routes() de pkg/<name>/router.go las rutas de un
* modelo nuevo. No hace nada si el handler ya esta enrutado. Devuelve false
* si router.go fue editado y ya no tiene la linea de anclaje: en ese caso
* las rutas hay que copiarlas a mano.
* @param path, handler, routes string
* @return bool, error
**/
func upsertRoutes(path, handler, routes string) (bool, error) {
	routerPath := strs.Format(`%s/router.go`, path)
	content, err := file.ReadFile(routerPath)
	if err != nil {
		return false, err
	}

	if strings.Contains(content, strs.Format(`rt.%s,`, handler)) {
		return true, nil
	}

	idx := strings.Index(content, routesAnchor)
	if idx == -1 {
		return false, nil
	}

	content = content[:idx] + "\n" + routes + content[idx:]
	return true, file.WriteFile(routerPath, content)
}

/**
* hasModels: Indica si el paquete ya inicializa modelos (Controller.Init
* llama a initModels), es decir, si se creo con schema y tiene base de datos.
* @param path string
* @return bool
**/
func hasModels(path string) bool {
	content, err := file.ReadFile(strs.Format(`%s/controller.go`, path))
	if err != nil {
		return false
	}

	return strings.Contains(content, "initModels(c.Db)")
}

func MakePkg(name, schema string) error {
	path, err := file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "event.go", modelEvent, name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "msg.go", modelMsg, name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "config.go", modelConfig, name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "rpc.go", modelhRpc, name)
	if err != nil {
		return err
	}

	if len(schema) > 0 {
		_, err = file.MakeFile(path, "controller.go", modelDbController, name)
		if err != nil {
			return err
		}

		schemaVar := strs.Append("schema", strs.Titlecase(schema), "")
		_, err = file.MakeFile(path, "schema.go", modelSchema, name, schemaVar, schema)
		if err != nil {
			return err
		}

		modelo := strs.Titlecase(name)
		err = upsertModelInit(path, name, modelo)
		if err != nil {
			return err
		}

		fileName := strs.Format(`h%s.go`, modelo)
		_, err = file.MakeFile(path, fileName, modelDbHandler, name, modelo, schemaVar, strs.Uppcase(modelo), strs.Lowcase(modelo))
		if err != nil {
			return err
		}

		title := strs.Titlecase(name)
		_, err = file.MakeFile(path, "router.go", modelDbRouter, name, title)
		if err != nil {
			return err
		}
	} else {
		_, err = file.MakeFile(path, "controller.go", modelController, name)
		if err != nil {
			return err
		}

		modelo := strs.Titlecase(name)
		fileName := strs.Format(`h%s.go`, modelo)
		_, err = file.MakeFile(path, fileName, modelHandler, name, modelo, strs.Lowcase(modelo))
		if err != nil {
			return err
		}

		_, err = file.MakeFile(path, "router.go", modelRouter, name, strs.Lowcase(name))
		if err != nil {
			return err
		}
	}

	return nil
}

/**
* MakeModel: Agrega un modelo (con schema, linq + CRUD) o un handler (sin
* schema) a un paquete existente en pkg/<packageName> y registra sus rutas
* en router.go. Devuelve false si las rutas no se pudieron registrar solas.
* @param packageName, modelo, schema string
* @return bool, error
**/
func MakeModel(packageName, modelo, schema string) (bool, error) {
	path := strs.Format(`./pkg/%s`, packageName)
	if !file.ExistPath(path) {
		return false, fmt.Errorf("no existe el paquete %s; crea primero el microservicio", path)
	}

	if len(schema) > 0 && !hasModels(path) {
		return false, fmt.Errorf(`el paquete %s se creo sin schema y no tiene base de datos; antes de agregar un modelo con schema:
  1. controller.go: agrega el campo Db *jdb.DB a Controller y, en Init, llama a initModels(c.Db) y devuelve su error
  2. internal/service/<servicio>/v1/api.go: carga la base con jdb.Load(), pasala como pkg.Controller{Db: db} y cierrala en Close()
  3. vuelve a ejecutar la opcion Modelo`, path)
	}

	if len(schema) > 0 {
		// El handler expone un metodo RPC en Services, definido en rpc.go
		_, err := file.MakeFile(path, "rpc.go", modelhRpc, packageName)
		if err != nil {
			return false, err
		}

		schemaVar := strs.Append("schema", strs.Titlecase(schema), "")
		_, err = file.MakeFile(path, "schema.go", modelSchema, packageName, schemaVar, schema)
		if err != nil {
			return false, err
		}

		modelo = strs.Titlecase(modelo)
		err = upsertModelInit(path, packageName, modelo)
		if err != nil {
			return false, err
		}

		fileName := strs.Format(`h%s.go`, modelo)
		_, err = file.MakeFile(path, fileName, modelDbHandler, packageName, modelo, schemaVar, strs.Uppcase(modelo), strs.Lowcase(modelo))
		if err != nil {
			return false, err
		}

		routes := fillTemplate(modelDbRoutes, modelo, strs.Lowcase(modelo))
		return upsertRoutes(path, strs.Format(`get%sById`, modelo), routes)
	}

	modelo = strs.Titlecase(modelo)
	fileName := strs.Format(`h%s.go`, modelo)
	_, err := file.MakeFile(path, fileName, modelHandler, packageName, modelo, strs.Lowcase(modelo))
	if err != nil {
		return false, err
	}

	routes := fillTemplate(modelRoutes, modelo, strs.Lowcase(modelo))
	return upsertRoutes(path, strs.Lowcase(modelo), routes)
}

/**
* MakeRpc: Genera pkg/<name>/rpc.go (Services + StartRpcServer) si no existe.
* @param name string
* @return error
**/
func MakeRpc(name string) error {
	path := strs.Format(`./pkg/%s`, name)
	if !file.ExistPath(path) {
		return fmt.Errorf("no existe el paquete %s; crea primero el microservicio", path)
	}

	_, err := file.MakeFile(path, "rpc.go", modelhRpc, name)
	if err != nil {
		return err
	}

	return nil
}

/**
* fillTemplate: Sustituye $1, $2, ... en model, igual que file.MakeFile.
* @param model string, args ...string
* @return string
**/
func fillTemplate(model string, args ...string) string {
	result := model
	for i, v := range args {
		result = strings.ReplaceAll(result, strs.Format(`$%d`, i+1), v)
	}

	return result
}
