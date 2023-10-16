package create

import (
	"log"
	"os/exec"

	"github.com/cgalvisleon/elvis/console"
)

func Command(coms []string) ([][]byte, error) {
	var result [][]byte
	for _, com := range coms {
		console.Log(com)

		out, err := exec.Command(com).Output()
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}

	console.Log(result)
	return result, nil
}
