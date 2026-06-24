package envar

import (
	"fmt"
	"os"
	"strconv"

	"github.com/celsiainternet/elvis/strs"
	_ "github.com/joho/godotenv/autoload"
)

type Store interface {
	Get(name string, _default string) string
}

var store Store

func Load(s Store) {
	store = s
}

/**
* MetaSet
* @param name string, _default string, description string, _var string
* @return string
**/
func MetaSet(name string, _default string, description, _var string) string {
	for i, arg := range os.Args[1:] {
		if arg == strs.Format("-%s", name) {
			val := os.Args[i+2]
			os.Setenv(_var, val)
			return val
		}
	}

	return _default
}

/**
* SetStr
* @param name string, _default string, usage string, _var string
* @return string
**/
func SetStr(name string, _default string, usage, _var string) string {
	return MetaSet(name, _default, usage, _var)
}

/**
* SetInt
* @param name string, _default int, usage string, _var string
* @return string
**/
func SetInt(name string, _default int, usage, _var string) int {
	result := MetaSet(name, fmt.Sprintf("%d", _default), usage, _var)
	val, err := strconv.Atoi(result)
	if err != nil {
		return _default
	}
	return val
}

/**
* SetInt64
* @param name string, _default int64, usage string, _var string
* @return int64
**/
func SetInt64(name string, _default int64, usage, _var string) int64 {
	result := MetaSet(name, strconv.FormatInt(_default, 10), usage, _var)
	val, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return _default
	}
	return val
}

/**
* SetBool
* @param name string, _default bool, usage string, _var string
**/
func SetBool(name string, _default bool, usage, _var string) bool {
	result := MetaSet(name, strconv.FormatBool(_default), usage, _var)

	val, err := strconv.ParseBool(result)
	if err != nil {
		return _default
	}

	return val
}

/**
* UpSetStr
* @param name string, value string
* @return string
**/
func UpSetStr(name string, value string) string {
	os.Setenv(name, value)
	return value
}

/**
* SetInt
* @param name string, value int
* @return int
**/
func UpSetInt(name string, value int) int {
	os.Setenv(name, strconv.Itoa(value))
	return value
}

/**
* UpSetFloat
* @param name string, value float64
* @return float64
**/
func UpSetFloat(name string, value float64) float64 {
	os.Setenv(name, strconv.FormatFloat(float64(value), 'f', -1, 64))
	return value
}

/**
* UpSetBool
* @param name string, value bool
* @return bool
**/
func UpSetBool(name string, value bool) bool {
	os.Setenv(name, strconv.FormatBool(value))
	return value
}

/**
* GetStr
* @param _default string, _var string
* @return string
**/
func GetStr(_default string, _var string) string {
	if store != nil {
		return store.Get(_default, _var)
	}

	result := os.Getenv(_var)

	if result == "" {
		return _default
	}

	return result
}

/**
* GetInt
* @param _default int, _var string
* @return int
**/
func GetInt(_default int, _var string) int {
	result := GetStr(strconv.Itoa(_default), _var)

	val, err := strconv.Atoi(result)
	if err != nil {
		return _default
	}

	return val
}

/**
* GetInt64
* @param int64 _default, string _var
* @return int64
**/
func GetInt64(_default int64, _var string) int64 {
	result := GetStr(strconv.FormatInt(_default, 10), _var)

	val, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return _default
	}

	return val
}

/**
* GetFloat64
* @param int64 _default, string _var
* @return int64
**/
func GetFloat64(_default float64, _var string) float64 {
	v := GetStr("0,0", _var)
	result, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return _default
	}

	return result
}

/**
* GetBool
* @param _default bool, _var string
* @return bool
**/
func GetBool(_default bool, _var string) bool {
	result := GetStr(strconv.FormatBool(_default), _var)

	val, err := strconv.ParseBool(result)
	if err != nil {
		return _default
	}

	return val
}
