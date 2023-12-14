package bus

import (
	"log"
	"os"
	"runtime"

	"github.com/cgalvisleon/elvis/console"
)

func init() {
	color := true
	if runtime.GOOS == "windows" {
		color = false
	}
	DefaultLogger = RequestLogger(&DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: !color})

	if err := DefineApimanager(); err != nil {
		console.PanicE(err)
	}

	console.LogK("BUS", "Init bus")
}
