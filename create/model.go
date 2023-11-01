package create

const modelDockerfile = `ARG GO_VERSION=1.18.4

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates openssl git tzdata
RUN update-ca-certificates

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN gofmt -w . && go build ./cmd/$1

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/$1 ./$1

ENTRYPOINT ["./$1"]
`

const modelMain = `package main

import (
	"os"
	"os/signal"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/envar"
	_ "github.com/joho/godotenv/autoload"
	serv "$1/internal/service/$2"
)

func main() {
	SetvarInt("port", 3000, "Port server", "PORT")
	SetvarInt("rpc", 0, "Port rpc server", "RPC_PORT")
	SetvarStr("dbhost", "localhost", "Database host", "DB_HOST")
	SetvarInt("dbport", 5432, "Database port", "DB_PORT")
	SetvarStr("dbname", "", "Database name", "DB_NAME")
	SetvarStr("dbuser", "", "Database user", "DB_USER")
	SetvarStr("dbpass", "", "Database password", "DB_PASSWORD")

	serv, err := serv.New()
	if err != nil {
		console.Fatal(err)
	}

	go serv.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	serv.Close()
}
`

const modelService = `package module

import (
	"net"
	"net/http"

	v1 "$1/internal/service/$2/v1"
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/envar"
	mw "github.com/cgalvisleon/elvis/middleware"
	. "github.com/cgalvisleon/elvis/utilities"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
)

type Server struct {
	http *http.Server
	rpc  *net.Listener
}

func New() (*Server, error) {
	server := Server{}

	/**
	 * HTTP
	 **/
	port := EnvarInt(3300, "PORT")

	if port != 0 {
		r := chi.NewRouter()

		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(mw.Telemetry)

		latest := v1.New()

		r.Mount("/", latest)
		r.Mount("/v1", latest)

		handler := cors.AllowAll().Handler(r)
		addr := Format(":%d", port)
		serv := &http.Server{
			Addr:    addr,
			Handler: handler,
		}

		server.http = serv
	}

	/**
	 * RPC
	 **/
	port = EnvarInt(4200, "RPC_PORT")

	if port != 0 {
		serv := v1.NewRpc(port)

		server.rpc = &serv
	}

	return &server, nil
}

func (serv *Server) Close() error {
	v1.Close()
	return nil
}

func (serv *Server) Start() {
	go func() {
		if serv.http == nil {
			return
		}

		svr := serv.http
		console.LogKF("Http", "Running on http://localhost%s", svr.Addr)
		console.Fatal(serv.http.ListenAndServe())
	}()

	go func() {
		if serv.rpc == nil {
			return
		}

		svr := *serv.rpc
		console.LogKF("RPC", "Running on tcp:localhost:%s", svr.Addr().String())
		http.Serve(svr, nil)
	}()

	v1.Banner()

	<-make(chan struct{})
}
`

const modelApi = `package v1

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"time"

	pkg "$1/pkg/$2"
	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/utilities"
	"github.com/cgalvisleon/elvis/ws"
	"github.com/dimiro1/banner"
	"github.com/go-chi/chi"
	"github.com/mattn/go-colorable"
)

func New() http.Handler {
	r := chi.NewRouter()

	_, err := cache.Load()
	if err != nil {
		panic(err)
	}

	_, err = event.Load()
	if err != nil {
		panic(err)
	}

	_, err = ws.Load()
	if err != nil {
		panic(err)
	}

	Db, err := jdb.Load()
	if err != nil {
		panic(err)
	}

	_pkg := &pkg.Router{
		Repository: &pkg.Controller{
			Db: Db,
		},
	}

	r.Mount(pkg.PackagePath, _pkg.Routes())

	return r
}

func Close() {
}

func NewRpc(port int) net.Listener {
	rpc.HandleHTTP()

	result, err := net.Listen("tcp", Address("0.0.0.0", port))
	if err != nil {
		panic(err)
	}

	return result
}

func Banner() {
	time.Sleep(3 * time.Second)
	templ := BannerTitle(pkg.PackageName, pkg.PackageVersion, 4)
	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
	fmt.Println()
}
`

const modelEvent = `package $1

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
)

func initEvents() {
	err := event.Stack("<channel>", HostName, eventAction)
	if err != nil {
		console.Error(err)
	}

}

func eventAction(m event.CreatedEvenMessage) {
	data, err := ToJson(m.Data)
	if err != nil {
		console.Error(err)
	}

	console.Log("eventAction", data)
}
`

