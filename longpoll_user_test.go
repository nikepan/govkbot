package govkbot

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestGetMessage(t *testing.T) {

	data := "[4, 606838, 1, 329007844, 1508267602, \"тест\"]"
	d := json.NewDecoder(strings.NewReader(data))
	d.UseNumber()
	var lp interface{}
	if err := d.Decode(&lp); err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", lp)
	message := GetLongPollMessage(lp.([]interface{}))
	if message.Body != "тест" {
		t.Error("wrong longpoll message")
	}
}

func TestUserLongPollServer_ParseLongPollMessages(t *testing.T) {
	SetAPI("", "test", "")
	data := `{"ts":1668805076,"updates":[[4,2105994,561,123456,1496404246,"hello",{"title":" ... "},{"attach1_type":"photo","attach1":"123456_417336473","attach2_type":"audio","attach2":"123456_456239018"}]]}`
	server := NewUserLongPollServer(false, longPollVersion, 25)
	messages, err := server.ParseLongPollMessages(data)
	if err != nil {
		log.Fatal(err)
	}
	for _, msg := range messages.Messages {
		if msg.Body == "" {
			t.Error("empty message")
		}
	}
	if len(messages.Messages) != 1 {
		t.Error("wrong messages count")
	}
	if messages.Messages[0].Body != "hello" {
		t.Error("wrong messages text", messages.Messages[0].Body)
	}
}
