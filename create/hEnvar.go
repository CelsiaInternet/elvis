package create

import "github.com/cgalvisleon/elvis/file"

func MakeEnv(packageName string) error {
	_, _ = file.MakeFile(".", ".env", modelEnvar, packageName)

	return nil
}
