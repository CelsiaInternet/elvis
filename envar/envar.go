package envar

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cgalvisleon/elvis/logs"
)

func appendStr(s1, s2 string) string {
	if len(s2) == 0 {
		return s1
	}
	if len(s1) == 0 {
		return s2
	}

	return fmt.Sprintf(`%s_%s`, strings.ToUpper(s1), strings.ToUpper(s2))
}

func MetavarInt(meta, name string, _default int, usage, _var string) int {
	result := _default
	ok := false
	for _, arg := range os.Args[1:] {
		if ok {
			val, err := strconv.Atoi(arg)
			if err != nil {
				logs.Errorf(`-%s in %s (default %d)`, name, usage, _default)
			}
			_var = appendStr(meta, _var)
			os.Setenv(_var, arg)
			result = val
			break
		} else if arg == fmt.Sprintf(`-%s`, name) {
			ok = true
		}
	}

	return result
}

func MetavarStr(meta, name string, _default string, usage, _var string) string {
	result := _default
	ok := false
	for _, arg := range os.Args[1:] {
		if ok {
			if arg == "" {
				logs.Errorf(`-%s in %s (default %s)`, name, usage, _default)
			}
			_var = appendStr(meta, _var)
			os.Setenv(_var, arg)
			result = arg
			break
		} else if arg == fmt.Sprintf(`-%s`, name) {
			ok = true
		}
	}

	return result
}

func SetvarInt(name string, _default int, usage, _var string) int {
	return MetavarInt("", name, _default, usage, _var)
}

func SetvarStr(name string, _default string, usage, _var string) string {
	return MetavarStr("", name, _default, usage, _var)
}

func EnvarStr(_default string, args ...string) string {
	var _var string
	if len(args) > 1 {
		_var = appendStr(args[0], args[1])
	} else if len(args) > 0 {
		_var = args[0]
	}
	result := os.Getenv(_var)

	if result == "" {
		result = _default
	}

	return result
}

func EnvarInt(_default int, args ...string) int {
	_var := EnvarStr(fmt.Sprintf(`%d`, _default), args...)
	result, err := strconv.Atoi(_var)
	if err != nil {
		return 0
	}

	return result
}
