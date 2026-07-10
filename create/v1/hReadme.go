package create

import "github.com/celsiainternet/elvis/file"

func MakeReadme(packageName string) error {
	_, err := file.MakeFile(".", "README.md", modelReadme, packageName)

	return err
}
