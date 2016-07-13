package main

import "github.com/nikepan/govkbot/vkbot"

var config configuration

func main() {
	config = configuration{ServiceURL: "http://vk1.mysocio.ru/api/"}

	readJSON("config.json", &config)

	vkbot.HandleMessage("/", anyHandler)
	vkbot.HandleMessage("/me", meHandler)
	vkbot.HandleMessage("/help", helpHandler)

	vkbot.HandleAction("chat_invite_user", inviteHandler)
	vkbot.HandleAction("chat_kick_user", kickHandler)
	vkbot.HandleAction("friend_add", addFriendHandler)
	vkbot.HandleAction("friend_delete", deleteFriendHandler)

	vkbot.Listen(config.VKToken, "", "")
}
