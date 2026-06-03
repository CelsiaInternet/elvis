package dt

import (
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
)

type Object struct {
	et.Item
	Key      string        `json:"key"`
	duration time.Duration `json:"-"`
}

/**
* newObject
* @param key string
* @return *Object
**/
func newObject(key string) *Object {
	duration := time.Duration(envar.GetInt(5, "CACHE_DURATION"))
	return &Object{
		Item: et.Item{
			Ok:     false,
			Result: et.Json{},
		},
		Key:      key,
		duration: duration * time.Minute,
	}
}

/**
* up
* @param data et.Json
* @return bool
**/
func (s *Object) up(data et.Item, save bool) {
	s.Ok = data.Ok
	s.Result = et.Json{}
	if !s.Ok {
		Drop(s.Key)
		return
	}

	for key, val := range data.Result {
		s.Set(key, val)
	}

	if save {
		s.save()
	}
}

/**
* save
* @return bool
**/
func (s *Object) save() bool {
	production := envar.GetBool(true, "PRODUCTION")
	if !production {
		return false
	}

	val := s.ToString()
	cache.Set("object:"+s.Key, val, s.duration)

	return true
}

/**
* Load
* @return error
*
 */
func (s *Object) load() error {
	item, err := cache.GetItem("object:" + s.Key)
	if err != nil {
		return err
	}

	s.up(item, false)

	return nil
}
