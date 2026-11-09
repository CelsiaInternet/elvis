package event

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/timezone"
)

const (
	ERROR_LEVEL_ERROR    = "error"
	ERROR_LEVEL_WARNING  = "warning"
	ERROR_LEVEL_CRITICAL = "critical"
	TICKET_OWNER_ADA     = "ADA"
	TICKET_OWNER_GLPI    = "GLPI"
)

var (
	errorLevels  = []string{ERROR_LEVEL_ERROR, ERROR_LEVEL_WARNING, ERROR_LEVEL_CRITICAL}
	ticketOwners = []string{TICKET_OWNER_ADA, TICKET_OWNER_GLPI}
)

/**
* ErrorReport: Error que un servicio reporta al registro central (catalog) por EVENT_ERROR.
* App, Message, OwnerKind y OwnerId son obligatorios. OwnerKind/OwnerId identifican la entidad
* afectada (lead, contrato, cliente...) y UserId al usuario que lo tuvo, si aplica. El ticket es
* opcional; TicketOwner es el sistema donde vive (ADA o GLPI) y define como lo abre el front.
**/
type ErrorReport struct {
	App         string    `json:"app"`
	Service     string    `json:"service"`
	Level       string    `json:"level"`
	Code        string    `json:"code"`
	Message     string    `json:"message"`
	Detail      et.Json   `json:"detail"`
	OwnerKind   string    `json:"owner_kind"`
	OwnerId     string    `json:"owner_id"`
	UserId      string    `json:"user_id"`
	TicketKind  string    `json:"ticket_kind"`
	TicketId    string    `json:"ticket_id"`
	TicketOwner string    `json:"ticket_owner"`
	OccurredAt  time.Time `json:"occurred_at"`
}

/**
* normalize: Quita espacios, Level en minusculas (error por defecto), TicketOwner en mayusculas
* y OccurredAt = ahora si no viene.
**/
func (s *ErrorReport) normalize() {
	for _, field := range []*string{&s.App, &s.Service, &s.Code, &s.Message, &s.OwnerKind, &s.OwnerId, &s.UserId, &s.TicketKind, &s.TicketId} {
		*field = strings.TrimSpace(*field)
	}

	s.Level = strings.ToLower(strings.TrimSpace(s.Level))
	if s.Level == "" {
		s.Level = ERROR_LEVEL_ERROR
	}

	s.TicketOwner = strings.ToUpper(strings.TrimSpace(s.TicketOwner))
	if s.Detail == nil {
		s.Detail = et.Json{}
	}

	if s.OccurredAt.IsZero() {
		s.OccurredAt = timezone.NowTime()
	}
}

/**
* Validate: Normaliza y valida el reporte.
* @return error
**/
func (s *ErrorReport) Validate() error {
	s.normalize()
	required := []struct{ name, val string }{
		{"app", s.App},
		{"message", s.Message},
		{"owner_kind", s.OwnerKind},
		{"owner_id", s.OwnerId},
	}
	for _, r := range required {
		if r.val == "" {
			return fmt.Errorf(ERR_REPORT_REQUIRED, r.name)
		}
	}

	if !slices.Contains(errorLevels, s.Level) {
		return fmt.Errorf(ERR_REPORT_LEVEL, s.Level)
	}

	if s.TicketId != "" && s.TicketOwner == "" {
		return fmt.Errorf(ERR_REPORT_REQUIRED, "ticket_owner")
	}

	if s.TicketOwner != "" && !slices.Contains(ticketOwners, s.TicketOwner) {
		return fmt.Errorf(ERR_REPORT_TICKET, s.TicketOwner)
	}

	return nil
}

/**
* ToJson
* @return et.Json
**/
func (s ErrorReport) ToJson() et.Json {
	return et.Json{
		"app":          s.App,
		"service":      s.Service,
		"level":        s.Level,
		"code":         s.Code,
		"message":      s.Message,
		"detail":       s.Detail,
		"owner_kind":   s.OwnerKind,
		"owner_id":     s.OwnerId,
		"user_id":      s.UserId,
		"ticket_kind":  s.TicketKind,
		"ticket_id":    s.TicketId,
		"ticket_owner": s.TicketOwner,
		"occurred_at":  s.OccurredAt,
	}
}

/**
* ParseErrorReport: Lee y valida el reporte que llega por EVENT_ERROR. occurred_at puede venir
* como RFC 3339 o sin zona ("2006-01-02T15:04:05", hora local de timezone).
* @param data et.Json
* @return ErrorReport, error
**/
func ParseErrorReport(data et.Json) (ErrorReport, error) {
	result := ErrorReport{
		App:         data.Str("app"),
		Service:     data.Str("service"),
		Level:       data.Str("level"),
		Code:        data.Str("code"),
		Message:     data.Str("message"),
		Detail:      data.Json("detail"),
		OwnerKind:   data.Str("owner_kind"),
		OwnerId:     data.Str("owner_id"),
		UserId:      data.Str("user_id"),
		TicketKind:  data.Str("ticket_kind"),
		TicketId:    data.Str("ticket_id"),
		TicketOwner: data.Str("ticket_owner"),
	}
	switch v := data["occurred_at"].(type) {
	case time.Time:
		result.OccurredAt = v
	case string:
		for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05"} {
			t, err := time.ParseInLocation(layout, v, timezone.Location())
			if err == nil {
				result.OccurredAt = t
				break
			}
		}
	}

	if err := result.Validate(); err != nil {
		return ErrorReport{}, err
	}

	return result, nil
}

/**
* Error: Publica el reporte en EVENT_ERROR sin bloquear a quien lo llama. Devuelve el error de
* validacion; si el reporte es valido la publicacion va en otra goroutine y no falla al llamador.
* Sin conexion a NATS no hace nada (como Publish).
* @param report ErrorReport
* @return error
**/
func Error(report ErrorReport) error {
	if err := report.Validate(); err != nil {
		return err
	}

	go Publish(EVENT_ERROR, report.ToJson())
	return nil
}
