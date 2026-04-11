package crontab

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
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

func New(tag string) (*Jobs, error) {
	loc := timezone.Location()

	_, err := event.Load()
	if err != nil {
		return nil, err
	}

	result := &Jobs{
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

	return result, nil
}

/**
* addEventJob
* @param jobType TypeJob, tag, spec, channel string, started bool, params et.Json, repetitions int, fn func(event.EvenMessage)
* @return error
**/
func (s *Jobs) addEventJob(jobType TypeJob, tag, spec, channel string, started bool, params et.Json, repetitions int, fn func(event.EvenMessage)) error {
	data := et.Json{
		"type":        jobType,
		"tag":         tag,
		"spec":        spec,
		"channel":     channel,
		"started":     started,
		"params":      params,
		"repetitions": repetitions,
	}

	err := event.Publish(EVENT_CRONTAB_REMOVE, et.Json{"tag": tag})
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	err = event.Publish(EVENT_CRONTAB_SET, data)
	if err != nil {
		return err
	}

	err = event.Stack(channel, fn)
	if err != nil {
		return err
	}

	return nil
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
	result, exists := s.Jobs[tag]
	s.mu.Unlock()
	if exists {
		return result, nil
	}

	result = newJob(s, tp, tag, spec, channel, params, repetitions)
	s.mu.Lock()
	s.Jobs[tag] = result
	s.mu.Unlock()

	logs.Logf(packageName, fmt.Sprintf("job:%s | status:add | type:%s | spec:%s", tag, tp, spec))

	if started {
		result.Start()
	}

	return result, nil
}

/**
* removeJob
* @param tag string
* @return bool
**/
func (s *Jobs) removeJob(tag string) bool {
	s.mu.Lock()
	job, exists := s.Jobs[tag]
	s.mu.Unlock()
	if !exists {
		return false
	}

	job.Stop()

	s.mu.Lock()
	delete(s.Jobs, tag)
	s.mu.Unlock()

	logs.Logf(packageName, fmt.Sprintf("job:%s removed", tag))
	return true
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
* stopJob
* @param tag string
* @return error
**/
func (s *Jobs) stopJob(tag string) error {
	s.mu.Lock()
	job, exists := s.Jobs[tag]
	s.mu.Unlock()
	if !exists {
		return fmt.Errorf("job not found")
	}

	job.Stop()
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

	err := s.eventInit()
	if err != nil {
		return err
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

	for _, job := range s.Jobs {
		job.Stop()
	}

	s.running = false

	logs.Logf(packageName, `Crontab stopped`)

	return nil
}
