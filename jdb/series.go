package jdb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/strs"
)

var (
	makedSeries bool
)

func defineSeries(db *sql.DB) error {
	if makedSeries {
		return nil
	}

	sql := `
  CREATE TABLE IF NOT EXISTS core.SERIES(
		SERIE VARCHAR(250) DEFAULT '',
		VALUE BIGINT DEFAULT 0,
		PRIMARY KEY(SERIE)
	);
	
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

	_, err := Query(db, sql)
	if err != nil {
		return console.Panic(err)
	}

	makedSeries = true

	return nil
}

func NextSerie(tag string) int {
	if !makedSeries {
		return 0
	}

	sql := `SELECT core.nextserie($1) AS SERIE;`

	item, err := QueryOne(conn.Db, sql, tag)
	if err != nil {
		console.Error(err)
		return 0
	}

	result := item.Int("serie")

	return result
}

func NextCode(tag, prefix string) string {
	num := NextSerie(tag)

	if len(prefix) == 0 {
		return strs.Format("%08v", num)
	} else {
		return strs.Format("%s%08v", prefix, num)
	}
}

func SetSerie(tag string, val int) (int, error) {
	if !makedSeries {
		return 0, nil
	}

	sql := `SELECT core.setserie($1, $2);`

	_, err := QueryOne(conn.Db, sql, tag, val)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func LastSerie(tag string) int {
	if !makedSeries {
		return 0
	}

	sql := `SELECT core.currserie($1) AS SERIE;`

	item, err := QueryOne(conn.Db, sql, tag)
	if err != nil {
		return 0
	}

	result := item.Int("serie")

	return result
}

func UUIndex(db *sql.DB, tag string) (int64, error) {
	now := time.Now()
	result := now.UnixMilli() * 10000
	replica, err := getVarInt(db, "REPLICA", 1)
	if err != nil {
		return 0, err
	}

	if replica < 10 {
		replica = replica * 1000
	} else if replica < 100 {
		replica = replica * 100
	} else {
		replica = replica * 10
	}

	result = result + replica
	key := fmt.Sprintf("%s:%d", tag, result)
	count := cache.Count(key, 1)

	result = result + count
	return result, nil
}
