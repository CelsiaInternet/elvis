package claim

import (
	"context"
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/golang-jwt/jwt/v4"
)

type ContextKey string

func (c ContextKey) String(ctx context.Context, def string) string {
	val := ctx.Value(c)
	result, ok := val.(string)
	if !ok {
		return def
	}

	return result
}

const (
	ServiceIdKey ContextKey = "serviceId"
	ClientIdKey  ContextKey = "clientId"
	AppKey       ContextKey = "app"
	NameKey      ContextKey = "name"
	KindKey      ContextKey = "kind"
	UsernameKey  ContextKey = "username"
	TokenKey     ContextKey = "token"
	ProjectIdKey ContextKey = "projectId"
	ProfileTpKey ContextKey = "profileTp"
	ModelKey     ContextKey = "model"
)

type AuthType string

const (
	BasicAuth   AuthType = "BasicAuth"
	BearerToken AuthType = "BearerToken"
)

type Claim struct {
	Salt     string        `json:"salt"`
	ID       string        `json:"id"`
	App      string        `json:"app"`
	Name     string        `json:"name"`
	Kind     string        `json:"kind"`
	Username string        `json:"username"`
	Device   string        `json:"device"`
	Duration time.Duration `json:"duration"`
	jwt.StandardClaims
}

/**
* ToJson
* @return et.Json
**/
func (c *Claim) ToJson() et.Json {
	return et.Json{
		"id":       c.ID,
		"app":      c.App,
		"name":     c.Name,
		"kind":     c.Kind,
		"username": c.Username,
		"device":   c.Device,
		"duration": c.Duration,
	}
}

/**
* GetTokenKey
* @param app string
* @param device string
* @param id string
* @return string
**/
func GetTokenKey(app, device, id string) string {
	str := strs.Append(app, device, "-")
	str = strs.Append(str, id, "-")
	str = strs.Format(`token:%s`, str)
	return utility.ToBase64(str)
}

/**
* NewToken
* @param id string
* @param app string
* @param name string
* @param kind AuthType
* @param username string
* @param device string
* @param duration time.Duration
* @return token string
* @return key string
* @return err error
**/
func NewToken(id, app, name string, kind AuthType, username, device string, duration time.Duration) (string, error) {
	secret := envar.GetStr("1977", "SECRET")
	c := Claim{}
	c.Salt = utility.GetOTP(6)
	c.ID = id
	c.App = app
	c.Name = name
	c.Kind = string(kind)
	c.Username = username
	c.Device = device
	c.Duration = duration
	_jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := _jwt.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	key := GetTokenKey(c.App, c.Device, c.ID)
	err = cache.Set(key, token, c.Duration)
	if err != nil {
		return "", err
	}

	return token, nil
}

/**
* DeleteToken
* @param app string
* @param device string
* @param id string
* @return error
**/
func DeleteToken(app, device, id string) error {
	key := GetTokenKey(app, device, id)
	_, err := cache.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

/**
* DeleteTokeByToken
* @param token string
* @return error
**/
func DeleteTokeByToken(token string) error {
	claim, err := ParceToken(token)
	if err != nil {
		return err
	}

	return DeleteToken(claim.App, claim.Device, claim.ID)
}

/**
* ParceToken
* @param token string
* @return *Claim
* @return error
**/
func ParceToken(token string) (*Claim, error) {
	secret := envar.GetStr("1977", "SECRET")
	jToken, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, console.Error(err)
	}

	if !jToken.Valid {
		return nil, console.Alert(MSG_TOKEN_INVALID)
	}

	claim, ok := jToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, console.Alert(MSG_REQUIRED_INVALID)
	}

	app, ok := claim["app"].(string)
	if !ok {
		return nil, console.Alert(ERR_INVALID_CLAIM)
	}

	id, ok := claim["id"].(string)
	if !ok {
		return nil, console.Alert(ERR_INVALID_CLAIM)
	}

	name, ok := claim["name"].(string)
	if !ok {
		return nil, console.Alert(ERR_INVALID_CLAIM)
	}

	kind, ok := claim["kind"].(string)
	if !ok {
		return nil, console.Alert(ERR_INVALID_CLAIM)
	}

	username, ok := claim["username"].(string)
	if !ok {
		return nil, console.Alert(ERR_INVALID_CLAIM)
	}

	device, ok := claim["device"].(string)
	if !ok {
		return nil, console.AlertF(MSG_TOKEN_INVALID_ATRIB, "device")
	}

	second, ok := claim["duration"].(float64)
	if !ok {
		return nil, console.AlertF(MSG_TOKEN_INVALID_ATRIB, "duration")
	}

	duration := time.Duration(second)

	return &Claim{
		ID:       id,
		App:      app,
		Name:     name,
		Kind:     kind,
		Username: username,
		Device:   device,
		Duration: duration,
	}, nil
}

/**
* GetFromToken
* @param token string
* @return *Claim
* @return error
**/
func GetFromToken(token string) (*Claim, error) {
	result, err := ParceToken(token)
	if err != nil {
		return nil, err
	}

	key := GetTokenKey(result.App, result.Device, result.ID)
	val, err := cache.Get(key, "")
	if err != nil {
		return nil, err
	}

	if val != token {
		return nil, err
	}

	err = cache.Set(key, token, result.Duration)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* SetToken
* @param app string
* @param device string
* @param id string
* @param token string
* @return error
**/
func SetToken(app, device, id, token string, duration time.Duration) error {
	key := GetTokenKey(app, device, id)
	err := cache.Set(key, token, duration)
	if err != nil {
		return err
	}

	return nil
}

/**
* GetUser
* @param r *http.Request
* @return et.Json
**/
func GetUser(r *http.Request) et.Json {
	now := utility.Now()
	ctx := r.Context()
	username := UsernameKey.String(ctx, "Anonimo")

	return et.Json{
		"date_of":  now,
		"username": username,
	}
}
