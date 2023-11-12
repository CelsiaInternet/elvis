package create

import "github.com/cgalvisleon/elvis/utility"

func MakeInternal(packageName, name string) error {
	_, err := utility.MakeFolder("internal", "data")
	if err != nil {
		return err
	}

	path, err := utility.MakeFolder("internal", "service", name)
	if err != nil {
		return err
	}

	_, err = utility.MakeFile(path, "service.go", modelService, packageName, name)
	if err != nil {
		return err
	}

	path, err = utility.MakeFolder("internal", "service", name, "v1")
	if err != nil {
		return err
	}

	_, err = utility.MakeFile(path, "api.go", modelApi, packageName, name)
	if err != nil {
		return err
	}

	return nil
}
