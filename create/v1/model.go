package create

const modelDockerfile = `# Versión de Go como argumento
ARG GO_VERSION=1.23

# Stage 1: Compilación (builder)
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder

# Argumentos para el sistema operativo y la arquitectura
ARG TARGETOS
ARG TARGETARCH

# Instalación de dependencias necesarias
RUN apk update && apk add --no-cache ca-certificates openssl git tzdata \
    && update-ca-certificates

# Configuración de las variables de entorno para la build
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

# Directorio de trabajo
WORKDIR /src

# Descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Formatear el código Go
RUN gofmt -w .

# Compilar el binario
RUN go build -v -o /$1 ./cmd/$1

# Cambiar permisos del binario
RUN chmod +x /$1

# Stage 2: Imagen final mínima
FROM scratch

# Copiar certificados y binario
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /$1 /$1

# Establecer el binario como punto de entrada
ENTRYPOINT ["/$1"]
`

const modelMain = `package main

import (
	"os"
	"os/signal"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	serv "$1/internal/service/$2"	
)

func main() {
	envar.SetInt("port", 3000, "Port server", "PORT")
	envar.SetInt("rpc", 4200, "Port rpc server", "RPC_PORT")
	envar.SetStr("dbhost", "localhost", "Database host", "DB_HOST")
	envar.SetInt("dbport", 5432, "Database port", "DB_PORT")
	envar.SetStr("dbname", "", "Database name", "DB_NAME")
	envar.SetStr("dbuser", "", "Database user", "DB_USER")
	envar.SetStr("dbpass", "", "Database password", "DB_PASSWORD")

	srv, err := serv.New()
	if err != nil {
		console.Fatal(err)
	}

	go srv.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	srv.Close()
}
`

const modelService = `package module

import (
	"fmt"
	"net/http"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/middleware"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	v1 "$1/internal/service/$2/v1"
)

type Server struct {
	http *http.Server
}

func New() (*Server, error) {
	server := Server{}

	port := envar.EnvarInt(3300, "PORT")
	if port == 0 {
		return nil, fmt.Errorf("variable PORT es requerida")
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	latest := v1.New()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.HTTPError(w, r, http.StatusNotFound, "404 Not Found")
	})

	r.Mount("/", latest)
	r.Mount("/v1", latest)

	handler := cors.AllowAll().Handler(r)
	addr := strs.Format(":%d", port)
	serv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	server.http = serv

	return &server, nil
}

func (serv *Server) Close() {
	v1.Close()

	console.LogK("Http", "Shutting down server...")
}

func (serv *Server) StartHttpServer() {
	if serv.http == nil {
		return
	}

	svr := serv.http
	console.LogKF("Http", "Running on http://localhost%s", svr.Addr)
	console.Fatal(serv.http.ListenAndServe())
}

func (serv *Server) Start() {
	go serv.StartHttpServer()

	v1.Banner()

	<-make(chan struct{})
}
`

const modelDbApi = `package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/jrpc"
	"github.com/celsiainternet/elvis/utility"
	"github.com/dimiro1/banner"
	"github.com/go-chi/chi/v5"
	"github.com/mattn/go-colorable"
	pkg "$1/pkg/$2"	
)

func New() http.Handler {
	r := chi.NewRouter()

	err := pkg.LoadConfig()
	if err != nil {
		console.Panic(err)
	}

	_, err = cache.Load()
	if err != nil {
		console.Panic(err)
	}

	_, err = event.Load()
	if err != nil {
		console.Panic(err)
	}

	db, err := jdb.Load()
	if err != nil {
		console.Panic(err)
	}

	_pkg := &pkg.Router{
		Repository: &pkg.Controller{
			Db: db,
		},
	}

	r.Mount(pkg.PackagePath, _pkg.Routes())

	return r
}

func Close() {
	jrpc.Close()
	cache.Close()
	event.Close()
}

func Banner() {
	time.Sleep(3 * time.Second)
	templ := utility.BannerTitle(pkg.PackageName, 4)
	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
	fmt.Println()
}
`

