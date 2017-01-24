package govkbot

import "log"

func IsDebugging() bool {
	return API.DEBUG
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		log.Printf("[VKBOT-DEBUG] "+format, values...)
	}
}
