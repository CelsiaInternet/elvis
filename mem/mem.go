package mem

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

type Mem struct {
	mu    sync.RWMutex
	items map[string]*Item
}

var conn *Mem

func Load() (*Mem, error) {
	result := &Mem{
		items: make(map[string]*Item),
	}

	logs.Logf("Mem", "Load memory cache")

	go result.sweeper()

	return result, nil
}

func init() {
	if conn != nil {
		return
	}

	var err error
	conn, err = Load()
	if err != nil {
		logs.Alert(err)
		return
	}
}

// sweeper runs a single goroutine that expires items instead of spawning one
// goroutine per Set call. It wakes every second, collects expired keys under a
// read lock, then removes them under a write lock with a fresh expiry check to
// avoid evicting a key that was just refreshed between the two lock windows.
func (c *Mem) sweeper() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		var expired []string

		c.mu.RLock()
		for key, item := range c.items {
			if !item.expiry.IsZero() && now.After(item.expiry) {
				expired = append(expired, key)
			}
		}
		c.mu.RUnlock()

		if len(expired) == 0 {
			continue
		}

		c.mu.Lock()
		for _, key := range expired {
			if item, ok := c.items[key]; ok {
				if !item.expiry.IsZero() && time.Now().After(item.expiry) {
					delete(c.items, key)
				}
			}
		}
		c.mu.Unlock()
	}
}

/**
* Type
* @return string
**/
func (c *Mem) Type() string {
	return "mem"
}

/**
* Set
* @param key string
* @param value string
* @param expiration time.Duration
* @return string
**/
func (c *Mem) Set(key string, value string, expiration time.Duration) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		item.Set(value)
	} else {
		item = New(key, value)
		c.items[key] = item
	}

	if expiration > 0 {
		item.expiry = time.Now().Add(expiration * time.Second)
	} else {
		item.expiry = time.Time{}
	}

	return value
}

/**
* Get
* @param key, def string
* @return string, error
**/
func (c *Mem) Get(key, def string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.items[key]; ok {
		return item.Str(), nil
	}

	return def, utility.NewError("IsNil")
}

/**
* Del
* @param key string
* @return bool
**/
func (c *Mem) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; !ok {
		return false
	}

	delete(c.items, key)

	return true
}

/**
* More
* @param key string
* @param expiration time.Duration
* @return int64
**/
func (c *Mem) More(key string, expiration time.Duration) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		newItem := New(key, "0")
		if expiration > 0 {
			newItem.expiry = time.Now().Add(expiration * time.Second)
		}
		c.items[key] = newItem
		return 0
	}

	result := item.Int64() + 1
	item.Value = strconv.FormatInt(result, 10)
	item.Dateupdate = time.Now()
	if expiration > 0 {
		item.expiry = time.Now().Add(expiration * time.Second)
	}

	return result
}

/**
* Clear removes all keys whose name contains the match substring.
* The regexp is compiled once before the loop instead of per-key.
* @param match string
**/
func (c *Mem) Clear(match string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if match == "" {
		c.items = make(map[string]*Item)
		return
	}

	pattern := fmt.Sprintf(".*%s.*", regexp.QuoteMeta(match))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	for key := range c.items {
		if re.MatchString(key) {
			delete(c.items, key)
		}
	}
}

func (c *Mem) Empty() {
	c.Clear("")
}

/**
* Len
* @return int
**/
func (c *Mem) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

/**
* Keys
* @return []string
**/
func (c *Mem) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))

	for key := range c.items {
		keys = append(keys, key)
	}

	return keys
}

/**
* Values
* @return []string
**/
func (c *Mem) Values() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	values := make([]string, 0, len(c.items))

	for _, item := range c.items {
		values = append(values, item.Str())
	}

	return values
}
