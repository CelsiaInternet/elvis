package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/instances"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
	"github.com/go-chi/chi"
)

var workFlows *WorkFlows

/**
* Load
* @return error
 */
func Load(db *jdb.DB, schemaName string) error {
	if workFlows != nil {
		return nil
	}

	_, err := event.Load()
	if err != nil {
		return err
	}

	workFlows = newWorkFlows()
	err = instances.Define(db, schemaName)
	if err != nil {
		return err
	}

	SetLoadInstance(instances.Load)
	SetSaveInstance(instances.Save)

	return nil
}

/**
* HealthCheck
* @return bool
**/
func HealthCheck() bool {
	if workFlows == nil {
		return false
	}

	return workFlows.healthCheck()
}

/**
* New
* @param tag, version, name, description string, fn FnContext, createdBy string
* @return *Flow
**/
func New(tag, version, name, description string, fn FnContext, stop bool, createdBy string) *Flow {
	if workFlows == nil {
		logs.Panic(MSG_WORKFLOWS_NOT_LOAD)
		return nil
	}

	return workFlows.newFlow(tag, version, name, description, fn, stop, createdBy)
}

/**
* Run
* @param instanceId, tag string, step int, tags et.Json, ctx et.Json, createdBy string
* @return et.Json, error
**/
func Run(instanceId, tag string, step int, tags et.Json, ctx et.Json, createdBy string) (et.Json, error) {
	if workFlows == nil {
		return et.Json{}, errors.New(MSG_WORKFLOWS_NOT_LOAD)
	}

	return workFlows.runInstance(instanceId, tag, step, tags, ctx, createdBy)
}

/**
* Reset
* @param instanceId, updatedBy string
* @return error
**/
func Reset(instanceId, updatedBy string) error {
	if workFlows == nil {
		return errors.New(MSG_WORKFLOWS_NOT_LOAD)
	}

	return workFlows.resetInstance(instanceId, updatedBy)
}

/**
* Rollback
* @param instanceId, updatedBy string
* @return et.Json, error
**/
func Rollback(instanceId, updatedBy string) (et.Json, error) {
	if workFlows == nil {
		return et.Json{}, errors.New(MSG_WORKFLOWS_NOT_LOAD)
	}

	return workFlows.rollback(instanceId, updatedBy)
}

/**
* Stop
* @param instanceId, updatedBy string
* @return error
**/
func Stop(instanceId, updatedBy string) error {
	if workFlows == nil {
		return errors.New(MSG_WORKFLOWS_NOT_LOAD)
	}

	return workFlows.stop(instanceId, updatedBy)
}

/**
* SetStatus
* @param instanceId, status, updatedBy string
* @return FlowStatus, error
**/
func Status(instanceId, status, updatedBy string) (FlowStatus, error) {
	if workFlows == nil {
		return "", errors.New(MSG_WORKFLOWS_NOT_LOAD)
	}

	if _, ok := FlowStatusList[FlowStatus(status)]; !ok {
		return "", fmt.Errorf("status %s no es valido", status)
	}

	instance, exists := workFlows.loadInstance(instanceId)
	if !exists {
		return "", fmt.Errorf("instance not found")
	}

	instance.setStatus(FlowStatus(status))
	return instance.Status, nil
}

/**
* DeleteFlow
* @param tag string
* @return (bool, error)
**/
func DeleteFlow(tag string) (bool, error) {
	if workFlows == nil {
		return false, errors.New(MSG_WORKFLOWS_NOT_LOAD)
	}

	return workFlows.deleteFlow(tag), nil
}

/**
* GetInstance
* @param instanceId string
* @return (*Instance, error)
**/
func GetInstance(instanceId string) (*Instance, error) {
	if workFlows == nil {
		return nil, errors.New(MSG_WORKFLOWS_NOT_LOAD)
	}

	instance, exists := workFlows.loadInstance(instanceId)
	if !exists {
		return nil, fmt.Errorf("instance not found")
	}

	return instance, nil
}

/**
* HttpGet
* @params w http.ResponseWriter, r *http.Request
**/
func HttpGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var instance Instance
	err := loadInstance(id, &instance)
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
	id := chi.URLParam(r, "id")
	var instance Instance
	err := loadInstance(id, &instance)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	body, err := response.GetBody(r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	status := body.Str("status")
	err = instance.setStatus(FlowStatus(status))
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
func HttpSetParams(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var instance Instance
	err := loadInstance(id, &instance)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	body, err := response.GetBody(r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	jsonData := instance.ToJson()
	for k, v := range body {
		keys := strings.Split(k, "->")
		jsonData = et.SetNested(jsonData, keys, v)
	}

	bt := jsonData.ToByte()
	err = json.Unmarshal(bt, &instance)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok:     true,
		Result: instance.ToJson(),
	})
}
