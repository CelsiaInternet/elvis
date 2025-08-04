package controlplane

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/celsiainternet/elvis/cache"
)

func NewStorage() *ControlPlane {
	return &ControlPlane{
		Nodes:    make(map[int]*NodeInfo),
		MaxNodes: 0,
		mu:       sync.RWMutex{},
	}
}

/**
* load
* @return error
**/
func load(name string) (*ControlPlane, error) {
	storage := NewStorage()
	bt, err := json.Marshal(storage)
	if err != nil {
		return nil, err
	}

	strs, err := cache.Get(fmt.Sprintf("control-plane/%s", name), string(bt))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(strs), &storage)
	if err != nil {
		return nil, err
	}

	storage.mu = sync.RWMutex{}

	return storage, nil
}

/**
* save
* @return error
**/
func save(name string, cp *ControlPlane) error {
	bt, err := json.Marshal(cp)
	if err != nil {
		return err
	}

	if err := cache.Set(fmt.Sprintf("control-plane/%s", name), string(bt), 0); err != nil {
		return err
	}

	return nil
}

func Reset(name string) error {
	if _, err := cache.Delete(fmt.Sprintf("control-plane/%s", name)); err != nil {
		return err
	}

	return nil
}
