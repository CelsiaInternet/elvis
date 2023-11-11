package create

import "github.com/cgalvisleon/elvis/utilities"

func MakeProject(name string) error {
	_, err := utilities.MakeFolder(name)
	if err != nil {
		return err
	}

	return nil
}
