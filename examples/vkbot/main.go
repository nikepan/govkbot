package main

import "github.com/nikepan/govkbot"

var config configuration

func main() {
	config = configuration{ServiceURL: "http://vk1.mysocio.ru/api/"}

	readJSON("config.json", &config)

	govkbot.HandleMessage("/", anyHandler)
	govkbot.HandleMessage("/me", meHandler)
	govkbot.HandleMessage("/help", helpHandler)

	govkbot.HandleAction("chat_invite_user", inviteHandler)
	govkbot.HandleAction("chat_kick_user", kickHandler)
	govkbot.HandleAction("friend_add", addFriendHandler)
	govkbot.HandleAction("friend_delete", deleteFriendHandler)

	govkbot.Listen(config.VKToken, "", "")
}
