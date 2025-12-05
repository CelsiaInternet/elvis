package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/redis/go-redis/v9"
)

var IncrSetTTLScript = redis.NewScript(`
    local ttl = tonumber(ARGV[1])

    local newVal = redis.call("INCR", KEYS[1])

    if ttl > 0 then
        redis.call("PEXPIRE", KEYS[1], ttl)
    end

    return newVal
`)

var DecrKeepTTLScript = redis.NewScript(`
    local val = redis.call("GET", KEYS[1])
    if not val then
        return -1
    end

    val = tonumber(val)

    if val > 0 then
        return redis.call("DECR", KEYS[1])
    else
        return -1
    end
`)

/**
* SetCtx
* @params ctx context.Context, key string, val string, second time.Duration
* @return error
**/
func SetCtx(ctx context.Context, key, val string, second time.Duration) error {
	if conn == nil {
		return fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.Set(ctx, key, val, second).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* ExpireCtx
* @params ctx context.Context, key string, second time.Duration
* @return error
**/
func ExpireCtx(ctx context.Context, key string, second time.Duration) error {
	if conn == nil {
		return fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	return conn.Expire(ctx, key, second).Err()
}

/**
* GetCtx
* @params ctx context.Context, key string
* @params def string
* @return string, error
**/
func GetCtx(ctx context.Context, key, def string) (string, error) {
	if conn == nil {
		return def, fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	result, err := conn.Get(ctx, key).Result()
	if err == redis.Nil {
		return def, nil
	} else if err != nil {
		return def, err
	}

	return result, nil
}

/**
* ExistsCtx
* @params ctx context.Context, key string
* @return bool
**/
func ExistsCtx(ctx context.Context, key string) bool {
	if conn == nil {
		logs.Alertm(msg.ERR_NOT_CACHE_SERVICE)
		return false
	}

	result, err := conn.Exists(ctx, key).Result()
	if err != nil {
		logs.Alertm(msg.ERR_NOT_CACHE_SERVICE)
		return false
	}

	return result == 1
}

/**
* DeleteCtx
* @params ctx context.Context, key string
* @return int64, error
**/
func DeleteCtx(ctx context.Context, key string) (int64, error) {
	if conn == nil {
		return 0, fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	intCmd := conn.Del(ctx, key)

	return intCmd.Val(), intCmd.Err()
}

/**
* IncrCtx
* @params ctx context.Context, key string, second time.Duration
* @return int64
**/
func IncrCtx(ctx context.Context, key string, second time.Duration) int64 {
	if conn == nil {
		return 0
	}

	result, err := IncrSetTTLScript.Run(
		ctx,
		conn,
		[]string{key},
		second.Milliseconds(),
	).Int64()
	if err != nil {
		return 0
	}

	return result
}

/**
* DecrCtx
* @params ctx context.Context, key string
* @return int64
**/
func DecrCtx(ctx context.Context, key string) int64 {
	if conn == nil {
		return 0
	}

	result, err := DecrKeepTTLScript.Run(
		ctx,
		conn,
		[]string{key},
	).Int64()
	if err != nil {
		return 0
	}

	return result
}

/**
* LPushCtx
* @params ctx context.Context, key string, val string
* @return error
**/
func LPushCtx(ctx context.Context, key string, val string) error {
	if conn == nil {
		return fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.RPush(ctx, key, val).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* LPushCtx
* @params ctx context.Context, key string, val string
* @return error
**/
func LRemCtx(ctx context.Context, key string, val string) error {
	if conn == nil {
		return fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.LRem(ctx, key, 1, val).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* LRangeCtx
* @params ctx context.Context, key string, start int64, stop int64
* @return []string, error
**/
func LRangeCtx(ctx context.Context, key string, start int64, stop int64) ([]string, error) {
	if conn == nil {
		return []string{}, fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	result, err := conn.LRange(ctx, key, start, stop).Result()

	return result, err
}

/**
* LTrimCtx
* @params ctx context.Context, key string, start int64, stop int64
* @return error
**/
func LTrimCtx(ctx context.Context, key string, start int64, stop int64) error {
	if conn == nil {
		return fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.LTrim(ctx, key, start, stop).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* HSetCtx
* @params ctx context.Context, key string, val map[string]string
* @return error
**/
func HSetCtx(ctx context.Context, key string, val map[string]string) error {
	if conn == nil {
		return fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.HSet(ctx, key, val).Err()
	if err != nil {
		return err
	}

	return nil
}

/**
* HGetCtx
* @params ctx context.Context, key string
* @return map[string]string, error
**/
func HGetCtx(ctx context.Context, key string) (map[string]string, error) {
	if conn == nil {
		return map[string]string{}, fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	result := conn.HGetAll(ctx, key).Val()

	return result, nil
}

/**
* HDeleteCtx
* @params ctx context.Context, key string
* @params atr string
* @return error
**/
func HDeleteCtx(ctx context.Context, key, atr string) error {
	if conn == nil {
		return fmt.Errorf(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.Do(ctx, "HDEL", key, atr).Err()
	if err != nil {
		return err
	}

	return nil
}
