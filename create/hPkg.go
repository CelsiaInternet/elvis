package create

import (
	"fmt"

	. "github.com/cgalvisleon/elvis/utilities"
)

func MakePkg(name, schema, schemaVar string) error {
	path, err := MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "event.go", modelEvent, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "model.go", modelModel, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "msg.go", modelMsg, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "controller.go", modelController, name)
	if err != nil {
		return err
	}

	title := Titlecase(name)
	_, err = MakeFile(path, "router.go", modelRouter, name, title)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "schema.go", modelSchema, name, schemaVar, schema)
	if err != nil {
		return err
	}

	return MakeModel(name, name, schemaVar)
}

func MakeModel(name, modelo, schemaVar string) error {
	path, err := MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	modelo = Titlecase(modelo)
	fileName := fmt.Sprintf(`h%s.go`, modelo)
	_, err = MakeFile(path, fileName, modelHandler, name, modelo, schemaVar, Uppcase(modelo))
	if err != nil {
		return err
	}

	return nil
}

func MakeRpc(name string) error {
	path, err := MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "hRpc.go", modelhRpc, name)
	if err != nil {
		return err
	}

	return nil
}
