package jfdb

import (
	"time"

	"github.com/cgalvisleon/elvis/elvis"
)

type User struct {
	Database     *Database
	Created_date time.Time
	Update_date  time.Time
	Id           string
	Usernaname   string
	Data         elvis.Json
	Index        Number
}

func NewUser() *User {
	return &User{}
}

func (s *User) Save() {

}
