package govkbot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// VkAPI - api config
type VkAPI struct {
	Token   string
	Url     string
	Ver     string
	Uid     int
	AdminID int
}

// Call - main api call method
func (api *VkAPI) Call(method string, parameters url.Values) []byte {
	p := "?" + parameters.Encode()
	//log.Printf("vk req: %+v\n", api.Url+method+p)
	parameters.Add("access_token", api.Token)
	parameters.Add("v", api.Ver)
	p = "?" + parameters.Encode()
	resp, err := http.Get(api.Url + method + p)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	//log.Printf("vk resp: %+v\n", string(buf))
	return buf
}

// GetMessages - get user messages
func (api *VkAPI) GetMessages() *Messages {

	p := url.Values{}

	buf := api.Call("messages.get", p)

	m := MessagesResponse{}
	json.Unmarshal(buf, &m)

	return &m.Response
}

// Me - get current user info
func (api *VkAPI) Me() *User {
	p := url.Values{}

	buf := api.Call("users.get", p)

	log.Printf("me: %+v\n", string(buf))

	u := UsersResponse{}
	json.Unmarshal(buf, &u)

	return u.Response[0]
}

// GetChatUsers - get chat users
func (api *VkAPI) GetChatUsers(chatID int) []*User {
	p := url.Values{}
	p.Add("chat_id", strconv.Itoa(chatID))
	p.Add("fields", "photo")

	buf := api.Call("messages.getChatUsers", p)

	log.Printf("users: %+v\n", string(buf))

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

	buf := api.Call("friends.getRequests", p)
	u := FriendRequestsResponse{}
	json.Unmarshal(buf, &u)

	return u.Response.Items
}

// AddFriend - add friend
func (api *VkAPI) AddFriend(uid int) bool {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(uid))

	buf := api.Call("friends.add", p)
	u := SimpleResponse{}
	json.Unmarshal(buf, &u)

	return u.Response == 1
}

// DeleteFriend - delete friend
func (api *VkAPI) DeleteFriend(uid int) bool {
	p := url.Values{}
	p.Add("user_id", strconv.Itoa(uid))

	buf := api.Call("friends.delete", p)
	u := SimpleResponse{}
	json.Unmarshal(buf, &u)

	return u.Response == 1
}

// User - get simple user info
func (api *VkAPI) User(uid int) (User, bool) {
	p := url.Values{}
	p.Add("user_ids", strconv.Itoa(uid))
	p.Add("fields", "sex")

	buf := api.Call("users.get", p)

	u := UsersResponse{}
	json.Unmarshal(buf, &u)
	if len(u.Response) > 0 {
		return *u.Response[0], true
	} else {
		return User{}, false
	}
}

// MarkAsRead - mark message as read
func (m Message) MarkAsRead() {

	p := url.Values{}
	p.Add("message_ids", strconv.Itoa(m.ID))

	_ = API.Call("messages.markAsRead", p)

}

// Reply - reply message
func (m Message) Reply(msg string) {
	p := url.Values{}
	if m.ChatID != 0 {
		p.Add("chat_id", strconv.Itoa(m.ChatID))
	} else {
		p.Add("user_id", strconv.Itoa(m.UserID))
	}
	p.Add("message", msg)

	_ = API.Call("messages.send", p)
}
