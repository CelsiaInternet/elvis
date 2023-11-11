package create

import utl "github.com/cgalvisleon/elvis/utilities"

func MakeWeb(name string) error {
	_, err := utl.MakeFolder(name, "web")
	if err != nil {
		return err
	}

	return nil
}
