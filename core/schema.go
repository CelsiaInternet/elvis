package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
)

var makedCore bool

func defineSchemaCore() error {
	if makedCore {
		return nil
	}

	db := jdb.DB(0)
	if db == nil {
		return console.PanicM("Database not found")
	}

	var err error
	err = jdb.CreateSchema(db.Db, "core")
	if err != nil {
		return err
	}

	sql := `
	CREATE OR REPLACE FUNCTION core.create_constraint_if_not_exists(
	s_name text,
	t_name text,
	c_name text,
	constraint_sql text) 
	RETURNS void AS $$
	BEGIN
		IF NOT EXISTS(
		SELECT constraint_name 
		FROM information_schema.table_constraints 
		WHERE UPPER(table_schema)=UPPER(s_name)
		AND UPPER(table_name)=UPPER(t_name)
		AND UPPER(constraint_name)=UPPER(c_name)) THEN
		 execute constraint_sql;
		END IF;
	END;
	$$ LANGUAGE 'plpgsql';`

	_, err = jdb.QDDL(sql)
	if err != nil {
		return console.Panic(err)
	}

	makedCore = true

	console.LogK("CORE", "Init core")

	return nil
}
