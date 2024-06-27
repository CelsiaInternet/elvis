package create

import "github.com/cgalvisleon/elvis/file"

func MakeEnv(packageName string) error {
	_, _ = file.MakeFile(".", ".env", modelEnvar, packageName)

	_, _ = file.MakeFile(".", ".env.local", modelEnvar, packageName)

	_, _ = file.MakeFile(".", ".env.prd", modelEnvar, packageName)

	_, _ = file.MakeFile(".", ".env.qa", modelEnvar, packageName)

	return nil
}
