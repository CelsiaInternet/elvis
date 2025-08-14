package create

import (
	"github.com/celsiainternet/elvis/file"
	"github.com/celsiainternet/elvis/strs"
)

func MakeInternal(packageName, name, schema string) error {
	modelsPath, err := file.MakeFolder("internal", "models", name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(modelsPath, "msg.go", modelMsg, name)
	if err != nil {
		return err
	}

	if len(schema) > 0 {
		schemaVar := strs.Append("schema", strs.Titlecase(schema), "")
		_, err = file.MakeFile(modelsPath, "schema.go", modelSchema, name, schemaVar, schema)
		if err != nil {
			return err
		}

		modelo := strs.Titlecase(name)
		modelFileName := strs.Format(`%s.go`, modelo)
		_, err = file.MakeFile(modelsPath, modelFileName, modelModel, name, modelo)
		if err != nil {
			return err
		}
	}

	servicePath, err := file.MakeFolder("internal", "service", name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(servicePath, "service.go", modelService, packageName, name)
	if err != nil {
		return err
	}

	v1Path, err := file.MakeFolder("internal", "service", name, "v1")
	if err != nil {
		return err
	}

	if len(schema) > 0 {
		_, err = file.MakeFile(v1Path, "api.go", modelDbApi, packageName, name)
		if err != nil {
			return err
		}
	} else {
		_, err = file.MakeFile(v1Path, "api.go", modelApi, packageName, name)
		if err != nil {
			return err
		}
	}

	return nil
}
