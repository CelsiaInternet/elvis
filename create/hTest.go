package create

import "github.com/cgalvisleon/elvis/utilities"

func MakeTest(name string) error {
	_, err := utilities.MakeFolder("test")
	if err != nil {
		return err
	}

	return nil
}
