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
	err := RouteMessage(&m)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestRouteAction(t *testing.T) {
	HandleAction("friend_add", baseHandler)
	SetAPI("", "test", "")
	m := Message{Body: ""}
	err := RouteAction(&m)
	if err != nil {
		t.Error(err.Error())
	}
}