/**
* Package dialect define el contrato Dialect que traduce
* columnas/limit/like al SQL de un motor especifico, mas un registry
* con patron factory para "cargar" el dialecto correcto por nombre en
* tiempo de ejecucion. Es un paquete independiente (no depende de
* jquery) para poder reutilizarse desde cualquier paquete que genere
* SQL, no solo el query builder.
*
* Cada motor soportado vive en su propio archivo dentro de este mismo
* paquete y se registra a si mismo desde su propio func init() llamando
* a Register. Motores incluidos: postgres.go (Postgres), sqlite.go
* (SQLite), mysql.go (MySQL), sqlserver.go (SQLServer) y oracle.go
* (Oracle). Agregar uno nuevo consiste en:
*
*  1. Crear dialect/<motor>.go con un tipo que implemente Dialect.
*  2. Registrarlo en su func init(): Register("<motor>", func() Dialect { ... }).
*
* Los paquetes consumidores (p.ej. jquery) nunca necesitan un switch
* nuevo: basta con que ese archivo exista en este paquete y con pedir
* el dialecto por nombre via Get.
**/
package dialect

import "fmt"

const (
	ERR_DIALECT_NOT_SUPPORTED = "dialecto no soportado (%s)"
)

/**
* Dialect: Reglas de sintaxis de un motor de base de datos.
**/
type Dialect interface {
	// Name identifica el dialecto (p.ej. "postgres").
	Name() string
	// QuoteIdent aplica el caracter de comillado de identificadores
	// del motor (p.ej. "col" en postgres, `col` en mysql, [col] en
	// sqlserver). Soporta identificadores calificados como "tabla.columna".
	QuoteIdent(name string) string
	// Like devuelve la palabra clave a usar para el operador LIKE
	// (postgres usa ILIKE para comparacion sin distinguir mayusculas).
	Like() string
	// LimitOffset arma la clausula de paginacion. rows es la cantidad
	// maxima de filas (LIMIT); offset es cuantas filas saltar.
	LimitOffset(rows, offset int) string
}

/**
* Factory: Construye una instancia nueva de un Dialect.
**/
type Factory func() Dialect

var registry = map[string]Factory{}

/**
* Register: Agrega o reemplaza la factory de un dialecto.
* @param name string, factory Factory
* @return void
**/
func Register(name string, factory Factory) {
	registry[name] = factory
}

/**
* Get: Devuelve una instancia nueva del dialecto registrado bajo name.
* @param name string
* @return Dialect, error
**/
func Get(name string) (Dialect, error) {
	factory, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf(ERR_DIALECT_NOT_SUPPORTED, name)
	}

	return factory(), nil
}

/**
* Registered: Nombres de los dialectos registrados.
* @return []string
**/
func Registered() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}

	return names
}
