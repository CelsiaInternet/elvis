package et

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/cgalvisleon/elvis/timezone"
)

type Any struct {
	value any
}

func NewAny(val any) *Any {
	return &Any{
		value: val,
	}
}

func (an *Any) IsNil() bool {
	switch an.value {
	case nil:
		return true
	case "<nil>":
		return true
	case "":
		return true
	case "0":
		return true
	case "0.0":
		return true
	}

	return false
}

func (an *Any) Set(val any) any {
	an.value = val
	return an.value
}

func (an *Any) Val() any {
	return an.value
}

func (an *Any) Str() string {
	result := fmt.Sprintf(`%v`, an.value)

	return result
}

func (an *Any) Int() int {
	switch v := an.value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case string:
		if v == "" {
			return 0
		}
		r, err := strconv.Atoi(v)
		if err != nil {
			log.Println("Any value int not conver:", reflect.TypeOf(v), "value:", v)
			return 0
		}
		return r
	default:
		log.Println("Any value int not conver:", reflect.TypeOf(v), "value:", v)
		return 0
	}
}

func (an *Any) Int64() int64 {
	switch v := an.value.(type) {
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case string:
		if v == "" {
			return 0
		}
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Println("Any value int64 not conver:", reflect.TypeOf(v), "value:", v)
			return 0
		}
		return r
	default:
		log.Println("Any value int64 not conver:", reflect.TypeOf(v), "value:", v)
		return 0
	}
}

func (an *Any) Num() float64 {
	switch v := an.value.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	case float32:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		log.Println("Any value number not conver:", reflect.TypeOf(v), "value:", v)
		return 0
	}
}

func (an *Any) Bool() bool {
	switch v := an.value.(type) {
	case bool:
		return v
	default:
		log.Println("Any value boolean not conver:", reflect.TypeOf(v), "value:", v)
		return false
	}
}

func (an *Any) Time() time.Time {
	_default := timezone.NowTime()
	switch v := an.value.(type) {
	case int:
		return _default
	case string:
		layout := "2006-01-02T15:04:05.000Z"
		result, err := time.Parse(layout, v)
		if err != nil {
			return _default
		}
		return result
	case time.Time:
		return v
	default:
		log.Println("Any value time not conver:", reflect.TypeOf(v), "value:", v)
		return _default
	}
}
