package jfdb

import (
	"time"

	"github.com/cgalvisleon/elvis/elvis"
)

type Schema struct {
	Database     *Database
	Created_date time.Time
	Update_date  time.Time
	Id           string
	Name         string
	Description  string
	Data         elvis.Json
	Collections  []*Collection
}
