package create

import (
	"fmt"

	utl "github.com/cgalvisleon/elvis/utilities"
)

func MakeCmd(packageName, name string) error {
	path, err := utl.MakeFolder("cmd", name)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "Dockerfile", modelDockerfile, name)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "main.go", modelMain, packageName, name)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCmd(packageName string) error {
	path := fmt.Sprintf(`./cmd/%s`, packageName)
	_, err := utl.RemoveFile(path)
	if err != nil {
		return err
	}

	path = fmt.Sprintf(`./internal/service/%s`, packageName)
	_, err = utl.RemoveFile(path)
	if err != nil {
		return err
	}

	path = fmt.Sprintf(`./internal/pkg/%s`, packageName)
	_, err = utl.RemoveFile(path)
	if err != nil {
		return err
	}

	path = fmt.Sprintf(`./internal/rest/%s.http`, packageName)
	_, err = utl.RemoveFile(path)
	if err != nil {
		return err
	}

	return nil
}
