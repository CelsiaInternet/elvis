package workflow

import (
	"errors"
	"fmt"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/reg"
)

var (
	errorInstanceNotFound = errors.New(MSG_INSTANCE_NOT_FOUND)
	errorInstanceExists   = errors.New(MSG_INSTANCE_EXISTS)
)

/**
* existInstance
* @param id, tag string
* @return bool, error
**/
func (s *WorkFlows) existInstance(id, tag string) (bool, error) {
	if s.Flows[tag] == nil {
		return false, fmt.Errorf(MSG_FLOW_NOT_FOUND)
	}

	if s.Instance[id] != nil {
		return true, nil
	}

	return s.Flows[tag].existInstance(id), nil
}

/**
* createInstance
* @param id, tag string, startId int, tags et.Json
* @return *Flow, error
**/
func (s *WorkFlows) createInstance(id, tag string, startId int, tags et.Json) (*Flow, error) {
	if exist, err := s.existInstance(id, tag); err != nil {
		return nil, err
	} else if exist {
		return nil, errorInstanceExists
	}

	result := s.Flows[tag].newInstance(id, tags)
	result.setCurrent(startId)
	s.Instance[result.Id] = result

	return result, nil
}

/**
* GetInstance
* @param id string
* @return *Flow, error
**/
func (s *WorkFlows) getInstance(id, tag string) (*Flow, error) {
	id = reg.GetUUID(id)
	if exist, err := s.existInstance(id, tag); err != nil {
		return nil, err
	} else if !exist {
		return nil, errorInstanceNotFound
	}

	result, err := s.Flows[tag].loadInstance(id)
	if err != nil {
		return nil, err
	}

	s.Instance[id] = result
	logs.Logf(packageName, MSG_INSTANCE_LOAD, id, result.Tag, result.Current)

	return result, nil
}

/**
* getOrCreateInstance
* @param id, tag string, tags et.Json
* @return *Flow, error
**/
func (s *WorkFlows) getOrCreateInstance(id, tag string, tags et.Json) (*Flow, error) {
	id = reg.GetUUID(id)
	if exist, err := s.existInstance(id, tag); err != nil {
		return nil, err
	} else if exist {
		return s.getInstance(id, tag)
	}

	return s.createInstance(id, tag, 0, tags)
}
