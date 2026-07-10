package create

import (
	"os/exec"
	"strings"
)

/**
* Command: Runs each shell command in coms and collects its stdout.
* @param coms []string
* @return [][]byte, error
**/
func Command(coms []string) ([][]byte, error) {
	var result [][]byte
	for _, com := range coms {
		args := strings.Fields(com)
		if len(args) == 0 {
			continue
		}

		out, err := exec.Command(args[0], args[1:]...).Output()
		if err != nil {
			return result, err
		}
		result = append(result, out)
	}

	return result, nil
}
