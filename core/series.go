package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/strs"
)

var (
	existSeries bool
	series      map[string]int
)

func DefineSeries() error {
	if err := DefineSchemaCore(); err != nil {
		return console.PanicE(err)
	}

	existSeries, _ = jdb.ExistTable(0, "core", "SERIES")
	if existSeries {
		return nil
	}

	sql := `  
  -- DROP TABLE IF EXISTS core.SERIES CASCADE;

  CREATE TABLE IF NOT EXISTS core.SERIES(
		DATE_MAKE TIMESTAMP DEFAULT NOW(),
		DATE_UPDATE TIMESTAMP DEFAULT NOW(),
		SERIE VARCHAR(250) DEFAULT '',
		VALUE BIGINT DEFAULT 0,
		PRIMARY KEY(SERIE)
	);`

	_, err := jdb.QDDL(sql)
	if err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
*
**/
func EnabledSeries() bool {
	if series == nil {
		series = make(map[string]int)
		existSeries, _ := jdb.ExistTable(0, "core", "SERIES")
		return existSeries
	}

	return existSeries
}

func NextVal(tag string) int {
	sql := `SELECT nextval($1) AS SERIE;`

	item, err := jdb.DBQueryOne(0, sql, tag)
	if err != nil {
		console.Error(err)
		return 0
	}
	result := item.Int("serie")

	return result
}

/**
*
**/
func GetSerie(tag string) int {
	if !EnabledSeries() {
		var result int
		tag = strs.Replace(tag, ".", "")
		if _, ok := series[tag]; ok {
			result = NextVal(tag)
		} else {
			ok, _ := jdb.ExistSerie(0, "public", tag)
			if !ok {
				jdb.CreateSerie(0, "public", tag)
				result = NextVal(tag)
			} else {
				result = NextVal(tag)
			}
		}

		series[tag] = result
		return result
	}

	db := 0
	if MasterIdx != db {
		db = MasterIdx
	}

	sql := `
		INSERT INTO core.SERIES AS A (SERIE, VALUE)
		VALUES ($1, 1)
		ON CONFLICT(SERIE) DO UPDATE SET
		DATE_UPDATE = NOW(),
		VALUE = A.VALUE + 1
		RETURNING VALUE AS SERIE;`

	item, err := jdb.DBQueryOne(db, sql, tag)
	if err != nil {
		console.Error(err)
		return 0
	}

	return item.Int("serie")
}

func GetCode(tag, prefix string) string {
	num := GetSerie(tag)

	if len(prefix) == 0 {
		return strs.Format("%08v", num)
	} else {
		return strs.Format("%s%08v", prefix, num)
	}
}

func GetSerieLast(tag string) int {
	db := 0
	if MasterIdx != db {
		db = MasterIdx
	}

	sql := `
  SELECT VALUE AS SERIE
  FROM core.SERIES
  WHERE SERIE=$1 LIMIT 1;`
	item, err := jdb.DBQueryOne(db, sql, tag)
	if err != nil {
		console.Error(err)
		return 0
	}

	if item.Ok {
		return item.Int("serie")
	}

	return 0
}

func SetSerieValue(db int, tag string, val int) (int, error) {
	sql := `
	INSERT INTO core.SERIES(SERIE, VALUE)
	VALUES ($1, $2)
	ON CONFLICT(SERIE)
	DO UPDATE SET
	DATE_UPDATE = NOW(),
	VALUE = $2
	WHERE VALUE < $2
	RETURNING VALUE AS SERIE;`

	item, err := jdb.DBQueryOne(db, sql, tag, val)
	if err != nil {
		return 0, err
	}

	return item.Int("serie"), nil
}

func SyncSeries(masterIdx int, c chan int) error {
	var ok bool = true
	var rows int = 30
	var page int = 1
	for ok {
		ok = false

		offset := (page - 1) * rows
		sql := strs.Format(`
		SELECT A.*
		FROM core.SERIES A
		ORDER BY A.SERIE
		LIMIT %d OFFSET %d;`, rows, offset)

		items, err := jdb.Query(sql)
		if err != nil {
			return err
		}

		for _, item := range items.Result {
			tag := item.Str("serie")
			val := item.Int("value")
			_, err = SetSerieValue(masterIdx, tag, val)
			if err != nil {
				return console.Error(err)
			}

			ok = true
		}

		page++
	}

	c <- masterIdx

	return nil
}
