package govkbot

import (
	"log"
)

const (
	vkAPIURL        = "https://api.vk.com/method/"
	vkAPIVer        = "5.131"
	messagesCount   = 200
	requestInterval = 400 // 3 requests per second VK limit
	longPollVersion = 3
)

// API - bot API
var API = newAPI()

var Bot = API.NewBot()

// SetDebug - enable/disable debug messages logging
func SetDebug(debug bool) {
	API.DEBUG = debug
}

func newAPI() *VkAPI {
	return &VkAPI{
		Token:           "",
		URL:             vkAPIURL,
		Ver:             vkAPIVer,
		MessagesCount:   messagesCount,
		RequestInterval: requestInterval,
		DEBUG:           false,
		HTTPS:           true,
	}
}

// SetToken - set bot token
func SetToken(token string) {
	API.Token = token
}

// SetAutoFriend - enables mutual auto friending
func SetAutoFriend(af bool) {
	Bot.SetAutoFriend(af)
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

// SetLang - sets VK response language. Default auto. Available: en, ru, ua, be, es, fi, de, it
func SetLang(lang string) {
	API.Lang = lang
}

// HandleMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
func HandleMessage(command string, handler func(*Message) string) {
	Bot.HandleMessage(command, handler)
}

// HandleAdvancedMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
func HandleAdvancedMessage(command string, handler func(*Message) Reply) {
	Bot.HandleAdvancedMessage(command, handler)
}

// HandleAction - add action handler.
// Function must return string to reply or "" (if no reply)
func HandleAction(command string, handler func(*Message) string) {
	Bot.HandleAction(command, handler)
}

// HandleError - add error handler
func HandleError(handler func(*Message, error)) {
	Bot.HandleError(handler)
}

func sendError(msg *Message, err error) {
	if Bot.errorHandler != nil {
		Bot.errorHandler(msg, err)
	} else {
		log.Fatalf("VKBot error: %+v\n", err.Error())
	}

}

// Listen - start server
func Listen(token string, url string, ver string, adminID int64) error {
	if API.Token == "" {
		SetAPI(token, url, ver)
	}
	API.AdminID = adminID
	if Bot.API.IsGroup() {
		return Bot.ListenGroup(API)
	}
	return Bot.ListenUser(API)
}

// NotifyAdmin - notify AdminID by VK
func NotifyAdmin(msg string) error {
	return API.NotifyAdmin(msg)
}
