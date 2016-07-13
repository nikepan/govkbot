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
func helpHandler(m *vkbot.Message) (reply string) {
  return "help received"
}

//vkbot.HandleMessage("/", anyHandler)
//vkbot.HandleMessage("/me", meHandler)
vkbot.HandleMessage("/help", helpHandler)

//vkbot.HandleAction("chat_invite_user", inviteHandler)
//vkbot.HandleAction("chat_kick_user", kickHandler)
//vkbot.HandleAction("friend_add", addFriendHandler)
//vkbot.HandleAction("friend_delete", deleteFriendHandler)

vkbot.Listen(config.VKToken, "", "")
```
