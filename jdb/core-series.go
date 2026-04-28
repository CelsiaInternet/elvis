package jdb

import (
	"errors"
	"fmt"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
)

func defineSeries(db *DB) error {
	exist, err := ExistTable(db, "core", "SERIES")
	if err != nil {
		return logs.Panice(err)
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
		return logs.Panice(err)
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
	 UPDATE core.SERIES SET
	 VALUE = VALUE + 1
	 WHERE SERIE = tag
	 RETURNING VALUE INTO result;
	 IF NOT FOUND THEN
	  INSERT INTO core.SERIES(SERIE, VALUE)
		VALUES (tag, 1)
		RETURNING VALUE INTO result;
	 END IF;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.setserie(tag VARCHAR(250), val BIGINT)
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 UPDATE core.SERIES SET
	 VALUE = val
	 WHERE SERIE = tag
	 RETURNING VALUE INTO result;
	 IF NOT FOUND THEN
	  INSERT INTO core.SERIES(SERIE, VALUE)
		VALUES (tag, val)
		RETURNING VALUE INTO result;	
	 END IF;

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
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.SERIES_AFTER_SET()
  RETURNS
    TRIGGER AS $$
	DECLARE
		TAG VARCHAR(250);
  BEGIN
	  SELECT CONCAT(TG_TABLE_SCHEMA, '.',  TG_TABLE_NAME) INTO TAG;
		PERFORM core.setserie(TAG, NEW.INDEX);

  	RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;
	`

	_, err := db.db.Exec(sql)
	if err != nil {
		return logs.Panice(err)
	}

	DefineSeries(db)

	return nil
}

/**
* NextSerie
* @param db *DB, tag string
* @return int64
**/
func NextSerie(db *DB, tag string) int64 {
	if !db.UseCore {
		return 0
	}

	sql := `SELECT core.nextserie($1) AS SERIE;`

	items, err := db.Query(sql, tag)
	if err != nil {
		logs.Error("jdb", err)
		return 0
	}

	item := items.First()
	if !item.Ok {
		return 0
	}

	return item.Int64("serie")
}

/**
* NextCode
* @param db *DB, tag string, prefix string
* @return string
**/
func NextCode(db *DB, tag, prefix string) string {
	num := NextSerie(db, tag)

	if len(prefix) == 0 {
		return strs.Format("%08v", num)
	} else {
		return strs.Format("%s%08v", prefix, num)
	}
}

/**
* SetSerie
* @param db *DB, tag string, val int
* @return int, error
**/
func SetSerie(db *DB, tag string, val int) (int, error) {
	sql := `SELECT core.setserie($1, $2);`

	_, err := db.Query(sql, tag, val)
	if err != nil {
		return 0, err
	}

	return val, nil
}

/**
* LastSerie
* @param db *DB, tag string
* @return int
**/
func LastSerie(db *DB, tag string) int {
	sql := `SELECT core.currserie($1) AS SERIE;`

	items, err := db.Query(sql, tag)
	if err != nil {
		return 0
	}

	item := items.First()
	if !item.Ok {
		return 0
	}

	result := item.Int("serie")

	return result
}

type series struct {
	db  *DB
	tag map[string]string
}

var seriesInstance *series

/**
* DefineSeries
* @param db *DB
* @return *series
**/
func DefineSeries(db *DB) *series {
	if seriesInstance == nil {
		seriesInstance = &series{
			db:  db,
			tag: make(map[string]string),
		}
	}

	return seriesInstance
}

/**
* NewSeries
* @param kind, tag, format string
**/
func (s *series) NewSeries(kind, tag, format string) {
	tg := fmt.Sprintf("%s:%s", kind, tag)
	s.tag[tg] = format
}

/**
* SetSeries
* @param kind, tag string, val int
* @return int, error
**/
func (s *series) SetSeries(kind, tag string, val int) (int, error) {
	tg := fmt.Sprintf("%s:%s", kind, tag)
	return SetSerie(s.db, tg, val)
}

/**
* GetSeries
* @param kind, tag string
* @return string
**/
func (s *series) GetSeries(kind, tag string) string {
	tg := fmt.Sprintf("%s:%s", kind, tag)
	format, ok := seriesInstance.tag[tg]
	if !ok {
		format = ""
		seriesInstance.tag[tg] = format
	}

	result := NextCode(seriesInstance.db, tg, format)
	return result
}

/**
* GetSeries
* @param kind, tag string
* @return string, error
**/
func GetSeries(kind, tag string) (string, error) {
	if seriesInstance == nil {
		return "", errors.New("Series not defined")
	}

	tg := fmt.Sprintf("%s:%s", kind, tag)
	format, ok := seriesInstance.tag[tg]
	if !ok {
		format = ""
		seriesInstance.tag[tg] = format
	}

	result := NextCode(seriesInstance.db, tg, format)
	return result, nil
}

/**
* SetSeries
* @param kind, tag string, val int
* @return int, error
**/
func SetSeries(kind, tag string, val int) (int, error) {
	if seriesInstance == nil {
		return 0, errors.New("Series not defined")
	}

	tg := fmt.Sprintf("%s:%s", kind, tag)
	return SetSerie(seriesInstance.db, tg, val)
}

/**
* GetLast
* @param kind, tag string
* @return int
**/
func (s *series) GetLast(kind, tag string) int {
	tg := fmt.Sprintf("%s:%s", kind, tag)
	return LastSerie(s.db, tg)
}
