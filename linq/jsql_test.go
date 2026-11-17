package linq

import (
	"strings"
	"testing"

	"github.com/celsiainternet/elvis/et"
)

var testQueryModel *Model

/**
* queryTestModel: Modelo de prueba (sin base) con columnas, un atributo en _data y una columna oculta.
* @return *Model
**/
func queryTestModel() *Model {
	if testQueryModel != nil {
		return testQueryModel
	}

	schema := &Schema{Name: "TEST"}
	model := NewModel(schema, "ITEMS", "Prueba de RunQuery", 1)
	model.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	model.DefineColum("_state", "", "VARCHAR(80)", "0")
	model.DefineColum("_id", "", "VARCHAR(80)", "-1")
	model.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	model.DefineColum("owner_id", "", "VARCHAR(80)", "-1")
	model.DefineColum("name", "", "VARCHAR(250)", "")
	model.DefineColum("description", "", "VARCHAR(250)", "")
	model.DefineColum("secret", "", "VARCHAR(250)", "")
	model.DefineColum("_data", "", "JSONB", "{}")
	model.DefineAtrib("service", "", "VARCHAR(80)", "-1")
	model.DefineHidden([]string{"secret"})
	testQueryModel = model

	return model
}

/**
* querySql: SQL que arma RunQuery para query y las condiciones fijas wheres, sin ejecutarlo.
* @param t *testing.T, query et.Json, wheres ...*Where
* @return string
**/
func querySql(t *testing.T, query et.Json, wheres ...*Where) string {
	t.Helper()
	model := queryTestModel()
	ql, err := queryLinq(model, query, wheres...)
	if err != nil {
		t.Fatalf("queryLinq: %v", err)
	}

	return ql.Sql()
}

/**
* assertContains: Falla si sql no contiene part.
* @param t *testing.T, sql, part string
**/
func assertContains(t *testing.T, sql, part string) {
	t.Helper()
	if !strings.Contains(sql, part) {
		t.Fatalf("se esperaba %q en:\n%s", part, sql)
	}
}

/**
* assertNotContains: Falla si sql contiene part.
* @param t *testing.T, sql, part string
**/
func assertNotContains(t *testing.T, sql, part string) {
	t.Helper()
	if strings.Contains(sql, part) {
		t.Fatalf("no se esperaba %q en:\n%s", part, sql)
	}
}

func TestRunQueryWithoutWheres(t *testing.T) {
	sql := querySql(t, et.Json{})
	assertContains(t, sql, "FROM TEST.ITEMS AS A")
	assertNotContains(t, sql, "WHERE")
	assertNotContains(t, sql, "'secret'")
}

func TestRunQueryEqEscapesQuotes(t *testing.T) {
	query := et.Json{
		"wheres": []any{
			map[string]any{"name": map[string]any{"eq": "O'Brien"}},
		},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "WHERE A.NAME = 'O''Brien'")
}

func TestRunQueryEscapesQuotedGoString(t *testing.T) {
	// et.Unquote haria strconv.Unquote de este texto y dejaria pasar la comilla; RunQuery no lo usa.
	query := et.Json{
		"wheres": []any{
			map[string]any{"name": map[string]any{"eq": `"\x27 OR 1=1 --"`}},
		},
	}
	sql := querySql(t, query)
	assertContains(t, sql, `WHERE A.NAME = '"\x27 OR 1=1 --"'`)
}

func TestRunQueryIn(t *testing.T) {
	query := et.Json{
		"wheres": []any{
			map[string]any{"owner_id": map[string]any{"in": []any{"cc:1", "id_deal:2", map[string]any{}}}},
		},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "WHERE A.OWNER_ID IN ('cc:1','id_deal:2')")
}

func TestRunQueryIgnoresUnknownHiddenAndUnsupported(t *testing.T) {
	query := et.Json{
		"wheres": []any{
			map[string]any{"nope": map[string]any{"eq": "x"}},
			map[string]any{"secret": map[string]any{"eq": "x"}},
			map[string]any{"name": map[string]any{"isNull": true}},
			map[string]any{"name": map[string]any{"eq": map[string]any{"a": 1}}},
			map[string]any{"name": map[string]any{"in": "no es lista"}},
		},
	}
	sql := querySql(t, query)
	assertNotContains(t, sql, "WHERE")
}

