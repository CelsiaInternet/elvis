package jfdb

import (
	"time"

	"github.com/cgalvisleon/elvis/elvis"
)

type Index struct {
	Collection   *Collection
	Created_date time.Time
	Update_date  time.Time
	Id           string
	Name         string
	Sorted       bool
	Atrib        string
	Filename     string
	Data         elvis.Json
	Index        Number
}
