package govkbot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// VkAPI - api config
type VkAPI struct {
	Token           string
	URL             string
	Ver             string
	UID             int
	AdminID         int
	MessagesCount   int
	RequestInterval int
	DEBUG           bool
}

// Call - main api call method
func (api *VkAPI) Call(method string, parameters url.Values) ([]byte, error) {
	p := "?" + parameters.Encode()
	if api.DEBUG {
		log.Printf("vk req: %+v\n", api.URL+method+p)
	}
	parameters.Add("access_token", api.Token)
	parameters.Add("v", api.Ver)
	p = "?" + parameters.Encode()
	resp, err := http.PostForm(api.URL+method, parameters)
	if err != nil {
		log.Println(err.Error())
		time.Sleep(time.Duration(time.Millisecond * time.Duration(api.RequestInterval)))
		return nil, err
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(time.Millisecond * time.Duration(api.RequestInterval)))
	if api.DEBUG {
		log.Printf("vk resp: %+v\n", string(buf))
	}

	u := SimpleResponse{}
	json.Unmarshal(buf, &u)
	if u.Error != nil {
		log.Printf("%+v\n", u.Error.ErrorMsg)
		return buf, errors.New(u.Error.ErrorMsg)
	}

	return buf, nil
}

// GetMessages - get user messages
func (api *VkAPI) GetMessages(count int, offset int) (*Messages, error) {

	p := url.Values{}
	p.Add("count", strconv.Itoa(count))
	p.Add("offset", strconv.Itoa(offset))

	buf, _ := api.Call("messages.get", p)

	m := MessagesResponse{}
	json.Unmarshal(buf, &m)
	if m.Error != nil {
		return &m.Response, errors.New(m.Error.ErrorMsg)
	}

	return &m.Response, nil
}

// Me - get current user info
func (api *VkAPI) Me() *User {
	p := url.Values{}

	buf, _ := api.Call("users.get", p)

	if api.DEBUG {
		log.Printf("me: %+v\n", string(buf))
	}
	u := UsersResponse{}
	json.Unmarshal(buf, &u)

	return u.Response[0]
}

// GetChatInfo - returns Chat info by id
func (api *VkAPI) GetChatInfo(chatID int) (*ChatInfo, *VKError) {
	p := url.Values{}
	p.Add("chat_id", strconv.Itoa(chatID))
	p.Add("fields", "photo,city,country")
	buf, _ := api.Call("messages.getChat", p)
	u := ChatInfoResponse{}
	json.Unmarshal(buf, &u)
	if u.Error != nil {
		return nil, u.Error
	}
	return &u.Response, nil
}

// GetChatUsers - get chat users
func (api *VkAPI) GetChatUsers(chatID int) []*User {
	p := url.Values{}
	p.Add("chat_id", strconv.Itoa(chatID))
	p.Add("fields", "photo")

	buf, _ := api.Call("messages.getChatUsers", p)

	if api.DEBUG {
		log.Printf("users: %+v\n", string(buf))
	}
	u := UsersResponse{}
	json.Unmarshal(buf, &u)

	return u.Response
}

// GetFriendRequests - get friend requests
func (api *VkAPI) GetFriendRequests(out bool) []int {
	p := url.Values{}
	if out {
		p.Add("out", "1")
	}

	buf, _ := api.Call("friends.getRequests", p)
	u := FriendRequestsResponse{}
	json.Unmarshal(buf, &u)

	return u.Response.Items
}

// AddFriend - add friend
func (api *VkAPI) AddFriend(uid int) bool {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(uid))

	buf, _ := api.Call("friends.add", p)
	u := SimpleResponse{}
	json.Unmarshal(buf, &u)

	return u.Response == 1
}

// DeleteFriend - delete friend
func (api *VkAPI) DeleteFriend(uid int) bool {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(uid))

	buf, _ := api.Call("friends.delete", p)
	u := SimpleResponse{}
	json.Unmarshal(buf, &u)

	return u.Response == 1
}

// User - get simple user info
func (api *VkAPI) User(uid int) (User, bool) {
	p := url.Values{}
	p.Add("user_ids", strconv.Itoa(uid))
	p.Add("fields", "sex")

	buf, _ := api.Call("users.get", p)

	u := UsersResponse{}
	json.Unmarshal(buf, &u)
	if len(u.Response) > 0 {
		return *u.Response[0], true
	}
	return User{}, false
}

// MarkAsRead - mark message as read
func (m Message) MarkAsRead() {

	p := url.Values{}
	p.Add("message_ids", strconv.Itoa(m.ID))

	_, _ = API.Call("messages.markAsRead", p)

}

//SendChatMessage sending a message to chat
func (api *VkAPI) SendChatMessage(chatID int, msg string) (err error) {
	p := url.Values{}
	p.Add("chat_id", strconv.Itoa(chatID))
	p.Add("message", msg)
	_, err = api.Call("messages.send", p)
	return err
}

//SendMessage sending a message to user
func (api *VkAPI) SendMessage(userID int, msg string) (err error) {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(userID))
	p.Add("message", msg)
	_, err = api.Call("messages.send", p)
	return err
}

// Reply - reply message
func (m Message) Reply(msg string) (err error) {
	p := url.Values{}
	if m.ChatID != 0 {
		p.Add("chat_id", strconv.Itoa(m.ChatID))
	} else {
		p.Add("user_id", strconv.Itoa(m.UserID))
	}
	//p.Add("forward_messages", strconv.Itoa(m.ID))
	p.Add("message", msg)

	_, err = API.Call("messages.send", p)

	if err != nil {
		log.Printf("%+v\n", err.Error())
	}
	return err
}

// NotifyAdmin - send notify to admin
func (api *VkAPI) NotifyAdmin(msg string) (err error) {
	if api.AdminID != 0 {
		return api.SendMessage(api.AdminID, msg)
	}
	return nil
}
