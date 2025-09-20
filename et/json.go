package et

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/timezone"
)

const TpObject = 1
const TpArray = 2

type JsonD struct {
	Type  int
	Value interface{}
}

type Json map[string]interface{}

func JsonToArrayJson(src map[string]interface{}) ([]Json, error) {
	result := []Json{}
	result = append(result, src)

	return result, nil
}

func Object(src interface{}) (Json, error) {
	var bt []byte
	var err error
	switch v := src.(type) {
	case string:
		bt = []byte(v)
	default:
		bt, err = json.Marshal(v)
		if err != nil {
			return Json{}, err
		}
	}

	result := Json{}
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return Json{}, err
	}

	return result, nil
}

func Array(src interface{}) ([]Json, error) {
	var bt []byte
	var err error
	switch v := src.(type) {
	case string:
		bt = []byte(v)
	default:
		bt, err = json.Marshal(src)
		if err != nil {
			return []Json{}, err
		}
	}

	result := []Json{}
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return []Json{}, err
	}

	return result, nil
}

func (s *Json) Scan(src interface{}) error {
	var ba []byte
	switch v := src.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return logs.Errorf("json.Scan", `Failed to unmarshal JSON type:%s`, reflect.TypeOf(v))
	}

	t := map[string]interface{}{}
	err := json.Unmarshal(ba, &t)
	if err != nil {
		return err
	}

	*s = Json(t)

	return nil
}

func (s *Json) ScanRows(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(cols))
	pointers := make([]interface{}, len(cols))
	for i := range values {
		pointers[i] = &values[i]
	}

	if err := rows.Scan(pointers...); err != nil {
		return err
	}

	result := make(Json)
	for i, col := range cols {
		src := values[i]
		switch v := src.(type) {
		case nil:
			result[col] = nil
		case []byte:
			var bt interface{}
			err = json.Unmarshal(v, &bt)
			if err == nil {
				result[col] = bt
				continue
			}
			result[col] = src
			logs.Debugf(`[]byte Col:%s Type:%v Value:%v`, col, reflect.TypeOf(v), v)
		default:
			result[col] = src
		}
	}

	*s = result

	return nil
}

func (s Json) IsEmpty() bool {
	return len(s) == 0
}

func (s Json) ToByte() []byte {
	result, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	return result
}

func (s Json) ToString() string {
	bt, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	result := string(bt)

	return result
}

func (s Json) ToEscapeHTML() string {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(s)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (s Json) ToUnquote() string {
	str := s.ToString()
	result := strs.Format(`'%v'`, str)

	return result
}

func (s Json) ToQuote() string {
	for k, v := range s {
		if str, ok := s["mensaje"].(string); ok {
			ustr, err := strconv.Unquote(`"` + str + `"`)
			if err != nil {
				s[k] = v
			} else {
				s[k] = ustr
			}
		} else {
			s[k] = v
		}
	}
	str := s.ToString()

	return str
}

func (s Json) ToItem(src interface{}) Item {
	s.Scan(src)
	return Item{
		Ok:     s.Bool("Ok"),
		Result: s.Json("Result"),
	}
}

func (s Json) ValAny(_default any, atribs ...string) any {
	return Val(s, _default, atribs...)
}

func (s Json) ValStr(_default string, atribs ...string) string {
	val := s.ValAny(_default, atribs...)

	switch v := val.(type) {
	case string:
		return v
	default:
		return strs.Format(`%v`, v)
	}
}

func (s Json) ValInt(_default int, atribs ...string) int {
	val := s.ValAny(_default, atribs...)

	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return _default
		}
		return i
	default:
		return _default
	}
}

func (s Json) ValInt64(_default int64, atribs ...string) int64 {
	val := s.ValAny(_default, atribs...)

	switch v := val.(type) {
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return _default
		}
		return i
	default:
		return _default
	}
}

func (s Json) ValNum(_default float64, atribs ...string) float64 {
	val := s.ValAny(_default, atribs...)

	switch v := val.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	case float32:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Println("ValNum value float not conver", reflect.TypeOf(v), v)
			return _default
		}
		return i
	default:
		log.Println("ValNum value is not float, type:", reflect.TypeOf(v), "value:", v)
		return _default
	}
}

