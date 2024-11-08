package daemon

import "github.com/celsiainternet/elvis/et"

type Systemd struct {
}

func New() RepositoryCMD {
	return &Systemd{}
}

func (s *Systemd) Version() string {
	result := "Version: 1.0.0"
	println(result)

	return result
}

func (s *Systemd) Help(key string) {
	if key == "" {
		println("Uso: daemon [opciones]")
		println("Opciones:")
		println("  --h, --help     Mostrar esta ayuda")
		println("  --v, --version  Mostrar la versión")
		println("  --s, --status   Mostrar el estado del servicio")
		println("  --r, --restart  Reiniciar el servicio")
		println("  --up   				 Actualizar el servicio")
		println("  --down   			 Detener el servicio")
		println("  --start   			 Iniciar el servicio")
	}
}

func (s *Systemd) Status() et.Json {
	return et.Json{}
}

func (s *Systemd) Start() et.Item {
	return et.Item{}
}

func (s *Systemd) Stop() et.Item {
	return et.Item{}
}

func (s *Systemd) Restart() et.Item {
	return et.Item{}
}

func init() {
	Registry("systemd", New())
}