const modelApi = `package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/jrpc"
	"github.com/celsiainternet/elvis/utility"
	"github.com/dimiro1/banner"
	"github.com/go-chi/chi/v5"
	"github.com/mattn/go-colorable"
	pkg "$1/pkg/$2"	
)

func New() http.Handler {
	r := chi.NewRouter()

	err := pkg.LoadConfig()
	if err != nil {
		console.Panic(err)
	}

	_, err = cache.Load()
	if err != nil {
		console.Panic(err)
	}

	_, err = event.Load()
	if err != nil {
		console.Panic(err)
	}

	db, err := jdb.Load()
	if err != nil {
		console.Panic(err)
	}

	_pkg := &pkg.Router{
		Repository: &pkg.Controller{
			Db: db,
		},
	}

	r.Mount(pkg.PackagePath, _pkg.Routes())

	return r
}

func Close() {
	jrpc.Close()
	cache.Close()
	event.Close()
}

func Banner() {
	time.Sleep(3 * time.Second)
	templ := utility.BannerTitle(pkg.PackageName, 4)
	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
	fmt.Println()
}
`

const modelEvent = `package $1

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/et"
)

func initEvents() {	
	err := event.Stack("<channel>", eventAction)
	if err != nil {
		console.Error(err)
	}

}

func eventAction(m event.EvenMessage) {
	data, err := et.ToJson(m.Data)
	if err != nil {
		console.Error(err)
	}

	console.Log("eventAction", data)
}
`

const modelModel = `package $1

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/jdb"
)

func initModels(db *jdb.DB) error {
	if err := Define$2(db); err != nil {
		return console.Panic(err)
	}

	return nil
}
`

const modelSchema = `package $1

import (
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"	
)

var $2 *linq.Schema

func defineSchema(db *jdb.DB) error {
	if $2 == nil {
		$2 = linq.NewSchema(db, "$3")
	}

	return nil
}
`

const modelhRpc = `package $1

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jrpc"
)

type Services struct{}

func StartRpcServer() {
	pkg, err := jrpc.Load(PackageName)
	if err != nil {
		console.Panic(err)
	}

	services := new(Services)
	err = jrpc.Mount(services)
	if err != nil {
		console.Fatal(err)
	}

	go pkg.Start()
}

func (c *Services) Version(require et.Json, response *et.Item) error {
	company := envar.EnvarStr("", "COMPANY")
	web := envar.EnvarStr("", "WEB")
	version := envar.EnvarStr("", "VERSION")
	help := envar.EnvarStr("", "RPC_HELP")
	response.Ok = true
	response.Result = et.Json{
		"methos":  "RPC",
		"version": version,
		"service": PackageName,
		"host":    HostName,
		"company": company,
		"web":     web,
		"help":    help,
	}

	return console.Rpc(response)
}

func (c *Services) Get$2ById(require et.Json, response *et.Item) error {
	id := require.Str("id")

	result, err := Get$2ById(id)
	if err != nil {
		return err
	}

	*response = result

	return console.Rpc(response)
}	
`

const modelMsg = `package $1

const (
	// MSG
	MSG_ATRIB_REQUIRED      = "Atributo requerido (%s)"
	MSG_VALUE_REQUIRED      = "Atributo requerido (%s) value:%s"
)
`

const modelConfig = `package $1

import (
	"github.com/celsiainternet/elvis/config"
	"github.com/celsiainternet/elvis/envar"
)

func LoadConfig() error {
	StartRpcServer()

	stage := envar.GetStr("local", "STAGE")
	return config.Load(stage, PackageName)
}
`

const modelDbController = `package $1

import (
	"context"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
)

type Controller struct {
	Db *jdb.DB
}

func (c *Controller) Version(ctx context.Context) (et.Json, error) {
	company := envar.GetStr("", "COMPANY")
	web := envar.GetStr("", "WEB")
	version := envar.EnvarStr("0.0.1", "VERSION")
  service := et.Json{
		"version": version,
		"service": PackageName,
		"host":    HostName,
		"company": company,
		"web":     web,
		"help":    "",
	}

	return service, nil
}

func (c *Controller) Init(ctx context.Context) {
	initModels(c.Db)
	initEvents()
}

type Repository interface {
	Version(ctx context.Context) (et.Json, error)
	Init(ctx context.Context)
}
`

