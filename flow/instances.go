package flow

import (
	"encoding/json"
	"fmt"

	"github.com/celsiainternet/elvis/cache"
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
func (s *WorkFlows) newInstance(id, tag string) (*Flow, error) {
	result, err := s.cloneInstance(id, tag)
	if err != nil {
		return nil, err
	}
	s.Instance[id] = result

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
	result.Retries = source.Retries
	result.RetryDelay = source.RetryDelay
	result.RetentionTime = source.RetentionTime
	result.Ctx = source.Ctx
	result.Ctxs = source.Ctxs
	result.Results = source.Results
	result.Rollbacks = source.Rollbacks
	result.LastRollback = source.LastRollback
	result.Attempt = source.Attempt
	result.TpConsistency = source.TpConsistency
	result.CreatedAt = source.CreatedAt
	result.UpdatedAt = source.UpdatedAt
	result.DoneAt = source.DoneAt
	result.Status = source.Status
	s.Instance[id] = result

	return result, nil
}
