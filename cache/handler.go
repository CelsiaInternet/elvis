package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/redis/go-redis/v9"
)

const IsNil = redis.Nil

/**
* SetCtx
* @params ctx context.Context
* @params key string
* @params val string
* @params second time.Duration
* @return error
**/
func SetCtx(ctx context.Context, key, val string, second time.Duration) error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	duration := second * time.Second

	err := conn.db.Set(ctx, key, val, duration).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* GetCtx
* @params ctx context.Context
* @params key string
* @params def string
* @return string, error
**/
func GetCtx(ctx context.Context, key, def string) (string, error) {
	if conn == nil {
		return def, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	result, err := conn.db.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return def, IsNil
	case err != nil:
		return def, err
	default:
		return result, nil
	}
}

/**
* DelCtx
* @params ctx context.Context
* @params key string
* @return int64, error
**/
func DelCtx(ctx context.Context, key string) (int64, error) {
	if conn == nil {
		return 0, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	intCmd := conn.db.Del(ctx, key)

	return intCmd.Val(), intCmd.Err()
}

/**
* HSetCtx
* @params ctx context.Context
* @params key string
* @params val map[string]string
* @return error
**/
func HSetCtx(ctx context.Context, key string, val map[string]string) error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.db.HSet(ctx, key, val).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* HGetCtx
* @params ctx context.Context
* @params key string
* @return map[string]string, error
**/
func HGetCtx(ctx context.Context, key string) (map[string]string, error) {
	if conn == nil {
		return map[string]string{}, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	result := conn.db.HGetAll(ctx, key).Val()

	return result, nil
}

/**
* HDelCtx
* @params ctx context.Context
* @params key string
* @params atr string
* @return error
**/
func HDelCtx(ctx context.Context, key, atr string) error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.db.Do(ctx, "HDEL", key, atr).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* Get
* @params key string
* @params def string
* @return string, error
**/
func Get(key, def string) (string, error) {
	if conn == nil {
		return def, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	return GetCtx(conn.ctx, key, def)
}

/**
* Del
* @params key string
* @return int64, error
**/
func Del(key string) (int64, error) {
	if conn == nil {
		return 0, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	return DelCtx(conn.ctx, key)
}

/**
* Set
* @params key string
* @params val interface{}
* @params second time.Duration
* @return error
**/
func Set(key string, val interface{}, second time.Duration) error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	switch v := val.(type) {
	case et.Json:
		return SetCtx(conn.ctx, key, v.ToString(), second)
	case et.Items:
		return SetCtx(conn.ctx, key, v.ToString(), second)
	case et.Item:
		return SetCtx(conn.ctx, key, v.ToString(), second)
	default:
		val, ok := val.(string)
		if ok {
			return SetCtx(conn.ctx, key, val, second)
		}
	}

	return nil
}

/**
* SetH
* @params key string
* @params val interface{}
* @return error
**/
func SetH(key string, val interface{}) error {
	return Set(key, val, time.Hour*1)
}

/**
* SetD
* @params key string
* @params val interface{}
* @return error
**/
func SetD(key string, val interface{}) error {
	return Set(key, val, time.Hour*24)
}

/**
* SetW
* @params key string
* @params val interface{}
* @return error
**/
func SetW(key string, val interface{}) error {
	return Set(key, val, time.Hour*24*7)
}

/**
* SetM
* @params key string
* @params val interface{}
* @return error
**/
func SetM(key string, val interface{}) error {
	return Set(key, val, time.Hour*24*30)
}

/**
* SetY
* @params key string
* @params val interface{}
* @return error
**/
func SetY(key string, val interface{}) error {
	return Set(key, val, time.Hour*24*365)
}

/**
* Empty
* @return error
**/
func Empty() error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	ctx := context.Background()
	iter := conn.db.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		DelCtx(ctx, key)
	}

	return nil
}

/**
* More
* @params key string
* @params second time.Duration
* @return int
**/
func More(key string, second time.Duration) int {
	n, err := Get(key, "")
	if err != nil {
		n = "0"
	}

	if n == "" {
		n = "0"
	}

	val, err := strconv.Atoi(n)
	if err != nil {
		return 0
	}

	val++
	Set(key, strs.Format(`%d`, val), second)

	return val
}

/**
* HSet
* @params key string
* @params val map[string]string
* @return error
**/
func HSet(key string, val map[string]string) error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	return HSetCtx(conn.ctx, key, val)
}

