package create

import (
	"fmt"

	"github.com/cgalvisleon/elvis/utilities"
)

func MakePkg(name, schema, schemaVar string) error {
	path, err := utilities.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, "event.go", modelEvent, name)
	if err != nil {
		return err
	}

	modelo := utilities.Titlecase(name)
	_, err = utilities.MakeFile(path, "model.go", modelModel, name, modelo)
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, "msg.go", modelMsg, name)
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, "controller.go", modelController, name)
	if err != nil {
		return err
	}

	title := utilities.Titlecase(name)
	_, err = utilities.MakeFile(path, "router.go", modelRouter, name, title)
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, "schema.go", modelSchema, name, schemaVar, schema)
	if err != nil {
		return err
	}

	return MakeModel(name, name, schemaVar)
}

func MakeModel(name, modelo, schemaVar string) error {
	path, err := utilities.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	modelo = utilities.Titlecase(modelo)
	fileName := fmt.Sprintf(`h%s.go`, modelo)
	_, err = utilities.MakeFile(path, fileName, modelHandler, name, modelo, schemaVar, utilities.Uppcase(modelo), utilities.Lowcase(modelo))
	if err != nil {
		return err
	}

	return nil
}

func MakeRpc(name string) error {
	path, err := utilities.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, "hRpc.go", modelhRpc, name)
	if err != nil {
		return err
	}

	return nil
}
