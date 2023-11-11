package create

import (
	"fmt"

	utl "github.com/cgalvisleon/elvis/utilities"
)

func MakePkg(name, schema, schemaVar string) error {
	path, err := utl.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "event.go", modelEvent, name)
	if err != nil {
		return err
	}

	modelo := utl.Titlecase(name)
	_, err = utl.MakeFile(path, "model.go", modelModel, name, modelo)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "msg.go", modelMsg, name)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "controller.go", modelController, name)
	if err != nil {
		return err
	}

	title := utl.Titlecase(name)
	_, err = utl.MakeFile(path, "router.go", modelRouter, name, title)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "schema.go", modelSchema, name, schemaVar, schema)
	if err != nil {
		return err
	}

	return MakeModel(name, name, schemaVar)
}

func MakeModel(name, modelo, schemaVar string) error {
	path, err := utl.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	modelo = utl.Titlecase(modelo)
	fileName := fmt.Sprintf(`h%s.go`, modelo)
	_, err = utl.MakeFile(path, fileName, modelHandler, name, modelo, schemaVar, utl.Uppcase(modelo), utl.Lowcase(modelo))
	if err != nil {
		return err
	}

	return nil
}

func MakeRpc(name string) error {
	path, err := utl.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "hRpc.go", modelhRpc, name)
	if err != nil {
		return err
	}

	return nil
}
