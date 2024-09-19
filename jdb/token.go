package jdb

import (
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/claim"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
)

/**
* DevelopToken
**/
func DevelopToken() {
	production := envar.EnvarBool(false, "PRODUCTION")
	if !production {
		device := "DevelopToken"
		key := claim.TokenKey(device, device, device)
		token, err := cache.Get(key, "")
		if err != nil {
			return
		}

		if token == "" {
			token, err = claim.NewToken(device, device, device, "requests", "Default token", device, time.Hour*24*90)
			if console.AlertE(err) != nil {
				return
			}
		}

		claim.SetToken(device, device, device, token)
		console.LogK(device, token)
	}
}
