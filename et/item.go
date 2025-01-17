package et

import (
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
)

type Item struct {
	Ok     bool `json:"ok"`
	Result Json `json:"result"`
}

func (it *Item) ScanRows(rows *sql.Rows) error {
	it.Ok = true
	it.Result = make(Json)
	it.Result.ScanRows(rows)

	return nil
}

func (it *Item) ValAny(_default any, atribs ...string) any {
	return Val(it.Result, _default, atribs...)
}

func (it *Item) ValStr(_default string, atribs ...string) string {
	return it.Result.ValStr(_default, atribs...)
}

func (it *Item) ValInt(_default int, atribs ...string) int {
	return it.Result.ValInt(_default, atribs...)
}

func (it *Item) ValNum(_default float64, atribs ...string) float64 {
	return it.Result.ValNum(_default, atribs...)
}

func (it *Item) ValBool(_default bool, atribs ...string) bool {
	return it.Result.ValBool(_default, atribs...)
}

func (it *Item) ValTime(_default time.Time, atribs ...string) time.Time {
	return it.Result.ValTime(_default, atribs...)
}

func (it *Item) ValJson(_default Json, atribs ...string) Json {
	return it.Result.ValJson(_default, atribs...)
}

func (it *Item) Uppcase(_default string, atribs ...string) string {
	result := Val(it.Result, _default, atribs...)

	switch v := result.(type) {
	case string:
		return strings.ToUpper(v)
	default:
		return strs.Format(`%v`, strings.ToUpper(_default))
	}
}

func (it *Item) Lowcase(_default string, atribs ...string) string {
	result := Val(it.Result, _default, atribs...)

	switch v := result.(type) {
	case string:
		return strings.ToLower(v)
	default:
		return strs.Format(`%v`, strings.ToLower(_default))
	}
}

func (it *Item) Titlecase(_default string, atribs ...string) string {
	result := Val(it.Result, _default, atribs...)

	switch v := result.(type) {
	case string:
		return strings.ToTitle(v)
	default:
		return strs.Format(`%v`, strings.ToTitle(_default))
	}
}

func (it *Item) Get(key string) interface{} {
	return it.Result.Get(key)
}

func (it *Item) Set(key string, val any) bool {
	return it.Result.Set(key, val)
}

func (it *Item) Del(key string) bool {
	return it.Result.Del(key)
}

func (it *Item) IsDiferent(new Json) bool {
	return IsDiferent(it.Result, new)
}

func (it *Item) IsChange(new Json) bool {
	return IsChange(it.Result, new)
}

func (it *Item) Any(_default any, atribs ...string) *Any {
	return it.Result.Any(_default, atribs...)
}

func (it *Item) Id() string {
	return it.Result.Id()
}

func (it *Item) IdT() string {
	return it.Result.IdT()
}

func (it *Item) State() string {
	return it.Result.State()
}

func (it *Item) Index() int {
	return it.Result.Index()
}

func (it *Item) Index64() int64 {
	return it.Result.Index64()
}

func (it *Item) Key(atribs ...string) string {
	return it.Result.Key(atribs...)
}

func (it *Item) Str(atribs ...string) string {
	return it.Result.Str(atribs...)
}

func (it *Item) Int(atribs ...string) int {
	return it.Result.Int(atribs...)
}

func (it *Item) Int64(atribs ...string) int64 {
	return it.Result.Int64(atribs...)
}

func (it *Item) Num(atribs ...string) float64 {
	return it.Result.Num(atribs...)
}

func (it *Item) Bool(atribs ...string) bool {
	return it.Result.Bool(atribs...)
}

func (it *Item) Time(atribs ...string) time.Time {
	return it.Result.Time(atribs...)
}

func (it *Item) Data(atribs ...string) JsonD {
	return it.Result.Data(atribs...)
}

func (it *Item) Json(atribs ...string) Json {
	val := Val(it.Result, Json{}, atribs...)

	switch v := val.(type) {
	case Json:
		return Json(v)
	case map[string]interface{}:
		return Json(v)
	default:
		logs.Errorf("Not Item.Json type (%v) value:%v", reflect.TypeOf(v), v)
		return Json{}
	}
}

func (it *Item) Array(atrib string) []Json {
	return it.Result.Array(atrib)
}

func (it *Item) ArrayStr(atrib string) []string {
	return it.Result.ArrayStr(atrib)
}

func (it *Item) ToString() string {
	return it.Result.ToString()
}

func (it *Item) ToJson() Json {
	return Json{
		"Ok":     it.Ok,
		"Result": it.Result,
	}
}

func (it *Item) ToByte() []byte {
	return Json{
		"Ok":     it.Ok,
		"Result": it.Result,
	}.ToByte()
}

func (it *Item) Consolidate(toField string, ruleOut ...string) Json {
	result := it.Result.Consolidate(toField, ruleOut...)

	return result
}

func (it *Item) ConsolidateAndUpdate(toField string, ruleOut []string, new Json) (Json, error) {
	return it.Result.ConsolidateAndUpdate(toField, ruleOut, new)
}
