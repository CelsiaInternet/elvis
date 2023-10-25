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

	_, err = MakeFile(path, "event.go", modelEventGo, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "model.go", modelModelGo, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "hRpc.go", modelhRpcGo, name)
	if err != nil {
		return err
	}	

	_, err = MakeFile(path, "msg.go", modelMsgGo, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "router.go", modelRouterGo, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "schema.go", modelSchemaGo, name, schemaVar, schema)
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
	fileName := fmt.Sprintf(`h%s.go`, Titlecase(modelo))
	_, err = MakeFile(path, fileName, modelHandlerGo, name, modelo, schemaVar, Uppcase(modelo))
	if err != nil {
		return err
	}

	return nil
}
