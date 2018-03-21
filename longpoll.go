package govkbot

import (
	"strconv"
	"net/url"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"strings"
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

type LongPollServer struct {
	Key    string
	Server string
	Ts     int
	Wait int
	Mode int
	Version int
	RequestInterval int
}

type LongPollServerResponse struct {
	Response LongPollServer
}

type LongPollUpdate []interface{}
type LongPollUpdateNum []int64

type LongPollResponse struct {
	Ts uint
	Updates []LongPollUpdate
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


func (api *VkAPI) GetLongPollServer(needPts bool, lpVersion int) (resp LongPollServer, err error) {
	r := LongPollServerResponse{}
	pts := 0
	if needPts {
		pts = 1
	}
	err = api.CallMethod(apiMessagesSend, H{
		"need_pts": strconv.Itoa(pts),
		"message": strconv.Itoa(lpVersion),
	}, &r)
	r.Response.Wait = DefaultWait
	r.Response.Mode = DefaultMode
	r.Response.Version = DefaultVersion
	r.Response.RequestInterval = api.RequestInterval
	return r.Response, err
}


func (server *LongPollServer) Request() (*LongPollResponse, error) {
	parameters := url.Values{}
	parameters.Add("act", "a_check")
	parameters.Add("ts", strconv.Itoa(server.Ts))
	parameters.Add("wait", strconv.Itoa(server.Wait))
	parameters.Add("version", strconv.Itoa(server.Version))
	query := "https://"+server.Server+parameters.Encode()
	resp, err := http.Get(query)
	if err != nil {
		debugPrint("%+v\n", err.Error())
		time.Sleep(time.Duration(time.Millisecond * time.Duration(server.RequestInterval)))
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(time.Millisecond * time.Duration(server.RequestInterval)))
	debugPrint("vk resp: %+v\n", string(buf))
	return GetLongPollResponse(buf)
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

func (server *LongPollServer) GetLongPollMessages() ([]*LongPollMessage, error) {
	resp, err := server.Request()
	if err != nil {
		return nil, err
	}
	messages := make([]*LongPollMessage, 1)
	for _, el := range resp.Updates {
		if el[0].(int) == 4 {
			messages = append(messages, GetLongPollMessage(el))
		}
	}
	return messages, nil
}