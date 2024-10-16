package linq

import (
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/strs"
)

func ddlColumn(col *Column) string {
	var result string

	_default := et.NewAny(col.Default)

	if _default.Str() == "NOW()" {
		result = strs.Append(`DEFAULT NOW()`, result, " ")
	} else {
		result = strs.Append(strs.Format(`DEFAULT %v`, et.Unquote(col.Default)), result, " ")
	}

	if col.Type == "SERIAL" {
		result = strs.Uppcase(col.Type)
	} else if len(col.Type) > 0 {
		result = strs.Append(strs.Uppcase(col.Type), result, " ")
	}
	if len(col.name) > 0 {
		result = strs.Append(strs.Uppcase(col.name), result, " ")
	}

	return result
}

func ddlIndex(col *Column) string {
	result := jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_IDX ON $1($4);`, col.Model.Table, strs.Uppcase(col.Model.Schema.Name), col.Model.Name, strs.Uppcase(col.name))
	if col.Low() == strs.Lowcase(SourceField) && col.Model.indexeSource {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_IDX ON $1 USING GIN($4);`, col.Model.Table, strs.Uppcase(col.Model.Schema.Name), col.Model.Name, strs.Uppcase(col.name))
	}

	return result
}

func ddlUniqueIndex(col *Column) string {
	result := jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $2_$3_$4_IDX ON $1($4);`, col.Model.Table, strs.Uppcase(col.Model.Schema.Name), col.Model.Name, strs.Uppcase(col.name))
	if col.Low() == strs.Lowcase(SourceField) && col.Model.indexeSource {
		result = jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $2_$3_$4_IDX ON $1 USING GIN($4);`, col.Model.Table, strs.Uppcase(col.Model.Schema.Name), col.Model.Name, strs.Uppcase(col.name))
	}

	return result
}

func ddlPrimaryKey(model *Model) string {
	primaryKeys := func() []string {
		var result []string
		for _, v := range model.PrimaryKeys {
			result = append(result, strs.Uppcase(v))
		}

		return result
	}

	return strs.Format(`PRIMARY KEY (%s)`, strings.Join(primaryKeys(), ", "))
}

func ddlForeignKeys(model *Model) string {
	var result string
	for _, ref := range model.ForeignKey {
		key := strs.Replace(model.Table, ".", "_") + "_" + ref.Fkey
		key = strs.Replace(key, "-", "_") + "_FKEY"
		key = strs.Lowcase(key)
		return strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) %s;`, model.Table, strs.Uppcase(key), strs.Uppcase(ref.Fkey), ref.DDL())
	}

	return result
}

func ddlSetSync(model *Model) string {
	result := jdb.SQLDDL(`
	DROP TRIGGER IF EXISTS RECORDS_BEFORE_INSERT ON $1 CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_INSERT
	BEFORE INSERT ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_INSERT();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_UPDATE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_UPDATE
	BEFORE UPDATE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_UPDATE();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_DELETE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_DELETE
	BEFORE DELETE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_DELETE();
	
	DROP TRIGGER IF EXISTS RECORDS_AFTER_INSERT ON $1 CASCADE;
	CREATE TRIGGER RECORDS_AFTER_INSERT
	AFTER INSERT ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_AFTER_INSERT();

	DROP TRIGGER IF EXISTS RECORDS_AFTER_UPDATE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_AFTER_UPDATE
	AFTER UPDATE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_AFTER_UPDATE();

	DROP TRIGGER IF EXISTS RECORDS_AFTER_DELETE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_AFTER_DELETE
	AFTER DELETE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_AFTER_DELETE();`, model.Table)

	result = strs.Replace(result, "\t", "")

	return result
}

func ddlSetRecyclig(model *Model) string {
	result := jdb.SQLDDL(`
  DROP TRIGGER IF EXISTS RECYCLING ON $1 CASCADE;
	CREATE TRIGGER RECYCLING
	AFTER UPDATE ON $1
	FOR EACH ROW WHEN (OLD._STATE!=NEW._STATE)
	EXECUTE PROCEDURE core.RECYCLING_UPDATE();`, model.Table)

	result = strs.Replace(result, "\t", "")

	return result
}

func ddlTable(model *Model) string {
	NewColumn(model, IdTFiled, "UUId", "VARCHAR(80)", "-1")

	var result string
	var columns string
	var indexs string
	var uniqueKeys string

	appedColumns := func(def string) {
		columns = strs.Append(columns, def, ",\n")
	}

	appendIndex := func(def string) {
		indexs = strs.Append(indexs, def, "\n")
	}

	appendUniqueKey := func(def string) {
		uniqueKeys = strs.Append(uniqueKeys, def, ", ")
	}

	for _, column := range model.Definition {
		if column.Tp == TpColumn {
			def := ddlColumn(column)
			appedColumns(def)
			if column.Indexed {
				if column.Unique {
					def := column.DDLUniqueIndex()
					appendUniqueKey(def)
				} else {
					def := column.DDLIndex()
					appendIndex(def)
				}
			}
		}
	}
	columns = strs.Append(columns, ",", "")
	columns = strs.Append(columns, ddlPrimaryKey(model), "\n")
	table := strs.Format("\nCREATE TABLE IF NOT EXISTS %s (\n%s);", model.Table, columns)
	result = strs.Append(result, table, "\n")
	result = strs.Append(result, uniqueKeys, "\n")
	result = strs.Append(result, indexs, "\n\n")
	foreign := ddlForeignKeys(model)
	result = strs.Append(result, foreign, "\n\n")
	sync := ddlSetSync(model)
	result = strs.Append(result, sync, "\n\n")
	if model.UseState {
		recicle := ddlSetRecyclig(model)
		result = strs.Append(result, recicle, "\n\n")
	}

	model.Ddl = result

	return model.Ddl
}

func dllMigration(model *Model) string {
	var fields []string

	table := model.Table
	model.Table = strs.Append(model.Schema.Name, "NEW_TABLE", ",")
	ddl := model.DDL()

	for _, column := range model.Definition {
		fields = append(fields, column.name)
	}

	insert := strs.Format(`INSERT INTO %s(%s) SELECT %s FROM %s;`, model.Name, strings.Join(fields, ", "), strings.Join(fields, ", "), table)

	drop := strs.Format(`DROP TABLE %s CASCADE;`, model.Name)

	alter := strs.Format(`ALTER TABLE %s RENAME TO %s;`, model.Name, table)

	result := strs.Format(`%s %s %s %s`, ddl, insert, drop, alter)

	return result
}
