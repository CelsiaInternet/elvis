package core

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/utilities"
)

func DefineSeries() error {
	if err := DefineCoreSchema(); err != nil {
		return console.PanicE(err)
	}

	exist, _ := ExistTable(0, "core", "SERIES")
	if exist {
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
	);

	CREATE OR REPLACE FUNCTION core.GET_SERIE(VSERIE VARCHAR(250))
	RETURNS
		BIGINT AS $$
	DECLARE
		RESULT BIGINT;
	BEGIN
		INSERT INTO core.SERIES AS A (SERIE, VALUE)
		VALUES (VSERIE, 1)
		ON CONFLICT(SERIE)
		DO UPDATE SET
		DATE_UPDATE=NOW(),
		VALUE=A.VALUE + 1
		RETURNING VALUE INTO RESULT;

		RETURN RESULT;
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.SET_SERIE(VSERIE VARCHAR(250), VAL BIGINT)
	RETURNS
		BIGINT AS $$
	DECLARE
		RESULT BIGINT;
	BEGIN
		INSERT INTO core.SERIES(SERIE, VALUE)
		VALUES (VSERIE, VAL)
		ON CONFLICT(SERIE)
		DO UPDATE SET
		DATE_UPDATE=NOW(),
		VALUE=VAL
		WHERE VALUE<VAL
		RETURNING VALUE INTO RESULT;
		
		RETURN RESULT;
	END;
	$$ LANGUAGE plpgsql;
  `

	_, err := QDDL(sql)
	if err != nil {
		return console.PanicE(err)
	}

	return nil
}

// Serires
func GetSerie(tag string) int {
	db := 0
	if MasterIdx != db {
		db = MasterIdx
	}

	sql := `SELECT core.GET_SERIE($1) AS SERIE;`
	item, err := DBQueryOne(db, sql, tag)
	if err != nil {
		console.Error(err)
		return 0
	}

	return item.Int("serie")
}

func GetCode(tag, prefix string) string {
	num := GetSerie(tag)

	if len(prefix) == 0 {
		return Format("%08v", num)
	} else {
		return Format("%s%08v", prefix, num)
	}
}

func GetSerieLast(tag string) int {
	db := 0
	if MasterIdx != db {
		db = MasterIdx
	}

	query := `
  SELECT VALUE AS SERIE
  FROM core.SERIES
  WHERE SERIE=$1 LIMIT 1;`
	item, err := DBQueryOne(db, query, tag)
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
	query := `SELECT core.SET_SERIE($1, $2) AS SERIE;`
	item, err := DBQueryOne(db, query, tag, val)
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
		sql := Format(`
		SELECT A.*
		FROM core.SERIES A
		ORDER BY A.SERIE
		LIMIT %d OFFSET %d;`, rows, offset)

		items, err := Query(sql)
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
