package create

import "github.com/cgalvisleon/elvis/utilities"

func MakeInternal(packageName, name string) error {
	_, err := utilities.MakeFolder("internal", "data")
	if err != nil {
		return err
	}

	path, err := utilities.MakeFolder("internal", "service", name)
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, "service.go", modelService, packageName, name)
	if err != nil {
		return err
	}

	path, err = utilities.MakeFolder("internal", "service", name, "v1")
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, "api.go", modelApi, packageName, name)
	if err != nil {
		return err
	}

	return nil
}
