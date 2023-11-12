package create

import (
	"fmt"

	"github.com/cgalvisleon/elvis/utility"
)

func MakeScripts(name string) error {
	path, err := utility.MakeFolder("scripts")
	if err != nil {
		return err
	}

	_, err = utility.MakeFile(path, fmt.Sprintf("%s.http", name), restHttp, name)
	if err != nil {
		return err
	}

	return nil
}
