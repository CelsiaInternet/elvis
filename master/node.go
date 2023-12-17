package master

import (
	"time"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
	e "github.com/cgalvisleon/elvis/json"
)

const NodeStatusIdle = 0
const NodeStatusActive = 1
const NodeStatusWorking = 2
const NodeStatusSync = 3
const NodeStatusError = 4

type Node struct {
	Db          int
	URL         string
	Date_make   time.Time `json:"date_make"`
	Date_update time.Time `json:"date_update"`
	Id          string    `json:"_id"`
	Mode        int       `json:"mode"`
	Data        e.Json    `json:"_data"`
	Status      int       `json:"status"`
	Index       int       `json:"index"`
}

func (n *Node) Scan(data *e.Json) error {
	n.Date_make = data.Time("date_make")
	n.Date_update = data.Time("date_update")
	n.Id = data.Str("_id")
	n.Mode = data.Int("mode")
	n.Data = data.Json("_data")
	n.Index = data.Int("index")
	n.Status = NodeStatusIdle

	return nil
}

func (c *Node) LatIndex() int {
	sql := `
	SELECT INDEX FROM core.MODE
	LIMIT 1;`

	item, err := jdb.DBQueryOne(c.Db, sql)
	if err != nil {
		return -1
	}

	return item.Index()
}

func (c *Node) GetSyncByIdT(idT string) (e.Item, error) {
	sql := `
  SELECT *
  FROM core.SYNC
  WHERE _IDT=$1
  LIMIT 1;`

	item, err := jdb.DBQueryOne(c.Db, sql, idT)
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func (c *Node) DelSyncByIndex(index int) error {
	sql := `
  DELETE FROM core.SYNC
  WHERE INDEX=$1;`

	_, err := jdb.DBQueryOne(c.Db, sql, index)
	if err != nil {
		return err
	}

	return nil
}

func NewNode(params *e.Json) (*Node, error) {
	result := &Node{}
	err := result.Scan(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DefineNodes() error {
	if master.InitNodes {
		return nil
	}

	sql := `
  -- DROP TABLE IF EXISTS core.NODES CASCADE;
	-- DROP TABLE IF EXISTS core.SYNC CASCADE;
	CREATE SCHEMA IF NOT EXISTS "core";

  CREATE TABLE IF NOT EXISTS core.NODES(
		DATE_MAKE TIMESTAMP DEFAULT NOW(),
		DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    _ID VARCHAR(80) DEFAULT '',
    MODE INTEGER DEFAULT 0,
		PASSWORD VARCHAR(250) DEFAULT '',
    _DATA JSONB DEFAULT '{}',
		INDEX SERIAL,
		PRIMARY KEY(_ID)
	);

	CREATE OR REPLACE FUNCTION core.NODES_UPSET()
	RETURNS
		TRIGGER AS $$
	BEGIN
		PERFORM pg_notify(
			'node',
			json_build_object(
				'action', TG_OP,
				'node', NEW._ID
			)::text
		);

	RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.NODES_DELETE()
	RETURNS
		TRIGGER AS $$
	BEGIN
		PERFORM pg_notify(
			'node',
			json_build_object(
				'action', TG_OP,
				'node', OLD._ID
			)::text
		);

	RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS NODES_UPSET ON core.NODES CASCADE;
	CREATE TRIGGER NODES_UPSET
	AFTER INSERT OR UPDATE ON core.NODES
	FOR EACH ROW
	EXECUTE PROCEDURE core.NODES_UPSET();

	DROP TRIGGER IF EXISTS NODES_DELETE ON core.NODES CASCADE;
	CREATE TRIGGER NODES_DELETE
	AFTER DELETE ON core.NODES
	FOR EACH ROW
	EXECUTE PROCEDURE core.NODES_DELETE();

  CREATE TABLE IF NOT EXISTS core.SYNC(
		DATE_MAKE TIMESTAMP DEFAULT NOW(),
		DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    ACTION VARCHAR(80) DEFAULT '',		
    _IDT VARCHAR(80) DEFAULT '-1',
		_DATA JSONB DEFAULT '{}',
		QUERY TEXT DEFAULT '',
		NODE VARCHAR(80) DEFAULT '',
    INDEX BIGINT DEFAULT 0,
		PRIMARY KEY (_IDT)
  );
	CREATE INDEX IF NOT EXISTS SYNC_NODE_IDX ON core.SYNC(NODE);
  CREATE INDEX IF NOT EXISTS SYNC_INDEX_IDX ON core.SYNC(INDEX);`

	_, err := jdb.QDDL(sql)
	if err != nil {
		return console.PanicE(err)
	}

	master.InitNodes = true

	go master.LoadNodes()

	return nil
}

/**
* Mode
*	Handler for CRUD data
 */
func GetNodeById(id string) (e.Item, error) {
	sql := `
	SELECT
	A._DATA||
  jsonb_build_object(
    'mode', A.MODE,
		'index', A.INDEX
  ) AS _DATA
	FROM core.NODES A
	WHERE A._ID=$1
	LIMIT 1;`

	item, err := jdb.QueryDataOne(sql, id)
	if err != nil {
		return e.Item{}, err
	}

	delete(item.Result, "password")

	return item, nil
}

func DeleteNodeById(id string) (e.Item, error) {
	sql := `
	DELETE FROM core.NODES	
	WHERE _ID=$1
	RETURNING *;`

	item, err := jdb.QueryDataOne(sql, id)
	if err != nil {
		return e.Item{}, err
	}

	delete(item.Result, "password")

	return item, nil
}

func AllNodes(search string, page, rows int) (e.List, error) {
	sql := `
	SELECT COUNT(*) AS COUNT
	FROM core.NODES A
	WHERE CONCAT('MODE:', A.MODE, ':DATA:', A._DATA::TEXT) ILIKE CONCAT('%', $1, '%');`

	all := jdb.QueryCount(sql, search)

	sql = `
	SELECT A._DATA||
	jsonb_build_object(
		'data_make', A.DATE_MAKE,
		'date_update', A.DATE_UPDATE,
		'_id', A._ID,
		'mode', A.MODE,
		'index', A.INDEX
	) AS _DATA
	FROM core.NODES A
	WHERE CONCAT('MODE:', A.MODE, ':DATA:', A._DATA::TEXT) ILIKE CONCAT('%', $1, '%')
	LIMIT $2 OFFSET $3;`

	offset := (page - 1) * rows
	items, err := jdb.Query(sql, search, rows, offset)
	if err != nil {
		return e.List{}, err
	}

	return items.ToList(all, page, rows), nil
}
