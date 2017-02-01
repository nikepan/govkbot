# Go VK bot package [![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/nikepan/govkbot)


This is simple VK.com bot API.

It can:

* Reply to messages (private and chat)  
* Send greetings message when invited to chat
* Add and remove mutual friends
* Send notifies to admin

Installatioin:

`go get github.com/nikepan/govkbot`

For work you need get VK access token with rights: messages,friends,offline.


To get token you need:

1. You can use any app id from https://vk.com/apps?act=wingames, for example 4775211 
 (You create own app and get app_id on page https://vk.com/editapp?act=create (standalone app))
2. Open in browser with logged in VK (you must use IP, where you want run bot)
```
 https://oauth.vk.com/authorize?client_id={{app_id}}&scope=offline,group,messages,friends&display=page&response_type=token&redirect_uri=https://oauth.vk.com/blank.html
 ```
3. Copy token query parameter from URL string. Token valid only for IP from what you get it.


# Quickstart

```Go
package main
import "github.com/nikepan/govkbot"
import "log"

var VKAdminID = 3759927
var VKToken = "efjr98j9fj8jf4j958jj4985jfj9joijerf0fj548jf94jfiroefije495jf48"

func helpHandler(m *govkbot.Message) (reply string) {
  return "help received"
}

func errorHandler(m *govkbot.Message, err error) {
  log.Fatal(err.Error())
}

func main() {
    //govkbot.HandleMessage("/", anyHandler)
    //govkbot.HandleMessage("/me", meHandler)
    govkbot.HandleMessage("/help", helpHandler)

    //govkbot.HandleAction("chat_invite_user", inviteHandler)
    //govkbot.HandleAction("chat_kick_user", kickHandler)
    //govkbot.HandleAction("friend_add", addFriendHandler)
    //govkbot.HandleAction("friend_delete", deleteFriendHandler)

    govkbot.HandleError(errorHandler)

    govkbot.SetAutoFriend(true) // enable auto accept/delete friends

    govkbot.SetDebug(true) // log debug messages

    // Optional Direct VK API access
    govkbot.SetAPI(VKToken, "", "") // Need only before Listen, if you use direct API
    me, _ := govkbot.API.Me() // call API method
    log.Printf("current user: %+v\n", me.FullName())
    // Optional end

    govkbot.Listen(VKToken, "", "", VKAdminID)
}
```
