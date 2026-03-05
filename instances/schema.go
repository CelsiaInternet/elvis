package instances

import (
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
)

func (s *Instance) defineSchema(db *jdb.DB, name string) error {
	if s.schema == nil {
		s.schema = linq.NewSchema(db, name)
	}

	return nil
}
