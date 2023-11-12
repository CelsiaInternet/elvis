package utility

import (
	"log"
	"reflect"
	"strconv"
)

type Any struct {
	value any
}

func NewAny(val any) *Any {
	return &Any{
		value: val,
	}
}

func (an *Any) Set(val any) any {
	an.value = val
	return an.value
}

func (an *Any) Val() any {
	return an.value
}

func (an *Any) String() string {
	result := Format(`%v`, an.value)

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
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Println("Int value int not conver", reflect.TypeOf(v), v)
			return 0
		}
		return i
	default:
		log.Println("Int value is not int, type:", reflect.TypeOf(v), "value:", v)
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
		log.Println("Json value number not conver", reflect.TypeOf(v))
		return 0
	}
}

func (an *Any) Bool() bool {
	switch v := an.value.(type) {
	case bool:
		return v
	default:
		log.Println("Json value boolean not conver", reflect.TypeOf(v))
		return false
	}
}
