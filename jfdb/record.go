package jfdb

import (
	"time"

	"github.com/cgalvisleon/elvis/elvis"
)

type Record struct {
	Collection   *Collection
	Created_date time.Time
	Update_date  time.Time
	Id           string
	Data         elvis.Json
	Index        Number
}