/**
* HGet
* @params key string
* @return map[string]string, error
**/
func HGet(key string) (map[string]string, error) {
	if conn == nil {
		return map[string]string{}, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	return HGetCtx(conn.ctx, key)
}

/**
* HSetAtrib
* @params key string
* @params atr string
* @params val string
* @return error
**/
func HSetAtrib(key, atr, val string) error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	return HSetCtx(conn.ctx, key, map[string]string{atr: val})
}

/**
* HGetAtrib
* @params key string
* @params atr string
* @return string, error
**/
func HGetAtrib(key, atr string) (string, error) {
	if conn == nil {
		return "", logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	atribs, err := HGetCtx(conn.ctx, key)
	if err != nil {
		return "", err
	}

	for k, v := range atribs {
		if k == atr {
			return v, nil
		}
	}

	return "", nil
}

/**
* HDel
* @params key string
* @params atr string
* @return error
**/
func HDel(key, atr string) error {
	if conn == nil {
		return logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	return HDelCtx(conn.ctx, key, atr)
}

/**
* SetVerify
* @params device string
* @params key string
* @params val string
* @return error
**/
func SetVerify(device, key, val string) error {
	key = strs.Format(`verify:%s:%s`, device, key)
	return Set(key, val, 5*60)
}

/**
* GetVerify
* @params device string
* @params key string
* @return string, error
**/
func GetVerify(device string, key string) (string, error) {
	key = strs.Format(`verify:%s:%s`, device, key)
	return Get(key, "")
}

/**
* DelVerify
* @params device string
* @params key string
* @return int64, error
**/
func DelVerify(device string, key string) (int64, error) {
	key = strs.Format(`verify:%s:%s`, device, key)
	return Del(key)
}

/**
* AllCache
* @params device string
* @params key string
* @params val string
* @return error
**/
func AllCache(search string, page, rows int) (et.List, error) {
	if conn == nil {
		return et.List{}, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	ctx := context.Background()
	var cursor uint64
	var count int64
	var items et.Items = et.Items{}
	offset := (page - 1) * rows
	cursor = uint64(offset)
	count = int64(rows)

	iter := conn.db.Scan(ctx, cursor, search, count).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		items.Result = append(items.Result, et.Json{"key": key})
		items.Count++
	}

	return items.ToList(items.Count, page, rows), nil
}

/**
* GetJson
* @params key string
* @return Json, error
**/
func GetJson(key string) (et.Json, error) {
	if conn == nil {
		return et.Json{}, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	_default := "{}"
	val, err := Get(key, _default)
	if err != nil {
		return et.Json{}, err
	}

	if val == _default {
		return et.Json{}, nil
	}

	var result et.Json
	err = result.Scan(val)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* GetItem
* @params key string
* @return Item, error
**/
func GetItem(key string) (et.Item, error) {
	if conn == nil {
		return et.Item{}, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	_default := "{}"
	val, err := Get(key, _default)
	if err != nil {
		return et.Item{}, err
	}

	if val == _default {
		return et.Item{}, nil
	}

	var result et.Json
	err = result.Scan(val)
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok:     true,
		Result: result,
	}, nil
}

/**
* GetItems
* @params key string
* @return Items, error
**/
func GetItems(key string) (et.Items, error) {
	if conn == nil {
		return et.Items{}, logs.Log(msg.ERR_NOT_CACHE_SERVICE)
	}

	_default := "[]"
	val, err := Get(key, _default)
	if err != nil {
		return et.Items{}, err
	}

	if val == _default {
		return et.Items{}, nil
	}

	var result et.Items
	err = result.Scan(val)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}
