package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var ProfileFolders *linq.Model

func DefineProfileFolders(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
	}

	if ProfileFolders != nil {
		return nil
	}

	ProfileFolders = linq.NewModel(SchemaModule, "PROFILE_FOLDERS", "Tabla de carpetas por perfil", 1)
	ProfileFolders.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	ProfileFolders.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	ProfileFolders.DefineColum("profile_tp", "", "VARCHAR(80)", "-1")
	ProfileFolders.DefineColum("folder_id", "", "VARCHAR(80)", "-1")
	ProfileFolders.DefineColum("index", "", "INTEGER", 0)
	ProfileFolders.DefinePrimaryKey([]string{"module_id", "profile_tp", "folder_id"})
	ProfileFolders.DefineIndex([]string{
		"date_make",
		"index",
	})
	ProfileFolders.DefineForeignKey("module_id", Modules.Column("_id"))
	ProfileFolders.DefineForeignKey("folder_id", Folders.Column("_id"))

	if err := ProfileFolders.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* Profile Folder
**/
func GetProfileFolderByIdT(idT string) (et.Item, error) {
	return ProfileFolders.Data().
		Where(ProfileFolders.Column("_idt").Eq(idT)).
		First()
}

/**
* GetProfileFolderById
* @param moduleId string
* @param profileTp string
* @param folderId string
* @return et.Item, error
**/
func GetProfileFolderById(moduleId, profileTp, folderId string) (et.Item, error) {
	return ProfileFolders.Data().
		Where(ProfileFolders.Column("module_id").Eq(moduleId)).
		And(ProfileFolders.Column("profile_tp").Eq(profileTp)).
		And(ProfileFolders.Column("folder_id").Eq(folderId)).
		First()
}

/**
* CheckProfileFolder
* @param moduleId string
* @param profileTp string
* @param folderId string
* @param chk bool
* @return et.Item, error
**/
func CheckProfileFolder(moduleId, profileTp, folderId string, chk bool) (et.Item, error) {
	if !utility.ValidId(moduleId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(profileTp) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profile_tp")
	}

	if !utility.ValidId(folderId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "folder_id")
	}

	if !chk {
		result, err := ProfileFolders.Delete().
			Where(ProfileFolders.Column("module_id").Eq(moduleId)).
			And(ProfileFolders.Column("profile_tp").Eq(profileTp)).
			And(ProfileFolders.Column("folder_id").Eq(folderId)).
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

	current, err := GetProfileFolderById(moduleId, profileTp, folderId)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		data := et.Json{}
		data.Set("module_id", moduleId)
		data.Set("profile_tp", profileTp)
		data.Set("folder_id", folderId)
		result, err := ProfileFolders.Insert(data).
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
