package create

import (
	"fmt"

	. "github.com/cgalvisleon/elvis/utilities"
)

func MkPMicroservice(packageName, name, author, schema string) error {
	progressNext(10)
	err := MakeProject(name)
	if err != nil {
		return err
	}

	progressNext(10)
	err = MkMicroservice(packageName, name, schema)
	if err != nil {
		return err
	}

	progressNext(10)
	err = MakeWeb(name)
	if err != nil {
		return err
	}

	progressNext(60)
	_, err = Command([]string{
		fmt.Sprintf("cd ./%s", name),
		fmt.Sprintf("go mod init github.com/%s/%s", author, name),
	})
	if err != nil {
		return err
	}

	progressNext(10)

	return nil
}

func MkMicroservice(packageName, name, schema string) error {
	progressNext(10)
	err := MakeCmd(packageName, name)
	if err != nil {
		return err
	}

	progressNext(10)
	err = MakeDeployments(name)
	if err != nil {
		return err
	}

	progressNext(10)
	err = MakeInternal(packageName, name)
	if err != nil {
		return err
	}

	progressNext(10)
	schemaVar := Format(`Schema%s`, Titlecase(schema))
	err = MakePkg(name, schema, schemaVar)
	if err != nil {
		return err
	}

	progressNext(10)
	err = MakeScripts(name)
	if err != nil {
		return err
	}

	progressNext(40)
	err = MakeTest(name)
	if err != nil {
		return err
	}

	progressNext(10)

	return nil
}

func MkMolue(name, modelo, schema string) error {
	progressNext(10)
	schemaVar := Format(`Schema%s`, Titlecase(schema))
	err := MakeModel(name, modelo, schemaVar)
	if err != nil {
		return err
	}

	progressNext(90)

	return nil
}

func MkRpc(name string) error {
	progressNext(10)
	err := MakeRpc(name)
	if err != nil {
		return err
	}

	progressNext(90)

	return nil
}
