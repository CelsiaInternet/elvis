package ws

import (
	"encoding/json"
	"time"

	"github.com/cgalvisleon/elvis/et"
	m "github.com/cgalvisleon/elvis/message"
	"github.com/cgalvisleon/elvis/utility"
)

type Message struct {
	Created_at time.Time `json:"created_at"`
	Id         string    `json:"id"`
	From       et.Json   `json:"from"`
	to         string
	Ignored    []string    `json:"ignored"`
	Tp         m.TpMessage `json:"tp"`
	Channel    string      `json:"channel"`
	Queue      string      `json:"queue"`
	Data       interface{} `json:"data"`
}

/**
* NewMessage
* @param et.Json
* @param interface{}
* @param m.TpMessage
* @return Message
**/
func NewMessage(from et.Json, message interface{}, tp m.TpMessage) Message {
	return Message{
		Created_at: time.Now(),
		Id:         utility.UUID(),
		From:       from,
		Data:       message,
		Tp:         tp,
		Ignored:    []string{},
	}
}

/**
* Type
* @return m.TpMessage
**/
func (e Message) Type() m.TpMessage {
	return e.Tp
}

/**
* ToString
* @return string
**/
func (e Message) ToString() string {
	j, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(j)
}

/**
* Encode
* @return []byte
* @return error
**/
func (e Message) Encode() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/**
* Json
* @return et.Json
* @return error
**/
func (e Message) Json() (et.Json, error) {
	result := et.Json{}
	err := result.Scan(e.Data)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* DecodeMessage
* @param []byte
* @return Message
* @return error
**/
func DecodeMessage(data []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		return Message{}, err
	}

	return m, nil
}
