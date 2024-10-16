package jdb

import (
	"github.com/celsiainternet/elvis/console"
)

func defineRecords(db *DB) error {
	exist, err := ExistTable(db, "core", "RECORDS")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return defineRecordsFunction(db)
	}

	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;

  CREATE TABLE IF NOT EXISTS core.RECORDS(
    DATE_MAKE TIMESTAMP DEFAULT NOW(),
		DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
		OPTION VARCHAR(80) DEFAULT '',
    _IDT VARCHAR(80) DEFAULT '-1',
    INDEX SERIAL,
    PRIMARY KEY (TABLE_SCHEMA, TABLE_NAME, _IDT)
  );    
  CREATE INDEX IF NOT EXISTS RECORDS_TABLE_SCHEMA_IDX ON core.RECORDS(TABLE_SCHEMA);
  CREATE INDEX IF NOT EXISTS RECORDS_TABLE_NAME_IDX ON core.RECORDS(TABLE_NAME);
	CREATE INDEX IF NOT EXISTS RECORDS_OPTION_IDX ON core.RECORDS(OPTION);
  CREATE INDEX IF NOT EXISTS RECORDS__IDT_IDX ON core.RECORDS(_IDT);  
	CREATE INDEX IF NOT EXISTS RECORDS_INDEX_IDX ON core.RECORDS(INDEX);`

	id := "define-records"
	_, err = db.Command(CommandDefine, id, sql)
	if err != nil {
		return console.Panic(err)
	}

	return defineRecordsFunction(db)
}

func defineRecordsFunction(db *DB) error {
	sql := `
	CREATE OR REPLACE FUNCTION core.RECORDS_BEFORE_INSERT()
  RETURNS
    TRIGGER AS $$  
  BEGIN
    IF NEW._IDT = '-1' THEN
      NEW._IDT = uuid_generate_v4();
		END IF;

		PERFORM pg_notify(
		'before',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,        
			'_idt', NEW._IDT
		)::text
		);
  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.RECORDS_BEFORE_UPDATE()
  RETURNS
    TRIGGER AS $$  
  BEGIN    
		PERFORM pg_notify(
		'before',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'_idt', NEW._IDT
		)::text
		);
  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION core.RECORDS_BEFORE_DELETE()
  RETURNS
    TRIGGER AS $$  
  BEGIN
		PERFORM pg_notify(
		'before',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'_idt', OLD._IDT
		)::text
		);
  RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.RECORDS_AFTER_INSERT()
  RETURNS
    TRIGGER AS $$  
  BEGIN
		INSERT INTO core.RECORDS(TABLE_SCHEMA, TABLE_NAME, OPTION, _IDT)
		VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, TG_OP, NEW._IDT);

		PERFORM pg_notify(
		'after',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,        
			'_idt', NEW._IDT
		)::text
		);
  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.RECORDS_AFTER_UPDATE()
  RETURNS
    TRIGGER AS $$  
  BEGIN
    UPDATE core.RECORDS SET
		DATE_UPDATE=NOW(),
		OPTION=TG_OP
		WHERE _IDT=NEW._IDT;

		PERFORM pg_notify(
		'after',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'_idt', NEW._IDT
		)::text
		);
  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION core.RECORDS_AFTER_DELETE()
  RETURNS
    TRIGGER AS $$  
  BEGIN
    DELETE FROM core.RECORDS
    WHERE TABLE_SCHEMA = TG_TABLE_SCHEMA
    AND TABLE_NAME = TG_TABLE_NAME
    AND _IDT = OLD._IDT;

		DELETE FROM core.COMMANDS
    WHERE OPTION = 'UPDATE'
    AND _ID = OLD._IDT;

		PERFORM pg_notify(
		'after',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'_idt', OLD._IDT
		)::text
		);
  RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;`

	id := "define-records-function"
	_, err := db.Command(CommandDefine, id, sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}
