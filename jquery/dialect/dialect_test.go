package dialect_test

import (
	"slices"
	"testing"

	"github.com/celsiainternet/elvis/jquery/dialect"
)

func TestGet_Postgres(t *testing.T) {
	d, err := dialect.Get(dialect.Postgres)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.Name() != dialect.Postgres {
		t.Fatalf("got %q, want %q", d.Name(), dialect.Postgres)
	}
}

func TestGet_Unregistered(t *testing.T) {
	_, err := dialect.Get("sqlserver")
	if err == nil {
		t.Fatal("expected error for unregistered dialect")
	}
}

func TestRegister_CustomDialect(t *testing.T) {
	const name = "fake-for-test"

	dialect.Register(name, func() dialect.Dialect {
		return &dialect.PostgresDialect{}
	})

	d, err := dialect.Get(name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name() != dialect.Postgres {
		t.Fatalf("got %q, want %q", d.Name(), dialect.Postgres)
	}
}

func TestRegistered_IncludesPostgres(t *testing.T) {
	names := dialect.Registered()

	if !slices.Contains(names, dialect.Postgres) {
		t.Fatalf("expected %q in %v", dialect.Postgres, names)
	}
}

func TestPostgresDialect_QuoteIdent(t *testing.T) {
	d := &dialect.PostgresDialect{}

	cases := map[string]string{
		"users":        `"users"`,
		"public.users": `"public"."users"`,
		"*":            "*",
		`"already"`:    `"already"`,
	}

	for in, want := range cases {
		got := d.QuoteIdent(in)
		if got != want {
			t.Errorf("QuoteIdent(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPostgresDialect_LimitOffset(t *testing.T) {
	d := &dialect.PostgresDialect{}

	if got := d.LimitOffset(0, 0); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
	if got := d.LimitOffset(10, 0); got != "LIMIT 10" {
		t.Fatalf("got %q, want %q", got, "LIMIT 10")
	}
	if got := d.LimitOffset(10, 20); got != "LIMIT 10 OFFSET 20" {
		t.Fatalf("got %q, want %q", got, "LIMIT 10 OFFSET 20")
	}
}

func TestPostgresDialect_Like(t *testing.T) {
	d := &dialect.PostgresDialect{}
	if d.Like() != "ILIKE" {
		t.Fatalf("got %q, want %q", d.Like(), "ILIKE")
	}
}
