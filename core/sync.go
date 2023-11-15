package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/utility"
)

var existSyncs bool

func DefineSync() error {
	existSyncs, _ := ExistTable(0, "core", "SYNCS")
	if existSyncs {
		return nil
	}

	if err := DefineCoreSchema(); err != nil {
		return console.PanicE(err)
	}

	sql := `
  -- DROP SCHEMA IF EXISTS core CASCADE;
  -- DROP TABLE IF EXISTS core.SYNCS CASCADE;

  CREATE TABLE IF NOT EXISTS core.SYNCS(
    DATE_MAKE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    ACTION VARCHAR(80) DEFAULT '',    
    _IDT VARCHAR(80) DEFAULT '-1',
    _DATA JSONB DEFAULT '{}',
    INDEX SERIAL,
    PRIMARY KEY (_IDT)
  );  
  CREATE INDEX IF NOT EXISTS SYNCS_INDEX_IDX ON core.SYNCS(INDEX);

  CREATE OR REPLACE FUNCTION core.SYNC_INSERT()
  RETURNS
    TRIGGER AS $$
  DECLARE
    SYNC BOOLEAN;
  BEGIN
    SYNC = NEW._IDT = '-1';
    IF SYNC THEN NEW._IDT = uuid_generate_v4(); END IF;    

    IF SYNC THEN
      INSERT INTO core.SYNCS(TABLE_SCHEMA, TABLE_NAME, ACTION, _DATA, _IDT)
      VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, TG_OP, row_to_json(NEW),  NEW._IDT);

      PERFORM pg_notify(
      'sync',
      json_build_object(
        '_idt', NEW._IDT
      )::text
      );
    END IF;

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION core.SYNC_UPDATE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    SYNC BOOLEAN;
  BEGIN
    SYNC = NEW._IDT != '-1';

    IF SYNC THEN
      INSERT INTO core.SYNCS(TABLE_SCHEMA, TABLE_NAME, ACTION, _DATA, _IDT)
      VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, TG_OP, row_to_json(NEW),  NEW._IDT)
      ON CONFLICT (_IDT) DO UPDATE SET
      ACTION = TG_OP,
      _DATA = row_to_json(NEW),
      INDEX = 0;

      PERFORM pg_notify(
      'sync',
      json_build_object(
        '_idt', NEW._IDT
      )::text
      );
    ELSE
      NEW._IDT = OLD._IDT;
      DELETE FROM core.SYNCS
      WHERE _IDT = OLD._IDT;
    END IF; 

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION core.SYNC_DELETE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    OK VARCHAR(80);
  BEGIN
    SELECT _IDT INTO OK FROM core.SYNCS WHERE _IDT = OLD._IDT LIMIT 1;
    IF FOUND THEN
      INSERT INTO core.SYNCS(TABLE_SCHEMA, TABLE_NAME, ACTION, _DATA, _IDT)
      VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, TG_OP, row_to_json(OLD),  OLD._IDT)
      ON CONFLICT (_IDT) DO UPDATE SET
      ACTION = TG_OP,
      _DATA = row_to_json(OLD),
      INDEX = 0;
      
      PERFORM pg_notify(
      'sync',
      json_build_object(
        '_idt', OLD._IDT
      )::text
      );
    END IF;

  RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;
  `

	_, err := jdb.QDDL(sql)
	if err != nil {
		return console.PanicE(err)
	}

	return nil
}

func SetSyncTrigger(schema, table string) error {
	exist, _ := ExistTable(0, "core", "SYNCS")
	if !exist {
		return nil
	}

	created, err := CreateColumn(0, schema, table, "_IDT", "VARCHAR(80)", "-1")
	if err != nil {
		return err
	}

	if created {
		tableName := utility.Append(utility.Lowcase(schema), utility.Uppcase(table), ".")
		sql := jdb.SQLDDL(`
    CREATE INDEX IF NOT EXISTS $2_IDT_IDX ON $1(_IDT);

    DROP TRIGGER IF EXISTS SYNC_INSERT ON $1 CASCADE;
    CREATE TRIGGER SYNC_INSERT
    BEFORE INSERT ON $1
    FOR EACH ROW
    EXECUTE PROCEDURE core.SYNC_INSERT();

    DROP TRIGGER IF EXISTS SYNC_UPDATE ON $1 CASCADE;
    CREATE TRIGGER SYNC_UPDATE
    BEFORE UPDATE ON $1
    FOR EACH ROW
    EXECUTE PROCEDURE core.SYNC_UPDATE();

    DROP TRIGGER IF EXISTS SYNC_DELETE ON $1 CASCADE;
    CREATE TRIGGER SYNC_DELETE
    BEFORE DELETE ON $1
    FOR EACH ROW
    EXECUTE PROCEDURE core.SYNC_DELETE();`, tableName, utility.Uppcase(table))

		_, err := jdb.QDDL(sql)
		if err != nil {
			return err
		}
	}

	return nil
}
