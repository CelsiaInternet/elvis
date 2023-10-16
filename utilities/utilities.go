package utilities

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

const NOT_FOUND = "Not found"
const FOUND = "Found"
const FOR_DELETE = "-2"
const OF_SYSTEM = "-1"
const ACTIVE = "0"
const ARCHIVED = "1"
const CANCELLED = "2"
const IN_PROCESS = "3"
const PENDING_APPROVAL = "4"
const APPROVAL = "5"
const REFUSED = "6"
const STOP = "Stop"
const CACHE_TIME = 60 * 60 * 24 * 1
const SELECt = "SELECT"
const INSERT = "INSERT"
const UPDATE = "UPDATE"
const DELETE = "DELETE"
const _STATE = "_STATE"
const BEFORE_INSERT = "BEFORE_INSERT"
const AFTER_INSERT = "AFTER_INSERT"
const BEFORE_UPDATE = "BEFORE_UPDATE"
const AFTER_UPDATE = "AFTER_UPDATE"
const BEFORE_STATE = "BEFORE_STATE"
const AFTER_STATE = "AFTER_STATE"
const BEFORE_DELETE = "BEFORE_DELETE"
const AFTER_DELETE = "AFTER_DELETE"
const VALUE_NOT_BOOL = "Value is not bolean"
const ROWS = 30

var ping = 0

func Ping() {
	ping++
	console.InfoF(`PING %d`, ping)
}

