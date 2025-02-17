package service

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

/**
* SendSms
* @param project_id, service_id string, contactNumbers []string, content string, params []et.Json, tp string, client et.Json
* @response et.Json
**/
func SendSms(project_id, service_id string, contactNumbers []string, content string, params []et.Json, tp string, client et.Json) et.Json {
	return event.Work("send/sms", et.Json{
		"project_id":      project_id,
		"service_id":      service_id,
		"contact_numbers": contactNumbers,
		"content":         content,
		"params":          params,
		"type":            tp,
		"client":          client,
	})
}

/**
* SendWhatsapp
* @param project_id, service_id string, template_id int, contactNumbers []string, params []et.Json, tp string, client et.Json
* @response et.Json
**/
func SendWhatsapp(project_id, service_id string, template_id int, contactNumbers []string, params []et.Json, tp string, client et.Json) et.Json {
	return event.Work("send/whatsapp", et.Json{
		"project_id":      project_id,
		"service_id":      service_id,
		"template_id":     template_id,
		"contact_numbers": contactNumbers,
		"params":          params,
		"type":            tp,
		"client":          client,
	})
}

/**
* SendEmail
* @param project_id, service_id string, to []et.Json, subject string, html_content string, params []et.Json, tp string, client et.Json
* @response et.Json
**/
func SendEmail(project_id, service_id string, to []et.Json, subject string, html_content string, params []et.Json, tp string, client et.Json) et.Json {
	return event.Work("send/email", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"to":         to,
		"subject":    subject,
		"content":    html_content,
		"params":     params,
		"type":       tp,
		"client":     client,
	})
}

/**
* SendEmailByTemplate
* @param project_id, service_id string, to []et.Json, subject string, template_id int, params []et.Json, tp string, client et.Json
* @response et.Json
**/
func SendEmailByTemplate(project_id, service_id string, to []et.Json, subject string, template_id int, params []et.Json, tp string, client et.Json) et.Json {
	return event.Work("send/email/template", et.Json{
		"project_id":  project_id,
		"service_id":  service_id,
		"to":          to,
		"subject":     subject,
		"template_id": template_id,
		"params":      params,
		"type":        tp,
		"client":      client,
	})
}

/**
* GenerateOtp
* @param project_id, service_id, channel, name, device string, length, duration int, client et.Json
* @response et.Json
**/
func GenerateOtpSms(project_id, service_id, phone_number, name, device string, length int, duration int, client et.Json) et.Json {
	return event.Work("generate/otp/sms", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    phone_number,
		"device":     device,
		"kind":       "sms",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"client":     client,
	})
}

/**
* GenerateOtpWhatsapp
* @param project_id, service_id, phone_number, name, device string, length, duration int, client et.Json
* @response et.Json
**/
func GenerateOtpWhatsapp(project_id, service_id, phone_number, name, device string, length int, duration int, client et.Json) et.Json {
	return event.Work("generate/otp/whatsapp", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    phone_number,
		"device":     device,
		"kind":       "whatsapp",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"client":     client,
	})
}

/**
* GenerateOtpEmail
* @param project_id, service_id, email, name, device string, length, duration int, client et.Json
* @response et.Json
**/
func GenerateOtpEmail(project_id, service_id, email, name, device string, length int, duration int, client et.Json) et.Json {
	return event.Work("generate/otp/email", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    email,
		"device":     device,
		"kind":       "email",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"client":     client,
	})
}

/**
* VerifyOtp
* @param channel, device, kind, code string, user et.Json
* @response et.Item, error
**/
func VerifyOtp(channel, device, kind, code string, user et.Json) (et.Item, error) {
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
		return et.Item{}, console.Alert(msg.MSG_CODE_UNVERIFY)
	}

	cache.DeleteVerify(device, key)

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": msg.MSG_EXECUTED_SUCCESS,
		},
	}, err
}
