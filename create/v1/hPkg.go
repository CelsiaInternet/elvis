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
* schema) a un paquete existente en pkg/<packageName>.
* @param packageName, modelo, schema string
* @return error
**/
func MakeModel(packageName, modelo, schema string) error {
	path := strs.Format(`./pkg/%s`, packageName)
	if !file.ExistPath(path) {
		return fmt.Errorf("no existe el paquete %s; crea primero el microservicio", path)
	}

	if len(schema) > 0 {
		// El handler expone un metodo RPC en Services, definido en rpc.go
		_, err := file.MakeFile(path, "rpc.go", modelhRpc, packageName)
		if err != nil {
			return err
		}

		schemaVar := strs.Append("schema", strs.Titlecase(schema), "")
		_, err = file.MakeFile(path, "schema.go", modelSchema, packageName, schemaVar, schema)
		if err != nil {
			return err
		}

		modelo = strs.Titlecase(modelo)
		err = upsertModelInit(path, packageName, modelo)
		if err != nil {
			return err
		}

		fileName := strs.Format(`h%s.go`, modelo)
		_, err = file.MakeFile(path, fileName, modelDbHandler, packageName, modelo, schemaVar, strs.Uppcase(modelo), strs.Lowcase(modelo))
		if err != nil {
			return err
		}
	} else {
		modelo = strs.Titlecase(modelo)
		fileName := strs.Format(`h%s.go`, modelo)
		_, err := file.MakeFile(path, fileName, modelHandler, packageName, modelo, strs.Lowcase(modelo))
		if err != nil {
			return err
		}
	}

	return nil
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
