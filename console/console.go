/**
*
**/
package console

import (
	"errors"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/strs"
	_ "github.com/joho/godotenv/autoload"
)

func NewError(message string) error {
	err := errors.New(message)

	return err
}

func NewErrorF(format string, args ...any) error {
	message := strs.Format(format, args...)
	err := NewError(message)

	return err
}

func LogC(kind string, color string, args ...any) string {
	event.Action("logs", map[string]interface{}{
		"kind": kind,
		"args": args,
	})
	return logs.Logln(kind, color, args...)
}

func LogK(kind string, args ...any) {
	LogC(kind, "Green", args...)
}

func LogKF(kind string, format string, args ...any) {
	message := strs.Format(format, args...)
	LogK(kind, message)
}

func Log(args ...any) {
	LogC("Log", "Green", args...)
}

func LogF(format string, args ...any) {
	message := strs.Format(format, args...)
	Log(message)
}

func Print(args ...any) {
	message := ""
	for i, arg := range args {
		if i == 0 {
			message = strs.Format("%v", arg)
		} else {
			message = strs.Format("%s, %v", message, arg)
		}
	}
	Log(message)
}

func Info(args ...any) {
	LogC("Info", "Blue", args...)
}

func InfoF(format string, args ...any) {
	message := strs.Format(format, args...)
	Info(message)
}

func Alert(message string) error {
	err := NewError(message)
	LogC("Alert", "Yellow", err)
	return err
}

func AlertF(format string, args ...any) error {
	message := strs.Format(format, args...)
	return Alert(message)
}

func Error(err error) error {
	var n int = 1
	var trces []string = []string{err.Error()}

	logs.Logln("ERROR", "Red", err.Error())

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
			trace := strs.Format("%s:%d func:%s", file, line, name)
			trces = append(trces, trace)
			logs.Logln("TRACE", "Red", trace)
		}
	}

	event.Action("logs", map[string]interface{}{
		"kind": "ERROR",
		"args": trces,
	})

	return err
}

func ErrorM(message string) error {
	err := NewError(message)
	return Error(err)
}

func ErrorF(format string, args ...any) error {
	message := strs.Format(format, args...)
	err := NewError(message)
	return Error(err)
}

func Fatal(v ...any) {
	LogC("Fatal", "Red", v...)
	os.Exit(1)
}

func FatalF(format string, args ...any) {
	message := strs.Format(format, args...)
	Fatal(message)
}

func Panic(v ...any) {
	LogC("Panic", "Red", v...)
	os.Exit(1)
}

func PanicE(err error) error {
	Panic(err)
	return err
}

func PanicM(message string) error {
	err := ErrorM(message)
	Panic(err)
	return err
}

func PanicF(format string, args ...any) error {
	err := ErrorF(format, args...)
	Panic(err)
	return err
}

func Ping() {
	Log("PING")
}

func Pong() {
	Log("PONG")
}
