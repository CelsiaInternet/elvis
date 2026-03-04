package instances

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/workflow"
	"github.com/go-chi/chi"
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
* loadInstance
* @param id string, v any
* @return error
**/
func loadInstance(id string, v any) error {
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
		return workflow.ErrorInstanceNotFound
	}

	scr, err := items.Byte("definition")
	if err != nil {
		return err
	}

	err = json.Unmarshal(scr, v)
	if err != nil {
		return err
	}

	return nil
}

/**
* saveInstance
* @param id string, tag string, definition []byte
* @return error
**/
func saveInstance(id, tag string, definition []byte) error {
	if Instances == nil {
		return nil
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
				"definition":  definition,
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
			"definition":  definition,
		}).
		Where(Instances.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return err
	}

	return nil
}

/**
* Status - Update the status of a suspension instance
* @param id string, status string, createdBy string
* @return et.Item, error
**/
func Status(id, status, createdBy string) (et.Item, error) {
	if !utility.ValidStr(status, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "status")
	}

	if !utility.ValidStr(id, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, jdb.KEY)
	}

	st := workflow.FlowStatus(status)
	if _, exists := workflow.FlowStatusList[st]; !exists {
		return et.Item{}, fmt.Errorf("invalid status: %s", status)
	}

	instance, err := Load(id)
	if err != nil {
		return et.Item{}, err
	}

	instance.SetStatus(st)

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": msg.RECORD_UPDATE,
		},
	}, nil
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

/**
* Load
* @param id string
* @return workflow.Instance, error
**/
func Load(id string) (*workflow.Instance, error) {
	var instance *workflow.Instance

	err := loadInstance(id, &instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

/**
* Save
* @param instance workflow.Instance
* @return error
**/
func Save(instance *workflow.Instance) error {
	bt, err := instance.Serialize()
	if err != nil {
		return err
	}

	err = saveInstance(instance.Id, instance.Tag, bt)
	if err != nil {
		return err
	}

	return nil
}

/**
* HttpGet
* @params w http.ResponseWriter, r *http.Request
**/
func HttpGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	instance, err := Load(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok:     true,
		Result: instance.ToJson(),
	})
}

/**
* HttpGetInstance
* @params w http.ResponseWriter, r *http.Request
**/
func HttpState(w http.ResponseWriter, r *http.Request) {
	body, err := response.GetBody(r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	id := body.Str("id")
	status := body.Str("status")
	username := claim.ClientName(r)
	result, err := Status(id, status, username)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* HttpGetInstance
* @params w http.ResponseWriter, r *http.Request
**/
func HttpSetParams(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	body, err := response.GetBody(r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	instance, err := Load(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	for k, v := range body {
		instance.SetParam(k, v)
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok:     true,
		Result: instance.ToJson(),
	})
}

func init() {
	workflow.SetLoadInstance(Load)
	workflow.SetSaveInstance(Save)
}
