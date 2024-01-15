package jfdb

import (
	"time"

	"github.com/cgalvisleon/elvis/elvis"
)

type Number int64

type Database struct {
	Jfdb         *Jfdb
	Created_date time.Time
	Update_date  time.Time
	Id           string
	Name         string
	Description  string
	Filename     string
	Data         elvis.Json
	Schemas      []*Schema
	Index        Number
}

func (s *Database) content() []byte {
	return []byte("")
}

func (s *Database) save() error {
	s.Update_date = time.Now()

	return nil
}

func (s *Jfdb) Index() int {
	return len(s.Databases)
}
