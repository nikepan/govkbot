package govkbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

const DefaultWait = 25
const (
	LongPollModeGetAttachments    = 2
	LongPollModeGetExtendedEvents = 8
	LongPollModeGetPts            = 32
	LongPollModeGetExtraData      = 64
	LongPollModeGetRandomID       = 128
)
const DefaultMode = LongPollModeGetAttachments
const DefaultVersion = 2
const ChatPrefix = 2000000000

type LongPollServer interface {
	Init() (err error)
	Request() ([]byte, error)
	GetLongPollMessages() ([]*Message, error)
	FilterReadMesages(messages []*Message) (result []*Message)
}

// LongPollServer - longpoll server structure
type UserLongPollServer struct {
	Key             string
	Server          string
	Ts              int
	Wait            int
	Mode            int
	Version         int
	RequestInterval int
	NeedPts         bool
	API             *VkAPI
	LpVersion       int
	ReadMessages    map[int]time.Time
}

// LongPollServerResponse - response format for longpoll info
type UserLongPollServerResponse struct {
	Response UserLongPollServer
}

type LongPollUpdate []interface{}
type LongPollUpdateNum []int64

type LongPollResponse struct {
	Ts       uint
	Messages []*Message
}

type Attachment struct {
	AttachType      string
	Attach          string
	Fwd             string
	From            int
	Geo             int
	GeoProvider     int
	Title           string
	AttachProductID int
	AttachPhoto     string
	AttachTitle     string
	AttachDesc      string
	AttachURL       string
	Emoji           bool
	FromAdmin       int
	SourceAct       string
	SourceMid       int
}

type LongPollMessage struct {
	MessageType int
	MessageID   int
	Flags       int
	PeerID      int64
	Timestamp   int64
	Text        string
	Attachments []Attachment
	RandomID    int
}

type FailResponse struct {
	Failed     int
	Ts         int
	MinVersion int `json:"min_version"`
	MaxVersion int `json:"max_version"`
}

// NewLongPollServer - get longpoll server
func NewUserLongPollServer(needPts bool, lpVersion int, requestInterval int) (resp *UserLongPollServer) {
	server := UserLongPollServer{}
	server.NeedPts = needPts
	server.Wait = DefaultWait
	server.Mode = DefaultMode
	server.Version = DefaultVersion
	server.RequestInterval = requestInterval
	server.LpVersion = lpVersion
	server.ReadMessages = make(map[int]time.Time)
	return &server
}

// Init - init longpoll server
func (server *UserLongPollServer) Init() (err error) {
	r := UserLongPollServerResponse{}
	pts := 0
	if server.NeedPts {
		pts = 1
	}
	err = API.CallMethod("messages.getLongPollServer", H{
		"need_pts": strconv.Itoa(pts),
		"message":  strconv.Itoa(server.LpVersion),
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

// Request - make request to longpoll server
func (server *UserLongPollServer) Request() ([]byte, error) {
	var err error

	if server.Server == "" {
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
	}

	parameters := url.Values{}
	parameters.Add("act", "a_check")
	parameters.Add("ts", strconv.Itoa(server.Ts))
	parameters.Add("wait", strconv.Itoa(server.Wait))
	parameters.Add("key", server.Key)
	parameters.Add("mode", strconv.Itoa(DefaultMode))
	parameters.Add("version", strconv.Itoa(server.Version))
	query := "https://" + server.Server + "?" + parameters.Encode()
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

	failResp := FailResponse{}
	err = json.Unmarshal(buf, &failResp)
	if err != nil {
		log.Printf("longpoll vk resp: %+v\n", string(buf))
		return nil, err
	}
	switch failResp.Failed {
	case 1:
		server.Ts = failResp.Ts
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
		return buf, nil
	}
}

// GetLongPollMessage - get message from longpoll json row
func GetLongPollMessage(resp []interface{}) *Message {
	message := Message{}
	mid, _ := resp[1].(json.Number).Int64()
	message.ID = int(mid)
	flags, _ := resp[2].(json.Number).Int64()
	message.Flags = int(flags)
	message.PeerID, _ = resp[3].(json.Number).Int64()
	message.Timestamp, _ = resp[4].(json.Number).Int64()
	message.Body = resp[5].(string)
	return &message
}

// GetLongPollMessages - get messages via longpoll
func (server *UserLongPollServer) GetLongPollMessages() ([]*Message, error) {
	resp, err := server.Request()
	if err != nil {
		return nil, err
	}
	messages, err := server.ParseLongPollMessages(string(resp))
	return messages.Messages, nil
}

func getJSONInt64(el interface{}) int64 {
	if el == nil {
		return 0
	}
	v, _ := el.(json.Number).Int64()
	return v
}

func getJSONInt(el interface{}) int {
	return int(getJSONInt64(el))
}

// ParseLongPollMessages - parse longpoll messages
func (server *UserLongPollServer) ParseLongPollMessages(j string) (*LongPollResponse, error) {
	d := json.NewDecoder(strings.NewReader(j))
	d.UseNumber()
	var lp interface{}
	if err := d.Decode(&lp); err != nil {
		return nil, err
	}
	lpMap := lp.(map[string]interface{})
	result := LongPollResponse{Messages: []*Message{}}
	ts, _ := lpMap["ts"].(json.Number).Int64()
	result.Ts = uint(ts)
	updates := lpMap["updates"].([]interface{})
	for _, event := range updates {
		el := event.([]interface{})
		eventType := getJSONInt(el[0])
		if eventType == 4 {
			out := getJSONInt(el[2]) & 2
			if out == 0 {
				msg := Message{}
				msg.ID = getJSONInt(el[1])
				msg.Body = el[5].(string)
				userID := el[6].(map[string]interface{})["from"]
				if userID != nil {
					msg.UserID, _ = strconv.Atoi(userID.(string))
				}
				msg.PeerID = getJSONInt64(el[3])
				if msg.UserID == 0 {
					msg.UserID = int(msg.PeerID)
				} else {
					msg.ChatID = int(msg.PeerID - ChatPrefix)
				}
				msg.Date = getJSONInt(el[4])
				fmt.Println(msg.Body)
				result.Messages = append(result.Messages, &msg)
			}
		}
	}
	if len(result.Messages) == 0 {
		fmt.Println(j)
	}
	result.Messages = server.FilterReadMesages(result.Messages)
	return &result, nil
}

// FilterReadMesages - filter read messages
func (server *UserLongPollServer) FilterReadMesages(messages []*Message) (result []*Message) {
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
