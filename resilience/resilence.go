package resilience

import (
	"net/http"
	"slices"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/service"
	"github.com/celsiainternet/elvis/utility"
)

type TpNotify string

const (
	TpNotifySms      TpNotify = "sms"
	TpNotifyEmail    TpNotify = "email"
	TpNotifyWhatsapp TpNotify = "whatsapp"
)

type Resilence struct {
	CreatedAt      time.Time
	Id             string
	Transactions   []*Transaction
	Attempts       int
	TimeAttempts   time.Duration
	NotifyType     TpNotify
	ContactNumbers []string
	Emails         []et.Json
	TemplateId     int
	Content        string
	Subject        string
	HtmlMessage    string
	Params         []et.Json
}

func (s *Resilence) Json() et.Json {
	transactions := make([]et.Json, len(s.Transactions))
	for i, transaction := range s.Transactions {
		transactions[i] = transaction.Json()
	}

	return et.Json{
		"id":              s.Id,
		"created_at":      s.CreatedAt,
		"transactions":    transactions,
		"attempts":        s.Attempts,
		"time_attempts":   s.TimeAttempts,
		"notify_type":     s.NotifyType,
		"contact_numbers": s.ContactNumbers,
		"emails":          s.Emails,
		"template_id":     s.TemplateId,
		"content":         s.Content,
		"subject":         s.Subject,
		"html_message":    s.HtmlMessage,
		"params":          s.Params,
	}
}

var resilience *Resilence

/**
* NewResilence
* @return *Resilience
 */
func NewResilence() *Resilence {
	attempts := envar.EnvarInt(3, "RESILIENCE_ATTEMPTS")
	timeAttempts := envar.EnvarNumber(30, "RESILIENCE_TIME_ATTEMPTS")

	return &Resilence{
		CreatedAt:      time.Now(),
		Id:             utility.UUID(),
		Transactions:   make([]*Transaction, 0),
		Attempts:       attempts,
		TimeAttempts:   time.Duration(timeAttempts) * time.Second,
		NotifyType:     TpNotifySms,
		ContactNumbers: make([]string, 0),
		Emails:         make([]et.Json, 0),
		TemplateId:     0,
		Content:        "",
		Subject:        "",
		HtmlMessage:    "",
		Params:         make([]et.Json, 0),
	}
}

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

func (s *Resilence) SetNotifyType(notifyType TpNotify) {
	s.NotifyType = notifyType
}

/**
* SetContactNumbers
* @param contactNumbers []string
 */
func (s *Resilence) SetContactNumbers(contactNumbers []string) {
	s.ContactNumbers = contactNumbers
}

/**
* SetEmails
* @param emails []et.Json
 */
func (s *Resilence) SetEmails(emails []et.Json) {
	s.Emails = emails
}

/**
* SetTemplateId
* @param templateId int
 */
func (s *Resilence) SetTemplateId(templateId int) {
	s.TemplateId = templateId
}

/**
* SetContentMessage
* @param content string, params []et.Json
 */
func (s *Resilence) SetContentMessage(content string, params []et.Json) {
	s.Content = content
	s.Params = params
}

/**
* SetSubject
* @param subject string
 */
func (s *Resilence) SetContentHtml(subject string, htmlMessage string, params []et.Json) {
	s.Subject = subject
	s.HtmlMessage = htmlMessage
	s.Params = params
}

/**
* Notify
* @param transaction *Transaction
 */
func (s *Resilence) Notify(transaction *Transaction) {
	projectId := envar.EnvarStr("-1", "PROJECT_ID")
	serviceId := utility.UUID()

	if s.NotifyType == TpNotifySms {
		service.SendSms(
			projectId,
			serviceId,
			s.ContactNumbers,
			s.Content,
			s.Params,
			service.TpTransactional, "resilience")
		return
	}

	if s.NotifyType == TpNotifyWhatsapp {
		service.SendWhatsapp(
			projectId,
			serviceId,
			s.TemplateId,
			s.ContactNumbers,
			s.Params,
			service.TpTransactional,
			"resilience",
		)
		return
	}

	service.SendEmail(
		projectId,
		serviceId,
		s.Emails,
		s.Subject,
		s.HtmlMessage,
		s.Params,
		service.TpTransactional,
		"resilience",
	)
}

/**
* Done
* @param transaction *Transaction
 */
func (s *Resilence) Done(transaction *Transaction) {
	idx := slices.IndexFunc(s.Transactions, func(t *Transaction) bool { return t.Id == transaction.Id })
	if idx != -1 {
		s.Transactions = append(s.Transactions[:idx], s.Transactions[idx+1:]...)
	}

	logs.Log("resilience", "done:", transaction.Json().ToString())
}

/**
* Run
* @param transaction *Transaction
 */
func (s *Resilence) Run(transaction *Transaction) {
	if s.TimeAttempts == 0 {
		return
	}

	time.AfterFunc(s.TimeAttempts, func() {
		if transaction.Status != StatusSuccess && transaction.Attempts < s.Attempts {
			_, err := transaction.Run()
			if err == nil {
				s.Done(transaction)
			} else {
				if transaction.Attempts == s.Attempts {
					s.Notify(transaction)
				} else {
					s.Run(transaction)
				}
			}
		}
	})
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
* GetById
* @param id string
* @return *Transaction
 */
func (s *Resilence) GetById(id string) *Transaction {
	idx := slices.IndexFunc(s.Transactions, func(t *Transaction) bool { return t.Id == id })
	if idx != -1 {
		return s.Transactions[idx]
	}

	return nil
}

/**
* GetByTag
* @param tag string
* @return *Transaction
 */
func (s *Resilence) GetByTag(tag string) *Transaction {
	idx := slices.IndexFunc(s.Transactions, func(t *Transaction) bool { return t.Tag == tag })
	if idx != -1 {
		return s.Transactions[idx]
	}

	return nil
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
