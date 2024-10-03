package utility

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/strs"
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
const DAY_SECOND = 60 * 60 * 24 * 1
const SELECt = "SELECT"
const INSERT = "INSERT"
const UPDATE = "UPDATE"
const DELETE = "DELETE"
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

var start = time.Now()
var ping = 0
var locks = make(map[string]*sync.RWMutex)
var count = make(map[string]int64)

func Ping() {
	ping++
	console.InfoF(`PING %d`, ping)
}

func Pong() {
	console.Info("PONG")
	ping = 0
}

func Now() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

/**
* GetOTP return a random number
* @param length int
* @return string
**/
func GetOTP(length int) string {
	const charset = "0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

/**
* UUID
* @return string
**/
func UUID() string {
	return uuid.NewString()
}

/**
* NewId
* @return string
**/
func NewId() string {
	return UUID()
}

/**
* GenId
* @param id string
* @return string
**/
func GenId(id string) string {
	if map[string]bool{"": true, "*": true, "new": true}[id] {
		return NewId()
	}

	return id
}

/**
* GenKey
* @param id string
* @return string
**/
func GenKey(id string) string {
	if map[string]bool{"": true, "-1": true, "*": true, "new": true}[id] {
		return uuid.NewString()
	}

	return id
}

/**
* More return the next value of a serie
* @param tag string
* @return int
**/
func More(tag string, expiration time.Duration) int64 {
	lock := locks[tag]
	if lock == nil {
		lock = &sync.RWMutex{}
		locks[tag] = lock
	}

	lock.Lock()
	defer lock.Unlock()

	n, ok := count[tag]
	if !ok {
		n = 0
	} else {
		n++
	}
	count[tag] = 0

	clean := func() {
		delete(count, tag)
		delete(locks, tag)
	}

	duration := expiration * time.Second
	if duration != 0 {
		go time.AfterFunc(duration, clean)
	}

	return n
}

/**
* UUIndex return the next value of a serie
* @param tag string
* @return int64
**/
func UUIndex(tag string) int64 {
	now := time.Now()
	result := now.UnixMilli() * 10000
	key := fmt.Sprintf("%s:%d", tag, result)
	n := More(key, 1*time.Second)
	result = result + int64(n)

	return result
}

/**
* Pointer return a string with the format collection/id
* @param collection string
* @param id string
* @return string
**/
func Pointer(collection string, id string) string {
	return strs.Format("%s/%s", collection, id)
}

/**
* Contains return true if the value is in the slice
* @param pointer string
* @return string
**/
func Contains(c []string, v string) bool {
	return slices.Contains(c, v)
}

/**
* ContainsInt return true if the value is in the slice
* @param pointer string
* @return string
**/
func ContainsInt(c []int, v int) bool {
	for _, i := range c {
		if i == v {
			return true
		}
	}

	return false
}

/**
* InStr return true if the value is in the slice
* @param pointer string
* @return string
**/
func InStr(val string, in []string) bool {
	ok := slices.Contains(in, val)

	return ok
}

/**
* InInt return true if the value is in the slice
* @param pointer string
* @return string
**/
func InInt(val string, in []string) bool {
	ok := slices.Contains(in, val)

	return ok
}

/**
* TimeDifference return the difference between two dates
* @param dateInt any
* @param dateEnd any
* @return time.Duration
**/
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

/**
* GeneratePortNumber return a random number
* @return int
**/
func GeneratePortNumber() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 1000
	max := 99999
	port := rand.Intn(max-min+1) + min

	return port
}

/**
* IsJsonBuild return true if the string is a json
* @param str string
* @return bool
**/
func IsJsonBuild(str string) bool {
	result := strings.Contains(str, "[")
	result = result && strings.Contains(str, "]")
	return result
}

/**
* FindIndex return the index of a value in a slice
* @param arr []string
* @param valor string
* @return int
**/
func FindIndex(arr []string, valor string) int {
	for i, v := range arr {
		if v == valor {
			return i
		}
	}
	return -1
}

