package envar

import (
	"os"

	"github.com/cgalvisleon/elvis/generic"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/strs"
)

func MetaSet(meta, name string, _default any, usage, _var string) *generic.Any {
	var result *generic.Any = generic.New(_default)
	ok := false
	for _, arg := range os.Args[1:] {
		if ok {
			if arg == "" {
				logs.Errorf(`-%s in %s (default %s)`, name, usage, _default)
			}
			_var = strs.AppendStr(meta, _var)
			os.Setenv(_var, arg)
			result.Set(arg)
			break
		} else if arg == strs.Format(`-%s`, name) {
			ok = true
		}
	}

	return result
}

func SetvarAny(name string, _default any, usage, _var string) *generic.Any {
	result := MetaSet("", name, _default, usage, _var)
	return result
}

func SetvarStr(name string, _default string, usage, _var string) string {
	result := MetaSet("", name, _default, usage, _var)
	return result.Str()
}

func SetvarInt(name string, _default int, usage, _var string) int {
	result := MetaSet("", name, _default, usage, _var)
	return result.Int()
}

func EnvarAny(_default any, args ...string) *generic.Any {
	var _var string
	if len(args) > 1 {
		_var = strs.AppendStr(args[0], args[1])
	} else if len(args) > 0 {
		_var = args[0]
	}

	val := os.Getenv(_var)
	var result *generic.Any = generic.New(val)
	if result.IsNil() {
		result.Set(_default)
	}

	return result
}

func EnvarStr(_default string, args ...string) string {
	return EnvarAny(_default, args...).Str()
}

func EnvarInt(_default int, args ...string) int {
	return EnvarAny(_default, args...).Int()
}
