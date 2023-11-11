package create

import (
	"fmt"

	"github.com/cgalvisleon/elvis/utilities"
)

func MakeScripts(name string) error {
	path, err := utilities.MakeFolder("scripts")
	if err != nil {
		return err
	}

	_, err = utilities.MakeFile(path, fmt.Sprintf("%s.http", name), restHttp, name)
	if err != nil {
		return err
	}

	return nil
}
