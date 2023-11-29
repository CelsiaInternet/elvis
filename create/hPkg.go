package create

import (
	"fmt"

	"github.com/cgalvisleon/elvis/utility"
)

func MakePkg(name, schema, schemaVar string) error {
	path, err := utility.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = utility.MakeFile(path, "event.go", modelEvent, name)
	if err != nil {
		return err
	}

	modelo := utility.Titlecase(name)
	_, err = utility.MakeFile(path, "model.go", modelModel, name, modelo)
	if err != nil {
		return err
	}

	_, err = utility.MakeFile(path, "msg.go", modelMsg, name)
	if err != nil {
		return err
	}

	_, err = utility.MakeFile(path, "controller.go", modelController, name)
	if err != nil {
		return err
	}

	title := utility.Titlecase(name)
	_, err = utility.MakeFile(path, "router.go", modelRouter, name, title)
	if err != nil {
		return err
	}

	if len(schema) > 0 {
		_, err = utility.MakeFile(path, "schema.go", modelSchema, name, schemaVar, schema)
		if err != nil {
			return err
		}
	}

	return MakeModel(name, name, schemaVar)
}

func MakeModel(name, modelo, schemaVar string) error {
	path, err := utility.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	modelo = utility.Titlecase(modelo)
	fileName := fmt.Sprintf(`h%s.go`, modelo)
	_, err = utility.MakeFile(path, fileName, modelHandler, name, modelo, schemaVar, utility.Uppcase(modelo), utility.Lowcase(modelo))
	if err != nil {
		return err
	}

	return nil
}

func MakeRpc(name string) error {
	path, err := utility.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = utility.MakeFile(path, "hRpc.go", modelhRpc, name)
	if err != nil {
		return err
	}

	return nil
}
