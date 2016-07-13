package govkbot

import (
	"log"
	"strings"
	"time"
)

// @TODO any substr handler
// @TODO long pooling
// @TODO replies from json files

// VKBot - bot config
type VKBot struct {
	msgRoutes    map[string]func(*Message) string
	actionRoutes map[string]func(*Message) string
	cmdHandlers  map[string]func(*Message) string
	msgHandlers  map[string]func(*Message) string
}

var bot = newBot()

//API - bot API
var API = newAPI()

func newBot() *VKBot {
	return &VKBot{
		msgRoutes:    make(map[string]func(*Message) string),
		actionRoutes: make(map[string]func(*Message) string)}
}

func newAPI() *VkAPI {
	return &VkAPI{Token: "", Url: "https://API.vk.com/method/", Ver: "5.52"}
}

// SetToken - set bot token
func SetToken(token string) {
	API.Token = token
}

// SetAPI - setup API config
func SetAPI(token string, url string, ver string) {
	API.Token = token
	if url != "" {
		API.Url = url
	}
	if ver != "" {
		API.Ver = ver
	}
}

// HandleMessage - add substr message handler
func HandleMessage(command string, handler func(*Message) string) {
	bot.msgRoutes[command] = handler
}

// HandleAction - add action handler
func HandleAction(command string, handler func(*Message) string) {
	bot.actionRoutes[command] = handler
}

// Listen - start server
func Listen(token string, url string, ver string) {
	SetAPI(token, url, ver)
	API.Uid = API.Me().ID

	go friendReceiver()

	c := time.Tick(3 * time.Second)
	for _ = range c {
		messages := API.GetMessages()
		for _, m := range messages.Items {
			if m.ReadState == 0 {
				message := strings.ToLower(m.Body)
				go m.MarkAsRead()
				if strings.HasPrefix(message, "/ ") {
					message = "/" + strings.TrimPrefix(message, "/ ")
				}
				if m.Action != "" {
					log.Printf(m.Action)
					for k, v := range bot.actionRoutes {
						if m.Action == k {
							log.Printf("success")
							msg := v(&m)
							if msg != "" {
								m.Reply(msg)
							}
						}
					}
				} else {
					for k, v := range bot.msgRoutes {
						if strings.HasPrefix(message, k) {
							msg := v(&m)
							if msg != "" {
								m.Reply(msg)
							}
						}
					}
				}
			}
		}
	}
}

func checkFriends() {
	uids := API.GetFriendRequests(false)
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
	uids = API.GetFriendRequests(true)
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
	checkFriends()
	c := time.Tick(30 * time.Second)
	for _ = range c {
		checkFriends()
	}
}
