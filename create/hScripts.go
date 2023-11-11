package create

import (
	"fmt"

	utl "github.com/cgalvisleon/elvis/utilities"
)

func MakeScripts(name string) error {
	path, err := utl.MakeFolder("scripts")
	if err != nil {
		return err
	}

	_, err = utl.MakeFile(path, fmt.Sprintf("%s.http", name), restHttp, name)
	if err != nil {
		return err
	}

	return nil
}
