package jdb

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/strs"
)

func defineSeries(db *DB) error {
	exist, err := ExistTable(db, "core", "SERIES")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return defineSeriesFunction(db)
	}

	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;

  CREATE TABLE IF NOT EXISTS core.SERIES(
		SERIE VARCHAR(250) DEFAULT '',
		VALUE BIGINT DEFAULT 0,
		PRIMARY KEY(SERIE)
	);`

	_, err = db.db.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return defineSeriesFunction(db)
}

func defineSeriesFunction(db *DB) error {
	sql := `
	CREATE OR REPLACE FUNCTION core.nextserie(tag VARCHAR(250))
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 INSERT INTO core.SERIES AS A (SERIE, VALUE)
	 SELECT tag, 1
	 ON CONFLICT (SERIE) DO UPDATE SET
	 VALUE = A.VALUE + 1
	 RETURNING VALUE INTO result;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.setserie(tag VARCHAR(250), val BIGINT)
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 INSERT INTO core.SERIES AS A (SERIE, VALUE)
	 SELECT tag, val
	 ON CONFLICT (SERIE) DO UPDATE SET
	 VALUE = val
	 WHERE A.VALUE < val
	 RETURNING VALUE INTO result;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.currserie(tag VARCHAR(250))
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 SELECT VALUE INTO result
	 FROM core.SERIES
	 WHERE SERIE = tag LIMIT 1;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;`

	_, err := db.db.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}

func NextSerie(db *DB, tag string) int64 {
	if !db.UseCore {
		return 0
	}

	sql := `SELECT core.nextserie($1) AS SERIE;`

	if db.dm != nil {
		rows, err := db.dm.Query(sql, tag)
		if err != nil {
			console.Error(err)
			return 0
		}
		defer rows.Close()

		item := rowsItem(rows)
		if !item.Ok {
			return 0
		}

		result := item.Int64("serie")

		return result
	}

	item, err := db.QueryOne(sql, tag)
	if err != nil {
		console.Error(err)
		return 0
	}
	if !item.Ok {
		return 0
	}

	result := item.Int64("serie")

	return result
}

func NextCode(db *DB, tag, prefix string) string {
	num := NextSerie(db, tag)

	if len(prefix) == 0 {
		return strs.Format("%08v", num)
	} else {
		return strs.Format("%s%08v", prefix, num)
	}
}

func SetSerie(db *DB, tag string, val int) (int, error) {
	sql := `SELECT core.setserie($1, $2);`

	_, err := db.QueryOne(sql, tag, val)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func LastSerie(db *DB, tag string) int {
	sql := `SELECT core.currserie($1) AS SERIE;`

	item, err := db.QueryOne(sql, tag)
	if err != nil {
		return 0
	}

	result := item.Int("serie")

	return result
}
