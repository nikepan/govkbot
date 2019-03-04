package govkbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// H - simple object struct
type H map[string]string

// VkAPI - api config
type VkAPI struct {
	Token           string
	URL             string
	Ver             string
	UID             int
	GroupID         int
	Lang            string
	HTTPS           bool
	AdminID         int
	MessagesCount   int
	RequestInterval int
	DEBUG           bool
}

const (
	apiUsersGet                       = "users.get"
	apiGroupsGet                      = "groups.getById"
	apiMessagesGet                    = "messages.get"
	apiMessagesGetChat                = "messages.getChat"
	apiMessagesGetConversationsById   = "messages.getConversationsById"
	apiMessagesGetChatUsers           = "messages.getChatUsers"
	apiMessagesGetConversationMembers = "messages.getConversationMembers"
	apiMessagesSend                   = "messages.send"
	apiMessagesMarkAsRead             = "messages.markAsRead"
	apiFriendsGetRequests             = "friends.getRequests"
	apiFriendsAdd                     = "friends.add"
	apiFriendsDelete                  = "friends.delete"
)

func (api *VkAPI) IsGroup() bool {
	if api.GroupID != 0 {
		return true
	} else if api.UID != 0 {
		return false
	}

	g, err := API.CurrentGroup()
	if err != nil || g.ID == 0 {
		u, err := API.Me()
		if err != nil || u == nil {
			fmt.Printf("Get current user/group error %+v\n", err)
		} else {
			api.UID = u.ID
		}
	} else {
		api.GroupID = g.ID
	}
	return api.GroupID != 0
}

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

// CallMethod - call VK API method by name to interfce
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

// CurrentGroup - get current group info
func (api *VkAPI) CurrentGroup() (*User, error) {

	r := UsersResponse{}
	err := api.CallMethod(apiGroupsGet, H{"fields": "screen_name"}, &r)

	if len(r.Response) > 0 {
		debugPrint("me: %+v - %+v\n", r.Response[0].ID, r.Response[0].ScreenName)
		return r.Response[0], err
	}
	return nil, err
}

// GetChatInfo - returns Chat info by id
func (api *VkAPI) GetUserChatInfo(chatID int) (*ChatInfo, error) {
	r := ChatInfoResponse{}
	err := api.CallMethod(apiMessagesGetChat, H{
		"chat_id": strconv.Itoa(chatID),
		"fields":  "photo,city,country,sex,bdate",
	}, &r)

	return &r.Response, err
}

// GetChatInfo - returns Chat info by id
func (api *VkAPI) GetConversation(chatID int) (*ChatInfo, error) {
	r := ConversationsResponse{}
	err := api.CallMethod(apiMessagesGetConversationsById, H{
		"peer_ids": strconv.Itoa(chatID),
		"extended": "1",
		"fields":   "photo,city,country,sex,bdate",
	}, &r)

	c := r.Response.Items[0]
	chat := ChatInfo{}
	chat.ID = c.Peer.ID
	chat.Type = c.Peer.Type
	chat.Title = c.ChatSettings.Title
	chat.AdminID = c.ChatSettings.OwnerID

	chat.Users = make([]*User, 0)
	for _, u := range r.Response.Profiles {
		user := User{}
		user.ID = u.ID
		user.FirstName = u.FirstName
		user.LastName = u.LastName
		user.Photo = u.Photo
		user.City = u.City
		user.Country = u.Country
		user.Sex = u.Sex
		chat.Users = append(chat.Users, &user)
	}

	return &chat, err
}

// GetChatInfo - returns Chat info by id
func (api *VkAPI) GetChatInfo(chatID int) (*ChatInfo, error) {
	if api.IsGroup() {
		return api.GetConversation(chatID)
	}
	return api.GetUserChatInfo(chatID)
}

func (api *VkAPI) GetChatFullInfo(chatID int) (*ChatInfo, error) {
	info, err := api.GetChatInfo(chatID)
	if err != nil {
		return nil, err
	}
	members, err := api.GetChatUsers(chatID)
	if err != nil {
		return nil, err
	}
	if info != nil && members != nil {
		info.Users = members
	}
	return info, err
}

// GetChatUsers - get chat users
func (api *VkAPI) GetUserChatUsers(chatID int) (users []*User, err error) {

	r := UsersResponse{}
	err = api.CallMethod(apiMessagesGetChatUsers, H{
		"chat_id": strconv.Itoa(chatID),
		"fields":  "photo,city,country,sex,bdate",
	}, &r)

	return r.Response, err
}

