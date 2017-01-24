package govkbot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"testing"
)

const (
	WRONG_VALUE_RETURNED = "Wrong value returned"
)

func TestCall(t *testing.T) {
	api := API
	buf, err := api.Call("utils.getServerTime", url.Values{})
	if err != nil {
		t.Error("no response from VK")
	}

	m := SimpleResponse{}
	json.Unmarshal(buf, &m)
	if m.Error != nil {
		t.Error(m.Error.Error())
	}
}

func TestVkAPI_Call(t *testing.T) {
	api := API
	SetDebug(true)
	log.SetOutput(ioutil.Discard)
	buf, err := api.Call("messages.get", url.Values{})
	if err == nil {
		t.Error("no error returned: " + string(buf))
	}
}

func TestVkAPI_Me(t *testing.T) {
	SetAPI("", "test", "")
	me := API.Me()
	if me.FullName() != "First Last" {
		t.Error(me.FullName())
	}
}

func TestVkAPI_GetChatInfo(t *testing.T) {
	SetAPI("", "test", "")
	chat, err := API.GetChatInfo(1)
	if err != nil {
		t.Error(err.Error())
	}
	if chat == nil {
		t.Error("Chat info == nil")
	}
}

func TestVkAPI_GetChatUsers(t *testing.T) {
	SetAPI("", "test", "")
	users, err := API.GetChatUsers(1)
	if err != nil {
		t.Error(err.Error())
	}
	if users == nil {
		t.Error("users == nil")
	}
}

func TestVkAPI_GetMessages(t *testing.T) {
	SetAPI("", "test", "")
	messages, err := API.GetMessages(100, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if messages == nil {
		t.Error("messages == nil")
	}
}

func TestVkAPI_GetFriendRequests(t *testing.T) {
	SetAPI("", "test", "")
	_, err := API.GetFriendRequests(false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestVkAPI_AddFriend(t *testing.T) {
	SetAPI("", "test", "")
	ok := API.AddFriend(1)
	if !ok {
		t.Error(WRONG_VALUE_RETURNED)
	}
}

func TestVkAPI_DeleteFriend(t *testing.T) {
	SetAPI("", "test", "")
	ok := API.DeleteFriend(1)
	if !ok {
		t.Error(WRONG_VALUE_RETURNED)
	}
}

func TestMessage_MarkAsRead(t *testing.T) {
	SetAPI("", "test", "")
	m := Message{}
	err := m.MarkAsRead()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMessage_Reply(t *testing.T) {
	SetAPI("", "test", "")
	m := Message{}
	mid, err := m.Reply("ok")
	if err != nil {
		t.Error(err.Error())
	}
	if mid == 0 {
		t.Error(WRONG_VALUE_RETURNED)
	}
	m = Message{ChatID: 1}
	mid, err = m.Reply("ok")
	if err != nil {
		t.Error(err.Error())
	}
	if mid == 0 {
		t.Error(WRONG_VALUE_RETURNED)
	}
}

func TestVkAPI_SendChatMessage(t *testing.T) {
	SetAPI("", "test", "")
	err := API.SendChatMessage(1, "ok")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestVkAPI_SendMessage(t *testing.T) {
	SetAPI("", "test", "")
	err := API.SendMessage(1, "ok")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestVkAPI_User(t *testing.T) {
	SetAPI("", "test", "")
	u, err := API.User(1)
	if err != nil {
		t.Error(err.Error())
	}
	if u.ID == 0 {
		t.Error(WRONG_VALUE_RETURNED)
	}
}

func TestVkAPI_NotifyAdmin(t *testing.T) {
	SetAPI("", "test", "")
	API.AdminID = 1
	err := API.NotifyAdmin("ok")
	if err != nil {
		t.Error(err.Error())
	}
}
