package console

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/timezone"
)

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

func Printl(kind string, color string, args ...any) string {
	kind = strings.ToUpper(kind)
	message := fmt.Sprint(args...)
	now := timezone.Now()
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

func NewError(message string) error {
	err := errors.New(message)

	return err
}

func NewErrorF(format string, args ...any) error {
	message := strs.Format(format, args...)
	err := NewError(message)

	return err
}

func LogK(kind string, args ...any) error {
	Printl(kind, "", args...)

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
	Printl("Rpc", "Blue", message...)

	return nil
}

func Debug(args ...any) error {
	Printl("Debug", "Cyan", args...)
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
	Printl("Info", "Blue", args...)
	return nil
}

func InfoF(format string, args ...any) error {
	message := strs.Format(format, args...)
	return Info(message)
}

func PrintFunctionName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullFuncName := runtime.FuncForPC(pc).Name()

	return fullFuncName
}

func Alert(message string) error {
	functionName := PrintFunctionName()
	err := NewError(message)
	Printl("Alert", "Yellow", err.Error(), " - ", functionName)
	return err
}

func AlertE(err error) error {
	functionName := PrintFunctionName()
	if err != nil {
		Printl("Alert", "Yellow", err.Error(), " - ", functionName)
	}
	return err
}

func AlertF(format string, args ...any) error {
	functionName := PrintFunctionName()
	err := NewError(fmt.Sprintf(format, args...))
	Printl("Alert", "Yellow", err.Error(), " - ", functionName)
	return err
}

func Error(err error) error {
	Printl("Error", "Red", err.Error())

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
	Printl("Fatal", "Red", v...)
	os.Exit(1)
}

func FatalF(format string, args ...any) {
	message := strs.Format(format, args...)
	Fatal(message)
}

func Panic(err error) error {
	Printl("Panic", "Red", err.Error())
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
