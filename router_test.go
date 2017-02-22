package govkbot

import (
	"testing"
)

func baseHandler(m *Message) (reply string) {
	return m.Body
}

func errorHandler(msg *Message, err error) {
	return
}

func TestHandleMessage(t *testing.T) {
	HandleMessage("/help", baseHandler)
	if len(bot.msgRoutes) == 0 {
		t.Error("Error adding message handler")
	}
}

func TestHandleAction(t *testing.T) {
	HandleAction("/help", baseHandler)
	if len(bot.actionRoutes) == 0 {
		t.Error("Error adding action handler")
	}
}

func TestHandleError(t *testing.T) {
	HandleError(errorHandler)
	if bot.errorHandler == nil {
		t.Error("Error set error handler")
	}
}

func TestSetAPI(t *testing.T) {
	token := "12345"
	SetAPI(token, "https://vk.com/api/", "5.52")
	if API.Token != token {
		t.Error("Error setup API")
	}
	SetLang("ru")
	if API.Lang != "ru" {
		t.Error("Error setup Lang")
	}
	SetLang("")
}

func TestSetToken(t *testing.T) {
	token := "12345"
	SetToken(token)
	if API.Token != token {
		t.Error("Error setup API")
	}
}

func TestSetAutoFriend(t *testing.T) {
	SetAutoFriend(true)
	if !bot.autoFriend {
		t.Error("Error set auto friend")
	}
	SetAutoFriend(false)
}

func TestSetDebug(t *testing.T) {
	SetDebug(true)
	if !API.DEBUG {
		t.Error("Error set debug mode")
	}
	SetDebug(false)
}

func TestRouteMessage(t *testing.T) {
	HandleMessage("/help", baseHandler)
	SetAPI("", "test", "")
	m := Message{Body: "/help"}
	replies, err := RouteMessage(&m)
	if err != nil {
		t.Error(err.Error())
	}
	if replies[0] != "/help" {
		t.Error(wrongValueReturned)
	}
}

func TestRouteAction(t *testing.T) {
	HandleAction("friend_add", baseHandler)
	SetAPI("", "test", "")
	m := Message{Action: "friend_add", Body: "ok"}
	replies, err := RouteAction(&m)
	if err != nil {
		t.Error(err.Error())
	}
	if replies[0] != "ok" {
		t.Error(wrongValueReturned)
	}
}

func TestRouteMessages(t *testing.T) {
	HandleMessage("/help", baseHandler)
	SetAPI("", "test", "")
	var messages []*Message
	m1 := Message{Body: "/help"}
	messages = append(messages, &m1)
	m2 := Message{Action: "friend_add", Body: "ok"}
	messages = append(messages, &m2)
	m3 := Message{Body: "skip"}
	messages = append(messages, &m3)
	replies := RouteMessages(messages)

	if replies[&m1][0] != "/help" {
		t.Error(wrongValueReturned)
	}
	if replies[&m2][0] != "ok" {
		t.Error(wrongValueReturned)
	}
}

func TestGetMessages(t *testing.T) {
	SetAPI("", "test", "")
	messages, err := GetMessages()
	if err != nil {
		t.Error(err.Error())
	}
	if len(messages) == 0 {
		t.Error("No messages")
	}
}

func TestCheckFriends(t *testing.T) {
	SetAPI("", "test", "")
	CheckFriends()
}

func TestMainRoute(t *testing.T) {
	SetAPI("", vkAPIURL, "")
	HandleError(errorHandler)
	MainRoute()
}

func TestNotifyAdmin(t *testing.T) {
	SetAPI("", "test", "")
	err := NotifyAdmin("ok")
	if err != nil {
		t.Error(err.Error())
	}
}
