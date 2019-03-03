# Go VK bot package
 
[![Build Status](https://travis-ci.org/nikepan/govkbot.svg?branch=master)](https://travis-ci.org/nikepan/govkbot)
[![codecov](https://codecov.io/gh/nikepan/govkbot/branch/master/graph/badge.svg)](https://codecov.io/gh/nikepan/govkbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikepan/govkbot)](https://goreportcard.com/report/github.com/nikepan/govkbot)
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/nikepan/govkbot)


This is simple VK.com bot API.


> At 2019-03-01 VK was restrict messages.send for user_tokens. This bot can work with group_token, and access to chat members if has admin rights in chat. You can use v1.0.1 also, if you need only user_token access.


It can:

* Reply to messages (private and chat)  
* Send greetings message when invited to chat
* Add and remove mutual friends
* Send notifies to admin

Installatioin:

`go get github.com/nikepan/govkbot`

For work you need get VK access token with rights: messages,friends,offline (see below).


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

func startHandler(m *govkbot.Message) (reply govkbot.Reply) {
	keyboard := govkbot.Keyboard{Buttons: make([][]govkbot.Button, 0)}
	button := govkbot.NewButton("/help", nil)
	row := make([]govkbot.Button, 0)
	row = append(row, button)
	keyboard.Buttons = append(keyboard.Buttons, row)

	return govkbot.Reply{Msg: availableCommands, Keyboard: &keyboard}
}

func errorHandler(m *govkbot.Message, err error) {
  log.Fatal(err.Error())
}

func main() {
    //govkbot.HandleMessage("/", anyHandler)
    //govkbot.HandleMessage("/me", meHandler)
    govkbot.HandleMessage("/help", helpHandler)
    govkbot.HandleAdvancedMessage("/start", startHandler)

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


# Getting token

You need standalone vk app_id. You can use any app_id from https://vk.com/apps?act=wingames, for example 4775211 
 (Or you can create own app and get app_id on page https://vk.com/editapp?act=create (standalone app))

You can get token from you server ip with this node.js package:
https://www.npmjs.com/package/vk-auth (you need login, pass and app_id)


To manual get token you need:

1. Open in browser with logged in VK (you must use IP, where you want run bot)
```
 https://oauth.vk.com/authorize?client_id={{app_id}}&scope=offline,groups,messages,friends&display=page&response_type=token&redirect_uri=https://oauth.vk.com/blank.html
 ```
2. Copy token query parameter from URL string. Token valid only for IP from what you get it.


If you receive validation check (for example, you use ip first time)
```json
{"error":{"error_code":17,"error_msg":"Validation required: please open redirect_uri in browser ...", 
"redirect_uri":"https://m.vk.com/login?act=security_check&api_hash=Qwerty1234567890"}}
```
you can use https://github.com/Yashko/vk-validation-node.
