package event

import (
	"strings"
	"testing"
	"time"

	"github.com/celsiainternet/elvis/et"
)

func TestParseErrorReport(t *testing.T) {
	base := func() et.Json {
		return et.Json{"app": " app.vcore ", "message": "falla", "owner_kind": "lead", "owner_id": "L-1"}
	}

	report, err := ParseErrorReport(base())
	if err != nil {
		t.Fatal(err)
	}
	if report.App != "app.vcore" || report.Level != ERROR_LEVEL_ERROR || report.OccurredAt.IsZero() || report.Detail == nil {
		t.Errorf("normalizado: %+v", report)
	}

	data := base()
	data["ticket_id"] = "45871"
	data["ticket_owner"] = "glpi"
	data["occurred_at"] = "2026-11-05T10:20:00"
	data["level"] = "CRITICAL"
	report, err = ParseErrorReport(data)
	if err != nil {
		t.Fatal(err)
	}
	if report.TicketOwner != TICKET_OWNER_GLPI || report.Level != ERROR_LEVEL_CRITICAL || report.OccurredAt.Hour() != 10 {
		t.Errorf("ticket: %+v", report)
	}

	cases := []struct {
		name string
		set  et.Json
		want string
	}{
		{"sin app", et.Json{"app": ""}, "app"},
		{"sin owner_id", et.Json{"owner_id": " "}, "owner_id"},
		{"level invalido", et.Json{"level": "info"}, "level"},
		{"ticket sin owner", et.Json{"ticket_id": "1"}, "ticket_owner"},
		{"ticket owner invalido", et.Json{"ticket_owner": "JIRA"}, "ADA or GLPI"},
	}
	for _, c := range cases {
		data := base()
		for k, v := range c.set {
			data[k] = v
		}
		if _, err := ParseErrorReport(data); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err=%v, want %q", c.name, err, c.want)
		}
	}

	// Ticket owner sin ticket_id es valido: el ticket puede no existir aun
	data = base()
	data["ticket_owner"] = "ADA"
	if _, err := ParseErrorReport(data); err != nil {
		t.Errorf("ticket_owner sin ticket_id: %v", err)
	}
}

func TestErrorWithoutConnection(t *testing.T) {
	if err := Error(ErrorReport{App: "a", Message: "m", OwnerKind: "k", OwnerId: "1", OccurredAt: time.Now()}); err != nil {
		t.Errorf("reporte valido sin NATS: %v", err)
	}
	if err := Error(ErrorReport{App: "a"}); err == nil {
		t.Error("reporte invalido aceptado")
	}
}
