package govkbot

import (
	"testing"
)

const (
	wrongValueReturned = "Wrong value returned"
)

//func TestCall(t *testing.T) {
//	r := SimpleResponse{}
//	err := API.CallMethod("utils.getServerTime", H{}, &r)
//	if err != nil {
//		t.Error(err.Error())
//	}
//}

//func TestVkAPI_Call(t *testing.T) {
//	SetDebug(true)
//	log.SetOutput(ioutil.Discard)
//	r := SimpleResponse{}
//	err := API.CallMethod("messages.get", H{}, &r)
//	SetDebug(false)
//	log.SetOutput(os.Stdout)
//	if err == nil {
//		t.Error("no error returned")
//	}
//}

func TestNoJSON(t *testing.T) {
	SetAPI("", "test", "")
	r := SimpleResponse{}
	err := API.CallMethod("nojson", H{}, &r)
	if err == nil {
		t.Error("no error returned")
	}
	if _, ok := err.(*ResponseError); !ok {
		t.Error("wrong error type")
	}
}

func TestVkError(t *testing.T) {
	SetAPI("", "test", "")
	r := SimpleResponse{}
	err := API.CallMethod("vkerr", H{}, &r)
	if err == nil {
		t.Error("no error returned")
	}
	if _, ok := err.(*VKError); !ok {
		t.Error("wrong error type")
	}
}

func TestVkAPI_Me(t *testing.T) {
	SetAPI("", "test", "")
	me, err := API.Me()
	if err != nil {
		t.Error(err.Error())
	}
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
		t.Error(wrongValueReturned)
	}
}

func TestVkAPI_DeleteFriend(t *testing.T) {
	SetAPI("", "test", "")
	ok := API.DeleteFriend(1)
	if !ok {
		t.Error(wrongValueReturned)
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

func TestVKBot_Reply(t *testing.T) {
	SetAPI("", "test", "")
	bot := API.NewBot()
	m := Message{}
	mid, err := bot.Reply(&m, Reply{Msg: "ok"})
	if err != nil {
		t.Error(err.Error())
	}
	if mid == 0 {
		t.Error(wrongValueReturned)
	}
	m = Message{ChatID: 1}
	mid, err = bot.Reply(&m, Reply{Msg: "ok"})
	if err != nil {
		t.Error(err.Error())
	}
	if mid == 0 {
		t.Error(wrongValueReturned)
	}
}

func TestVkAPI_SendChatMessage(t *testing.T) {
	SetAPI("", "test", "")
	_, err := API.SendChatMessage(1, "ok")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestVkAPI_SendMessage(t *testing.T) {
	SetAPI("", "test", "")
	_, err := API.SendMessage(1, "ok")
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
		t.Error(wrongValueReturned)
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

func TestMessage_GetMentions(t *testing.T) {
	testStr := "/who [id373336876|@sociobesed]"
	m := Message{Body: testStr}
	mentions := m.GetMentions()
	if len(mentions) != 1 {
		t.Error("wrong mentions count")
	}
	if mentions[0].ID != 373336876 {
		t.Error("wrong mention")
	}
}