func (s Json) ValBool(_default bool, atribs ...string) bool {
	val := s.ValAny(_default, atribs...)

	switch v := val.(type) {
	case bool:
		return v
	case int:
		return v == 1
	case string:
		v = strings.ToLower(v)
		switch v {
		case "true":
			return true
		case "false":
			return false
		default:
			log.Println("ValBool value is not bool, type:", reflect.TypeOf(v), "value:", v)
			return _default
		}
	default:
		log.Println("ValBool value is not bool, type:", reflect.TypeOf(v), "value:", v)
		return _default
	}
}

func (s Json) ValTime(_default time.Time, atribs ...string) time.Time {
	val := s.ValAny(_default, atribs...)

	switch v := val.(type) {
	case int:
		return _default
	case string:
		layout := "2006-01-02T15:04:05.000Z"
		result, err := time.Parse(layout, v)
		if err != nil {
			return _default
		}
		return result
	case time.Time:
		return v
	default:
		log.Println("ValTime value is not time, type:", reflect.TypeOf(v), "value:", v)
		return _default
	}
}

func (s Json) ValJson(_default Json, atribs ...string) Json {
	val := s.ValAny(_default, atribs...)

	switch v := val.(type) {
	case Json:
		return v
	default:
		log.Println("ValTime value is not json, type:", reflect.TypeOf(v), "value:", v)
		return _default
	}
}

func (s Json) ValArray(defaultVal []interface{}, atribs ...string) []interface{} {
	var result []interface{}
	val := s.ValAny(defaultVal, atribs...)

	switch v := val.(type) {
	case []interface{}:
		return v
	case []Json:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	case []map[string]interface{}:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	case []string:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	case []int:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	case []int64:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	case []float64:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	case []float32:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	case []bool:
		for _, item := range v {
			result = append(result, item)
		}

		return result
	default:
		src := fmt.Sprintf(`%v`, v)
		err := json.Unmarshal([]byte(src), &result)
		if err != nil {
			err := fmt.Errorf(`valor: %v error:%v type:%T`, val, err.Error(), val)
			logs.Alert(err)
			return defaultVal
		}

		return result
	}
}

func (s Json) Any(_default any, atribs ...string) *Any {
	result := Val(s, _default, atribs...)
	return NewAny(result)
}

func (s Json) Id() string {
	return s.ValStr("-1", "_id")
}

func (s Json) IdT() string {
	return s.ValStr("-1", "_idT")
}

func (s Json) State() string {
	return s.ValStr("-1", "_state")
}

func (s Json) Index() int {
	return s.ValInt(-1, "index")
}

func (s Json) Index64() int64 {
	return s.ValInt64(-1, "index")
}

func (s Json) Key(atribs ...string) string {
	return s.ValStr("-1", atribs...)
}

func (s Json) Str(atribs ...string) string {
	return s.ValStr("", atribs...)
}

func (s Json) Int(atribs ...string) int {
	return s.ValInt(0, atribs...)
}

func (s Json) Int64(atribs ...string) int64 {
	return s.ValInt64(0, atribs...)
}

func (s Json) Num(atribs ...string) float64 {
	return s.ValNum(0.00, atribs...)
}

func (s Json) Bool(atribs ...string) bool {
	return s.ValBool(false, atribs...)
}

