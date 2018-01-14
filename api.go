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

type H map[string]string

// VkAPI - api config
type VkAPI struct {
	Token           string
	URL             string
	Ver             string
	UID             int
	Lang            string
	HTTPS           bool
	AdminID         int
	MessagesCount   int
	RequestInterval int
	DEBUG           bool
}

const (
	apiUsersGet             = "users.get"
	apiMessagesGet          = "messages.get"
	apiMessagesGetChat      = "messages.getChat"
	apiMessagesGetChatUsers = "messages.getChatUsers"
	apiMessagesSend         = "messages.send"
	apiMessagesMarkAsRead   = "messages.markAsRead"
	apiFriendsGetRequests   = "friends.getRequests"
	apiFriendsAdd           = "friends.add"
	apiFriendsDelete        = "friends.delete"
)

// Call - main api call method
func (api *VkAPI) Call(method string, params map[string]string) ([]byte, error) {
	debugPrint("vk req: %+v params: %+v\n", api.URL+method, params)
	params["access_token"] = api.Token
	params["v"] = api.Ver
	if api.Lang != "" {
		params["lang"] = api.Lang
	}
	if api.HTTPS {
		params["https"] = "1"
	}

	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

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
	buf, err := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(time.Millisecond * time.Duration(api.RequestInterval)))
	debugPrint("vk resp: %+v\n", string(buf))

	return buf, err
}

func (api *VkAPI) CallMethod(method string, params map[string]string, result interface{}) error {
	buf, err := api.Call(method, params)
	if err != nil {
		return err
	}
	r := ErrorResponse{}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return &ResponseError{errors.New("vkapi: vk response is not json"), string(buf)}
	}
	if r.Error != nil {
		debugPrint("%+v\n", r.Error.ErrorMsg)
		return r.Error
	}

	err = json.Unmarshal(buf, result)
	return err
}

// GetMessages - get user messages (up to 200)
func (api *VkAPI) GetMessages(count int, offset int) (*Messages, error) {

	m := MessagesResponse{}
	err := api.CallMethod(apiMessagesGet, H{
		"count":  strconv.Itoa(count),
		"offset": strconv.Itoa(offset),
	}, &m)

	return &m.Response, err
}

// Me - get current user info
func (api *VkAPI) Me() (*User, error) {

	r := UsersResponse{}
	err := api.CallMethod(apiUsersGet, H{"fields": "screen_name"}, &r)

	if len(r.Response) > 0 {
		debugPrint("me: %+v - %+v\n", r.Response[0].ID, r.Response[0].ScreenName)
		return r.Response[0], err
	}
	return nil, err
}

// GetChatInfo - returns Chat info by id
func (api *VkAPI) GetChatInfo(chatID int) (*ChatInfo, error) {
	r := ChatInfoResponse{}
	err := api.CallMethod(apiMessagesGetChat, H{
		"chat_id": strconv.Itoa(chatID),
		"fields":  "photo,city,country",
	}, &r)

	return &r.Response, err
}

// GetChatUsers - get chat users
func (api *VkAPI) GetChatUsers(chatID int) (users []*User, err error) {

	r := UsersResponse{}
	err = api.CallMethod(apiMessagesGetChatUsers, H{
		"chat_id": strconv.Itoa(chatID),
		"fields":  "photo",
	}, &r)

	return r.Response, err
}

// GetFriendRequests - get friend requests
func (api *VkAPI) GetFriendRequests(out bool) (friends []int, err error) {
	p := H{}
	if out {
		p["out"] = "1"
	}

	r := FriendRequestsResponse{}
	err = api.CallMethod(apiFriendsGetRequests, p, &r)

	return r.Response.Items, err
}

// AddFriend - add friend
func (api *VkAPI) AddFriend(uid int) bool {

	r := SimpleResponse{}
	err := api.CallMethod(apiFriendsAdd, H{"user_id": strconv.Itoa(uid)}, &r)
	if err != nil {
		return false
	}

	return r.Response == 1
}

// DeleteFriend - delete friend
func (api *VkAPI) DeleteFriend(uid int) bool {

	u := FriendDeleteResponse{}
	err := api.CallMethod(apiFriendsDelete, H{"user_id": strconv.Itoa(uid)}, &u)

	if err != nil {
		return false
	}

	return u.Response["success"] == 1
}

// User - get simple user info
func (api *VkAPI) User(uid int) (*User, error) {

	r := UsersResponse{}
	err := api.CallMethod(apiUsersGet, H{
		"user_ids": strconv.Itoa(uid),
		"fields":   "sex,screen_name",
	}, &r)

	if err != nil {
		return nil, err
	}
	if len(r.Response) > 0 {
		return r.Response[0], err
	}
	return nil, errors.New("no users returned")
}

// MarkAsRead - mark message as read
func (m Message) MarkAsRead() (err error) {

	r := SimpleResponse{}
	err = API.CallMethod(apiMessagesMarkAsRead, H{"message_ids": strconv.Itoa(m.ID)}, &r)
	return err
}

//SendChatMessage sending a message to chat
func (api *VkAPI) SendChatMessage(chatID int, msg string) (id int, err error) {
	r := SimpleResponse{}
	err = api.CallMethod(apiMessagesSend, H{
		"chat_id": strconv.Itoa(chatID),
		"message": msg,
	}, &r)
	return r.Response, err
}

//SendMessage sending a message to user
func (api *VkAPI) SendMessage(userID int, msg string) (id int, err error) {
	r := SimpleResponse{}
	if msg != "" {
		err = api.CallMethod(apiMessagesSend, H{
			"user_id": strconv.Itoa(userID),
			"message": msg,
		}, &r)
	}
	return r.Response, err
}

// Reply - reply message
func (m Message) Reply(msg string) (id int, err error) {
	if m.ChatID != 0 {
		return API.SendChatMessage(m.ChatID, msg)
	}
	return API.SendMessage(m.UserID, msg)
}

// NotifyAdmin - send notify to admin
func (api *VkAPI) NotifyAdmin(msg string) (err error) {
	if api.AdminID != 0 {
		_, err = api.SendMessage(api.AdminID, msg)
	}
	return err
}
