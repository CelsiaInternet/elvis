package create

import (
	"fmt"

	. "github.com/cgalvisleon/elvis/utilities"
)

func MakeRest(name string) error {
	path, err := MakeFolder("rest")
	if err != nil {
		return err
	}

	_, err = MakeFile(path, fmt.Sprintf("%s.http", name), restHttp, name)
	if err != nil {
		return err
	}

	return nil
}