func TestRunQueryOperators(t *testing.T) {
	query := et.Json{
		"wheres": []any{
			map[string]any{"name": map[string]any{"neg": "a"}},
			map[string]any{"and": map[string]any{"name": map[string]any{"like": "%b%"}}},
			map[string]any{"date_make": map[string]any{"more": "2026-01-01", "less": "2026-12-31"}},
			map[string]any{"date_make": map[string]any{"moreEq": "2026-02-01", "lessEq": "2026-11-30"}},
		},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "WHERE A.NAME != 'a'")
	assertContains(t, sql, "AND A.NAME ILIKE '%b%'")
	assertContains(t, sql, "AND A.DATE_MAKE < '2026-12-31'")
	assertContains(t, sql, "AND A.DATE_MAKE > '2026-01-01'")
	assertContains(t, sql, "AND A.DATE_MAKE <= '2026-11-30'")
	assertContains(t, sql, "AND A.DATE_MAKE >= '2026-02-01'")
}

func TestRunQueryNumbersAndAtrib(t *testing.T) {
	query := et.Json{
		"wheres": []any{
			map[string]any{"service": map[string]any{"eq": "OCI"}},
			map[string]any{"name": map[string]any{"eq": float64(5)}},
		},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "A._DATA#>>'{service}' = 'OCI'")
	assertContains(t, sql, "AND A.NAME = 5")
}

func TestRunQueryOrIsGroupedAfterFixedWheres(t *testing.T) {
	model := queryTestModel()
	projects := model.Column("project_id").In("-1", "p1")
	state := model.Column("_state").Eq("0")
	query := et.Json{
		"wheres": []any{
			map[string]any{"name": map[string]any{"eq": "a"}},
			map[string]any{"or": map[string]any{"name": map[string]any{"eq": "b"}}},
		},
	}
	sql := querySql(t, query, projects, state)
	assertContains(t, sql, "WHERE A.PROJECT_ID IN ('-1','p1')")
	assertContains(t, sql, "AND A._STATE = '0'")
	assertContains(t, sql, "AND (A.NAME = 'a' OR A.NAME = 'b')")
}

func TestRunQueryWithoutOrIsNotGrouped(t *testing.T) {
	model := queryTestModel()
	state := model.Column("_state").Eq("0")
	query := et.Json{
		"wheres": []any{
			map[string]any{"owner_id": map[string]any{"eq": "cc:1"}},
		},
	}
	sql := querySql(t, query, state)
	assertContains(t, sql, "WHERE A._STATE = '0'\nAND A.OWNER_ID = 'cc:1'")
	assertNotContains(t, sql, "(A.OWNER_ID")
}

func TestRunQueryConcat(t *testing.T) {
	query := et.Json{
		"wheres": []any{
			map[string]any{"concat": map[string]any{"fields": []any{"name", "nope", "description"}, "like": "%juan perez%"}},
		},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "WHERE CONCAT(A.NAME, ' ', A.DESCRIPTION)::TEXT ILIKE '%juan perez%'")
}

func TestRunQueryConcatWithoutFields(t *testing.T) {
	query := et.Json{
		"wheres": []any{
			map[string]any{"concat": map[string]any{"fields": []any{"nope", "secret"}, "like": "%x%"}},
		},
	}
	sql := querySql(t, query)
	assertNotContains(t, sql, "WHERE")
}

func TestRunQueryOrderBy(t *testing.T) {
	query := et.Json{
		"order_by": map[string]any{"asc": []any{"name"}, "desc": []any{"date_make"}},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "ORDER BY A.NAME ASC, A.DATE_MAKE DESC")
}

func TestRunQueryOrderByInvalid(t *testing.T) {
	model := queryTestModel()
	for _, name := range []string{"nope", "secret"} {
		query := et.Json{
			"order_by": map[string]any{"desc": []any{name}},
		}
		_, err := queryLinq(model, query)
		if err == nil {
			t.Fatalf("order_by %s: se esperaba error", name)
		}
	}
}

