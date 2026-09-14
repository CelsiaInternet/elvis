package config

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jrpc"
)

type ConfigStore struct {
	projectId  string
	packegName string
	stage      string
}

func NewConfigStore(projectId string, packegName string, stage string) *ConfigStore {
	return &ConfigStore{
		projectId:  projectId,
		packegName: packegName,
		stage:      stage,
	}
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




