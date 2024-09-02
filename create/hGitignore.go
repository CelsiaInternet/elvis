package create

import "github.com/cgalvisleon/elvis/file"

func MakeGitignore(packageName string) error {
	_, _ = file.MakeFile(".", ".gitignore", modelGitignore)

	return nil
}
