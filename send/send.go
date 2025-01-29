package send

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/strs"
)

func SMS(serviceId string, data et.Json) (et.Item, error) {

	channel := strs.Format("send/sms")
	event.Work(channel, data)
	message := strs.Format("SMS sent to service %s", serviceId)

	return et.Item{
		Ok: true,
		Result: et.Json{
			"service_id": serviceId,
			"message":    message,
		},
	}, nil
}
