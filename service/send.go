package service

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/utility"
)

/**
* SendSms
* @param project_id string, service_id string, contactNumbers []string, content string, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendSms(project_id, service_id string, contactNumbers []string, content string, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id = reg.GetUUID(service_id)
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
* @param project_id string, service_id string, template_id int, contactNumbers []string, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendWhatsapp(project_id, service_id string, template_id int, contactNumbers []string, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id = reg.GetUUID(service_id)
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
* @param project_id string, service_id string, to []et.Json, subject string, html_content string, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendEmail(project_id, service_id string, to []et.Json, subject string, html_content string, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id = reg.GetUUID(service_id)
	result := event.Work("send/email", et.Json{
		"project_id":   project_id,
		"service_id":   service_id,
		"to":           to,
		"subject":      subject,
		"html_content": html_content,
		"params":       params,
		"type":         tp.String(),
		"client_id":    clientId,
	})

	result["service_id"] = service_id
	return result
}

/**
* SendEmailByTemplate
* @param project_id string, service_id string, groups_id []string, subject string, template_id int, params []et.Json, tp TpMessage, clientId string
* @response et.Json
**/
func SendEmailByTemplate(project_id, service_id string, groups_id []string, subject, template_id string, params []et.Json, tp TpMessage, clientId string) et.Json {
	service_id = reg.GetUUID(service_id)
	result := event.Work("send/emailbytemplate", et.Json{
		"project_id":  project_id,
		"service_id":  service_id,
		"groups_id":   groups_id,
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
