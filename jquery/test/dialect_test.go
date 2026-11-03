package test

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
	_, err := dialect.Get("db2")
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

func TestRegistered_IncludesBuiltinDialects(t *testing.T) {
	names := dialect.Registered()

	for _, want := range []string{dialect.Postgres, dialect.SQLite, dialect.MySQL, dialect.SQLServer, dialect.Oracle} {
		if !slices.Contains(names, want) {
			t.Fatalf("expected %q in %v", want, names)
		}
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

func TestSQLiteDialect_QuoteIdent(t *testing.T) {
	d := &dialect.SQLiteDialect{}

	cases := map[string]string{
		"users":        `"users"`,
		"public.users": `"public"."users"`,
		"*":            "*",
	}

	for in, want := range cases {
		if got := d.QuoteIdent(in); got != want {
			t.Errorf("QuoteIdent(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSQLiteDialect_LikeAndLimitOffset(t *testing.T) {
	d := &dialect.SQLiteDialect{}

	if d.Like() != "LIKE" {
		t.Fatalf("got %q, want %q", d.Like(), "LIKE")
	}
	if got := d.LimitOffset(0, 0); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
	if got := d.LimitOffset(10, 20); got != "LIMIT 10 OFFSET 20" {
		t.Fatalf("got %q, want %q", got, "LIMIT 10 OFFSET 20")
	}
}

func TestMySQLDialect_QuoteIdent(t *testing.T) {
	d := &dialect.MySQLDialect{}

	cases := map[string]string{
		"users":        "`users`",
		"public.users": "`public`.`users`",
		"*":            "*",
	}

	for in, want := range cases {
		if got := d.QuoteIdent(in); got != want {
			t.Errorf("QuoteIdent(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMySQLDialect_LikeAndLimitOffset(t *testing.T) {
	d := &dialect.MySQLDialect{}

	if d.Like() != "LIKE" {
		t.Fatalf("got %q, want %q", d.Like(), "LIKE")
	}
	if got := d.LimitOffset(10, 0); got != "LIMIT 10" {
		t.Fatalf("got %q, want %q", got, "LIMIT 10")
	}
}

func TestSQLServerDialect_QuoteIdent(t *testing.T) {
	d := &dialect.SQLServerDialect{}

	cases := map[string]string{
		"users":        "[users]",
		"public.users": "[public].[users]",
		"*":            "*",
	}

	for in, want := range cases {
		if got := d.QuoteIdent(in); got != want {
			t.Errorf("QuoteIdent(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSQLServerDialect_LikeAndLimitOffset(t *testing.T) {
	d := &dialect.SQLServerDialect{}

	if d.Like() != "LIKE" {
		t.Fatalf("got %q, want %q", d.Like(), "LIKE")
	}
	if got := d.LimitOffset(0, 0); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
	if got := d.LimitOffset(10, 20); got != "OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY" {
		t.Fatalf("got %q, want %q", got, "OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY")
	}
}

func TestOracleDialect_QuoteIdent(t *testing.T) {
	d := &dialect.OracleDialect{}

	cases := map[string]string{
		"users":        `"users"`,
		"public.users": `"public"."users"`,
		"*":            "*",
	}

	for in, want := range cases {
		if got := d.QuoteIdent(in); got != want {
			t.Errorf("QuoteIdent(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestOracleDialect_LikeAndLimitOffset(t *testing.T) {
	d := &dialect.OracleDialect{}

	if d.Like() != "LIKE" {
		t.Fatalf("got %q, want %q", d.Like(), "LIKE")
	}
	if got := d.LimitOffset(0, 0); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
	if got := d.LimitOffset(10, 20); got != "OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY" {
		t.Fatalf("got %q, want %q", got, "OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY")
	}
}
