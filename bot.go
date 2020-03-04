package govkbot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// VKBot - bot config
type VKBot struct {
	msgRoutes        map[string]msgRoute
	actionRoutes     map[string]func(*Message) string
	cmdHandlers      map[string]func(*Message) string
	msgHandlers      map[string]func(*Message) string
	errorHandler     func(*Message, error)
	LastMsg          int
	lastUserMessages map[int]int
	lastChatMessages map[int]int
	autoFriend       bool
	IgnoreBots       bool
	LongPoll         LongPollServer
	API              *VkAPI
}

type msgRoute struct {
	SimpleHandler func(*Message) string
	Handler       func(*Message) Reply
}

// NewBot - create new instance of bot
func (api *VkAPI) NewBot() *VKBot {
	if api.IsGroup() {
		return &VKBot{
			msgRoutes:        make(map[string]msgRoute),
			actionRoutes:     make(map[string]func(*Message) string),
			lastUserMessages: make(map[int]int),
			lastChatMessages: make(map[int]int),
			LongPoll:         NewGroupLongPollServer(API.RequestInterval),
			API:              api,
		}
	}
	return &VKBot{
		msgRoutes:        make(map[string]msgRoute),
		actionRoutes:     make(map[string]func(*Message) string),
		lastUserMessages: make(map[int]int),
		lastChatMessages: make(map[int]int),
		LongPoll:         NewUserLongPollServer(false, longPollVersion, API.RequestInterval),
		API:              api,
	}
}

// ListenUser - listen User VK API (deprecated)
func (bot *VKBot) ListenUser(api *VkAPI) error {
	bot.LongPoll = NewUserLongPollServer(false, longPollVersion, API.RequestInterval)
	go bot.friendReceiver()

	c := time.Tick(3 * time.Second)
	for range c {
		bot.MainRoute()
	}
	return nil
}

// ListenGroup - listen group VK API
func (bot *VKBot) ListenGroup(api *VkAPI) error {
	bot.LongPoll = NewGroupLongPollServer(API.RequestInterval)
	c := time.Tick(3 * time.Second)
	for range c {
		bot.MainRoute()
	}
	return nil
}

// HandleMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
func (bot *VKBot) HandleMessage(command string, handler func(*Message) string) {
	bot.msgRoutes[command] = msgRoute{SimpleHandler: handler}
}

// HandleAdvancedMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
func (bot *VKBot) HandleAdvancedMessage(command string, handler func(*Message) Reply) {
	bot.msgRoutes[command] = msgRoute{Handler: handler}
}

// HandleAction - add action handler.
// Function must return string to reply or "" (if no reply)
func (bot *VKBot) HandleAction(command string, handler func(*Message) string) {
	bot.actionRoutes[command] = handler
}

// HandleError - add error handler
func (bot *VKBot) HandleError(handler func(*Message, error)) {
	bot.errorHandler = handler
}

// SetAutoFriend - auto add friends
func (bot *VKBot) SetAutoFriend(af bool) {
	bot.autoFriend = af
}

// GetMessages - request unread messages from VK (more than 200)
func (bot *VKBot) GetMessages() ([]*Message, error) {
	var allMessages []*Message
	lastMsg := bot.LastMsg
	offset := 0
	var err error
	var messages *Messages
	for {
		messages, err = API.GetMessages(API.MessagesCount, offset)
		if len(messages.Items) > 0 {
			if messages.Items[0].ID > lastMsg {
				lastMsg = messages.Items[0].ID
			}
		}
		allMessages = append(allMessages, messages.Items...)
		if bot.LastMsg > 0 {
			if len(messages.Items) > 0 {
				if messages.Items[len(messages.Items)-1].ID <= bot.LastMsg {
					bot.LastMsg = lastMsg
					break
				}
			} else {
				break
			}
			offset += API.MessagesCount
		} else {
			bot.LastMsg = lastMsg
			break
		}
	}
	if offset > 0 {
		API.NotifyAdmin("many messages in interval. offset: " + strconv.Itoa(offset))
	}
	return allMessages, err
}

