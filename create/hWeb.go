package create

import "github.com/cgalvisleon/elvis/utilities"

func MakeWeb(name string) error {
	_, err := utilities.MakeFolder(name, "web")
	if err != nil {
		return err
	}

	return nil
}