const modelController = `package $1

import (
	"context"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
)

type Controller struct {
	Db *jdb.DB
}

func (c *Controller) Version(ctx context.Context) (et.Json, error) {
	company := envar.GetStr("", "COMPANY")
	web := envar.GetStr("", "WEB")
	version := envar.EnvarStr("0.0.1", "VERSION")
  service := et.Json{
		"version": version,
		"service": PackageName,
		"host":    HostName,
		"company": company,
		"web":     web,
		"help":    "",
	}

	return service, nil
}

func (c *Controller) Init(ctx context.Context) {
	initEvents()
}

type Repository interface {
	Version(ctx context.Context) (et.Json, error)
	Init(ctx context.Context)
}
`

const modelDbRouter = `package $1

import (
	"context"
	"net/http"
	"os"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/middleware"
	"github.com/celsiainternet/elvis/response"
	er "github.com/celsiainternet/elvis/router"
	"github.com/celsiainternet/elvis/strs"
	"github.com/go-chi/chi/v5"
)

var PackageName = "$1"
var PackageTitle = "$1"
var PackagePath = envar.GetStr("/api/$1", "PATH_URL")
var PackageVersion = envar.EnvarStr("0.0.1", "VERSION")
var HostName, _ = os.Hostname()

type Router struct {
	Repository Repository
}

func (rt *Router) Routes() http.Handler {
	defaultHost := strs.Format("http://%s", HostName)
	var host = strs.Format("%s:%d", envar.GetStr(defaultHost, "HOST"), envar.GetInt(3300, "PORT"))

	r := chi.NewRouter()

	er.PublicRoute(r, er.Get, "/version", rt.version, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/routes", rt.routes, PackageName, PackagePath, host)
	// $2
	er.ProtectRoute(r, er.Get, "/{id}", rt.get$2ById, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Post, "/", rt.upSert$2, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Put, "/state/{id}", rt.state$2, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Delete, "/{id}", rt.delete$2, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/", rt.all$2, PackageName, PackagePath, host)

	ctx := context.Background()
	rt.Repository.Init(ctx)
	middleware.SetServiceName(PackageName)

	console.LogKF(PackageName, "Router version:%s", PackageVersion)
	return r
}

func (rt *Router) version(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := rt.Repository.Version(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

func (rt *Router) routes(w http.ResponseWriter, r *http.Request) {
	_routes := er.GetRoutes()
	routes := []et.Json{}
	for _, route := range _routes {
		routes = append(routes, et.Json{
			"method": route.Str("method"),
			"path":   route.Str("path"),
		})
	}

	result := et.Items{
		Ok:     true,
		Count:  len(routes),
		Result: routes,
	}

	response.ITEMS(w, r, http.StatusOK, result)
}
`

const modelRouter = `package $1

import (
	"context"
	"net/http"
	"os"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/middleware"
	"github.com/celsiainternet/elvis/response"
	er "github.com/celsiainternet/elvis/router"
	"github.com/celsiainternet/elvis/strs"
	"github.com/go-chi/chi/v5"
)

var PackageName = "$1"
var PackageTitle = "$1"
var PackagePath = envar.GetStr("/api/$1", "PATH_URL")
var PackageVersion = envar.EnvarStr("0.0.1", "VERSION")
var HostName, _ = os.Hostname()

type Router struct {
	Repository Repository
}

func (rt *Router) Routes() http.Handler {
	defaultHost := strs.Format("http://%s", HostName)
	var host = strs.Format("%s:%d", envar.GetStr(defaultHost, "HOST"), envar.GetInt(3300, "PORT"))

	r := chi.NewRouter()

	er.PublicRoute(r, er.Get, "/version", rt.version, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/routes", rt.routes, PackageName, PackagePath, host)
	// $2
	er.ProtectRoute(r, er.Post, "/", rt.$2, PackageName, PackagePath, host)
	
	ctx := context.Background()
	rt.Repository.Init(ctx)
	middleware.SetServiceName(PackageName)

	console.LogKF(PackageName, "Router version:%s", PackageVersion)
	return r
}

func (rt *Router) version(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := rt.Repository.Version(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

func (rt *Router) routes(w http.ResponseWriter, r *http.Request) {
	_routes := er.GetRoutes()
	routes := []et.Json{}
	for _, route := range _routes {
		routes = append(routes, et.Json{
			"method": route.Str("method"),
			"path":   route.Str("path"),
		})
	}

	result := et.Items{
		Ok:     true,
		Count:  len(routes),
		Result: routes,
	}

	response.ITEMS(w, r, http.StatusOK, result)
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

const modelDbHandler = `package $1

