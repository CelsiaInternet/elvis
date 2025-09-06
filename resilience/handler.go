package resilience

import (
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
)

/**
* Load
* @return error
 */
func Load() error {
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

	resilience = NewResilence()
	return nil
}

/**
* AddCustom
* @param id, tag, description string, totalAttempts int, timeAttempts time.Duration, fn interface{}, fnArgs ...interface{}
* @return *Attempt
 */
func AddCustom(id, tag, description string, totalAttempts int, timeAttempts time.Duration, fn interface{}, fnArgs ...interface{}) *Attempt {
	if resilience == nil {
		logs.Log("resilience", "resilience is nil")
		return nil
	}

	result := NewAttempt(id, tag, description, totalAttempts, timeAttempts, fn, fnArgs...)
	resilience.Attempts = append(resilience.Attempts, result)
	logs.Log("resilience", "add:", result.Json().ToString())
	resilience.Notify(result)
	resilience.Run(result)

	return result
}

/**
* Add
* @param tag, description string, fn interface{}, fnArgs ...interface{}
* @return *Attempt
 */
func Add(id, tag, description string, fn interface{}, fnArgs ...interface{}) *Attempt {
	return AddCustom(id, tag, description, resilience.TotalAttempts, resilience.TimeAttempts, fn, fnArgs...)
}

/**
* HealthCheck
* @return bool
**/
func HealthCheck() bool {
	if resilience == nil {
		return false
	}

	return resilience.HealthCheck()
}

/**
* HttpGetResilience
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilience(w http.ResponseWriter, r *http.Request) {
	if resilience == nil {
		response.JSON(w, r, http.StatusServiceUnavailable, et.Json{
			"message": "resilience is not initialized",
		})
		return
	}

	data := resilience.Json()
	response.JSON(w, r, http.StatusOK, data)
}

/**
* HttpGetResilienceById
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilienceById(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	id := body.Str("id")
	attempt := resilience.GetById(id)
	if attempt == nil {
		response.JSON(w, r, http.StatusNotFound, et.Json{
			"message": "attempt not found",
		})
		return
	}

	response.JSON(w, r, http.StatusOK, attempt.Json())
}

/**
* HttpGetResilienceByTag
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilienceByTag(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	tag := body.Str("tag")
	attempt := resilience.GetByTag(tag)
	if attempt == nil {
		response.JSON(w, r, http.StatusNotFound, et.Json{
			"message": "attempt not found",
		})
		return
	}

	response.JSON(w, r, http.StatusOK, attempt.Json())
}
