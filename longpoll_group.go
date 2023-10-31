package govkbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// LongPollServer - longpoll server structure
type GroupLongPollServer struct {
	Key             string
	Server          string
	Ts              string
	Wait            int
	Mode            int
	Version         int
	RequestInterval int
	NeedPts         bool
	API             *VkAPI
	LpVersion       int
	ReadMessages    map[int64]time.Time
}

type GroupLongPollServerResponse struct {
	Response GroupLongPollServer
}

type GroupLongPollMessage struct {
	Type   string
	Object struct {
		Date                  int                    `json:"date"`
		FromID                int64                  `json:"from_id"`
		ID                    int64                  `json:"id"`
		Out                   int                    `json:"out"`
		PeerID                int                    `json:"peer_id"`
		Text                  string                 `json:"text"`
		ConversationMessageID int                    `json:"conversation_message_id"`
		FwdMessages           []GroupLongPollMessage `json:"fwd_messages"`
		Important             bool                   `json:"important"`
		RandomID              int64                  `json:"random_id"`
		//Attachments [] `json:"attachments"`
		IsHidden bool `json:"is_hidden"`
		Action   struct {
			Type     string
			MemberID int64 `json:"member_id"`
		}
	}
	GroupID int64
}

type GroupLongPollEvent struct {
	Type    string
	GroupID int64
}

type GroupFailResponse struct {
	Failed     int
	Date       int         `json:"date,omitempty"`
	Ts         interface{} `json:"ts,omitempty"`
	MinVersion int         `json:"min_version"`
	MaxVersion int         `json:"max_version"`
}

func (g *GroupFailResponse) TS() (string, error) {
	switch g.Ts.(type) {
	case string:
		return g.Ts.(string), nil
	case float64:
		return strconv.FormatFloat(g.Ts.(float64), 'f', 0, 64), nil
	default:
		return "", errors.New("unknown type")
	}
}

// NewLongPollServer - get longpoll server
func NewGroupLongPollServer(requestInterval int) (resp *GroupLongPollServer) {
	server := GroupLongPollServer{}
	server.Wait = DefaultWait
	server.Mode = DefaultMode
	server.Version = DefaultVersion
	server.RequestInterval = requestInterval
	server.ReadMessages = make(map[int64]time.Time)
	return &server
}

// Init - init longpoll server
func (server *GroupLongPollServer) Init() (err error) {
	r := GroupLongPollServerResponse{}
	err = API.CallMethod("groups.getLongPollServer", H{
		"group_id": strconv.FormatInt(API.GroupID, 10),
	}, &r)
	server.Wait = DefaultWait
	server.Mode = DefaultMode
	server.Version = DefaultVersion
	server.RequestInterval = API.RequestInterval
	server.Server = r.Response.Server
	server.Ts = r.Response.Ts
	server.Key = r.Response.Key
	return err
}

type GroupLongPollResponse struct {
	Ts       string
	Messages []*Message
}

// Request - make request to longpoll server
func (server *GroupLongPollServer) Request() ([]byte, error) {
	var err error

	if server.Server == "" {
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
	}

	parameters := url.Values{}
	parameters.Add("act", "a_check")
	parameters.Add("ts", server.Ts)
	parameters.Add("wait", strconv.Itoa(server.Wait))
	parameters.Add("key", server.Key)
	query := server.Server + "?" + parameters.Encode()
	if server.Server == "test" {
		content, err := ioutil.ReadFile("./mocks/longpoll.json")
		return content, err
	}
	resp, err := http.Get(query)
	if err != nil {
		debugPrint("%+v\n", err.Error())
		time.Sleep(time.Duration(time.Millisecond * time.Duration(server.RequestInterval)))
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(time.Millisecond * time.Duration(server.RequestInterval)))
	//debugPrint("longpoll vk resp: %+v\n", string(buf))

	failResp := GroupFailResponse{}
	err = json.Unmarshal(buf, &failResp)
	if err != nil {
		return nil, err
	}
	switch failResp.Failed {
	case 1:
		server.Ts, err = failResp.TS()
		if err != nil {
			log.Printf("error ts: %+v\n", err)
		}
		return server.Request()
	case 2:
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
		return server.Request()
	case 3:
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
		return server.Request()
	case 4:
		return nil, errors.New("vkapi: wrong longpoll version")
	default:
		server.Ts, err = failResp.TS()
		if err != nil {
			log.Printf("error ts: %+v\n", err)
		}
		return buf, nil
	}
}

