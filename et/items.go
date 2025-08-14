package et

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
)

type Items struct {
	Ok     bool   `json:"ok"`
	Count  int    `json:"count"`
	Result []Json `json:"result"`
}

func (it *Items) Scan(src interface{}) error {
	var ba []byte
	switch v := src.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return logs.Errorf(`json/Scan - Failed to unmarshal JSON value:%s`, src)
	}

	var t []Json
	err := json.Unmarshal(ba, &t)
	if err != nil {
		return err
	}

	*it = Items{
		Ok:     len(t) > 0,
		Count:  len(t),
		Result: t,
	}

	return nil
}

func (it *Items) First() Item {
	if !it.Ok {
		return Item{}
	}

	return Item{
		Ok:     true,
		Result: it.Result[0],
	}
}

func (it *Items) Add(item Json) *Items {
	it.Result = append(it.Result, item)
	it.Count = len(it.Result)
	it.Ok = it.Count > 0

	return it
}

func (s *Items) AddMany(items []Json) {
	(*s).Result = append((*s).Result, items...)
	(*s).Count = len((*s).Result)
	(*s).Ok = (*s).Count > 0
}

func (it *Items) ValAny(idx int, _default any, atribs ...string) any {
	item := it.Result[idx]
	if item == nil {
		return _default
	}

	return item.ValAny(_default, atribs...)
}

func (it *Items) ValStr(idx int, _default string, atribs ...string) string {
	item := it.Result[idx]
	if item == nil {
		return _default
	}

	return item.ValStr(_default, atribs...)
}

func (it *Items) Uppcase(idx int, _default string, atribs ...string) string {
	item := it.Result[idx]
	if item == nil {
		return _default
	}

	result := Val(item, _default, atribs...)

	switch v := result.(type) {
	case string:
		return strings.ToUpper(v)
	default:
		return strs.Format(`%v`, strings.ToUpper(_default))
	}
}

func (it *Items) Lowcase(idx int, _default string, atribs ...string) string {
	item := it.Result[idx]
	if item == nil {
		return _default
	}

	result := Val(item, _default, atribs...)

	switch v := result.(type) {
	case string:
		return strings.ToLower(v)
	default:
		return strs.Format(`%v`, strings.ToLower(_default))
	}
}

func (it *Items) Titlecase(idx int, _default string, atribs ...string) string {
	item := it.Result[idx]
	if item == nil {
		return _default
	}

	result := Val(it.Result[idx], _default, atribs...)

	switch v := result.(type) {
	case string:
		return strings.ToTitle(v)
	default:
		return strs.Format(`%v`, strings.ToTitle(_default))
	}
}

func (it *Items) Get(idx int, key string) interface{} {
	item := it.Result[idx]
	if item == nil {
		return nil
	}

	return it.Result[idx].Get(key)
}

func (it *Items) Set(idx int, key string, val interface{}) bool {
	item := it.Result[idx]
	if item == nil {
		return false
	}

	return item.Set(key, val)
}

func (it *Items) Del(idx int, key string) bool {
	item := it.Result[idx]
	if item == nil {
		return false
	}

	return item.Del(key)
}

func (it *Items) Id(idx int) string {
	item := it.Result[idx]
	if item == nil {
		return ""
	}

	return item.Id()
}

func (it *Items) IdT(idx int) string {
	item := it.Result[idx]
	if item == nil {
		return ""
	}

	return item.IdT()
}

func (it *Items) Key(idx int, atribs ...string) string {
	item := it.Result[idx]
	if item == nil {
		return ""
	}

	return item.Key(atribs...)
}

func (it *Items) Str(idx int, atribs ...string) string {
	item := it.Result[idx]
	if item == nil {
		return ""
	}

	return item.Str(atribs...)
}

func (it *Items) Int(idx int, atribs ...string) int {
	item := it.Result[idx]
	if item == nil {
		return 0
	}

	return item.Int(atribs...)
}

func (it *Items) Num(idx int, atribs ...string) float64 {
	item := it.Result[idx]
	if item == nil {
		return 0
	}

	return item.Num(atribs...)
}

func (it *Items) Bool(idx int, atribs ...string) bool {
	item := it.Result[idx]
	if item == nil {
		return false
	}

	return item.Bool(atribs...)
}

func (it *Items) Json(idx int, atribs ...string) Json {
	item := it.Result[idx]
	if item == nil {
		return Json{}
	}

	val := Val(item, Json{}, atribs...)

	switch v := val.(type) {
	case Json:
		return Json(v)
	case map[string]interface{}:
		return Json(v)
	default:
		logs.Errorf("Not Items.Json type (%v) value:%v", reflect.TypeOf(v), v)
		return Json{}
	}
}

func (it *Items) ToByte() []byte {
	return Json{
		"Ok":     it.Ok,
		"Count":  it.Count,
		"Result": it.Result,
	}.ToByte()
}

func (it *Items) ToString() string {
	return ArrayToString(it.Result)
}

func (it *Items) ToJson() Json {
	return Json{
		"Ok":     it.Ok,
		"Count":  it.Count,
		"Result": it.Result,
	}
}

func (it *Items) ToList(all, page, rows int) List {
	var start int
	var end int
	count := it.Count

	if count <= 0 {
		start = 0
		end = 0
	} else {
		offset := (page - 1) * rows

		if offset > 0 {
			start = offset + 1
			end = offset + count
		} else {
			start = 1
			end = count
		}
	}

	return List{
		Rows:   rows,
		All:    all,
		Count:  count,
		Page:   page,
		Start:  start,
		End:    end,
		Result: it.Result,
	}
}