const modelModel = `package $1

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/module"
)

func initModels() error {
	if err := InitCore(); err != nil {
		return console.PanicE(err)
	}
	if err := InitModules(); err != nil {
		return console.PanicE(err)
	}

	return nil
}
`

const modelSchema = `package $1

import . "github.com/cgalvisleon/elvis/linq"

var $2 *Schema

func defineSchema() error {
	if $2 != nil {
		return nil
	}

	$2 = NewSchema(0, "$3")

	return nil
}
`

const modelhRpc = `package $1

import (
	"net/rpc"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/json"
)

type Service Item

func InitRpc() error {
	service := new(Service)

	err := rpc.Register(service)
	if err != nil {
		return console.Error(err)
	}

	return nil
}

func (c *Service) Version(require []byte, response *[]byte) error {
	result := Item{
		Ok: true,
		Result: Json{
			"service": PackageName,
			"host":    HostName,
			"help":    "",
		},
	}

	*response = result.ToByte()

	return nil
}
`

const modelMsg = `package $1

const (
	// MSG
	MSG = "Message"
)
`

const modelController = `package $1
import (
	"context"

	. "github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
)

type Controller struct {
	Db *jdb.Conn
}

func (c *Controller) Version(ctx context.Context) (Json, error) {
	company := EnvarStr("", "COMPANY")
	web := EnvarStr("", "WEB")
	version := EnvarStr("", "VERSION")
  service := Json{
		"version": version,
		"service": PackageName,
		"host":    HostName,
		"company": company,
		"web":     web,
		"help":    "",
	}

  event.EventPublish("service/starat", service)

	return service, nil
}

func (c *Controller) Init(ctx context.Context) {
	initModels()
	initEvents()
}

type Repository interface {
	Version(ctx context.Context) (Json, error)
	Init(ctx context.Context)
}
`

const modelRouter = `package $1

import (
	"context"
	"net/http"
	"os"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/response"
	. "github.com/cgalvisleon/elvis/router"
	"github.com/go-chi/chi"
)

var PackageName = "$1"
var PackageTitle = "$1"
var PackagePath = "/api/$1"
var PackageVersion = EnvarStr("VERSION")
var HostName, _ = os.Hostname()

type Router struct {
	Repository Repository
}

func (rt *Router) Routes() http.Handler {
	r := chi.NewRouter()

	PublicRoute(r, Get, "/version", rt.Version)
	// $2
	ProtectRoute(r, Get, "/$1/{id}", rt.Get$2ById)
	ProtectRoute(r, Post, "/$1", rt.UpSert$2)
	ProtectRoute(r, Put, "/$1/state/{id}", rt.State$2)
	ProtectRoute(r, Delete, "/$1/{id}", rt.Delete$2)
	ProtectRoute(r, Get, "/$1/all", rt.All$2)

	ctx := context.Background()
	rt.Repository.Init(ctx)

	console.LogKF(PackageName, "Router version:%s", PackageVersion)
	return r
}

func (rt *Router) Version(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := rt.Repository.Version(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}
`

const restHttp = `@host=localhost:3300
@token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IlVTRVIuQURNSU4iLCJhcHAiOiJEZXZvcHMtSW50ZXJuZXQiLCJuYW1lIjoiQ2VzYXIgR2FsdmlzIExlw7NuIiwia2luZCI6ImF1dGgiLCJ1c2VybmFtZSI6Iis1NzMxNjA0Nzk3MjQiLCJkZXZpY2UiOiJkZXZlbG9wIiwiZHVyYXRpb24iOjI1OTIwMDB9.dexIOute7r9o_P8U3t6l9RihN8BOnLl4xpoh9QbQI4k

###
GET /auth HTTP/1.1
Host: {{host}}/version
Authorization: Bearer {{token}}

###
POST /api/test/test HTTP/1.1
Host: {{host}}
Content-Type: application/json
Authorization: Bearer {{token}}
Content-Length: 227

{
}
`