func (s Json) Byte(atribs ...string) ([]byte, error) {
	value := s.ValAny("", atribs...)
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (s Json) Time(atribs ...string) time.Time {
	_default := timezone.NowTime()
	return s.ValTime(_default, atribs...)
}

func (s Json) Duration(atribs ...string) (time.Duration, error) {
	str := s.ValStr("", atribs...)
	timDuration, err := strs.StrToTime(str)
	if err != nil {
		return 0, err
	}

	return timDuration.Sub(timezone.NowTime()), nil
}

func (s Json) Data(atrib ...string) JsonD {
	val := Val(s, nil, atrib...)
	if val == nil {
		return JsonD{
			Type:  TpObject,
			Value: Json{},
		}
	}

	switch v := val.(type) {
	case Json:
		return JsonD{
			Type:  TpObject,
			Value: v,
		}
	case map[string]interface{}:
		return JsonD{
			Type:  TpObject,
			Value: Json(v),
		}
	case []Json:
		return JsonD{
			Type:  TpArray,
			Value: v,
		}
	case []interface{}:
		return JsonD{
			Type:  TpArray,
			Value: v,
		}
	default:
		logs.Errorf("json.Data", "Atrib:%s Type:%v Value:%v", atrib, reflect.TypeOf(v), v)
		return JsonD{
			Type:  TpObject,
			Value: Json{},
		}
	}
}

func (s Json) Json(atrib string) Json {
	val := Val(s, nil, atrib)
	if val == nil {
		return Json{}
	}

	switch v := val.(type) {
	case Json:
		return Json(v)
	case map[string]interface{}:
		return Json(v)
	case []interface{}:
		result := Json{
			atrib: v,
		}

		return result
	default:
		logs.Errorf("json/Json - Atrib:%s Type:%v Value:%v", atrib, reflect.TypeOf(v), v)
		return Json{}
	}
}

/**
* Array
* @param atrib ...string
* @return []interface{}
**/
func (s Json) Array(atrib ...string) []interface{} {
	return s.ValArray([]interface{}{}, atrib...)
}

/**
* ArrayJson
* @param atrib string
* @return []Json
**/
func (s Json) ArrayJson(atrib string) []Json {
	val := Val(s, nil, atrib)
	if val == nil {
		return []Json{}
	}

	data, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return []Json{}
	}

	var result []Json
	err = json.Unmarshal(data, &result)
	if err != nil {
		return []Json{}
	}

	return result
}

func (s Json) ArrayStr(atrib string) []string {
	val := Val(s, nil, atrib)
	if val == nil {
		return []string{}
	}

	data, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return []string{}
	}

	var result []string
	err = json.Unmarshal(data, &result)
	if err != nil {
		return []string{}
	}

	return result
}

func (s Json) Update(fromJson Json) error {
	var result bool = false
	for k, new := range fromJson {
		v := s[k]

		if v == nil {
			s[k] = new
		} else if new != nil {
			if !result && reflect.DeepEqual(v, new) {
				result = true
			}
			s[k] = new
		}
	}

	return nil
}

func (s Json) IsDiferent(new Json) bool {
	return IsDiferent(s, new)
}

func (s Json) IsChanged(from Json) bool {
	for key, fromValue := range from {
		if s[key] == nil {
			return true
		}

		if strings.EqualFold(fmt.Sprintf(`%v`, s[key]), fmt.Sprintf(`%v`, fromValue)) {
			return true
		}
	}

	return false
}

/**
* Get
* @param key string
* @return interface{}
**/
func (s Json) Get(key string) interface{} {
	v, ok := s[key]
	if !ok {
		return nil
	}

	return v
}

func (s Json) Set(key string, val interface{}) bool {
	key = strs.Lowcase(key)
	if s[key] != nil {
		s[key] = val
		return true
	}

	s[key] = val

	return false
}

func (s *Json) Append(obj Json) *Json {
	var result Json = *s
	for k, v := range obj {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}

	return &result
}

func (s Json) Del(key string) bool {
	if _, ok := s[key]; !ok {
		return false
	}

	delete(s, key)
	return true
}

func (s Json) ExistKey(key string) bool {
	return s[key] != nil
}

func (s Json) Consolidate(toField string, ruleOut ...string) Json {
	FindIndex := func(arr []string, valor string) int {
		for i, v := range arr {
			if v == valor {
				return i
			}
		}
		return -1
	}

	result := s
	if s.ExistKey(toField) {
		result = s.Json(toField)
	}

	for k, v := range s {
		if k != toField {
			idx := FindIndex(ruleOut, k)
			if idx == -1 {
				result[k] = v
			}
		}
	}

	return result
}

func (s Json) ConsolidateAndUpdate(toField string, ruleOut []string, new Json) (Json, error) {
	result := s.Consolidate(toField, ruleOut...)
	err := result.Update(new)
	if err != nil {
		return Json{}, nil
	}

	return result, nil
}

func (s Json) Clone() Json {
	result := Json{}
	for k, v := range s {
		result[k] = v
	}

	return result
}
