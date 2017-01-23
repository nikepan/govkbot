package govkbot

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestCall(t *testing.T) {
	api := API
	buf, err := api.Call("utils.getServerTime", url.Values{})
	if err != nil {
		t.Error("no response from VK")
	}

	m := SimpleResponse{}
	json.Unmarshal(buf, &m)
	if m.Error != nil {
		t.Error("Error response from VK")
	}
}
