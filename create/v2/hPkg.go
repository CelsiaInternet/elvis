package create

import (
	"github.com/celsiainternet/elvis/file"
	"github.com/celsiainternet/elvis/strs"
)

func MakePkg(name, schema string) error {
	pkgPath, err := file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(pkgPath, "event.go", modelEvent, name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(pkgPath, "config.go", modelConfig, name)
	if err != nil {
		return err
	}

	if len(schema) > 0 {
		_, err = file.MakeFile(pkgPath, "controller.go", modelDbController, name)
		if err != nil {
			return err
		}

		modelo := strs.Titlecase(name)
		routerFileName := strs.Format(`router-%s.go`, modelo)
		_, err = file.MakeFile(pkgPath, routerFileName, modelDbModelRouter, name, modelo, strs.Uppcase(modelo), strs.Lowcase(modelo))
		if err != nil {
			return err
		}

		modelo = strs.Titlecase(modelo)
		_, err = file.MakeFile(pkgPath, "rpc.go", modelhRpc, name, modelo)
		if err != nil {
			return err
		}

		title := strs.Titlecase(name)
		_, err = file.MakeFile(pkgPath, "router.go", modelDbRouter, name, title)
		if err != nil {
			return err
		}
	} else {
		_, err = file.MakeFile(pkgPath, "controller.go", modelController, name)
		if err != nil {
			return err
		}

		modelo := strs.Titlecase(name)
		routerFileName := strs.Format(`router-%s.go`, modelo)
		_, err = file.MakeFile(pkgPath, routerFileName, modelDbModelRouter, name, modelo, strs.Uppcase(modelo), strs.Lowcase(modelo))
		if err != nil {
			return err
		}

		_, err = file.MakeFile(pkgPath, "router.go", modelRouter, name, strs.Lowcase(name))
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
		_, _ = file.MakeFile(path, "schema.go", modelSchema, packageName, schemaVar, schema)

		modelo := strs.Titlecase(modelo)
		_, _ = file.MakeFile(path, "model.go", modelModel, packageName, modelo)

		modelo = strs.Titlecase(modelo)
		fileName := strs.Format(`h%s.go`, modelo)
		_, err := file.MakeFile(path, fileName, modelDbHandler, packageName, modelo, schemaVar, strs.Uppcase(modelo), strs.Lowcase(modelo))
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
