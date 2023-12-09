package bus

import "github.com/cgalvisleon/elvis/linq"

var SchemaBus *linq.Schema

func defineSchema() error {
	if SchemaBus != nil {
		return nil
	}

	SchemaBus = linq.NewSchema(0, "bus")

	return nil
}
