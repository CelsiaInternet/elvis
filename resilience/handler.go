package resilience

import (
	"net/http"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
)

/**
* Load
* @return error
 */
func Load() error {
	if resilience != nil {
		return nil
	}

	_, err := cache.Load()
	if err != nil {
		return err
	}

	_, err = event.Load()
	if err != nil {
		return err
	}

	resilience = NewResilence()
	return nil
}

/**
* Add
* @param tag, description string, fn interface{}, fnArgs ...interface{}
* @return *Transaction
 */
func Add(tag, description string, fn interface{}, fnArgs ...interface{}) *Transaction {
	if resilience == nil {
		logs.Log("resilience", "resilience is nil")
		return nil
	}

	result := NewTransaction(tag, description, fn, fnArgs...)
	resilience.Transactions = append(resilience.Transactions, result)
	logs.Log("resilience", "add:", result.Json().ToString())
	resilience.Notify(result)
	resilience.Run(result)

	return result
}

/**
* SetNotifyType
* @param notifyType TpNotify
 */
func SetNotifyType(notifyType TpNotify) {
	resilience.NotifyType = notifyType
}

/**
* SetContactNumbers
* @param contactNumbers []string
 */
func SetContactNumbers(contactNumbers []string) {
	resilience.ContactNumbers = contactNumbers
}

/**
* SetEmails
* @param emails []et.Json
 */
func SetEmails(emails []et.Json) {
	resilience.Emails = emails
}

/**
* SetTemplateId
* @param templateId int
 */
func SetTemplateId(templateId int) {
	resilience.TemplateId = templateId
}

/**
* SetContentSMS
* @param content string, params []et.Json
 */
func SetContentSMS(content string, params []et.Json) {
	resilience.SetContentSMS(content, params)
}

/**
* SetContentEmail
* @param subject string, htmlMessage string, params []et.Json
 */
func SetContentEmail(subject string, htmlMessage string, params []et.Json) {
	resilience.SetContentEmail(subject, htmlMessage, params)
}

/**
* HttpGetResilience
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilience(w http.ResponseWriter, r *http.Request) {
	if resilience == nil {
		response.JSON(w, r, http.StatusServiceUnavailable, et.Json{
			"message": "resilience is not initialized",
		})
		return
	}

	data := resilience.Json()
	response.JSON(w, r, http.StatusOK, data)
}

/**
* HttpGetResilienceById
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilienceById(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	id := body.Str("id")
	transaction := resilience.GetById(id)
	if transaction == nil {
		response.JSON(w, r, http.StatusNotFound, et.Json{
			"message": "transaction not found",
		})
		return
	}

	response.JSON(w, r, http.StatusOK, transaction.Json())
}

/**
* HttpGetResilienceByTag
* @param w http.ResponseWriter, r *http.Request
**/
func HttpGetResilienceByTag(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	tag := body.Str("tag")
	transaction := resilience.GetByTag(tag)
	if transaction == nil {
		response.JSON(w, r, http.StatusNotFound, et.Json{
			"message": "transaction not found",
		})
		return
	}

	response.JSON(w, r, http.StatusOK, transaction.Json())
}
