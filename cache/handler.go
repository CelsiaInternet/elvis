package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	js "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/msg"
)

func SetCtx(ctx context.Context, key, val string, second time.Duration) error {
	if conn == nil {
		return logs.Errorm(msg.ERR_NOT_CACHE_SERVICE)
	}

	duration := second * time.Second

	err := conn.db.Set(ctx, key, val, duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetCtx(ctx context.Context, key string) (string, error) {
	if conn == nil {
		return "", logs.Errorm(msg.ERR_NOT_CACHE_SERVICE)
	}

	result, err := conn.db.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func DelCtx(ctx context.Context, key string) (int64, error) {
	if conn == nil {
		return 0, logs.Errorm(msg.ERR_NOT_CACHE_SERVICE)
	}

	intCmd := conn.db.Del(ctx, key)

	return intCmd.Val(), intCmd.Err()
}

func HSetCtx(ctx context.Context, key string, val map[string]string) error {
	if conn == nil {
		return logs.Errorm(msg.ERR_NOT_CACHE_SERVICE)
	}

	for k, v := range val {
		err := conn.db.HSet(ctx, key, k, v).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func HGetCtx(ctx context.Context, key string) (map[string]string, error) {
	if conn == nil {
		return map[string]string{}, logs.Errorm(msg.ERR_NOT_CACHE_SERVICE)
	}

	result := conn.db.HGetAll(ctx, key).Val()

	return result, nil
}

func HDelCtx(ctx context.Context, key, atr string) error {
	if conn == nil {
		return logs.Errorm(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.db.Do(ctx, "HDEL", key, atr).Err()
	if err != nil {
		return err
	}

	return nil
}

func Set(key, val string, second time.Duration) error {
	return SetCtx(conn.ctx, key, val, second)
}

func Get(key string) (string, error) {
	return GetCtx(conn.ctx, key)
}

func Del(key string) (int64, error) {
	return DelCtx(conn.ctx, key)
}

func Empty() error {
	ctx := context.Background()
	iter := conn.db.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		DelCtx(ctx, key)
	}

	return nil
}

func More(key string, second time.Duration) int {
	n, err := Get(key)
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
	Set(key, fmt.Sprintf(`%d`, val), second)

	return val
}

func HSet(key string, val map[string]string) error {
	return HSetCtx(conn.ctx, key, val)
}

func HGet(key string) (map[string]string, error) {
	return HGetCtx(conn.ctx, key)
}

func HDel(key, atr string) error {
	return HDelCtx(conn.ctx, key, atr)
}

func SetVerify(device, key, val string) error {
	key = fmt.Sprintf(`verify:%s/%s`, device, key)
	return Set(key, val, 5*60)
}

func GetVerify(device string, key string) (string, error) {
	key = fmt.Sprintf(`verify:%s/%s`, device, key)
	return Get(key)
}

func DelVerify(device string, key string) (int64, error) {
	key = fmt.Sprintf(`verify:%s/%s`, device, key)
	return Del(key)
}

func AllCache(search string, page, rows int) (js.List, error) {
	ctx := context.Background()
	var cursor uint64
	var count int64
	var items js.Items = js.Items{}
	offset := (page - 1) * rows
	cursor = uint64(offset)
	count = int64(rows)

	iter := conn.db.Scan(ctx, cursor, search, count).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		items.Result = append(items.Result, js.Json{"key": key})
		items.Count++
	}

	return items.ToList(items.Count, page, rows), nil
}
