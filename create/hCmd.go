package create

import . "github.com/cgalvisleon/elvis/utilities"

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
