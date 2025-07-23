package resilience

import (
	"slices"
	"time"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/service"
	"github.com/celsiainternet/elvis/utility"
)

type TpNotify int

const (
	TpNotifySms TpNotify = iota
	TpNotifyEmail
	TpNotifyWhatsapp
)

func (s TpNotify) String() string {
	return []string{"sms", "email", "whatsapp"}[s]
}

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
* SetContentSMS
* @param content string, params []et.Json
 */
func (s *Resilence) SetContentSMS(content string, params []et.Json) {
	s.Content = content
	s.Params = params
}

/**
* SetSubject
* @param subject string
 */
func (s *Resilence) SetContentEmail(subject string, htmlMessage string, params []et.Json) {
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
			service.TpTransactional,
			"resilience",
		)
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
