package jdb

import "github.com/celsiainternet/elvis/logs"

func defineModel(db *DB) error {
	exist, err := ExistTable(db, "core", "RECORDS")
	if err != nil {
		return logs.Panice(err)
	}

	if exist {
		return defineModelFunction(db)
	}

	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;

  CREATE TABLE IF NOT EXISTS core.MODELS(
    DATE_MAKE TIMESTAMP DEFAULT NOW(),
		DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
		_DATA JSONB DEFAULT '{}',
    _IDT VARCHAR(80) DEFAULT '-1',
    INDEX SERIAL,
    PRIMARY KEY (TABLE_SCHEMA, TABLE_NAME)
  );    
  CREATE INDEX IF NOT EXISTS MODELS_TABLE_SCHEMA_IDX ON core.MODELS(TABLE_SCHEMA);
  CREATE INDEX IF NOT EXISTS MODELS_TABLE_NAME_IDX ON core.MODELS(TABLE_NAME);
	CREATE INDEX IF NOT EXISTS MODELS__DATA_IDX ON core.MODELS USING GIN(_DATA);
  CREATE INDEX IF NOT EXISTS MODELS__IDT_IDX ON core.MODELS(_IDT);  
	CREATE INDEX IF NOT EXISTS MODELS_INDEX_IDX ON core.MODELS(INDEX);`

	_, err = db.db.Exec(sql)
	if err != nil {
		return logs.Panice(err)
	}

	return defineModelFunction(db)
}

func defineModelFunction(db *DB) error {
	sql := ``
	_, err := db.db.Exec(sql)
	if err != nil {
		return logs.Panice(err)
	}

	return nil
}
