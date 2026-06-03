package authorization

import (
	"errors"
	"fmt"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/dt"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
)

type Authorization struct {
	schema *linq.Schema
	model  *linq.Model
}

var (
	inb            *Authorization
	ErrorSetAuthor = fmt.Errorf(msg.RECORD_NOT_FOUND)
)

/**
* Load
* @param db *jdb.DB, schema, name string
* @return (*Authorization, error)
**/
func Load(db *jdb.DB, schema, name string) (*Authorization, error) {
	if inb != nil {
		return inb, nil
	}

	var err error
	inb, err = Define(db, schema, name)
	if err != nil {
		return nil, err
	}

	return inb, nil
}

/**
* Define
* @param db *jdb.DB, schema, name string
* @return (*Instance, error)
**/
func Define(db *jdb.DB, schema, name string) (*Authorization, error) {
	schemaObj, err := defineSchema(db, schema)
	if err != nil {
		return nil, console.Panic(err)
	}

	if name == "" {
		name = "instances"
	}

	model := linq.NewModel(schemaObj, name, "Tabla", 1)
	model.DefineColum("created_at", "", "TIMESTAMP", "NOW()")
	model.DefineColum("project_id", "", "VARCHAR(80)", "")
	model.DefineColum("profile_id", "", "VARCHAR(80)", "")
	model.DefineColum("method", "", "VARCHAR(80)", "")
	model.DefineColum("path", "", "VARCHAR(250)", "")
	model.DefinePrimaryKey([]string{"project_id", "profile_id", "method", "path"})

	if err := model.Init(); err != nil {
		return nil, err
	}

	return &Authorization{
		schema: schemaObj,
		model:  model,
	}, nil
}

/**
* Author
* @param projectId, profileId, method, path string
* @return et.Item, error
**/
func (s *Authorization) Author(projectId, profileId, method, path string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", projectId, profileId, method, path)
	result := dt.Get(key)
	if result.Ok {
		return result.Bool("ok"), nil
	}

	item, err := s.model.
		Select().
		Where(s.model.Column("project_id").Eq(projectId)).
		And(s.model.Column("profile_id").Eq(profileId)).
		And(s.model.Column("method").Eq(method)).
		And(s.model.Column("path").Eq(path)).
		First()
	if err != nil {
		return false, err
	}

	dt.Up(key, et.Item{Ok: item.Ok, Result: et.Json{"ok": item.Ok}})
	return item.Ok, nil
}

/**
* RemoveAuthor
* @param projectId, profileId, method, path string
* @return error
**/
func (s *Authorization) RemoveAuthor(projectId, profileId, method, path string) error {
	key := fmt.Sprintf("%s:%s:%s:%s", projectId, profileId, method, path)
	dt.Drop(key)

	_, err := s.model.
		Delete().
		Where(s.model.Column("project_id").Eq(projectId)).
		And(s.model.Column("profile_id").Eq(profileId)).
		And(s.model.Column("method").Eq(method)).
		And(s.model.Column("path").Eq(path)).
		All()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	event.Publish(EVENT_DEL_AUTHORIZATION, et.Json{key: key})
	return nil
}

/**
* StateAuthorizationes
* @param id, stateId, createdBy string
* @return et.Item, error
**/
func (s *Authorization) SetAuthor(projectId, profileId, method, path string) error {
	if !utility.ValidStr(method, 0, []string{""}) {
		return fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "method")
	}
	if !utility.ValidStr(path, 0, []string{""}) {
		return fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "path")
	}

	key := fmt.Sprintf("%s:%s:%s:%s", projectId, profileId, method, path)
	now := timezone.Now()
	_, err := s.model.
		Insert(et.Json{
			"created_at": now,
			"project_id": projectId,
			"profile_id": profileId,
			"method":     method,
			"path":       path,
		}).
		All()
	if err != nil {
		return err
	}

	dt.Drop(key)

	return nil
}

/**
* SetPath
* @params method, path string
* @return error
**/
func (s *Authorization) SetPath(method, path string) error {
	err := s.SetAuthor("", "", method, path)
	if err != nil && !errors.Is(err, ErrorSetAuthor) {
		return err
	}

	return nil
}
