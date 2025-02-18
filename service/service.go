package service

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

// Type message
type TpMessage int

const (
	TpTransactional TpMessage = iota
	TpComercial
)

func (tp TpMessage) String() string {
	switch tp {
	case TpTransactional:
		return "transactional"
	case TpComercial:
		return "comercial"
	default:
		return "unknown"
	}
}

/**
* GetId
* @param client_id, kind, description string
* @response string
**/
func GetId(client_id, kind, description string) string {
	now := utility.Now()
	result := utility.UUID()
	data := et.Json{
		"created_at":  now,
		"service_id":  result,
		"client_id":   client_id,
		"kind":        kind,
		"description": description,
	}
	event.Work("service/client", data)
	cache.SetH(result, data)

	return result
}

/**
* SetStatus
* @param serviceId string, status et.Json
**/
func SetStatus(serviceId string, status et.Json) {
	event.Work("service/status", et.Json{
		"service_id": serviceId,
		"status":     status,
	})

	cache.SetH(serviceId, status)
}

/**
* GetStatus
* @param serviceId string
* @response et.Json, error
**/
func GetStatus(serviceId string) (et.Json, error) {
	return cache.GetJson(serviceId)
}

/**
* SendSms
* @param project_id string, contactNumbers []string, content string, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendSms(project_id string, contactNumbers []string, content string, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id := GetId(clientId, "sms", "Send SMS")
	result := event.Work("send/sms", et.Json{
		"project_id":      project_id,
		"service_id":      service_id,
		"contact_numbers": contactNumbers,
		"content":         content,
		"params":          params,
		"type":            tp.String(),
		"client_id":       clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* SendWhatsapp
* @param project_id string, template_id int, contactNumbers []string, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendWhatsapp(project_id string, template_id int, contactNumbers []string, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id := GetId(clientId, "whatsapp", "Send Whatsapp")
	result := event.Work("send/whatsapp", et.Json{
		"project_id":      project_id,
		"service_id":      service_id,
		"template_id":     template_id,
		"contact_numbers": contactNumbers,
		"params":          params,
		"type":            tp.String(),
		"client_id":       clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* SendEmail
* @param project_id string, to []et.Json, subject string, html_content string, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendEmail(project_id string, to []et.Json, subject string, html_content string, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id := GetId(clientId, "email", "Send Email")
	result := event.Work("send/email", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"to":         to,
		"subject":    subject,
		"content":    html_content,
		"params":     params,
		"type":       tp.String(),
		"client_id":  clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* SendEmailByTemplate
* @param project_id string, to []et.Json, subject string, template_id int, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendEmailByTemplate(project_id string, to []et.Json, subject string, template_id int, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id := GetId(clientId, "email", "Send Email By Template")
	result := event.Work("send/email/template", et.Json{
		"project_id":  project_id,
		"service_id":  service_id,
		"to":          to,
		"subject":     subject,
		"template_id": template_id,
		"params":      params,
		"type":        tp.String(),
		"client_id":   clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* SendOtp
* @param project_id string, channel, name, device string, length, duration int, clientId string
* @response et.Json
**/
func SendOtpSms(project_id string, phone_number, name, device string, length int, duration int, clientId string) et.Json {
	service_id := GetId(clientId, "sms otp", "Send OTP by SMS")
	result := event.Work("generate/otp/sms", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    phone_number,
		"device":     device,
		"kind":       "sms",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"client_id":  clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* SendOtpWhatsapp
* @param project_id string, phone_number, name, device string, length, duration int, clientId string
* @response et.Json
**/
func SendOtpWhatsapp(project_id string, phone_number, name, device string, length int, duration int, clientId string) et.Json {
	service_id := GetId(clientId, "whatsapp otp", "Send OTP by Whatsapp")
	result := event.Work("generate/otp/whatsapp", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    phone_number,
		"device":     device,
		"kind":       "whatsapp",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"client_id":  clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* SendOtpEmail
* @param project_id string, email, name, device string, length, duration int, clientId string
* @response et.Json
**/
func SendOtpEmail(project_id string, email, name, device string, length int, duration int, clientId string) et.Json {
	service_id := GetId(clientId, "email otp", "Send OTP by Email")
	result := event.Work("generate/otp/email", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    email,
		"device":     device,
		"kind":       "email",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"client_id":  clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* VerifyOtp
* @param channel, device, kind, code, clientId string
* @response et.Item, error
**/
func VerifyOtp(channel, device, kind, code, clientId string) (et.Item, error) {
	if !utility.ValidStr(channel, 1, []string{}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "channel")
	}

	if !utility.ValidStr(kind, 0, []string{}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "kind")
	}

	if !utility.ValidStr(device, 1, []string{}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "device")
	}

	if !utility.ValidStr(code, 1, []string{}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "code")
	}

	var kinds map[string]bool = map[string]bool{
		"sms":      true,
		"email":    true,
		"whatsapp": true,
	}
	if !kinds[kind] {
		return et.Item{}, console.Alert(msg.MSG_VALIDATE_KIND)
	}

	key := kind + channel
	codeVerify, err := cache.GetVerify(device, key)
	if codeVerify != code {
		GetId(clientId, "otp denied", "Verify OTP denied")
		return et.Item{}, console.Alert(msg.MSG_CODE_UNVERIFY)
	}

	GetId(clientId, "otp", "Verify OTP success")
	cache.DeleteVerify(device, key)

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": msg.MSG_EXECUTED_SUCCESS,
		},
	}, err
}
