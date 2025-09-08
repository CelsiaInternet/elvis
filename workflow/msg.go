package workflow

const (
	MSG_FLOW_CREATED              = "Flujo creado Tag:%s version:%s name:%s"
	MSG_START_WORKFLOW            = "Start the workflow"
	MSG_INSTANCE_FAILED           = "Instance fallido:%s Tag:%s status:%s, error:%s"
	MSG_INSTANCE_STATUS           = "Instance:%s Tag:%s status:%s"
	MSG_INSTANCE_GOTO             = "Instance %s Tag:%s ir al step:%d, message:%s"
	MSG_INSTANCE_ALREADY_DONE     = "flow already done"
	MSG_INSTANCE_ALREADY_RUNNING  = "flow already running"
	MSG_INSTANCE_WORKFLOWS_IS_NIL = "workFlows is nil"
	MSG_ID_REQUIRED               = "id is required"
	MSG_INSTANCE_EXPRESSION_TRUE  = "Resultado de la expresion es true"
	MSG_INSTANCE_EXPRESSION_FALSE = "Resultado de la expresion es false"
	MSG_INSTANCE_ROLLBACK         = "Esta intentando hacer rollback de un step que no existe"
	MSG_INSTANCE_ROLLBACK_STEP    = "haciendo rollback del step:%d"
	MSG_INSTANCE_STEP_CREATED     = "Step creado:%d name:%s Tag:%s"
	MSG_INSTANCE_ROLLBACK_CREATED = "Rollback creado:%d name:%s Tag:%s"
	MSG_INSTANCE_CONSISTENCY      = "Consistencia definida Tag:%s consistency:%s"
	MSG_INSTANCE_RESILIENCE       = "Resilencia definida Tag:%s totalAttempts:%d timeAttempts:%s retentionTime:%s"
	MSG_INSTANCE_IFELSE           = "IfElse definido step:%d name:%s expresion:%s ? %d : %d Tag:%s"
)
