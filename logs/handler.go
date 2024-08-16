package logs

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/msg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var IsNil = mongo.ErrNoDocuments

// Estructura que representa el documento almacenado en MongoDB
type MongoDocument struct {
	Key        string            `bson:"key"`
	Value      string            `bson:"value"`
	Attributes map[string]string `bson:"attributes,omitempty"`
}

func Log(kind string, args ...any) error {
	console.Printl(kind, "", args...)
	return nil
}

func Logf(kind string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	console.Printl(kind, "", message)
}

func Traces(kind, color string, err error) ([]string, error) {
	var n int = 1
	var traces []string = []string{err.Error()}

	console.Printl(kind, color, err.Error())

	for {
		pc, file, line, more := runtime.Caller(n)
		if !more {
			break
		}
		n++
		function := runtime.FuncForPC(pc)
		name := function.Name()
		list := strings.Split(name, ".")
		if len(list) > 0 {
			name = list[len(list)-1]
		}
		if !slices.Contains([]string{"ErrorM", "ErrorF"}, name) {
			trace := fmt.Sprintf("%s:%d func:%s", file, line, name)
			traces = append(traces, trace)
			console.Printl("TRACE", color, trace)
		}
	}

	return traces, err
}

func Alert(err error) error {
	if err != nil {
		console.Printl("Alert", "Yellow", err.Error())
	}

	return err
}

func Alertm(message string) error {
	return Alert(errors.New(message))
}

func Alertf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)

	return Alertm(message)
}

func Error(err error) error {
	_, err = Traces("Error", "red", err)

	return err
}

func Errorm(message string) error {
	err := errors.New(message)
	return Error(err)
}

func Errorf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	err := errors.New(message)
	return Error(err)
}

func Info(v ...any) {
	console.Printl("Info", "Blue", v...)
}

func Infof(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	console.Printl("Info", "Blue", message)
}

func Fatal(v ...any) {
	console.Printl("Fatal", "Red", v...)
	os.Exit(1)
}

func Panic(v ...any) {
	console.Printl("Panic", "Red", v...)
	os.Exit(1)
}

func Ping() {
	console.Printl("PING", "")
}

func Pong() {
	console.Printl("PONG", "")
}

func Debug(v ...any) {
	console.Printl("Debug", "Cyan", v...)
}

func Debugf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	console.Printl("Debug", "Cyan", message)
}

/**
* MongoDB database conexion and registering logs
**/

func set(collection string, key string, val et.Json) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)

	filter := et.Json{"key": key}
	update := et.Json{
		"$set": et.Json{
			"value": val,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(conn.ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func Get(collection string, key, def string) (string, error) {
	if conn == nil {
		return def, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)
	filter := et.Json{"key": key}
	var result MongoDocument

	err := coll.FindOne(conn.ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return def, IsNil
	} else if err != nil {
		return def, err
	}

	return result.Value, nil
}

func Del(collection string, key string) (int64, error) {
	if conn == nil {
		return 0, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)
	filter := et.Json{"key": key}
	deleteResult, err := coll.DeleteOne(conn.ctx, filter)
	if err != nil {
		return 0, err
	}

	return deleteResult.DeletedCount, nil
}

func Set(collection string, key string, val interface{}) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)
	var valStr string

	switch v := val.(type) {
	case et.Json:
		valStr = v.ToString()
	case et.Items:
		valStr = v.ToString()
	case et.Item:
		valStr = v.ToString()
	default:
		var ok bool
		valStr, ok = val.(string)
		if !ok {
			return Log(msg.ERR_INVALID_TYPE)
		}
	}

	filter := et.Json{"key": key}
	update := et.Json{"$set": et.Json{"value": valStr}}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(conn.ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func HSet(collection string, key string, val map[string]string) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)
	filter := et.Json{"key": key}
	update := et.Json{"$set": et.Json{"value": val}}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(conn.ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func HGet(collection string, key string) (map[string]string, error) {
	if conn == nil {
		return nil, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)
	filter := et.Json{"key": key}
	var result et.Json

	err := coll.FindOne(conn.ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, IsNil
	} else if err != nil {
		return nil, err
	}

	value, ok := result["value"].(map[string]string)
	if !ok {
		return nil, fmt.Errorf("invalid type for value")
	}

	return value, nil
}

func HSetAtrib(collection string, key, atr, val string) error {
	return HSet(collection, key, map[string]string{atr: val})
}

func HGetAtrib(collection string, key, atr string) (string, error) {
	atribs, err := HGet(collection, key)
	if err != nil {
		return "", err
	}

	return atribs[atr], nil
}

func HDel(collection string, key, atr string) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)
	filter := et.Json{"key": key}
	update := et.Json{"$unset": et.Json{"value." + atr: ""}}

	_, err := coll.UpdateOne(conn.ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func Empty(collection string) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	coll := conn.db.Collection(collection)
	_, err := coll.DeleteMany(conn.ctx, et.Json{})
	return err
}
