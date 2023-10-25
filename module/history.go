package module

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/linq"
)

var Historys *Model

func DefineHistorys() error {
	if err := defineSchema(); err != nil {
		return console.PanicE(err)
	}

	if Historys != nil {
		return nil
	}

	Historys = NewModel(SchemaModule, "HISTORY", "Tabla de historicos", 1)
	Historys.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Historys.DefineColum("table_schema", "", "VARCHAR(80)", "")
	Historys.DefineColum("table_name", "", "VARCHAR(80)", "")
	Historys.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Historys.DefineColum("_data", "", "JSONB", "{}")
	Historys.DefineColum("index", "", "INTEGER", 0)
	Historys.DefineIndex([]string{
		"date_make",
		"table_schema",
		"table_name",
		"_id",
		"index",
	})
	Historys.UseRecycle = false

	return InitModel(Historys)
}