import (
	"fmt"
	"net/http"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
	"github.com/go-chi/chi/v5"
)

var $2 *linq.Model

func Define$2(db *jdb.DB) error {
	if err := defineSchema(db); err != nil {
		return console.Panic(err)
	}

	if $2 != nil {
		return nil
	}

	$2 = linq.NewModel($3, "$4", "Tabla", 1)
	$2.DefineColum("created_at", "", "TIMESTAMP", "NOW()")
	$2.DefineColum("updated_at", "", "TIMESTAMP", "NOW()")
	$2.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	$2.DefineColum(jdb.KEY, "", "VARCHAR(80)", "-1")
	$2.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	$2.DefineColum("name", "", "VARCHAR(250)", "")
	$2.DefineColum("description", "", "TEXT", "")
	$2.DefineColum("_data", "", "JSONB", "{}")
	$2.DefineColum("index", "", "INTEGER", 0)
	$2.DefinePrimaryKey([]string{jdb.KEY})
	$2.DefineIndex([]string{
		"created_at",
		"updated_at",
		"_state",
		"project_id",
		"name",
		"index",
	})
	$2.DefineRequired([]string{
		"name:Atributo requerido (name)",
	})
	$2.IntegrityAtrib(true)
	$2.IndexSource(true)
	$2.Trigger(linq.BeforeInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	$2.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	$2.Trigger(linq.BeforeUpdate, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	$2.Trigger(linq.AfterUpdate, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	$2.Trigger(linq.BeforeDelete, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	$2.Trigger(linq.AfterDelete, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	$2.OnListener = func(data et.Json) {
		console.Debug(data.ToString())
	}
	
	if err := $2.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
*	Get$2ById
* @param id string
* @return et.Item, error
**/
func Get$2ById(id string) (et.Item, error) {
	result, err := $2.Data().
		Where($2.Column(jdb.KEY).Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil	
}

/**
* Insert$2
* @params project_id, id, name, description string
* @params data et.Json
* @return et.Item, error
**/
func Insert$2(project_id, id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidId(project_id) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "name")
	}

	if !utility.ValidId(id) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, jdb.KEY)
	}

	current, err := $2.Data().
		Where($2.Column(jdb.KEY).Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if current.Ok {
		return et.Item{Ok: false, Result: current.Result}, nil
	}

	id = utility.GenKey(id)
	now := utility.Now()
	data["created_at"] = now
	data["updated_at"] = now
	data["project_id"] = project_id
	data[jdb.KEY] = id
	data["name"] = name
	data["description"] = description
	item, err := $2.Insert(data).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* UpSert$2
* @param project_id string
* @param id string
* @param data et.Json
* @return et.Item, error
**/
func UpSert$2(project_id, id, name, description string, data et.Json) (et.Item, error) {
	current, err := Insert$2(project_id, id, name, description, data)
	if err != nil {
		return et.Item{}, err
	}
	
	if current.Ok {
		return current, nil
	}

	current_state := current.Key("_state")
	if current_state != utility.ACTIVE {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_UPDATE)
	}

	id = current.Str(jdb.KEY)
	now := utility.Now()
	data["updated_at"] = now
	data["project_id"] = project_id
	data[jdb.KEY] = id
	data["name"] = name
	data["description"] = description
	result, err := $2.Update(data).
		Where($2.Column(jdb.KEY).Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* State$2
* @param id, state string
* @return et.Item, error
**/
func State$2(id, state string) (et.Item, error) {
	if !utility.ValidId(state) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := $2.Data("_state").
		Where($2.Column(jdb.KEY).Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_FOUND)
	}

	old_state := current.Key("_state")
	if old_state == state {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_CHANGE)		
	}

	return $2.Update(et.Json{
		"_state":   state,
	}).
		Where($2.Column(jdb.KEY).Eq(id)).
		CommandOne()	
}

/**
* Delete$2
* @param id string
* @return et.Item, error
**/
func Delete$2(id string) (et.Item, error) {
	return State$2(id, utility.FOR_DELETE)
}

/**
* All$2
* @param project_id, state, search string
* @param page, rows int
* @param _select string
* @return et.List, error
**/
func All$2(project_id, state, search string, page, rows int, _select string) (et.List, error) {	
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return $2.Data(_select).
			Where($2.Column("project_id").In("-1", project_id)).
			And($2.Concat("NAME:", $2.Column("name"), "DESCRIPTION:", $2.Column("description"), "DATA:", $2.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy($2.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return $2.Data(_select).
			Where($2.Column("_state").Neg(state)).
			And($2.Column("project_id").In("-1", project_id)).
			OrderBy($2.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return $2.Data(_select).
			Where($2.Column("_state").In("-1", state)).
			And($2.Column("project_id").In("-1", project_id)).
			OrderBy($2.Column("name"), true).
			List(page, rows)
	} else {
		return $2.Data(_select).
			Where($2.Column("_state").Eq(state)).
			And($2.Column("project_id").In("-1", project_id)).
			OrderBy($2.Column("name"), true).
			List(page, rows)
	}
}

/**
* upSert$2
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) upSert$2(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	project_id := body.Str("project_id")
	id := body.Str("id")
	name := body.Str("name")
	description := body.Str("description")

	result, err := UpSert$2(project_id, id, name, description, body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* get$2ById
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) get$2ById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := Get$2ById(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* state$2
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) state$2(w http.ResponseWriter, r *http.Request) {
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

/**
* delete$2
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) delete$2(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := Delete$2(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* all$2
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) all$2(w http.ResponseWriter, r *http.Request) {
	query := response.GetQuery(r)
	project_id := query.Str("project_id")
	state := query.Str("state")
	search := query.Str("search")
	page := query.ValInt(1, "page")
	rows := query.ValInt(30, "rows")
	_select := query.Str("select")

	result, err := All$2(project_id, state, search, page, rows, _select)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

/** Copy this code to router.go
	// $2
	er.ProtectRoute(r, er.Get, "/$5/{id}", rt.get$2ById, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Post, "/$5", rt.upSert$2, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Put, "/$5/state/{id}", rt.state$2, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Delete, "/$5/{id}", rt.delete$2, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/$5/all", rt.all$2, PackageName, PackagePath, host)
**/

/** Copy this code to func initModel in model.go
	if err := Define$2(db); err != nil {
		return console.Panic(err)
	}
**/
`

const modelHandler = `package $1

import (
	"net/http"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/response"
)

func $2(project_id, id string, params et.Json) (et.Item, error) {

	return et.Item{}, nil
}


/**
* Router
**/
func (rt *Router) $3(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	project_id := body.Str("project_id")
	id := body.Str("id")	

	result, err := $2(project_id, id, body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/** Copy this code to router.go
	// $2
	er.ProtectRoute(r, er.Post, "/$3", rt.$2, PackageName, PackagePath, Host)	
**/
`

const modelReadme = `
## Project $1

## Create project

go mod init github.com/$1/api

### Dependencias

go get github.com/celsiainternet/elvis@v1.1.2

### Crear projecto, microservicios, modelos

go run github.com/celsiainternet/elvis/cmd/create-go create

### Run project

gofmt -w . && go run ./cmd/$1 -port 3400 -rpc 4400
`

const modelEnvar = `APP=
PORT=3300
VERSION=0.0.0
COMPANY=Company
PATH_URL=
WEB=
PRODUCTION=false
HOST=localhost

# DB
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=test
DB_USER=test
DB_PASSWORD=test
DB_APPLICATION_NAME=test

# REDIS
REDIS_HOST=localhost:6379
REDIS_PASSWORD=test
REDIS_DB=0

# NATS
NATS_HOST=localhost:4222

# CALM
SECRET=test

`

const modelDeploy = `version: "3"

networks:
  $3:
    external: true

services:
  $1:
    image: $1:latest
    logging:
      driver: "json-file"
      options:
        max-size: "1m"
        max-file: "2"
    networks:
      - $3
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.$1.rule=PathPrefix($2)"
      - "traefik.http.services.$1.loadbalancer.server.port=3300"
    deploy:
      replicas: 1
    environment:
      - "APP=Celsia Internet"
      - "PORT=3300"
      - "VERSION=1.0.1"
      - "COMPANY=Celsia Internet"
      - "WEB=https://www.app.com"
      - "PATH_URL=/api/$1"
      - "PRODUCTION=true"
      - "HOST=stack"
      # DB
      - "DB_DRIVER=postgres"
      - "DB_HOST="
      - "DB_PORT=5432"
      - "DB_NAME=internet"
      - "DB_USER=internet"
      - "DB_PASSWORD="
      - "DB_APPLICATION_NAME=$1"
      # REDIS
      - "REDIS_HOST="
      - "REDIS_PASSWORD="
      - "REDIS_DB=0"
      # NATS
      - "NATS_HOST=nats:4222"
      # CALM
      - "SECRET="
      # RPC
      - "PORT_RPC=4200"
`

const modelOkeStatefulSet = `# Definition del Service
apiVersion: v1
kind: Service
metadata:
  name: $ROLE
  namespace: $NS
spec:
  selector:
    role: $ROLE
  ports:
    - name: http
      protocol: TCP
      port: $PORT # Puerto en el Service
      targetPort: $PORT # Puerto en el Pod
    - name: rpc
      protocol: TCP
      port: 4200 # Puerto en el Service (RPC)
      targetPort: 4200 # Puerto en el Pod
  type: ClusterIP

---
# Definition del Deployment
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: $ROLE
  namespace: $NS
spec:
  replicas: $REPLICAS # Número de réplicas
  revisionHistoryLimit: $HISTORY_LIMIT # Limitar de revisiones
  selector:
    matchLabels:
      role: $ROLE # Selector que debe coincidir con los labels de los pods
  template:
    metadata:
      labels:
        role: $ROLE
    spec:
      imagePullSecrets:
        - name: dockerhub-secret
      containers:
        - name: $ROLE
          image: docker.io/$IMAGE
          imagePullPolicy: Always # Siempre descarga la imagen más reciente
          resources:
            requests:
              cpu: "$CPU_REQUEST" # Recursos mínimos requeridos
              memory: "$MEMORY_REQUEST"
            limits:
              cpu: "$CPU_LIMIT" # Recursos máximos permitidos
              memory: "$MEMORY_LIMIT"
          ports:
            - containerPort: $PORT # Puerto expuesto en el contenedor
          env:
            - name: APP
              value: "Celsia Internet"
            - name: PORT
              value: "$PORT"
            - name: VERSION
              value: "$RELEASE"
            - name: PRODUCTION
              value: "$PRODUCTION"
            - name: RPC_HOST
              value: "$ROLE"
            - name: RPC_PORT
              value: "4200"
            - name: PROJECT_ID
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: PROJECT_ID
            - name: PROJECT_NAME
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: PROJECT_NAME
            - name: RT_URL
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: RT_URL
            - name: WS_USERNAME
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: WS_USERNAME
            - name: WS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: WS_PASSWORD
            - name: USER_ADMIN
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: USER_ADMIN
            - name: USER_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: USER_PASSWORD
            - name: AUTHORIZATION_METHOD
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: AUTHORIZATION_METHOD
            - name: COMPANY
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: COMPANY
            - name: STAGE
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: STAGE
            - name: WEB
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: WEB
            - name: PATH_URL
              value: "$PATH_URL"
            - name: HOST
              value: "$HOST"
            - name: REQUESTS_LIMIT
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: REQUESTS_LIMIT
            - name: DEBUG
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DEBUG
            # DB
            - name: DB_DRIVER
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_DRIVER
            - name: DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_HOST
            - name: DB_PORT
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_PORT
            - name: DB_NAME
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_NAME
            - name: DB_USER
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_USER
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: DB_PASSWORD
            - name: DB_APPLICATION_NAME
              value: $ROLE
            # REDIS
            - name: REDIS_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: REDIS_HOST
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: REDIS_PASSWORD
            - name: REDIS_DB
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: REDIS_DB
            # NATS
            - name: NATS_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: NATS_HOST
            - name: NATS_USER
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: NATS_USER
            - name: NATS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: NATS_PASSWORD
            # SECRET
            - name: SECRET
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: SECRET
            # OSS DB
            - name: ORA_DB_HOST
              value: "$ORA_DB_HOST"
            - name: ORA_DB_PORT
              value: "$ORA_DB_PORT"
            - name: ORA_DB_NAME
              value: "$ORA_DB_NAME"
            - name: ORA_DB_USER
              value: "$ORA_DB_USER"
            - name: ORA_DB_PASSWORD
              value: "$ORA_DB_PASSWORD"
            - name: ORA_DB_SERVICE_NAME_ORACLE
              value: "$ORA_DB_SERVICE_NAME_ORACLE"
            - name: ORA_DB_SSL_ORACLE
              value: "$ORA_DB_SSL_ORACLE"
            - name: ORA_DB_SSL_VERIFY_ORACLE
              value: "$ORA_DB_SSL_VERIFY_ORACLE"
              # OSS
            - name: OSS_URL_ESECURITY_PROXY
              value: "$OSS_URL_ESECURITY_PROXY"
            - name: OSS_URL_LOGIN
              value: "$OSS_URL_LOGIN"
            - name: OSS_PASSWORD
              value: "$OSS_PASSWORD"
            - name: OSS_USER
              value: "$OSS_USER"
            # TokenCRM
            - name: CRM_TOKEN
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: CRM_TOKEN
            - name: CRM_URLCLIENT
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: CRM_URLCLIENT
            # Failure
            - name: PROJECTID_FAILURE
              value: "-1"
            - name: URL_API_FAILURE
              value: "https://api.celsiainternet.com"
            - name: TOKEN_API
              value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImU3ZDA1M2U2LWY1ZjQtNGRkZS1iYmIyLTMxZDU5YjFmZDE4OCIsImFwcCI6IiIsIm5hbWUiOiJJbnRlcmF4YSIsImtpbmQiOiJ0b2tlbiIsInVzZXJuYW1lIjoiIiwiZGV2aWNlIjoicmVxdWVzdHMiLCJkdXJhdGlvbiI6MH0.GmTc2VOpT2wgAu5NEqQwt78vUOuZi3WHYffUX9B19g8"
            # CAPTCHA
            - name: RECAPTCHA_KEY
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: RECAPTCHA_KEY
            - name: RECAPTCHA_URL
              value: "https://challenges.cloudflare.com/turnstile/v0/siteverify"
            - name: FOLDER_REPORT_REVIEW
              value: "prueba/updatedata/reporte"
`

const modelOkeDeployment = `# Definition del Service
apiVersion: v1
kind: Service
metadata:
  name: $ROLE
  namespace: $NS
spec:
  selector:
    role: $ROLE
  ports:
    - name: http
      protocol: TCP
      port: $PORT # Puerto en el Service
      targetPort: $PORT # Puerto en el Pod
    - name: rpc
      protocol: TCP
      port: 4200 # Puerto en el Service (RPC)
      targetPort: 4200 # Puerto en el Pod
  type: ClusterIP

---
# Definition del Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $ROLE
  namespace: $NS
spec:
  replicas: $REPLICAS # Número de réplicas
  revisionHistoryLimit: $HISTORY_LIMIT # Limitar de revisiones
  selector:
    matchLabels:
      role: $ROLE # Selector que debe coincidir con los labels de los pods
  template:
    metadata:
      labels:
        role: $ROLE
    spec:
      imagePullSecrets:
        - name: dockerhub-secret
      containers:
        - name: $ROLE
          image: docker.io/$IMAGE
          imagePullPolicy: Always # Siempre descarga la imagen más reciente
          resources:
            requests:
              cpu: "$CPU_REQUEST" # Recursos mínimos requeridos
              memory: "$MEMORY_REQUEST"
            limits:
              cpu: "$CPU_LIMIT" # Recursos máximos permitidos
              memory: "$MEMORY_LIMIT"
          ports:
            - containerPort: $PORT # Puerto expuesto en el contenedor
          env:
            - name: APP
              value: "Celsia Internet"
            - name: PORT
              value: "$PORT"
            - name: VERSION
              value: "$RELEASE"
            - name: PRODUCTION
              value: "$PRODUCTION"
            - name: RPC_HOST
              value: "$ROLE"
            - name: RPC_PORT
              value: "4200"
            - name: PROJECT_ID
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: PROJECT_ID
            - name: PROJECT_NAME
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: PROJECT_NAME
            - name: RT_URL
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: RT_URL
            - name: WS_USERNAME
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: WS_USERNAME
            - name: WS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: WS_PASSWORD
            - name: USER_ADMIN
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: USER_ADMIN
            - name: USER_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: USER_PASSWORD
            - name: AUTHORIZATION_METHOD
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: AUTHORIZATION_METHOD
            - name: COMPANY
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: COMPANY
            - name: STAGE
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: STAGE
            - name: WEB
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: WEB
            - name: PATH_URL
              value: "$PATH_URL"
            - name: HOST
              value: "$HOST"
            - name: REQUESTS_LIMIT
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: REQUESTS_LIMIT
            - name: DEBUG
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DEBUG
            # DB
            - name: DB_DRIVER
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_DRIVER
            - name: DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_HOST
            - name: DB_PORT
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_PORT
            - name: DB_NAME
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_NAME
            - name: DB_USER
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: DB_USER
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: DB_PASSWORD
            - name: DB_APPLICATION_NAME
              value: $ROLE
            # REDIS
            - name: REDIS_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: REDIS_HOST
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: REDIS_PASSWORD
            - name: REDIS_DB
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: REDIS_DB
            # NATS
            - name: NATS_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: NATS_HOST
            - name: NATS_USER
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: NATS_USER
            - name: NATS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: NATS_PASSWORD
            # CAPTCHA
            - name: RECAPTCHA_KEY
              value: "0x4AAAAAABtV1cqhc94ExTeOCdTd4EcS4UU"
            - name: RECAPTCHA_URL
              value: "https://challenges.cloudflare.com/turnstile/v0/siteverify"
            # SECRET
            - name: SECRET
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: SECRET
            # OSS DB
            - name: ORA_DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: ORA_DB_HOST
            - name: ORA_DB_PORT
              value: "1521"
            - name: ORA_DB_NAME
              value: ""
            - name: ORA_DB_USER
              value: "emalaeaprd"
            - name: ORA_DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: ORA_DB_PASSWORD
            - name: ORA_DB_SERVICE_NAME_ORACLE
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: ORA_DB_SERVICES
            - name: ORA_DB_SSL_ORACLE
              value: "true"
            - name: ORA_DB_SSL_VERIFY_ORACLE
              value: "false"
            # OSS
            - name: OSS_URL_ESECURITY_PROXY
              value: "$OSS_URL_ESECURITY_PROXY"
            - name: OSS_URL_LOGIN
              value: "$OSS_URL_LOGIN"
            - name: OSS_PASSWORD
              value: "$OSS_PASSWORD"
            - name: OSS_USER
              value: "$OSS_USER"
            # TokenCRM
            - name: CRM_TOKEN
              valueFrom:
                secretKeyRef:
                  name: secret-$NS
                  key: CRM_TOKEN
            - name: CRM_URLCLIENT
              valueFrom:
                configMapKeyRef:
                  name: config-$NS
                  key: CRM_URLCLIENT
            # Failure
            - name: PROJECTID_FAILURE
              value: "-1"
            - name: URL_API_FAILURE
              value: "https://api.celsiainternet.com"
            - name: TOKEN_API
              value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImU3ZDA1M2U2LWY1ZjQtNGRkZS1iYmIyLTMxZDU5YjFmZDE4OCIsImFwcCI6IiIsIm5hbWUiOiJJbnRlcmF4YSIsImtpbmQiOiJ0b2tlbiIsInVzZXJuYW1lIjoiIiwiZGV2aWNlIjoicmVxdWVzdHMiLCJkdXJhdGlvbiI6MH0.GmTc2VOpT2wgAu5NEqQwt78vUOuZi3WHYffUX9B19g8"
            - name: FOLDER_REPORT_REVIEW
              value: "updatedata/reporte"

  strategy:
    type: RollingUpdate # Actualización gradual (opcional)
    rollingUpdate:
      maxUnavailable: $MAX_PODS_AVAILABLE # Número máximo de pods no disponibles durante la actualización
      maxSurge: $MAX_PODS_SURGE # Número máximo de pods adicionales creados durante la actualización
`

const modelGitignore = `# Created by https://www.toptal.com/developers/gitignore/api/go
# Edit at https://www.toptal.com/developers/gitignore?templates=go

### Go ###
# If you prefer the allow list template instead of the deny list, see community template:
# https://github.com/github/gitignore/blob/main/community/Golang/Go.AllowList.gitignore
#
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
.env
data
build
sql
.vscode
deployments/oke.yml

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# Credencial acces token to AWS server
*.pem`
