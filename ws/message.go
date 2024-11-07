package ws

import (
	"encoding/json"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
)

type TpMessage int

const (
	TpPing           TpMessage = iota // 0
	TpSetFrom                         // 1
	TpSubscribe                       // 2
	TpQueueSubscribe                  // 3
	TpStack                           // 4
	TpUnsubscribe                     // 5
	TpPublish                         // 6
	TpDirect                          // 7
	TpConnect                         // 8
	TpDisconnect                      // 9
)

func (s TpMessage) String() string {
	switch s {
	case TpPing:
		return "Ping"
	case TpSetFrom:
		return "Set id and name"
	case TpSubscribe:
		return "Subscribe"
	case TpQueueSubscribe:
		return "Queue"
	case TpStack:
		return "Stack"
	case TpUnsubscribe:
		return "Unsubscribe"
	case TpPublish:
		return "Publish"
	case TpDirect:
		return "Direct"
	case TpConnect:
		return "Connect"
	case TpDisconnect:
		return "Disconnect"
	default:
		return "Unknown"
	}
}

func ToTpMessage(s string) TpMessage {
	switch s {
	case "Ping":
		return TpPing
	case "SetFrom":
		return TpSetFrom
	case "Subscribe":
		return TpSubscribe
	case "Queue":
		return TpQueueSubscribe
	case "Stack":
		return TpStack
	case "Unsubscribe":
		return TpUnsubscribe
	case "Publish":
		return TpPublish
	case "Direct":
		return TpDirect
	default:
		return -1
	}
}

type Message struct {
	Created_at time.Time   `json:"created_at"`
	Id         string      `json:"id"`
	From       et.Json     `json:"from"`
	To         string      `json:"to"`
	Ignored    []string    `json:"ignored"`
	Data       interface{} `json:"data"`
	Tp         TpMessage   `json:"tp"`
	Channel    string      `json:"channel"`
	Queue      string      `json:"queue"`
}

/**
* NewMessage
* @param et.Json
* @param interface{}
* @param TpMessage
* @return Message
**/
func NewMessage(from et.Json, message interface{}, tp TpMessage) Message {
	return Message{
		Created_at: timezone.NowTime(),
		Id:         utility.UUID(),
		From:       from,
		Ignored:    []string{},
		Data:       message,
		Tp:         tp,
	}
}

/**
* Encode return the message as byte array
* @return []byte
**/
func (e Message) Encode() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/**
* ToJson return the message as et.Json
* @return et.Json
**/
func (e Message) ToJson() et.Json {
	result, err := et.Object(e)
	if err != nil {
		return et.Json{}
	}

	result["tp"] = et.Json{
		"code": e.Tp,
		"name": e.Tp.String(),
	}

	return result
}

/**
* ToString return the message as string
* @return string
**/
func (e Message) ToString() string {
	result := e.ToJson()

	return result.ToString()
}

/**
* DecodeMessage
* @param []byte
* @return Message
**/
func DecodeMessage(data []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		return Message{}, err
	}

	return m, nil
}