func TestRunQuerySelects(t *testing.T) {
	query := et.Json{
		"selects": []any{"_id", "name"},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "'_id', A._ID")
	assertContains(t, sql, "'name', A.NAME")
	assertNotContains(t, sql, "'description'")

	model := queryTestModel()
	for _, name := range []string{"nope", "secret"} {
		invalid := et.Json{
			"selects": []any{name},
		}
		_, err := queryLinq(model, invalid)
		if err == nil {
			t.Fatalf("selects %s: se esperaba error", name)
		}
	}
}

func TestRunQueryPageAndRows(t *testing.T) {
	cases := []struct {
		query et.Json
		page  int
		rows  int
	}{
		{et.Json{}, 1, QUERY_DEFAULT_ROWS},
		{et.Json{"page": 3, "rows": 20}, 3, 20},
		{et.Json{"page": -1, "limit": 50}, 1, 50},
		{et.Json{"rows": 5000}, 1, QUERY_MAX_ROWS},
	}

	for _, item := range cases {
		page := queryPage(item.query)
		rows := queryRows(item.query)
		if page != item.page || rows != item.rows {
			t.Fatalf("%v: page %d rows %d, se esperaba %d %d", item.query, page, rows, item.page, item.rows)
		}
	}
}

func TestRunQuerySearch(t *testing.T) {
	model := queryTestModel()
	projects := model.Column("project_id").In("-1", "p1")
	query := et.Json{
		"wheres": []any{
			map[string]any{"name": map[string]any{"eq": "a"}},
			map[string]any{"or": map[string]any{"name": map[string]any{"eq": "b"}}},
		},
		"search": map[string]any{"value": "  Ana  Gó'mez ", "fields": []any{"name", "service"}},
	}
	sql := querySql(t, query, projects)
	assertContains(t, sql, "WHERE A.PROJECT_ID IN ('-1','p1')")
	assertContains(t, sql, "AND (A.NAME = 'a' OR A.NAME = 'b')")
	assertContains(t, sql, "AND CONCAT(A.NAME, ' ', A._DATA#>>'{service}')::TEXT ILIKE '%Ana%'")
	assertContains(t, sql, "AND CONCAT(A.NAME, ' ', A._DATA#>>'{service}')::TEXT ILIKE '%Gó''mez%'")
}

func TestRunQuerySearchEscapesWildcards(t *testing.T) {
	query := et.Json{
		"search": map[string]any{"value": `50%_a\b`, "fields": []any{"name"}},
	}
	sql := querySql(t, query)
	assertContains(t, sql, `ILIKE '%50\%\_a\\b%'`)
}

func TestRunQuerySearchMaxWords(t *testing.T) {
	query := et.Json{
		"search": map[string]any{"value": "w1 w2 w3 w4 w5 w6", "fields": []any{"name"}},
	}
	sql := querySql(t, query)
	assertContains(t, sql, "'%w5%'")
	assertNotContains(t, sql, "w6")
}

func TestRunQuerySearchEmpty(t *testing.T) {
	cases := []et.Json{
		{"search": map[string]any{"value": "   ", "fields": []any{"name"}}},
		{"search": map[string]any{"value": ""}},
		{"search": "texto suelto"},
	}
	for _, query := range cases {
		sql := querySql(t, query)
		assertNotContains(t, sql, "WHERE")
	}
}

func TestRunQuerySearchInvalid(t *testing.T) {
	model := queryTestModel()
	long := strings.Repeat("a", SEARCH_MAX_LENGTH+1)
	cases := []et.Json{
		{"value": "x"},
		{"value": "x", "fields": []any{}},
		{"value": "x", "fields": []any{"name", "name", "name", "name", "name", "name"}},
		{"value": "x", "fields": []any{"nope"}},
		{"value": "x", "fields": []any{"secret"}},
		{"value": "x", "fields": []any{"_data"}},
		{"value": long, "fields": []any{"name"}},
	}
	for _, search := range cases {
		_, err := queryLinq(model, et.Json{"search": search})
		if err == nil {
			t.Fatalf("search %v: se esperaba error", search)
		}
	}
}
