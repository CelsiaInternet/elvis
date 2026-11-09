package event

const (
	ERR_ENV_REQUIRED     = "variables de entorno requerida (%s)"
	ERR_NOT_CONNECT      = "not connect nats"
	ERR_PARAM_NOT_FOUND  = "param not found"
	ERR_CLIENT_ID_EMPTY  = "client id is empty"
	ERR_CHANNEL_REQUIRED = "channel is required"
	PARAMS_UPDATED       = "Params updated"
	ERR_REPORT_REQUIRED  = "error report: %s is required"
	ERR_REPORT_LEVEL     = "error report: level must be error, warning or critical (%s)"
	ERR_REPORT_TICKET    = "error report: ticket_owner must be ADA or GLPI (%s)"
)
