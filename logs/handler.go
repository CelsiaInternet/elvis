package logs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/msg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var IsNil = mongo.ErrNoDocuments

// Estructura que representa el documento almacenado en MongoDB
type MongoDocument struct {
	Key        string            `bson:"key"`
	Value      string            `bson:"value"`
	Attributes map[string]string `bson:"attributes,omitempty"`
	Expiration time.Time         `bson:"expiration,omitempty"`
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
var useColor = true

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
		useColor = false
	}
}

func log(kind string, color string, args ...any) string {
	kind = strings.ToUpper(kind)
	message := fmt.Sprint(args...)
	now := time.Now().Format("2006/01/02 15:04:05")
	var result string

	switch color {
	case "Reset":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + message + Reset
	case "Red":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Red + message + Reset
	case "Green":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Green + message + Reset
	case "Yellow":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Yellow + message + Reset
	case "Blue":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Blue + message + Reset
	case "Purple":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Purple + message + Reset
	case "Cyan":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Cyan + message + Reset
	case "Gray":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Gray + message + Reset
	case "White":
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + White + message + Reset
	default:
		result = now + Purple + fmt.Sprintf(" [%s]: ", kind) + Reset + Green + message + Reset
	}

	println(result)

	return result
}

func Log(kind string, args ...any) error {
	log(kind, "", args...)
	return nil
}

func Logf(kind string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	log(kind, "", message)
}

func Traces(kind, color string, err error) ([]string, error) {
	var n int = 1
	var traces []string = []string{err.Error()}

	log(kind, color, err.Error())

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
			log("TRACE", color, trace)
		}
	}

	return traces, err
}

func Alert(err error) error {
	if err != nil {
		log("Alert", "Yellow", err.Error())
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
	log("Info", "Blue", v...)
}

func Infof(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	log("Info", "Blue", message)
}

func Fatal(v ...any) {
	log("Fatal", "Red", v...)
	os.Exit(1)
}

func Panic(v ...any) {
	log("Panic", "Red", v...)
	os.Exit(1)
}

func Ping() {
	log("PING", "")
}

func Pong() {
	log("PONG", "")
}

func Debug(v ...any) {
	log("Debug", "Cyan", v...)
}

func Debugf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	log("Debug", "Cyan", message)
}

func SetCtx(ctx context.Context, collection *mongo.Collection, key, val string, second time.Duration) error {
	if collection == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	duration := second * time.Second
	expiration := time.Now().Add(duration)

	filter := bson.M{"key": key}
	update := bson.M{
		"$set": bson.M{
			"value":      val,
			"expiration": expiration,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func GetCtx(ctx context.Context, collection *mongo.Collection, key, def string) (string, error) {
	if collection == nil {
		return def, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	filter := bson.M{"key": key}
	var result MongoDocument

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return def, IsNil
	} else if err != nil {
		return def, err
	}

	if time.Now().After(result.Expiration) {
		return def, IsNil
	}

	return result.Value, nil
}

func DelCtx(ctx context.Context, collection *mongo.Collection, key string) (int64, error) {
	if collection == nil {
		return 0, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	filter := bson.M{"key": key}
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return deleteResult.DeletedCount, nil
}

func HSetCtx(ctx context.Context, collection *mongo.Collection, key string, val map[string]string) error {
	if collection == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	filter := bson.M{"key": key}
	update := bson.M{
		"$set": bson.M{
			"attributes": val,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func HGetCtx(ctx context.Context, collection *mongo.Collection, key string) (map[string]string, error) {
	if collection == nil {
		return map[string]string{}, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	filter := bson.M{"key": key}
	var result MongoDocument

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return map[string]string{}, IsNil
	} else if err != nil {
		return map[string]string{}, err
	}

	return result.Attributes, nil
}

func HDelCtx(ctx context.Context, collection *mongo.Collection, key, atr string) error {
	if collection == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	filter := bson.M{"key": key}
	update := bson.M{
		"$unset": bson.M{
			"attributes." + atr: "",
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func Get(key, def string) (string, error) {
	if conn == nil {
		return def, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	return GetCtx(conn.ctx, conn.collection, key, def)
}

func Set(key string, val interface{}, second time.Duration) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	switch v := val.(type) {
	case et.Json:
		return SetCtx(conn.ctx, conn.collection, key, v.ToString(), second)
	case et.Items:
		return SetCtx(conn.ctx, conn.collection, key, v.ToString(), second)
	case et.Item:
		return SetCtx(conn.ctx, conn.collection, key, v.ToString(), second)
	default:
		valStr, ok := val.(string)
		if ok {
			return SetCtx(conn.ctx, conn.collection, key, valStr, second)
		}
	}

	return nil
}

func Del(key string) (int64, error) {
	if conn == nil {
		return 0, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	return DelCtx(conn.ctx, conn.collection, key)
}

func HSet(key string, val map[string]string) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	return HSetCtx(conn.ctx, conn.collection, key, val)
}

func HGet(key string) (map[string]string, error) {
	if conn == nil {
		return map[string]string{}, Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	return HGetCtx(conn.ctx, conn.collection, key)
}

func HSetAtrib(key, atr, val string) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	return HSetCtx(conn.ctx, conn.collection, key, map[string]string{atr: val})
}

func HGetAtrib(key, atr string) (string, error) {
	if conn == nil {
		return "", Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	atribs, err := HGetCtx(conn.ctx, conn.collection, key)
	if err != nil {
		return "", err
	}

	return atribs[atr], nil
}

func HDel(key, atr string) error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	return HDelCtx(conn.ctx, conn.collection, key, atr)
}

func Empty() error {
	if conn == nil {
		return Log(msg.ERR_NOT_COLLETION_MONGO)
	}

	_, err := conn.collection.DeleteMany(context.Background(), bson.M{})
	return err
}
