package govkbot

import (
	"testing"
	"encoding/json"
	"log"
	"strings"
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
	if message.Text != "тест" {
		t.Error("wrong longpoll message")
	}
}

func TestLongPollServer_ParseLongPollMessages(t *testing.T) {
	SetAPI("", "test", "")
	data := `{"ts":1668805076,"updates":[[4,2105994,561,123456,1496404246,"hello",{"title":" ... "},{"attach1_type":"photo","attach1":"123456_417336473","attach2_type":"audio","attach2":"123456_456239018"}]]}`
	server := API.GetLongPollServer(false, longPollVersion)
	messages, err := server.ParseLongPollMessages(data)
	if err != nil {
		log.Fatal(err)
	}
	for _, msg := range messages {
		if msg.Body == "" {
			t.Error("empty message")
		}
	}
	if len(messages) != 1 {
		t.Error("wrong messages count")
	}
	if messages[0].Body != "hello" {
		t.Error("wrong messages text", messages[0].Body)
	}
}