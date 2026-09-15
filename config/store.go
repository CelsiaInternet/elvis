package config

import (
	"fmt"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jrpc"
)

const (
	EVENT_SET_CONFIG  = "event:set:config"
	EVENT_ONCE_CONFIG = "event:once:config"
)

type ConfigStore struct {
	projectId  string
	packegName string
	stage      string
}

func NewConfigStore(projectId string, packegName string, stage string) *ConfigStore {
	result := &ConfigStore{
		projectId:  projectId,
		packegName: packegName,
		stage:      stage,
	}
	result.initEvent()
	return result
}

func (c *ConfigStore) Get(_default string, name string) string {
	data := et.Json{
		"project_id":   c.projectId,
		"stage":        c.stage,
		"package_name": c.packegName,
		"key":          name,
	}
	item, err := jrpc.CallItem("config.Services.GetConfigKey", data)
	if err != nil {
		return _default
	}

	value := item.Result.Str(name)
	if value == "" {
		return _default
	}

	return value
}

func (c *ConfigStore) SetConfig(name string, value string) {

	data := et.Json{
		"project_id":   c.projectId,
		"stage":        c.stage,
		"package_name": c.packegName,
		"description":  "",
		"config":       et.Json{name: value},
	}

	jrpc.CallItem("config.Services.SetConfig", data)
}

func (c *ConfigStore) initEvent() {
	channel := fmt.Sprintf("%s:%s", EVENT_ONCE_CONFIG, c.packegName)
	event.Stack(channel, func(message event.EvenMessage) {
		config := envar.GetConfig()
		event.Publish(EVENT_SET_CONFIG, et.Json{
			"project_id":   c.projectId,
			"stage":        c.stage,
			"package_name": c.packegName,
			"config":       config,
		})
	})
}
