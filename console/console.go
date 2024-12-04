package console

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/realtime"
	"github.com/celsiainternet/elvis/stdrout"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
)

func printLn(kind string, color string, args ...any) {
	stdrout.Printl(kind, color, args...)

	if realtime.Active() {
		realtime.Publish("logs", et.Json{
			"kind":    kind,
			"message": fmt.Sprint(args...),
		})
	}
}

func LogK(kind string, args ...any) error {
	printLn(kind, "", args...)

	return nil
}

func LogKF(kind string, format string, args ...any) error {
	message := strs.Format(format, args...)
	return LogK(kind, message)
}

func Log(args ...any) error {
	return LogK("Log", args...)
}

func LogF(format string, args ...any) error {
	message := strs.Format(format, args...)
	return Log(message)
}

func Rpc(args ...any) error {
	pc, _, _, _ := runtime.Caller(1)
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcName := fullFuncName[strings.LastIndex(fullFuncName, "/")+1:]
	message := append([]any{funcName, ":"}, args...)
	printLn("Rpc", "Blue", message...)

	return nil
}

func Debug(args ...any) error {
	printLn("Debug", "Cyan", args...)
	return nil
}

func DebugF(format string, args ...any) error {
	message := strs.Format(format, args...)
	return Debug(message)
}

func Print(args ...any) error {
	message := ""
	for i, arg := range args {
		if i == 0 {
			message = strs.Format("%v", arg)
		} else {
			message = strs.Format("%s, %v", message, arg)
		}
	}
	return Log(message)
}

func Info(args ...any) error {
	printLn("Info", "Blue", args...)
	return nil
}

func InfoF(format string, args ...any) error {
	message := strs.Format(format, args...)
	return Info(message)
}

func Alert(message string) error {
	functionName := utility.PrintFunctionName()
	err := errors.New(message)
	printLn("Alert", "Yellow", err.Error(), " - ", functionName)
	return err
}

func AlertE(err error) error {
	functionName := utility.PrintFunctionName()
	if err != nil {
		printLn("Alert", "Yellow", err.Error(), " - ", functionName)
	}
	return err
}

func AlertF(format string, args ...any) error {
	functionName := utility.PrintFunctionName()
	message := fmt.Sprintf(format, args...)
	err := errors.New(message)
	printLn("Alert", "Yellow", err.Error(), " - ", functionName)
	return err
}

func Error(err error) error {
	printLn("Error", "Red", err.Error())

	return err
}

func ErrorM(message string) error {
	err := errors.New(message)
	return Error(err)
}

func ErrorF(format string, args ...any) error {
	message := strs.Format(format, args...)
	err := errors.New(message)
	return Error(err)
}

func Fatal(v ...any) {
	printLn("Fatal", "Red", v...)
	os.Exit(1)
}

func FatalF(format string, args ...any) {
	message := strs.Format(format, args...)
	Fatal(message)
}

func Panic(err error) error {
	printLn("Panic", "Red", err.Error())
	os.Exit(1)

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
	Log("PONG")
}

func Pong() {
	Log("PING")
}
