package govkbot

import (
	"encoding/json"
	"github.com/nikepan/govkbot"
	"net/url"
	"testing"
)

func TestCall(t *testing.T) {
	api := govkbot.API
	buf, err := api.Call("utils.getServerTime", url.Values{})
	if err != nil {
		t.Error("no response from VK")
	}

	m := govkbot.SimpleResponse{}
	json.Unmarshal(buf, &m)
	if m.Error != nil {
		t.Error("Error response from VK")
	}
}
