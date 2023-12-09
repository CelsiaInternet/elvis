package bus

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/cgalvisleon/elvis/console"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/response"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/go-chi/chi"
)

type Endpoint struct {
	Ok     bool
	Method string
	Scheme string
	Host   string
	Path   string
	Proto  string
	Header http.Header
	Params []e.Json
	Query  url.Values
	Body   e.Json
}

func (c *Endpoint) Define() e.Json {
	return e.Json{
		"method": c.Method,
		"scheme": c.Scheme,
		"host":   c.Host,
		"path":   c.Path,
		"proto":  c.Proto,
		"header": c.Header,
		"params": c.Params,
		"query":  c.Query,
		"body":   c.Body,
	}
}

var DefaultLogger func(next http.Handler) http.Handler

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

const (
	Ldate     = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                     // the time in the local time zone: 01:23:23
	LstdFlags = Ldate | Ltime // initial values for the standard logger
)

func ApiManager(next http.Handler) http.Handler {
	return DefaultLogger(next)
}

type LogFormatter interface {
	NewLogEntry(r *http.Request) LogEntry
}

type LogEntry interface {
	Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{})
	Panic(v interface{}, stack []byte)
}

type LoggerInterface interface {
	Print(v ...interface{})
}

type DefaultLogFormatter struct {
	Logger  LoggerInterface
	NoColor bool
}

func (l *DefaultLogFormatter) NewLogEntry(r *http.Request) LogEntry {
	useColor := !l.NoColor
	entry := &defaultLogEntry{
		DefaultLogFormatter: l,
		request:             r,
		buf:                 &bytes.Buffer{},
		useColor:            useColor,
	}

	entry.buf.WriteString("from ")
	entry.buf.WriteString(r.RemoteAddr)
	entry.buf.WriteString(" - ")

	return entry
}

type defaultLogEntry struct {
	*DefaultLogFormatter
	request  *http.Request
	buf      *bytes.Buffer
	useColor bool
}

func (l *defaultLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.Print(l.buf.String())
}

func (l *defaultLogEntry) Panic(v interface{}, stack []byte) {
	console.Ping()
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func RequestLogger(f LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			host := r.Host
			path := r.URL.Path
			proto := r.Proto

			defer func() {

			}()

			header := r.Header
			var params []e.Json = []e.Json{}
			ctx := r.Context()
			if rctx := chi.RouteContext(ctx); rctx != nil {
				for k := len(rctx.URLParams.Keys) - 1; k >= 0; k-- {
					key := rctx.URLParams.Keys[k]
					if key == "*" {
						continue
					}

					value := rctx.URLParams.Values[k]
					params = append(params, e.Json{
						"key":   key,
						"value": value,
					})
				}
			}
			query := r.URL.Query()
			body, _ := response.GetBody(r)
			endpoint := &Endpoint{
				Ok:     false,
				Method: method,
				Scheme: scheme,
				Host:   host,
				Path:   path,
				Proto:  proto,
				Header: header,
				Params: params,
				Query:  query,
				Body:   body,
			}

			endpoint = FindPath(endpoint)
			if endpoint.Ok {
				ExecutePath(endpoint, w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}
}

func FindPath(endpoint *Endpoint) *Endpoint {
	if endpoint.Path == "/api/version" {
		endpoint.Ok = true
	}

	return endpoint
}

func ExecutePath(endpoint *Endpoint, w http.ResponseWriter, r *http.Request) {
	apiUrl := utility.Format(`%s://%s%s`, endpoint.Scheme, endpoint.Host, endpoint.Path)
	body := endpoint.Body
	bodyParams := []byte(body.ToString())
	req, err := http.NewRequest(endpoint.Method, apiUrl, bytes.NewBuffer(bodyParams))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	req.Header = endpoint.Header
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		response.HTTPError(w, r, res.StatusCode, err.Error())
		return
	}
	defer res.Body.Close()

	e, err := json.Marshal(res.Body)
	if err != nil {
		response.HTTPError(w, r, res.StatusCode, err.Error())
		return
	}

	response.WriteResponse(w, res.StatusCode, e)
}

func init() {
	color := true
	if runtime.GOOS == "windows" {
		color = false
	}
	DefaultLogger = RequestLogger(&DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: !color})
}
