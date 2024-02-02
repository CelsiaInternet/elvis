package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/strs"
)

func DefineSync() error {
	if err := DefineSchemaCore(); err != nil {
		return console.Panic(err)
	}

	existSyncs, _ := jdb.ExistTable(0, "core", "SYNCS")
	if existSyncs {
		return nil
	}

	sql := `
  -- DROP SCHEMA IF EXISTS core CASCADE;
  -- DROP TABLE IF EXISTS core.SYNCS CASCADE;

  CREATE TABLE IF NOT EXISTS core.SYNCS(
    DATE_MAKE TIMESTAMP DEFAULT NOW(),
    DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    _IDT VARCHAR(80) DEFAULT '-1',
    ACTION VARCHAR(80) DEFAULT '',
    _ID VARCHAR(80) DEFAULT '-1',
    _SYNC BOOLEAN DEFAULT FALSE,    
    INDEX SERIAL,
    PRIMARY KEY (TABLE_SCHEMA, TABLE_NAME, _IDT)
  );  
  CREATE INDEX IF NOT EXISTS SYNCS_INDEX_IDX ON core.SYNCS(INDEX);

  CREATE OR REPLACE FUNCTION core.SYNC_INSERT()
  RETURNS
    TRIGGER AS $$
  BEGIN
    IF NEW._IDT = '-1' THEN
      NEW._IDT = uuid_generate_v4();

      INSERT INTO core.SYNCS(TABLE_SCHEMA, TABLE_NAME, _IDT, ACTION, _ID)
      VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, NEW._IDT, TG_OP, uuid_generate_v4());

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
  BEGIN
    IF NEW._IDT = '-1' THEN
      NEW._IDT = OLD._IDT;
    ELSE
     INSERT INTO core.SYNCS(TABLE_SCHEMA, TABLE_NAME, _IDT, ACTION)
     VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, NEW._IDT, TG_OP)
		 ON CONFLICT(TABLE_SCHEMA, TABLE_NAME, _IDT) DO UPDATE SET
     DATE_UPDATE = NOW(),
     ACTION = TG_OP,
     _SYNC = FALSE,
     _ID = uuid_generate_v4();

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

  CREATE OR REPLACE FUNCTION core.SYNC_DELETE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    VINDEX INTEGER;
  BEGIN
    SELECT INDEX INTO VINDEX
    FROM core.SYNCS
    WHERE TABLE_SCHEMA = TG_TABLE_SCHEMA
    AND TABLE_NAME = TG_TABLE_NAME
    AND _IDT = OLD._IDT
    LIMIT 1;
    IF FOUND THEN
      UPDATE core.SYNCS SET
      DATE_UPDATE = NOW(),
      ACTION = TG_OP,
      _SYNC = FALSE,
      _ID = uuid_generate_v4()
      WHERE INDEX = VINDEX;
      
      PERFORM pg_notify(
      'sync',
      json_build_object(
        '_idt', OLD._IDT
      )::text
      );
    END IF;

  RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;`

	_, err := jdb.QDDL(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}

func SetSyncTrigger(schema, table string) error {
	exist, _ := jdb.ExistTable(0, "core", "SYNCS")
	if !exist {
		return nil
	}

	created, err := jdb.CreateColumn(0, schema, table, "_IDT", "VARCHAR(80)", "-1")
	if err != nil {
		return err
	}

	if created {
		tableName := strs.Append(strs.Lowcase(schema), strs.Uppcase(table), ".")
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
    EXECUTE PROCEDURE core.SYNC_DELETE();`, tableName, strs.Uppcase(table))

		_, err := jdb.QDDL(sql)
		if err != nil {
			return err
		}
	}

	return nil
}
