package govkbot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

const (
	API_USERS_GET               = "users.get"
	API_MESSAGES_GET            = "messages.get"
	API_MESSAGES_GET_CHAT       = "messages.getChat"
	API_MESSAGES_GET_CHAT_USERS = "messages.getChatUsers"
	API_MESSAGES_SEND           = "messages.send"
	API_MESSAGES_MARK_AS_READ   = "messages.markAsRead"
	API_FRIENDS_GET_REQUESTS    = "friends.getRequests"
	API_FRIENDS_ADD             = "friends.add"
	API_FRIENDS_DELETE          = "friends.delete"
)

// Call - main api call method
func (api *VkAPI) Call(method string, parameters url.Values) ([]byte, error) {
	p := "?" + parameters.Encode()
	debugPrint("vk req: %+v\n", api.URL+method+p)
	parameters.Add("access_token", api.Token)
	parameters.Add("v", api.Ver)

	if api.URL == "test" {
		content, err := ioutil.ReadFile("./mocks/" + method + ".json")
		return content, err
	}
	resp, err := http.PostForm(api.URL+method, parameters)
	if err != nil {
		debugPrint("%+v\n", err.Error())
		time.Sleep(time.Duration(time.Millisecond * time.Duration(api.RequestInterval)))
		return nil, err
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(time.Millisecond * time.Duration(api.RequestInterval)))
	debugPrint("vk resp: %+v\n", string(buf))

	u := SimpleResponse{}
	json.Unmarshal(buf, &u)
	if u.Error != nil {
		debugPrint("%+v\n", u.Error.ErrorMsg)
		return buf, errors.New(u.Error.ErrorMsg)
	}

	return buf, nil
}

// GetMessages - get user messages
func (api *VkAPI) GetMessages(count int, offset int) (*Messages, error) {

	p := url.Values{}
	p.Add("count", strconv.Itoa(count))
	p.Add("offset", strconv.Itoa(offset))

	buf, _ := api.Call(API_MESSAGES_GET, p)

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

	buf, _ := api.Call(API_USERS_GET, p)
	debugPrint("me: %+v\n", string(buf))

	u := UsersResponse{}
	json.Unmarshal(buf, &u)
	if len(u.Response) > 0 {
		return u.Response[0]
	}
	return nil
}

// GetChatInfo - returns Chat info by id
func (api *VkAPI) GetChatInfo(chatID int) (*ChatInfo, *VKError) {
	p := url.Values{}
	p.Add("chat_id", strconv.Itoa(chatID))
	p.Add("fields", "photo,city,country")
	buf, _ := api.Call(API_MESSAGES_GET_CHAT, p)
	u := ChatInfoResponse{}
	json.Unmarshal(buf, &u)
	if u.Error != nil {
		return nil, u.Error
	}
	return &u.Response, nil
}

// GetChatUsers - get chat users
func (api *VkAPI) GetChatUsers(chatID int) (users []*User, err error) {
	p := url.Values{}
	p.Add("chat_id", strconv.Itoa(chatID))
	p.Add("fields", "photo")

	buf, err := api.Call(API_MESSAGES_GET_CHAT_USERS, p)
	if err != nil {
		return nil, err
	}

	debugPrint("users: %+v\n", string(buf))
	u := UsersResponse{}
	json.Unmarshal(buf, &u)

	return u.Response, nil
}

// GetFriendRequests - get friend requests
func (api *VkAPI) GetFriendRequests(out bool) (friends []int, err error) {
	p := url.Values{}
	if out {
		p.Add("out", "1")
	}

	buf, err := api.Call(API_FRIENDS_GET_REQUESTS, p)
	u := FriendRequestsResponse{}
	json.Unmarshal(buf, &u)

	return u.Response.Items, err
}

// AddFriend - add friend
func (api *VkAPI) AddFriend(uid int) bool {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(uid))

	buf, _ := api.Call(API_FRIENDS_ADD, p)
	u := SimpleResponse{}
	json.Unmarshal(buf, &u)

	return u.Response == 1
}

// DeleteFriend - delete friend
func (api *VkAPI) DeleteFriend(uid int) bool {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(uid))

	buf, _ := api.Call(API_FRIENDS_DELETE, p)
	u := FriendDeleteResponse{}
	json.Unmarshal(buf, &u)

	ok := u.Response["success"] == 1

	return ok
}

// User - get simple user info
func (api *VkAPI) User(uid int) (*User, error) {
	p := url.Values{}
	p.Add("user_ids", strconv.Itoa(uid))
	p.Add("fields", "sex")

	buf, err := api.Call(API_USERS_GET, p)
	if err != nil {
		return nil, err
	}

	u := UsersResponse{}
	json.Unmarshal(buf, &u)
	if len(u.Response) > 0 {
		return u.Response[0], nil
	}
	return nil, errors.New("no users returned")
}

// MarkAsRead - mark message as read
func (m Message) MarkAsRead() (err error) {

	p := url.Values{}
	p.Add("message_ids", strconv.Itoa(m.ID))

	_, err = API.Call(API_MESSAGES_MARK_AS_READ, p)
	return err

}

//SendChatMessage sending a message to chat
func (api *VkAPI) SendChatMessage(chatID int, msg string) (err error) {
	p := url.Values{}
	p.Add("chat_id", strconv.Itoa(chatID))
	p.Add("message", msg)
	_, err = api.Call(API_MESSAGES_SEND, p)
	return err
}

//SendMessage sending a message to user
func (api *VkAPI) SendMessage(userID int, msg string) (err error) {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(userID))
	p.Add("message", msg)
	_, err = api.Call(API_MESSAGES_SEND, p)
	return err
}

// Reply - reply message
func (m Message) Reply(msg string) (id int, err error) {
	p := url.Values{}
	if m.ChatID != 0 {
		p.Add("chat_id", strconv.Itoa(m.ChatID))
	} else {
		p.Add("user_id", strconv.Itoa(m.UserID))
	}
	//p.Add("forward_messages", strconv.Itoa(m.ID))
	p.Add("message", msg)

	buf, err := API.Call(API_MESSAGES_SEND, p)
	r := SimpleResponse{}
	json.Unmarshal(buf, &r)

	return r.Response, err
}

// NotifyAdmin - send notify to admin
func (api *VkAPI) NotifyAdmin(msg string) (err error) {
	if api.AdminID != 0 {
		return api.SendMessage(api.AdminID, msg)
	}
	return nil
}
