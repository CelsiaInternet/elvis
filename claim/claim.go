package claim

import (
	"net/http"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/golang-jwt/jwt/v4"
)

type Claim struct {
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
* TokenKey
* @param app string
* @param device string
* @param id string
* @return string
**/
func TokenKey(app, device, id string) string {
	str := strs.Append(app, device, "-")
	str = strs.Append(str, id, "-")
	return strs.Format(`token:%s`, str)
}

/**
* NewToken
* @param id string
* @param app string
* @param name string
* @param kind string
* @param username string
* @param device string
* @param duration time.Duration
* @return token string
* @return key string
* @return err error
**/
func NewToken(id, app, name, kind, username, device string, duration time.Duration) (string, error) {
	secret := envar.GetStr("1977", "SECRET")
	c := Claim{}
	c.ID = id
	c.App = app
	c.Name = name
	c.Kind = kind
	c.Username = username
	c.Device = device
	c.Duration = duration
	_jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := _jwt.SignedString([]byte(secret))
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
	key := TokenKey(app, device, id)
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
	secret := envar.GetStr("1977", "SECRET")
	jToken, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	claim, ok := jToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	app, ok := claim["app"].(string)
	if !ok {
		return nil
	}

	id, ok := claim["id"].(string)
	if !ok {
		return nil
	}

	device, ok := claim["device"].(string)
	if !ok {
		return nil
	}

	return DeleteToken(app, device, id)
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
		return nil, err
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

	key := TokenKey(result.App, result.Device, result.ID)
	c, err := cache.Get(key, "")
	if err != nil {
		return nil, console.Alert(MSG_TOKEN_INVALID)
	}

	if c != token {
		return nil, console.Alert(MSG_TOKEN_INVALID)
	}

	err = cache.Set(key, c, result.Duration)
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
func SetToken(app, device, id, token string) error {
	key := TokenKey(app, device, id)
	err := cache.Set(key, token, 0)
	if err != nil {
		return err
	}

	return nil
}

/**
* GetClient
* @param r *http.Request
* @return et.Json
**/
func GetClient(r *http.Request) et.Json {
	now := utility.Now()
	ctx := r.Context()

	return et.Json{
		"date_of":   now,
		"client_id": et.NewAny(ctx.Value("clientId")).Str(),
		"username":  et.NewAny(ctx.Value("username")).Str(),
		"name":      et.NewAny(ctx.Value("name")).Str(),
	}
}
