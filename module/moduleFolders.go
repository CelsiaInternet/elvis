package module

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

var ModelFolders *linq.Model

func DefineModuleFolders(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return logs.Panice(err)
	}

	if ModelFolders != nil {
		return nil
	}

	ModelFolders = linq.NewModel(SchemaModule, "MODULE_FOLDERS", "Tabla de folders por modulo", 1)
	ModelFolders.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	ModelFolders.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	ModelFolders.DefineColum("folder_id", "", "VARCHAR(80)", "-1")
	ModelFolders.DefineColum("index", "", "INTEGER", 0)
	ModelFolders.DefinePrimaryKey([]string{"module_id", "folder_id"})
	ModelFolders.DefineForeignKey("module_id", Modules.Col("_id"))
	ModelFolders.DefineForeignKey("folder_id", Folders.Col("_id"))
	Modules.DefineIndex([]string{
		"date_make",
		"index",
	})

	if err := ModelFolders.Init(); err != nil {
		return logs.Panice(err)
	}

	return nil
}

/**
* GetModuleFolderByIdT
* @param _idt string
* @return et.Item, error
**/
func GetModuleFolderByIdT(_idt string) (et.Item, error) {
	return ModelFolders.Data().
		Where(ModelFolders.Column("_idt").Eq(_idt)).
		First()
}

/**
* GetModuleFolderById
* @param module_id string
* @param folder_id string
* @return et.Item, error
**/
func GetModuleFolderById(module_id, folder_id string) (et.Item, error) {
	return ModelFolders.Data().
		Where(ModelFolders.Column("module_id").Eq(module_id)).
		And(ModelFolders.Column("folder_id").Eq(folder_id)).
		First()
}

// Check folder that module
func CheckModuleFolder(module_id, folder_id string, chk bool) (et.Item, error) {
	if !utility.ValidId(module_id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(folder_id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "folder_id")
	}

	if !chk {
		result, err := ModelFolders.Delete().
			Where(ModelFolders.Column("module_id").Eq(module_id)).
			And(ModelFolders.Column("folder_id").Eq(folder_id)).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return et.Item{
			Ok: result.Ok,
			Result: et.Json{
				"message": utility.OkOrNot(result.Ok, msg.RECORD_DELETE, msg.RECORD_NOT_DELETE),
			},
		}, nil
	}

	current, err := GetModuleFolderById(module_id, folder_id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		data := et.Json{}
		data.Set("module_id", module_id)
		data.Set("folder_id", folder_id)

		result, err := ModelFolders.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return et.Item{
			Ok: result.Ok,
			Result: et.Json{
				"message": utility.OkOrNot(result.Ok, msg.RECORD_CREATE, msg.RECORD_NOT_CREATE),
			},
		}, nil
	}

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": msg.RECORD_FOUND,
		},
	}, nil
}
