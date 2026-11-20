package service

import (
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
)

const (
	// Paquete RPC por defecto de catalog, que guarda el registro central de errores (Soporte)
	DEFAULT_CATALOG_RPC = "catalog.Services"
	// Metodo RPC de catalog que registra un caso de soporte
	SUPPORT_RPC_METHOD = "RegisterError"
	// event.Work de respaldo cuando catalog no resuelve por RPC (catalog lo atiende con event.Stack)
	EVENT_SUPPORT_ERROR = "support/error"
)

/**
* RegisterError: Registra un caso de soporte en catalog por RPC; si no resuelve, publica event.Work EVENT_SUPPORT_ERROR.
* @param report event.ErrorReport
* @return et.Item, error
**/
func RegisterError(report event.ErrorReport) (et.Item, error) {
	if err := report.Validate(); err != nil {
		return et.Item{}, err
	}

	pkg := envar.GetStrOs(DEFAULT_CATALOG_RPC, "CATALOG_RPC")
	return dispatchRpc(pkg, SUPPORT_RPC_METHOD, EVENT_SUPPORT_ERROR, report.ToJson())
}
