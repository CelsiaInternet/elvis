package create

import (
	"fmt"

	"github.com/cgalvisleon/elvis/file"
)

func MakeScripts(name string) error {
	path, err := file.MakeFolder("scripts")
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, fmt.Sprintf("%s.http", name), restHttp, name)
	if err != nil {
		return err
	}

	return nil
}