// GetLongPollMessages - get messages via longpoll
func (server *GroupLongPollServer) GetLongPollMessages() ([]*Message, error) {
	resp, err := server.Request()
	if err != nil {
		return nil, err
	}
	messages, err := server.ParseLongPollMessages(string(resp))
	return messages.Messages, nil
}

func (server *GroupLongPollServer) ParseMessage(obj map[string]interface{}) (Message, error) {
	msg := Message{}
	msg.ID = getJSONInt64(obj["id"])
	if obj["text"] != nil {
		msg.Body = obj["text"].(string)
	} else {
		msg.Body = ""
		fmt.Printf("error parse message: %+v\n", obj)
		return msg, errors.New("error parse message")
	}
	userID := getJSONInt64(obj["from_id"])
	if userID != 0 {
		msg.UserID = userID
	}
	msg.PeerID = getJSONInt64(obj["peer_id"])
	if int64(msg.UserID) == msg.PeerID {
		msg.ChatID = 0
	} else {
		msg.ChatID = msg.PeerID
	}
	msg.Date = getJSONInt(obj["date"])
	fwd, ok := obj["fwd_messages"]
	if ok {
		for _, m := range fwd.([]interface{}) {
			fwdMsg, err := server.ParseMessage(m.(map[string]interface{}))
			if err != nil {
				fmt.Printf("error parse fwd message: %+v\n", err)
			} else {
				msg.FwdMessages = append(msg.FwdMessages, fwdMsg)
			}
		}
	}
	return msg, nil
}

// ParseLongPollMessages - parse longpoll messages
func (server *GroupLongPollServer) ParseLongPollMessages(j string) (*GroupLongPollResponse, error) {
	//fmt.Printf("\n>>>>>>>>>>>>>updates0: %+v\n\n", j)
	d := json.NewDecoder(strings.NewReader(j))
	d.UseNumber()
	var lp interface{}
	if err := d.Decode(&lp); err != nil {
		return nil, err
	}
	lpMap := lp.(map[string]interface{})
	result := GroupLongPollResponse{Messages: []*Message{}}
	result.Ts = lpMap["ts"].(string)
	updates := lpMap["updates"].([]interface{})
	//fmt.Printf("\n>>>>>>>>>>>>>updates: %+v\n\n", updates)
	for _, event := range updates {
		//el := event.(interface{})
		eventType := event.(map[string]interface{})["type"].(string)
		if eventType == "message_new" {
			obj := event.(map[string]interface{})["object"].(map[string]interface{})
			out := getJSONInt(obj["out"])
			if out == 0 {
				debugPrint("new message: %+v\n", obj)
				msg, err := server.ParseMessage(obj)
				result.Messages = append(result.Messages, &msg)
				if err != nil {
					fmt.Printf("error parse message: %+v \n>>>%+v \n>>>>> %+v\n", err, obj, msg)
				}
			}
		}
	}
	if len(result.Messages) == 0 {
		debugPrint(j)
	}
	if len(result.Messages) > 0 {
		debugPrint("new messages: ts: %+v = %+v\n", result.Ts, len(result.Messages))
	}
	// result.Messages = server.FilterReadMesages(result.Messages)
	// fmt.Printf("\n>>>>>>>>>>>>>messages2: %+v\n\n", result)
	return &result, nil
}

// FilterReadMesages - filter read messages
func (server *GroupLongPollServer) FilterReadMesages(messages []*Message) (result []*Message) {
	for _, m := range messages {
		t, ok := server.ReadMessages[m.ID]
		if ok {
			if time.Since(t).Minutes() > 1 {
				delete(server.ReadMessages, m.ID)
			}
		} else {
			result = append(result, m)
			server.ReadMessages[m.ID] = time.Now()
		}
	}
	return result
}
