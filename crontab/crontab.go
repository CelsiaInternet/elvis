package crontab

import (
	"fmt"
	"os"
	"sync"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
	"github.com/robfig/cron/v3"
)

const (
	packageName = "crontab"
)

var (
	hostName, _  = os.Hostname()
	ErrJobExists = fmt.Errorf("job already exists")
)

type Jobs struct {
	Id       string          `json:"id"`
	Tag      string          `json:"tag"`
	HostName string          `json:"host_name"`
	Jobs     map[string]*Job `json:"jobs"`
	cronJobs *cron.Cron      `json:"-"`
	running  bool            `json:"-"`
	mu       *sync.Mutex     `json:"-"`
}

func New(tag string) *Jobs {
	loc := timezone.Location()
	return &Jobs{
		Id:       utility.UUID(),
		Tag:      tag,
		HostName: hostName,
		Jobs:     make(map[string]*Job),
		cronJobs: cron.New(
			cron.WithSeconds(),
			cron.WithLocation(loc),
		),
		mu: &sync.Mutex{},
	}
}

/**
* addJob
* @param tp TypeJob, tag, spec, channel string, started bool, params et.Json, repetitions int
* @return *Job, error
**/
func (s *Jobs) addJob(tp TypeJob, tag, spec, channel string, started bool, params et.Json, repetitions int) (*Job, error) {
	if !utility.ValidStr(tag, 0, []string{"", " "}) {
		return nil, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "tag")
	}

	s.mu.Lock()
	result, ok := s.Jobs[tag]
	s.mu.Unlock()
	if ok {
		if !result.Started {
			result.Spec = spec
			result.Stop()
		}

		return result, nil
	}

	result = &Job{
		Type:        tp,
		Tag:         tag,
		Channel:     channel,
		Params:      params,
		Spec:        spec,
		Started:     false,
		Status:      Pending,
		HostName:    hostName,
		Attempts:    0,
		Repetitions: repetitions,
		idx:         -1,
		jobs:        s,
		mu:          &sync.Mutex{},
	}

	s.mu.Lock()
	s.Jobs[tag] = result
	s.mu.Unlock()
	if started {
		err := result.Start()
		if err != nil {
			return nil, err
		}
	} else {
		result.setStatus(Pending)
	}

	return result, nil
}

/**
* removeJob
* @param tag string
* @return error
**/
func (s *Jobs) removeJob(tag string) error {
	s.mu.Lock()
	job, exists := s.Jobs[tag]
	s.mu.Unlock()
	if !exists {
		return fmt.Errorf("job not found")
	}

	job.Stop()

	s.mu.Lock()
	delete(s.Jobs, tag)
	s.mu.Unlock()

	return nil
}

/**
* startJob
* @param tag string
* @return error
**/
func (s *Jobs) startJob(tag string) error {
	s.mu.Lock()
	job, exists := s.Jobs[tag]
	s.mu.Unlock()
	if !exists {
		return fmt.Errorf("job not found")
	}

	err := job.Start()
	if err != nil {
		return err
	}

	return nil
}

/**
* Start
* @return error
**/
func (s *Jobs) start() error {
	if s.cronJobs == nil {
		return fmt.Errorf("crontab not initialized")
	}

	if s.running {
		return nil
	}

	s.cronJobs.Start()
	s.running = true

	logs.Logf(packageName, `Crontab started`)

	return nil
}

/**
* stop
* @return error
**/
func (s *Jobs) stop() error {
	if s.cronJobs == nil {
		return fmt.Errorf("crontab not initialized")
	}

	s.cronJobs.Stop()
	s.running = false

	logs.Logf(packageName, `Crontab stopped`)

	return nil
}
