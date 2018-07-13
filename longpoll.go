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
const vkTsDiff = 180116722

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
}

type LongPollServerResponse struct {
	Response LongPollServer
}

type LongPollUpdate []interface{}
type LongPollUpdateNum []int64

type LongPollResponse struct {
	Ts uint
	Updates []interface{}
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

type HistoryMessage struct {
	Date int `json:"date"`
	FromID int `json:"from_id"`
	ID int `json:"id"`
	Out int `json:"out"`
	PeerID int64 `json:"peer_id"`
	Text string `json:"text"`
	ConversationMessageId int `json:"converstion_message_id"`
	FwdMessages []LongPollMessage `json:"fwd_messages"`
	Important bool `json:"important"`
	RandomID int `json:"random_id"`
	Attachments []Attachment `json:"attachments"`
	IsHidden bool `json:"is_hidden"`
}

type HistoryResponse struct {
	Messages []HistoryMessage
}

type HistoryReader struct {
	ts int64
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

func VKTimestamp(ts int64) int64 {
	return ts + vkTsDiff
}

func GetMessagesHistory(ts int64) ([]HistoryMessage, error) {
	params := make(map[string]string)
	params["ts"] = strconv.FormatInt(ts, 10)
	resp := HistoryResponse{}
	err := API.CallMethod("messages.getLongPollHistory", params, &resp)
	return resp.Messages, err
}

func (h *HistoryReader) GetMessages() ([]*Message, error) {
	ts := VKTimestamp(time.Now().Unix())
	if h.ts == 0 {
		h.ts = VKTimestamp(time.Now().Unix())
	}
	messages, err := GetMessagesHistory(h.ts)
	h.ts = ts
	return ConvertMessages(messages), err
}

func ConvertMessages(messages []HistoryMessage) []*Message {
	result := []*Message{}
	for _, message := range messages {
		msg := Message{}
		msg.Body = message.Text
		msg.UserID = message.FromID
		msg.ChatID = int(message.PeerID)
		msg.ID = message.ID
		result = append(result, &msg)
	}
	return result
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
	parameters.Add("mode", "2")
	parameters.Add("version", strconv.Itoa(server.Version))
	query := "https://"+server.Server+"?"+parameters.Encode()
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


func GetLongPollMessage(resp []interface{}) *LongPollMessage {
	message := LongPollMessage{}
	mt, _ := resp[0].(json.Number).Int64()
	message.MessageType = int(mt)
	mid, _ := resp[1].(json.Number).Int64()
	message.MessageID = int(mid)
	flags, _ := resp[2].(json.Number).Int64()
	message.Flags = int(flags)
	message.PeerID, _ = resp[3].(json.Number).Int64()
	message.Timestamp, _ = resp[4].(json.Number).Int64()
	message.Text =resp[5].(string)
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

func (server *LongPollServer) ParseLongPollMessages(j string) ([]*Message, error) {
	//fmt.Println(j)
	count := gjson.Get(j, "updates.#")
	result := []*Message{}
	for i := 0; i < int(count.Int()); i++ {
		eventType := gjson.Get(j, "updates."+strconv.Itoa(i)+".0")
		if eventType.Int() == 4 {
			out := gjson.Get(j, "updates."+strconv.Itoa(i)+".2").Int() & 2
			if out == 0 {
				msg := Message{}
				msg.ID = int(gjson.Get(j, "updates."+strconv.Itoa(i)+".1").Int())
				msg.Body = gjson.Get(j, "updates."+strconv.Itoa(i)+".5").String()
				msg.UserID = int(gjson.Get(j, "updates."+strconv.Itoa(i)+".6.from").Int())
				msg.PeerID = int(gjson.Get(j, "updates."+strconv.Itoa(i)+".3").Int())
				if msg.UserID == 0 {
					msg.UserID = msg.PeerID
				} else {
					msg.ChatID = msg.PeerID - 2000000000
				}
				msg.Date = int(gjson.Get(j, "updates."+strconv.Itoa(i)+".4").Int())
				result = append(result, &msg)
				if msg.UserID == 3759927 {
					fmt.Println(msg)
					fmt.Println(j)
				}
			}
		}
	}
	return result, nil
}