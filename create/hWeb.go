package create

import "github.com/cgalvisleon/elvis/utility"

func MakeWeb(name string) error {
	_, err := utility.MakeFolder(name, "web")
	if err != nil {
		return err
	}

	return nil
}