const modelHandler = `package $1

import (
	"net/http"
	"strconv"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/response"
	. "github.com/cgalvisleon/elvis/utilities"
	"github.com/go-chi/chi"
)

var $2 *Model

func Define$2() error {
	if err := defineSchema(); err != nil {
		return console.PanicE(err)
	}

	if $2 != nil {
		return nil
	}

	$2 = NewModel($3, "$4", "Tabla de tipo", 1)
	$2.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	$2.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	$2.DefineColum("_state", "", "VARCHAR(80)", ACTIVE)
	$2.DefineColum("_id", "", "VARCHAR(80)", "-1")
	$2.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	$2.DefineColum("name", "", "VARCHAR(250)", "")
	$2.DefineColum("description", "", "TEXT", "")
	$2.DefineColum("_data", "", "JSONB", "{}")
	$2.DefineColum("index", "", "INTEGER", 0)
	$2.DefinePrimaryKey([]string{"_id"})
	$2.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"project_id",
		"name",
		"index",
	})
	$2.IntegrityAtrib(true)
	$2.Trigger(BeforeInsert, func(model *Model, old, new *Json, data Json) {

	})
	$2.Trigger(AfterInsert, func(model *Model, old, new *Json, data Json) {

	})
	$2.Trigger(BeforeUpdate, func(model *Model, old, new *Json, data Json) {

	})
	$2.Trigger(AfterUpdate, func(model *Model, old, new *Json, data Json) {

	})
	$2.Trigger(BeforeDelete, func(model *Model, old, new *Json, data Json) {

	})
	$2.Trigger(AfterDelete, func(model *Model, old, new *Json, data Json) {

	})
	
	if err := InitModel($2); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
*	Handler for CRUD data
 */
func Get$2ById(id string) (Item, error) {
	return $2.Select().
		Where($2.Column("_id").Eq(id)).
		First()
}

func UpSert$2(projectId, id, name, description string) (Item, error) {
	if !ValidId(id) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "_id")
	}

	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "name")
	}

	id = GenId(id)
	data := Json{}
	data["project_id"] = projectId
	data["_id"] = id
	data["name"] = name
	data["description"] = description
	return $2.Upsert(data).
		Where($2.Column("_id").Eq(id)).
		Command()
}

func State$2(id, state string) (Item, error) {
	if !ValidId(state) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "state")
	}

	return $2.Upsert(Json{
		"_state": state,
	}).
		Where($2.Column("_id").Eq(id)).
		And($2.Column("_state").Neg(state)).
		Command()
}

func Delete$2(id string) (Item, error) {
	return State$2(id, FOR_DELETE)
}

func All$2(projectId, state, search string, page, rows int, _select string) (List, error) {	
	if state == "" {
		state = ACTIVE
	}

	auxState := state

	cols := StrToColN(_select)

	if auxState == "*" {
		state = FOR_DELETE

		return $2.Select(cols).
			Where($2.Column("_state").Neg(state)).
			And($2.Column("project_id").In("-1", projectId)).
			And($2.Concat("NAME:", $2.Column("name"), ":DESCRIPTION", $2.Column("description"), ":DATA:", $2.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy($2.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return $2.Select(cols).
			Where($2.Column("_state").In("-1", state)).
			And($2.Column("project_id").In("-1", projectId)).
			And($2.Concat("NAME:", $2.Column("name"), ":DESCRIPTION", $2.Column("description"), ":DATA:", $2.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy($2.Column("name"), true).
			List(page, rows)
	} else {
		return $2.Select(cols).
			Where($2.Column("_state").Eq(state)).
			And($2.Column("project_id").In("-1", projectId)).
			And($2.Concat("NAME:", $2.Column("name"), ":DESCRIPTION", $2.Column("description"), ":DATA:", $2.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy($2.Column("name"), true).
			List(page, rows)
	}
}

/**
* Router
**/
func (rt *Router) UpSert$2(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	projectId := body.Str("project_id")
	id := body.Str("id")
	name := body.Str("name")
	data := body.Json("description")

	result, err := UpSert$2(projectId, id, name, data)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

func (rt *Router) Get$2ById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := Get$2ById(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

func (rt *Router) State$2(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	body, _ := response.GetBody(r)
	state := body.Str("state")

	result, err := State$2(id, state)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

func (rt *Router) Delete$2(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := Delete$2(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

func (rt *Router) All$2(w http.ResponseWriter, r *http.Request) {
	project_id := r.URL.Query().Get("project_id")
	state := r.URL.Query().Get("state")
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	rowsStr := r.URL.Query().Get("rows")
	_select := r.URL.Query().Get("select")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	rows, err := strconv.Atoi(rowsStr)
	if err != nil {
		rows = 10
	}

	result, err := All$2(project_id, state, search, page, rows, _select)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

/**
	// $2
	ProtectRoute(r, Get, "/$1/{id}", rt.Get$2ById)
	ProtectRoute(r, Post, "/$1", rt.UpSert$2)
	ProtectRoute(r, Put, "/$1/state/{id}", rt.State$2)
	ProtectRoute(r, Delete, "/$1/{id}", rt.Delete$2)
	ProtectRoute(r, Get, "/$1/all", rt.All$2)
**/
`
