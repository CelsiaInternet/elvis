package create

import "github.com/celsiainternet/elvis/file"

func MakeGitignore(packageName string) error {
	_, err := file.MakeFile(".", ".gitignore", modelGitignore)

	return err
}
