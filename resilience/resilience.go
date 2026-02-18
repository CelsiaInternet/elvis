package resilience

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/reg"
)

/**
* Instance
* @param id, tag, description string, totalAttempts int, timeAttempts time.Duration, tags et.Json, team string, level string, fn interface{}, fnArgs ...interface{}
* @return Instance
 */
func NewInstance(id, tag, description string, totalAttempts int, timeAttempts time.Duration, tags et.Json, team string, level string, fn interface{}, fnArgs ...interface{}) *Instance {
	id = reg.GetUUID(id)
	result := &Instance{
		CreatedAt:     time.Now(),
		Id:            id,
		Tag:           tag,
		Description:   description,
		fn:            fn,
		fnArgs:        fnArgs,
		fnResult:      []reflect.Value{},
		TotalAttempts: totalAttempts,
		TimeAttempts:  timeAttempts,
		Tags:          tags,
		Team:          team,
		Level:         level,
		stop:          false,
	}
	result.setStatus(StatusPending)

	return result
}

/**
* LoadById
* @param id string
* @return *Instance, error
**/
func LoadById(id string) (*Instance, error) {
	key := fmt.Sprintf("resilience:%s", id)
	exists := cache.Exists(key)
	if !exists {
		return nil, fmt.Errorf(MSG_INSTANCE_NOT_FOUND)
	}

	bt, err := cache.Get(key, "")
	if err != nil {
		return nil, err
	}

	var result Instance
	err = json.Unmarshal([]byte(bt), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
