# Go VK bot package

This is simple VK.com bot API.

It can:

* Reply to messages (private and chat)  
* Send greetings message when invited to chat
* Add and remove mutual friends
* Send notifies to admin

Installatioin:

`go get github.com/nikepan/govkbot`

For work you need get VK access token with rights: messages,friends,offline.

You can get it by this url in browser:

https://oauth.vk.com/authorize?client_id={{app_id}}&scope=offline,group,messages,friends&display=page&response_type=token&redirect_uri=https://oauth.vk.com/blank.html

app_id you can get on page https://vk.com/editapp?act=create (standalone app)

Usage example:

```Go
func helpHandler(m *govkbot.Message) (reply string) {
  return "help received"
}

//govkbot.HandleMessage("/", anyHandler)
//govkbot.HandleMessage("/me", meHandler)
govkbot.HandleMessage("/help", helpHandler)

//govkbot.HandleAction("chat_invite_user", inviteHandler)
//govkbot.HandleAction("chat_kick_user", kickHandler)
//govkbot.HandleAction("friend_add", addFriendHandler)
//govkbot.HandleAction("friend_delete", deleteFriendHandler)

govkbot.Listen(config.VKToken, "", "")
```
