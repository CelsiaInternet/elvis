package test

import (
	"strings"
	"testing"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jquery"
	"github.com/celsiainternet/elvis/jquery/dialect"
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

func TestJQuery_GroupBy(t *testing.T) {
	query := et.Json{
		"from":     "users",
		"select":   []string{"name", "count(*)"},
		"group_by": []string{"name"},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT "name", COUNT(*) FROM "users" GROUP BY "name"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_Having(t *testing.T) {
	query := et.Json{
		"from":     "users",
		"select":   []string{"name", "count(*)"},
		"group_by": []string{"name"},
		"having": et.Json{
			"count(*)": et.Json{"more": 1},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT "name", COUNT(*) FROM "users" GROUP BY "name" HAVING COUNT(*) > 1`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_HavingWithAndOr(t *testing.T) {
	query := et.Json{
		"from":     "users",
		"select":   []string{"name", "count(*)"},
		"group_by": []string{"name"},
		"having": et.Json{
			"name": et.Json{"eq": "cesar"},
			"and": []et.Json{
				{"count(*)": et.Json{"eq": 30}},
			},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT "name", COUNT(*) FROM "users" GROUP BY "name" HAVING COUNT(*) = 30 AND "name" = 'cesar'`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_SelectAggregations(t *testing.T) {
	query := et.Json{
		"from":   "orders",
		"select": []string{"count(*)", "count()", "max(price)", "min(price)", "sum(price)", "COUNT(id)"},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT COUNT(*), COUNT(*), MAX("price"), MIN("price"), SUM("price"), COUNT("id") FROM "orders"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_FullQueryWithGroupByAndHaving(t *testing.T) {
	query := et.Json{
		"from":   "table",
		"select": []string{"id", "name", "age", "count(*)"},
		"wheres": et.Json{
			"name": et.Json{"eq": "cesar"},
			"and": []et.Json{
				{"age": et.Json{"eq": 30}},
			},
			"or": []et.Json{
				{"age": et.Json{"more": 45}},
			},
		},
		"group_by": []string{"name"},
		"having": et.Json{
			"name": et.Json{"eq": "cesar"},
			"and": []et.Json{
				{"count(*)": et.Json{"eq": 30}},
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

	want := `SELECT "id", "name", "age", COUNT(*) FROM "table" WHERE "age" = 30 AND "name" = 'cesar' AND "age" > 45 GROUP BY "name" HAVING COUNT(*) = 30 AND "name" = 'cesar' ORDER BY "name" ASC, "age" DESC LIMIT 100`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_FromWithAlias(t *testing.T) {
	query := et.Json{
		"from": "users:A",
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" AS "A"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_FromWithSchemaQualifiedAlias(t *testing.T) {
	query := et.Json{
		"from": "public.users:A",
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "public"."users" AS "A"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_ColumnReferenceValue(t *testing.T) {
	query := et.Json{
		"from": "users:A",
		"join": et.Json{
			"to": "roles:B",
			"on": et.Json{
				"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}},
			},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" AS "A" JOIN "roles" AS "B" ON "B"."user_id" = "A"."id"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_LeftJoinWithAndCondition(t *testing.T) {
	query := et.Json{
		"from": "users:A",
		"join": et.Json{
			"type": "left",
			"to":   "roles:B",
			"on": et.Json{
				"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}},
				"and": []et.Json{
					{"B.role": et.Json{"neg": "admin"}},
				},
			},
		},
		"select": []string{"A.id", "A.name"},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT "A"."id", "A"."name" FROM "users" AS "A" LEFT JOIN "roles" AS "B" ON "B"."user_id" = "A"."id" AND "B"."role" != 'admin'`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_MultipleJoins(t *testing.T) {
	query := et.Json{
		"from": "users:A",
		"join": []et.Json{
			{
				"type": "inner",
				"to":   "roles:B",
				"on":   et.Json{"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}}},
			},
			{
				"type": "right",
				"to":   "departments:C",
				"on":   et.Json{"C.id": et.Json{"eq": et.Json{"col": "A.department_id"}}},
			},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT * FROM "users" AS "A" INNER JOIN "roles" AS "B" ON "B"."user_id" = "A"."id" RIGHT JOIN "departments" AS "C" ON "C"."id" = "A"."department_id"`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_JoinMissingTo(t *testing.T) {
	query := et.Json{
		"from": "users:A",
		"join": et.Json{
			"on": et.Json{"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}}},
		},
	}

	if _, err := jquery.JQuery(query); err == nil {
		t.Fatal("expected error for join missing 'to'")
	}
}

func TestJQuery_JoinMissingOn(t *testing.T) {
	query := et.Json{
		"from": "users:A",
		"join": et.Json{
			"to": "roles:B",
		},
	}

	if _, err := jquery.JQuery(query); err == nil {
		t.Fatal("expected error for join missing 'on'")
	}
}

func TestJQuery_JoinInvalidType(t *testing.T) {
	query := et.Json{
		"from": "users:A",
		"join": et.Json{
			"type": "outer",
			"to":   "roles:B",
			"on":   et.Json{"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}}},
		},
	}

	if _, err := jquery.JQuery(query); err == nil {
		t.Fatal("expected error for invalid join type")
	}
}

func TestJQuery_FullQueryWithJoinFromFeaturesExample(t *testing.T) {
	query := et.Json{
		"from": "users:A",
		"join": et.Json{
			"to": "roles:B",
			"on": et.Json{
				"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}},
				"and": []et.Json{
					{"B.role": et.Json{"neg": "admin"}},
				},
			},
		},
		"select": []string{"A.id", "A.name", "A.age", "count(*)"},
		"wheres": et.Json{
			"A.name": et.Json{"eq": "cesar"},
			"and": []et.Json{
				{"A.age": et.Json{"eq": 30}},
			},
			"or": []et.Json{
				{"A.age": et.Json{"more": 45}},
			},
		},
		"limit": et.Json{
			"page": 1,
			"rows": 100,
		},
		"order_by":      []string{"A.name"},
		"order_by_desc": []string{"A.age"},
		"group_by":      []string{"A.name"},
		"having": et.Json{
			"A.name": et.Json{"eq": "cesar"},
			"and": []et.Json{
				{"count(*)": et.Json{"eq": 30}},
			},
		},
	}

	sql, err := jquery.JQuery(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `SELECT "A"."id", "A"."name", "A"."age", COUNT(*) FROM "users" AS "A" JOIN "roles" AS "B" ON "B"."user_id" = "A"."id" AND "B"."role" != 'admin' WHERE "A"."name" = 'cesar' AND "A"."age" = 30 AND "A"."age" > 45 GROUP BY "A"."name" HAVING "A"."name" = 'cesar' AND COUNT(*) = 30 ORDER BY "A"."name" ASC, "A"."age" DESC LIMIT 100`
	if sql != want {
		t.Fatalf("got %q, want %q", sql, want)
	}
}

func TestJQuery_FullFeatureSetAcrossDialects(t *testing.T) {
	// Exercises alias, join, aggregation, group by, having and
	// where/limit together on every registered dialect, to confirm
	// the JQueryBuilder additions (join, group_by, having, {"col":..})
	// go through Dialect.QuoteIdent/Like/LimitOffset consistently and
	// not just for the default (postgres) path.
	newQuery := func(dialectName string) et.Json {
		return et.Json{
			"dialect": dialectName,
			"from":    "users:A",
			"join": et.Json{
				"type": "left",
				"to":   "roles:B",
				"on": et.Json{
					"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}},
				},
			},
			"select":   []string{"A.id", "A.name", "count(*)"},
			"group_by": []string{"A.name"},
			"having": et.Json{
				"count(*)": et.Json{"more": 1},
			},
			"wheres": et.Json{
				"A.name": et.Json{"like": "%cesar%"},
			},
			"limit": et.Json{"page": 2, "rows": 10},
		}
	}

	cases := []struct {
		dialect string
		want    string
	}{
		{
			"postgres",
			`SELECT "A"."id", "A"."name", COUNT(*) FROM "users" AS "A" LEFT JOIN "roles" AS "B" ON "B"."user_id" = "A"."id" WHERE "A"."name" ILIKE '%cesar%' GROUP BY "A"."name" HAVING COUNT(*) > 1 LIMIT 10 OFFSET 10`,
		},
		{
			"sqlite",
			`SELECT "A"."id", "A"."name", COUNT(*) FROM "users" AS "A" LEFT JOIN "roles" AS "B" ON "B"."user_id" = "A"."id" WHERE "A"."name" LIKE '%cesar%' GROUP BY "A"."name" HAVING COUNT(*) > 1 LIMIT 10 OFFSET 10`,
		},
		{
			"mysql",
			"SELECT `A`.`id`, `A`.`name`, COUNT(*) FROM `users` AS `A` LEFT JOIN `roles` AS `B` ON `B`.`user_id` = `A`.`id` WHERE `A`.`name` LIKE '%cesar%' GROUP BY `A`.`name` HAVING COUNT(*) > 1 LIMIT 10 OFFSET 10",
		},
		{
			"sqlserver",
			`SELECT [A].[id], [A].[name], COUNT(*) FROM [users] AS [A] LEFT JOIN [roles] AS [B] ON [B].[user_id] = [A].[id] WHERE [A].[name] LIKE '%cesar%' GROUP BY [A].[name] HAVING COUNT(*) > 1 OFFSET 10 ROWS FETCH NEXT 10 ROWS ONLY`,
		},
		{
			"oracle",
			`SELECT "A"."id", "A"."name", COUNT(*) FROM "users" AS "A" LEFT JOIN "roles" AS "B" ON "B"."user_id" = "A"."id" WHERE "A"."name" LIKE '%cesar%' GROUP BY "A"."name" HAVING COUNT(*) > 1 OFFSET 10 ROWS FETCH NEXT 10 ROWS ONLY`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.dialect, func(t *testing.T) {
			sql, err := jquery.JQuery(newQuery(tc.dialect))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if sql != tc.want {
				t.Fatalf("got %q, want %q", sql, tc.want)
			}
		})
	}
}

func TestJQuery_AllJoinTypesAcrossAllDialects(t *testing.T) {
	// Cross product of every JoinType ("join"/"inner"/"left"/"right")
	// against every registered dialect, confirming both the join
	// keyword and the identifier quoting are correct in all 20
	// combinations. Expected identifiers are computed from each
	// dialect's own QuoteIdent, so this stays correct even if a
	// dialect's quoting rules change later.
	joinTypes := []struct {
		jsonType string
		keyword  string
	}{
		{"join", "JOIN"},
		{"inner", "INNER JOIN"},
		{"left", "LEFT JOIN"},
		{"right", "RIGHT JOIN"},
	}

	dialectNames := []string{"postgres", "sqlite", "mysql", "sqlserver", "oracle"}

	for _, dialectName := range dialectNames {
		d, err := dialect.Get(dialectName)
		if err != nil {
			t.Fatalf("unexpected error getting dialect %q: %v", dialectName, err)
		}

		for _, jt := range joinTypes {
			t.Run(dialectName+"/"+jt.jsonType, func(t *testing.T) {
				query := et.Json{
					"dialect": dialectName,
					"from":    "users:A",
					"join": et.Json{
						"type": jt.jsonType,
						"to":   "roles:B",
						"on": et.Json{
							"B.user_id": et.Json{"eq": et.Json{"col": "A.id"}},
						},
					},
				}

				sql, err := jquery.JQuery(query)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				want := "SELECT * FROM " + d.QuoteIdent("users") + " AS " + d.QuoteIdent("A") +
					" " + jt.keyword + " " + d.QuoteIdent("roles") + " AS " + d.QuoteIdent("B") +
					" ON " + d.QuoteIdent("B.user_id") + " = " + d.QuoteIdent("A.id")

				if sql != want {
					t.Fatalf("got %q, want %q", sql, want)
				}
			})
		}
	}
}
