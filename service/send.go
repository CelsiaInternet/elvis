package service

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/utility"
)

/**
* SendSms
* @param project_id string, service_id string, contactNumbers []string, content string, params []et.Json, tp TpMessage, user et.Json
* @return et.Json
**/
func SendSms(project_id, service_id string, contactNumbers []string, content string, params []et.Json, tp TpMessage, user et.Json) (et.Item, error) {
	service_id = reg.GetUUID(service_id)
	return sendRpc("SendSMS", "send/sms", et.Json{
		"project_id":      project_id,
		"service_id":      service_id,
		"contact_numbers": contactNumbers,
		"content":         content,
		"params":          params,
		"type":            tp.String(),
		"user":            user,
	})
}

/**
* SendWhatsapp
* @param project_id string, service_id string, template_id int, contactNumbers []string, params []et.Json, tp TpMessage, user et.Json
* @return et.Json
**/
func SendWhatsapp(project_id, service_id string, template_id int, contactNumbers []string, params []et.Json, tp TpMessage, user et.Json) (et.Item, error) {
	service_id = reg.GetUUID(service_id)
	return sendRpc("SendWhatsapp", "send/whatsapp", et.Json{
		"project_id":      project_id,
		"service_id":      service_id,
		"template_id":     template_id,
		"contact_numbers": contactNumbers,
		"params":          params,
		"type":            tp.String(),
		"user":            user,
	})
}

/**
* SendEmail
* @param project_id string, service_id string, to []et.Json, subject string, html_content string, params []et.Json, tp TpMessage, user et.Json
* @return et.Json
**/
func SendEmail(project_id, service_id, smtp string, to []et.Json, subject string, html_content string, params []et.Json, tp TpMessage, user et.Json) (et.Item, error) {
	service_id = reg.GetUUID(service_id)
	return sendRpc("SendEmail", "send/email", et.Json{
		"project_id":   project_id,
		"service_id":   service_id,
		"smtp":         smtp,
		"to":           to,
		"subject":      subject,
		"html_content": html_content,
		"params":       params,
		"type":         tp.String(),
		"user":         user,
	})
}

/**
* SendEmailByTemplate
* @param project_id string, service_id string, groups []string, subject string, template_id int, params []et.Json, tp TpMessage, user et.Json
* @return et.Json
**/
func SendEmailByTemplate(project_id, service_id, smtp string, groups []string, subject, template_id string, params []et.Json, tp TpMessage, user et.Json) (et.Item, error) {
	service_id = reg.GetUUID(service_id)
	return sendRpc("SendEmailByTemplate", "send/emailbytemplate", et.Json{
		"project_id":  project_id,
		"service_id":  service_id,
		"smtp":        smtp,
		"groups":      groups,
		"subject":     subject,
		"template_id": template_id,
		"params":      params,
		"type":        tp.String(),
		"user":        user,
	})
}

/**
* SendEmailToByTemplate
* @param project_id, service_id, name, email string, subject, template_id string, params []et.Json, tp TpMessage, user et.Json
* @return et.Json
**/
func SendEmailToByTemplate(project_id, service_id, smtp, name, email string, subject, template_id string, params []et.Json, tp TpMessage, user et.Json) (et.Item, error) {
	service_id = reg.GetUUID(service_id)
	return sendRpc("SendEmailToByTemplate", "send/emailtobytemplate", et.Json{
		"project_id":  project_id,
		"service_id":  service_id,
		"smtp":        smtp,
		"name":        name,
		"email":       email,
		"subject":     subject,
		"template_id": template_id,
		"params":      params,
		"type":        tp.String(),
		"user":        user,
	})
}

/**
* SendOtp
* @param project_id string, channel, name, device string, length, duration int, user et.Json
* @return et.Json
**/
func SendOtpSms(project_id string, phone_number, name, device string, length int, duration int, user et.Json) (et.Item, error) {
	service_id := reg.UUID()
	return otpRpc("GenerateOTPSMS", "generate/otp/sms", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    phone_number,
		"device":     device,
		"kind":       "sms",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"user":       user,
	})
}

/**
* SendOtpWhatsapp
* @param project_id string, phone_number, name, device string, length, duration int, user et.Json
* @return et.Json
**/
func SendOtpWhatsapp(project_id string, phone_number, name, device string, length int, duration int, user et.Json) (et.Item, error) {
	service_id := reg.UUID()
	return otpRpc("GenerateOTPWhatsapp", "generate/otp/whatsapp", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    phone_number,
		"device":     device,
		"kind":       "whatsapp",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"user":       user,
	})
}

/**
* SendOtpEmail
* @param project_id string, email, name, device string, length, duration int, user et.Json
* @return et.Json
**/
func SendOtpEmail(project_id string, email, name, device string, length int, duration int, user et.Json) (et.Item, error) {
	service_id := reg.UUID()
	return otpRpc("GenerateOTPEmail", "generate/otp/email", et.Json{
		"project_id": project_id,
		"service_id": service_id,
		"channel":    email,
		"device":     device,
		"kind":       "email",
		"name":       name,
		"length":     length,
		"duration":   duration,
		"user":       user,
	})
}

/**
* VerifyOtp
* @param channel, device, kind, code, clientId string
* @return et.Item, error
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
