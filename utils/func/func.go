package _func

import (
	"runtime/debug"

	"github.com/mapprotocol/fe-backend/resource/log"
)

func Go(fn func(), name string) {
	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to " + name)

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover" + name)
			}
		}()

		fn()
	}()
}
