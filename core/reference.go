package core

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/linq"
)

func DefineReference() error {
	if err := DefineCoreSchema(); err != nil {
		return console.PanicE(err)
	}

	exist, _ := ExistTable(0, "core", "REFERENCES")
	if exist {
		return nil
	}

	sql := `  
  -- DROP SCHEMA IF EXISTS core CASCADE;
  -- DROP TABLE IF EXISTS core.REFERENCES CASCADE;

  CREATE TABLE IF NOT EXISTS core.REFERENCES(
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
		_ID VARCHAR(80) DEFAULT '-1',
		COUNT INT DEFAULT 0,    
    INDEX SERIAL,
		PRIMARY KEY (TABLE_SCHEMA, TABLE_NAME, _ID)
  );
  CREATE INDEX IF NOT EXISTS REFERENCES_INDEX_IDX ON core.REFERENCES(INDEX);`

	_, err := QDDL(sql)
	if err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* After reference
**/
func SetReferences(references []*ReferenceValue) {
	for _, ref := range references {
		if ref.Key == "" {
			continue
		}

		sql := `
		INSERT INTO core.REFERENCES AS A (TABLE_SCHEMA, TABLE_NAME, _ID, COUNT)
		VALUES($1, $2, $3, 1)
		ON CONFLICT (TABLE_SCHEMA, TABLE_NAME, _ID) DO UPDATE SET
		COUNT = A.COUNT + $4
		RETURNING INDEX;`

		_, err := QueryOne(sql, ref.Schema, ref.Table, ref.Key, ref.Op)
		if err != nil {
			return
		}
	}
}