/**
* OkOrNot return the value of the condition
* @param condition bool
* @param ok interface{}
* @param not interface{}
* @return interface{}
**/
func OkOrNot(condition bool, ok interface{}, not interface{}) interface{} {
	if condition {
		return ok
	} else {
		return not
	}
}

/**
* ExtractMencion return the mentions in a string
* @param str string
* @return []string
**/
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

/**
* Quote return the value in a string format
* @param val interface{}
* @return any
**/
func Quote(val interface{}) any {
	switch v := val.(type) {
	case string:
		return strs.Format(`'%s'`, v)
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
		return strs.Format(`'%s'`, v.Format("2006-01-02 15:04:05"))
	case []interface{}:
		var r string
		for _, _v := range v {
			q := Quote(_v).(string)
			if len(r) == 0 {
				r = q
			} else {
				r = strs.Format(`%v, %v`, r, q)
			}
		}
		return strs.Format(`'[%s]'`, r)
	case map[string]interface{}:
		var r string
		for k, _v := range v {
			q := Quote(_v).(string)
			if len(r) == 0 {
				r = strs.Format(`"%v": %v`, k, q)
			} else {
				r = strs.Format(`%v, "%v": %v`, r, k, q)
			}
		}
		return strs.Format(`'%s'`, r)
	case []map[string]interface{}:
		var r string
		for _, _v := range v {
			q := Quote(_v).(string)
			if len(r) == 0 {
				r = q
			} else {
				r = strs.Format(`%v, %v`, r, q)
			}
		}
		return strs.Format(`'[%s]'`, r)
	case nil:
		return "NULL"
	default:
		logs.Errorf("Not quote type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

/**
* Params return the value in a string format
* @param str string
* @param args ...any
* @return string
**/
func Params(str string, args ...any) string {
	var result string = str
	for i, v := range args {
		p := strs.Format(`$%d`, i+1)
		rp := strs.Format(`%v`, v)
		result = strs.Replace(result, p, rp)
	}

	return result
}

/**
* ParamQuote return the value in a string format
* @param str string
* @param args ...any
* @return string
**/
func ParamQuote(str string, args ...any) string {
	for i, arg := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`%v`, Quote(arg))
		str = strings.ReplaceAll(str, old, new)
	}

	return str
}

/**
* Address return the value in a string format
* @param host string
* @param port int
* @return string
**/
func Address(host string, port int) string {
	return strs.Format("%s:%d", host, port)
}

/**
* BannerTitle return the value in a string format
* @param name string
* @param version string
* @param size int
* @return string
**/
func BannerTitle(name, version string, size int) string {
	return strs.Format(`{{ .Title "%s V%s" "" %d }}`, name, version, size)
}

/**
* GoMod return the value in a string format
* @param atrib string
* @return string
* @return error
**/
func GoMod(atrib string) (string, error) {
	var result string
	rutaArchivoGoMod := "./go.mod"

	contenido, err := os.ReadFile(rutaArchivoGoMod)
	if err != nil {
		return "", err
	}

	lineas := strings.Split(string(contenido), "\n")
	for _, linea := range lineas {
		if strings.HasPrefix(linea, atrib) {
			partes := strings.Fields(linea)
			if len(partes) > 1 {
				result = partes[1]
				break
			}
		}
	}

	return result, nil
}

/**
* StartTime
**/
func StartTime() {
	start = time.Now()
}

/**
* Duration
**/
func Duration() {
	duration := time.Since(start) // Calcula la duración
	console.DebugF("La función tardó %v en ejecutarse\n", duration)
}

func GitVersion(idx int) string {
	result := "v0.0.0"
	attr := strs.Format("--abbrev=%d", idx)
	out, err := exec.Command("git", "describe", "--tags", attr).Output()
	if err == nil {
		result = string(out)
	}

	return result
}
