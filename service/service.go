package service

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/utility"
)

const STEP_REUSE = "reuse"

func New(ownerId, tag string) string {
	if !utility.ValidKey(ownerId) {
		return ""
	}

	result := utility.UUID()
	event.Publish("service/new", et.Json{
		"created_at": utility.Now(),
		"owner_id":   ownerId,
		"tag":        tag,
		"_id":        result,
	})

	return result
}

func Get(ownerId, tag string) string {
	if !utility.ValidKey(ownerId) {
		return ""
	}

	key := cache.GenKey("service", ownerId, tag)
	if !cache.Exists(key) {
		result := New(ownerId, tag)
		cache.SetH(key, result)
		return result
	}

	result, err := cache.Get(key, "-1")
	if err != nil {
		result := New(ownerId, tag)
		cache.SetH(key, result)
		return result
	}

	event.Publish("service/reused", et.Json{
		"created_at": utility.Now(),
		"_id":        result,
	})

	return result
}

func Step(_id, step string) string {
	if !utility.ValidKey(_id) {
		return ""
	}

	event.Publish("service/step", et.Json{
		"created_at": utility.Now(),
		"_id":        _id,
		"step":       step,
	})

	return _id
}

func End(_id string) string {
	if !utility.ValidKey(_id) {
		return ""
	}

	event.Publish("service/end", et.Json{
		"created_at": utility.Now(),
		"_id":        _id,
	})

	return _id
}
