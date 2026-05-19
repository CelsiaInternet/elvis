package stdrout

import (
	"fmt"
	"runtime"
	"strings"

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

func colorFor(color string) string {
	switch color {
	case "Reset":
		return Reset
	case "Red":
		return Red
	case "Green":
		return Green
	case "Yellow":
		return Yellow
	case "Blue":
		return Blue
	case "Purple":
		return Purple
	case "Cyan":
		return Cyan
	case "Gray":
		return Gray
	case "White":
		return White
	default:
		return Green
	}
}

func Printl(kind string, color string, args ...any) string {
	kind = strings.ToUpper(kind)
	message := fmt.Sprint(args...)
	now := timezone.Now()

	var b strings.Builder
	b.Grow(len(now) + len(Purple) + 2 + len(kind) + 4 + len(Reset) + len(colorFor(color)) + len(message) + len(Reset))
	b.WriteString(now)
	b.WriteString(Purple)
	b.WriteString(" [")
	b.WriteString(kind)
	b.WriteString("]: ")
	b.WriteString(Reset)
	b.WriteString(colorFor(color))
	b.WriteString(message)
	b.WriteString(Reset)

	result := b.String()
	println(result)
	return result
}
