package crontab

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/robfig/cron/v3"
)

type JobStatus string

const (
	Pending  JobStatus = "pending"
	Prepared JobStatus = "prepared"
	Running  JobStatus = "running"
	Done     JobStatus = "done"
	Failed   JobStatus = "failed"
	Finished JobStatus = "finished"
)

type TypeJob string

const (
	CronJob     TypeJob = "cronJob"
	ScheduleJob TypeJob = "scheduleJob"
)

type Job struct {
	Type        TypeJob       `json:"type"`
	Tag         string        `json:"tag"`
	Channel     string        `json:"channel"`
	Params      et.Json       `json:"params"`
	Spec        string        `json:"spec"`
	Started     bool          `json:"started"`
	Status      JobStatus     `json:"status"`
	HostName    string        `json:"host_name"`
	Attempts    int           `json:"attempts"`
	Repetitions int           `json:"repetitions"`
	Duration    time.Duration `json:"duration"`
	idx         cron.EntryID  `json:"-"`
	shot        *time.Timer   `json:"-"`
	jobs        *Jobs         `json:"-"`
	mu          *sync.Mutex   `json:"-"`
}

/**
* Serialize
* @return ([]byte, error)
**/
func (s *Job) serialize() ([]byte, error) {
	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Job) ToJson() et.Json {
	bt, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* save
* @return error
**/
func (s *Job) Save() error {
	if saveInstance == nil {
		return nil
	}

	return saveInstance(s.Tag, s.Tag, s)
}

/**
* setStatus
* @param status JobStatus
* @return void
**/
func (s *Job) setStatus(status JobStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Status = status
	logs.Logf(packageName, fmt.Sprintf("Job %s status:%s host:%s attempt:%d", s.Tag, s.Status, s.HostName, s.Attempts))
	go s.Save()
}

/**
* Start
* @return error
**/
func (s *Job) Start() error {
	fn := func() {
		s.Attempts++
		err := event.Publish(s.Channel, s.Params)
		if err != nil {
			s.setStatus(Failed)
		}
		s.setStatus(Running)
		if s.Repetitions != 0 && s.Attempts >= s.Repetitions {
			s.Finish()
		} else {
			s.setStatus(Pending)
		}
	}

	if s.Type == CronJob {
		id, err := s.jobs.cronJobs.AddFunc(s.Spec, fn)
		if err != nil {
			return err
		}

		s.idx = id
	} else {
		now := timezone.NowTime()
		shotTime, err := timezone.Parse("2006-01-02T15:04:05", s.Spec)
		if err != nil {
			return err
		}
		if shotTime.After(now) {
			duration := shotTime.Sub(now)
			s.Duration = duration
			s.shot = time.AfterFunc(duration, fn)
			return s.Save()
		} else if s.shot != nil {
			s.Stop()
		}
	}

	s.Started = true

	return nil
}

/**
* Stop
* @return error
**/
func (s *Job) Stop() {
	if !s.Started {
		return
	}

	s.Started = false
	time.AfterFunc(1*time.Second, func() {
		if s.Type == CronJob {
			s.jobs.cronJobs.Remove(s.idx)
			s.idx = -1
		} else if s.shot != nil {
			s.shot.Stop()
		}
		s.setStatus(Pending)
	})
}

/**
* Finish
* @return error
**/
func (s *Job) Finish() {
	s.Stop()
	time.AfterFunc(1*time.Second, func() {
		s.setStatus(Finished)
	})
}
