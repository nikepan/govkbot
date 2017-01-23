package govkbot

import (
//	"encoding/json"
//	"net/url"
	"testing"
	"net/url"
	"encoding/json"
)

func TestCall(t *testing.T) {
	SetAPI("","test","")
	api := API
	buf, err := api.Call("utils.getServerTime", url.Values{})
	if err != nil {
		t.Error("no response from VK")
	}

	m := SimpleResponse{}
	json.Unmarshal(buf, &m)
	if m.Error != nil {
		t.Error("Error returned")
	}
}

func TestVkAPI_Me(t *testing.T) {
	SetAPI("","test","")
	me := API.Me()
	if me.FullName() != "First Last" {
		t.Error(me.FullName())
	}
}

func TestVkAPI_GetChatInfo(t *testing.T) {
	SetAPI("","test","")
	chat, err := API.GetChatInfo(1)
	if err != nil {
		t.Error("Error returned")
	}
	if chat == nil {
		t.Error("Chat info == nil")
	}
}

func TestVkAPI_GetChatUsers(t *testing.T) {
	SetAPI("","test","")
	users, err := API.GetChatUsers(1)
	if err != nil {
		t.Error("Error returned")
	}
	if users == nil {
		t.Error("users == nil")
	}
}

func TestVkAPI_GetMessages(t *testing.T) {
	SetAPI("","test","")
	messages, err := API.GetMessages(100, 0)
	if err != nil {
		t.Error("Error returned")
	}
	if messages == nil {
		t.Error("messages == nil")
	}
}

func TestVkAPI_GetFriendRequests(t *testing.T) {
	SetAPI("","test","")
	_, err := API.GetFriendRequests(false)
	if err != nil {
		t.Error("Error returned")
	}
}

func TestVkAPI_AddFriend(t *testing.T) {
	SetAPI("","test","")
	ok := API.AddFriend(1)
	if !ok {
		t.Error("Wrong returned value")
	}
}

func TestVkAPI_DeleteFriend(t *testing.T) {
	SetAPI("","test","")
	ok := API.AddFriend(1)
	if !ok {
		t.Error("Wrong returned value")
	}
}

func TestMessage_MarkAsRead(t *testing.T) {
	SetAPI("","test","")
	m := Message{}
	err := m.MarkAsRead()
	if err != nil {
		t.Error("Error returned")
	}
}

func TestMessage_Reply(t *testing.T) {
	SetAPI("","test","")
	m := Message{}
	err := m.Reply("ok")
	if err != nil {
		t.Error("Error returned")
	}
}

func TestVkAPI_SendChatMessage(t *testing.T) {
	SetAPI("","test","")
	err := API.SendChatMessage(1, "ok")
	if err != nil {
		t.Error("Error returned")
	}
}

func TestVkAPI_SendMessage(t *testing.T) {
	SetAPI("","test","")
	err := API.SendMessage(1, "ok")
	if err != nil {
		t.Error("Error returned")
	}
}

func TestVkAPI_User(t *testing.T) {
	SetAPI("","test","")
	u, err := API.User(1)
	if err != nil {
		t.Error("Error returned")
	}
	if u.ID == 0 {
		t.Error("Wrong value returned")
	}
}

func TestVkAPI_NotifyAdmin(t *testing.T) {
	SetAPI("","test","")
	err := API.NotifyAdmin("ok")
	if err != nil {
		t.Error("Error returned")
	}
}