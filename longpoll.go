package govkbot

import (
	"strconv"
	"net/url"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/tidwall/gjson"
	"github.com/labstack/gommon/log"
	"fmt"
	"errors"
)

const DefaultWait = 25
const (
	LongPollModeGetAttachments = 2
	LongPollModeGetExtendedEvents = 8
	LongPollModeGetPts = 32
	LongPollModeGetExtraData = 64
	LongPollModeGetRandomID = 128
)
const DefaultMode = LongPollModeGetAttachments
const DefaultVersion = 2
const ChatPrefix = 2000000000

type LongPollServer struct {
	Key    string
	Server string
	Ts     int
	Wait int
	Mode int
	Version int
	RequestInterval int
	NeedPts bool
	Api *VkAPI
	LpVersion int
	ReadMessages map[int]time.Time
}

type LongPollServerResponse struct {
	Response LongPollServer
}

type LongPollUpdate []interface{}
type LongPollUpdateNum []int64

type LongPollResponse struct {
	Ts uint
	Messages []*Message
}

type Attachment struct {
	AttachType string
	Attach string
	Fwd string
	From int
	Geo int
	GeoProvider int
	Title string
	AttachProductID int
	AttachPhoto string
	AttachTitle string
	AttachDesc string
	AttachURL string
	Emoji bool
	FromAdmin int
	SourceAct string
	SourceMid int
}

type LongPollMessage struct {
	MessageType int
	MessageID int
	Flags int
	PeerID int64
	Timestamp int64
	Text string
	Attachments []Attachment
	RandomID int
}

type FailResponse struct {
	Failed int
	Ts     int
	MinVersion int `json:"min_version"`
	MaxVersion int `json:"max_version"`
}


func (api *VkAPI) GetLongPollServer(needPts bool, lpVersion int) (resp *LongPollServer) {
	server := LongPollServer{}
	server.NeedPts = needPts
	server.Wait = DefaultWait
	server.Mode = DefaultMode
	server.Version = DefaultVersion
	server.RequestInterval = api.RequestInterval
	server.LpVersion =lpVersion
	server.ReadMessages = make(map[int]time.Time)
	return &server
}

func (server *LongPollServer) Init() (err error) {
	r := LongPollServerResponse{}
	pts := 0
	if server.NeedPts {
		pts = 1
	}
	err = API.CallMethod("messages.getLongPollServer", H{
		"need_pts": strconv.Itoa(pts),
		"message": strconv.Itoa(server.LpVersion),
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

func (server *LongPollServer) Request() ([]byte, error) {
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
	query := "https://"+server.Server+"?"+parameters.Encode()
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
	debugPrint("longpoll vk resp: %+v\n", string(buf))

	failResp := FailResponse{}
	err = json.Unmarshal(buf, &failResp)
	if err != nil {
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
	return buf, nil
}

func GetLongPollResponse(buf []byte) (*LongPollResponse, error) {
	d := json.NewDecoder(strings.NewReader(string(buf)))
	d.UseNumber()
	var lp interface{}
	if err := d.Decode(&lp); err != nil {
		return nil, err
	}
	lpMap := lp.(map[string]interface{})
	result := LongPollResponse{}
	ts, _ := lpMap["ts"].(json.Number).Int64()
	result.Ts = uint(ts)
	result.Updates = lpMap["updates"].([]interface{})
	return &result, nil
}


func GetLongPollMessage(resp []interface{}) *Message {
	message := Message{}
	mid, _ := resp[1].(json.Number).Int64()
	message.ID = int(mid)
	flags, _ := resp[2].(json.Number).Int64()
	message.Flags = int(flags)
	message.PeerID, _ = resp[3].(json.Number).Int64()
	message.Timestamp, _ = resp[4].(json.Number).Int64()
	message.Body =resp[5].(string)
	return &message
}

func (server *LongPollServer) GetLongPollMessages() ([]*Message, error) {
	resp, err := server.Request()
	if err != nil {
		return nil, err
	}
	messages, err := server.ParseLongPollMessages(string(resp))
	return messages, nil
}

func (server *LongPollServer) ParseLongPollMessages(j string) (*LongPollResponse, error) {
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
		eventType := el[0].(int)
		if eventType == 4 {
			out := el[2].(int) & 2
			if out == 0 {
				msg := Message{}
				msg.ID = el[1].(int)
				msg.Body = el[5].(string)
				msg.UserID = el[6].(map[string]interface{})["from"].(int)
				msg.PeerID = el[3].(int64)
				if msg.UserID == 0 {
					msg.UserID = int(msg.PeerID)
				} else {
					msg.ChatID = int(msg.PeerID - ChatPrefix)
				}
				msg.Date = int(gjson.Get(j, "updates."+strconv.Itoa(i)+".4").Int())
				result.Messages = append(result.Messages, &msg)
				if msg.UserID == 3759927 {
					fmt.Println(msg)
					fmt.Println(j)
				}
			}
		}
	}
	if len(result.Messages) == 0 {
		fmt.Println(j)
	}
	result.Messages = server.FilterReadMesages(result.Messages)
	return &result, nil
}

func (server *LongPollServer) FilterReadMesages(messages []*Message) (result []*Message) {
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