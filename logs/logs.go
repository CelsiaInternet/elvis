package logs

import (
	"errors"
	"fmt"
	"os"

	"github.com/celsiainternet/elvis/stdrout"
	"github.com/celsiainternet/elvis/utility"
)

func printLn(kind string, color string, args ...any) {
	stdrout.Printl(kind, color, args...)
}

func Log(kind string, args ...any) error {
	printLn(kind, "", args...)
	return nil
}

func Logf(kind string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	printLn(kind, "", message)
}

func Alert(err error) error {
	functionName := utility.PrintFunctionName()
	if err != nil {
		printLn("Alert", "Yellow", err.Error(), " - ", functionName)
	}

	return err
}

func Alertm(message string) error {
	functionName := utility.PrintFunctionName()
	err := errors.New(message)
	printLn("Alert", "Yellow", err.Error(), " - ", functionName)
	return err
}

func Alertf(format string, args ...any) error {
	functionName := utility.PrintFunctionName()
	message := fmt.Sprintf(format, args...)
	err := errors.New(message)
	printLn("Alert", "Yellow", err.Error(), " - ", functionName)
	return err
}

func Traces(err error) error {
	_, err = utility.Traces("Error", "red", err)

	return err
}

func Error(err error) error {
	printLn("Error", "red", err.Error())

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
	printLn("Info", "Blue", v...)
}

func Infof(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	printLn("Info", "Blue", message)
}

func Fatal(v ...any) {
	printLn("Fatal", "Red", v...)
	os.Exit(1)
}

func Panic(v ...any) {
	printLn("Panic", "Red", v...)
	os.Exit(1)
}

func Panice(v error) error {
	printLn("Panic", "Red", v.Error())
	os.Exit(1)

	return v
}

func Panicf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	err := errors.New(message)

	return Panice(err)
}

func Panicm(v string) error {
	err := errors.New(v)

	return Panice(err)
}

func Ping(args ...any) {
	printLn("PONG", "Cyan", args...)
}

func Pong(args ...any) {
	printLn("PING", "Cyan", args...)
}

func Debug(v ...any) {
	printLn("Debug", "Cyan", v...)
}

func Debugf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	printLn("Debug", "Cyan", message)
}
