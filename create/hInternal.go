package create

import utl "github.com/cgalvisleon/elvis/utilities"

func MakeInternal(packageName, name string) error {
	_, err := utl.MakeFolder("internal", "data")
	if err != nil {
		return err
	}

	path, err := utl.MakeFolder("internal", "service", name)
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "service.go", modelService, packageName, name)
	if err != nil {
		return err
	}

	path, err = utl.MakeFolder("internal", "service", name, "v1")
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, "api.go", modelApi, packageName, name)
	if err != nil {
		return err
	}

	return nil
}
