package crontab

import (
	"encoding/json"
	"slices"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
)

type Storage struct {
	Jobs    []*Job
	Version string
}

func NewStorage() *Storage {
	return &Storage{
		Jobs:    make([]*Job, 0),
		Version: "v0.0.1",
	}
}

/**
* storage
* @return error
**/
func (s *Jobs) storage() (*Storage, error) {
	storage := NewStorage()
	bt, err := json.Marshal(storage)
	if err != nil {
		return nil, err
	}

	src, err := cache.Get(s.storageKey, string(bt))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(src), &storage)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

/**
* load
* @param isServer bool
* @return error
**/
func (s *Jobs) load(isServer bool) error {
	_, err := cache.Load()
	if err != nil {
		return err
	}

	_, err = event.Load()
	if err != nil {
		return err
	}

	s.isServer = isServer
	if !isServer {
		logs.Logf(packageName, `Crontab loaded`)

		return nil
	}

	storage, err := s.storage()
	if err != nil {
		return err
	}

	for _, job := range storage.Jobs {
		_, err := s.addEventJob(job.Id, job.Name, job.Spec, job.Channel, job.Params, false)
		if err != nil {
			continue
		}
	}

	logs.Logf(packageName, `Crontab loaded`)

	return nil
}

/**
* save
* @return error
**/
func (s *Jobs) save() error {
	if !s.isServer {
		return nil
	}

	if !cache.IsLoad() {
		return nil
	}

	storage, err := s.storage()
	if err != nil {
		return err
	}

	for i, job := range s.jobs {
		if job.delete {
			idx := slices.IndexFunc(storage.Jobs, func(s *Job) bool { return s.Id == job.Id })
			if idx != -1 {
				storage.Jobs = slices.Delete(storage.Jobs, idx, idx+1)
			}
			s.jobs = slices.Delete(s.jobs, i, i+1)
			continue
		}

		idx := slices.IndexFunc(storage.Jobs, func(s *Job) bool { return s.Id == job.Id })
		if idx != -1 {
			storage.Jobs[idx] = job
		}
		storage.Jobs = append(storage.Jobs, job)
	}

	bt, err := json.Marshal(storage)
	if err != nil {
		return err
	}

	cache.Set(s.storageKey, string(bt), 0)

	return nil
}
