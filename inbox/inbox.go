package inbox

import (
	"fmt"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/dt"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
)

type Inbox struct {
	schema *linq.Schema
	model  *linq.Model
}

var inb *Inbox

/**
* Load
* @param db *jdb.DB, schema, name string
* @return (*Inbox, error)
**/
func Load(db *jdb.DB, schema, name string) (*Inbox, error) {
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
func Define(db *jdb.DB, schema, name string) (*Inbox, error) {
	schemaObj, err := defineSchema(db, schema)
	if err != nil {
		return nil, console.Panic(err)
	}

	if name == "" {
		name = "instances"
	}

	model := linq.NewModel(schemaObj, name, "Tabla", 1)
	model.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	model.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	model.DefineColum("project_id", "", "VARCHAR(80)", "")
	model.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	model.DefineColum("_id", "", "VARCHAR(80)", "-1")
	model.DefineColum("user_id", "", "VARCHAR(80)", "-1")
	model.DefineColum("app_id", "", "VARCHAR(80)", "-1")
	model.DefineColum("kind", "", "VARCHAR(80)", "-1")
	model.DefineColum("code", "", "VARCHAR(80)", "-1")
	model.DefineColum("title", "", "VARCHAR(250)", "-1")
	model.DefineColum("_data", "", "JSONB", "{}")
	model.DefinePrimaryKey([]string{"_id"})
	model.DefineIndex([]string{
		"date_make",
		"date_update",
		"project_id",
		"_state",
		"user_id",
		"app_id",
		"kind",
		"code",
		"title",
	})

	if err := model.Init(); err != nil {
		return nil, err
	}

	return &Inbox{
		schema: schemaObj,
		model:  model,
	}, nil
}

/**
* GetInboxesById
* @param id string
* @return et.Item, error
**/
func (s *Inbox) GetInboxesById(id string) (et.Item, error) {
	if s.model == nil {
		return et.Item{}, fmt.Errorf("model not found")
	}

	result, err := s.model.
		Data().
		Where(s.model.Column("_id").Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* GetInboxesByCode
* @param code string
* @return et.Items, error
**/
func (s *Inbox) GetInboxesByCode(code string) (et.Items, error) {
	if s.model == nil {
		return et.Items{}, fmt.Errorf("model not found")
	}

	result, err := s.model.
		Data().
		Where(s.model.Column("code").Eq(code)).
		All()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* GetInboxesByMy
* @param userId, appId, inbox, status string, page, rows int
* @return et.Items, error
**/
func (s *Inbox) GetInboxesByMy(userId, appId, inbox, status string, page, rows int) (et.Items, error) {
	ql := s.model.
		Data().
		Where(s.model.Column("_id").Eq(userId)).
		And(s.model.Column("app_id").Eq(appId)).
		And(s.model.Column("inbox").Eq(inbox))
	if status == "0" {
		ql = ql.And(s.model.Column("_status").In(status, "-2"))
	} else {
		ql = ql.And(s.model.Column("_status").Eq(status))
	}

	result, err := ql.
		OrderBy(s.model.Column("updated_at"), false).
		Page(page, rows)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* GenInboxesCode
* @param projectId string
* @return string, error
**/
func (s *Inbox) GenInboxesCode(projectId string) (string, error) {
	code, err := jdb.GetSeries("services", projectId)
	if err != nil {
		return "", err
	}

	return code, nil
}

/**
* UpsertInboxes
* @param projectId, id string, userId, appId, inbox string, data et.Json, createdBy string
* @return et.Item, error
**/
func (s *Inbox) UpsertInboxes(projectId, id, userId, appId, inbox string, data et.Json, createdBy string) (et.Item, error) {
	if !utility.ValidStr(projectId, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidStr(id, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "_id")
	}

	if !utility.ValidStr(userId, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "user_id")
	}

	if !utility.ValidStr(appId, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "app_id")
	}

	if !utility.ValidStr(inbox, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "inbox")
	}

	id = reg.GetUUID(id)
	current, err := s.model.
		Data().
		Where(s.model.Column("_id").Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		now := timezone.Now()
		data["project_id"] = projectId
		data["_id"] = id
		data["inbox"] = inbox
		data["created_at"] = now
		data["updated_at"] = now
		data["status_id"] = utility.ACTIVE
		data["app_id"] = appId
		data["user_id"] = userId
		data["created_by"] = createdBy
		code, err := jdb.GetSeries(inbox, projectId)
		if err == nil {
			data["code"] = code
		}
		result, err := s.model.
			Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return result, nil
	}

	now := timezone.Now()
	data["project_id"] = projectId
	data["_id"] = id
	data["inbox"] = inbox
	data["updated_at"] = now
	data["updated_by"] = createdBy
	code, err := jdb.GetSeries(inbox, projectId)
	if err == nil {
		data["code"] = code
	}
	result, err := s.model.
		Update(data).
		Where(s.model.Column("_id").Eq(id)).
		And(s.model.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* StateInboxes
* @param id, stateId, createdBy string
* @return et.Item, error
**/
func (s *Inbox) StateInboxes(id, stateId, createdBy string) (et.Item, error) {
	if !utility.ValidStr(stateId, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "_state")
	}

	if !utility.ValidStr(id, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "_id")
	}

	result, err := s.model.
		Update(et.Json{
			"_state":     stateId,
			"updated_by": createdBy,
		}).
		Where(s.model.Column("_id").Eq(id)).
		And(s.model.Column("_state").Neg(stateId)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	dt.Drop(id)

	return et.Item{
		Ok: result.Ok,
		Result: et.Json{
			"message": msg.RECORD_UPDATE,
		},
	}, nil
}
