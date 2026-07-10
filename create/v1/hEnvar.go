package create

import "github.com/celsiainternet/elvis/file"

func MakeEnv(packageName string) error {
	_, err := file.MakeFile(".", ".env", modelEnvar, packageName)

	return err
}
