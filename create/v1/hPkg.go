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

	block := strs.Format("\n\tif err := %s; err != nil {\n\t\treturn console.Panic(err)\n\t}\n", call)
	content = content[:idx] + block + content[idx:]

	return file.WriteFile(modelPath, content)
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

		path, err := file.MakeFolder("pkg", name)
		if err != nil {
			return err
		}

		fileName := strs.Format(`h%s.go`, modelo)
		_, err = file.MakeFile(path, fileName, modelDbHandler, name, modelo, schemaVar, strs.Uppcase(modelo), strs.Lowcase(modelo))
		if err != nil {
			return err
		}

		modelo = strs.Titlecase(modelo)
		_, err = file.MakeFile(path, "rpc.go", modelhRpc, name, modelo)
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

func MakeModel(packageName, modelo, schema string) error {
	path := strs.Format(`./pkg/%s`, packageName)

	if len(schema) > 0 {
		schemaVar := strs.Append("schema", strs.Titlecase(schema), "")
		_, err := file.MakeFile(path, "schema.go", modelSchema, packageName, schemaVar, schema)
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

func MakeRpc(name, modelo string) error {
	path, err := file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	modelo = strs.Titlecase(modelo)
	_, err = file.MakeFile(path, "rpc.go", modelhRpc, name, modelo)
	if err != nil {
		return err
	}

	return nil
}
