package create

import "github.com/cgalvisleon/elvis/file"

func MakeReadme(packageName string) error {
	_, _ = file.MakeFile(".", "README.md", modelReadme, packageName)

	return nil
}
