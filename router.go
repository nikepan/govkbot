package govkbot

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// @TODO auth by user/password
// @TODO any substr handler
// @TODO long pooling
// @TODO replies from json files

// VKBot - bot config
type VKBot struct {
	msgRoutes        map[string]func(*Message) string
	actionRoutes     map[string]func(*Message) string
	cmdHandlers      map[string]func(*Message) string
	msgHandlers      map[string]func(*Message) string
	errorHandler     func(*Message, error)
	LastMsg          int
	markedMessages   map[int]*Message
	lastUserMessages map[int]int
	lastChatMessages map[int]int
	autoFriend       bool
}

var bot = newBot()

//API - bot API
var API = newAPI()

// SetDebug - enable/disable debug messages logging
func SetDebug(debug bool) {
	API.DEBUG = debug
}

func newBot() *VKBot {
	return &VKBot{
		msgRoutes:        make(map[string]func(*Message) string),
		actionRoutes:     make(map[string]func(*Message) string),
		markedMessages:   make(map[int]*Message),
		lastUserMessages: make(map[int]int),
		lastChatMessages: make(map[int]int),
	}
}

func newAPI() *VkAPI {
	return &VkAPI{
		Token:           "",
		URL:             "https://api.vk.com/method/",
		Ver:             "5.52",
		MessagesCount:   200,
		RequestInterval: 400, // 3 requests per second VK limit
		DEBUG:           false,
	}
}

// SetToken - set bot token
func SetToken(token string) {
	API.Token = token
}

// SetAutoFriend - enables mutual auto friending
func SetAutoFriend(af bool) {
	bot.autoFriend = af
}

// SetAPI - setup API config
func SetAPI(token string, url string, ver string) {
	SetToken(token)
	if url != "" {
		API.URL = url
	}
	if ver != "" {
		API.Ver = ver
	}
}

// HandleMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
// You can use m.Reply(string), if need more replies in handler
func HandleMessage(command string, handler func(*Message) string) {
	bot.msgRoutes[command] = handler
}

// HandleAction - add action handler.
// Function must return string to reply or "" (if no reply)
// You can use m.Reply(string), if need more replies in handler
func HandleAction(command string, handler func(*Message) string) {
	bot.actionRoutes[command] = handler
}

// HandleError - add error handler
func HandleError(handler func(*Message, error)) {
	bot.errorHandler = handler
}

func getMessages() ([]*Message, error) {
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

func sendError(msg *Message, err error) {
	if bot.errorHandler != nil {
		bot.errorHandler(msg, err)
	} else {
		log.Fatalf("VKBot error: %+v\n", err.Error())
	}

}

//RouteAction routes an action
func RouteAction(m *Message) (err error) {
	if m.Action != "" {
		debugPrint("route action: %+v\n", m.Action)
		for k, v := range bot.actionRoutes {
			if m.Action == k {
				bot.markedMessages[m.ID] = m
				msg := v(m)
				if msg != "" {
					_, err = m.Reply(msg)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// RouteMessage routes single message
func RouteMessage(m *Message) (err error) {
	message := strings.TrimSpace(strings.ToLower(m.Body))
	if strings.HasPrefix(message, "/ ") {
		message = "/" + strings.TrimPrefix(message, "/ ")
	}
	if m.Action != "" {
		err = RouteAction(m)
		if err != nil {
			return err
		}
	} else {
		marked := false
		for k, v := range bot.msgRoutes {
			if strings.HasPrefix(message, k) {
				msg := v(m)
				if msg != "" {
					_, err = m.Reply(msg)
					if err != nil {
						return err
					}
					marked = true
					_, ok := bot.markedMessages[m.ID]
					if ok {
						delete(bot.markedMessages, m.ID)
					}
				} else {
					if !marked {
						bot.markedMessages[m.ID] = m
					}
				}
			}
		}
	}
	return nil
}

// RouteMessages routes inbound messages
func RouteMessages(messages []*Message) {
	var err error
	for _, m := range messages {
		//if m.ID <= bot.LastMsg {
		//	break
		//}
		if m.ReadState == 0 {
			err = RouteMessage(m)
			if err != nil {
				sendError(m, err)
			}
		}
	}
}

// Listen - start server
func Listen(token string, url string, ver string, adminID int) {
	SetAPI(token, url, ver)
	API.AdminID = adminID
	API.UID = API.Me().ID

	go friendReceiver()

	c := time.Tick(3 * time.Second)
	for range c {
		bot.markedMessages = make(map[int]*Message)
		messages, err := getMessages()
		if err != nil {
			sendError(nil, err)
		}
		RouteMessages(messages)

		for _, m := range bot.markedMessages {
			m.MarkAsRead()
		}
	}
}

// CheckFriends checking friend invites and mathes and deletes mutual
func CheckFriends() {
	uids, _ := API.GetFriendRequests(false)
	if len(uids) > 0 {
		for _, uid := range uids {
			API.AddFriend(uid)
			for k, v := range bot.actionRoutes {
				if k == "friend_add" {
					m := Message{Action: "friend_add", UserID: uid}
					v(&m)
				}
			}
		}
	}
	uids, _ = API.GetFriendRequests(true)
	if len(uids) > 0 {
		for _, uid := range uids {
			API.DeleteFriend(uid)
			for k, v := range bot.actionRoutes {
				if k == "friend_delete" {
					m := Message{Action: "friend_delete", UserID: uid}
					v(&m)
				}
			}
		}
	}
}

func friendReceiver() {
	CheckFriends()
	c := time.Tick(30 * time.Second)
	for range c {
		CheckFriends()
	}
}

// NotifyAdmin - notify AdminID by VK
func NotifyAdmin(msg string) error {
	return API.NotifyAdmin(msg)
}
