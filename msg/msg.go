package msg

const (
	MSG_ATRIB_REQUIRED      = "Atributo requerido (%s)"
	MSG_VALUE_REQUIRED      = "Atributo requerido (%s) value:%s"
	MSG_USER_INVALID        = "Usuario no valido, atrib (%s)"
	MSG_MAIL_001            = "Verifica tu dirección de correo electrónico"
	MSG_MAIL_002            = "Código de validación"
	MSG_MOBILE_VALIDATION   = "%s: Código de validación %s. Su código expira en 3 minutos."
	MSG_MOBILE_SEND_VERIFY  = "Enviamos un código a tu número celular para verificar tu identidad."
	MSG_CODE_INCORRECT      = "Código de verificacion incorrecto"
	MSG_MOBILE_SIGNIN       = "%s: Bienvenido %s, acabas de iniciar sesion."
	MSG_MOBILE_WELCOME      = "Bienvenido %s, acabas de ser registrado como usuario en %s, con el perfil de %s."
	MSG_ADMIN_WELCOME       = "Bienvenido %s, acabas de ser registrado como administrador de la base de datos de %s."
	MSG_SESION_BEGIN        = "Bienvanido %s, acabas de inciar sesión en %s"
	MSG_MAIN_CONNECT        = "Connect to main db"
	ERR_RECORS_STATE        = "Not edit for register state. (%s)"
	ERR_ENV_REQUIRED        = "Variables de entorno requerida (%s)"
	ERR_NOT_CACHE_SERVICE   = "Not cache service"
	ERR_COMM                = "Not connect db"
	NOT_SELECT_DRIVE        = "Not select drive"
	NOT_CONNECT_DB          = "Not connect db"
	NOT_INIT_CORE           = "Not init core schema"
	NOT_INIT_MIGRATION      = "Not init migration"
	NOT_MAIN_CONNECT        = "Not connect main db"
	NOT_MAIN_SERIE          = "Not exists serie in main db"
	NOT_TOKEN_NOT_EXISTS    = "Token not exists"
	RECORD_NOT_FOUND        = "Record not found"
	RECORD_FOUND            = "Record found"
	RECORD_NOT_ACTIVE       = "Record not active:%s"
	RECORD_NOT_CREATE       = "Record not create"
	RECORD_DELETE           = "Record delete"
	RECORD_NOT_DELETE       = "Record not delete"
	RECORD_NOT_UPDATE       = "Record not update"
	RECORD_NOT_CHANGE       = "Record not change"
	RECORD_CREATE           = "Create record"
	RECORD_UPDATE           = "Update record"
	RECORD_CANCELED         = "Record canceled"
	RECORD_ARCHIVED         = "Record archived"
	RECORD_ON_DELETE        = "Record on deleted"
	RECORD_IN_PROCESS       = "Record in process"
	RECORD_BATCH_LOAD       = "Record batch loaded"
	RECORD_DUPLICATE        = "Record duplicate"
	ERR_WHERE_NOT_DEFINED   = "Condition where not defined"
	ERR_DB_NOT_EXISTS       = "Conection to database not exist"
	ERR_DB_INDEX_NOT_EXISTS = "Database not exist index:(%d)"
	ERR_MIGRATION_ESCHEMA   = "Migration scehma faillure"
	ERR_MIGRATION_MODEL     = "Migration model faillure"
	ERR_NOT_NATS_SERVICE    = "Not nats service"
	MODEL_NOT_FOUND         = "Model not found:(%s)"
	TABLE_RECORD_FOUND      = "Record found, table name:%s"
	ERR_SQL                 = "SQL Error, %s\nQuery:%s"
	STEP_VERITY             = "Verify"
	MODULE_NOT_FOUND        = "Module not found"
	PROFILE_NOT_FOUND       = "Profile not found - %s module %s"
	PROJECT_NOT_FOUND       = "Project not found - %s"
	USER_FONUND             = "User found"
	USER_NOT_FONUND         = "User not found"
	SYSTEM_NOT_HAVE_ADMIN   = "You don't have admin user"
	SYSTEM_HAVE_ADMIN       = "You have an admin user"
	SYSTEM_USER_NOT_FOUNT   = "User not found"
	MARTER_NOT_FOUNT        = "Master DB not found:%s"
	TABLE_NOT_FOUND         = "Table not found:%s.%s"
	RPC_NOT_FOUND           = "RPC not found"
	ERR_NOT_FONUND          = "Document not found"
	ERR_NOT_COLLETION_MONGO = "Collection not found"
	ERR_INVALID_TYPE        = "Invalid type"
)
