package daemon

import (
	"os"

	"github.com/celsiainternet/elvis/et"
)

type TypeCommand string

const (
	CMD_Help    TypeCommand = "help"
	CMD_Version TypeCommand = "version"
	CMD_Status  TypeCommand = "status"
	CMD_Start   TypeCommand = "start"
	CMD_Stop    TypeCommand = "stop"
	CMD_Restart TypeCommand = "restart"
	CMD_Conf    TypeCommand = "conf"
)

func ToTypeCommand(val string) TypeCommand {
	switch val {
	case "version", "--v":
		return CMD_Version
	case "status", "--s":
		return CMD_Status
	case "start", "--start":
		return CMD_Start
	case "stop", "--stop":
		return CMD_Stop
	case "restart", "--restart":
		return CMD_Restart
	case "conf", "--conf":
		return CMD_Conf
	}

	return CMD_Help
}

type RepositoryCMD interface {
	Help(key string)
	Version() string
	Status() et.Json
	Start() et.Item
	Stop() et.Item
	Restart() et.Item
}

var apps map[string]RepositoryCMD
var app string = "systemd"

func main() {
	if len(os.Args) < 2 {
		help("")
		return
	}

	_app := apps[app]
	command := os.Args[1]
	switch ToTypeCommand(command) {
	case CMD_Version:
		_app.Version()
	case CMD_Help:
		if len(os.Args) > 2 {
			_app.Help(os.Args[2])
		} else {
			_app.Help("")
		}
	case CMD_Status:
		_app.Status()
	case CMD_Start:
		_app.Start()
	case CMD_Stop:
		_app.Stop()
	case CMD_Restart:
		_app.Restart()
	}
}

func help(val string) {
	if val == "" {
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

func Registry(name string, cmd RepositoryCMD) {
	if apps == nil {
		apps = make(map[string]RepositoryCMD)
	}

	apps[name] = cmd
	app = name
}
