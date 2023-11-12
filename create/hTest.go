package create

import "github.com/cgalvisleon/elvis/utility"

func MakeTest(name string) error {
	_, err := utility.MakeFolder("test")
	if err != nil {
		return err
	}

	return nil
}