// GetChatUsers - get chat users
func (api *VkAPI) GetConversationMembers(chatID int) (users []*User, err error) {

	r := MembersResponse{}

	err = api.CallMethod(apiMessagesGetConversationMembers, H{
		"peer_id": strconv.Itoa(chatID),
		"fields":  "photo,city,country,sex,bdate",
	}, &r)

	users = make([]*User, 0)
	for _, u := range r.Response.Profiles {
		user := User{}
		user.ID = u.ID
		user.FirstName = u.FirstName
		user.LastName = u.LastName
		user.ScreenName = u.ScreenName
		user.Photo = u.Photo
		user.Sex = u.Sex
		user.City = u.City
		user.Country = u.Country
		user.BDate = u.BDate
		for _, i := range r.Response.Items {
			if i.MemberID == u.ID {
				user.IsAdmin = i.IsAdmin
				user.IsOwner = i.IsOwner
				user.InvitedBy = i.InvitedBy
				users = append(users, &user)
				break
			}
		}
	}

	return users, err
}

func (api *VkAPI) GetChatUsers(chatID int) (users []*User, err error) {
	if api.IsGroup() {
		return api.GetConversationMembers(chatID)
	}
	return api.GetUserChatUsers(chatID)
}

func FindUser(users []*User, ID int) *User {
	for _, u := range users {
		if u.ID == ID {
			return u
		}
	}
	return nil
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
		"fields":   "sex,screen_name, city, country, bdate",
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

	//r := SimpleResponse{}
	//err = API.CallMethod(apiMessagesMarkAsRead, H{"message_ids": strconv.Itoa(m.ID)}, &r)
	return nil
}

func (m Message) GetMentions() []Mention {
	mentions := make([]Mention, 0)
	ok := false
	ustr := ""
	uname := ""
	msg := m.Body
	i := strings.Index(msg, "[id")
	imax := len(msg)
	if i >= 0 {
		i += 2
		for {
			i++
			if i >= imax {
				break
			}
			if '0' <= msg[i] && msg[i] <= '9' {
				ustr += string(msg[i])
			} else {
				ok = true
				break
			}
		}
		w := 1
		var runeValue rune
		for {
			i += w
			if i >= imax {
				ok = false
				break
			}
			runeValue, w = utf8.DecodeRuneInString(msg[i:])
			if msg[i] != ']' {
				if msg[i] != '|' {
					uname += fmt.Sprintf("%c", runeValue) //string(msg[i])
				}
			} else {
				ok = true
				break
			}
		}
	}
	if ok {
		u, err := strconv.Atoi(ustr)
		if err == nil {
			mentions = append(mentions, Mention{u, uname})
		}
	}
	return mentions
}

func (api *VkAPI) GetRandomID() string {
	return strconv.FormatUint(uint64(rand.Uint32()), 10)
}

//SendAdvancedPeerMessage sending a message to chat
func (api *VkAPI) SendAdvancedPeerMessage(peerID int64, message Reply) (id int, err error) {
	r := SimpleResponse{}
	params := H{
		"peer_id":          strconv.FormatInt(peerID, 10),
		"message":          message.Msg,
		"dont_parse_links": "1",
		"random_id":        api.GetRandomID(),
	}
	if message.Keyboard != nil {
		keyboard, err := json.Marshal(message.Keyboard)
		if err != nil {
			fmt.Printf("ERROR encode keyboard %+v\n", message.Keyboard)
		} else {
			params["keyboard"] = string(keyboard)
		}
	}
	err = api.CallMethod(apiMessagesSend, params, &r)
	return r.Response, err
}

//SendPeerMessage sending a message to chat
func (api *VkAPI) SendPeerMessage(peerID int64, msg string) (id int, err error) {
	r := SimpleResponse{}
	err = api.CallMethod(apiMessagesSend, H{
		"peer_id":          strconv.FormatInt(peerID, 10),
		"message":          msg,
		"dont_parse_links": "1",
		"random_id":        api.GetRandomID(),
	}, &r)
	return r.Response, err
}

//SendChatMessage sending a message to chat
func (api *VkAPI) SendChatMessage(chatID int, msg string) (id int, err error) {
	r := SimpleResponse{}
	err = api.CallMethod(apiMessagesSend, H{
		"chat_id":          strconv.Itoa(chatID),
		"message":          msg,
		"dont_parse_links": "1",
		"random_id":        api.GetRandomID(),
	}, &r)
	return r.Response, err
}

//SendMessage sending a message to user
func (api *VkAPI) SendMessage(userID int, msg string) (id int, err error) {
	r := SimpleResponse{}
	if msg != "" {
		err = api.CallMethod(apiMessagesSend, H{
			"user_id":          strconv.Itoa(userID),
			"message":          msg,
			"dont_parse_links": "1",
			"random_id":        api.GetRandomID(),
		}, &r)
	}
	return r.Response, err
}

func NewButton(label string, payload interface{}) Button {
	button := Button{}
	button.Action.Type = "text"
	button.Action.Label = label
	button.Action.Payload = "{}"
	if payload != nil {
		jPayoad, err := json.Marshal(payload)
		if err == nil {
			button.Action.Payload = string(jPayoad)
		}
	}
	button.Color = "default"
	return button
}

// NotifyAdmin - send notify to admin
func (api *VkAPI) NotifyAdmin(msg string) (err error) {
	if api.AdminID != 0 {
		_, err = api.SendMessage(api.AdminID, msg)
	}
	return err
}
