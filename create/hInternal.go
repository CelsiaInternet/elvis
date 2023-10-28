package create

import . "github.com/cgalvisleon/elvis/utilities"

func MakeInternal(packageName, name string) error {
	_, err := MakeFolder("internal", "data")
	if err != nil {
		return err
	}

	path, err := MakeFolder("internal", "service", name)
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "service.go", modelService, packageName, name)
	if err != nil {
		return err
	}

	path, err = MakeFolder("internal", "service", name, "v1")
	if err != nil {
		return err
	}

	_, err = MakeFile(path, "api.go", modelApi, packageName, name)
	if err != nil {
		return err
	}

	return nil
}
