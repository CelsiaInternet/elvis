package create

import "github.com/cgalvisleon/elvis/utility"

func MakeProject(name string) error {
	_, err := utility.MakeFolder(name)
	if err != nil {
		return err
	}

	return nil
}
