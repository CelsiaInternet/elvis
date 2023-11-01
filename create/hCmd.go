package create

import (
	"fmt"

	. "github.com/cgalvisleon/elvis/utilities"
)

func MakeCmd(packageName, name string) error {
	path, err := MakeFolder("cmd", name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "Dockerfile", modelDockerfile, name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "main.go", modelMain, packageName, name)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCmd(packageName string) error {
	path := fmt.Sprintf(`./cmd/%s`, packageName)
	_, err := RemoveFile(path)
	if err != nil {
		return err
	}

	path = fmt.Sprintf(`./internal/service/%s`, packageName)
	_, err = RemoveFile(path)
	if err != nil {
		return err
	}

	path = fmt.Sprintf(`./internal/pkg/%s`, packageName)
	_, err = RemoveFile(path)
	if err != nil {
		return err
	}

	path = fmt.Sprintf(`./internal/rest/%s.http`, packageName)
	_, err = RemoveFile(path)
	if err != nil {
		return err
	}

	return nil
}
