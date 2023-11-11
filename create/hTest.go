package create

import utl "github.com/cgalvisleon/elvis/utilities"

func MakeTest(name string) error {
	_, err := utl.MakeFolder("test")
	if err != nil {
		return err
	}

	return nil
}
