package main

import (
	"fmt"
	"log"

	"github.com/nikepan/govkbot/v2"
)

type configuration struct {
	VKToken      string
	ServiceURL   string
	ServiceToken string
	AdminID      int64
}

func getMeMessage(uid int64) (reply string) {
	me, _ := govkbot.API.Me()
	return fmt.Sprintf("You: %+v %+v", me.FirstName, me.LastName)
}

func anyHandler(m *govkbot.Message) (reply string) {
	notifyAdmin(fmt.Sprintf("Command %+v by user vk.com/id%+v in chat %+v", m.Body, m.UserID, m.PeerID))
	return reply
}

func meHandler(m *govkbot.Message) (reply string) {
	return getMeMessage(m.UserID)
}

func helpHandler(m *govkbot.Message) (reply govkbot.Reply) {
	keyboard := govkbot.Keyboard{Buttons: make([][]govkbot.Button, 0)}
	button := govkbot.NewButton("/me", nil)
	row := make([]govkbot.Button, 0)
	row = append(row, button)
	keyboard.Buttons = append(keyboard.Buttons, row)

	return govkbot.Reply{Msg: availableCommands, Keyboard: &keyboard}
}

func errorHandler(msg *govkbot.Message, err error) {
	if _, ok := err.(*govkbot.VKError); !ok {
		notifyAdmin("VK ERROR: " + err.Error()) // err.(govkbot.VKError).ErrorCode
	}
	notifyAdmin("ERROR: " + err.Error())
}

func inviteHandler(m *govkbot.Message) (reply string) {
	log.Printf("invite: %+v %+v %+v\n", m.ActionMID, govkbot.API.UID, m.ActionMID == govkbot.API.UID)
	if m.ActionMID == govkbot.API.UID {
		go m.MarkAsRead()
		notifyAdmin(fmt.Sprintf("I'm invited to chat %+v )", m.Title))
		reply = replyGreet()
	} else {
		log.Printf("greet: %+v %+v\n", m.ActionMID, m)
		reply = greetUser(m.ActionMID)
	}
	return reply
}

func kickHandler(m *govkbot.Message) (reply string) {
	if m.ActionMID == govkbot.API.UID {
		go m.MarkAsRead()
		notifyAdmin(fmt.Sprintf("I'm kicked from chat %+v (", m.Title))
	}
	return reply
}

func greetUser(uid int64) (reply string) {
	u, err := govkbot.API.User(uid)
	if err == nil {
		reply = fmt.Sprintf("Hello, %+v", u.FullName())
	}
	return reply
}

func replyGreet() (reply string) {
	reply = "Hi all. I'am bot\n" + availableCommands
	return reply
}

func addFriendHandler(m *govkbot.Message) (reply string) {
	log.Printf("friend %+v added\n", m.UserID)
	notifyAdmin(fmt.Sprintf("user vk.com/id%+v add me to friends", m.UserID))
	return reply
}

func deleteFriendHandler(m *govkbot.Message) (reply string) {
	log.Printf("friend %+v deleted\n", m.UserID)
	notifyAdmin(fmt.Sprintf("user vk.com/id%+v delete me from friends", m.UserID))
	return reply
}

func notifyAdmin(msg string) {
	err := govkbot.NotifyAdmin(msg)
	if err != nil {
		log.Printf("VK Admin Notify ERROR: %+v\n", msg)
	}
}
