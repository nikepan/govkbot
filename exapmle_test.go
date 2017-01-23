package govkbot_test

import (
	"github.com/nikepan/govkbot"
	"log"
	"fmt"
)


func helpHandler(m *govkbot.Message) (reply string) {
	return "Available commands: /help, /me\nYou message" + m.Body
}

func errorHandler(msg *govkbot.Message, err error) {
	log.Fatal(err.Error())
}

func addFriendHandler(m *govkbot.Message) (reply string) {
	log.Printf("friend %+v added\n", m.UserID)
	govkbot.NotifyAdmin(fmt.Sprintf("user vk.com/id%+v add me to friends", m.UserID))
	return reply
}


func ExampleListen() {

	//govkbot.HandleMessage("/", anyHandler) // any commands starts with "/"
	//govkbot.HandleMessage("/me", meHandler)
	govkbot.HandleMessage("/help", helpHandler)  // any commands starts with "/help"

	//govkbot.HandleAction("chat_invite_user", inviteHandler)
	//govkbot.HandleAction("chat_kick_user", kickHandler)
	govkbot.HandleAction("friend_add", addFriendHandler)
	//govkbot.HandleAction("friend_delete", deleteFriendHandler)

	govkbot.HandleError(errorHandler)

	govkbot.SetAutoFriend(true) // enable auto accept/delete friends

	govkbot.SetDebug(true) // log debug messages

	// Optional Direct VK API access
	govkbot.SetAPI("!!!!VK_TOKEN!!!!", "", "") // Need only before Listen, if you use direct API
	me := govkbot.API.Me() // call API method
	log.Printf("current user: %+v\n", me.FullName())
	// Optional end

	govkbot.Listen("!!!!VK_TOKEN!!!!", "", "", 12345678) // start bot

}