func Pong() {
	console.Info("PONG")
	ping = 0
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCodeVerify(length int) string {
	const charset = "0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func GenId(id string) string {
	if map[string]bool{"": true, "*": true, "new": true}[id] {
		return uuid.NewString()
	}

	return id
}

func NewId() string {
	return GenId("")
}

func Format(format string, args ...any) string {
	var result string
	result = fmt.Sprintf(format, args...)

	return result
}

func FormatUppCase(format string, args ...any) string {
	result := Format(format, args...)

	return Uppcase(result)
}

func FormatLowCase(format string, args ...any) string {
	result := Format(format, args...)

	return Lowcase(result)
}

func Replace(str string, old string, new string) string {
	return strings.ReplaceAll(str, old, new)
}

func ReplaceAll(str string, olds []string, new string) string {
	var result string = str
	for _, str := range olds {
		result = strings.ReplaceAll(result, str, new)
	}

	return result
}

func Name(str string) string {
	regex := `[0-9\s]+`
	pattern := regexp.MustCompile(regex)
	return pattern.ReplaceAllString(str, "_")
}

func Trim(str string) string {
	return strings.Trim(str, " ")
}

func NotSpace(str string) string {
	return Replace(str, " ", "")
}

func Uppcase(s string) string {
	return strings.ToUpper(s)
}

func Lowcase(s string) string {
	return strings.ToLower(s)
}

func Titlecase(str string) string {
	var result string
	var ok bool
	for i, char := range str {
		s := fmt.Sprintf("%c", char)
		if i == 0 {
			s = strings.ToUpper(s)
		} else if s == "" {
			ok = true
		} else if ok {
			ok = false
			s = strings.ToUpper(s)
		}

		result = Append(result, s, "")
	}

	return result
}

func Pointer(collection string, id string) string {
	return Format("%s/%s", collection, id)
}

func Contains(c []string, v string) bool {
	return slices.Contains(c, v)
}

func ContainsInt(c []int, v int) bool {
	for _, i := range c {
		if i == v {
			return true
		}
	}

	return false
}

func ValidStr(val string, min int, notIn []string) bool {
	v := Replace(val, " ", "")
	ok := len(v) <= min
	if ok {
		return !ok
	}

	ok = Contains(notIn, val)

	return !ok
}

func ValidIn(val string, min int, in []string) bool {
	v := Replace(val, " ", "")
	ok := len(v) <= min
	if ok {
		return !ok
	}

	ok = Contains(in, val)

	return ok
}

func ValidId(val string) bool {
	return !Contains([]string{"", "*", "new"}, val)
}

func ValidInt(val int, notIn []int) bool {
	ok := slices.Contains(notIn, val)

	return !ok
}

func ValidNum(val float64, notIn []float64) bool {
	ok := slices.Contains(notIn, val)

	return !ok
}

func ValidName(val string) bool {
	regex := `^[a-zA-Z\s\']+`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(val)
}

func ValidEmail(val string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(val)
}

func ValidPhone(val string) bool {
	regex := `^\d{10}$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(val)
}

func ValidUUID(val string) bool {
	regex := `^(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(val)
}

func ValidCode(val string) bool {
	ok := len(val) < 6
	if ok {
		return !ok
	}

	regex := `^-?\d+$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(val)
}

func InStr(val string, in []string) bool {
	ok := slices.Contains(in, val)

	return ok
}

func InInt(val string, in []string) bool {
	ok := slices.Contains(in, val)

	return ok
}

func StrToBool(val string) (bool, error) {
	if Lowcase(val) == "true" {
		return true, nil
	} else if Lowcase(val) == "false" {
		return false, nil
	}

	return false, console.ErrorM(VALUE_NOT_BOOL)
}

func Empty(str1, str2 string) string {
	if len(str1) == 0 {
		return str2
	}

	return str1
}

func Append(str1, str2, sp string) string {
	if len(str1) == 0 {
		return str2
	}
	if len(str2) == 0 {
		return str1
	}

	return Format(`%s%s%s`, str1, sp, str2)
}

func AppendAny(val1, val2 any, sp string) string {
	any1 := NewAny(val1)
	any2 := NewAny(val2)

	if len(any1.String()) == 0 {
		return any2.String()
	}
	if len(any2.String()) == 0 {
		return any1.String()
	}

	return Format(`%v%s%v`, any1, sp, any2)
}

func Split(str, sep string) []string {
	return strings.Split(str, sep)
}

func GetSplitIndex(str, sep string, idx int) string {
	split := Split(str, sep)
	if idx < 0 {
		idx = len(split) + idx
	}

	if idx < len(split) {
		return split[idx]
	}

	return ""
}

func ApendAny(space string, strs ...any) string {
	var result string = ""
	for i, s := range strs {
		if i == 0 {
			result = fmt.Sprintf(`%v`, s)
		} else if len(result) == 0 && len(fmt.Sprint(s)) > 0 {
			result = fmt.Sprintf(`%v`, s)
		} else if len(result) > 0 && len(fmt.Sprint(s)) > 0 {
			result = fmt.Sprintf(`%s%v%v`, result, space, s)
		}
	}

	return result
}

func TimeDifference(dateInt, dateEnd any) time.Duration {
	var result time.Time
	layout := "2006-01-02T15:04:05.000Z"

	if dateInt == 0 {
		return result.Sub(result)
	}
	if dateEnd == 0 {
		return result.Sub(result)
	}
	_dateInt, err := time.Parse(layout, fmt.Sprint(dateInt))
	if err != nil {
		return result.Sub(result)
	}

	_dateEnd, err := time.Parse(layout, fmt.Sprint(dateEnd))
	if err != nil {
		return result.Sub(result)
	}

	return _dateInt.Sub(_dateEnd)
}

func RemoveAcents(str string) string {
	str = strings.ReplaceAll(str, "á", "a")
	str = strings.ReplaceAll(str, "é", "e")
	str = strings.ReplaceAll(str, "í", "i")
	str = strings.ReplaceAll(str, "ó", "o")
	str = strings.ReplaceAll(str, "ú", "u")

	str = strings.ReplaceAll(str, "Á", "A")
	str = strings.ReplaceAll(str, "É", "E")
	str = strings.ReplaceAll(str, "Í", "I")
	str = strings.ReplaceAll(str, "Ó", "O")
	str = strings.ReplaceAll(str, "Ú", "U")
	return str
}

func GeneratePortNumber() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 1000
	max := 99999
	port := rand.Intn(max-min+1) + min

	return port
}

func IsJsonBuild(str string) bool {
	result := strings.Contains(str, "[")
	result = result && strings.Contains(str, "]")
	return result
}

func FindIndex(arr []string, valor string) int {
	for i, v := range arr {
		if v == valor {
			return i
		}
	}
	return -1
}

func OkOrNot(condition bool, ok interface{}, not interface{}) interface{} {
	if condition {
		return ok
	} else {
		return not
	}
}

func ExtractMencion(str string) []string {
	patron := `@([a-zA-Z0-9_]+)`
	expresionRegular := regexp.MustCompile(patron)
	mencions := expresionRegular.FindAllString(str, -1)
	unique := make(map[string]bool)
	result := []string{}

	for _, val := range mencions {
		if !unique[val] {
			unique[val] = true
			result = append(result, val)
		}
	}

	return result
}

func Quote(val interface{}) any {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf(`'%s'`, v)
	case int:
		return v
	case float64:
		return v
	case float32:
		return v
	case int16:
		return v
	case int32:
		return v
	case int64:
		return v
	case bool:
		return v
	case time.Time:
		return fmt.Sprintf(`'%s'`, v.Format("2006-01-02 15:04:05"))
	case []interface{}:
		var r string
		for _, _v := range v {
			q := Quote(_v).(string)
			if len(r) == 0 {
				r = q
			} else {
				r = fmt.Sprintf(`%v, %v`, r, q)
			}
		}
		return fmt.Sprintf(`'[%s]'`, r)
	case map[string]interface{}:
		var r string
		for k, _v := range v {
			q := Quote(_v).(string)
			if len(r) == 0 {
				r = fmt.Sprintf(`"%v": %v`, k, q)
			} else {
				r = fmt.Sprintf(`%v, "%v": %v`, r, k, q)
			}
		}
		return fmt.Sprintf(`'%s'`, r)
	case []map[string]interface{}:
		var r string
		for _, _v := range v {
			q := Quote(_v).(string)
			if len(r) == 0 {
				r = q
			} else {
				r = fmt.Sprintf(`%v, %v`, r, q)
			}
		}
		return fmt.Sprintf(`'[%s]'`, r)
	case nil:
		return fmt.Sprintf(`%s`, "NULL")
	default:
		logs.Errorf("Not quote type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

func Params(str string, args ...any) string {
	var result string = str
	for i, v := range args {
		p := Format(`$%d`, i+1)
		rp := Format(`%v`, v)
		result = Replace(result, p, rp)
	}

	return result
}

func ParamQuote(str string, args ...any) string {
	for i, arg := range args {
		old := Format(`$%d`, i+1)
		new := Format(`%v`, Quote(arg))
		str = strings.ReplaceAll(str, old, new)
	}

	return str
}

func Address(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func BannerTitle(name, version string, size int) string {
	return fmt.Sprintf(`{{ .Title "%s V%s" "" %d }}`, name, version, size)
}

func ModuleName() (string, error) {
	var result string
	rutaArchivoGoMod := "./go.mod"

	contenido, err := os.ReadFile(rutaArchivoGoMod)
	if err != nil {
		return "", err
	}

	lineas := strings.Split(string(contenido), "\n")
	for _, linea := range lineas {
		if strings.HasPrefix(linea, "module") {
			partes := strings.Fields(linea)
			if len(partes) > 1 {
				result = partes[1]
				break
			}
		}
	}

	return result, nil
}
