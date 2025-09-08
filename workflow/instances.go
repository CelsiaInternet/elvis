package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/logs"
)

/**
* cloneInstance
* @param id, tag string
* @return *Flow, error
**/
func (s *WorkFlows) cloneInstance(id, tag string) (*Flow, error) {
	if s.Flows[tag] == nil {
		return nil, fmt.Errorf("flow not found")
	}

	result := s.Flows[tag].cloneInstance(id)
	return result, nil
}

/**
* newInstance
* @param id, tag string
* @return *Flow, error
**/
func (s *WorkFlows) newInstance(id, tag string, startId int) (*Flow, error) {
	result, err := s.cloneInstance(id, tag)
	if err != nil {
		return nil, err
	}
	result.Current = startId
	s.Instance[result.Id] = result

	return result, nil
}

/**
* getInstance
* @param id string
* @return *Flow, error
**/
func (s *WorkFlows) getInstance(id string) (*Flow, error) {
	if s.Instance[id] != nil {
		return s.Instance[id], nil
	}

	if !cache.Exists(id) {
		return nil, nil
	}

	source := &Flow{}
	bt, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	src, err := cache.Get(id, string(bt))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(src), &source)
	if err != nil {
		return nil, err
	}

	result, err := s.cloneInstance(id, source.Tag)
	if err != nil {
		return nil, err
	}

	result.Current = source.Current
	result.TotalAttempts = source.TotalAttempts
	result.TimeAttempts = source.TimeAttempts
	result.RetentionTime = source.RetentionTime
	result.Ctxs = source.Ctxs
	result.Results = source.Results
	result.Rollbacks = source.Rollbacks
	result.LastRollback = source.LastRollback
	result.TpConsistency = source.TpConsistency
	result.CreatedAt = source.CreatedAt
	result.UpdatedAt = source.UpdatedAt
	result.DoneAt = source.DoneAt
	result.setCtx(source.Ctx)
	result.setStatus(source.Status)
	s.Instance[id] = result
	logs.Logf(packageName, "Instancia load:%s tag:%s currentStep:%d", id, result.Tag, result.Current)

	return result, nil
}
