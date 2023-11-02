/**
*
**/
package console

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/logs"
	_ "github.com/joho/godotenv/autoload"
)

func NewError(message string) error {
	err := errors.New(message)

	return err
}

func NewErrorF(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	err := NewError(message)

	return err
}

func FuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc)
	return function.Name()
}

func LogC(kind string, color string, args ...any) string {
	event.EventPublish("logs", map[string]interface{}{
		"kind": kind,
		"args": args,
	})
	return logs.Logln(kind, color, args...)
}

func LogK(kind string, args ...any) {
	LogC(kind, "Green", args...)
}

func LogKF(kind string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	LogK(kind, message)
}

func Log(args ...any) {
	LogC("Log", "Green", args...)
}

func LogF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	Log(message)
}

func Print(args ...any) {
	message := ""
	for i, arg := range args {
		if i == 0 {
			message = fmt.Sprintf("%v", arg)
		} else {
			message = fmt.Sprintf("%s, %v", message, arg)
		}
	}
	Log(message)
}

func Info(args ...any) {
	LogC("Info", "Blue", args...)
}

func InfoF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	Info(message)
}

func Alert(message string) error {
	err := NewError(message)
	LogC("Alert", "Yellow", err)
	return err
}

func AlertF(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return Alert(message)
}

func Error(err error) error {
	LogC("ERROR", "Red", fmt.Sprintf(`%s - %s`, FuncName(), err.Error()))
	return err
}

func ErrorM(message string) error {
	err := NewError(message)
	LogC("ERROR", "Red", fmt.Sprintf(`%s - %s`, FuncName(), err.Error()))
	return err
}

func ErrorF(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	err := NewError(message)
	LogC("ERROR", "Red", fmt.Sprintf(`%s - %s`, FuncName(), err.Error()))
	return err
}

func Fatal(v ...any) {
	LogC("Fatal", "Red", v...)
	os.Exit(1)
}

func FatalF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
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
