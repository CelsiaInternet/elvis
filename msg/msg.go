package msg

const (
	MSG_ATRIB_REQUIRED      = "Atributo requerido (%s)"
	MSG_VALUE_REQUIRED      = "Atributo requerido (%s) value:%s"
	MSG_USER_INVALID        = "Usuario no valido, atrib (%s)"
	MSG_USER_NOT_FOUND      = "Usuario no encontrado"
	MSG_MAIL_001            = "Verifica tu dirección de correo electrónico"
	MSG_MAIL_002            = "Código de validación"
	MSG_MOBILE_VALIDATION   = "%s: Código de validación %s. Su código expira en 3 minutos."
	MSG_MOBILE_SEND_VERIFY  = "Enviamos un código a tu número celular para verificar tu identidad."
	MSG_CODE_INCORRECT      = "Código de verificacion incorrecto"
	MSG_MOBILE_SIGNIN       = "%s: Bienvenido %s, acabas de iniciar sesion."
	MSG_MOBILE_WELCOME      = "Bienvenido %s, acabas de ser registrado como usuario en %s, con el perfil de %s."
	MSG_ADMIN_WELCOME       = "Bienvenido %s, acabas de ser registrado como administrador de la base de datos de %s."
	MSG_SESION_BEGIN        = "Bienvanido %s, acabas de inciar sesión en %s"
	MSG_MAIN_CONNECT        = "Conectado a la base de datos principal"
	ERR_RECORS_STATE        = "no se puede editar el registro por el estado. (%s)"
	ERR_ENV_REQUIRED        = "variables de entorno requerida (%s)"
	ERR_NOT_CACHE_SERVICE   = "no hay servicio de caching"
	NOT_SELECT_DRIVE        = "Driver no seleccionado"
	NOT_CONNECT_DB          = "No connectado a la db"
	NOT_INIT_CORE           = "Schema no iniciado"
	NOT_INIT_MIGRATION      = "Migración no iniciada"
	NOT_MAIN_CONNECT        = "No connectado a la bs principal"
	NOT_MAIN_SERIE          = "Serie no existe en la base principal"
	NOT_TOKEN_NOT_EXISTS    = "Token no existe"
	RECORD_NOT_FOUND        = "Registro no encontrado"
	RECORD_FOUND            = "Registro encontrado"
	RECORD_NOT_ACTIVE       = "Registro no activo:%s"
	RECORD_NOT_CREATE       = "Registro no creado"
	RECORD_DELETE           = "Registro eliminado"
	RECORD_NOT_DELETE       = "Registro no eliminado"
	RECORD_NOT_UPDATE       = "Registro no actualizado"
	RECORD_NOT_CHANGE       = "Registro no cambiado"
	RECORD_CREATE           = "Registro creado"
	RECORD_UPDATE           = "Registro actualizado"
	RECORD_CANCELED         = "Registro cancelado"
	RECORD_ARCHIVED         = "Registro archivado"
	RECORD_ON_DELETE        = "Registro en eliminación"
	RECORD_IN_PROCESS       = "Registro en proceso"
	RECORD_BATCH_LOAD       = "Carga de registros"
	RECORD_DUPLICATE        = "Registro duplicado"
	RECORD_IS_SYSTEM        = "Registro del sistema"
	ERR_WHERE_NOT_DEFINED   = "Where no definido"
	ERR_DB_NOT_EXISTS       = "Database no existe"
	ERR_DB_INDEX_NOT_EXISTS = "Index not exist:(%s)"
	ERR_MIGRATION_ESCHEMA   = "Falló la migración del esquema"
	ERR_MIGRATION_MODEL     = "Falló la migración del modelo"
	ERR_NOT_NATS_SERVICE    = "No hay servicio de nats"
	MODEL_NOT_FOUND         = "Modelo no encontrado:(%s)"
	TABLE_RECORD_FOUND      = "Registro encontrado en la tabla:(%s)"
	ERR_SQL                 = "Error SQL:(%s)\nQuery:%s"
	STEP_VERITY             = "Paso de verificación"
	MODULE_NOT_FOUND        = "Módulo no encontrado:(%s)"
	PROFILE_NOT_FOUND       = "Perfil no encontrado:(%s)"
	PROJECT_NOT_FOUND       = "Proyecto no encontrado:(%s)"
	USER_FONUND             = "Usuario encontrado"
	USER_NOT_FONUND         = "Usuario no encontrado"
	SYSTEM_NOT_HAVE_ADMIN   = "No tienes un usuario admin"
	SYSTEM_HAVE_ADMIN       = "Ya tienes un usuario admin"
	MARTER_NOT_FOUNT        = "Master data no encontrado"
	TABLE_NOT_FOUND         = "Tabla no encontrada:(%s.%s)"
	RPC_NOT_FOUND           = "RPC no encontrado"
	ERR_NOT_FONUND          = "Documento no encontrado"
	ERR_NOT_COLLETION_MONGO = "Collection no encontrada"
	ERR_INVALID_TYPE        = "Tipo invalido"
	PASSWORD_NOT_MATCH      = "Contraseña no coincide"
)
