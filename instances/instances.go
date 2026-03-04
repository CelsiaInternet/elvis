package instances

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/utility"
)

var Instances *linq.Model

func Define(db *jdb.DB, schemaName string) error {
	if err := defineSchema(db, schemaName); err != nil {
		return console.Panic(err)
	}

	if Instances != nil {
		return nil
	}

	Instances = linq.NewModel(schema, "INSTANCES", "Tabla", 1)
	Instances.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Instances.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Instances.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Instances.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Instances.DefineColum("tag", "", "VARCHAR(80)", "-1")
	Instances.DefineColum("definition", "", "BYTEA", "")
	Instances.DefinePrimaryKey([]string{"_id"})
	Instances.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"index",
	})

	if err := Instances.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* Load
* @param id string, dest any
* @return error
**/
func Load(id string, dest any) error {
	if Instances == nil {
		return fmt.Errorf("model not found")
	}

	items, err := Instances.
		Data().
		Where(Instances.Column("_id").Eq(id)).
		First()
	if err != nil {
		return err
	}

	if !items.Ok {
		return errors.New("Instance not found")
	}

	scr, err := items.Byte("definition")
	if err != nil {
		return err
	}

	err = json.Unmarshal(scr, dest)
	if err != nil {
		return err
	}

	return nil
}

/**
* Save
* @param id string, tag string, definition []byte
* @return error
**/
func Save(id, tag string, obj any) error {
	if Instances == nil {
		return nil
	}

	bt, ok := obj.([]byte)
	if !ok {
		var err error
		bt, err = json.Marshal(obj)
		if err != nil {
			return err
		}
	}

	items, err := Instances.
		Data().
		Where(Instances.Column("_id").Eq(id)).
		First()
	if err != nil {
		return err
	}

	now := utility.Now()
	if !items.Ok {
		_, err := Instances.
			Insert(et.Json{
				"date_make":   now,
				"date_update": now,
				"_state":      utility.ACTIVE,
				"_id":         id,
				"tag":         tag,
				"definition":  bt,
			}).
			CommandOne()
		if err != nil {
			return err
		}

		return nil
	}

	_, err = Instances.
		Update(et.Json{
			"date_update": now,
			"_id":         id,
			"tag":         tag,
			"definition":  bt,
		}).
		Where(Instances.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return err
	}

	return nil
}

/**
* Delete
* @param id string
* @return error
**/
func Delete(id string) error {
	if Instances == nil {
		return nil
	}

	_, err := Instances.
		Delete().
		Where(Instances.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return err
	}

	return nil
}
