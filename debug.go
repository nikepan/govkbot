package govkbot

import "log"

func isDebugging() bool {
	return API.DEBUG
}

func debugPrint(format string, values ...interface{}) {
	if isDebugging() {
		log.Printf("[VKBOT-DEBUG] "+format, values...)
	}
}
