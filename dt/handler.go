package dt

import (
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
)

/**
* Up
* @param key string, data et.Item
* @return Object
**/
func Up(key string, data et.Item) Object {
	obj := newObject(key)
	obj.up(data, data.Ok)

	return *obj
}

/**
* UpWithDuration
* @param key string, data et.Item, duration time.Duration
* @return Object
**/
func UpWithDuration(key string, data et.Item, duration time.Duration) Object {
	obj := newObject(key)
	obj.duration = duration
	obj.up(data, data.Ok)

	return *obj
}

/**
* Get
* @param key string
* @return Object
**/
func Get(key string) Object {
	obj := newObject(key)
	obj.load()

	return *obj
}

/**
* Drop
* @param key string
**/
func Drop(key string) {
	cache.Delete("object:" + key)
}
