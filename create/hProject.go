package create

import utl "github.com/cgalvisleon/elvis/utilities"

func MakeProject(name string) error {
	_, err := utl.MakeFolder(name)
	if err != nil {
		return err
	}

	return nil
}
