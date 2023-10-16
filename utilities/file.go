package utilities

import (
	"errors"
	"os"
	"strings"
)

func ExistPath(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func MakeFile(folder, name, model string, args ...any) (string, error) {
	path := Format(`%s/%s`, folder, name)

	if ExistPath(path) {
		return "", errors.New("File found")
	}

	_content := Params(model, args...)
	content := []byte(_content)
	err := os.WriteFile(path, content, 0666)
	if err != nil {
		return "", err
	}

	return path, nil
}

func MakeFolder(names ...string) (string, error) {
	var path string
	for _, name := range names {
		path = Append(path, name, "/")

		if !ExistPath(path) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return path, err
			}
		}
	}

	return path, nil
}

func RemoveFile(path string) (bool, error) {
	file := path
	if _, err := os.Stat(file); os.IsNotExist(err) {
		if err != nil {
			return false, err
		}

		return true, nil
	} else {
		os.Remove(file)
		return true, nil
	}
}

func ExtencionFile(filename string) string {
	lst := strings.Split(filename, ".")
	n := len(lst)
	if n > 1 {
		return lst[n-1]
	}

	return ""
}