//RouteAction routes an action
func (bot *VKBot) RouteAction(m *Message) (replies []string, err error) {
	if m.Action != "" {
		debugPrint("route action: %+v\n", m.Action)
		for k, v := range bot.actionRoutes {
			if m.Action == k {
				msg := v(m)
				if msg != "" {
					replies = append(replies, msg)
				}
			}
		}
	}
	return replies, nil
}

// RouteMessage routes single message
func (bot *VKBot) RouteMessage(m *Message) (replies []Reply, err error) {
	message := strings.TrimSpace(strings.ToLower(m.Body))
	if HasPrefix(message, "/ ") {
		message = "/" + TrimPrefix(message, "/ ")
	}
	if m.Action != "" {
		actionReplies, err := bot.RouteAction(m)
		for _, r := range actionReplies {
			replies = append(replies, Reply{Msg: r})
		}
		return replies, err
	}
	for k, v := range bot.msgRoutes {
		if HasPrefix(message, k) {
			if v.Handler != nil {
				reply := v.Handler(m)
				if reply.Msg != "" || reply.Keyboard != nil {
					replies = append(replies, reply)
				}
			} else {
				msg := v.SimpleHandler(m)
				if msg != "" {
					replies = append(replies, Reply{Msg: msg})
				}

			}
		}
	}
	return replies, nil
}

// RouteMessages routes inbound messages
func (bot *VKBot) RouteMessages(messages []*Message) (result map[*Message][]Reply) {
	result = make(map[*Message][]Reply)
	for _, m := range messages {
		if m.ReadState == 0 {
			if bot.IgnoreBots && m.UserID < 0 {
				continue
			}
			replies, err := bot.RouteMessage(m)
			if err != nil {
				sendError(m, err)
			}
			if len(replies) > 0 {
				result[m] = replies
			}
		}
	}
	return result
}

// MainRoute - main router func. Working cycle Listen.
func (bot *VKBot) MainRoute() {
	messages, err := bot.LongPoll.GetLongPollMessages()
	if err != nil {
		sendError(nil, err)
	}
	fmt.Println("inbox: ", messages)
	replies := bot.RouteMessages(messages)
	for m, msgs := range replies {
		for _, reply := range msgs {
			fmt.Println("outbox: ", reply.Msg)
			if reply.Msg != "" || reply.Keyboard != nil {
				_, err = bot.Reply(m, reply)
				if err != nil {
					log.Printf("Error sending message: '%+v'\n", reply)
					sendError(m, err)
					_, err = bot.Reply(m, Reply{Msg: "Cant send message, maybe wrong/china letters?"})
					if err != nil {
						sendError(m, err)
					}
				}
			}
		}
	}
}

// Reply - reply message
func (bot *VKBot) Reply(m *Message, reply Reply) (id int, err error) {
	if m.PeerID != 0 {
		return bot.API.SendAdvancedPeerMessage(m.PeerID, reply)
	}
	if m.ChatID != 0 {
		return bot.API.SendChatMessage(m.ChatID, reply.Msg)
	}
	return bot.API.SendMessage(m.UserID, reply.Msg)
}

// CheckFriends checking friend invites and matсhes and deletes mutual
func (bot *VKBot) CheckFriends() {
	uids, _ := bot.API.GetFriendRequests(false)
	if len(uids) > 0 {
		for _, uid := range uids {
			bot.API.AddFriend(uid)
			for k, v := range bot.actionRoutes {
				if k == "friend_add" {
					m := Message{Action: "friend_add", UserID: uid}
					v(&m)
				}
			}
		}
	}
	uids, _ = bot.API.GetFriendRequests(true)
	if len(uids) > 0 {
		for _, uid := range uids {
			bot.API.DeleteFriend(uid)
			for k, v := range bot.actionRoutes {
				if k == "friend_delete" {
					m := Message{Action: "friend_delete", UserID: uid}
					v(&m)
				}
			}
		}
	}
}

func (bot *VKBot) friendReceiver() {
	if bot.API.UID > 0 {
		bot.CheckFriends()
		c := time.Tick(30 * time.Second)
		for range c {
			bot.CheckFriends()
		}
	}
}
