package create

import (
	"github.com/cgalvisleon/elvis/strs"
)

func MkProject(packageName, name, author, schema string) error {
	ProgressNext(10)
	err := MakeProject(name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MkMicroservice(packageName, name, schema)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeWeb(name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeReadme(name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeEnv(name)
	if err != nil {
		return err
	}

	ProgressNext(50)
	_, err = Command([]string{
		strs.Format("cd ./%s", name),
		strs.Format("go mod init github.com/%s/%s", author, name),
	})
	if err != nil {
		return err
	}

	ProgressNext(10)

	return nil
}

func MkMicroservice(packageName, name, schema string) error {
	ProgressNext(10)
	err := MakeCmd(packageName, name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeDeployments(name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeInternal(packageName, name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakePkg(name, schema)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeScripts(name)
	if err != nil {
		return err
	}

	ProgressNext(40)
	err = MakeTest(name)
	if err != nil {
		return err
	}

	ProgressNext(10)

	return nil
}

func MkMolue(name, modelo, schema string) error {
	ProgressNext(10)
	if len(schema) > 0 {
		err := MakeDbHandler(name, modelo, schema)
		if err != nil {
			return err
		}
	} else {
		err := MakeHandler(name, modelo)
		if err != nil {
			return err
		}
	}

	ProgressNext(90)

	return nil
}

func MkRpc(name string) error {
	ProgressNext(10)
	err := MakeRpc(name)
	if err != nil {
		return err
	}

	ProgressNext(90)

	return nil
}

/**
*
**/
func DeleteMicroservice(packageName string) error {
	ProgressNext(10)
	err := DeleteCmd(packageName)
	if err != nil {
		return err
	}

	ProgressNext(80)

	ProgressNext(10)

	return nil
}
