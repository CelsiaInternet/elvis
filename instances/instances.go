package instances

import (
	"encoding/json"
	"fmt"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/utility"
)

type Instance struct {
	schema *linq.Schema
	model  *linq.Model
}

var instance *Instance

func Define(db *jdb.DB, schema, name string) (*Instance, error) {
	if instance != nil {
		return instance, nil
	}

	instance = &Instance{}

	if err := instance.defineSchema(db, schema); err != nil {
		return nil, console.Panic(err)
	}

	if name == "" {
		name = "instances"
	}

	instance.model = linq.NewModel(instance.schema, name, "Tabla", 1)
	instance.model.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	instance.model.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	instance.model.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	instance.model.DefineColum("_id", "", "VARCHAR(80)", "-1")
	instance.model.DefineColum("tag", "", "VARCHAR(80)", "-1")
	instance.model.DefineColum("definition", "", "BYTEA", "")
	instance.model.DefinePrimaryKey([]string{"_id"})
	instance.model.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"index",
	})

	if err := instance.model.Init(); err != nil {
		return nil, err
	}

	return instance, nil
}

/**
* Load
* @param id string, dest any
* @return bool, error
**/
func (s *Instance) Get(id string, dest any) (bool, error) {
	if s.model == nil {
		return false, fmt.Errorf("model not found")
	}

	items, err := s.model.
		Data().
		Where(s.model.Column("_id").Eq(id)).
		First()
	if err != nil {
		return false, err
	}

	if !items.Ok {
		return false, nil
	}

	scr, err := items.Byte("definition")
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(scr, dest)
	if err != nil {
		return false, err
	}

	return true, nil
}

/**
* Save
* @param id string, tag string, definition []byte
* @return error
**/
func (s *Instance) Set(id, tag string, obj any) error {
	if s.model == nil {
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

	items, err := s.model.
		Data().
		Where(s.model.Column("_id").Eq(id)).
		First()
	if err != nil {
		return err
	}

	now := utility.Now()
	if !items.Ok {
		_, err := s.model.
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

	_, err = s.model.
		Update(et.Json{
			"date_update": now,
			"_id":         id,
			"tag":         tag,
			"definition":  bt,
		}).
		Where(s.model.Column("_id").Eq(id)).
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
func (s *Instance) Delete(id string) error {
	if s.model == nil {
		return nil
	}

	_, err := s.model.
		Delete().
		Where(s.model.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return err
	}

	return nil
}
