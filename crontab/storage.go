package crontab

import (
	"encoding/json"

	"github.com/celsiainternet/elvis/cache"
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
* loadByCache
* @return error
**/
func (s *Jobs) loadByCache() error {
	storage, err := s.storage()
	if err != nil {
		return err
	}

	for _, job := range storage.Jobs {
		_, err := s.addEventJob(job.Id, job.Name, job.Spec, job.Channel, job.Started, job.Params)
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
	// if !cache.IsLoad() {
	// 	return nil
	// }

	// storage, err := s.storage()
	// if err != nil {
	// 	return err
	// }

	// for _, job := range s.jobs {
	// 	idx := slices.IndexFunc(storage.Jobs, func(e *Job) bool { return e.Id == job.Id })
	// 	if idx != -1 {
	// 		storage.Jobs[idx] = job
	// 	}
	// 	storage.Jobs = append(storage.Jobs, job)
	// }

	// bt, err := json.Marshal(storage)
	// if err != nil {
	// 	return err
	// }

	// cache.Set(s.storageKey, string(bt), 0)

	return nil
}
