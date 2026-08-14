package jquery_test

import (
	"strings"
	"testing"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jquery"
)

func TestJQuery_SelectAll(t *testing.T) {
	query := et.Json{
		"from": "users",
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_SelectColumns(t *testing.T) {
	query := et.Json{
		"from":   "users",
		"select": []string{"id", "name", "age"},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT "id", "name", "age" FROM "users"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_MissingFrom(t *testing.T) {
	_, err := jquery.JQuery(et.Json{})
	if err == nil {
		t.Fatal("expected error for missing 'from'")
	}
}

func TestJQuery_SimpleWhereEq(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"name": et.Json{"eq": "cesar"},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" WHERE "name" = 'cesar'`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_MultipleColumnsAreAnded(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"name": et.Json{"eq": "cesar"},
			"age":  et.Json{"eq": 30},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// sortedKeys orders "age" before "name" alphabetically.
	want := `SELECT * FROM "users" WHERE "age" = 30 AND "name" = 'cesar'`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_MultipleOperatorsSameColumn(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"age": et.Json{"more": 10, "less": 20},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" WHERE "age" < 20 AND "age" > 10`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_AndOrFromFeaturesExample(t *testing.T) {
	// Mirrors the shape from features.md, confirmed with the user:
	// "and"/"or" are arrays of condition objects.
	query := et.Json{
		"from": "table",
		"wheres": et.Json{
			"name": et.Json{"eq": "cesar"},
			"and": []et.Json{
				{"age": et.Json{"eq": 30}},
			},
			"or": []et.Json{
				{"age": et.Json{"more": 45}},
			},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "table" WHERE "age" = 30 AND "name" = 'cesar' AND "age" > 45`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_OrGroupWithMultipleItemsIsWrappedInParens(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"or": []et.Json{
				{"age": et.Json{"less": 10}},
				{"age": et.Json{"more": 45}},
			},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" WHERE ("age" < 10 OR "age" > 45)`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_NestedAndOr(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"or": []et.Json{
				{
					"and": []et.Json{
						{"age": et.Json{"more_eq": 18}},
						{"age": et.Json{"less_eq": 30}},
					},
				},
				{"name": et.Json{"eq": "admin"}},
			},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" WHERE (("age" >= 18 AND "age" <= 30) OR "name" = 'admin')`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_AllOperators(t *testing.T) {
	cases := []struct {
		name string
		ops  et.Json
		want string
	}{
		{"eq", et.Json{"eq": "x"}, `"c" = 'x'`},
		{"neg", et.Json{"neg": "x"}, `"c" != 'x'`},
		{"less", et.Json{"less": 5}, `"c" < 5`},
		{"less_eq", et.Json{"less_eq": 5}, `"c" <= 5`},
		{"more", et.Json{"more": 5}, `"c" > 5`},
		{"more_eq", et.Json{"more_eq": 5}, `"c" >= 5`},
		{"like", et.Json{"like": "%x%"}, `"c" ILIKE '%x%'`},
		{"in", et.Json{"in": []any{1, 2, 3}}, `"c" IN (1, 2, 3)`},
		{"not_in", et.Json{"not_in": []any{1, 2}}, `"c" NOT IN (1, 2)`},
		{"is", et.Json{"is": true}, `"c" IS true`},
		{"is_not", et.Json{"is_not": true}, `"c" IS NOT true`},
		{"null", et.Json{"null": true}, `"c" IS NULL`},
		{"not_null", et.Json{"not_null": true}, `"c" IS NOT NULL`},
		{"between", et.Json{"between": []any{1, 10}}, `"c" BETWEEN 1 AND 10`},
		{"not_between", et.Json{"not_between": []any{1, 10}}, `"c" NOT BETWEEN 1 AND 10`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			query := et.Json{
				"from": "t",
				"wheres": et.Json{
					"c": tc.ops,
				},
			}

			sql, err := jquery.JQuery(query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := `SELECT * FROM "t" WHERE ` + tc.want
			if sql != want {
				t.Fatalf("got %q, want %q", sql, want)
			}
		})
	}
}

func TestJQuery_InvalidOperator(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"name": et.Json{"bogus": "x"},
		},
	}

	_, err := jquery.JQuery(query)
	if err == nil {
		t.Fatal("expected error for invalid operator")
	}
}

func TestJQuery_BetweenRequiresTwoValues(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"age": et.Json{"between": []any{1, 2, 3}},
		},
	}

	_, err := jquery.JQuery(query)
	if err == nil {
		t.Fatal("expected error for between with wrong arity")
	}
}

func TestJQuery_InRequiresArray(t *testing.T) {
	query := et.Json{
		"from": "users",
		"wheres": et.Json{
			"age": et.Json{"in": 5},
		},
	}

	_, err := jquery.JQuery(query)
	if err == nil {
		t.Fatal("expected error for in with non-array value")
	}
}

func TestJQuery_OrderBy(t *testing.T) {
	query := et.Json{
		"from":          "users",
		"order_by":      []string{"name"},
		"order_by_desc": []string{"age"},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" ORDER BY "name" ASC, "age" DESC`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_LimitOffsetFromPage(t *testing.T) {
	query := et.Json{
		"from": "users",
		"limit": et.Json{
			"page": 3,
			"rows": 50,
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" LIMIT 50 OFFSET 100`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_LimitFirstPageHasNoOffset(t *testing.T) {
	query := et.Json{
		"from": "users",
		"limit": et.Json{
			"page": 1,
			"rows": 100,
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" LIMIT 100`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_FullQueryFromFeaturesExample(t *testing.T) {
	query := et.Json{
		"from":   "table",
		"select": []string{"id", "name", "age"},
		"wheres": et.Json{
			"name": et.Json{"eq": "cesar"},
			"and": []et.Json{
				{"age": et.Json{"eq": 30}},
			},
			"or": []et.Json{
				{"age": et.Json{"more": 45}},
			},
		},
		"limit": et.Json{
			"page": 1,
			"rows": 100,
		},
		"order_by":      []string{"name"},
		"order_by_desc": []string{"age"},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT "id", "name", "age" FROM "table" WHERE "age" = 30 AND "name" = 'cesar' AND "age" > 45 ORDER BY "name" ASC, "age" DESC LIMIT 100`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_QualifiedIdentifiers(t *testing.T) {
	query := et.Json{
		"from":   "public.users",
		"select": []string{"users.id"},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(sql, `"public"."users"`) {
		t.Fatalf("expected quoted qualified from, got %q", sql)
	}
	if !strings.Contains(sql, `"users"."id"`) {
		t.Fatalf("expected quoted qualified column, got %q", sql)
	}
}

func TestJQuery_AcceptsRawJSONUnmarshaledMaps(t *testing.T) {
	// Simulates a query arriving as real JSON text (json.Unmarshal
	// produces map[string]interface{}/[]interface{}, not et.Json/[]et.Json),
	// which is the far more common real-world path than hand-built et.Json.
	raw := map[string]any{
		"from": "users",
		"wheres": map[string]any{
			"and": []any{
				map[string]any{"age": map[string]any{"eq": float64(30)}},
			},
		},
	}

	sql, err := jquery.JQuery(et.Json(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" WHERE "age" = 30`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}
