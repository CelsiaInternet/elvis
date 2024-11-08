package logs

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/celsiainternet/elvis/console"
	"go.mongodb.org/mongo-driver/mongo"
)

var IsNil = mongo.ErrNoDocuments

func NewError(message string) error {
	return errors.New(message)
}

func NewErrorf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return NewError(message)
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
	err := NewError(message)
	return Alert(err)
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
	err := NewError(message)
	return Error(err)
}

func Errorf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	err := NewError(message)
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

func Panice(v error) error {
	console.Printl("Panic", "Red", v.Error())
	os.Exit(1)

	return v
}

func Panicf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	err := NewError(message)

	return Panice(err)
}

func Panicm(v string) error {
	err := NewError(v)

	return Panice(err)
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
