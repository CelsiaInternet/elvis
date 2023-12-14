package core

import (
	"github.com/cgalvisleon/elvis/console"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

type Collection struct {
	Name      string
	Id        string
	ProjectId string
	Result    e.Json
}

var Collections *linq.Model

func DefineCollection() error {
	if err := defineSchema(); err != nil {
		return console.PanicE(err)
	}

	if Collections != nil {
		return nil
	}

	Collections = linq.NewModel(SchemaCore, "COLLECTION", "Tabla de colecciones", 1)
	Collections.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Collections.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Collections.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Collections.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Collections.DefineColum("collection", "", "VARCHAR(80)", "")
	Collections.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Collections.DefineColum("_data", "", "JSONB", "{}")
	Collections.DefineColum("expiration", "", "INTEGER", 0)
	Collections.DefineColum("index", "", "INTEGER", 0)
	Collections.DefinePrimaryKey([]string{"collection", "_id"})
	Collections.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"project_id",
		"expiration",
		"index",
	})
	Collections.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		collection := new.Str("collection")

		if collection == "__telemetry" {
			return nil
		}

		item, err := Collections.Select().
			Where(Collections.Column("collection").Eq("__telemetry")).
			And(Collections.Column("_id").Eq(collection)).
			First()
		if err != nil {
			return err
		}

		projectId := new.Key("project_id")
		count := item.Int("count")
		count++
		data = e.Json{}
		data["collection"] = "__telemetry"
		data["project_id"] = projectId
		data["_id"] = collection
		data["expiration"] = 0
		data["count"] = count
		item, err = Collections.Upsert(data).
			Where(Collections.Column("collection").Eq("__telemetry")).
			And(Collections.Column("_id").Eq(collection)).
			Command()
		if err != nil {
			return err
		}

		return nil
	})
	Collections.Trigger(linq.AfterDelete, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		collection := old.Str("collection")

		if collection == "__telemetry" {
			return nil
		}

		item, err := Collections.Select().
			Where(Collections.Column("collection").Eq("__telemetry")).
			And(Collections.Column("_id").Eq(collection)).
			First()
		if err != nil {
			return err
		}

		projectId := old.Key("project_id")
		count := item.Int("count")
		count--
		data = e.Json{}
		data["collection"] = "__telemetry"
		data["project_id"] = projectId
		data["_id"] = collection
		data["expiration"] = 0
		data["count"] = count
		item, err = Collections.Upsert(data).
			Where(Collections.Column("collection").Eq("__telemetry")).
			And(Collections.Column("_id").Eq(collection)).
			Command()
		if err != nil {
			return err
		}

		return nil
	})

	return InitModel(Collections)
}

func GetCollection(collection, id string) *Collection {
	item, err := GetCollectionById(collection, id)
	if err != nil {
		return &Collection{}
	}

	return &Collection{
		Name:      collection,
		Id:        id,
		ProjectId: item.Key("project_id"),
		Result:    item.Result,
	}
}

func (c *Collection) Set(atrib string, val any) error {
	c.Result.Set(atrib, val)
	item, err := UpSertCollection(c.Name, c.ProjectId, c.Id, c.Result)
	if err != nil {
		return err
	}

	c.Result = item.Result
	return nil
}

func (c *Collection) Int(atribs ...string) int {
	return c.Result.Int(atribs...)
}

func (c *Collection) Str(atribs ...string) string {
	return c.Result.Str(atribs...)
}

/**
* Collection
*	Handler for CRUD data
 */
func GetCollectionById(collection, id string) (e.Item, error) {
	return Collections.Select().
		Where(Collections.Column("collection").Eq(collection)).
		And(Collections.Column("_id").Eq(id)).
		First()
}

func GetCollectionByIndex(collection string, idx int) (e.Item, error) {
	return Collections.Select().
		Where(Collections.Column("collection").Eq(collection)).
		And(Collections.Column("index").Eq(idx)).
		First()
}

func UpSertCollection(collection, projectId, id string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(collection, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "collection")
	}

	if projectId == "" {
		projectId = "-1"
	}

	id = utility.GenId(id)
	data["collection"] = collection
	data["project_id"] = projectId
	data["_id"] = id
	data["expiration"] = 0
	return Collections.Upsert(data).
		Where(Collections.Column("collection").Eq(collection)).
		And(Collections.Column("_id").Eq(id)).
		Command()
}

func StateCollection(collection, id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	return Collections.Upsert(e.Json{
		"_state": state,
	}).
		Where(Collections.Column("collection").Eq(collection)).
		And(Collections.Column("_id").Eq(id)).
		And(Collections.Column("_state").Neg(state)).
		Command()
}

func DeleteCollection(collection, id string) (e.Item, error) {
	return Collections.Delete().
		Where(Collections.Column("collection").Eq(collection)).
		And(Collections.Column("_id").Eq(id)).
		Command()
}

func AllCollections(projectId, collection, state, search string, page, rows int) (e.List, error) {
	if !utility.ValidId(projectId) {
		return e.List{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidStr(collection, 0, []string{""}) {
		return e.List{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "collection")
	}

	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if auxState == "*" {
		state = utility.FOR_DELETE

		return Collections.Select().
			Where(Collections.Column("collection").Eq(collection)).
			And(Collections.Column("_state").Neg(state)).
			And(Collections.Column("project_id").In("-1", projectId)).
			And(Collections.Column("_data").Cast("TEXT").Like("%"+search+"%")).
			OrderBy(Collections.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Collections.Select().
			Where(Collections.Column("collection").Eq(collection)).
			And(Collections.Column("_state").In("-1", state)).
			And(Collections.Column("project_id").In("-1", projectId)).
			And(Collections.Column("_data").Cast("TEXT").Like("%"+search+"%")).
			OrderBy(Collections.Column("name"), true).
			List(page, rows)
	} else {
		return Collections.Select().
			Where(Collections.Column("collection").Eq(collection)).
			And(Collections.Column("_state").Eq(state)).
			And(Collections.Column("project_id").In("-1", projectId)).
			And(Collections.Column("_data").Cast("TEXT").Like("%"+search+"%")).
			OrderBy(Collections.Column("name"), true).
			List(page, rows)
	}
}
