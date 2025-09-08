package resilience

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/response"
	"github.com/go-chi/chi"
)

var resilience map[string]*Instance

/**
* load
* @return error
 */
func load() error {
	if resilience != nil {
		return nil
	}

	_, err := cache.Load()
	if err != nil {
		return err
	}

	_, err = event.Load()
	if err != nil {
		return err
	}

	initEvents()

	resilience = make(map[string]*Instance)

	return nil
}

/**
* HealthCheck
* @return bool
**/
func HealthCheck() bool {
	if err := load(); err != nil {
		return false
	}

	if !cache.HealthCheck() {
		return false
	}

	if !event.HealthCheck() {
		return false
	}

	return true
}

/**
* AddCustom
* @param id, tag, description string, totalAttempts int, timeAttempts, retentionTime time.Duration, fn interface{}, fnArgs ...interface{}
* @return *Instance
 */
func AddCustom(id, tag, description string, totalAttempts int, timeAttempts, retentionTime time.Duration, fn interface{}, fnArgs ...interface{}) *Instance {
	if err := load(); err != nil {
		return nil
	}

	result := NewInstance(id, tag, description, totalAttempts, timeAttempts, retentionTime, fn, fnArgs...)
	resilience[id] = result
	result.runAttempt()

	return result
}

/**
* Add
* @param tag, description string, fn interface{}, fnArgs ...interface{}
* @return *Instance
 */
func Add(id, tag, description string, fn interface{}, fnArgs ...interface{}) *Instance {
	totalAttempts := envar.EnvarInt(3, "RESILIENCE_TOTAL_ATTEMPTS")
	timeAttempts := envar.EnvarNumber(30, "RESILIENCE_TIME_ATTEMPTS")
	retentionTime := envar.EnvarNumber(10, "RESILIENCE_RETENTION_TIME")

	return AddCustom(id, tag, description, totalAttempts, time.Duration(timeAttempts)*time.Second, time.Duration(retentionTime)*time.Minute, fn, fnArgs...)
}

/**
* Stop
* @param id string
* @return error
 */
func Stop(id string) error {
	if err := load(); err != nil {
		return err
	}

	if _, ok := resilience[id]; !ok {
		return fmt.Errorf(MSG_ID_NOT_FOUND)
	}

	resilience[id].setStop()

	return nil
}

/**
* Restart
* @param id string
* @return error
 */
func Restart(id string) error {
	if err := load(); err != nil {
		return err
	}

	if _, ok := resilience[id]; !ok {
		return fmt.Errorf(MSG_ID_NOT_FOUND)
	}

	resilience[id].setRestart()

	return nil
}

/**
* HttpGetResilienceById
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilienceById(w http.ResponseWriter, r *http.Request) {
	if err := load(); err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, MSG_RESILIENCE_NOT_INITIALIZED)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.HTTPError(w, r, http.StatusBadRequest, MSG_ID_REQUIRED)
		return
	}

	res, err := loadById(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok:     true,
		Result: res.ToJson(),
	})
}

/**
* HttpGetResilienceStop
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilienceStop(w http.ResponseWriter, r *http.Request) {
	if err := load(); err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, MSG_RESILIENCE_NOT_INITIALIZED)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.HTTPError(w, r, http.StatusBadRequest, MSG_ID_REQUIRED)
		return
	}

	res, ok := resilience[id]
	if !ok {
		response.HTTPError(w, r, http.StatusNotFound, MSG_ID_NOT_FOUND)
		return
	}

	result := res.setStop()
	response.ITEM(w, r, http.StatusOK, result)
}

/**
* HttpGetResilienceRestart
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilienceRestart(w http.ResponseWriter, r *http.Request) {
	if err := load(); err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, MSG_RESILIENCE_NOT_INITIALIZED)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.HTTPError(w, r, http.StatusBadRequest, MSG_ID_REQUIRED)
		return
	}

	res, ok := resilience[id]
	if !ok {
		response.HTTPError(w, r, http.StatusNotFound, MSG_ID_NOT_FOUND)
		return
	}

	result := res.setRestart()
	response.ITEM(w, r, http.StatusOK, result)
}
