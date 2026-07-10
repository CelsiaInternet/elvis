package create

import (
	"github.com/celsiainternet/elvis/file"
	"github.com/celsiainternet/elvis/strs"
)

/**
* MakeCmd
* @param packageName, name string
* @return error
**/
func MakeCmd(packageName, name string) error {
	path, err := file.MakeFolder("cmd", name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "Dockerfile", modelDockerfile, name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "main.go", modelMain, packageName, name)
	if err != nil {
		return err
	}

	return nil
}

/**
* DeleteCmd
* @param name string
* @return error
**/
func DeleteCmd(name string) error {
	path := strs.Format(`./cmd/%s`, name)
	_, err := file.RemoveFile(path)
	if err != nil {
		return err
	}

	path = strs.Format(`./deployments/%s`, strs.Lowcase(name))
	_, err = file.RemoveFile(path)
	if err != nil {
		return err
	}

	path = strs.Format(`./internal/service/%s`, name)
	_, err = file.RemoveFile(path)
	if err != nil {
		return err
	}

	path = strs.Format(`./pkg/%s`, name)
	_, err = file.RemoveFile(path)
	if err != nil {
		return err
	}

	path = strs.Format(`./scripts/%s.http`, name)
	_, err = file.RemoveFile(path)
	if err != nil {
		return err
	}

	path = strs.Format(`./www/%s`, name)
	_, err = file.RemoveFile(path)
	if err != nil {
		return err
	}

	return nil
}
