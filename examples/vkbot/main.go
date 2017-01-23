package main

import (
	"github.com/nikepan/govkbot"
	"log"
)

var config configuration

func main() {
	config = configuration{ServiceURL: "http://vk1.mysocio.ru/api/"}

	readJSON("config.json", &config)

	govkbot.HandleMessage("/", anyHandler) // any commands starts with "/"
	govkbot.HandleMessage("/me", meHandler)
	govkbot.HandleMessage("/help", helpHandler)

	govkbot.HandleAction("chat_invite_user", inviteHandler)
	govkbot.HandleAction("chat_kick_user", kickHandler)
	govkbot.HandleAction("friend_add", addFriendHandler)
	govkbot.HandleAction("friend_delete", deleteFriendHandler)

	govkbot.HandleError(errorHandler)

	govkbot.SetAutoFriend(true) // enable auto accept/delete friends

	govkbot.SetDebug(true) // log debug messages

	// Optional Direct VK API access
	govkbot.SetAPI(config.VKToken, "", "") // Need only before Listen, if you use direct API
	me := govkbot.API.Me() // call API method
	log.Printf("current user: %+v\n", me.FullName())
	// Optional end

	govkbot.Listen(config.VKToken, "", "", config.AdminID) // start bot
}
